package transport

import (
	"cm/gen/authpb"
	"cm/services/sso/internal/endpoints"
	"cm/services/sso/internal/interfaces"
	"cm/services/sso/internal/logger"
	"cm/services/sso/internal/middlewares"
	"context"

	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type AuthServer struct {
	authpb.UnimplementedAuthenticationServer
	register kitgrpc.Handler
	login    kitgrpc.Handler
}

func NewGRPCServer(svc interfaces.Service, logger *logger.Logger) *AuthServer {
	ep := endpoints.MakeEndpoints(svc)

	return &AuthServer{
		register: kitgrpc.NewServer(
			middlewares.CustomChain(logger)(ep.RegisterEndpoint),
			decodeRegisterRequest,
			encodeRegisterResponse,
			kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
		login: kitgrpc.NewServer(
			middlewares.CustomChain(logger)(ep.LoginEndpoint),
			decodeLoginRequest,
			encodeLoginResponse,
			kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
	}
}

func (s *AuthServer) Register(ctx context.Context, request *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	_, resp, err := s.register.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*authpb.RegisterResponse), nil
}

func (s *AuthServer) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*authpb.LoginResponse), nil
}
