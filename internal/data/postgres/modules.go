package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	modulesTableName  = "modules"
	modulesIdColumn   = modulesTableName + ".id"
	modulesNameColumn = modulesTableName + ".name"
)

var modulesColumns = []string{
	modulesIdColumn,
	modulesNameColumn,
}

type ModulesQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
}

func NewModulesQ(db *pgdb.DB) data.Modules {
	return &ModulesQ{
		db:            db.Clone(),
		selectBuilder: sq.Select("*").From(modulesTableName),
		deleteBuilder: sq.Delete(modulesTableName),
	}
}

func (q ModulesQ) New() data.Modules {
	return NewModulesQ(q.db)
}

func (q ModulesQ) Insert(module data.Module) error {
	clauses := structs.Map(module)

	query := sq.Insert(modulesTableName).SetMap(clauses).Suffix("ON CONFLICT (name) DO NOTHING")

	err := q.db.Exec(query)

	return err
}

func (q ModulesQ) Select() ([]data.Module, error) {
	var result []data.Module

	err := q.db.Select(&result, q.selectBuilder)

	return result, err
}

func (q ModulesQ) FilterByNames(names ...string) data.Modules {
	equalNames := sq.Eq{modulesNameColumn: names}
	q.selectBuilder = q.selectBuilder.Where(equalNames)
	q.deleteBuilder = q.deleteBuilder.Where(equalNames)

	return q
}

func (q ModulesQ) Get() (*data.Module, error) {
	var result data.Module

	err := q.db.Get(&result, q.selectBuilder)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q ModulesQ) Delete() error {
	var deleted []data.Module

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}
