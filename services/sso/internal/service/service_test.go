package service

import (
	"cm/internal/entities"
	"cm/internal/publicauth"

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

var expected = entities.User{
	Name:         "test",
	Email:        "test@testmail.com",
	Password:     "123",
	IsConsultant: false,
}

func TestRegister(t *testing.T) {
	postgresDB, err := testhelpers.PreparePostgres(t)
	if err != nil {
		t.Fatal(err)
	}

	redisDB, err := testhelpers.PrepareRedis(t)
	if err != nil {
		t.Fatal(err)
	}

	cache := cache.New(redisDB)
	repo := repository.New(postgresDB)

	svc := New(repo, cache)

	token, err := svc.Register(expected)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := publicauth.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, claims)

	gotFromPostgres := entities.User{}
	err = postgresDB.Get(&gotFromPostgres, "SELECT * FROM users WHERE email = $1", expected.Email)
	if err != nil {
		t.Fatal(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(gotFromPostgres.Password), []byte(expected.Password)); err != nil {
		t.Fatalf("expected and got passwords are not equal")
	}

	assert.Equal(t, expected.Email, gotFromPostgres.Email)

	gotFromRedis := entities.User{}
	result, err := redisDB.Get(context.Background(), expected.Email).Result()
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal([]byte(result), &gotFromRedis)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expected.Email, gotFromRedis.Email)
}

func TestLogin(t *testing.T) {
	postgresDB, err := testhelpers.PreparePostgres(t)
	if err != nil {
		t.Fatal(err)
	}

	redisDB, err := testhelpers.PrepareRedis(t)
	if err != nil {
		t.Fatal(err)
	}

	cache := cache.New(redisDB)

	repo := repository.New(postgresDB)

	svc := New(repo, cache)

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

	gotClaims, err := publicauth.ValidateToken(gotToken)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, gotClaims)
	// assert.Equal(t, expected.Id, gotClaims.Id)
}
