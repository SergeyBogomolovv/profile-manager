package repo

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/SergeyBogomolovv/profile-manager/common/e"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type userRepo struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &userRepo{db: db, qb: qb}
}

func (r *userRepo) Save(ctx context.Context, user domain.User) error {
	m := map[string]any{
		"user_id": user.ID,
	}
	if user.Email != "" {
		m["email"] = user.Email
	}
	if user.TelegramID != 0 {
		m["telegram_id"] = user.TelegramID
	}
	query, args := r.qb.Insert("users").SetMap(m).MustSql()

	_, err := r.execContext(ctx, query, args...)
	return e.WrapIfErr(err, "failed to save user")
}

func (r *userRepo) IsExists(ctx context.Context, userID string) (bool, error) {
	query, args := r.qb.Select("TRUE").From("users").Where(sq.Eq{"user_id": userID}).MustSql()
	var exists bool
	err := r.getContext(ctx, &exists, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return exists, e.WrapIfErr(err, "failed to check user exists")
}

func (r *userRepo) GetByID(ctx context.Context, userID string) (domain.User, error) {
	query, args := r.qb.Select("*").From("users").Where(sq.Eq{"user_id": userID}).MustSql()
	var user User
	err := r.getContext(ctx, &user, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user.ToDomain(), e.WrapIfErr(err, "failed to get user")
}

func (r *userRepo) GetByTelegramID(ctx context.Context, telegramID int64) (domain.User, error) {
	query, args := r.qb.Select("*").From("users").Where(sq.Eq{"telegram_id": telegramID}).MustSql()
	var user User
	err := r.getContext(ctx, &user, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user.ToDomain(), e.WrapIfErr(err, "failed to get user")
}

func (r *userRepo) Update(ctx context.Context, user domain.User) error {
	m := map[string]any{"email": user.Email, "telegram_id": user.TelegramID}
	if user.Email == "" {
		m["email"] = nil
	}
	if user.TelegramID == 0 {
		m["telegram_id"] = nil
	}
	query, args := r.qb.Update("users").SetMap(m).Where(sq.Eq{"user_id": user.ID}).MustSql()
	_, err := r.execContext(ctx, query, args...)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
		return domain.ErrAccountAlreadyExists
	}
	return e.WrapIfErr(err, "failed to update user")
}

func (r *userRepo) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx := transaction.ExtractTx(ctx)
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return r.db.ExecContext(ctx, query, args...)
}

func (r *userRepo) getContext(ctx context.Context, dest any, query string, args ...any) error {
	tx := transaction.ExtractTx(ctx)
	if tx != nil {
		return tx.GetContext(ctx, dest, query, args...)
	}
	return r.db.GetContext(ctx, dest, query, args...)
}
