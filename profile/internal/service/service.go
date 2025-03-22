package service

import (
	"context"
	"strings"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
)

type ImageRepo interface {
	UploadAvatar(ctx context.Context, userID string, body []byte) (string, error)
	DeleteAvatar(ctx context.Context, url string) error
}

type ProfileRepo interface {
	Create(ctx context.Context, profile domain.Profile) error
	ProfileByID(ctx context.Context, id string) (domain.Profile, error)
	Update(ctx context.Context, profile *domain.Profile) error
	UsernameExists(ctx context.Context, username string) (bool, error)
}

type profileService struct {
	repo   ProfileRepo
	images ImageRepo
}

func NewProfileService(repo ProfileRepo, images ImageRepo) *profileService {
	return &profileService{repo: repo, images: images}
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
		return domain.Profile{}, err
	}
	if dto.Username != "" && profile.Username != dto.Username {
		if err := s.checkUsername(ctx, dto.Username); err != nil {
			return domain.Profile{}, err
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
	if dto.Avatar != nil {
		profile.Avatar, err = s.images.UploadAvatar(ctx, profile.UserID, dto.Avatar)
		if err != nil {
			return domain.Profile{}, err
		}
	}

	if err := s.repo.Update(ctx, &profile); err != nil {
		return domain.Profile{}, err
	}
	return profile, nil
}

func (s *profileService) checkUsername(ctx context.Context, username string) error {
	ex, err := s.repo.UsernameExists(ctx, username)
	if err != nil {
		return err
	}
	if ex {
		return domain.ErrUsernameExists
	}
	return nil
}
