package middlewares

import (
	"cm/internal/log"

	"github.com/go-kit/kit/endpoint"
)

func GrpcCustomChain(logger *log.Logger) endpoint.Middleware {
	return endpoint.Chain(Recoverer(logger), GrpcLogger(logger))
}

func DefaultCustomChain(logger *log.Logger) endpoint.Middleware {
	return endpoint.Chain(Recoverer(logger))
}
