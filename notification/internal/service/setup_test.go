package service_test

import (
	"context"
	"testing"

	txMocks "github.com/SergeyBogomolovv/profile-manager/common/transaction/mocks"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/service"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_VerifyTelegram(t *testing.T) {
	type args struct {
		token      string
		telegramID int64
	}
	type MockBehavior func(
		tx *txMocks.TxManager,
		users *mocks.SetupUserRepo,
		subs *mocks.SetupSubsRepo,
		tokens *mocks.SetupTokenRepo,
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
			mockBehavior: func(tx *txMocks.TxManager, users *mocks.SetupUserRepo, subs *mocks.SetupSubsRepo, tokens *mocks.SetupTokenRepo, args args) {
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
			mockBehavior: func(tx *txMocks.TxManager, users *mocks.SetupUserRepo, subs *mocks.SetupSubsRepo, tokens *mocks.SetupTokenRepo, args args) {
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
			mockBehavior: func(tx *txMocks.TxManager, users *mocks.SetupUserRepo, subs *mocks.SetupSubsRepo, tokens *mocks.SetupTokenRepo, args args) {
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
			users := mocks.NewSetupUserRepo(t)
			subs := mocks.NewSetupSubsRepo(t)
			tokens := mocks.NewSetupTokenRepo(t)
			svc := service.NewSetupService(tx, users, tokens, subs)
			tc.mockBehavior(tx, users, subs, tokens, tc.args)
			err := svc.LinkTelegram(context.Background(), tc.args.token, tc.args.telegramID)
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestService_GenerateToken(t *testing.T) {
	type MockBehavior func(
		users *mocks.SetupUserRepo,
		tokens *mocks.SetupTokenRepo,
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
			mockBehavior: func(users *mocks.SetupUserRepo, tokens *mocks.SetupTokenRepo, userID string) {
				users.EXPECT().IsExists(mock.Anything, userID).Return(true, nil)
				tokens.EXPECT().Create(mock.Anything, userID).Return("token", nil)
			},
			want:    "token",
			wantErr: nil,
		},
		{
			name:   "user not found",
			userID: uuid.NewString(),
			mockBehavior: func(users *mocks.SetupUserRepo, tokens *mocks.SetupTokenRepo, userID string) {
				users.EXPECT().IsExists(mock.Anything, userID).Return(false, nil)
			},
			want:    "",
			wantErr: domain.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx := txMocks.NewTxManager(t)
			users := mocks.NewSetupUserRepo(t)
			subs := mocks.NewSetupSubsRepo(t)
			tokens := mocks.NewSetupTokenRepo(t)
			svc := service.NewSetupService(tx, users, tokens, subs)
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
