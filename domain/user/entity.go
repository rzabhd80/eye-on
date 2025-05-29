package user

import (
	"context"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/exchange"
	"github.com/rzabhd80/eye-on/domain/exchangeCredentials"
	"github.com/rzabhd80/eye-on/internal/database/models"
	"github.com/rzabhd80/eye-on/internal/helpers"
)

type User struct {
	UserRepo         *UserRepository
	ExchangeRepo     *exchange.ExchangeRepository
	ExchangeCredRepo *exchangeCredentials.ExchangeCredentialRepository
	jwtParser        *helpers.JWTParser
}

func (user *User) Register(ctx context.Context, request RegisterRequest) (*AuthResponse, *ErrorResponse) {
	if userWithEmail, err := user.UserRepo.GetByEmail(ctx, request.Email); err == nil && userWithEmail != nil {
		return nil, &ErrorResponse{Error: "email already exists"}
	}
	if userWithUsername, err := user.UserRepo.GetByUsername(ctx, request.Username); err == nil && userWithUsername != nil {
		return nil, &ErrorResponse{Error: "Username already exists"}
	}
	hashedPassword, err := helpers.HashPassword(request.Password)
	if err != nil {
		return nil, &ErrorResponse{Error: "internal server error"}
	}
	createdUser := models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
		IsActive: true,
	}
	err = user.UserRepo.Create(ctx, &createdUser)
	if err != nil {
		return nil, &ErrorResponse{Error: "internal server error"}
	}
	token, err := user.jwtParser.GenerateJWT(&createdUser)
	if err != nil {
		return nil, &ErrorResponse{Error: "internal server error"}
	}

	return &AuthResponse{
		Token: token,
		User:  createdUser,
	}, nil
}

func (user *User) Login(ctx context.Context, request LoginRequest) (*AuthResponse, *ErrorResponse) {
	userByUsername, err := user.UserRepo.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, &ErrorResponse{Error: "invalid username"}
	}

	if err := helpers.VerifyHashedPassword(userByUsername.Password, request.Password); err != nil {
		return nil, &ErrorResponse{Error: "invalid credentials"}
	}

	token, err := user.jwtParser.GenerateJWT(userByUsername)
	if err != nil {
		return nil, &ErrorResponse{Error: "internal server error"}
	}

	return &AuthResponse{
		Token: token,
		User:  *userByUsername,
	}, nil

}

func (user *User) CreateExchangeCredential(ctx context.Context, request ExchangeCredentialRequest, userId uuid.UUID) (
	*ExchangeCredentialResponse, *ErrorResponse) {

	exchangeReg, err := user.ExchangeRepo.GetByName(ctx, request.ExchangeName)
	if err != nil || exchangeReg == nil {
		return nil, &ErrorResponse{Error: "Exchange Not Found "}
	}
	existingCredentials, err := user.ExchangeCredRepo.GetByUserAndExchange(ctx, userId, exchangeReg.ID)
	if err != nil || existingCredentials != nil {
		return nil, &ErrorResponse{Error: "Exchange Credentials Already Exists "}
	}
	encryptedSecretKey := hex.EncodeToString([]byte(request.SecretKey))
	encryptedPassphrase := ""
	if request.AccessKey != "" {
		encryptedPassphrase = hex.EncodeToString([]byte(request.AccessKey))
	}
	credential := models.ExchangeCredential{
		UserID:     userId,
		ExchangeID: exchangeReg.ID,
		Label:      request.Label,
		APIKey:     request.APIKey,
		SecretKey:  encryptedSecretKey,
		AccessKey:  encryptedPassphrase,
		IsActive:   true,
		IsTestnet:  request.IsTestnet,
	}
	err = user.ExchangeCredRepo.Create(ctx, &credential)
	if err != nil {
		return nil, &ErrorResponse{Error: "Internal Server Error"}
	}
	return &ExchangeCredentialResponse{
		ID:         credential.ID,
		ExchangeID: exchangeReg.ID,
		Label:      request.Label,
		APIKey:     request.APIKey,
		IsActive:   true,
		IsTestnet:  false,
		LastUsed:   nil,
		Exchange:   *exchangeReg,
	}, nil
}
