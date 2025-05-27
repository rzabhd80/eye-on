package orderBook

import "github.com/rzabhd80/eye-on/internal/database/models"

type GetOrderBookRequest struct {
	Symbol string `json:"symbol"`
	Limit  int    `json:"limit,omitempty"` // Number of price levels, defaults to 100
}

type OrderBookSchema struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type GetOrderBookResponse struct {
	Symbol    string                     `json:"symbol"`
	Bids      []models.OrderBookSnapshot `json:"bids"` // Buying orders (price descending)
	Asks      []models.OrderBookSnapshot `json:"asks"` // Selling orders (price ascending)
	Timestamp int64                      `json:"timestamp"`
}
