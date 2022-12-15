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

func (q *ModulesQ) Create(module data.ModulePermission) (data.ModulePermission, error) {
	clauses := structs.Map(module)

	query := sq.Insert(modulesTableName).SetMap(clauses).Suffix("returning *")

	var result data.ModulePermission
	err := q.db.Get(&result, query)

	return result, err
}

func (q *ModulesQ) Select() ([]data.ModulePermission, error) {
	var result []data.ModulePermission

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *ModulesQ) Get() (*data.ModulePermission, error) {
	var result data.ModulePermission

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

func (q *ModulesQ) FilterByModule(moduleName string) data.Modules {
	q.sql = q.sql.Where(sq.Eq{"module_name": moduleName})

	return q
}

func (q *ModulesQ) FilterById(id int64) data.Modules {
	q.sql = q.sql.Where(sq.Eq{"id": id})

	return q
}

func (q *ModulesQ) ResetFilters() data.Modules {
	q.sql = selectedModulesTable

	return q
}
