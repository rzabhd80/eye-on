package bitpin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/balance"
	"github.com/rzabhd80/eye-on/domain/exchange/bitpin"
	"github.com/rzabhd80/eye-on/domain/order"
	"github.com/rzabhd80/eye-on/domain/orderBook"
	"strings"
	"time"
)

type BitpinService struct {
	Exchange *bitpin.BitpinExchange
}

func (service *BitpinService) GetBalance(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uuid.UUID)
	var request balance.GetBalanceRequest
	if err := c.QueryParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(bitpin.ErrorResponse{Error: "Bad Request Format"})
	}
	balanceSnapshots, err := service.Exchange.GetBalance(c.Context(), userId, &request.Asset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(bitpin.ErrorResponse{Error: err.Error()})
	}

	balances := make([]balance.StandardBalanceResponse, 0, len(balanceSnapshots))
	for _, balanceIns := range balanceSnapshots {
		available := balanceIns.Available
		total := balanceIns.Total
		frozen := total - available

		balances = append(balances, balance.StandardBalanceResponse{
			Asset:  strings.ToUpper(balanceIns.Currency),
			Free:   available,
			Locked: frozen,
			Total:  total,
		})
	}
	return c.Status(fiber.StatusOK).JSON(balances)
}

func (service *BitpinService) GetOrderBook(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uuid.UUID)
	var request orderBook.StandardOrderBookRequest
	if err := c.ParamsParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(bitpin.ErrorResponse{Error: "HINT:Bitpin access token expires every " +
			"15 min. Refresh itBad Request Format"})
	}
	orderBookHistory, err := service.Exchange.GetOrderBook(c.Context(), request.Symbol, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(bitpin.ErrorResponse{Error: err.Error()})
	}
	history := orderBook.StandardOrderBookResponse{
		Symbol:    orderBookHistory.Symbol,
		Bids:      orderBookHistory.Bids,
		Asks:      orderBookHistory.Asks,
		Timestamp: time.Now().Format(time.RFC850),
	}
	return c.Status(fiber.StatusOK).JSON(history)
}

func (service *BitpinService) PlaceOrder(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uuid.UUID)
	var request order.StandardOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(bitpin.ErrorResponse{Error: "HINT:Bitpin access token expires every" +
			" 15 min. Refresh itBad Request Format"})
	}
	orderHistory, err := service.Exchange.PlaceOrder(c.Context(), &request, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(bitpin.ErrorResponse{Error: err.Error()})
	}
	response := order.StandardOrderResponse{
		ID:         orderHistory.ID.String(),
		Symbol:     orderHistory.TradingPair.Symbol,
		Side:       order.OrderSide(orderHistory.Side),
		Type:       order.OrderType(orderHistory.Type),
		Quantity:   orderHistory.Quantity,
		Price:      orderHistory.Price,
		Status:     order.OrderStatus(orderHistory.Status),
		CreatedAt:  orderHistory.CreatedAt,
		UpdatedAt:  orderHistory.UpdatedAt,
		ExchangeID: orderHistory.ExchangeID.String(),
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (service *BitpinService) cancelOrder(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uuid.UUID)
	var request order.CancelOrderRequest
	if err := c.ParamsParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(bitpin.ErrorResponse{Error: "Bad Request Format missing orderId as url param"})
	}

	resultErr := service.Exchange.CancelOrder(c.Context(), &request.OrderId, userId, nil)
	if resultErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(bitpin.ErrorResponse{Error: resultErr.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(map[string]string{"message": "success"})
}
