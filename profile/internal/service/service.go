package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
)

type ProfileRepo interface {
	Create(ctx context.Context, profile domain.Profile) error
	ProfileByID(ctx context.Context, id string) (domain.Profile, error)
	Update(ctx context.Context, profile domain.Profile) error
	UsernameExists(ctx context.Context, username string) (bool, error)
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

func (s *profileService) GetProfile(ctx context.Context, userID string) (domain.Profile, error) {
	return s.repo.ProfileByID(ctx, userID)
}

func (s *profileService) Update(ctx context.Context, userID string, dto domain.UpdateProfileDTO) (domain.Profile, error) {
	profile, err := s.repo.ProfileByID(ctx, userID)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("failed to get profile: %w", err)
	}
	if dto.Username != "" && profile.Username != dto.Username {
		ex, err := s.repo.UsernameExists(ctx, dto.Username)
		if err != nil {
			return domain.Profile{}, fmt.Errorf("failed to check username: %w", err)
		}
		if ex {
			return domain.Profile{}, domain.ErrUsernameExists
		}
		profile.Username = dto.Username
	}
	if dto.BirthDate != "" {
		profile.BirthDate = dto.BirthDate
	}
	if dto.FirstName != "" {
		profile.FirstName = dto.FirstName
	}
	if dto.LastName != "" {
		profile.LastName = dto.LastName
	}
	if dto.Gender != "" {
		profile.Gender = domain.UserGender(dto.Gender)
	}
	if err := s.repo.Update(ctx, profile); err != nil {
		return domain.Profile{}, err
	}
	return profile, nil
}
