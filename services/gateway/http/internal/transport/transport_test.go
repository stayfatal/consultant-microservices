package transport

import (
	"bytes"
	"cm/internal/entities"
	"cm/internal/publicauth"
	"cm/services/gateway/http/internal/models"
	"cm/services/gateway/http/internal/service"
	transport "cm/services/gateway/http/internal/transport/http"
	"encoding/json"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var user = entities.User{
	Name:         "test",
	Email:        "transporttest@testmail.com",
	Password:     "123",
	IsConsultant: false,
}

func TestRegister(t *testing.T) {
	svc, err := service.New()
	if err != nil {
		t.Fatal(err)
	}

	srv := transport.NewGatewayServer(svc)
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go http.Serve(l, srv)

	client := &http.Client{Timeout: time.Second}

	marshalledUser, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "http://gatewayhttp:3000/auth/register", bytes.NewBuffer(marshalledUser))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, resp)

	var response models.RegistrationResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Error, "")

	claims, err := publicauth.ValidateToken(response.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, claims)
}

func TestLogin(t *testing.T) {
	svc, err := service.New()
	if err != nil {
		t.Fatal(err)
	}

	srv := transport.NewGatewayServer(svc)
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go http.Serve(l, srv)

	client := &http.Client{Timeout: time.Second}

	marshalledUser, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://gatewayhttp:3000/auth/login", bytes.NewBuffer(marshalledUser))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, resp)

	var response models.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Error, "")

	claims, err := publicauth.ValidateToken(response.Token)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, claims)
}
