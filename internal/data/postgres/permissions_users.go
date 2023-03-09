package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const PermissionUsersTableName = "permissions_users"

type PermissionUsersQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedPermissionUsersTable = sq.Select("*").From(PermissionUsersTableName)

func NewPermissionUsersQ(db *pgdb.DB) data.PermissionUsers {
	return &PermissionUsersQ{
		db:  db.Clone(),
		sql: selectedPermissionUsersTable,
	}
}

func (q *PermissionUsersQ) New() data.PermissionUsers {
	return NewPermissionUsersQ(q.db)
}

func (q *PermissionUsersQ) Create(module data.PermissionUser) (data.PermissionUser, error) {
	clauses := structs.Map(module)

	query := sq.Insert(PermissionUsersTableName).SetMap(clauses).Suffix("returning *")

	var result data.PermissionUser
	err := q.db.Get(&result, query)

	return result, err
}

func (q *PermissionUsersQ) Select() ([]data.PermissionUser, error) {
	var result []data.PermissionUser

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *PermissionUsersQ) Get() (*data.PermissionUser, error) {
	var result data.PermissionUser

	err := q.db.Get(&result, q.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *PermissionUsersQ) Delete(PermissionUser data.PermissionUser) error {
	query := sq.Delete(PermissionUsersTableName).Where(sq.And{
		sq.Eq{"permission_id": PermissionUser.PermissionId},
		sq.Eq{"user_id": PermissionUser.UserId},
	})

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

func (q *PermissionUsersQ) FilterByPermissionId(permissionId int64) data.PermissionUsers {
	q.sql = q.sql.Where(sq.Eq{"permission_id": permissionId})

	return q
}

func (q *PermissionUsersQ) FilterByUserId(userId int64) data.PermissionUsers {
	q.sql = q.sql.Where(sq.Eq{"user_id": userId})

	return q
}
