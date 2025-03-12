package controller_test

import (
	"context"
	"log/slog"
	"testing"

	pb "github.com/SergeyBogomolovv/profile-manager/common/api/profile"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/controller"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/controller/mocks"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGRPCController_GetProfile(t *testing.T) {
	type MockBehavior func(svc *mocks.ProfileService, req *pb.GetProfileRequest)

	userId := uuid.NewString()
	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.GetProfileRequest
		want         string
		wantErr      bool
	}{
		{
			name: "success",
			req:  &pb.GetProfileRequest{UserId: userId},
			mockBehavior: func(svc *mocks.ProfileService, req *pb.GetProfileRequest) {
				svc.EXPECT().GetProfile(mock.Anything, req.UserId).
					Return(domain.Profile{UserID: req.UserId}, nil).Once()
			},
			want:    userId,
			wantErr: false,
		},
		{
			name:         "invalid id",
			req:          &pb.GetProfileRequest{UserId: "invalid_id"},
			mockBehavior: func(svc *mocks.ProfileService, req *pb.GetProfileRequest) {},
			want:         "",
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewProfileService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewGRPCController(slog.Default(), svc)
			got, err := c.GetProfile(context.Background(), tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, got.UserId, tc.want)
		})
	}
}
