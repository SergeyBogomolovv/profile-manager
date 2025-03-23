package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/SergeyBogomolovv/profile-manager/common/e"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type tokenRepo struct {
	db *redis.Client
}

func NewTokenRepo(db *redis.Client) *tokenRepo {
	return &tokenRepo{db: db}
}

func (r *tokenRepo) Create(ctx context.Context, userID string) (string, error) {
	token := uuid.NewString()
	return token, r.db.Set(ctx, tokenKey(token), userID, domain.TokenTTL).Err()
}

func (r *tokenRepo) CheckUserID(ctx context.Context, token string) (string, error) {
	userID, err := r.db.Get(ctx, tokenKey(token)).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrInvalidToken
	}
	if err != nil {
		return "", fmt.Errorf("failed to get user id: %w", err)
	}
	return userID, nil
}

func (r *tokenRepo) Revoke(ctx context.Context, token string) error {
	return e.WrapIfErr(r.db.Del(ctx, tokenKey(token)).Err(), "failed to revoke token")
}

func tokenKey(token string) string {
	return fmt.Sprintf("telegram_token:%s", token)
}
