package exchange

type OrderSide string
type OrderType string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

type CreateOrderRequest struct {
	Symbol   string    `json:"symbol"`          // e.g., "BTCUSDT"
	Side     OrderSide `json:"side"`            // "buy" or "sell"
	Type     OrderType `json:"type"`            // "market" or "limit"
	Quantity string    `json:"quantity"`        // Amount to buy/sell
	Price    string    `json:"price,omitempty"` // Required for limit orders
}

type CreateOrderResponse struct {
	OrderID   string `json:"order_id"`
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Type      string `json:"type"`
	Quantity  string `json:"quantity"`
	Price     string `json:"price"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}
