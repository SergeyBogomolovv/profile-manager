package repo

import (
	"database/sql"

	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
)

type Profile struct {
	UserID    string         `db:"user_id"`
	Username  string         `db:"username"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	BirthDate sql.NullString `db:"birth_date"`
	Gender    string         `db:"gender"`
	Avatar    sql.NullString `db:"avatar"`
}

func (p Profile) ToDomain() domain.Profile {
	return domain.Profile{
		UserID:    p.UserID,
		Username:  p.Username,
		FirstName: p.FirstName.String,
		LastName:  p.LastName.String,
		BirthDate: p.BirthDate.String,
		Gender:    domain.UserGender(p.Gender),
		Avatar:    p.Avatar.String,
	}
}
