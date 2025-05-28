package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/user"
)

type UserAuthService struct {
	User *user.User
}

func (service *UserAuthService) Register(c *fiber.Ctx) error {
	var requestBody user.RegisterRequest = user.RegisterRequest{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(user.ErrorResponse{Error: "Bad Request Format"})
	}
	response, err := service.User.Register(c.Context(), requestBody)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (service *UserAuthService) Login(c *fiber.Ctx) error {
	var requestBody user.LoginRequest = user.LoginRequest{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(user.ErrorResponse{Error: "Bad Request Format"})
	}
	response, err := service.User.Login(c.Context(), requestBody)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (service *UserAuthService) CreateExchangeCredential(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uuid.UUID)
	var requestBody user.ExchangeCredentialRequest = user.ExchangeCredentialRequest{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(user.ErrorResponse{Error: "Bad Request Format"})
	}
	response, err := service.User.CreateExchangeCredential(c.Context(), requestBody, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
