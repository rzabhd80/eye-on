package middleware

import (
	"fmt"
	"github.com/rzabhd80/eye-on/domain/order"
)

func ValidateOrderRequest(req *order.StandardOrderRequest) error {
	if req.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}

	if req.Side != order.OrderSideBuy && req.Side != order.OrderSideSell {
		return fmt.Errorf("invalid side: must be 'buy' or 'sell'")
	}

	if req.Type != order.OrderTypeMarket && req.Type != order.OrderTypeLimit {
		return fmt.Errorf("invalid type: must be 'market' or 'limit'")
	}

	// For limit orders, price is required
	if req.Type == order.OrderTypeLimit && req.Price == nil {
		return fmt.Errorf("price is required for limit orders")
	}

	// Must have either Quantity OR (BaseAmount/QuoteAmount)
	hasQuantity := req.Quantity != nil && *req.Quantity > 0
	hasBaseAmount := req.BaseAmount != nil && *req.BaseAmount > 0
	hasQuoteAmount := req.QuoteAmount != nil && *req.QuoteAmount > 0

	if !hasQuantity && !hasBaseAmount && !hasQuoteAmount {
		return fmt.Errorf("must specify either quantity or base_amount/quote_amount")
	}

	return nil
}
