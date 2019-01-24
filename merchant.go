package gateway

import (
	"context"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
)

type Merchant struct {
	client *Client
	key    string
	secret string
}

func (c *Client) Merchant(key, secret string) *Merchant {
	return &Merchant{
		key:    key,
		secret: secret,
		client: c.Group("merchant"),
	}
}

func NewMerchant(key, secret, apiBase string) *Merchant {
	return NewClient(apiBase).Merchant(key, secret)
}

// auth

type merchantAuth struct {
	*Merchant
	memberID string
	expire   time.Duration
}

func (m *merchantAuth) GenerateToken(method, uri string, body []byte) string {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(m.expire).Unix(),
		"sign": signRequest(method, uri, body),
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

	return t
}

func (m *Merchant) PresignMember(memberID string, expire time.Duration) Authenticator {
	return &merchantAuth{
		Merchant: m,
		expire:   expire,
		memberID: memberID,
	}
}

func (m *Merchant) Presign(expire time.Duration) Authenticator {
	return m.PresignMember("", expire)
}

func (m *Merchant) CreateMember(ctx context.Context) (*MemberView, *MemberSessionView, error) {
	data, err := m.client.POST("/member/new").Auth(m.Presign(time.Minute)).Do(ctx)
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

func (m *Merchant) LoginMember(ctx context.Context, id string, expire time.Duration) (*MemberView, *MemberSessionView, error) {
	data, err := m.client.POST("/member/login").
		P("id", id).
		P("expire", int64(expire.Seconds())).
		Auth(m.Presign(time.Minute)).
		Do(ctx)
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
func (m *Merchant) ClearUserSessions(ctx context.Context, memberID string, sessionKey ...string) error {
	req := m.client.POST("/member/logout").P("id", memberID)
	if len(sessionKey) > 0 {
		req = req.P("session_key", sessionKey[0])
	}
	data, err := req.Auth(m.Presign(time.Minute)).Do(ctx)
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
