package repo

import (
	sq "github.com/Masterminds/squirrel"
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
