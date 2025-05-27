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
	ExchangeName string       `json:"exchange_name" validate:"required"`
	Label        string       `json:"label" validate:"required"`
	APIKey       string       `json:"api_key" validate:"required"`
	SecretKey    string       `json:"secret_key" validate:"required"`
	Passphrase   string       `json:"passphrase,omitempty"`
	IsTestnet    bool         `json:"is_testnet"`
	Permissions  models.JSONB `json:"permissions,omitempty"`
}

type ExchangeCredentialResponse struct {
	ID          uuid.UUID       `json:"id"`
	ExchangeID  uuid.UUID       `json:"exchange_id"`
	Label       string          `json:"label"`
	APIKey      string          `json:"api_key"`
	IsActive    bool            `json:"is_active"`
	IsTestnet   bool            `json:"is_testnet"`
	Permissions models.JSONB    `json:"permissions"`
	LastUsed    *time.Time      `json:"last_used,omitempty"`
	Exchange    models.Exchange `json:"exchange,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
