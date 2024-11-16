package transport

import (
	"cm/gen/authpb"
	"cm/internal/log"
	"cm/internal/middlewares"
	"cm/services/sso/internal/endpoints"
	"cm/services/sso/internal/interfaces"
	"context"

	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type serverApi struct {
	authpb.UnimplementedAuthenticationServer
	register kitgrpc.Handler
	login    kitgrpc.Handler
}

func NewGRPCServer(svc interfaces.Service, logger *log.Logger) authpb.AuthenticationServer {
	ep := endpoints.MakeEndpoints(svc)

	return &serverApi{
		register: kitgrpc.NewServer(
			middlewares.GrpcCustomChain(logger)(ep.RegisterEndpoint),
			decodeRegisterRequest,
			encodeRegisterResponse,
			kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
		login: kitgrpc.NewServer(
			middlewares.GrpcCustomChain(logger)(ep.LoginEndpoint),
			decodeLoginRequest,
			encodeLoginResponse,
			kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		),
	}
}

func (sa *serverApi) Register(ctx context.Context, request *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	_, resp, err := sa.register.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*authpb.RegisterResponse), nil
}

func (sa *serverApi) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	_, resp, err := sa.login.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*authpb.LoginResponse), nil
}
