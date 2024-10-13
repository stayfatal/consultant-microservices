package transport

import (
	"cm/services/sso/internal/endpoints"
	"cm/services/sso/internal/interfaces"
	"cm/services/sso/internal/transport/pb"
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServer
	register kitgrpc.Handler
	login    kitgrpc.Handler
}

func NewGRPCServer(svc interfaces.Service) *AuthServer {
	ep := endpoints.MakeEndpoints(svc)

	return &AuthServer{
		register: kitgrpc.NewServer(
			ep.RegisterEndpoint,
			decodeRegisterRequest,
			encodeRegisterResponse,
		),
		login: kitgrpc.NewServer(
			ep.LoginEndpoint,
			decodeLoginRequest,
			encodeLoginResponse,
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
