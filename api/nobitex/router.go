package nobitex

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rzabhd80/eye-on/api/middleware"
	"github.com/rzabhd80/eye-on/internal/helpers"
)

type Router struct {
	Service *NobitexService
	Parser  *helpers.JWTParser
}

func (router *Router) SetUserRouter(fiberRouter *fiber.App) {
	group := fiberRouter.Group("/exchange/nobitex")
	group.Use(middleware.JWTAuthMiddleware(*router.Service.Exchange.UserRepo, router.Parser))
	group.Post("/order", router.Service.PlaceOrder)
	group.Delete("/order/:orderId", router.Service.cancelOrder)
	group.Get("/orderBook/:symbol", router.Service.GetOrderBook)
	group.Get("/balance/:symbol", router.Service.GetBalance)
}
