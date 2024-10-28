package interfaces

import (
	"cm/services/entities"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Service interface {
	Register(user entities.User) (string, error)
	Login(user entities.User) (string, error)
}

type Repository interface {
	CreateUser(user entities.User) (int, error)
	GetUserByEmail(user entities.User) (entities.User, error)
}

type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

type CacheDB interface {
	SetUser(user entities.User) error
	GetUser(user entities.User) (entities.User, error)
}
