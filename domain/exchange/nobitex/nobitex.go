package nobitex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/balance"
	"github.com/rzabhd80/eye-on/domain/exchange"
	"github.com/rzabhd80/eye-on/domain/exchangeCredentials"
	"github.com/rzabhd80/eye-on/domain/order"
	"github.com/rzabhd80/eye-on/domain/orderBook"
	"github.com/rzabhd80/eye-on/domain/traidingPair"
	"github.com/rzabhd80/eye-on/domain/user"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"github.com/rzabhd80/eye-on/internal/helpers"
	"net/http"
	"strconv"
	"time"
)

type NobitexExchange struct {
	NobitexExchangeModel   *models.Exchange
	ExchangeRepo           *exchange.ExchangeRepository
	ExchangeCredentialRepo *exchangeCredentials.ExchangeCredentialRepository
	UserRepo               *user.UserRepository
	TradingPairRepo        *traidingPair.TradingPairRepository
	OrderRepo              *order.OrderRepository
	OrderBookRepo          *orderBook.OrderBookSnapshotRepository
	BalanceRepo            *balance.BalanceSnapshotRepository
	Request                *helpers.Request
}

func (exchange *NobitexExchange) Name() string                   { return exchange.NobitexExchangeModel.Name }
func (exchange *NobitexExchange) Ping(ctx context.Context) error { return nil }
func (exchange *NobitexExchange) GetBalance(ctx context.Context, userId uuid.UUID, symbol *string) ([]models.BalanceSnapshot, error) {
	if symbol == nil {
		return nil, errors.New("Symbol cannot be null")
	}
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.NobitexExchangeModel.ID)
	if creds == nil {
		return nil, fmt.Errorf("credentials are required")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}
	request := exchange.Request
	symbolBody := map[string]string{
		"currency": *symbol,
	}
	marshaledBody, err := json.Marshal(symbolBody)
	if err != nil {
		return nil, err
	}
	respBody, body, err := request.MakeRequest(ctx, "POST", "/users/wallets/balance", marshaledBody, &models.ExchangeCredential{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
		IsTestnet: creds.IsTestnet,
	}, exchange.NobitexExchangeModel.BaseURL, false, true, helpers.ApiKeyAuth)
	if err != nil {
		return nil, err
	}
	if respBody.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("balance request failed: %s", body)
	}
	balanceResp := struct {
		Balance string `json:"balance"`
		Status  string `json:"status"`
	}{}
	if err := json.Unmarshal(body, &balanceResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	total, err := strconv.ParseFloat(balanceResp.Balance, 64)
	if err != nil {
		return nil, err
	}
	balanceSnapshot := []models.BalanceSnapshot{models.BalanceSnapshot{
		BaseModel:    models.BaseModel{ID: uuid.New()},
		UserID:       userId,
		ExchangeID:   exchange.NobitexExchangeModel.ID,
		Total:        total,
		Available:    total,
		SnapshotTime: time.Now(),
	}}
	return balanceSnapshot, nil
}
func (exchange *NobitexExchange) GetOrderBook(ctx context.Context, symbol string, userId uuid.UUID) (*models.OrderBookSnapshot, error) {
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.NobitexExchangeModel.ID)
	if creds == nil {
		return nil, fmt.Errorf("credentials are required")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}
	tradePair, err := exchange.TradingPairRepo.GetByExchangeAndSymbol(ctx, exchange.NobitexExchangeModel.ID, symbol)
	if tradePair == nil {
		return nil, fmt.Errorf("this symbol is not for this exchange ")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}

	request := exchange.Request
	endpoint := fmt.Sprintf("/v3/orderbook/%s", symbol)

	respBody, body, err := request.MakeRequest(ctx, "GET", endpoint, nil, nil,
		exchange.NobitexExchangeModel.BaseURL, false, false, helpers.ApiKeyAuth)
	if err != nil {
		return nil, err
	}
	orderBookResponse := struct {
		Status         string     `json:"status"`
		LastUpdate     int64      `json:"lastUpdate"`
		LastTradePrice string     `json:"lastTradePrice"`
		Asks           [][]string `json:"asks"`
		Bids           [][]string `json:"bids"`
	}{}
	if err := json.Unmarshal(body, &orderBookResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if respBody.StatusCode != http.StatusAccepted && respBody.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("order book request failed: %s", respBody.Body)
	}

	bids := make([]orderBook.StandardOrderLevel, 0, len(orderBookResponse.Bids))
	for _, bid := range orderBookResponse.Bids {
		if len(bid) >= 2 {
			price, _ := strconv.ParseFloat(bid[0], 64)
			quantity, _ := strconv.ParseFloat(bid[1], 64)
			bids = append(bids, orderBook.StandardOrderLevel{
				Price:    price,
				Quantity: quantity,
			})
		}
	}

	asks := make([]orderBook.StandardOrderLevel, 0, len(orderBookResponse.Asks))
	for _, ask := range orderBookResponse.Asks {
		if len(ask) >= 2 {
			price, _ := strconv.ParseFloat(ask[0], 64)
			quantity, _ := strconv.ParseFloat(ask[1], 64)
			asks = append(asks, orderBook.StandardOrderLevel{
				Price:    price,
				Quantity: quantity,
			})
		}
	}

	orderbookInstance := models.OrderBookSnapshot{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		ExchangeID:    exchange.NobitexExchangeModel.ID,
		TradingPairID: tradePair.ID,
		Symbol:        symbol,
		Bids: models.JSONB{
			"data": bids,
		},
		Asks: models.JSONB{
			"data": asks,
		},
		SnapshotTime: time.Now(),
	}
	err = exchange.OrderBookRepo.Create(ctx, &orderbookInstance)
	if err != nil {
		return nil, err
	}
	return &orderbookInstance, nil
}
func (exchange *NobitexExchange) PlaceOrder(ctx context.Context, req *order.StandardOrderRequest, userId uuid.UUID) (*models.OrderHistory, error) {
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.NobitexExchangeModel.ID)
	if creds == nil {
		return nil, fmt.Errorf("credentials are required")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}
	tradePair, err := exchange.TradingPairRepo.GetByExchangeAndSymbol(ctx, exchange.NobitexExchangeModel.ID, req.Symbol)
	if err != nil {
		return nil, errors.New("symbol not found for this exchange")
	}

	helper := &helpers.OrderCalculationHelper{}
	orderData, err := helper.ConvertToNobitexFormat(req)

	body, err := json.Marshal(orderData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	request := exchange.Request
	respBody, body, err := request.MakeRequest(ctx, "POST", "/market/orders/add", body, &models.ExchangeCredential{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
		IsTestnet: creds.IsTestnet,
	}, exchange.NobitexExchangeModel.BaseURL, false, true, helpers.ApiKeyAuth)
	if err != nil {
		return nil, err
	}

	exchangeOrderResponse := struct {
		Status string `json:"status"`
		Order  struct {
			Type            string    `json:"type"`
			SrcCurrency     string    `json:"srcCurrency"`
			DstCurrency     string    `json:"dstCurrency"`
			Price           string    `json:"price"`
			Amount          string    `json:"amount"`
			TotalPrice      string    `json:"totalPrice"`
			MatchedAmount   float64   `json:"matchedAmount"`
			UnmatchedAmount string    `json:"unmatchedAmount"`
			ID              int64     `json:"id"`
			Status          string    `json:"status"`
			Partial         bool      `json:"partial"`
			Fee             float64   `json:"fee"`
			CreatedAt       time.Time `json:"created_at"`
			ClientOrderID   string    `json:"clientOrderId"`
		} `json:"order"`
	}{}
	if err := json.Unmarshal(body, &exchangeOrderResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if respBody.StatusCode != http.StatusOK && respBody.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("order creation failed: %s", exchangeOrderResponse.Status)
	}

	// Convert to standard response

	// Map Bitpin status to standard status
	var status string
	switch exchangeOrderResponse.Order.Status {
	case "pending":
		status = "pending"
	case "partial":
		status = "partial"
	case "filled":
		status = "failed"
	case "cancelled":
		status = "cancelled"
	default:
		status = "new"
	}
	quantity, err := strconv.ParseFloat(exchangeOrderResponse.Order.Amount, 64)
	if err != nil {
		return nil, err
	}
	priceReturned, err := strconv.ParseFloat(exchangeOrderResponse.Order.Price, 64)
	if err != nil {
		return nil, err
	}
	orderHistory := models.OrderHistory{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		UserID:               userId,
		ExchangeCredentialID: creds.ID,
		ExchangeID:           exchange.NobitexExchangeModel.ID,
		TradingPairID:        tradePair.ID,
		ClientOrderID:        strconv.FormatInt(exchangeOrderResponse.Order.ID, 10) + userId.String(),
		ExchangeOrderID:      strconv.FormatInt(exchangeOrderResponse.Order.ID, 10),
		Side:                 exchangeOrderResponse.Order.Type,
		Type:                 "market",
		Quantity:             quantity,
		Price:                &priceReturned,
		Status:               status,
	}
	err = exchange.OrderRepo.Create(ctx, &orderHistory)
	if err != nil {
		return nil, err
	}

	return &orderHistory, nil
}

func (exchange *NobitexExchange) CancelOrder(ctx context.Context, orderID string, userId uuid.UUID) error {
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.NobitexExchangeModel.ID)
	if creds == nil {
		return fmt.Errorf("credentials are required")
	}
	if err != nil {
		return errors.New("Internal Server Error")
	}
	request := exchange.Request
	requestBody := map[string]string{
		"Order":  orderID,
		"Status": "Canceled",
	}
	requestBodyJson, err := json.Marshal(requestBody)
	respBody, body, err := request.MakeRequest(ctx, "POST", "/market/orders/update-status", requestBodyJson,
		&models.ExchangeCredential{
			APIKey:    creds.APIKey,
			SecretKey: creds.SecretKey,
			IsTestnet: creds.IsTestnet,
		}, exchange.NobitexExchangeModel.BaseURL, true, false, helpers.ApiKeyAuth)
	if err != nil {
		return err
	}
	if respBody.StatusCode != http.StatusOK && respBody.StatusCode != http.StatusAccepted {
		return fmt.Errorf("order cancellation failed: %s", respBody.Body)
	}

	var cancelResp struct {
		Status        string `json:"status"`
		UpdatedStatus string `json:"updatedStatus"`
	}
	if err := json.Unmarshal(body, &cancelResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
