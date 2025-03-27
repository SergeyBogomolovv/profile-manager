package service

import (
	"context"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/mailer"
	"golang.org/x/sync/errgroup"
)

type TokenRepo interface {
	Create(ctx context.Context, userID string) (string, error)
	CheckUserID(ctx context.Context, token string) (string, error)
	Revoke(ctx context.Context, token string) error
}

type Sender interface {
	SendLoginNotification(telegramID int64, data domain.LoginNotification) error
}

type UserRepo interface {
	Save(ctx context.Context, user domain.User) error
	IsExists(ctx context.Context, userID string) (bool, error)
	GetByID(ctx context.Context, userID string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type SubscriptionRepo interface {
	SubscriptionsByUser(ctx context.Context, userID string) ([]domain.Subscription, error)
	Save(ctx context.Context, userID string, subType domain.SubscriptionType) error
	IsExists(ctx context.Context, userID string, subType domain.SubscriptionType) (bool, error)
}

type service struct {
	txManager     transaction.TxManager
	mailer        mailer.Mailer
	users         UserRepo
	sender        Sender
	tokens        TokenRepo
	subscriptions SubscriptionRepo
}

func New(txManager transaction.TxManager, mailer mailer.Mailer, sender Sender, users UserRepo, tokens TokenRepo, subscriptions SubscriptionRepo) *service {
	return &service{mailer: mailer, users: users, txManager: txManager, sender: sender, tokens: tokens, subscriptions: subscriptions}
}

func (s *service) SendLoginNotification(ctx context.Context, data events.UserLogin) error {
	subscriptions, err := s.subscriptions.SubscriptionsByUser(ctx, data.ID)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	notification := domain.LoginNotification{IP: data.IP, Time: data.Time.Format("2006-01-02 15:04:05"), Type: data.Type}
	for _, sub := range subscriptions {
		if sub.Enabled {
			switch sub.Type {
			case domain.SubscriptionTypeEmail:
				eg.Go(func() error {
					return s.mailer.SendLoginEmail(sub.User.Email, notification)
				})
			case domain.SubscriptionTypeTelegram:
				eg.Go(func() error {
					return s.sender.SendLoginNotification(sub.User.TelegramID, notification)
				})
			}
		}
	}

	return eg.Wait()
}

func (s *service) HandleRegister(ctx context.Context, data events.UserRegister) error {
	return s.txManager.Run(ctx, func(ctx context.Context) error {
		if err := s.users.Save(ctx, domain.User{ID: data.ID, Email: data.Email}); err != nil {
			return err
		}
		if err := s.subscriptions.Save(ctx, data.ID, domain.SubscriptionTypeEmail); err != nil {
			return err
		}
		return s.mailer.SendRegisterEmail(data.Email)
	})
}

func (s *service) VerifyTelegram(ctx context.Context, token string, telegramID int64) error {
	return s.txManager.Run(ctx, func(ctx context.Context) error {
		userID, err := s.tokens.CheckUserID(ctx, token)
		if err != nil {
			return err
		}
		user, err := s.users.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		user.TelegramID = telegramID
		if err := s.users.Update(ctx, user); err != nil {
			return err
		}
		subExists, err := s.subscriptions.IsExists(ctx, userID, domain.SubscriptionTypeTelegram)
		if err != nil {
			return err
		}
		if !subExists {
			if err := s.subscriptions.Save(ctx, userID, domain.SubscriptionTypeTelegram); err != nil {
				return err
			}
		}
		return s.tokens.Revoke(ctx, token)
	})
}

func (s *service) GenerateToken(ctx context.Context, userID string) (string, error) {
	isExists, err := s.users.IsExists(ctx, userID)
	if err != nil {
		return "", err
	}
	if !isExists {
		return "", domain.ErrUserNotFound
	}
	return s.tokens.Create(ctx, userID)
}
