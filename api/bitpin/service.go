package bitpin

import (
	"context"
	"github.com/rzabhd80/eye-on/domain/balance"
	"github.com/rzabhd80/eye-on/domain/exchange"
	"github.com/rzabhd80/eye-on/domain/exchangeCredentials"
	"github.com/rzabhd80/eye-on/domain/order"
	"github.com/rzabhd80/eye-on/domain/orderBook"
	"github.com/rzabhd80/eye-on/domain/user"
)

type BitpinExchange struct {
	ExchangeRepo           *exchange.ExchangeRepository
	ExchangeCredentialRepo *exchangeCredentials.ExchangeCredentialRepository
	UserREpo               *user.UserRepository
}

func (exchange *BitpinExchange) Name() string                   { return "" }
func (exchange *BitpinExchange) Ping(ctx context.Context) error { return nil }
func (exchange *BitpinExchange) GetBalance(ctx context.Context) ([]balance.Balance, error) {
	return nil, nil
}
func (exchange *BitpinExchange) GetOrderBook(ctx context.Context, symbol string) (*orderBook.OrderBook, error) {
	return nil, nil
}
func (exchange *BitpinExchange) PlaceOrder(ctx context.Context, req *order.OrderRequest) (*order.Order, error) {
	return nil, nil
}

func (exchange *BitpinExchange) GetOrder(ctx context.Context, symbol, orderID string) (*order.Order, error) {
	return nil, nil
}

func (exchange *BitpinExchange) CancelOrder(ctx context.Context, symbol, orderID string) error {
	return nil
}
func (exchange *BitpinExchange) GetOpenOrders(ctx context.Context, symbol string) ([]order.Order, error) {
	return nil, nil
}

func (exchange *BitpinExchange) GetOrderHistory(ctx context.Context, symbol string, limit int) ([]order.Order, error) {
	return nil, nil
}
