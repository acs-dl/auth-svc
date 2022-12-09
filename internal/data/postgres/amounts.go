package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/Auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const amountsTableName = "amounts"

type AmountsQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedAmountsTable = sq.Select("*").From(amountsTableName)

func NewAmountsQ(db *pgdb.DB) data.Amounts {
	return &AmountsQ{
		db:  db.Clone(),
		sql: selectedAmountsTable,
	}
}

func (q *AmountsQ) New() data.Amounts {
	return NewAmountsQ(q.db)
}

func (q *AmountsQ) Add(columnName string) error {
	err := q.db.ExecRaw(fmt.Sprintf("UPDATE %s SET %s = %s + 1", amountsTableName, columnName, columnName))

	return err
}

func (q *AmountsQ) Get() (*data.Amount, error) {
	var result data.Amount

	err := q.db.Get(&result, q.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}
