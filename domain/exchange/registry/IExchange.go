package registry

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/balance"
	"github.com/rzabhd80/eye-on/domain/order"
	"github.com/rzabhd80/eye-on/domain/orderBook"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
}
type Symbol struct {
	Symbol     string
	BaseAsset  string
	QuoteAsset string
}
type ExchangeConfig struct {
	Name        string
	DisplayName string
	BaseURL     string
	APIKey      string
	SecretKey   string
	Passphrase  string
	IsTestnet   bool
	RateLimit   int
	Timeout     time.Duration
	RetryConfig RetryConfig
	Features    map[string]interface{} // Will be stored as JSONB
	UserID      uuid.UUID              // Required for ExchangeCredential
	Label       string                 // defaults to "Default"
	symbols     []Symbol
}

type IExchange interface {
	Name() string
	Ping(ctx context.Context) error
	GetBalance(ctx context.Context) ([]balance.Balance, error)
	GetOrderBook(ctx context.Context, symbol string) (*orderBook.OrderBook, error)
	PlaceOrder(ctx context.Context, req *order.OrderRequest) (*order.Order, error)
	GetOrder(ctx context.Context, symbol, orderID string) (*order.Order, error)
	CancelOrder(ctx context.Context, symbol, orderID string) error
	GetOpenOrders(ctx context.Context, symbol string) ([]order.Order, error)
	GetOrderHistory(ctx context.Context, symbol string, limit int) ([]order.Order, error)
}
