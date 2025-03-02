package service

import (
	"context"
	"fmt"

	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func signJwt(userID string, secretKey []byte) (string, error) {
	claims := domain.NewTokenClaims(userID)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) createTokens(ctx context.Context, userID uuid.UUID) (domain.Tokens, error) {
	refreshToken, err := s.tokens.Create(ctx, userID)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("failed to create refresh token: %w", err)
	}
	accessToken, err := signJwt(userID.String(), s.jwtSecret)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("failed to sign access token: %w", err)
	}
	return domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
