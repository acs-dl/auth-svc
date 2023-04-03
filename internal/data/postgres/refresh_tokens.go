package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	refreshTokensTableName       = "refresh_tokens"
	refreshTokensTokenColumn     = refreshTokensTableName + ".token"
	refreshTokensValidTillColumn = refreshTokensTableName + ".valid_till"
)

type RefreshTokensQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
}

var selectedRefreshTokensTable = sq.Select("*").From(refreshTokensTableName)

func NewRefreshTokensQ(db *pgdb.DB) data.RefreshTokens {
	return &RefreshTokensQ{
		db:            db.Clone(),
		selectBuilder: selectedRefreshTokensTable,
		deleteBuilder: sq.Delete(refreshTokensTableName),
	}
}

func (q RefreshTokensQ) New() data.RefreshTokens {
	return NewRefreshTokensQ(q.db)
}

func (q RefreshTokensQ) Create(token data.RefreshToken) error {
	clauses := structs.Map(token)

	query := sq.Insert(refreshTokensTableName).SetMap(clauses)

	err := q.db.Exec(query)

	return err
}

func (q RefreshTokensQ) Select() ([]data.RefreshToken, error) {
	var result []data.RefreshToken

	err := q.db.Select(&result, q.selectBuilder)

	return result, err
}

func (q RefreshTokensQ) Get() (*data.RefreshToken, error) {
	var result data.RefreshToken

	err := q.db.Get(&result, q.selectBuilder)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q RefreshTokensQ) Delete() error {
	var deleted []data.RefreshToken

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		return errors.Errorf("no such token")
	}

	return nil
}

func (q RefreshTokensQ) FilterByTokens(tokens ...string) data.RefreshTokens {
	equalTokens := sq.Eq{refreshTokensTokenColumn: tokens}
	q.selectBuilder = q.selectBuilder.Where(equalTokens)
	q.deleteBuilder = q.deleteBuilder.Where(equalTokens)

	return q
}

func (q RefreshTokensQ) FilterByLowerValidTill(expiresAtUnix int64) data.RefreshTokens {
	lowerTime := sq.Lt{refreshTokensValidTillColumn: expiresAtUnix}
	q.selectBuilder = q.selectBuilder.Where(lowerTime)
	q.deleteBuilder = q.deleteBuilder.Where(lowerTime)

	return q
}
