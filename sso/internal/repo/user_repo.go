package repo

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &userRepo{
		db: db,
		qb: qb,
	}
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query, args := r.qb.Select("*").From("users").Where(sq.Eq{"email": email}).MustSql()
	var user User
	if err := r.getContext(ctx, &user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user.ToDomain(), nil
}

func (r *userRepo) Create(ctx context.Context, email string) (domain.User, error) {
	query, args := r.qb.Insert("users").Columns("email").Values(email).Suffix("RETURNING *").MustSql()
	var user User
	if err := r.getContext(ctx, &user, query, args...); err != nil {
		return domain.User{}, err
	}
	return user.ToDomain(), nil
}

func (r *userRepo) AddAccount(ctx context.Context, userID uuid.UUID, provider domain.AccountType, password []byte) (domain.Account, error) {
	query, args := r.qb.
		Insert("accounts").
		Columns("user_id", "provider", "password").
		Values(userID, provider, password).
		Suffix("RETURNING *").MustSql()
	var account Account
	if err := r.getContext(ctx, &account, query, args...); err != nil {
		return domain.Account{}, err
	}
	return account.ToDomain(), nil
}

func (r *userRepo) GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	query, args := r.qb.Select("*").From("users").Where(sq.Eq{"user_id": userID}).MustSql()
	var user User
	if err := r.getContext(ctx, &user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user.ToDomain(), nil
}

func (r *userRepo) AccountByID(ctx context.Context, userID uuid.UUID, provider domain.AccountType) (domain.Account, error) {
	query, args := r.qb.Select("*").From("accounts").Where(sq.Eq{"user_id": userID, "provider": provider}).MustSql()
	var account Account
	if err := r.getContext(ctx, &account, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Account{}, domain.ErrAccountNotFound
		}
		return domain.Account{}, err
	}
	return account.ToDomain(), nil
}

func (r *userRepo) getContext(ctx context.Context, dest any, query string, args ...any) error {
	tx := transaction.ExtractTx(ctx)
	if tx != nil {
		return tx.GetContext(ctx, dest, query, args...)
	}
	return r.db.GetContext(ctx, dest, query, args...)
}
