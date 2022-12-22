package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const modulesTableName = "modules"

var modulesColumns = []string{"modules.id", "modules.name"}

type ModulesQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedModulesTable = sq.Select("*").From(modulesTableName)

func NewModulesQ(db *pgdb.DB) data.Modules {
	return &ModulesQ{
		db:  db.Clone(),
		sql: selectedModulesTable,
	}
}

func (q *ModulesQ) New() data.Modules {
	return NewModulesQ(q.db)
}

func (q *ModulesQ) Create(module data.Module) (*data.Module, error) {
	clauses := structs.Map(module)

	query := sq.Insert(modulesTableName).SetMap(clauses).Suffix("ON CONFLICT DO NOTHING")

	err := q.db.Exec(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec insert")
	}

	var result data.Module
	err = q.db.Get(&result, selectedModulesTable.Where(sq.Eq{"name": module.Name}))

	return &result, err
}

func (q *ModulesQ) Select() ([]data.Module, error) {
	var result []data.Module

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *ModulesQ) Get() (*data.Module, error) {
	var result data.Module

	err := q.db.Get(&result, q.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *ModulesQ) Delete(moduleName string) error {
	query := sq.Delete(modulesTableName).Where(sq.Eq{"module_name": moduleName})

	result, err := q.db.ExecWithResult(query)
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.New("no such module")
	}

	return nil
}
