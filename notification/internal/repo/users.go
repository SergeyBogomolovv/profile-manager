package repo

import (
	sq "github.com/Masterminds/squirrel"
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
