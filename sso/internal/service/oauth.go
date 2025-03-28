package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/google/uuid"
)

func (s *authService) OAuth(ctx context.Context, info domain.OAuthUserInfo, provider domain.AccountType, ip string) (domain.Tokens, error) {
	var user domain.User
	var account domain.Account

	err := s.txManager.Run(ctx, func(ctx context.Context) (err error) {
		usr, added, err := s.ensureUser(ctx, info.Email)
		if err != nil {
			return fmt.Errorf("failed to ensure user: %w", err)
		}
		user = usr
		account, err = s.ensureAccount(ctx, user.ID, provider)
		if err != nil {
			return fmt.Errorf("failed to ensure account: %w", err)
		}
		if !added {
			return nil
		}
		return s.broker.PublishUserRegister(events.UserRegister{
			ID:     user.ID.String(),
			Email:  user.Email,
			Name:   info.Name,
			Avatar: info.Picture,
		})
	})
	if err != nil {
		return domain.Tokens{}, err
	}

	if account.UserID != user.ID {
		return domain.Tokens{}, domain.ErrInvalidCredentials
	}

	tokens, err := s.createTokens(ctx, user.ID)
	if err != nil {
		return domain.Tokens{}, err
	}
	err = s.broker.PublishUserLogin(events.UserLogin{
		ID:   user.ID.String(),
		IP:   ip,
		Time: time.Now(),
		Type: string(provider),
	})
	if err != nil {
		return domain.Tokens{}, err
	}
	return tokens, nil
}

func (s *authService) ensureUser(ctx context.Context, email string) (domain.User, bool, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return user, false, nil
	}

	if errors.Is(err, domain.ErrUserNotFound) {
		user, err := s.users.Create(ctx, email)
		if err != nil {
			return domain.User{}, false, err
		}
		return user, true, nil
	}

	return domain.User{}, false, err
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
