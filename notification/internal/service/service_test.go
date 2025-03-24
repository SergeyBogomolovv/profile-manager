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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_HandleRegister(t *testing.T) {
	type MockBehavior func(
		tx *txMocks.TxManager,
		subscriptions *mocks.SubscriptionRepo,
		users *mocks.UserRepo,
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
				subscriptions *mocks.SubscriptionRepo,
				users *mocks.UserRepo,
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
				subscriptions *mocks.SubscriptionRepo,
				users *mocks.UserRepo,
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
				subscriptions *mocks.SubscriptionRepo,
				users *mocks.UserRepo,
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
			users := mocks.NewUserRepo(t)
			subscriptions := mocks.NewSubscriptionRepo(t)
			tokens := mocks.NewTokenRepo(t)
			mailer := mailMocks.NewMailer(t)
			sender := mocks.NewSender(t)
			svc := service.New(tx, mailer, sender, users, tokens, subscriptions)
			tc.mockBehavior(tx, subscriptions, users, mailer, tc.data)
			err := svc.HandleRegister(context.Background(), tc.data)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestService_SendLoginNotification(t *testing.T) {
	type MockBehavior func(subs *mocks.SubscriptionRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin)

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
			mockBehavior: func(subs *mocks.SubscriptionRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin) {
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
			mockBehavior: func(subs *mocks.SubscriptionRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin) {
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
			mockBehavior: func(subs *mocks.SubscriptionRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin) {
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
			mockBehavior: func(subs *mocks.SubscriptionRepo, sender *mocks.Sender, mailer *mailMocks.Mailer, data events.UserLogin) {
				subs.EXPECT().SubscriptionsByUser(mock.Anything, data.ID).Return(nil, assert.AnError)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := txMocks.NewTxManager(t)
			users := mocks.NewUserRepo(t)
			subs := mocks.NewSubscriptionRepo(t)
			tokens := mocks.NewTokenRepo(t)
			mailer := mailMocks.NewMailer(t)
			sender := mocks.NewSender(t)
			svc := service.New(tx, mailer, sender, users, tokens, subs)
			tc.mockBehavior(subs, sender, mailer, tc.data)
			err := svc.SendLoginNotification(context.Background(), tc.data)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestService_VerifyTelegram(t *testing.T) {
	type args struct {
		token      string
		telegramID int64
	}
	type MockBehavior func(
		tx *txMocks.TxManager,
		users *mocks.UserRepo,
		subs *mocks.SubscriptionRepo,
		tokens *mocks.TokenRepo,
		args args,
	)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "need to create subscription",
			args: args{
				token:      "token",
				telegramID: 123,
			},
			want: nil,
			mockBehavior: func(tx *txMocks.TxManager, users *mocks.UserRepo, subs *mocks.SubscriptionRepo, tokens *mocks.TokenRepo, args args) {
				tx.EXPECT().Run(mock.Anything, mock.AnythingOfType("func(context.Context) error")).RunAndReturn(
					func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					},
				)
				tokens.EXPECT().CheckUserID(mock.Anything, args.token).Return("user_id", nil)
				users.EXPECT().GetByID(mock.Anything, "user_id").Return(domain.User{ID: "user_id"}, nil)
				users.EXPECT().Update(mock.Anything, domain.User{ID: "user_id", TelegramID: 123}).Return(nil)
				subs.EXPECT().IsExists(mock.Anything, "user_id", domain.SubscriptionTypeTelegram).Return(false, nil)
				subs.EXPECT().Save(mock.Anything, "user_id", domain.SubscriptionTypeTelegram).Return(nil)
				tokens.EXPECT().Revoke(mock.Anything, args.token).Return(nil)
			},
		},
		{
			name: "subscription exists",
			args: args{
				token:      "token",
				telegramID: 123,
			},
			want: nil,
			mockBehavior: func(tx *txMocks.TxManager, users *mocks.UserRepo, subs *mocks.SubscriptionRepo, tokens *mocks.TokenRepo, args args) {
				tx.EXPECT().Run(mock.Anything, mock.AnythingOfType("func(context.Context) error")).RunAndReturn(
					func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					},
				)
				tokens.EXPECT().CheckUserID(mock.Anything, args.token).Return("user_id", nil)
				users.EXPECT().GetByID(mock.Anything, "user_id").Return(domain.User{ID: "user_id"}, nil)
				users.EXPECT().Update(mock.Anything, domain.User{ID: "user_id", TelegramID: 123}).Return(nil)
				subs.EXPECT().IsExists(mock.Anything, "user_id", domain.SubscriptionTypeTelegram).Return(true, nil)
				tokens.EXPECT().Revoke(mock.Anything, args.token).Return(nil)
			},
		},
		{
			name: "invalid token",
			args: args{
				token:      "token",
				telegramID: 123,
			},
			mockBehavior: func(tx *txMocks.TxManager, users *mocks.UserRepo, subs *mocks.SubscriptionRepo, tokens *mocks.TokenRepo, args args) {
				tx.EXPECT().Run(mock.Anything, mock.AnythingOfType("func(context.Context) error")).RunAndReturn(
					func(ctx context.Context, f func(context.Context) error) error {
						return f(ctx)
					},
				)
				tokens.EXPECT().CheckUserID(mock.Anything, args.token).Return("", domain.ErrInvalidToken)
			},
			want: domain.ErrInvalidToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := txMocks.NewTxManager(t)
			users := mocks.NewUserRepo(t)
			subs := mocks.NewSubscriptionRepo(t)
			tokens := mocks.NewTokenRepo(t)
			mailer := mailMocks.NewMailer(t)
			sender := mocks.NewSender(t)
			svc := service.New(tx, mailer, sender, users, tokens, subs)
			tc.mockBehavior(tx, users, subs, tokens, tc.args)
			err := svc.VerifyTelegram(context.Background(), tc.args.token, tc.args.telegramID)
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestService_GenerateToken(t *testing.T) {
	type MockBehavior func(
		users *mocks.UserRepo,
		tokens *mocks.TokenRepo,
		userID string,
	)

	testCases := []struct {
		name         string
		userID       string
		mockBehavior MockBehavior
		want         string
		wantErr      error
	}{
		{
			name:   "success",
			userID: uuid.NewString(),
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, userID string) {
				users.EXPECT().IsExists(mock.Anything, userID).Return(true, nil)
				tokens.EXPECT().Create(mock.Anything, userID).Return("token", nil)
			},
			want:    "token",
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: uuid.NewString(),
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, userID string) {
				users.EXPECT().IsExists(mock.Anything, userID).Return(false, nil)
			},
			want:    "",
			wantErr: domain.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := txMocks.NewTxManager(t)
			users := mocks.NewUserRepo(t)
			subs := mocks.NewSubscriptionRepo(t)
			tokens := mocks.NewTokenRepo(t)
			mailer := mailMocks.NewMailer(t)
			sender := mocks.NewSender(t)
			svc := service.New(tx, mailer, sender, users, tokens, subs)
			tc.mockBehavior(users, tokens, tc.userID)
			token, err := svc.GenerateToken(context.Background(), tc.userID)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, token, tc.want)
		})
	}
}
