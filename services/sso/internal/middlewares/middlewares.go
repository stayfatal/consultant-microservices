package middlewares

import (
	customlog "cm/services/sso/internal/logger"
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func Logger(logger *customlog.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			tStart := time.Now()
			method, _ := grpc.Method(ctx)

			resp, err := next(ctx, request)

			logger.Info().Str("method", method).Str("duration", time.Since(tStart).String()).Msg("Done")

			return resp, err
		}
	}
}

func Recoverer() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			resp, err := next(ctx, request)
			if err := recover(); err != nil {
				log.Error().Msgf("Recovered: %v", err)
			}
			return resp, err
		}
	}
}
