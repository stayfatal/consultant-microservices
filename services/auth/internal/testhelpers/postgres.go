package testhelpers

import (
	"cm/libs/config"
	"testing"

	"github.com/jmoiron/sqlx"
)

func PreparePostgres(t *testing.T) (*sqlx.Tx, error) {
	db, err := config.NewPostgresDB()
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
