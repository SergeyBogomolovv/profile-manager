package service

import (
	"context"
	"errors"

	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
)

type SetupTokenRepo interface {
	Create(ctx context.Context, userID string) (string, error)
	CheckUserID(ctx context.Context, token string) (string, error)
	Revoke(ctx context.Context, token string) error
}

type SetupUserRepo interface {
	IsExists(ctx context.Context, userID string) (bool, error)
	GetByID(ctx context.Context, userID string) (domain.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type SetupSubsRepo interface {
	Save(ctx context.Context, userID string, subType domain.SubscriptionType) error
	IsExists(ctx context.Context, userID string, subType domain.SubscriptionType) (bool, error)
	Update(ctx context.Context, userID string, subType domain.SubscriptionType, enabled bool) error
	Delete(ctx context.Context, userID string, subType domain.SubscriptionType) error
}

type setupService struct {
	txManager transaction.TxManager
	users     SetupUserRepo
	tokens    SetupTokenRepo
	subs      SetupSubsRepo
}

func NewSetupService(txManager transaction.TxManager, users SetupUserRepo, tokens SetupTokenRepo, subs SetupSubsRepo) *setupService {
	return &setupService{txManager: txManager, users: users, tokens: tokens, subs: subs}
}

func (s *setupService) LinkTelegram(ctx context.Context, token string, telegramID int64) error {
	return s.txManager.Run(ctx, func(ctx context.Context) error {
		userID, err := s.tokens.CheckUserID(ctx, token)
		if err != nil {
			return err
		}
		user, err := s.users.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		if user.TelegramID == telegramID {
			return domain.ErrActionDontNeeded
		}
		user.TelegramID = telegramID
		if err := s.users.Update(ctx, user); err != nil {
			return err
		}
		subExists, err := s.subs.IsExists(ctx, userID, domain.SubscriptionTypeTelegram)
		if err != nil {
			return err
		}
		if !subExists {
			if err := s.subs.Save(ctx, userID, domain.SubscriptionTypeTelegram); err != nil {
				return err
			}
		}
		return s.tokens.Revoke(ctx, token)
	})
}

func (s *setupService) UnlinkTelegram(ctx context.Context, telegramID int64) error {
	return s.txManager.Run(ctx, func(ctx context.Context) error {
		user, err := s.users.GetByTelegramID(ctx, telegramID)
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrActionDontNeeded
		}
		if err != nil {
			return err
		}
		if err := s.subs.Delete(ctx, user.ID, domain.SubscriptionTypeTelegram); err != nil {
			return err
		}
		user.TelegramID = 0
		return s.users.Update(ctx, user)
	})
}

func (s *setupService) UpdateSubscriptionStatus(ctx context.Context, telegramID int64, subType domain.SubscriptionType, enabled bool) error {
	user, err := s.users.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return err
	}
	return s.subs.Update(ctx, user.ID, subType, enabled)
}

func (s *setupService) GenerateToken(ctx context.Context, userID string) (string, error) {
	isExists, err := s.users.IsExists(ctx, userID)
	if err != nil {
		return "", err
	}
	if !isExists {
		return "", domain.ErrUserNotFound
	}
	return s.tokens.Create(ctx, userID)
}
