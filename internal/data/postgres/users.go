package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/acs/auth/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const usersTableName = "users"

type UsersQ struct {
	db  *pgdb.DB
	sql sq.SelectBuilder
}

var selectedUsersTable = sq.Select("*").From(usersTableName)

func NewUsersQ(db *pgdb.DB) data.Users {
	return &UsersQ{
		db:  db.Clone(),
		sql: selectedUsersTable,
	}
}

func (q *UsersQ) New() data.Users {
	return NewUsersQ(q.db)
}

func (q *UsersQ) Get() (*data.User, error) {
	var result data.User

	err := q.db.Get(&result, q.sql)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *UsersQ) FilterByEmail(email string) data.Users {
	q.sql = q.sql.Where(sq.Eq{"email": email})

	return q
}

func (q *UsersQ) FilterById(id int64) data.Users {
	q.sql = q.sql.Where(sq.Eq{"id": id})

	return q
}
