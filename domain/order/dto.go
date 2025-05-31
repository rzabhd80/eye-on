package order

import (
	"time"
)

//OrderRequest & OrderResponse unified across all exchanges

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

type OrderStatus string

const (
	NEW       OrderStatus = "new"
	FILLED    OrderStatus = "filled"
	PARTIALLY OrderStatus = "partially_filled"
	CANCELED  OrderStatus = "canceled"
	REJECTED  OrderStatus = "rejected"
)

type StandardOrderRequest struct {
	Symbol        string    `json:"symbol" validate:"required"`
	Side          OrderSide `json:"side" validate:"required,oneof=buy sell"`
	Type          OrderType `json:"type" validate:"required,oneof=market limit"`
	Quantity      *float64  `json:"quantity,omitempty"` // Amount of base asset
	BaseCurrency  string    `json:"base_currency,omitempty"`
	QuoteCurrency string    `json:"Quote_currency,omitempty"`
	BaseAmount    *float64  `json:"base_amount,omitempty"`     // Amount of base asset (e.g., BTC amount)
	QuoteAmount   *float64  `json:"quote_amount,omitempty"`    // Amount of quote asset (e.g., USDT amount)
	Price         *float64  `json:"price,omitempty"`           // Price per unit
	StopPrice     *float64  `json:"stop_price,omitempty"`      // For stop orders
	TimeInForce   string    `json:"time_in_force,omitempty"`   // GTC, IOC, FOK, etc.
	ClientOrderId string    `json:"client_order_id,omitempty"` // Client-specified order ID
}

// StandardOrderResponse represents unified order response
type StandardOrderResponse struct {
	ID         string      `json:"id"`
	Symbol     string      `json:"symbol"`
	Side       OrderSide   `json:"side"`
	Type       OrderType   `json:"type"`
	Quantity   float64     `json:"quantity"`
	Price      *float64    `json:"price,omitempty"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	ExchangeID string      `json:"exchange_id"`
}

type CreateOrderRequest struct {
	Symbol   string    `json:"symbol"`          // e.g., "BTCUSDT"
	Side     OrderSide `json:"side"`            // "buy" or "sell"
	Type     OrderType `json:"type"`            // "market" or "limit"
	Quantity string    `json:"quantity"`        // Amount to buy/sell
	Price    string    `json:"price,omitempty"` // Required for limit orders
}

type CancelOrderRequest struct {
	OrderId string   `json:"orderId"`
	Hours   *float64 `json:"hours,omitempty"`
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
