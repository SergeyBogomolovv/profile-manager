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
			return err
		}
		account, err = s.ensureAccount(ctx, user.ID, provider)
		if err != nil {
			return err
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

	if errors.Is(err, domain.ErrUserNotFound) {
		user, err = s.users.Create(ctx, email)
		if err != nil {
			return domain.User{}, fmt.Errorf("failed to create user: %w", err)
		}
		return user, nil
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *authService) ensureAccount(ctx context.Context, userID uuid.UUID, provider domain.AccountType) (domain.Account, error) {
	account, err := s.users.AccountByID(ctx, userID, provider)

	if errors.Is(err, domain.ErrAccountNotFound) {
		account, err := s.users.AddAccount(ctx, userID, provider, nil)
		if err != nil {
			return domain.Account{}, fmt.Errorf("failed to add account: %w", err)
		}
		// TODO: send data to rabbitmq
		return account, nil
	}

	if err != nil {
		return domain.Account{}, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}
