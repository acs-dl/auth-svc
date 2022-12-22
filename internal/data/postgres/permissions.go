package postgres

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const permissionsTableName = "permissions"

var permissionsColumns = []string{"permissions.id", "permissions.module_id", "permissions.name"}

type PermissionsQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedPermissionsTable = sq.Select(permissionsColumns...).From(permissionsTableName)

func NewPermissionsQ(db *pgdb.DB) data.Permissions {
	return &PermissionsQ{
		db:  db.Clone(),
		sql: selectedPermissionsTable,
	}
}

func (q *PermissionsQ) New() data.Permissions {
	return NewPermissionsQ(q.db)
}

func (q *PermissionsQ) Create(permission data.Permission) (*data.Permission, error) {
	clauses := structs.Map(permission)

	query := sq.Insert(permissionsTableName).SetMap(clauses).Suffix("ON CONFLICT DO NOTHING")

	err := q.db.Exec(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec insert")
	}

	var result data.Permission
	err = q.db.Get(&result, selectedPermissionsTable.Where(sq.Eq{"name": permission.Name}))

	return &result, err
}

func (q *PermissionsQ) Select() ([]data.ModulePermission, error) {
	var result []data.ModulePermission

	q.sql = q.sql.GroupBy(permissionsColumns...)
	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *PermissionsQ) Get() (*data.ModulePermission, error) {
	var result data.ModulePermission

	q.sql = q.sql.GroupBy(permissionsColumns...)
	err := q.db.Get(&result, q.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *PermissionsQ) Delete(permission data.Permission) error {
	query := sq.Delete(permissionsTableName).Where(
		sq.Eq{"name": permission.Name},
		sq.Eq{"module_id": permission.ModuleId})

	result, err := q.db.ExecWithResult(query)
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.New("no such permission")
	}

	return nil
}

func (q *PermissionsQ) FilterByModuleName(moduleName string) data.Permissions {
	q.sql = q.sql.Where(sq.Eq{fmt.Sprintf("%s.name", modulesTableName): moduleName})

	return q
}

func (q *PermissionsQ) FilterByPermissionId(permissionId int64) data.Permissions {
	q.sql = q.sql.Where(sq.Eq{fmt.Sprintf("%s.id", permissionsTableName): permissionId})

	return q
}

func (q *PermissionsQ) ResetFilters() data.Permissions {
	q.sql = selectedPermissionsTable

	return q
}

func (q *PermissionsQ) WithModules() data.Permissions {
	q.sql = sq.Select().Columns(fmt.Sprintf("%s.id", permissionsTableName), fmt.Sprintf("%s.module_id", permissionsTableName)).
		Column("permissions.name as permission_name").From(permissionsTableName)
	q.sql = q.sql.
		LeftJoin(
			fmt.Sprintf(
				"%s ON %s.module_id = %s.id ",
				modulesTableName, permissionsTableName, modulesTableName)).
		Column(fmt.Sprintf("%s.id", modulesTableName)).
		Column(fmt.Sprintf("%s.name as module_name", modulesTableName)).
		GroupBy(modulesColumns...)
	return q
}
