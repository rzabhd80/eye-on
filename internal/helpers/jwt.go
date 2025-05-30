package helpers

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rzabhd80/eye-on/internal/database/models"
	envCofig "github.com/rzabhd80/eye-on/internal/envConfig"
	"time"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

type JWTParser struct {
	EnvConf *envCofig.AppConfig
}

func (jwtParser *JWTParser) GenerateSecureKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}

func (jwtParser *JWTParser) GenerateJWT(userInstance *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userInstance.ID,
		Username: userInstance.Username,
		Email:    userInstance.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtParser.EnvConf.JWTKey))
}

func (jwtParser *JWTParser) ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtParser.EnvConf.JWTKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
