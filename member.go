package gateway

import (
	"context"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fox-one/httpclient"
	jsoniter "github.com/json-iterator/go"
)

type MemberClient struct {
	*Client

	key    string
	secret string
}

func (c *Client) Member() *MemberClient {
	return &MemberClient{
		Client: c.Group("member"),
	}
}

func NewMemberClient(apiBase string) *MemberClient {
	return NewClient(apiBase).Member()
}

func (c *MemberClient) WithSession(key, secret string) *MemberClient {
	return &MemberClient{
		Client: c.Client,
		key:    key,
		secret: secret,
	}
}

// auth

type memberAuth struct {
	*MemberClient

	pin    string
	nonce  string
	expire time.Duration
}

func (m *memberAuth) PrepareAuth(req *httpclient.Request) {
	if m.nonce != "" {
		req.P(nonceKey, m.nonce)
	}
}

func (m *memberAuth) token(method, uri string, body []byte) string {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(m.expire).Unix(),
		"sign": signRequest(method, uri, string(body)),
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

func (m *memberAuth) Auth(req *httpclient.Request, method, uri string, body []byte) {
	req.AddToken(m.token(method, uri, body))
}

func (m *MemberClient) PresignWithPin(pin string, expire time.Duration) *memberAuth {
	return &memberAuth{
		MemberClient: m,
		pin:          pin,
		nonce:        newNonce(),
		expire:       expire,
	}
}

func (m *MemberClient) Presign(expire time.Duration) *memberAuth {
	return &memberAuth{MemberClient: m, expire: expire}
}

func (m *MemberClient) Sign(method, uri string, body []byte, expire time.Duration) string {
	return m.Presign(expire).token(method, uri, body)
}

func (m *MemberClient) SignWithPin(pin, method, uri string, body []byte, expire time.Duration) string {
	return m.PresignWithPin(pin, expire).token(method, uri, body)
}

func (m *MemberClient) MemberInfo(ctx context.Context) (*MemberView, error) {
	data, err := m.GET("/info").Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
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

func (m *MemberClient) Validate(ctx context.Context, method, uri, body, token string) (string, error) {
	data, err := m.POST("/validate").
		P("method", method).
		P("uri", uri).
		P("body", body).
		P("token", token).
		Do(ctx).Bytes()
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
