package transport

import (
	"cm/services/sso/internal/endpoints"
	"cm/services/sso/internal/interfaces"
	"cm/services/sso/internal/logger"
	"cm/services/sso/internal/middlewares"
	"cm/services/sso/internal/transport/pb"
	"context"

	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServer
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

func (s *AuthServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	_, resp, err := s.register.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.RegisterResponse), nil
}

func (s *AuthServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, resp, err := s.login.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LoginResponse), nil
}
