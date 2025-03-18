package middlewares

import (
	"cm/libs/log"

	"github.com/go-kit/kit/endpoint"
)

func GrpcCustomChain(logger *log.Logger) endpoint.Middleware {
	return endpoint.Chain(Recoverer(logger), GrpcLogger(logger))
}

func DefaultCustomChain(logger *log.Logger) endpoint.Middleware {
	return endpoint.Chain(Recoverer(logger))
}
