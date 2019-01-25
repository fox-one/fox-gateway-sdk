package gateway

import (
	"time"

	"github.com/fox-one/httpclient"
	uuid "github.com/satori/go.uuid"
)

const (
	timestampKey = "_ts"
	nonceKey     = "_nonce"
)

type Client struct {
	client *httpclient.Client
}

func (c *Client) Group(group string) *Client {
	return &Client{c.client.Group(group)}
}

func NewClient(apiBase string) *Client {
	c := httpclient.NewClient(apiBase)
	c.OnRequest = func(req *httpclient.Request, method, uri string) {
		req.P(timestampKey, time.Now().Unix())
		req.P(nonceKey, uuid.Must(uuid.NewV4()).String())
	}

	return &Client{c}
}
