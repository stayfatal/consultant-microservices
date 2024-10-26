package repository

import (
	"cm/services/sso/internal/models"
	"cm/services/sso/internal/testhelpers"
	"context"
	"testing"

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

	_, err = repo.CreateUser(expected)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetUserByEmail(expected)
	if err != nil {
		t.Fatal(err)
	}

	expected.Id = got.Id
	expected.CreatedAt = got.CreatedAt
	assert.Equal(t, expected, got)
}
