package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/acs-dl/auth-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	usersTableName   = "users"
	usersIdColumn    = usersTableName + ".id"
	usersEmailColumn = usersTableName + ".email"
)

type UsersQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	deleteBuilder sq.DeleteBuilder
}

var selectedUsersTable = sq.Select("*").From(usersTableName)

func NewUsersQ(db *pgdb.DB) data.Users {
	return &UsersQ{
		db:            db.Clone(),
		selectBuilder: selectedUsersTable,
		deleteBuilder: sq.Delete(usersTableName),
	}
}

func (q UsersQ) New() data.Users {
	return NewUsersQ(q.db)
}

func (q UsersQ) Get() (*data.User, error) {
	var result data.User

	err := q.db.Get(&result, q.selectBuilder)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q UsersQ) FilterByEmails(emails ...string) data.Users {
	equalEmails := sq.Eq{usersEmailColumn: emails}
	q.selectBuilder = q.selectBuilder.Where(equalEmails)
	q.deleteBuilder = q.deleteBuilder.Where(equalEmails)

	return q
}

func (q UsersQ) FilterByIds(ids ...int64) data.Users {
	equalIds := sq.Eq{usersIdColumn: ids}
	q.selectBuilder = q.selectBuilder.Where(equalIds)
	q.deleteBuilder = q.deleteBuilder.Where(equalIds)

	return q
}
