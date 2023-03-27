package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const refreshTokensTableName = "refresh_tokens"

type RefreshTokensQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedRefreshTokensTable = sq.Select("*").From(refreshTokensTableName)

func NewRefreshTokensQ(db *pgdb.DB) data.RefreshTokens {
	return &RefreshTokensQ{
		db:  db.Clone(),
		sql: selectedRefreshTokensTable,
	}
}

func (q *RefreshTokensQ) New() data.RefreshTokens {
	return NewRefreshTokensQ(q.db)
}

func (q *RefreshTokensQ) Create(token data.RefreshToken) error {
	clauses := structs.Map(token)

	query := sq.Insert(refreshTokensTableName).SetMap(clauses)

	err := q.db.Exec(query)

	return err
}

func (q *RefreshTokensQ) Select() ([]data.RefreshToken, error) {
	var result []data.RefreshToken

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *RefreshTokensQ) Get() (*data.RefreshToken, error) {
	var result data.RefreshToken

	err := q.db.Get(&result, q.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *RefreshTokensQ) Delete(token string) error {
	var deleted []data.RefreshToken

	query := sq.Delete(refreshTokensTableName).
		Where(sq.Eq{
			"token": token,
		}).
		Suffix("RETURNING *")

	err := q.db.Select(&deleted, query)
	if err != nil {
		return err
	}
	if len(deleted) == 0 {
		return errors.Errorf("no such token")
	}

	return nil
}

func (q *RefreshTokensQ) FilterByToken(token string) data.RefreshTokens {
	q.sql = q.sql.Where(sq.Eq{"token": token})

	return q
}

func (q *RefreshTokensQ) FilterByValidTill(expiresAtUnix int64) data.RefreshTokens {
	q.sql = q.sql.Where(sq.Lt{"valid_till": expiresAtUnix})

	return q
}
