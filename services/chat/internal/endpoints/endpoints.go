package endpoints

import (
	"cm/services/chat/internal/interfaces"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AddConsultantEndpoint endpoint.Endpoint
}

func MakeEndpoints(svc interfaces.Service) *Endpoints {
	return &Endpoints{
		makeAddConsultantEndpoint(),
	}
}

func makeAddConsultantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// code
		return
	}
}
