package repository

import (
	"cm/services/sso/config"
	"cm/services/sso/internal/models"
	"context"
	"log"
	"strconv"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func configureTest() (testcontainers.Container, *sqlx.DB, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.3-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "mypass",
			"POSTGRES_DB":       "prod_consultant_db",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	port, err := container.MappedPort(context.Background(), nat.Port("5432/tcp"))
	if err != nil {
		return nil, nil, err
	}
	log.Printf("***%s***\n", port.Port())

	cfg.POSTGRES_PORT, err = strconv.Atoi(port.Port())
	if err != nil {
		return nil, nil, err
	}
	// cfg.POSTGRES_PORT = 5432
	log.Printf("***%d***\n", cfg.POSTGRES_PORT)
	db, err := config.NewPostgresDb(cfg)
	if err != nil {
		return nil, nil, err
	}

	return container, db, nil
}

func gratefulTestStop(t *testing.T, container testcontainers.Container, db *sqlx.DB) {
	defer func() {
		err := testcontainers.TerminateContainer(container)
		if err != nil {
			t.Fatal(err)
		}
	}()
	err := db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {
	container, db, err := configureTest()
	if err != nil {
		t.Fatal(err)
	}
	defer gratefulTestStop(t, container, db)

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
	container, db, err := configureTest()
	if err != nil {
		t.Fatal(err)
	}
	defer gratefulTestStop(t, container, db)

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
