package dbtest

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func PrepareTestingDB() (*sqlx.DB, *sqlx.Tx, error) {
	conn := "user=postgres port=80 password=mypass dbname=prod_consultant_testing_db sslmode=disable"

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, nil, err
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, nil, err
	}

	table := `CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(225) NOT NULL,
    is_consultant BOOLEAN,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = tx.Exec(table)
	if err != nil {
		return nil, nil, err
	}

	return db, tx, nil
}

func ClearTestingDB(t *testing.T, db *sqlx.DB, tx *sqlx.Tx) {
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := tx.Rollback(); err != nil {
		t.Error(err)
	}
}
