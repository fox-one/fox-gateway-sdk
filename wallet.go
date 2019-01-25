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

// AssetList model
type AssetList []Asset

// Implements the AssetList error interface
func (assetList AssetList) Error() string {
	return "asset list error"
}

// Snapshot model
type Snapshot struct {
	SnapshotID string      `json:"snapshot_id"`
	TraceID    string      `json:"trace_id"`
	WalletID   string      `json:"wallet_id"`
	AssetID    string      `json:"asset_id"`
	OpponentID string      `json:"opponent_id"`
	Source     string      `json:"source"`
	Amount     string      `json:"amount"`
	Memo       string      `json:"memo"`
	MemberID   string      `json:"member_id"`
	Service    string      `json:"service"`
	Label      string      `json:"label"`
	Data       interface{} `json:"data"`
	CreatedAt  int64       `json:"created_at"`
	Asset      Asset       `json:"asset"`
}

// Implements the error interface
func (snapshot Snapshot) Error() string {
	return "snapshot error"
}

// SnapshotList model
type SnapshotList []Snapshot

// Implements the AssetList error interface
func (snapshotList SnapshotList) Error() string {
	return "snapshot list error"
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

// WalletList model
type WalletList []Wallet

// Implements the error interface
func (walletList WalletList) Error() string {
	return "wallet list error"
}

// GetAssets 取用户的资产
func (memberService *MemberService) GetAssets(ctx context.Context) (AssetList, error) {
	resp, err := memberService.GET("/assets").
		Auth(memberService.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var assets AssetList

	err = jsoniter.Unmarshal(resp, &assets)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return assets, nil
}

// GetAssetsByID 读取用户某币种余额
func (memberService *MemberService) GetAssetsByID(ctx context.Context, assetID string) (AssetList, error) {
	resp, err := memberService.GET("/assets").
		P("asset_id", assetID).
		Auth(memberService.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var assets AssetList

	err = jsoniter.Unmarshal(resp, &assets)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return assets, nil
}

// GetSnapshots 读取转账记录列表
func (memberService *MemberService) GetSnapshots(ctx context.Context, assetID string,
	cursor string, limit int, order string) (SnapshotList, error) {
	resp, err := memberService.GET("/snapshots").
		P("asset_id", assetID).
		P("cursor", cursor).
		P("limit", limit).
		P("order", order).
		Auth(memberService.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	var snapshots SnapshotList

	err = jsoniter.Unmarshal(resp, &snapshots)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return snapshots, nil
}

// Withdraw 提现
func (memberService *MemberService) Withdraw(ctx context.Context, body string) (Snapshot, error) {
	resp, err := memberService.POST("/withdraw").
		P("body", body).
		Auth(memberService.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		log.Error(err)
		// return nil, err
	}

	var snapshot Snapshot

	err = jsoniter.Unmarshal(resp, &snapshot)

	if err != nil {
		log.Error(err)
		// return err
	}

	return snapshot, nil
}

// Transfer 转账
func (memberService *MemberService) Transfer(ctx context.Context, body string) (Snapshot, error) {
	resp, err := memberService.POST("/transfer").
		P("body", body).
		Auth(memberService.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		log.Error(err)
		// return nil, err
	}

	var snapshot Snapshot

	err = jsoniter.Unmarshal(resp, &snapshot)

	if err != nil {
		log.Error(err)
		// return err
	}

	return snapshot, nil
}
