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
)

type subscriptionRepo struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func NewSubscriptionRepo(db *sqlx.DB) *subscriptionRepo {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &subscriptionRepo{db: db, qb: qb}
}

func (r *subscriptionRepo) SubscriptionsByUser(ctx context.Context, userID string) ([]domain.Subscription, error) {
	query, args := r.qb.
		Select("user_id", "email", "telegram_id", "type", "enabled").
		From("subscriptions").
		Join("users USING(user_id)").
		Where(sq.Eq{"subscriptions.user_id": userID}).
		MustSql()

	var subscriptions []SubscriptionWithUser
	if err := r.db.SelectContext(ctx, &subscriptions, query, args...); err != nil {
		return nil, e.Wrap(err, "failed to get subscriptions")
	}

	res := make([]domain.Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		res[i] = sub.ToDomain()
	}
	return res, nil
}

func (r *subscriptionRepo) Save(ctx context.Context, userID string, subType domain.SubscriptionType) error {
	query, args := r.qb.
		Insert("subscriptions").
		Columns("user_id", "type").
		Values(userID, subType).MustSql()

	_, err := r.execContext(ctx, query, args...)
	return e.WrapIfErr(err, "failed to save subscription")
}

func (r *subscriptionRepo) IsExists(ctx context.Context, userID string, subType domain.SubscriptionType) (bool, error) {
	query, args := r.qb.Select("TRUE").From("subscriptions").Where(sq.Eq{"user_id": userID, "type": subType}).MustSql()
	var exists bool
	err := r.getContext(ctx, &exists, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return exists, e.WrapIfErr(err, "failed to check subscription exists")
}

func (r *subscriptionRepo) Update(ctx context.Context, userID string, subType domain.SubscriptionType, enabled bool) error {
	query, args := r.qb.
		Update("subscriptions").
		Set("enabled", enabled).
		Where(sq.Eq{"user_id": userID, "type": subType}).
		MustSql()

	_, err := r.execContext(ctx, query, args...)
	return e.WrapIfErr(err, "failed to update subscription")
}

func (r *subscriptionRepo) Delete(ctx context.Context, userID string, subType domain.SubscriptionType) error {
	query, args := r.qb.
		Delete("subscriptions").
		Where(sq.Eq{"user_id": userID, "type": subType}).
		MustSql()

	_, err := r.execContext(ctx, query, args...)
	return e.WrapIfErr(err, "failed to delete subscription")
}

func (r *subscriptionRepo) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx := transaction.ExtractTx(ctx)
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return r.db.ExecContext(ctx, query, args...)
}

func (r *subscriptionRepo) getContext(ctx context.Context, dest any, query string, args ...any) error {
	tx := transaction.ExtractTx(ctx)
	if tx != nil {
		return tx.GetContext(ctx, dest, query, args...)
	}
	return r.db.GetContext(ctx, dest, query, args...)
}
