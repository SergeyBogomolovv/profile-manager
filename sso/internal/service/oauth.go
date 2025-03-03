package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/google/uuid"
)

func (s *authService) OAuth(ctx context.Context, info domain.OAuthUserInfo, provider domain.AccountType) (domain.Tokens, error) {
	var user domain.User
	var account domain.Account

	err := s.txManager.Run(ctx, func(ctx context.Context) (err error) {
		user, err = s.ensureUser(ctx, info.Email)
		if err != nil {
			return fmt.Errorf("failed to ensure user: %w", err)
		}
		account, err = s.ensureAccount(ctx, user.ID, provider)
		if err != nil {
			return fmt.Errorf("failed to ensure account: %w", err)
		}
		return nil
	})
	if err != nil {
		return domain.Tokens{}, err
	}

	if account.UserID != user.ID {
		return domain.Tokens{}, domain.ErrInvalidCredentials
	}

	return s.createTokens(ctx, user.ID)
}

func (s *authService) ensureUser(ctx context.Context, email string) (domain.User, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return user, nil
	}

	if errors.Is(err, domain.ErrUserNotFound) {
		return s.users.Create(ctx, email)
	}

	return domain.User{}, err
}

func (s *authService) ensureAccount(ctx context.Context, userID uuid.UUID, provider domain.AccountType) (domain.Account, error) {
	account, err := s.users.AccountByID(ctx, userID, provider)
	if err == nil {
		return account, nil
	}

	if errors.Is(err, domain.ErrAccountNotFound) {
		return s.users.AddAccount(ctx, userID, provider, nil)
	}

	return domain.Account{}, err
}
