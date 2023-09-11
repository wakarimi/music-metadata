package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Queryer interface {
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Queryx(query string, arg ...interface{}) (*sqlx.Rows, error)
}
