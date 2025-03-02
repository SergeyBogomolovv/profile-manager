package repo

import (
	"time"

	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"user_id"`
	Email        string    `db:"email"`
	RegisteredAt time.Time `db:"registered_at"`
}

type Account struct {
	UserID   uuid.UUID          `db:"user_id"`
	Provider domain.AccountType `db:"provider"`
	Password []byte             `db:"password"`
}

type RefreshToken struct {
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
