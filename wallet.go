package gateway

import (
	"context"
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// GetAssets 取用户的资产
func (m *MemberService) ReadAssets(ctx context.Context, chain int) ([]*WalletUserAssetView, error) {
	result := m.GET("/assets").P("chain", chain).
		Auth(m.Presign(time.Minute)).
		Do(ctx)

	data, err := result.Bytes()
	if err != nil {
		return nil, err
	}

	var assets []*WalletUserAssetView
	if err = jsoniter.Unmarshal(data, &assets); err == nil {
		return assets, nil
	}

	var e Err
	if jsoniter.Unmarshal(data, &e) == nil && e.Code > 0 {
		return nil, e
	}

	_, status := result.Status()
	return nil, errors.New(status)
}

// GetAssetsByID 读取用户某币种余额
func (m *MemberService) ReadAsset(ctx context.Context, assetID string) (*WalletUserAssetView, error) {
	data, err := m.GET("/asset/" + assetID).
		Auth(m.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))

	var resp struct {
		Err
		Asset *WalletUserAssetView `json:"asset"`
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.Asset, nil
}

// GetSnapshots 读取转账记录列表
func (m *MemberService) ReadSnapshots(ctx context.Context, assetID string, cursor string, limit int, order string) ([]*WalletSnapshotView, *Pagination, error) {
	path := "/snapshots"
	if assetID != "" {
		path += "/" + assetID
	}

	result := m.GET(path).
		P("cursor", cursor).
		P("limit", limit).
		P("order", order).
		Auth(m.Presign(time.Minute)).
		Do(ctx)

	data, err := result.Bytes()
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Err
		Snapshots  []*WalletSnapshotView `json:"snapshots"`
		Pagination *Pagination           `json:"pagination"`
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, nil, err
	}

	if resp.Code > 0 {
		return nil, nil, resp.Err
	}

	return resp.Snapshots, resp.Pagination, nil
}

func (m *MemberService) ReadSnapshot(ctx context.Context, id string) (*WalletSnapshotView, error) {
	result := m.GET("/snapshot/" + id).
		Auth(m.Presign(time.Minute)).
		Do(ctx)

	data, err := result.Bytes()
	if err != nil {
		return nil, err
	}

	var resp struct {
		Err
		*WalletSnapshotView
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.WalletSnapshotView, nil
}

type WalletAssetOperation struct {
	AssetID string `json:"asset_id"`
	Amount  string `json:"amount"`
	Memo    string `json:"memo"`
	TraceID string `json:"trace_id"`
}

type WalletWithdrawOperation struct {
	WalletAssetOperation

	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`
}

// Withdraw 提现
func (m *MemberService) Withdraw(ctx context.Context, op *WalletWithdrawOperation) (*WalletSnapshotView, error) {
	data, err := m.POST("/withdraw").
		Body(op).
		Auth(m.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		return nil, err
	}

	var resp struct {
		Err
		*WalletSnapshotView
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.WalletSnapshotView, nil
}

type WalletTransferOperation struct {
	WalletAssetOperation

	OpponentID string `json:"opponent_id"`
}

// Transfer 转账
func (m *MemberService) Transfer(ctx context.Context, op *WalletTransferOperation) (*WalletSnapshotView, error) {
	data, err := m.POST("/transfer").
		Body(op).
		Auth(m.Presign(time.Minute)).
		Do(ctx).Bytes()

	if err != nil {
		return nil, err
	}

	var resp struct {
		Err
		*WalletSnapshotView
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.WalletSnapshotView, nil
}

// wallet public

func (c *Client) WithdrawFee(ctx context.Context, assetId string, publicKey string, accountName, accountTag string) (*WalletAssetView, string, error) {
	req := c.GET("/wallet/withdraw-fee").
		Q("asset_id", assetId).
		Q("public_key", publicKey).
		Q("account_name", accountName).
		Q("account_tag", accountTag)

	data, err := req.Do(ctx).Bytes()
	if err != nil {
		return nil, "", err
	}

	var resp struct {
		Err
		Asset *WalletAssetView `json:"fee_asset"`
		Fee   string           `json:"fee_amount"`
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, resp.Fee, err
	}

	if resp.Code > 0 {
		return nil, resp.Fee, resp.Err
	}

	return resp.Asset, resp.Fee, nil
}
