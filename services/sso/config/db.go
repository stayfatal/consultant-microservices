package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func NewDb(cfg Config) (*sqlx.DB, error) {
	conn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSslMode)

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, errors.Wrap(err, "trying to conn to db")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "pinging db")
	}

	return db, nil
}
