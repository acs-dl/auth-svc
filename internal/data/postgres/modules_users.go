package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const modulesUsersTableName = "modules_users"

type ModulesUsersQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedModulesUsersTable = sq.Select("*").From(modulesUsersTableName)

func NewModulesUsersQ(db *pgdb.DB) data.ModulesUsers {
	return &ModulesUsersQ{
		db:  db.Clone(),
		sql: selectedModulesUsersTable,
	}
}

func (q *ModulesUsersQ) New() data.ModulesUsers {
	return NewModulesUsersQ(q.db)
}

func (q *ModulesUsersQ) Create(module data.ModuleUser) (data.ModuleUser, error) {
	clauses := structs.Map(module)

	query := sq.Insert(modulesUsersTableName).SetMap(clauses).Suffix("returning *")

	var result data.ModuleUser
	err := q.db.Get(&result, query)

	return result, err
}

func (q *ModulesUsersQ) Select() ([]data.ModuleUser, error) {
	var result []data.ModuleUser

	err := q.db.Select(&result, q.sql)

	return result, err
}

func (q *ModulesUsersQ) Delete(moduleUser data.ModuleUser) error {
	query := sq.Delete(modulesUsersTableName).Where(sq.And{
		sq.Eq{"module_id": moduleUser.ModuleId},
		sq.Eq{"user_id": moduleUser.UserId},
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

func (q *ModulesUsersQ) FilterByModuleId(moduleId int64) data.ModulesUsers {
	q.sql = q.sql.Where(sq.Eq{"module_id": moduleId})

	return q
}

func (q *ModulesUsersQ) FilterByUserId(userId int64) data.ModulesUsers {
	q.sql = q.sql.Where(sq.Eq{"user_id": userId})

	return q
}
