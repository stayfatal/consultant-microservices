package repository

import (
	"cm/services/sso/internal/models"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, err := sqlx.Open("sqlite3", ":memory:?_busy_timeout=5000&cache=shared")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	table := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(225) NOT NULL,
		is_consultant BOOLEAN,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = tx.Exec(table)
	if err != nil {
		t.Error(err)
	}

	repo := New(tx)

	expected := models.User{
		Name:         "test",
		Email:        "test@gmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	id, err := repo.CreateUser(expected)
	if err != nil {
		t.Error(err)
	}

	got := models.User{}

	err = tx.Get(&got, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		t.Error(err)
	}

	expected.Id = got.Id
	expected.CreatedAt = got.CreatedAt
	if ok := assert.Equal(t, expected, got); !ok {
		t.Errorf("expected %v got %v", expected, got)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, err := sqlx.Open("sqlite3", ":memory:?_busy_timeout=5000&cache=shared")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	tx, err := db.Beginx()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()

	table := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(225) NOT NULL,
		is_consultant BOOLEAN,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = tx.Exec(table)
	if err != nil {
		t.Error(err)
	}

	repo := New(tx)

	expected := models.User{
		Name:         "test",
		Email:        "test@gmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	_, err = repo.CreateUser(expected)
	if err != nil {
		t.Error(err)
	}

	got, err := repo.GetUserByEmail(expected)
	if err != nil {
		t.Error(err)
	}

	expected.Id = got.Id
	expected.CreatedAt = got.CreatedAt
	if ok := assert.Equal(t, expected, got); !ok {
		t.Errorf("expected %v got %v", expected, got)
	}
}
