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

func (router *Router) SetUserRouter(fiberRouter fiber.Router) {
	fiberRouter.Group("/exchange/bitpin")
	fiberRouter.Use(middleware.JWTAuthMiddleware(*router.Service.exchange.UserRepo, *router.Parser))
	fiberRouter.Post("/order", router.Service.PlaceOrder)
	fiberRouter.Delete("/order", router.Service.cancelOrder)
	fiberRouter.Delete("/orderBook", router.Service.GetOrderBook)
	fiberRouter.Get("/balance", router.Service.GetBalance)
}
