package exchange

import (
	"context"
	"github.com/rzabhd80/eye-on/internal/database/models"
)

type Exchange struct {
	exchange models.Exchange
}

func (exchange *Exchange) Ping(ctx context.Context) error {
	return nil
}
func (exchange *Exchange) Name() string { return exchange.Name() }

func (exchange *Exchange) GetBalance(ctx context.Context) {

}
