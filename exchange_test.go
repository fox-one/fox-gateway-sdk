package gateway

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeMarkets(t *testing.T) {
	c := NewExchangeClient(apiBase)
	ctx := context.Background()

	// assets
	{
		assets, err := c.MarketAssets(ctx)
		assert.Nil(t, err)
		assert.NotEmpty(t, assets)
	}

	// pairs
	{
		pairs, err := c.MarketPairs(ctx)
		assert.Nil(t, err)
		assert.NotEmpty(t, pairs)
	}
}
