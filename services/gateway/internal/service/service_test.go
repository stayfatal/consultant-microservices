package service

import (
	"cm/libs/config"
	"cm/libs/entities"
	"cm/libs/publicauth"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLogin(t *testing.T) {
	cfg, err := config.LoadConfigs()
	if err != nil {
		t.Fatal(err)
	}

	svc, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	expected := entities.User{
		Name:         "test",
		Email:        fmt.Sprintf("test%s@gmail.com", uuid.New().String()),
		Password:     "123",
		IsConsultant: false,
	}

	regResp, err := svc.Register(expected)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, regResp)

	assert.Equal(t, regResp.Error, "")

	regClaims, err := publicauth.ValidateToken(regResp.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, regClaims)
	assert.Equal(t, expected.Email, regClaims.Email)

	loginResp, err := svc.Login(expected)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, loginResp)

	assert.Equal(t, loginResp.Error, "")

	loginClaims, err := publicauth.ValidateToken(loginResp.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, regClaims)
	assert.Equal(t, expected.Email, loginClaims.Email)
}
