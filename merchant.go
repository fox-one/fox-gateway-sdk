package gateway

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fox-one/httpclient"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

type MerchantClient struct {
	*Client
	key    string
	secret string
}

func (c *Client) Merchant() *MerchantClient {
	return &MerchantClient{
		Client: c.Group("merchant"),
	}
}

func NewMerchantClient(apiBase string) *MerchantClient {
	return NewClient(apiBase).Merchant()
}

func (c *MerchantClient) WithSession(key, secret string) *MerchantClient {
	return &MerchantClient{
		Client: c.Client,
		key:    key,
		secret: secret,
	}
}

// auth
type merchantAuth struct {
	*MerchantClient
	memberID string
	expire   time.Duration
}

func (m *merchantAuth) token(method, uri string, body []byte) string {
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

	return t
}

func (m *merchantAuth) Auth(req *httpclient.Request, method, uri string, body []byte) {
	req.AddToken(m.token(method, uri, body))
}

func (m *MerchantClient) PresignMember(memberID string, expire time.Duration) *merchantAuth {
	return &merchantAuth{
		MerchantClient: m,
		expire:         expire,
		memberID:       memberID,
	}
}

func (m *MerchantClient) Presign(expire time.Duration) *merchantAuth {
	return m.PresignMember("", expire)
}

func (m *MerchantClient) Sign(method, uri string, body []byte, expire time.Duration) string {
	return m.Presign(expire).token(method, uri, body)
}

func (m *MerchantClient) SignMember(memberID, method, uri string, body []byte, expire time.Duration) string {
	return m.PresignMember(memberID, expire).token(method, uri, body)
}

// member

func (m *MerchantClient) CreateMember(ctx context.Context) (*MemberView, *MemberSessionView, error) {
	data, err := m.POST("/member/new").Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
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
	data, err := m.POST("/member/login").
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
	req := m.POST("/member/logout").P("id", memberID)
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
