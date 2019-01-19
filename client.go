package gateway

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
)

var defaultHttpClient = &http.Client{
	Transport: &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			dialer := net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			return dialer.Dial(network, addr)
		},
	},
}

type Client struct {
	apiBase string

	client *http.Client
}

func NewClient(apiBase string) *Client {
	return &Client{
		apiBase: apiBase,
		client:  defaultHttpClient,
	}
}

type Authenticator interface {
	GenerateToken(method, uri string, body []byte) string
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
		params: map[string]interface{}{
			"_ts":    time.Now().Unix(),
			"_nonce": uuid.Must(uuid.NewV4()).String(),
		},
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
	r.params["_nonce"] = nonce
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
	u := &url.URL{}
	u.Path = r.uri

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

	u.Scheme = "https"
	u.Host = r.c.apiBase

	request, err := http.NewRequest(r.method, u.String(), bytes.NewBuffer(body))
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
