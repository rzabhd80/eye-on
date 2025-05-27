package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rzabhd80/eye-on/domain/user"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"github.com/rzabhd80/eye-on/internal/helpers"
	"strings"
)

func JWTAuthMiddleware(userRepo user.UserRepository, jwtParser helpers.JWTParser) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(user.ErrorResponse{
				Error: "Authorization header required",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(user.ErrorResponse{
				Error: "Bearer token required",
			})
		}
		parser := helpers.JWTParser{}
		claims, err := parser.ParseJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(user.ErrorResponse{
				Error: "Invalid token",
			})
		}
		var foundUser models.User
		if _, err := userRepo.GetByID(c.Context(), claims.UserID); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(user.ErrorResponse{
				Error: "User not found or inactive",
			})
		}

		c.Locals("user", foundUser)
		c.Locals("user_id", foundUser.ID)
		return c.Next()
	}
}
