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
	assetID := "43d61dcd-e413-450d-80b8-101d5e903357"
	ctx := context.Background()
	asset, fee, err := c.WithdrawFee(ctx, assetID, pk, "")
	if assert.Nil(t, err) {
		assert.Equal(t, assetID, asset.AssetID)
		assert.NotEmpty(t, asset.Price)
		f, _ := decimal.NewFromString(fee)
		assert.True(t, f.IsPositive())
	}
}

func TestSearchUser(t *testing.T) {
	c := NewClient("https://dev-gateway.fox.one")
	ctx := context.Background()
	user, err := c.SearchWalletUser(ctx, "8017d200-7870-4b82-b53f-74bae1d2dad7")
	if assert.Nil(t, err) {
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Fullname)
		assert.NotEmpty(t, user.Avatar)
		assert.Equal(t, "yiplee", user.Fullname)
	}
}
