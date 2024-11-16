package transport

import (
	"cm/gen/chatpb"
	"cm/services/chat/internal/endpoints"
	"cm/services/chat/internal/interfaces"
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type serverApi struct {
	chatpb.UnimplementedChatServer
	addConsultant kitgrpc.Handler
}

func NewGrpcServer(svc interfaces.Service) chatpb.ChatServer {
	eps := endpoints.MakeEndpoints(svc)
	return &serverApi{addConsultant: kitgrpc.NewServer(
		eps.AddConsultantEndpoint,
		nil,
		nil,
	)}
}

func (sa *serverApi) AddConsultant(ctx context.Context, request *chatpb.AddConsultantRequest) (*chatpb.AddConsultantResponse, error) {
	_, resp, err := sa.addConsultant.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*chatpb.AddConsultantResponse), nil
}
