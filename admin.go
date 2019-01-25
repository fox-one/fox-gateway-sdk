package gateway

import (
	"context"

	jsoniter "github.com/json-iterator/go"
)

// Admin admin
type AdminClient struct {
	client *Client
}

// NewAdmin new admin
func NewAdminClient(apiBase string) *AdminClient {
	return NewClient(apiBase).Admin()
}

func (c *Client) Admin() *AdminClient {
	return &AdminClient{c.Group("admin")}
}

// Validate validate
func (a *AdminClient) Validate(ctx context.Context, method, uri, body, token string) (*AdminUserView, error) {
	data, err := a.client.POST("/admin/validate").
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
