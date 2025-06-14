package user

import (
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"time"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type ExchangeCredentialRequest struct {
	ExchangeName string `json:"exchange_name" validate:"required"`
	Label        string `json:"label" validate:"required"`
	APIKey       string `json:"api_key" validate:"required"`
	SecretKey    string `json:"secret_key" validate:"required"`
	AccessKey    string `json:"access_key,omitempty"`
	IsTestnet    bool   `json:"is_testnet"`
}

type ExchangeCredentialUpdateRequest struct {
	ExchangeName string `json:"exchange_name" validate:"required"`
	Label        string `json:"label" validate:"required"`
	APIKey       string `json:"api_key" validate:"required"`
	SecretKey    string `json:"secret_key,omitempty"`
	AccessKey    string `json:"access_key,omitempty"`
	RefreshToken string `json:"refresh_key,omitempty"`
	IsActive     string `json:"is_active,omitempty"`
	IsTestnet    bool   `json:"is_testnet"`
}

type ExchangeCredentialResponse struct {
	ID         uuid.UUID       `json:"id"`
	ExchangeID uuid.UUID       `json:"exchange_id"`
	Label      string          `json:"label"`
	APIKey     string          `json:"api_key"`
	IsActive   bool            `json:"is_active"`
	IsTestnet  bool            `json:"is_testnet"`
	AccessKey  *string         `json:"access_key,omitempty"`
	LastUsed   *time.Time      `json:"last_used,omitempty"`
	Exchange   models.Exchange `json:"exchange,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
