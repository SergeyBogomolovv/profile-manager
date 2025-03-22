package repo

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/notification/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func New(db *sqlx.DB) *userRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &userRepo{db: db, qb: qb}
}

func (r *userRepo) Subscriptions(ctx context.Context, userID string) ([]domain.Subscription, error) {
	query, args := r.qb.
		Select("user_id", "email", "telegram_id", "type", "enabled").
		From("subscriptions").
		Join("users USING(user_id)").
		Where(sq.Eq{"subscriptions.user_id": userID}).
		MustSql()

	var subscriptions []SubscriptionWithUser
	if err := r.db.SelectContext(ctx, &subscriptions, query, args...); err != nil {
		return nil, err
	}

	res := make([]domain.Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		res[i] = sub.ToDomain()
	}
	return res, nil
}

func (r *userRepo) SaveUser(ctx context.Context, user domain.User) error {
	m := map[string]any{
		"user_id": user.ID,
	}
	if user.Email != "" {
		m["email"] = user.Email
	}
	if user.TelegramID != 0 {
		m["telegram_id"] = user.TelegramID
	}
	q := r.qb.Insert("users").SetMap(m)
	query, args := q.MustSql()
	_, err := r.execContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (r *userRepo) SaveSubscription(ctx context.Context, userID string, subType domain.SubscriptionType) error {
	query, args := r.qb.
		Insert("subscriptions").
		Columns("user_id", "type").
		Values(userID, subType).MustSql()

	_, err := r.execContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to save subscription: %w", err)
	}
	return nil
}

func (r *userRepo) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx := transaction.ExtractTx(ctx)
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return r.db.ExecContext(ctx, query, args...)
}
