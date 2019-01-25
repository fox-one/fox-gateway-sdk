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

// wallet

// Asset asset
type WalletAssetView struct {
	AssetID  string `json:"asset_id"`
	AssetKey string `json:"asset_key,omitempty"`
	ChainID  string `json:"chain_id"`

	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	IconURL string `json:"icon_url"`
}

// UserAddress user address
type WalletUserAddressView struct {
	UserID  string `json:"user_id"`
	ChainID string `json:"chain_id"`

	PublicKey   string `json:"public_key"`
	AccountName string `json:"account_name"`
	AccountTag  string `json:"account_tag"`

	Confirmations  int     `json:"confirmations"`
	Capitalization float64 `json:"capitalization"`
}

// WalletUserAssetView wallet user asset view
type WalletUserAssetView struct {
	AssetID           string                 `json:"asset_id"`
	Balance           string                 `json:"balance"`
	TransactionAmount string                 `json:"transaction_amount"`
	TransactionCount  int64                  `json:"transaction_count"`
	Asset             *WalletAssetView       `json:"asset,omitempty"`
	Address           *WalletUserAddressView `json:"address,omitempty"`
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
