package service

import (
	"cm/services/entities"
	"cm/services/sso/internal/auth"
	"cm/services/sso/internal/cache"
	"cm/services/sso/internal/repository"
	"cm/services/sso/internal/testhelpers"
	"context"
	"encoding/json"
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

	expected := entities.User{
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

	assert.Equal(t, claims.Id, expected.Id)

	gotFromPostgres := entities.User{}
	err = postgresDB.Get(&gotFromPostgres, "SELECT * FROM users WHERE id = $1", expected.Id)
	if err != nil {
		t.Fatal(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(gotFromPostgres.Password), []byte(expected.Password)); err != nil {
		t.Fatalf("expected and got passwords are not equal")
	}

	gotFromPostgres.Password = expected.Password
	gotFromPostgres.Id = expected.Id
	gotFromPostgres.CreatedAt = expected.CreatedAt

	assert.Equal(t, expected, gotFromPostgres)

	gotFromRedis := entities.User{}
	result, err := redisDB.Get(context.Background(), expected.Email).Result()
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal([]byte(result), &gotFromRedis)
	if err != nil {
		t.Fatal(err)
	}

	gotFromRedis.Password = expected.Password
	gotFromRedis.Id = expected.Id
	gotFromRedis.CreatedAt = expected.CreatedAt
	assert.Equal(t, expected, gotFromRedis)
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

	expected := entities.User{
		Name:         "test",
		Email:        "test@testmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	var (
		id int
	)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(expected.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = postgresDB.QueryRow("INSERT INTO users (name,email,password,is_consultant) VALUES ($1,$2,$3,$4) RETURNING id", expected.Name, expected.Email, hashedPass, expected.IsConsultant).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	expected.Id = id

	gotToken, err := svc.Login(expected)
	if err != nil {
		t.Fatal(err)
	}

	gotClaims, err := auth.ValidateToken(gotToken)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected.Id, gotClaims.Id)
}
