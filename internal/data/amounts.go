package data

type Amounts interface {
	New() Amounts

	Add(columnName string) error
	Get() (*Amount, error)
}

type Amount struct {
	access  int64 `db:"access" structs:"access"`
	refresh int64 `db:"refresh" structs:"refresh"`
}
