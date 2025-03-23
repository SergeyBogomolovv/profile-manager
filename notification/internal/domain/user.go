package domain

import "errors"

type User struct {
	ID         string
	Email      string
	TelegramID int64
}

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
)
