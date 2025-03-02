package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type tokensRepo struct {
	db *redis.Client
}

func NewTokensRepo(db *redis.Client) *tokensRepo {
	return &tokensRepo{db: db}
}

func (r *tokensRepo) Create(ctx context.Context, userID uuid.UUID) (string, error) {
	payload := RefreshToken{
		UserID:    userID,
		ExpiresAt: time.Now().Add(domain.RefreshTokenTTL),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal refresh token: %w", err)
	}

	token := uuid.New().String()
	if err := r.db.Set(ctx, token, data, domain.RefreshTokenTTL).Err(); err != nil {
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}
	return token, nil
}

func (r *tokensRepo) UserID(ctx context.Context, token string) (uuid.UUID, error) {
	data, err := r.db.Get(ctx, token).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return uuid.Nil, domain.ErrInvalidToken
		}
		return uuid.Nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	var payload RefreshToken
	if err := json.Unmarshal(data, &payload); err != nil {
		return uuid.Nil, fmt.Errorf("failed to unmarshal refresh token: %w", err)
	}
	if payload.ExpiresAt.Before(time.Now()) {
		return uuid.Nil, domain.ErrInvalidToken
	}
	return payload.UserID, nil
}
