package gateway

import (
	"context"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fox-one/httpclient"
	jsoniter "github.com/json-iterator/go"
)

type MerchantClient struct {
	*Client
	key    string
	secret string
}

func (c *Client) Merchant(key, secret string) *MerchantClient {
	return &MerchantClient{
		key:    key,
		secret: secret,
		Client: c.Group("merchant"),
	}
}

func NewMerchantClient(key, secret, apiBase string) *MerchantClient {
	return NewClient(apiBase).Merchant(key, secret)
}

// auth
type merchantAuth struct {
	*MerchantClient
	memberID string
	expire   time.Duration
}

func (m *merchantAuth) Auth(req *httpclient.Request, method, uri string, body []byte) {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(m.expire).Unix(),
		"sign": signRequest(method, uri, string(body)),
		"key":  m.key,
	}

	if len(m.memberID) > 0 {
		claims["mem"] = m.memberID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(m.secret))
	if err != nil {
		log.Panicln("sign merchent token", err)
	}

	req.AddToken(t)
}

func (m *MerchantClient) PresignMember(memberID string, expire time.Duration) httpclient.Authenticator {
	return &merchantAuth{
		MerchantClient: m,
		expire:         expire,
		memberID:       memberID,
	}
}

func (m *MerchantClient) Presign(expire time.Duration) httpclient.Authenticator {
	return m.PresignMember("", expire)
}

// member

func (m *MerchantClient) CreateMember(ctx context.Context) (*MemberView, *MemberSessionView, error) {
	data, err := m.client.POST("/member/new").Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Err

		Member  *MemberView        `json:"member"`
		Session *MemberSessionView `json:"session"`
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, nil, err
	}

	if resp.Code > 0 {
		return nil, nil, resp.Err
	}

	return resp.Member, resp.Session, nil
}

func (m *MerchantClient) LoginMember(ctx context.Context, id string, expire time.Duration) (*MemberView, *MemberSessionView, error) {
	data, err := m.client.POST("/member/login").
		P("id", id).
		P("expire", int64(expire.Seconds())).
		Auth(m.Presign(time.Minute)).
		Do(ctx).Bytes()
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Err

		Member  *MemberView        `json:"member"`
		Session *MemberSessionView `json:"session"`
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, nil, err
	}

	if resp.Code > 0 {
		return nil, nil, resp.Err
	}

	return resp.Member, resp.Session, nil
}

// ClearUserSessions clear user session
func (m *MerchantClient) ClearUserSessions(ctx context.Context, memberID string, sessionKey ...string) error {
	req := m.client.POST("/member/logout").P("id", memberID)
	if len(sessionKey) > 0 {
		req = req.P("session_key", sessionKey[0])
	}
	data, err := req.Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
	if err != nil {
		return err
	}

	var resp struct {
		Err
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return err
	}

	if resp.Code > 0 {
		return resp.Err
	}

	return nil
}
