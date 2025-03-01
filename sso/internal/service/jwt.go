package service

import (
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func signJwt(userID string, secretKey []byte) (string, error) {
	claims := domain.NewTokenClaims(userID)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}
