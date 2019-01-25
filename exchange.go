package gateway

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

type ExchangeClient struct {
	client *Client
}

func (c *Client) Exchange() *ExchangeClient {
	return &ExchangeClient{c.Group("exchange")}
}

func NewExchangeClient(apiBase string) *ExchangeClient {
	return NewClient(apiBase).Exchange()
}

func (c *ExchangeClient) MarketAssets(ctx context.Context) ([]*ExchangeAssetView, error) {
	r := c.client.GET("/market/assets").Do(ctx)
	data, err := r.Bytes()
	if err != nil {
		return nil, err
	}

	assets := []*ExchangeAssetView{}
	if jsoniter.Unmarshal(data, &assets) == nil {
		return assets, nil
	}

	if e := decodeErr(data); e.Code > 0 {
		return nil, e
	}

	_, status := r.Status()
	return nil, errors.New(status)
}

func (c *ExchangeClient) MarketPairs(ctx context.Context) ([]*ExchangePairView, error) {
	r := c.client.GET("/market/pairs").Do(ctx)
	data, err := r.Bytes()
	if err != nil {
		return nil, err
	}

	pairs := []*ExchangePairView{}
	if jsoniter.Unmarshal(data, &pairs) == nil {
		return pairs, nil
	}

	if e := decodeErr(data); e.Code > 0 {
		return nil, e
	}

	_, status := r.Status()
	return nil, errors.New(status)
}
