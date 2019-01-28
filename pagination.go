package gateway

type Pagination struct {
	Next    string `json:"next_cursor"`
	HasNext bool   `json:"has_next"`
}
