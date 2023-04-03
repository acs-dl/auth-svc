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
	permissionsTableName      = "permissions"
	permissionsIdColumn       = permissionsTableName + ".id"
	permissionsModuleIdColumn = permissionsTableName + ".module_id"
	permissionsNameColumn     = permissionsTableName + ".name"
	permissionsStatusColumn   = permissionsTableName + ".status"
)

var (
	permissionsColumns = []string{
		permissionsIdColumn,
		permissionsModuleIdColumn,
		permissionsNameColumn,
		permissionsStatusColumn,
	}
	selectedPermissionsTable = sq.Select(permissionsColumns...).From(permissionsTableName)
)

type PermissionsQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
}

func NewPermissionsQ(db *pgdb.DB) data.Permissions {
	return &PermissionsQ{
		db:            db.Clone(),
		selectBuilder: selectedPermissionsTable,
		deleteBuilder: sq.Delete(permissionsTableName),
	}
}

func (q PermissionsQ) New() data.Permissions {
	return NewPermissionsQ(q.db)
}

func (q PermissionsQ) Upsert(permission data.Permission) error {
	clauses := structs.Map(permission)

	query := sq.Insert(permissionsTableName).SetMap(clauses).
		Suffix("ON CONFLICT (module_id, name, status) DO NOTHING")

	err := q.db.Exec(query)

	return err
}

func (q PermissionsQ) Select() ([]data.ModulePermission, error) {
	var result []data.ModulePermission

	q.selectBuilder = q.selectBuilder.GroupBy(permissionsColumns...)
	err := q.db.Select(&result, q.selectBuilder)

	return result, err
}

func (q PermissionsQ) Get() (*data.ModulePermission, error) {
	var result data.ModulePermission

	q.selectBuilder = q.selectBuilder.GroupBy(permissionsColumns...)
	err := q.db.Get(&result, q.selectBuilder)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q PermissionsQ) Delete() error {
	var deleted []data.Permission

	err := q.db.Select(&deleted, q.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		return errors.Errorf("no rows deleted")
	}

	return nil
}

func (q PermissionsQ) FilterByStatus(status data.UserStatus) data.Permissions {
	equalStatus := sq.Eq{permissionsStatusColumn: status}
	q.selectBuilder = q.selectBuilder.Where(equalStatus)
	q.deleteBuilder = q.deleteBuilder.Where(equalStatus)

	return q
}

func (q PermissionsQ) WithModules() data.Permissions {
	q.selectBuilder = sq.Select().Columns(permissionsIdColumn, permissionsModuleIdColumn).
		Column(permissionsNameColumn + " as permission_name").From(permissionsTableName)

	q.selectBuilder = q.selectBuilder.
		LeftJoin(modulesTableName + " ON " + permissionsModuleIdColumn + " = " + modulesIdColumn).
		Column(modulesIdColumn).
		Column(modulesNameColumn + " as module_name").
		GroupBy(modulesColumns...)

	return q
}
