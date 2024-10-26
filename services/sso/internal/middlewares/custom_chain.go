package middlewares

import (
	"cm/services/sso/internal/logger"

	"github.com/go-kit/kit/endpoint"
)

func CustomChain(logger *logger.Logger) endpoint.Middleware {
	return endpoint.Chain(Recoverer(logger), Logger(logger))
}
