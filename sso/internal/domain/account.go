package domain

import (
	"errors"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeGoogle      AccountType = "google"
	AccountTypeCredentials AccountType = "credentials"
)

type Account struct {
	UserID   uuid.UUID
	Provider AccountType
	Password []byte
}

var (
	ErrAccountNotFound = errors.New("account not found")
)
