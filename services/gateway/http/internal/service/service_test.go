package service

import (
	"cm/internal/entities"
	"cm/internal/publicauth"
	"testing"

	"github.com/stretchr/testify/assert"
)

var user = entities.User{
	Name:         "test",
	Email:        "servicetest@testmail.com",
	Password:     "123",
	IsConsultant: false,
}

func TestRegister(t *testing.T) {
	svc, err := New()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := svc.Register(user)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, resp)

	assert.Equal(t, resp.Error, "")

	claims, err := publicauth.ValidateToken(resp.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, claims)
}

func TestLogin(t *testing.T) {
	svc, err := New()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := svc.Login(user)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, resp)

	assert.Equal(t, resp.Error, "")

	claims, err := publicauth.ValidateToken(resp.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, claims)
}
