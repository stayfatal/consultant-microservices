package transport

import (
	"cm/services/gateway/http/internal/endpoints"
	"cm/services/gateway/http/internal/interfaces"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewGatewayServer(svc interfaces.Service) *mux.Router {
	ep := endpoints.MakeEndpoints(svc)

	r := mux.NewRouter()

	r.Handle("/auth/register", kithttp.NewServer(
		ep.Register,
		decodeRegisterRequest,
		encodeRegisterResponse,
	)).Methods("POST")

	r.Handle("/auth/login", kithttp.NewServer(
		ep.Login,
		decodeLoginRequest,
		encodeLoginResponse,
	)).Methods("GET")

	return r
}
