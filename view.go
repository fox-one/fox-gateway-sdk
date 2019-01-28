package gateway

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
	Label    string `json:"label"`
	MemberID string `json:"member_id"`
	Service  string `json:"service"`
	WalletID string `json:"wallet_id"`
}

// wallet

// Asset asset
type WalletAssetView struct {
	AssetID string `json:"asset_id"`
	ChainID string `json:"chain_id"`

	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	IconURL string `json:"icon_url"`

	Price    string `json:"price"`
	PriceUSD string `json:"price_usd"`
	Change   string `json:"change"`
}

type WalletUserAssetView struct {
	WalletAssetView

	Balance     string `json:"balance"`
	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`
}

// Snapshot model
type WalletSnapshotView struct {
	SnapshotID string `json:"snapshot_id"`
	TraceID    string `json:"trace_id"`
	WalletID   string `json:"wallet_id"`
	AssetID    string `json:"asset_id"`
	OpponentID string `json:"opponent_id"`
	Source     string `json:"source"`
	Amount     string `json:"amount"`
	Memo       string `json:"memo"`
	MemberID   string `json:"member_id"`
	Service    string `json:"service"`
	Label      string `json:"label"`
	CreatedAt  int64  `json:"created_at"`

	Asset WalletAssetView `json:"asset"`
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
