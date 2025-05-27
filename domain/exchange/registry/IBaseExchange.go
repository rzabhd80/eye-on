package registry

import (
	"context"
	"github.com/rzabhd80/eye-on/domain/balance"
	"github.com/rzabhd80/eye-on/domain/order"
	"github.com/rzabhd80/eye-on/domain/orderBook"
)

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
