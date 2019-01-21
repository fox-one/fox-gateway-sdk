package gateway

import (
	"context"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
)

type Member struct {
	*Client

	key    string
	secret string
}

func NewMember(key, secret, apiBase string) *Member {
	return &Member{
		key:    key,
		secret: secret,
		Client: NewClient(apiBase),
	}
}

// auth

type memberAuth struct {
	*Member

	pin    string
	nonce  string
	expire time.Duration
}

func (m *memberAuth) GenerateToken(method, uri string, body []byte) string {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(m.expire).Unix(),
		"sign": signRequest(method, uri, body),
		"key":  m.key,
	}

	if len(m.pin) > 0 {
		payload := map[string]interface{}{
			"p": m.pin,
			"n": m.nonce,
		}

		data, _ := jsoniter.Marshal(payload)

		aeskey := MD5(m.key)
		aesiv := []byte(m.secret)
		pinToken, err := Encrypt(data, aeskey, aesiv)
		if err != nil {
			log.Panicln("encrypt pin", err)
		}

		claims["pin"] = pinToken
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(m.secret))
	if err != nil {
		log.Panicln("sign merchent token", err)
	}

	return t
}

func (m *Member) PresignWithPin(pin, nonce string, expire time.Duration) Authenticator {
	return &memberAuth{
		Member: m,
		pin:    pin,
		nonce:  nonce,
		expire: expire,
	}
}

func (m *Member) Presign(expire time.Duration) Authenticator {
	return m.PresignWithPin("", "", expire)
}

func (m *Member) MemberInfo(ctx context.Context) (*MemberView, error) {
	data, err := m.GET("/member/info").Auth(m.Presign(time.Minute)).Do(ctx)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Err
		Member *MemberView `json:"member"`
	}
	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.Member, nil
}

func (m *Member) Validate(ctx context.Context, method, uri, body, token string) (string, error) {
	data, err := m.POST("/member/validate").
		P("method", method).
		P("uri", uri).
		P("body", body).
		P("token", token).
		Do(ctx)
	if err != nil {
		return "", err
	}

	var resp struct {
		Err
		MemberID string `json:"member_id"`
	}
	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return "", err
	}

	if resp.Code > 0 {
		return "", resp.Err
	}

	return resp.MemberID, nil
}
