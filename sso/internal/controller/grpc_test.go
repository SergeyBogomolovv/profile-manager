package controller_test

import (
	"context"
	"testing"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/sso"
	"github.com/SergeyBogomolovv/profile-manager/common/testutils"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/controller"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/controller/mocks"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGRPCController_Register(t *testing.T) {
	type args struct {
		req *pb.RegisterRequest
	}

	type MockBehavior func(svc *mocks.AuthService, req *pb.RegisterRequest)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         *pb.RegisterResponse
		wantErr      bool
	}{
		{
			name: "success",
			args: args{req: &pb.RegisterRequest{Email: "xLb3u@example.com", Password: "password"}},
			mockBehavior: func(svc *mocks.AuthService, req *pb.RegisterRequest) {
				svc.EXPECT().Register(mock.Anything, req.Email, req.Password).Return("user_id", nil).Once()
			},
			want:    &pb.RegisterResponse{UserId: "user_id"},
			wantErr: false,
		},
		{
			name:         "invalid email",
			mockBehavior: func(svc *mocks.AuthService, req *pb.RegisterRequest) {},
			args:         args{req: &pb.RegisterRequest{Email: "invalid-email", Password: "password"}},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "user already exists",
			args: args{req: &pb.RegisterRequest{Email: "xLb3u@example.com", Password: "password"}},
			mockBehavior: func(svc *mocks.AuthService, req *pb.RegisterRequest) {
				svc.EXPECT().Register(mock.Anything, req.Email, req.Password).Return("", domain.ErrUserAlreadyExists).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to register user",
			args: args{req: &pb.RegisterRequest{Email: "xLb3u@example.com", Password: "password"}},
			mockBehavior: func(svc *mocks.AuthService, req *pb.RegisterRequest) {
				svc.EXPECT().Register(mock.Anything, req.Email, req.Password).Return("", assert.AnError).Once()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewAuthService(t)
			controller := controller.NewGRPCController(testutils.NewTestLogger(), svc)
			tc.mockBehavior(svc, tc.args.req)
			got, err := controller.Register(context.Background(), tc.args.req)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGRPCController_Login(t *testing.T) {
	type args struct {
		req *pb.LoginRequest
	}

	type MockBehavior func(svc *mocks.AuthService, req *pb.LoginRequest)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         *pb.TokensResponse
		wantErr      bool
	}{
		{
			name: "success",
			args: args{req: &pb.LoginRequest{Email: "xLb3u@example.com", Password: "password"}},
			mockBehavior: func(svc *mocks.AuthService, req *pb.LoginRequest) {
				svc.EXPECT().Login(mock.Anything, req.Email, req.Password).Return(domain.Tokens{AccessToken: "access_token", RefreshToken: "refresh_token"}, nil).Once()
			},
			want:    &pb.TokensResponse{AccessToken: "access_token", RefreshToken: "refresh_token"},
			wantErr: false,
		},
		{
			name:         "invalid email",
			mockBehavior: func(svc *mocks.AuthService, req *pb.LoginRequest) {},
			args:         args{req: &pb.LoginRequest{Email: "invalid-email", Password: "password"}},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "invalid credentials",
			args: args{req: &pb.LoginRequest{Email: "xLb3u@example.com", Password: "invalid-password"}},
			mockBehavior: func(svc *mocks.AuthService, req *pb.LoginRequest) {
				svc.EXPECT().Login(mock.Anything, req.Email, req.Password).Return(domain.Tokens{}, domain.ErrInvalidCredentials).Once()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewAuthService(t)
			controller := controller.NewGRPCController(testutils.NewTestLogger(), svc)
			tc.mockBehavior(svc, tc.args.req)
			got, err := controller.Login(context.Background(), tc.args.req)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGRPCController_Refresh(t *testing.T) {
	type args struct {
		req *pb.RefreshRequest
	}

	type MockBehavior func(svc *mocks.AuthService, req *pb.RefreshRequest)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         *pb.AccessTokenResponse
		wantErr      bool
	}{
		{
			name: "success",
			args: args{req: &pb.RefreshRequest{RefreshToken: "refresh_token"}},
			mockBehavior: func(svc *mocks.AuthService, req *pb.RefreshRequest) {
				svc.EXPECT().Refresh(mock.Anything, req.RefreshToken).Return("access_token", nil).Once()
			},
			want:    &pb.AccessTokenResponse{AccessToken: "access_token"},
			wantErr: false,
		},
		{
			name: "invalid token",
			mockBehavior: func(svc *mocks.AuthService, req *pb.RefreshRequest) {
				svc.EXPECT().Refresh(mock.Anything, req.RefreshToken).Return("", domain.ErrInvalidToken).Once()
			},
			args:    args{req: &pb.RefreshRequest{RefreshToken: "invalid_token"}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewAuthService(t)
			controller := controller.NewGRPCController(testutils.NewTestLogger(), svc)
			tc.mockBehavior(svc, tc.args.req)
			got, err := controller.Refresh(context.Background(), tc.args.req)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
