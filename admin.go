package gateway

import (
	"context"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fox-one/httpclient"
	jsoniter "github.com/json-iterator/go"
)

// Admin admin
type AdminClient struct {
	client *Client

	key    string
	secret string
}

// NewAdmin new admin
func NewAdminClient(apiBase string) *AdminClient {
	return NewClient(apiBase).Admin()
}

func (c *Client) Admin() *AdminClient {
	return &AdminClient{client: c.Group("admin")}
}

func (c *AdminClient) WithSession(key, secret string) *AdminClient {
	return &AdminClient{
		client: c.client,
		key:    key,
		secret: secret,
	}
}

type adminAuth struct {
	*AdminClient
	expire time.Duration
}

func (a *adminAuth) token(method, uri string, body []byte) string {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(a.expire).Unix(),
		"sign": signRequest(method, uri, string(body)),
		"key":  a.key,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(a.secret))
	if err != nil {
		log.Panicln("sign merchent token", err)
	}

	return t
}

func (a *adminAuth) Auth(req *httpclient.Request, method, uri string, body []byte) {
	req.AddToken(a.token(method, uri, body))
}

func (c *AdminClient) Presign(expire time.Duration) *adminAuth {
	return &adminAuth{
		AdminClient: c,
		expire:      expire,
	}
}

// Validate validate
func (a *AdminClient) Validate(ctx context.Context, token, method, uri string, body []byte) (*AdminUserView, error) {
	data, err := a.client.POST("/validate").
		P("method", method).
		P("uri", uri).
		P("body", body).
		P("token", token).
		Do(ctx).Bytes()
	if err != nil {
		return nil, err
	}

	var resp struct {
		Err
		AdminUser *AdminUserView `json:"admin"`
	}
	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.AdminUser, nil
}

func (a *AdminClient) CreateAdmin(ctx context.Context, name, password string) (*AdminUserView, *AdminSessionView, error) {
	data, err := a.client.POST("/create-admin").
		P("username", name).P("password", password).
		Auth(a.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Err
		AdminUser        *AdminUserView    `json:"admin"`
		AdminSessionView *AdminSessionView `json:"session"`
	}
	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, nil, err
	}

	if resp.Code > 0 {
		return nil, nil, resp.Err
	}

	return resp.AdminUser, resp.AdminSessionView, nil
}
