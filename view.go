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
