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
