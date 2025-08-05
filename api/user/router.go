package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rzabhd80/eye-on/api/middleware"
	"github.com/rzabhd80/eye-on/internal/helpers"
)

type Router struct {
	Service *UserAuthService
	Parser  *helpers.JWTParser
}

func (router *Router) SetUserRouter(fiberRouter *fiber.App) {
	groupRouter := fiberRouter.Group("/user")
	groupRouter.Post("/register", router.Service.Register)
	groupRouter.Post("/login", router.Service.Login)
	groupRouter.Post("/exchangeCredentials", middleware.JWTAuthMiddleware(
		*router.Service.User.UserRepo, router.Parser), router.Service.CreateExchangeCredential)
	groupRouter.Put("/exchangeCredentials", middleware.JWTAuthMiddleware(
		*router.Service.User.UserRepo, router.Parser), router.Service.UpdateExchangeCredentials)
}
