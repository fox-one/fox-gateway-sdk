package gateway

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const (
	timestampKey = "_ts"
	nonceKey     = "_nonce"
)

type Client struct {
	scheme string
	host   string
	path   string

	client *http.Client

	newRequestHandler func(req *Request)
}

func parseAPIBase(apiBase string) (*url.URL, error) {
	if !strings.HasPrefix(apiBase, "http") {
		apiBase = "https://" + apiBase
	}

	return url.Parse(apiBase)
}

func NewClient(apiBase string, paths ...string) *Client {
	u, err := parseAPIBase(apiBase)
	if err != nil {
		panic(err)
	}

	return &Client{
		scheme: u.Scheme,
		host:   u.Host,
		path:   path.Join(paths...),
		client: http.DefaultClient,
	}
}

func (c *Client) Group(uri string) *Client {
	return &Client{
		host:   c.host,
		path:   path.Join(c.path, uri),
		client: c.client,
	}
}

func (c *Client) HandleNewRequest(fn func(*Request)) {
	c.newRequestHandler = fn
}

type Authenticator interface {
	GenerateToken(method, uri string, body []byte) string
}

type AuthenticatorParams interface {
	AuthParams() map[string]interface{}
}

type Request struct {
	c *Client

	method  string
	uri     string
	params  map[string]interface{}
	headers http.Header

	auth Authenticator
}

func (c *Client) req(method, uri string) *Request {
	return &Request{
		c:       c,
		method:  method,
		uri:     uri,
		headers: http.Header{},
	}
}

func (c *Client) GET(uri string) *Request {
	return c.req(http.MethodGet, uri)
}

func (c *Client) DELETE(uri string) *Request {
	return c.req(http.MethodDelete, uri)
}

func (c *Client) PUT(uri string) *Request {
	return c.req(http.MethodPut, uri)
}

func (c *Client) POST(uri string) *Request {
	return c.req(http.MethodPost, uri)
}

func (r *Request) H(key, value string) *Request {
	r.headers.Add(key, value)
	return r
}

func (r *Request) P(key string, value interface{}) *Request {
	if value != nil {
		r.params[key] = value
	} else {
		delete(r.params, key)
	}

	return r
}

func (r *Request) Nonce(nonce string) *Request {
	r.params[nonceKey] = nonce
	return r
}

func (r *Request) Auth(auth Authenticator) *Request {
	r.auth = auth
	return r
}

type tokenStringAuth string

func (token tokenStringAuth) GenerateToken(method, uri string, body []byte) string {
	return string(token)
}

func (r *Request) WithTokenString(token string) *Request {
	return r.Auth(tokenStringAuth(token))
}

func (r *Request) Do(ctx context.Context) ([]byte, error) {
	u, err := url.Parse(path.Join(r.c.path, r.uri))
	if err != nil {
		return nil, err
	}

	if authParams, ok := r.auth.(AuthenticatorParams); ok {
		for k, v := range authParams.AuthParams() {
			r.P(k, v)
		}
	}

	var body []byte

	switch r.method {
	case http.MethodPut, http.MethodPost:
		body, _ = jsoniter.Marshal(r.params)
	default:
		query := u.Query()
		for k, v := range r.params {
			value := fmt.Sprint(v)
			query.Add(k, value)
		}
		u.RawQuery = query.Encode()
	}

	if r.auth != nil {
		token := r.auth.GenerateToken(r.method, u.String(), body)
		r.H("Authorization", "Bearer "+token)
	}

	u.Scheme = r.c.scheme
	u.Host = r.c.host

	request, err := http.NewRequest(r.method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header = r.headers
	request = request.WithContext(ctx)
	resp, err := r.c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
