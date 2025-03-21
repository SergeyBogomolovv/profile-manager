package repo

import (
	"database/sql"
	"time"

	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID      `db:"id"`
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

type Subscription struct {
	UserID  uuid.UUID               `db:"user_id"`
	Type    domain.SubscriptionType `db:"type"`
	Enabled bool                    `db:"enabled"`
}

func (s Subscription) ToDomain() domain.Subscription {
	return domain.Subscription{
		UserID:  s.UserID.String(),
		Type:    s.Type,
		Enabled: s.Enabled,
	}
}
