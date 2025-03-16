package service

import (
	"context"
	"fmt"

	"github.com/SergeyBogomolovv/profile-manager/common/auth"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/google/uuid"
)

func (s *authService) createTokens(ctx context.Context, userID uuid.UUID) (domain.Tokens, error) {
	refreshToken, err := s.tokens.Create(ctx, userID)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("failed to create refresh token: %w", err)
	}
	accessToken, err := s.signJwt(userID)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("failed to sign access token: %w", err)
	}
	return domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

const issuer = "sso"

func (s *authService) signJwt(userID uuid.UUID) (string, error) {
	return auth.SignJWT(userID.String(), s.jwtSecret, domain.AccessTokenTTL, issuer)
}
