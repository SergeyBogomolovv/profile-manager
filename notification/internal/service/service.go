package service

import "github.com/SergeyBogomolovv/profile-manager/notification/internal/mailer"

type UserRepo interface{}

type service struct {
	mailer mailer.Mailer
	users  UserRepo
}

func New(mailer mailer.Mailer, users UserRepo) *service {
	return &service{mailer: mailer, users: users}
}
