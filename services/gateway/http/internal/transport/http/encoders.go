package transport

import (
	"cm/services/gateway/http/internal/models"
	"cm/services/gen/authpb"
	"context"
	"encoding/json"
	"net/http"
)

func encodeRegisterResponse(_ context.Context, w http.ResponseWriter, i interface{}) error {
	resp := i.(*authpb.RegisterResponse)
	response := models.RegistrationResponse{
		Token: resp.Token,
		Error: resp.Error,
	}
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, i interface{}) error {
	resp := i.(*authpb.LoginResponse)
	response := models.LoginResponse{
		Token: resp.Token,
		Error: resp.Error,
	}
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
