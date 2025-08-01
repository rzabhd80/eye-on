package user

import (
	"context"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/domain/exchange"
	"github.com/rzabhd80/eye-on/domain/exchangeCredentials"
	"github.com/rzabhd80/eye-on/internal/database/models"
	envCofig "github.com/rzabhd80/eye-on/internal/envConfig"
	"github.com/rzabhd80/eye-on/internal/helpers"
	"strconv"
)

type User struct {
	UserRepo         *UserRepository
	ExchangeRepo     *exchange.ExchangeRepository
	ExchangeCredRepo *exchangeCredentials.ExchangeCredentialRepository
	JwtParser        *helpers.JWTParser
	EnvConf          *envCofig.AppConfig
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
	token, err := user.JwtParser.GenerateJWT(&createdUser)
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

	token, err := user.JwtParser.GenerateJWT(userByUsername)
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
	if err == nil && existingCredentials != nil {
		return nil, &ErrorResponse{Error: "Exchange Credentials Already Exists "}
	}
	ecryptionKey := user.EnvConf.EncryptionKey
	encryptedApiKey, err := helpers.EncryptAPIKey(request.APIKey, ecryptionKey)
	if err != nil {
		return nil, &ErrorResponse{Error: "Internal Server Error "}
	}
	encryptedSecretKey := hex.EncodeToString([]byte(request.SecretKey))
	encryptedAccKey := ""
	if request.AccessKey != "" {
		encryptedAccKey, err = helpers.EncryptAPIKey(request.AccessKey, ecryptionKey)
		if err != nil {
			return nil, &ErrorResponse{Error: "Internal Server Error "}
		}
	}
	encyptedRefreshKey := ""
	if request.RefreshKey != "" {
		encyptedRefreshKey, err = helpers.EncryptAPIKey(request.RefreshKey, ecryptionKey)
		if err != nil {
			return nil, &ErrorResponse{Error: "Internal Server Error "}
		}
	}
	credential := models.ExchangeCredential{
		UserID:     userId,
		ExchangeID: exchangeReg.ID,
		Label:      request.Label,
		APIKey:     encryptedApiKey,
		SecretKey:  encryptedSecretKey,
		IsActive:   true,
		IsTestnet:  request.IsTestnet,
	}
	if request.AccessKey != "" {
		credential.AccessKey = encryptedAccKey
	}
	if request.RefreshKey != "" {
		credential.RefreshKey = encyptedRefreshKey
	}
	err = user.ExchangeCredRepo.Update(ctx, &credential)
	if err != nil {
		return nil, &ErrorResponse{Error: "Internal Server Error"}
	}
	response := &ExchangeCredentialResponse{
		ID:         credential.ID,
		ExchangeID: exchangeReg.ID,
		Label:      request.Label,
		APIKey:     request.APIKey,
		IsActive:   true,
		IsTestnet:  false,
		LastUsed:   nil,
		Exchange:   *exchangeReg,
	}
	if request.RefreshKey != "" {
		response.RefreshKey = &request.RefreshKey
	}
	if request.AccessKey != "" {
		response.AccessKey = &request.AccessKey
	}
	return response, nil
}

func (user *User) UpdateExchangeCredential(ctx context.Context, request ExchangeCredentialUpdateRequest, userId uuid.UUID) (
	*ExchangeCredentialResponse, *ErrorResponse) {

	exchangeReg, err := user.ExchangeRepo.GetByName(ctx, request.ExchangeName)
	if err != nil || exchangeReg == nil {
		return nil, &ErrorResponse{Error: "Exchange Not Found "}
	}
	existingCredentials, err := user.ExchangeCredRepo.GetByUserAndExchange(ctx, userId, exchangeReg.ID)
	if err != nil || existingCredentials == nil {
		return nil, &ErrorResponse{Error: "Exchange Credentials Not Found "}
	}
	ecryptionKey := user.EnvConf.EncryptionKey
	encryptedApiKey, err := helpers.EncryptAPIKey(request.APIKey, ecryptionKey)
	if err != nil {
		return nil, &ErrorResponse{Error: "Internal Server Error "}
	}
	var encryptedSecretKey string
	if request.SecretKey != "" {
		encryptedSecretKey = hex.EncodeToString([]byte(request.SecretKey))
	}

	encryptedAccKey := ""
	if request.AccessKey != "" {
		encryptedAccKey, err = helpers.EncryptAPIKey(request.AccessKey, ecryptionKey)
		if err != nil {
			return nil, &ErrorResponse{Error: "Internal Server Error "}
		}
	}
	var active bool
	if request.IsActive != "" {
		active, err = strconv.ParseBool(request.IsActive)
	}
	var refreshTokenEnc string
	if request.RefreshToken != "" {
		refreshTokenEnc, err = helpers.EncryptAPIKey(request.RefreshToken, ecryptionKey)
		if err != nil {
			return nil, &ErrorResponse{Error: err.Error()}
		}
	}
	if request.APIKey != "" {
		existingCredentials.APIKey = encryptedApiKey
	}
	if request.AccessKey != "" {
		existingCredentials.AccessKey = encryptedAccKey
	}
	if request.IsActive != "" {
		existingCredentials.IsActive = active
	}
	if request.SecretKey != "" {
		existingCredentials.SecretKey = encryptedSecretKey
	}
	if request.RefreshToken != "" {
		existingCredentials.RefreshKey = refreshTokenEnc
	}
	err = user.ExchangeCredRepo.Update(ctx, existingCredentials)
	if err != nil {
		return nil, &ErrorResponse{Error: "Internal Server Error"}
	}

	response := &ExchangeCredentialResponse{
		ID:         existingCredentials.ID,
		ExchangeID: exchangeReg.ID,
		Label:      request.Label,
		APIKey:     request.APIKey,
		IsActive:   existingCredentials.IsActive,
		IsTestnet:  false,
		LastUsed:   nil,
		Exchange:   *exchangeReg,
	}
	if request.AccessKey != "" {
		response.AccessKey = &request.AccessKey
	}
	if request.RefreshToken != "" {
		response.RefreshKey = &request.RefreshToken
	}
	return response, nil
}
