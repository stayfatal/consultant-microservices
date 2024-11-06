package transport

import (
	"bytes"
	"cm/internal/entities"
	"cm/internal/log"
	"cm/internal/publicauth"
	"cm/services/gateway/http/config"
	"cm/services/gateway/http/internal/models"
	"cm/services/gateway/http/internal/service"
	transport "cm/services/gateway/http/internal/transport/http"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLogin(t *testing.T) {
	cfg, err := config.LoadServiceConfig()
	if err != nil {
		t.Fatal(err)
	}

	svc, err := service.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	logger := log.New()
	srv := transport.NewGatewayServer(svc, logger)

	go http.ListenAndServe(":3005", srv)

	client := &http.Client{Timeout: time.Second}

	expected := entities.User{
		Name:         "test",
		Email:        fmt.Sprintf("test%s@gmail.com", uuid.New().String()),
		Password:     "123",
		IsConsultant: false,
	}

	marshalledUser, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	regReq, err := http.NewRequest("POST", "http://localhost:3005/auth/register", bytes.NewBuffer(marshalledUser))
	if err != nil {
		t.Fatal(err)
	}

	regResp, err := client.Do(regReq)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, regResp)

	var regResponse models.RegistrationResponse
	err = json.NewDecoder(regResp.Body).Decode(&regResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, regResponse.Error, "")

	regClaims, err := publicauth.ValidateToken(regResponse.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, regClaims)
	assert.Equal(t, expected.Email, regClaims.Email)

	loginReq, err := http.NewRequest("GET", "http://localhost:3005/auth/login", bytes.NewBuffer(marshalledUser))
	if err != nil {
		t.Fatal(err)
	}

	loginResp, err := client.Do(loginReq)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, loginResp)

	var loginResponse models.LoginResponse
	err = json.NewDecoder(loginResp.Body).Decode(&loginResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, loginResponse.Error, "")

	loginClaims, err := publicauth.ValidateToken(loginResponse.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, loginClaims)
	assert.Equal(t, expected.Email, loginClaims.Email)
}
