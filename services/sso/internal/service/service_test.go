package service

import (
	"cm/internal/entities"
	"cm/internal/publicauth"
	"fmt"

	"cm/services/sso/internal/cache"
	"cm/services/sso/internal/repository"
	"cm/services/sso/internal/testhelpers"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

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

	expected := entities.User{
		Name:         "test",
		Email:        fmt.Sprintf("test%s@gmail.com", uuid.New().String()),
		Password:     "123",
		IsConsultant: false,
	}

	token, err := svc.Register(expected)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := publicauth.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, claims)
	assert.Equal(t, expected.Email, claims.Email)

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

	expected := entities.User{
		Name:         "test",
		Email:        fmt.Sprintf("test%s@gmail.com", uuid.New().String()),
		Password:     "123",
		IsConsultant: false,
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(expected.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	_, err = postgresDB.Exec("INSERT INTO users (name,email,password,is_consultant) VALUES ($1,$2,$3,$4)", expected.Name, expected.Email, hashedPass, expected.IsConsultant)
	if err != nil {
		t.Fatal(err)
	}

	gotToken, err := svc.Login(expected)
	if err != nil {
		t.Fatal(err)
	}

	gotClaims, err := publicauth.ValidateToken(gotToken)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, gotClaims)
	assert.Equal(t, expected.Email, gotClaims.Email)
}
