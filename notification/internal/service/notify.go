package service

import (
	"context"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/mailer"
	"golang.org/x/sync/errgroup"
)

type Sender interface {
	SendLoginNotification(telegramID int64, data domain.LoginNotification) error
}

type NotifyUserRepo interface {
	Save(ctx context.Context, user domain.User) error
	IsExists(ctx context.Context, userID string) (bool, error)
	GetByID(ctx context.Context, userID string) (domain.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
}

type NotifySubsRepo interface {
	SubscriptionsByUser(ctx context.Context, userID string) ([]domain.Subscription, error)
	Save(ctx context.Context, userID string, subType domain.SubscriptionType) error
	IsExists(ctx context.Context, userID string, subType domain.SubscriptionType) (bool, error)
	Update(ctx context.Context, userID string, subType domain.SubscriptionType, enabled bool) error
	Delete(ctx context.Context, userID string, subType domain.SubscriptionType) error
}

type service struct {
	txManager transaction.TxManager
	mailer    mailer.Mailer
	users     NotifyUserRepo
	sender    Sender
	subs      NotifySubsRepo
}

func NewNotifyService(txManager transaction.TxManager, mailer mailer.Mailer, sender Sender, users NotifyUserRepo, subs NotifySubsRepo) *service {
	return &service{mailer: mailer, users: users, txManager: txManager, sender: sender, subs: subs}
}

func (s *service) SendLoginNotification(ctx context.Context, data events.UserLogin) error {
	subscriptions, err := s.subs.SubscriptionsByUser(ctx, data.ID)
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
		if err := s.subs.Save(ctx, data.ID, domain.SubscriptionTypeEmail); err != nil {
			return err
		}
		return s.mailer.SendRegisterEmail(data.Email)
	})
}
