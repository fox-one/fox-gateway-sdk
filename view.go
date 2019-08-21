package gateway

import (
	"github.com/shopspring/decimal"
)

// MemberView member
type MemberView struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	IsPinSet  bool   `json:"is_pin_set"`
}

// MemberSessionView member session
type MemberSessionView struct {
	Key       string `json:"key"`
	Secret    string `json:"secret"`
	CreatedAt int64  `json:"created_at"`
	ExpiredAt int64  `json:"expired_at"`
}

// AdminUserView admin user
type AdminUserView struct {
	ID        uint   `json:"id"`
	CreatedAt int64  `json:"created_at"`
	Username  string `json:"username"`
	Merchant  string `json:"merchant"`
}

// AdminSessionView session of admin
type AdminSessionView struct {
	Key       string `json:"key"`
	Secret    string `json:"secret"`
	CreatedAt int64  `json:"created_at"`
	ExpiredAt int64  `json:"expired_at"`
}

// MemberWalletView member Wallet
type MemberWalletView struct {
	Label      string `json:"label"`
	MemberID   string `json:"member_id"`
	Service    string `json:"service"`
	WalletID   string `json:"wallet_id"`
	SessionID  string `json:"session_id,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
}

// wallet

type WalletUserView struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
	Avatar   string `json:"avatar"`
}

// Asset asset
type WalletAssetView struct {
	AssetID string `json:"asset_id"`
	ChainID string `json:"chain_id"`

	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	IconURL string `json:"icon_url"`

	Price     string `json:"price"`
	PriceUSD  string `json:"price_usd"`
	PriceBTC  string `json:"price_btc"`
	Change    string `json:"change"`
	ChangeUSD string `json:"change_usd"`
	ChangeBTC string `json:"change_btc"`
}

// WalletChainView wallet chain view
type WalletChainView struct {
	WalletAssetView

	Fee           decimal.Decimal `json:"fee"`
	Confirmations int             `jons:"confirmations"`
}

type WalletUserAssetView struct {
	WalletAssetView

	Balance     string `json:"balance"`
	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`

	Chain *WalletChainView `json:"chain"`
}

// WithdrawAddressView withdraw address
type WithdrawAddressView struct {
	AddressID string `json:"address_id,omitempty"`
	AssetID   string `json:"asset_id"`

	PublicKey string `json:"public_key,omitempty"`
	Label     string `json:"label,omitempty"`

	AccountName string `json:"account_name,omitempty"`
	AccountTag  string `json:"account_tag,omitempty"`
}

// Snapshot model
type WalletSnapshotView struct {
	SnapshotID      string      `json:"snapshot_id"`
	TraceID         string      `json:"trace_id"`
	WalletID        string      `json:"wallet_id"`
	AssetID         string      `json:"asset_id"`
	OpponentID      string      `json:"opponent_id"`
	Source          string      `json:"source"`
	Amount          string      `json:"amount"`
	Memo            string      `json:"memo"`
	MemberID        string      `json:"member_id"`
	Service         string      `json:"service"`
	Label           string      `json:"label"`
	CreatedAt       int64       `json:"created_at"`
	TransactionHash string      `json:"transaction_hash"`
	Sender          string      `json:"sender"`
	Receiver        string      `json:"receiver"`
	ExtraData       interface{} `json:"data"`

	Asset WalletAssetView `json:"asset"`
}

// WalletPendingDepositView pending deposit
type WalletPendingDepositView struct {
	Type string `json:"type"`

	TransactionID   string `json:"transaction_id"`
	TransactionHash string `json:"transaction_hash"`
	CreatedAt       int64  `json:"created_at"`

	AssetID       string `json:"asset_id,omitempty"`
	ChainID       string `json:"chain_id,omitempty"`
	Amount        string `json:"amount"`
	Confirmations int    `json:"confirmations"`
	Threshold     int    `json:"threshold"`

	BrokerID    string `json:"broker_id"`
	UserID      string `json:"user_id"`
	Sender      string `json:"sender"`
	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`
}

// Exchange

type ExchangeAssetView struct {
	Symbol    string `json:"symbol,omitempty"`
	AssetID   string `json:"asset_id,omitempty"`
	Name      string `json:"name,omitempty"`
	ChainID   string `json:"chain_id,omitempty"`
	Icon      string `json:"icon_url,omitempty"`
	Type      string `json:"type,omitempty"`
	Precision int    `json:"precision,omitempty"`
}

type ExchangePairView struct {
	Symbol           string `json:"symbol,omitempty"`
	Logo             string `json:"logo,omitempty"`
	BaseAssetId      string `json:"base_asset_id,omitempty"`
	BaseAssetSymbol  string `json:"base_asset,omitempty"`
	QuoteAssetId     string `json:"quote_asset_id,omitempty"`
	QuoteAssetSymbol string `json:"quote_asset,omitempty"`
	PricePrecision   int    `json:"price_precision,omitempty"`
	AmountPrecision  int    `json:"amount_precision,omitempty"`
	Status           string `json:"status,omitempty"`
	BaseMinAmount    string `json:"base_min_amount,omitempty"`
	BaseMaxAmount    string `json:"base_max_amount,omitempty"`
	QuoteMinAmount   string `json:"quote_min_amount,omitempty"`
	QuoteMaxAmount   string `json:"quote_max_amount,omitempty"`
}

func (v *ExchangePairView) Market() string {
	return v.BaseAssetId + "-" + v.QuoteAssetId
}
