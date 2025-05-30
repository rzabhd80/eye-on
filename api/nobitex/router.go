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
	group := fiberRouter.Group("/Exchange/bitpin")
	group.Use(middleware.JWTAuthMiddleware(*router.Service.Exchange.UserRepo, router.Parser))
	group.Post("/order", router.Service.PlaceOrder)
	group.Delete("/order", router.Service.cancelOrder)
	group.Delete("/orderBook", router.Service.GetOrderBook)
	group.Get("/balance", router.Service.GetBalance)
}
