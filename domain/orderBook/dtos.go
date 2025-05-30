package orderBook

import (
	"github.com/rzabhd80/eye-on/internal/database/models"
)

type GetOrderBookRequest struct {
	Symbol string `json:"symbol"`
	Limit  int    `json:"limit,omitempty"` // Number of price levels, defaults to 100
}

type StandardOrderBookRequest struct {
	Symbol string `json:"symbol"`
}

type StandardOrderBookResponse struct {
	Symbol    string       `json:"symbol"`
	Bids      models.JSONB `json:"bids"` // Buying orders (price descending)
	Asks      models.JSONB `json:"asks"` // Selling orders (price ascending)
	Timestamp string       `json:"timestamp"`
}

// StandardOrderLevel represents price level in order book
type StandardOrderLevel struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}
