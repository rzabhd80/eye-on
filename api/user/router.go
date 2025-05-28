package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rzabhd80/eye-on/api/middleware"
	"github.com/rzabhd80/eye-on/internal/helpers"
)

type Router struct {
	service *UserAuthService
	parser  *helpers.JWTParser
}

func (router *Router) userRouter(fiberRouter fiber.Router) {
	fiberRouter.Group("/User")
	fiberRouter.Post("/register", router.service.Register)
	fiberRouter.Post("/login", router.service.Login)
	fiberRouter.Post("/exchangeCredentials", middleware.JWTAuthMiddleware(
		*router.service.User.UserRepo, *router.parser), router.service.CreateExchangeCredential)
}
