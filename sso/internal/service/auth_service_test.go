package service_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/service"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	type args struct {
		email    string
		password string
	}

	type MockBehavior func(users *mocks.UserRepo, tokens *mocks.TokenRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         domain.Tokens
		wantErr      error
	}{
		{
			name: "success",
			args: args{email: "email@email.com", password: "password"},
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, args args) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.password), bcrypt.DefaultCost)
				require.NoError(t, err)
				userID := uuid.New()
				users.EXPECT().GetByEmail(mock.Anything, args.email).Return(domain.User{ID: userID}, nil)
				users.EXPECT().
					AccountByID(mock.Anything, userID, domain.AccountTypeCredentials).
					Return(domain.Account{Password: hashedPassword}, nil)
				tokens.EXPECT().Create(mock.Anything, userID).Return("token", nil)
			},
			want:    domain.Tokens{RefreshToken: "token"},
			wantErr: nil,
		},
		{
			name: "account not exists",
			args: args{email: "email@email.com", password: "password"},
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, args args) {
				userID := uuid.New()
				users.EXPECT().GetByEmail(mock.Anything, args.email).Return(domain.User{ID: userID}, nil)
				users.EXPECT().AccountByID(mock.Anything, userID, domain.AccountTypeCredentials).Return(domain.Account{}, domain.ErrAccountNotFound)
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "user not exists",
			args: args{email: "email@email.com", password: "password"},
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, args args) {
				users.EXPECT().GetByEmail(mock.Anything, args.email).Return(domain.User{}, domain.ErrUserNotFound)
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			args: args{email: "email@email.com", password: "password"},
			mockBehavior: func(users *mocks.UserRepo, tokens *mocks.TokenRepo, args args) {
				userID := uuid.New()
				users.EXPECT().GetByEmail(mock.Anything, args.email).Return(domain.User{ID: userID}, nil)
				users.EXPECT().
					AccountByID(mock.Anything, userID, domain.AccountTypeCredentials).
					Return(domain.Account{Password: []byte("wrongpass")}, nil)
			},
			wantErr: domain.ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := mocks.NewUserRepo(t)
			tokenRepo := mocks.NewTokenRepo(t)
			broker := mocks.NewBroker(t)
			svc := service.NewAuthService(broker, nil, userRepo, tokenRepo, []byte("secret"))
			tc.mockBehavior(userRepo, tokenRepo, tc.args)
			tokens, err := svc.Login(context.Background(), tc.args.email, tc.args.password)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NotEmpty(t, tokens.AccessToken)
			assert.Equal(t, tokens.RefreshToken, tc.want.RefreshToken)
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	type args struct {
		refreshToken string
		userID       uuid.UUID
	}

	type MockBehavior func(tokens *mocks.TokenRepo, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      error
	}{
		{
			name: "success",
			args: args{refreshToken: "token", userID: uuid.New()},
			mockBehavior: func(tokens *mocks.TokenRepo, args args) {
				tokens.EXPECT().UserID(mock.Anything, args.refreshToken).Return(args.userID, nil)
			},
			wantErr: nil,
		},
		{
			name: "ivalid token",
			args: args{refreshToken: "token"},
			mockBehavior: func(tokens *mocks.TokenRepo, args args) {
				tokens.EXPECT().UserID(mock.Anything, args.refreshToken).Return(uuid.Nil, domain.ErrInvalidToken)
			},
			wantErr: domain.ErrInvalidToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenRepo := mocks.NewTokenRepo(t)
			secret := []byte("secret")
			broker := mocks.NewBroker(t)
			svc := service.NewAuthService(broker, nil, nil, tokenRepo, secret)
			tc.mockBehavior(tokenRepo, tc.args)
			accessToken, err := svc.Refresh(context.Background(), tc.args.refreshToken)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			claims, err := service.VerifyJWT(accessToken, secret)
			require.NoError(t, err)
			assert.Equal(t, claims.UserID, tc.args.userID.String())
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	type MockBehavior func(tokens *mocks.TokenRepo, token string)
	testCases := []struct {
		name         string
		token        string
		mockBehavior MockBehavior
		want         error
	}{
		{
			name:  "success",
			token: "token",
			mockBehavior: func(tokens *mocks.TokenRepo, token string) {
				tokens.EXPECT().Revoke(mock.Anything, token).Return(nil)
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenRepo := mocks.NewTokenRepo(t)
			broker := mocks.NewBroker(t)
			svc := service.NewAuthService(broker, nil, nil, tokenRepo, []byte("secret"))
			tc.mockBehavior(tokenRepo, tc.token)
			err := svc.Logout(context.Background(), tc.token)
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	t.Skip()
}

func TestAuthService_OAuth(t *testing.T) {
	t.Skip()
}
