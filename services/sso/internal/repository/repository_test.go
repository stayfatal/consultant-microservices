package repository

import (
	"cm/services/sso/internal/models"
	"cm/services/sso/internal/testhelpers"
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	container, db, err := testhelpers.ConfigurePostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupPostgresContainer(t, container, db)

	repo := New(db)

	expected := models.User{
		Name:         "test",
		Email:        "test@gmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	id, err := repo.CreateUser(expected)
	if err != nil {
		t.Fatal(err)
	}

	got := models.User{}

	err = db.Get(&got, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		t.Fatal(err)
	}

	expected.Id = got.Id
	expected.CreatedAt = got.CreatedAt
	assert.Equal(t, expected, got)
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	container, db, err := testhelpers.ConfigurePostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupPostgresContainer(t, container, db)

	repo := New(db)

	expected := models.User{
		Name:         "test",
		Email:        "test@gmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	var (
		id         int
		created_at time.Time
	)
	err = db.QueryRow("INSERT INTO users (name,email,password,is_consultant) VALUES ($1,$2,$3,$4) RETURNING id, created_at", expected.Name, expected.Email, expected.Password, expected.IsConsultant).Scan(&id, &created_at)
	if err != nil {
		t.Fatal(err)
	}
	expected.Id = id
	expected.CreatedAt = created_at

	got, err := repo.GetUserByEmail(expected)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, got)
}
