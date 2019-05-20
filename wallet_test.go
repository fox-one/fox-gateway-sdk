package gateway

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWithdrawFee(t *testing.T) {
	c := NewClient("https://dev-gateway.fox.one")
	pk := "0x5f59c030ecc49c2eab15f37de7edf26329ad40999"
	assetId := "43d61dcd-e413-450d-80b8-101d5e903357"
	ctx := context.Background()
	asset, fee, err := c.WithdrawFee(ctx, assetId, pk, "", "")
	if assert.Nil(t, err) {
		assert.Equal(t, assetId, asset.AssetID)
		assert.NotEmpty(t, asset.Price)
		f, _ := decimal.NewFromString(fee)
		assert.True(t, f.IsPositive())
	}
}
