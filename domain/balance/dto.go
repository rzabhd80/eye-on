package balance

type GetBalanceRequest struct {
	Asset string `json:"asset,omitempty"` // If empty, get all balances
}

type StandardBalanceResponse struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
	Total  float64 `json:"total"`
}
