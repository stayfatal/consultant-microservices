package config

import (
	"cm/services/sso/utils"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	viperr "github.com/spf13/viper"
)

type PostgresConfig struct {
	IMAGE    string
	USER     string
	PASSWORD string
	DB_NAME  string
	SSL_MODE string
	PORT     int
	HOST     string
}

func LoadPostgresConfig() (*PostgresConfig, error) {
	path, err := utils.GetPath("services/sso/config/postgres_config.env")
	if err != nil {
		return nil, err
	}
	viperr.SetConfigFile(path)
	err = viperr.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "reading config")
	}

	cfg := &PostgresConfig{
		IMAGE:    viperr.GetString("IMAGE"),
		USER:     viperr.GetString("USER"),
		PASSWORD: viperr.GetString("PASSWORD"),
		DB_NAME:  viperr.GetString("DB_NAME"),
		SSL_MODE: viperr.GetString("SSL_MODE"),
		PORT:     viperr.GetInt("PORT"),
		HOST:     viperr.GetString("HOST"),
	}

	return cfg, nil
}

func NewPostgresDb(cfg *PostgresConfig) (*sqlx.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.HOST, cfg.PORT, cfg.USER, cfg.PASSWORD, cfg.DB_NAME, cfg.SSL_MODE)
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
