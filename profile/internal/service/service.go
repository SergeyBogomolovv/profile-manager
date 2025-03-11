package service

import (
	"context"
	"strings"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
)

type ProfileRepo interface {
	Create(ctx context.Context, profile domain.Profile) error
}

type profileService struct {
	repo ProfileRepo
}

func NewProfileService(repo ProfileRepo) *profileService {
	return &profileService{repo: repo}
}

func (s *profileService) Create(ctx context.Context, user events.UserRegister) error {
	username, _, _ := strings.Cut(user.Email, "@")
	profile := domain.Profile{
		UserID:    user.ID,
		Username:  username,
		FirstName: user.Name,
		Avatar:    user.Avatar,
	}
	return s.repo.Create(ctx, profile)
}
