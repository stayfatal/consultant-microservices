package service

import (
	"cm/services/sso/internal/auth"
	"cm/services/sso/internal/models"
	"cm/services/sso/internal/repository"
	"cm/services/sso/internal/utils"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	db, err := utils.PrepareTestingDB()
	if err != nil {
		t.Fatalf("cant prepare db for tests %v", err)
	}
	defer utils.ClearTestingDB(db)

	repo := repository.New(db)

	svc := New(repo)

	expected := models.User{
		Id:           1,
		Name:         "test",
		Email:        "test@testmail.com",
		Password:     "123",
		IsConsultant: false,
	}

	token, err := svc.Register(expected)
	if err != nil {
		t.Error(err)
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Error(err)
	}

	if claims.Id != expected.Id {
		t.Errorf("expected token id  %d got %d", expected.Id, claims.Id)
	}

	got := models.User{}
	err = db.Get(&got, "SELECT * FROM users WHERE id = $1", expected.Id)
	if err != nil {
		t.Error(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(expected.Password)); err != nil {
		t.Errorf("expected and got passwords are not equal")
	}

	got.Password = expected.Password
	got.Id = expected.Id
	got.CreatedAt = expected.CreatedAt

	assert.Equal(t, expected, got)
}

func TestLogin(t *testing.T) {
	db, err := utils.PrepareTestingDB()
	if err != nil {
		t.Fatalf("cant prepare db for tests %v", err)
	}
	defer utils.ClearTestingDB(db)

	repo := repository.New(db)

	svc := New(repo)

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

	assert.Equal(t, expToken, gotToken)

	assert.Equal(t, expClaims.Id, gotClaims.Id)
}
