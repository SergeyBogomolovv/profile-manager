package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
)

type ProfileRepo interface{}

type profileService struct {
	repo ProfileRepo
}

func NewProfileService(repo ProfileRepo) *profileService {
	return &profileService{repo: repo}
}

func (s *profileService) Create(ctx context.Context, user events.UserRegister) error {
	fmt.Printf("Create user: %+v\n", user)
	return errors.New("not implemented")
}
