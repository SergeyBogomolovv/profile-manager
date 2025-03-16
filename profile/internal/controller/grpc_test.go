package controller_test

import (
	"context"
	"log/slog"
	"testing"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/profile"
	"github.com/SergeyBogomolovv/profile-manager/common/auth"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/controller"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/controller/mocks"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestGRPCController_GetProfile(t *testing.T) {
	type MockBehavior func(svc *mocks.ProfileService, userID string)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		userID       string
		wantErr      bool
	}{
		{
			name:   "success",
			userID: uuid.NewString(),
			mockBehavior: func(svc *mocks.ProfileService, userID string) {
				svc.EXPECT().GetProfile(mock.Anything, userID).
					Return(domain.Profile{UserID: userID}, nil).Once()
			},
			wantErr: false,
		},
		{
			name:         "invalid id",
			userID:       "invalid",
			mockBehavior: func(svc *mocks.ProfileService, userID string) {},
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewProfileService(t)
			tc.mockBehavior(svc, tc.userID)
			c := controller.NewGRPCController(slog.Default(), svc)
			md := metadata.New(map[string]string{auth.UserIdKey: tc.userID})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			got, err := c.GetProfile(ctx, &pb.GetProfileRequest{})
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, got.UserId, tc.userID)
		})
	}
}

func TestGRPCController_UpdateProfile(t *testing.T) {
	type args struct {
		req    *pb.UpdateProfileRequest
		userID string
	}
	type MockBehavior func(svc *mocks.ProfileService, args args)

	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		args         args
		want         string
		wantErr      bool
	}{
		{
			name:    "success",
			args:    args{req: &pb.UpdateProfileRequest{Username: "new"}, userID: uuid.NewString()},
			want:    "new",
			wantErr: false,
			mockBehavior: func(svc *mocks.ProfileService, args args) {
				svc.EXPECT().Update(mock.Anything, args.userID, domain.UpdateProfileDTO{Username: args.req.Username}).
					Return(domain.Profile{Username: args.req.Username}, nil).Once()
			},
		},
		{
			name:         "ivalid date",
			args:         args{req: &pb.UpdateProfileRequest{BirthDate: "12312334"}},
			wantErr:      true,
			mockBehavior: func(svc *mocks.ProfileService, args args) {},
		},
		{
			name:         "invalid gender",
			args:         args{req: &pb.UpdateProfileRequest{Gender: "helicopter"}},
			wantErr:      true,
			mockBehavior: func(svc *mocks.ProfileService, args args) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewProfileService(t)
			tc.mockBehavior(svc, tc.args)
			c := controller.NewGRPCController(slog.Default(), svc)
			md := metadata.New(map[string]string{auth.UserIdKey: tc.args.userID})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			got, err := c.UpdateProfile(ctx, tc.args.req)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, got.Username, tc.want)
		})
	}
}
