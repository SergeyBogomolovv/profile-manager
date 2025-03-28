package service_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	txMocks "github.com/SergeyBogomolovv/profile-manager/common/transaction/mocks"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	mailMocks "github.com/SergeyBogomolovv/profile-manager/notification/internal/mailer/mocks"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/service"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_HandleRegister(t *testing.T) {
	type MockBehavior func(
		tx *txMocks.TxManager,
		subscriptions *mocks.NotifySubsRepo,
		users *mocks.NotifyUserRepo,
		mailer *mailMocks.Mailer,
		data events.UserRegister,
	)

	testCases := []struct {
		name         string
		data         events.UserRegister
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			data: events.UserRegister{
				ID:    "user123",
				Email: "user@example.com",
			},
			mockBehavior: func(
				tx *txMocks.TxManager,
				subscriptions *mocks.NotifySubsRepo,
				users *mocks.NotifyUserRepo,
				mailer *mailMocks.Mailer,
				data events.UserRegister,
			) {
				tx.EXPECT().Run(mock.Anything, mock.AnythingOfType("func(context.Context) error")).RunAndReturn(
					func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					},
				)
				subscriptions.EXPECT().Save(mock.Anything, data.ID, domain.SubscriptionTypeEmail).Return(nil)
				users.EXPECT().Save(mock.Anything, domain.User{ID: data.ID, Email: data.Email}).Return(nil)
				mailer.EXPECT().SendRegisterEmail(data.Email).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "failed to save user",
			data: events.UserRegister{
				ID:    "user123",
				Email: "user@example.com",
			},
			mockBehavior: func(
				tx *txMocks.TxManager,
				subscriptions *mocks.NotifySubsRepo,
				users *mocks.NotifyUserRepo,
				mailer *mailMocks.Mailer,
				data events.UserRegister,
			) {
				tx.EXPECT().Run(mock.Anything, mock.AnythingOfType("func(context.Context) error")).RunAndReturn(
					func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					},
				)
				users.EXPECT().Save(mock.Anything, domain.User{ID: data.ID, Email: data.Email}).Return(assert.AnError)
			},
			wantErr: assert.AnError,
		},
		{
			name: "failed to save subscription",
			data: events.UserRegister{
				ID:    "user123",
				Email: "user@example.com",
			},
			mockBehavior: func(
				tx *txMocks.TxManager,
				subscriptions *mocks.NotifySubsRepo,
				users *mocks.NotifyUserRepo,
				mailer *mailMocks.Mailer,
				data events.UserRegister,
			) {
				tx.EXPECT().Run(mock.Anything, mock.AnythingOfType("func(context.Context) error")).RunAndReturn(
					func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					},
				)
				users.EXPECT().Save(mock.Anything, domain.User{ID: data.ID, Email: data.Email}).Return(nil)
				subscriptions.EXPECT().Save(mock.Anything, data.ID, domain.SubscriptionTypeEmail).Return(assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := txMocks.NewTxManager(t)
			users := mocks.NewNotifyUserRepo(t)
			subs := mocks.NewNotifySubsRepo(t)
			mailer := mailMocks.NewMailer(t)
			sender := mocks.NewSender(t)
			svc := service.NewNotifyService(tx, mailer, sender, users, subs)
			tc.mockBehavior(tx, subs, users, mailer, tc.data)
			err := svc.HandleRegister(context.Background(), tc.data)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestService_SendLoginNotification(t *testing.T) {
	type MockBehavior func(subs *mocks.NotifySubsRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin)

	testCases := []struct {
		name         string
		data         events.UserLogin
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			data: events.UserLogin{
				ID: "user123",
			},
			mockBehavior: func(subs *mocks.NotifySubsRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin) {
				subs.EXPECT().SubscriptionsByUser(mock.Anything, data.ID).Return([]domain.Subscription{
					{User: domain.User{Email: "user@example.com"}, Type: domain.SubscriptionTypeEmail, Enabled: true},
					{User: domain.User{TelegramID: 123}, Type: domain.SubscriptionTypeTelegram, Enabled: true},
				}, nil)
				noti := domain.LoginNotification{
					IP:   data.IP,
					Time: data.Time.Format("2006-01-02 15:04:05"),
					Type: data.Type,
				}
				mailer.EXPECT().SendLoginEmail("user@example.com", noti).Return(nil)
				sender.EXPECT().SendLoginNotification(int64(123), noti).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "email only",
			data: events.UserLogin{
				ID: "user123",
			},
			mockBehavior: func(subs *mocks.NotifySubsRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin) {
				subs.EXPECT().SubscriptionsByUser(mock.Anything, data.ID).Return([]domain.Subscription{
					{User: domain.User{Email: "user@example.com"}, Type: domain.SubscriptionTypeEmail, Enabled: true},
				}, nil)
				mailer.EXPECT().SendLoginEmail("user@example.com",
					domain.LoginNotification{
						IP:   data.IP,
						Time: data.Time.Format("2006-01-02 15:04:05"),
						Type: data.Type,
					}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "telegram only",
			data: events.UserLogin{
				ID: "user123",
			},
			mockBehavior: func(subs *mocks.NotifySubsRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin) {
				subs.EXPECT().SubscriptionsByUser(mock.Anything, data.ID).Return([]domain.Subscription{
					{User: domain.User{TelegramID: 123}, Type: domain.SubscriptionTypeTelegram, Enabled: true},
				}, nil)
				sender.EXPECT().SendLoginNotification(int64(123),
					domain.LoginNotification{
						IP:   data.IP,
						Time: data.Time.Format("2006-01-02 15:04:05"),
						Type: data.Type,
					}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "failed",
			data: events.UserLogin{
				ID: "user123",
			},
			mockBehavior: func(subs *mocks.NotifySubsRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin) {
				subs.EXPECT().SubscriptionsByUser(mock.Anything, data.ID).Return(nil, assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := txMocks.NewTxManager(t)
			users := mocks.NewNotifyUserRepo(t)
			subs := mocks.NewNotifySubsRepo(t)
			mailer := mailMocks.NewMailer(t)
			sender := mocks.NewSender(t)
			svc := service.NewNotifyService(tx, mailer, sender, users, subs)
			tc.mockBehavior(subs, sender, mailer, tc.data)
			err := svc.SendLoginNotification(context.Background(), tc.data)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
