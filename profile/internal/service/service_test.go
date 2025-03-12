package service_test

import (
	"context"
	"testing"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/service"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProfileService_Create(t *testing.T) {
	type MockBehavior func(profiles *mocks.ProfileRepo, data events.UserRegister)

	testCases := []struct {
		name         string
		data         events.UserRegister
		mockBehavior MockBehavior
		want         error
	}{
		{
			name: "succes",
			data: events.UserRegister{
				ID:    "id",
				Email: "user@email.com",
			},
			mockBehavior: func(profiles *mocks.ProfileRepo, data events.UserRegister) {
				profiles.EXPECT().Create(mock.Anything, domain.Profile{
					UserID:   data.ID,
					Username: "user",
				}).Return(nil)
			},
			want: nil,
		},
		{
			name: "failed",
			data: events.UserRegister{
				ID:    "id",
				Email: "user@email.com",
			},
			mockBehavior: func(profiles *mocks.ProfileRepo, data events.UserRegister) {
				profiles.EXPECT().Create(mock.Anything, domain.Profile{
					UserID:   data.ID,
					Username: "user",
				}).Return(assert.AnError)
			},
			want: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			profiles := mocks.NewProfileRepo(t)
			svc := service.NewProfileService(profiles)
			tc.mockBehavior(profiles, tc.data)
			err := svc.Create(context.Background(), tc.data)
			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestProfileService_Update(t *testing.T) {
	type args struct {
		userID string
		dto    domain.UpdateProfileDTO
	}

	type MockBehavior func(profiles *mocks.ProfileRepo, args args)
	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         domain.Profile
		wantErr      error
	}{
		{
			name: "success",
			args: args{
				userID: "id",
				dto: domain.UpdateProfileDTO{
					Username: "username",
				},
			},
			mockBehavior: func(profiles *mocks.ProfileRepo, args args) {
				profiles.EXPECT().ProfileByID(mock.Anything, args.userID).Return(domain.Profile{Username: "user"}, nil)
				profiles.EXPECT().UsernameExists(mock.Anything, args.dto.Username).Return(false, nil)
				profiles.EXPECT().Update(mock.Anything, domain.Profile{Username: "username"}).Return(nil)
			},
			want:    domain.Profile{Username: "username"},
			wantErr: nil,
		},
		{
			name: "profile not found",
			args: args{
				userID: "not valid",
			},
			mockBehavior: func(profiles *mocks.ProfileRepo, args args) {
				profiles.EXPECT().ProfileByID(mock.Anything, args.userID).Return(domain.Profile{}, domain.ErrProfileNotFound)
			},
			want:    domain.Profile{},
			wantErr: domain.ErrProfileNotFound,
		},
		{
			name: "username already exists",
			args: args{
				userID: "id",
				dto: domain.UpdateProfileDTO{
					Username: "new user",
				},
			},
			mockBehavior: func(profiles *mocks.ProfileRepo, args args) {
				profiles.EXPECT().ProfileByID(mock.Anything, args.userID).Return(domain.Profile{Username: "user"}, nil)
				profiles.EXPECT().UsernameExists(mock.Anything, args.dto.Username).Return(true, nil)
			},
			want:    domain.Profile{},
			wantErr: domain.ErrUsernameExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			profiles := mocks.NewProfileRepo(t)
			svc := service.NewProfileService(profiles)
			tc.mockBehavior(profiles, tc.args)
			got, err := svc.Update(context.Background(), tc.args.userID, tc.args.dto)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, got, tc.want)
		})
	}
}
