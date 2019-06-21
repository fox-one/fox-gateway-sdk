package gateway

import (
	"context"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fox-one/httpclient"
	jsoniter "github.com/json-iterator/go"
)

type MemberClient struct {
	*Client

	MemberAuth
}

func (c *Client) Member() *MemberClient {
	return &MemberClient{
		Client: c.Group("member"),
	}
}

func NewMemberClient(apiBase string) *MemberClient {
	return NewClient(apiBase).Member()
}

func (c *MemberClient) WithAuth(auth MemberAuth) *MemberClient {
	return &MemberClient{
		Client:     c.Client,
		MemberAuth: auth,
	}
}

func (c *MemberClient) WithSession(key, secret string) *MemberClient {
	return c.WithAuth(&memberSessionAuth{key, secret})
}

// auth

type memberAuth struct {
	key    string
	secret string
	pin    string
	expire time.Duration
}

func (m *memberAuth) token(method, uri string, body []byte) string {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(m.expire).Unix(),
		"sign": signRequest(method, uri, string(body)),
		"key":  m.key,
	}

	if len(m.pin) > 0 {
		payload := map[string]interface{}{
			"t": time.Now().Unix(),
			"n": newNonce(),
			"p": m.pin,
		}

		data, _ := jsoniter.Marshal(payload)
		pinToken, err := rsaEncrypt(data)

		if err != nil {
			log.Panic(err)
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

func (m *MemberClient) VerifyPin(ctx context.Context, pin string) error {
	data, err := m.POST("/pin").Auth(m.PresignWithPin(pin, time.Minute)).Do(ctx).Bytes()
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	var resp Err

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return err
	}

	return gatewayErr(resp)
}

func (m *MemberClient) UpdatePin(ctx context.Context, pin, newPin string) error {
	pinToken, err := rsaEncrypt([]byte(newPin))
	if err != nil {
		return err
	}

	data, err := m.PUT("/pin").Auth(m.PresignWithPin(pin, time.Minute)).P("pin", pinToken).Do(ctx).Bytes()
	if err != nil {
		return err
	}

	var resp Err

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return err
	}

	return gatewayErr(resp)
}

func (m *MemberClient) UpdateProfile(ctx context.Context, fullname, avatar string) (*MemberView, error) {
	data, err := m.PUT("/profile").Auth(m.Presign(time.Minute)).P("fullname", fullname).P("avatar", avatar).Do(ctx).Bytes()
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
