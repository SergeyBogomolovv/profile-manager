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

func (u User) ToDomain() domain.User {
	return domain.User{
		ID:    u.ID,
		Email: u.Email,
	}
}

type Account struct {
	UserID   uuid.UUID          `db:"user_id"`
	Provider domain.AccountType `db:"provider"`
	Password []byte             `db:"password"`
}

func (a Account) ToDomain() domain.Account {
	return domain.Account{
		UserID:   a.UserID,
		Provider: a.Provider,
		Password: a.Password,
	}
}

type RefreshToken struct {
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
