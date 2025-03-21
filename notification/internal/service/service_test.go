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
	type MockBehavior func(tx *txMocks.TxManager, users *mocks.UserRepo, mailer *mailMocks.Mailer, data events.UserRegister)

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
			mockBehavior: func(tx *txMocks.TxManager, users *mocks.UserRepo, mailer *mailMocks.Mailer, data events.UserRegister) {
				tx.EXPECT().Run(mock.Anything, mock.AnythingOfType("func(context.Context) error")).RunAndReturn(
					func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					},
				)
				users.EXPECT().SaveSubscription(mock.Anything, data.ID, domain.SubscriptionTypeEmail).Return(nil)
				users.EXPECT().SaveUser(mock.Anything, domain.User{ID: data.ID, Email: data.Email}).Return(nil)
				mailer.EXPECT().SendRegisterEmail(data.Email).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "failed",
			data: events.UserRegister{
				ID:    "user123",
				Email: "user@example.com",
			},
			mockBehavior: func(tx *txMocks.TxManager, users *mocks.UserRepo, mailer *mailMocks.Mailer, data events.UserRegister) {
				tx.EXPECT().Run(mock.Anything, mock.AnythingOfType("func(context.Context) error")).RunAndReturn(
					func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					},
				)
				users.EXPECT().SaveSubscription(mock.Anything, data.ID, domain.SubscriptionTypeEmail).Return(nil)
				users.EXPECT().SaveUser(mock.Anything, domain.User{ID: data.ID, Email: data.Email}).Return(assert.AnError)
				mailer.EXPECT().SendRegisterEmail(data.Email).Return(nil)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := txMocks.NewTxManager(t)
			users := mocks.NewUserRepo(t)
			mailer := mailMocks.NewMailer(t)
			svc := service.New(tx, mailer, users)
			tc.mockBehavior(tx, users, mailer, tc.data)
			err := svc.HandleRegister(context.Background(), tc.data)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestService_SendLoginNotification(t *testing.T) {
	type MockBehavior func(users *mocks.UserRepo, mailer *mailMocks.Mailer, data events.UserLogin)

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
			mockBehavior: func(users *mocks.UserRepo, mailer *mailMocks.Mailer, data events.UserLogin) {
				users.EXPECT().Subscriptions(mock.Anything, data.ID).Return([]domain.Subscription{
					{User: domain.User{Email: "user@example.com"}, Type: domain.SubscriptionTypeEmail, Enabled: true},
					{User: domain.User{TelegramID: 123}, Type: domain.SubscriptionTypeTelegram, Enabled: true},
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
			name: "email only",
			data: events.UserLogin{
				ID: "user123",
			},
			mockBehavior: func(users *mocks.UserRepo, mailer *mailMocks.Mailer, data events.UserLogin) {
				users.EXPECT().Subscriptions(mock.Anything, data.ID).Return([]domain.Subscription{
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
			mockBehavior: func(users *mocks.UserRepo, mailer *mailMocks.Mailer, data events.UserLogin) {
				users.EXPECT().Subscriptions(mock.Anything, data.ID).Return([]domain.Subscription{
					{User: domain.User{TelegramID: 123}, Type: domain.SubscriptionTypeTelegram, Enabled: true},
				}, nil)
			},
			wantErr: nil,
		},
		{
			name: "failed",
			data: events.UserLogin{
				ID: "user123",
			},
			mockBehavior: func(users *mocks.UserRepo, mailer *mailMocks.Mailer, data events.UserLogin) {
				users.EXPECT().Subscriptions(mock.Anything, data.ID).Return(nil, assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := txMocks.NewTxManager(t)
			users := mocks.NewUserRepo(t)
			mailer := mailMocks.NewMailer(t)
			svc := service.New(tx, mailer, users)
			tc.mockBehavior(users, mailer, tc.data)
			err := svc.SendLoginNotification(context.Background(), tc.data)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
