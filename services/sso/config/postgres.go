package config

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func NewPostgresDb(cfg *Config) (*sqlx.DB, error) {
	conn := fmt.Sprintf("host=host.docker.internal port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.POSTGRES_PORT, cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_NAME, cfg.POSTGRES_SSL_MODE)
	log.Println(conn)
	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, errors.Wrap(err, "trying to conn to db")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "pinging db")
	}

	//need to be deleted asap
	table := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(225) NOT NULL,
		is_consultant BOOLEAN,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(table)
	if err != nil {
		return nil, err
	}
	// ...

	return db, nil
}
