package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"

	jsoniter "github.com/json-iterator/go"
)

func TestUnmarshalAssets(t *testing.T) {
	type resp struct {
		Err
		WalletUserAssetsView
	}

	r := resp{}
	r.Code = 200
	r.Msg = "error"

	data, _ := jsoniter.MarshalToString(r)
	assert.Empty(t, data)

	assets := WalletUserAssetsView{}
	data, _ = jsoniter.MarshalToString(assets)
	assert.Empty(t, data)

	r = resp{}
	err := jsoniter.UnmarshalFromString(data, &r)
	assert.Nil(t, err)
}
