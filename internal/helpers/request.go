package helpers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rzabhd80/eye-on/domain/order"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthToken string

const (
	ApiKeyAuth      AuthToken = "ApiKey"
	ApiAccToken     AuthToken = "ApiAccToken"
	ApiRefreshToken AuthToken = "ApiRefreshToken"
)

type Request struct {
	client           *http.Client
	rateLimiter      chan struct{}
	symbolMap        map[string]string
	reverseSymbolMap map[string]string
}

func NewRequest(timeout time.Duration) *Request {
	return &Request{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (n *Request) MakeRequest(ctx context.Context, method, endpoint string, body []byte,
	creds *models.ExchangeCredential, baseURL string, addBearer bool, addTokenPhrase bool, apiKey AuthToken) (*http.Response, []byte, error) {
	url := baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if creds != nil && creds.AccessKey != "" && apiKey == ApiAccToken {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		var authToken string
		if addBearer {
			authToken = "Bearer " + creds.AccessKey
		} else if addTokenPhrase {
			authToken = "Token " + creds.AccessKey
		} else {
			authToken = creds.APIKey
		}
		req.Header.Set("Authorization", authToken)
		req.Header.Set("X-Timestamp", timestamp)
	}

	if creds != nil && creds.APIKey != "" && apiKey == ApiKeyAuth {
		var authToken string
		if addBearer {
			authToken = "Bearer " + creds.APIKey
		} else if addTokenPhrase {
			authToken = "Token " + creds.APIKey
		} else {
			authToken = creds.APIKey
		}
		req.Header.Set("Authorization", authToken)

	}
	fmt.Printf("request body: %s\n", string(body))
	resp, err := n.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	fmt.Printf("response body: %s\n", string(respBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response: %w", err)
	}

	return resp, respBody, nil
}

// OrderCalculationHelper helps calculate amounts for different exchange formats
type OrderCalculationHelper struct{}

// ValidateOrderRequest validates that the order request has the required fields
func (h *OrderCalculationHelper) ValidateOrderRequestForNobitex(req *order.StandardOrderRequest) error {
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
	if req.BaseCurrency == "" || req.QuoteCurrency == "" {
		return fmt.Errorf("nobitex expects src and dest currencies'")
	}

	req.BaseCurrency, req.QuoteCurrency = strings.ToLower(req.BaseCurrency), strings.ToLower(req.QuoteCurrency)
	if req.BaseCurrency == "irt" {
		req.BaseCurrency = "rls"
	}
	if req.QuoteCurrency == "irt" {
		req.QuoteCurrency = "rls"
	}

	if req.BaseCurrency != "rls" && req.BaseCurrency != "usdt" {
		return fmt.Errorf("nobitex only supports buying usdt and rials")
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

// GetQuantityForExchange calculates the appropriate quantity for exchanges that need single quantity
func (h *OrderCalculationHelper) GetQuantityForExchange(req *order.StandardOrderRequest) (float64, error) {
	// If quantity is directly specified, use it
	if req.Quantity != nil && *req.Quantity > 0 {
		return *req.Quantity, nil
	}

	// If base_amount is specified, use it as quantity (most common case)
	if req.BaseAmount != nil && *req.BaseAmount > 0 {
		return *req.BaseAmount, nil
	}

	// If only quote_amount is specified and we have price, calculate base amount
	if req.QuoteAmount != nil && *req.QuoteAmount > 0 && req.Price != nil && *req.Price > 0 {
		return *req.QuoteAmount / *req.Price, nil
	}

	return 0, fmt.Errorf("cannot determine quantity from provided amounts")
}

// GetBaseAmountForExchange calculates base amount for exchanges that need it separately
func (h *OrderCalculationHelper) GetBaseAmountForExchange(req *order.StandardOrderRequest) (float64, error) {
	// If base_amount is directly specified, use it
	if req.BaseAmount != nil && *req.BaseAmount > 0 {
		return *req.BaseAmount, nil
	}

	// If quantity is specified, use it as base amount
	if req.Quantity != nil && *req.Quantity > 0 {
		return *req.Quantity, nil
	}

	// If only quote_amount is specified and we have price, calculate base amount
	if req.QuoteAmount != nil && *req.QuoteAmount > 0 && req.Price != nil && *req.Price > 0 {
		return *req.QuoteAmount / *req.Price, nil
	}

	return 0, fmt.Errorf("cannot determine base amount from provided amounts")
}

// GetQuoteAmountForExchange calculates quote amount for exchanges that need it separately
func (h *OrderCalculationHelper) GetQuoteAmountForExchange(req *order.StandardOrderRequest) (float64, error) {
	// If quote_amount is directly specified, use it
	if req.QuoteAmount != nil && *req.QuoteAmount > 0 {
		return *req.QuoteAmount, nil
	}

	// Calculate from base amount and price
	var baseAmount float64

	if req.BaseAmount != nil && *req.BaseAmount > 0 {
		baseAmount = *req.BaseAmount
	} else if req.Quantity != nil && *req.Quantity > 0 {
		baseAmount = *req.Quantity
	} else {
		return 0, fmt.Errorf("cannot determine quote amount without base amount or quantity")
	}

	if req.Price == nil || *req.Price <= 0 {
		return 0, fmt.Errorf("cannot determine quote amount without price")
	}

	return baseAmount * *req.Price, nil
}

func (h *OrderCalculationHelper) ParseSymbolParts(symbol string) (base, quote string, err error) {
	// Handle different symbol formats
	if strings.Contains(symbol, "_") {
		parts := strings.Split(symbol, "_")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid symbol format: %s", symbol)
		}
		return strings.ToUpper(parts[0]), strings.ToUpper(parts[1]), nil
	}

	if strings.Contains(symbol, "-") {
		parts := strings.Split(symbol, "-")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid symbol format: %s", symbol)
		}
		return strings.ToUpper(parts[0]), strings.ToUpper(parts[1]), nil
	}

	// For symbols like BTCUSDT, try to parse common patterns
	symbol = strings.ToUpper(symbol)

	// Common quote currencies to try
	quoteCurrencies := []string{"USDT", "USDC", "BTC", "ETH", "BNB", "IRT", "RLS"}

	for _, quote := range quoteCurrencies {
		if strings.HasSuffix(symbol, quote) {
			base := strings.TrimSuffix(symbol, quote)
			if len(base) > 0 {
				return base, quote, nil
			}
		}
	}

	return "", "", fmt.Errorf("cannot parse symbol: %s", symbol)
}

func (h *OrderCalculationHelper) ConvertToNobitexFormat(req *order.StandardOrderRequest) (map[string]interface{}, error) {
	if err := h.ValidateOrderRequestForNobitex(req); err != nil {
		return nil, err
	}

	quantity, err := h.GetQuantityForExchange(req)
	if err != nil {
		return nil, err
	}

	orderType := "buy"
	if req.Side == order.OrderSideSell {
		orderType = "sell"
	}
	if req.ClientOrderId == "" {
		return nil, errors.New("nobitex expects client order id")
	}

	orderData := map[string]interface{}{
		"type":          orderType,
		"srcCurrency":   strings.ToLower(req.QuoteCurrency),
		"dstCurrency":   strings.ToLower(req.BaseCurrency),
		"amount":        fmt.Sprintf("%.8f", quantity),
		"clientOrderId": req.ClientOrderId,
	}

	if req.Type == order.OrderTypeLimit && req.Price != nil {
		orderData["price"] = fmt.Sprintf("%.8f", *req.Price)
	}

	return orderData, nil
}

func (h *OrderCalculationHelper) ValidateOrderRequestForBitpin(req *order.StandardOrderRequest) error {
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

// ConvertToBitpinFormat converts standard request to Bitpin format
func (h *OrderCalculationHelper) ConvertToBitpinFormat(req *order.StandardOrderRequest) (map[string]interface{}, error) {
	if err := h.ValidateOrderRequestForBitpin(req); err != nil {
		return nil, err
	}

	baseAmount, err := h.GetBaseAmountForExchange(req)
	if err != nil {
		return nil, err
	}

	quoteAmount, err := h.GetQuoteAmountForExchange(req)
	if err != nil {
		return nil, err
	}

	// Convert symbol to Bitpin format (e.g., BTCUSDT -> BTC_USDT)

	orderData := map[string]interface{}{
		"symbol":       req.Symbol,
		"type":         strings.ToLower(string(req.Type)),
		"side":         strings.ToLower(string(req.Side)),
		"base_amount":  fmt.Sprintf("%.8f", baseAmount),
		"quote_amount": fmt.Sprintf("%.8f", quoteAmount),
	}

	if req.Type == order.OrderTypeLimit && req.Price != nil {
		orderData["price"] = fmt.Sprintf("%.8f", *req.Price)
	}

	if req.StopPrice != nil {
		orderData["stop_price"] = fmt.Sprintf("%.8f", *req.StopPrice)
	}

	if req.ClientOrderId != "" {
		orderData["identifier"] = req.ClientOrderId
	}

	return orderData, nil
}
