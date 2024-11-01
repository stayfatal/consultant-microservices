package transport

import (
	"cm/internal/log"
	"cm/internal/middlewares"
	"cm/services/gateway/http/internal/endpoints"
	"cm/services/gateway/http/internal/interfaces"
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewGatewayServer(svc interfaces.Service, logger *log.Logger) *mux.Router {
	ep := endpoints.MakeEndpoints(svc)

	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeErrorResponse),
		kithttp.ServerFinalizer(makeFinalizer(logger)),
	}

	r.Handle("/auth/register", kithttp.NewServer(
		middlewares.DefaultCustomChain(logger)(ep.Register),
		decodeRegisterRequest,
		encodeRegisterResponse,
		options...,
	)).Methods("POST")

	r.Handle("/auth/login", kithttp.NewServer(
		middlewares.DefaultCustomChain(logger)(ep.Login),
		decodeLoginRequest,
		encodeLoginResponse,
		options...,
	)).Methods("GET")

	return r
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
}

func makeFinalizer(logger *log.Logger) kithttp.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		logger.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("ip", r.RemoteAddr).
			Str("Status code", fmt.Sprintf("%d", code)).
			Msg("")
	}
}
