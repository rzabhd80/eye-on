package balance

import "github.com/rzabhd80/eye-on/internal/database/models"

type GetBalanceRequest struct {
	Asset string `json:"asset,omitempty"` // If empty, get all balances
}

type BalanceSchema struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
	Total  string `json:"total"`
}

type GetBalanceResponse struct {
	Balances  []models.BalanceSnapshot `json:"balances"`
	Timestamp int64                    `json:"timestamp"`
}
