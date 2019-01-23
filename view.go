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
