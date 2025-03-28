package transport

import (
	"cm/gen/chatpb"
	"cm/libs/log"
	"cm/services/matchmaking/internal/endpoints"
	"cm/services/matchmaking/internal/interfaces"
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type serverApi struct {
	chatpb.UnimplementedChatServer
	addConsultant kitgrpc.Handler
}

func NewGrpcServer(svc interfaces.Service, logger *log.Logger) chatpb.ChatServer {
	eps := endpoints.MakeEndpoints(svc)
	return &serverApi{addConsultant: kitgrpc.NewServer(
		eps.AddConsultantEndpoint,
		decodeAddConsultantRequest,
		encodeAddConsultantResponse,
		kitgrpc.ServerErrorLogger(logger),
	)}
}

func (sa *serverApi) AddConsultant(ctx context.Context, request *chatpb.AddConsultantRequest) (*chatpb.AddConsultantResponse, error) {
	_, resp, err := sa.addConsultant.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*chatpb.AddConsultantResponse), nil
}
