package order

import "time"

//OrderRequest & OrderResponse unified across all exchanges

type OrderRequest struct {
	Symbol      string   `json:"symbol"`
	Side        string   `json:"side"` // buy/sell
	Type        string   `json:"type"` // limit/market
	Quantity    float64  `json:"quantity"`
	Price       *float64 `json:"price,omitempty"`
	TimeInForce string   `json:"time_in_force,omitempty"` // GTC, IOC, FOK
	ClientID    string   `json:"client_id,omitempty"`
}

type OrderResponse struct {
	ID            string    `json:"id"`
	ClientOrderID string    `json:"client_order_id"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Type          string    `json:"type"`
	Quantity      float64   `json:"quantity"`
	Price         *float64  `json:"price,omitempty"`
	Status        string    `json:"status"`
	ExecutedQty   float64   `json:"executed_qty"`
	ExecutedPrice float64   `json:"executed_price"`
	Commission    float64   `json:"commission"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
