package bitpin

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
	envCofig "github.com/rzabhd80/eye-on/internal/envConfig"
	"github.com/rzabhd80/eye-on/internal/helpers"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BitpinExchange struct {
	BitpinExchangeModel    *models.Exchange
	ExchangeRepo           *exchange.ExchangeRepository
	ExchangeCredentialRepo *exchangeCredentials.ExchangeCredentialRepository
	UserRepo               *user.UserRepository
	TradingPairRepo        *traidingPair.TradingPairRepository
	OrderRepo              *order.OrderRepository
	OrderBookRepo          *orderBook.OrderBookSnapshotRepository
	BalanceRepo            *balance.BalanceSnapshotRepository
	Request                *helpers.Request
	EnvConf                *envCofig.AppConfig
}

func (exchange *BitpinExchange) Name() string                   { return exchange.BitpinExchangeModel.Name }
func (exchange *BitpinExchange) Ping(ctx context.Context) error { return nil }
func (exchange *BitpinExchange) GetBalance(ctx context.Context, userId uuid.UUID, sign *string) ([]models.BalanceSnapshot, error) {
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.BitpinExchangeModel.ID)
	if creds == nil {
		return nil, fmt.Errorf("credentials are required")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}
	request := exchange.Request

	respBody, body, err := request.MakeRequest(ctx, "GET", "/api/v1/wlt/wallets/", nil, &models.ExchangeCredential{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
		AccessKey: creds.AccessKey,
		IsTestnet: creds.IsTestnet,
	}, exchange.BitpinExchangeModel.BaseURL, true, false, helpers.ApiAccToken)
	if err != nil {
		return nil, err
	}

	var balanceResp []struct {
		ID      int    `json:"id"`
		Asset   string `json:"asset"`
		Balance string `json:"balance"`
		Frozen  string `json:"frozen"`
		Service string `json:"service"`
	}

	if respBody.StatusCode != http.StatusOK && respBody.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("API error. Exchange %s: said: status %d, body: %s", exchange.Name(),
			respBody.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &balanceResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	balances := make([]balance.StandardBalanceResponse, 0, len(balanceResp))
	balanceSnapshot := make([]models.BalanceSnapshot, 0, len(balanceResp))
	for _, balanceIns := range balanceResp {

		frozen, _ := strconv.ParseFloat(balanceIns.Frozen, 64)
		total, _ := strconv.ParseFloat(balanceIns.Balance, 64)
		available := total - frozen

		balances = append(balances, balance.StandardBalanceResponse{
			Asset:  strings.ToUpper(balanceIns.Asset),
			Free:   available,
			Locked: frozen,
			Total:  total,
		})
		balanceSnapshot = append(balanceSnapshot, models.BalanceSnapshot{
			BaseModel:    models.BaseModel{ID: uuid.New()},
			UserID:       userId,
			ExchangeID:   exchange.BitpinExchangeModel.ID,
			Total:        total,
			Available:    available,
			SnapshotTime: time.Now(),
		})
	}

	return balanceSnapshot, nil
}
func (exchange *BitpinExchange) GetOrderBook(ctx context.Context, symbol string, userId uuid.UUID) (*models.OrderBookSnapshot, error) {
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.BitpinExchangeModel.ID)
	if creds == nil {
		return nil, fmt.Errorf("credentials are required")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}
	tradePair, err := exchange.TradingPairRepo.GetByExchangeAndSymbol(ctx, exchange.BitpinExchangeModel.ID, symbol)
	if tradePair == nil {
		return nil, fmt.Errorf("this symbol is not for this exchange ")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}

	request := exchange.Request
	endpoint := fmt.Sprintf("/api/v1/mth/orderbook/%s/", symbol)

	respBody, body, err := request.MakeRequest(ctx, "GET", endpoint, nil, nil,
		exchange.BitpinExchangeModel.BaseURL, false, false, helpers.ApiAccToken)
	if err != nil {
		return nil, err
	}
	orderBookResponse := struct {
		Bids [][]string `json:"bids"`
		Asks [][]string `json:"asks"`
	}{}
	if err := json.Unmarshal(body, &orderBookResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	if respBody.StatusCode != http.StatusOK && respBody.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("API error. Exchange %s said: status %d, body: %s", exchange.Name(),
			respBody.StatusCode, string(body))
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
		ExchangeID:    exchange.BitpinExchangeModel.ID,
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

// RenewAccessToken Renews Bitpin access token
func (exchange *BitpinExchange) RenewAccessToken(ctx context.Context, userId uuid.UUID) (
	*models.ExchangeCredential, error) {
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.BitpinExchangeModel.ID)
	if creds == nil {
		return nil, fmt.Errorf("credentials are required")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}
	var body map[string]interface{} = map[string]interface{}{"refresh": creds.RefreshKey}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	respBody, pureBody, err := exchange.Request.MakeRequest(ctx, "POST", "/api/v1/usr/refresh_token/",
		jsonBody, creds, exchange.BitpinExchangeModel.BaseURL, false, false, helpers.ApiRefreshToken)
	if err != nil {
		return nil, err
	}
	if respBody.StatusCode != http.StatusOK && respBody.StatusCode != http.StatusAccepted &&
		respBody.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error. Exchange %s said: status %d, body: %s", exchange.Name(),
			respBody.StatusCode, string(pureBody))

	}
	expectedResponse := struct {
		Access string `json:"access"`
	}{}
	if err := json.Unmarshal(pureBody, &expectedResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	creds.AccessKey, err = helpers.EncryptAPIKey(expectedResponse.Access, exchange.EnvConf.EncryptionKey)
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}

	creds.RefreshKey, err = helpers.EncryptAPIKey(creds.RefreshKey, exchange.EnvConf.EncryptionKey)
	if err != nil {
		return nil, err
	}
	
	creds.APIKey, err = helpers.EncryptAPIKey(creds.APIKey, exchange.EnvConf.EncryptionKey)
	if err != nil {
		return nil, err
	}

	updateErr := exchange.ExchangeCredentialRepo.Update(ctx, creds)
	if updateErr != nil {
		return nil, updateErr
	}
	return creds, nil
}

func (exchange *BitpinExchange) PlaceOrder(ctx context.Context, req *order.StandardOrderRequest, userId uuid.UUID) (*models.OrderHistory, error) {
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.BitpinExchangeModel.ID)
	if creds == nil {
		return nil, fmt.Errorf("credentials are required")
	}
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}
	helper := &helpers.OrderCalculationHelper{}
	orderData, err := helper.ConvertToBitpinFormat(req)

	tradePair, err := exchange.TradingPairRepo.GetByExchangeAndSymbol(ctx, exchange.BitpinExchangeModel.ID, req.Symbol)
	if err != nil {
		return nil, errors.New("symbol not found for this exchange")
	}
	body, err := json.Marshal(orderData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	request := exchange.Request
	respBody, body, err := request.MakeRequest(ctx, "POST", "/api/v1/odr/orders/", body, &models.ExchangeCredential{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
		AccessKey: creds.AccessKey,
		IsTestnet: creds.IsTestnet,
	}, exchange.BitpinExchangeModel.BaseURL, true, false, helpers.ApiAccToken)
	if err != nil {
		return nil, err
	}

	exchangeOrderResponse := struct {
		ID                int64      `json:"id"`
		Symbol            string     `json:"symbol"`
		Type              string     `json:"type"`
		Side              string     `json:"side"`
		Price             string     `json:"price"`
		StopPrice         *string    `json:"stop_price"`
		OCOTargetPrice    *string    `json:"oco_target_price"`
		BaseAmount        string     `json:"base_amount"`
		QuoteAmount       string     `json:"quote_amount"`
		Identifier        *string    `json:"identifier"`
		State             string     `json:"state"`
		ClosedAt          *time.Time `json:"closed_at"`
		CreatedAt         time.Time  `json:"created_at"`
		DealedBaseAmount  string     `json:"dealed_base_amount"`
		DealedQuoteAmount string     `json:"dealed_quote_amount"`
		ReqToCancel       bool       `json:"req_to_cancel"`
		Commission        string     `json:"commission"`
	}{}
	if respBody.StatusCode != http.StatusOK && respBody.StatusCode != http.StatusAccepted &&
		respBody.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error. Exchange %s said: status %d, body: %s", exchange.Name(),
			respBody.StatusCode, string(body))

	}
	if err := json.Unmarshal(body, &exchangeOrderResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert to standard response

	// Map Bitpin status to standard status
	var status string
	switch exchangeOrderResponse.State {
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
	quantity, err := strconv.ParseFloat(exchangeOrderResponse.QuoteAmount, 64)
	if err != nil {
		return nil, err
	}
	priceReturned, err := strconv.ParseFloat(exchangeOrderResponse.Price, 64)
	if err != nil {
		return nil, err
	}
	orderHistory := models.OrderHistory{
		BaseModel: models.BaseModel{
			ID: uuid.New(),
		},
		UserID:               userId,
		ExchangeCredentialID: creds.ID,
		ExchangeID:           exchange.BitpinExchangeModel.ID,
		TradingPairID:        tradePair.ID,
		ClientOrderID:        strconv.FormatInt(exchangeOrderResponse.ID, 10) + userId.String(),
		ExchangeOrderID:      strconv.FormatInt(exchangeOrderResponse.ID, 10),
		Side:                 exchangeOrderResponse.Side,
		Type:                 exchangeOrderResponse.Type,
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

func (exchange *BitpinExchange) CancelOrder(ctx context.Context, orderID *string, userId uuid.UUID, hours *float64) error {
	creds, err := exchange.ExchangeCredentialRepo.GetByUserAndExchange(ctx, userId, exchange.BitpinExchangeModel.ID)
	orderId, err := uuid.Parse(*orderID)
	if err != nil {
		return errors.New("malformed orderId")
	}
	orderData, err := exchange.OrderRepo.GetByID(ctx, orderId)
	if err != nil {
		return errors.New("order record was not found")
	}

	if creds == nil {
		return fmt.Errorf("credentials are required")
	}
	if err != nil {
		return errors.New("Internal Server Error")
	}
	request := exchange.Request
	endpoint := fmt.Sprintf("/api/v1/odr/orders/%s/", orderData.ExchangeOrderID)
	respBody, body, err := request.MakeRequest(ctx, "DELETE", endpoint, nil, &models.ExchangeCredential{
		APIKey:    creds.APIKey,
		SecretKey: creds.SecretKey,
		AccessKey: creds.AccessKey,
		IsTestnet: creds.IsTestnet,
	}, exchange.BitpinExchangeModel.BaseURL, true, false, helpers.ApiAccToken)
	if err != nil {
		return err
	}

	if respBody.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API error. Exchange %s said: status %d, body: %s", exchange.Name(),
			respBody.StatusCode, string(body))
	}

	return nil
}
