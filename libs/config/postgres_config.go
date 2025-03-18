package config

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sqlx.DB, error) {
	viper.AutomaticEnv()
	params := map[string]string{
		"USER":     "",
		"PASSWORD": "",
		"DB":       "",
		"HOST":     "",
		"PORT":     "",
		"SSL_MODE": "",
	}

	for key := range params {
		paramKey := fmt.Sprintf("POSTGRES_%s", key)

		err := viper.BindEnv(paramKey)
		if err != nil {
			return nil, errors.Wrap(err, "binding env var")
		}
		params[key] = viper.GetString(paramKey)
	}

	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			params["HOST"],
			params["PORT"],
			params["USER"],
			params["PASSWORD"],
			params["DB"],
			params["SSL_MODE"],
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "opening db")
	}

	return db, errors.Wrap(db.Ping(), "pinging db")
}
