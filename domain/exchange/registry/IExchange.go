package registry

import (
	"context"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/order"
	"github.com/rzabhd80/eye-on/internal/database/models"
)

type Symbol struct {
	Symbol     string
	BaseAsset  string
	QuoteAsset string
}
type ExchangeConfig struct {
	Name          string
	DisplayName   string
	BaseURL       string
	RateLimit     int
	Features      map[string]interface{} // Will be stored as JSONB
	SymbolFactory ISymbolFactory
}

type ISymbolFactory interface {
	RegisterExchangeSymbols(bitpinExchange *models.Exchange) *[]models.TradingPair
}
type IExchange interface {
	Name() string
	Ping(ctx context.Context) error
	GetBalance(ctx context.Context, userId uuid.UUID, sign *string) ([]models.BalanceSnapshot, error)
	GetOrderBook(ctx context.Context, symbol string, userId uuid.UUID) (*models.OrderBookSnapshot, error)
	PlaceOrder(ctx context.Context, req *order.StandardOrderRequest, userId uuid.UUID) (*models.OrderHistory, error)
	CancelOrder(ctx context.Context, orderID string, userId uuid.UUID) error
}
