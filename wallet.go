package gateway

import (
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

// Asset Model
type Asset struct {
	AssetID string `json:"asset_id"`
	ChainID string `json:"chain_id"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	IconURL string `json:"icon_url"`
}

// Implements the error interface
func (asset Asset) Error() string {
	return asset.AssetID
}

// AssetList xx
type AssetList []Asset

// Implements the error interface
func (assetList AssetList) Error() string {
	return "asset"
}

// Wallet Model
type Wallet struct {
	Label    string `json:"label"`
	MemberID string `json:"member_id"`
	Service  string `json:"service"`
	WalletID string `json:"wallet_id"`
}

// Implements the error interface
func (wallet Wallet) Error() string {
	return wallet.Label
}

// GetAssets 取用户的资产
func (memberService *MemberService) GetAssets(ctx context.Context, memberID string) ([]Asset, error) {
	resp, err := memberService.GET("/assets").Auth(memberService.Presign(time.Minute)).Do(ctx).Bytes()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var assets []Asset

	err = jsoniter.Unmarshal(resp, &assets)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return assets, nil
}
