package repo

import (
	"database/sql"
	"time"

	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID      `db:"user_id"`
	Email      sql.NullString `db:"email"`
	TelegramID sql.NullInt64  `db:"telegram_id"`
	CreatedAt  time.Time      `db:"created_at"`
}

func (u User) ToDomain() domain.User {
	return domain.User{
		ID:         u.ID.String(),
		Email:      u.Email.String,
		TelegramID: u.TelegramID.Int64,
	}
}

type SubscriptionWithUser struct {
	UserID     uuid.UUID               `db:"user_id"`
	Email      sql.NullString          `db:"email"`
	TelegramID sql.NullInt64           `db:"telegram_id"`
	Type       domain.SubscriptionType `db:"type"`
	Enabled    bool                    `db:"enabled"`
}

func (s SubscriptionWithUser) ToDomain() domain.Subscription {
	return domain.Subscription{
		User: domain.User{
			ID:         s.UserID.String(),
			Email:      s.Email.String,
			TelegramID: s.TelegramID.Int64,
		},
		Type:    s.Type,
		Enabled: s.Enabled,
	}
}
