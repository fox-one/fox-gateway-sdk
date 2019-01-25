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

func (c *Client) Member(key, secret string) *MemberClient {
	return &MemberClient{
		key:    key,
		secret: secret,
		Client: c.Group("member"),
	}
}

func NewMemberClient(key, secret, apiBase string) *MemberClient {
	return NewClient(apiBase).Member(key, secret)
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

func (m *memberAuth) Auth(req *httpclient.Request, method, uri string, body []byte) {
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

	req.AddToken(t)
}

func (m *MemberClient) PresignWithPin(pin, nonce string, expire time.Duration) httpclient.Authenticator {
	return &memberAuth{
		MemberClient: m,
		pin:          pin,
		nonce:        nonce,
		expire:       expire,
	}
}

func (m *MemberClient) Presign(expire time.Duration) httpclient.Authenticator {
	return &memberAuth{MemberClient: m, expire: expire}
}

func (m *MemberClient) MemberInfo(ctx context.Context) (*MemberView, error) {
	data, err := m.client.GET("/info").Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
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
	data, err := m.client.POST("/validate").
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
