package transport

import (
	"cm/services/gateway/internal/models"
	"context"
	"encoding/json"
	"net/http"
)

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req.User)
	return req, err
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req.User)
	return req, err
}
