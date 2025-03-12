package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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
	qb := r.qb.Insert("profiles").Columns("user_id", "username").Values(profile.UserID, profile.Username)

	if profile.FirstName != "" {
		qb = qb.Columns("first_name").Values(profile.FirstName)
	}
	if profile.Avatar != "" {
		qb = qb.Columns("avatar").Values(profile.Avatar)
	}

	query, args := qb.MustSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}
	return nil
}

func (r *profileRepo) ProfileByID(ctx context.Context, id string) (domain.Profile, error) {
	query, args := r.qb.Select("*").From("profiles").Where(sq.Eq{"user_id": id}).MustSql()
	var profile Profile
	if err := r.db.GetContext(ctx, &profile, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Profile{}, domain.ErrProfileNotFound
		}
		return domain.Profile{}, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile.ToDomain(), nil
}

func (r *profileRepo) Update(ctx context.Context, profile domain.Profile) error {
	query, args := r.qb.Update("profiles").
		Set("username", profile.Username).
		Set("first_name", profile.FirstName).
		Set("last_name", profile.LastName).
		Set("birth_date", profile.BirthDate).
		Set("gender", profile.Gender).Where(sq.Eq{"user_id": profile.UserID}).MustSql()
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
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
		return false, fmt.Errorf("failed to check username exists")
	}
	return ex, nil
}
