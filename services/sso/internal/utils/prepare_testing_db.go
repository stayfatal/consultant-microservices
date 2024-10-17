package utils

import "github.com/jmoiron/sqlx"

func PrepareTestingDB() (*sqlx.DB, error) {
	conn := "user=postgres port=80 password=mypass dbname=prod_consultant_testing_db sslmode=disable"

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

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

	return db, nil
}

func ClearTestingDB(db *sqlx.DB) error {
	defer db.Close()
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		return err
	}
	return nil
}
