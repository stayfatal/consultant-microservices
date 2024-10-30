package testhelpers

import (
	"cm/services/sso/config"
	"testing"

	"github.com/jmoiron/sqlx"
)

func PreparePostgres(t *testing.T) (*sqlx.Tx, error) {
	cfg, err := config.LoadPostgresConfig()
	if err != nil {
		return nil, err
	}

	db, err := config.NewPostgresDb(cfg)
	if err != nil {
		return nil, err
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil {
			t.Error(err)
		}

		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	})

	return tx, nil
}
