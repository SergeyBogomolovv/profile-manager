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
			req:  &pb.GetProfileRequest{},
			mockBehavior: func(svc *mocks.ProfileService, req *pb.GetProfileRequest) {
				svc.EXPECT().GetProfile(mock.Anything, "UPDATE").
					Return(domain.Profile{UserID: "UPDATE"}, nil).Once()
			},
			want:    userId,
			wantErr: false,
		},
		{
			name:         "invalid id",
			req:          &pb.GetProfileRequest{},
			mockBehavior: func(svc *mocks.ProfileService, req *pb.GetProfileRequest) {},
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

func TestGRPCController_UpdateProfile(t *testing.T) {
	type MockBehavior func(svc *mocks.ProfileService, req *pb.UpdateProfileRequest)

	// userId := uuid.NewString()
	testCases := []struct {
		name         string
		mockBehavior MockBehavior
		req          *pb.UpdateProfileRequest
		want         string
		wantErr      bool
	}{
		{
			name:    "success",
			req:     &pb.UpdateProfileRequest{Username: "new", BirthDate: "1990-10-10"},
			want:    "new",
			wantErr: false,
			mockBehavior: func(svc *mocks.ProfileService, req *pb.UpdateProfileRequest) {
				svc.EXPECT().Update(mock.Anything, "UPDATE", domain.UpdateProfileDTO{Username: req.Username, BirthDate: req.BirthDate}).
					Return(domain.Profile{Username: req.Username, BirthDate: req.BirthDate}, nil).Once()
			},
		},
		{
			name:         "ivalid date",
			req:          &pb.UpdateProfileRequest{BirthDate: "12312334"},
			wantErr:      true,
			mockBehavior: func(svc *mocks.ProfileService, req *pb.UpdateProfileRequest) {},
		},
		{
			name:         "invalid gender",
			req:          &pb.UpdateProfileRequest{Gender: "helicopter"},
			wantErr:      true,
			mockBehavior: func(svc *mocks.ProfileService, req *pb.UpdateProfileRequest) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewProfileService(t)
			tc.mockBehavior(svc, tc.req)
			c := controller.NewGRPCController(slog.Default(), svc)
			got, err := c.UpdateProfile(context.Background(), tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, got.Username, tc.want)
		})
	}
}
