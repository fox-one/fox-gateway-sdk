package gateway

import (
	"context"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
)

type Merchant struct {
	*Client

	key    string
	secret string
}

func NewMerchant(key, secret, apiBase string) *Merchant {
	return &Merchant{
		key:    key,
		secret: secret,
		Client: NewClient(apiBase),
	}
}

// auth

type merchantAuth struct {
	*Merchant
	expire time.Duration
}

func (m *merchantAuth) GenerateToken(method, uri string, body []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(m.expire).Unix(),
		"sign": signRequest(method, uri, body),
		"key":  m.key,
	})
	t, err := token.SignedString([]byte(m.secret))
	if err != nil {
		log.Panicln("sign merchent token", err)
	}

	return t
}

func (m *Merchant) Presign(expire time.Duration) Authenticator {
	return &merchantAuth{
		Merchant: m,
		expire:   expire,
	}
}

func (m *Merchant) CreateMember(ctx context.Context) (*MemberView, *MemberSessionView, error) {
	data, err := m.POST("/member/new").Auth(m.Presign(time.Minute)).Do(ctx)
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
	data, err := m.POST("/member/login").
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
