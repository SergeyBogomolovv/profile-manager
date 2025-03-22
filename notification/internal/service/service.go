package service

import (
	"context"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/mailer"
	"golang.org/x/sync/errgroup"
)

type UserRepo interface {
	Subscriptions(ctx context.Context, userID string) ([]domain.Subscription, error)
	SaveSubscription(ctx context.Context, userID string, subType domain.SubscriptionType) error
	SaveUser(ctx context.Context, user domain.User) error
}

type service struct {
	txManager transaction.TxManager
	mailer    mailer.Mailer
	users     UserRepo
}

func New(txManager transaction.TxManager, mailer mailer.Mailer, users UserRepo) *service {
	return &service{mailer: mailer, users: users, txManager: txManager}
}

func (s *service) SendLoginNotification(ctx context.Context, data events.UserLogin) error {
	subscriptions, err := s.users.Subscriptions(ctx, data.ID)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, sub := range subscriptions {
		if sub.Enabled {
			switch sub.Type {
			case domain.SubscriptionTypeEmail:
				eg.Go(func() error {
					return s.sendLoginEmail(sub.User.Email, data)
				})
			case domain.SubscriptionTypeTelegram:
				eg.Go(func() error {
					return s.sendLoginTelegram(sub.User.TelegramID, data)
				})
			}
		}
	}

	return eg.Wait()
}

func (s *service) sendLoginEmail(to string, data events.UserLogin) error {
	return s.mailer.SendLoginEmail(to, domain.LoginNotification{IP: data.IP, Time: data.Time.Format("2006-01-02 15:04:05"), Type: data.Type})
}

func (s *service) sendLoginTelegram(id int64, data events.UserLogin) error {
	// TODO: implement
	return nil
}

func (s *service) HandleRegister(ctx context.Context, data events.UserRegister) error {
	return s.txManager.Run(ctx, func(ctx context.Context) error {
		eg, ctx := errgroup.WithContext(ctx)
		eg.Go(func() error {
			return s.users.SaveSubscription(ctx, data.ID, domain.SubscriptionTypeEmail)
		})
		eg.Go(func() error {
			return s.users.SaveUser(ctx, domain.User{ID: data.ID, Email: data.Email})
		})
		eg.Go(func() error {
			return s.mailer.SendRegisterEmail(data.Email)
		})
		return eg.Wait()
	})
}
