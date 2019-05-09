package gateway

import (
	"time"

	"github.com/fox-one/httpclient"
)

const (
	timestampKey = "_ts"
	nonceKey     = "_nonce"
)

type Client struct {
	*httpclient.Client
}

func (c *Client) Group(group string) *Client {
	return &Client{c.Client.Group(group)}
}

func NewClient(apiBase string) *Client {
	c := httpclient.NewClient(apiBase)
	c.OnRequest = func(req *httpclient.Request, method, uri string) {
		req.Q(timestampKey, time.Now().Unix())
		req.Q(nonceKey, newNonce())
	}

	return &Client{c}
}
