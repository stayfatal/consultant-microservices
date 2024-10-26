package service

import (
	"cm/services/sso/internal/auth"
	"cm/services/sso/internal/cache"
	"cm/services/sso/internal/models"
	"cm/services/sso/internal/repository"
	"cm/services/sso/internal/testhelpers"
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	posgresContainer, postgresDB, err := testhelpers.ConfigurePostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupPostgresContainer(t, posgresContainer, postgresDB)

	redisContainer, redisDB, err := testhelpers.ConfigureRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupRedisContainer(t, redisContainer, redisDB)

	cache := cache.New(redisDB)
	repo := repository.New(postgresDB)

	svc := New(repo, cache)

	expected := models.User{
		Id:           1,
		Name:         "test",
		Email:        "test@testmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	token, err := svc.Register(expected)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if claims.Id != expected.Id {
		t.Fatalf("expected token id  %d got %d", expected.Id, claims.Id)
	}

	got := models.User{}
	err = postgresDB.Get(&got, "SELECT * FROM users WHERE id = $1", expected.Id)
	if err != nil {
		t.Fatal(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(expected.Password)); err != nil {
		t.Fatalf("expected and got passwords are not equal")
	}

	got.Password = expected.Password
	got.Id = expected.Id
	got.CreatedAt = expected.CreatedAt

	assert.Equal(t, expected, got)
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	postgresContainer, postgresDB, err := testhelpers.ConfigurePostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupPostgresContainer(t, postgresContainer, postgresDB)

	redisContainer, redisDB, err := testhelpers.ConfigureRedisContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer testhelpers.CleanupRedisContainer(t, redisContainer, redisDB)

	cache := cache.New(redisDB)

	repo := repository.New(postgresDB)

	svc := New(repo, cache)

	user := models.User{
		Id:           1,
		Name:         "test",
		Email:        "test@testmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	expToken, err := svc.Register(user)
	if err != nil {
		t.Fatal(err)
	}

	expClaims, err := auth.ValidateToken(expToken)
	if err != nil {
		t.Fatal(err)
	}

	gotToken, err := svc.Login(user)
	if err != nil {
		t.Fatal(err)
	}

	gotClaims, err := auth.ValidateToken(gotToken)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expClaims.Id, gotClaims.Id)
}
