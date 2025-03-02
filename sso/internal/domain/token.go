package domain

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

const AccessTokenTTL = time.Hour
const RefreshTokenTTL = time.Hour * 24 * 7

func NewTokenClaims(userID string) TokenClaims {
	return TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "sso",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		},
	}
}
