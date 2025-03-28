package repo

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/SergeyBogomolovv/profile-manager/common/e"
	"github.com/SergeyBogomolovv/profile-manager/profile/internal/domain"
	"github.com/jmoiron/sqlx"
)

type profileRepo struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func NewProfileRepo(db *sqlx.DB) *profileRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &profileRepo{db: db, qb: qb}
}

func (r *profileRepo) Create(ctx context.Context, profile domain.Profile) error {
	m := map[string]any{
		"user_id":  profile.UserID,
		"username": profile.Username,
	}

	if profile.FirstName != "" {
		m["first_name"] = profile.FirstName
	}
	if profile.Avatar != "" {
		m["avatar"] = profile.Avatar
	}

	query, args := r.qb.Insert("profiles").SetMap(m).MustSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	return e.WrapIfErr(err, "failed to create profile")
}

func (r *profileRepo) ProfileByID(ctx context.Context, id string) (domain.Profile, error) {
	query, args := r.qb.Select("*").From("profiles").Where(sq.Eq{"user_id": id}).MustSql()
	var profile Profile
	if err := r.db.GetContext(ctx, &profile, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Profile{}, domain.ErrProfileNotFound
		}
		return domain.Profile{}, e.Wrap(err, "failed to get profile")
	}
	return profile.ToDomain(), nil
}

func (r *profileRepo) Update(ctx context.Context, profile *domain.Profile) error {
	q := r.qb.Update("profiles").
		Set("username", profile.Username).
		Set("first_name", profile.FirstName).
		Set("last_name", profile.LastName).
		Set("gender", profile.Gender).
		Set("avatar", profile.Avatar)

	if profile.BirthDate != "" {
		q = q.Set("birth_date", profile.BirthDate)
	}
	query, args := q.Where(sq.Eq{"user_id": profile.UserID}).Suffix("RETURNING *").MustSql()
	var p Profile
	if err := r.db.GetContext(ctx, &p, query, args...); err != nil {
		return e.Wrap(err, "failed to update profile")
	}
	*profile = p.ToDomain()
	return nil
}

func (r *profileRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	query, args := r.qb.Select("TRUE").From("profiles").Where(sq.Eq{"username": username}).MustSql()
	var ex bool
	err := r.db.GetContext(ctx, &ex, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, e.Wrap(err, "failed to check username exists")
	}
	return ex, nil
}
