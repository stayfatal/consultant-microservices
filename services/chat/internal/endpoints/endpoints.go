package endpoints

import (
	"cm/services/chat/internal/interfaces"
	"cm/services/chat/internal/models"
	"context"

	"github.com/pkg/errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AddConsultantEndpoint endpoint.Endpoint
}

func MakeEndpoints(svc interfaces.Service) *Endpoints {
	return &Endpoints{
		makeAddConsultantEndpoint(svc),
	}
}

func makeAddConsultantEndpoint(svc interfaces.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*models.AddConsultantRequest)
		if !ok {
			err := errors.New("type assertion error")
			return &models.AddConsultantResponse{Error: err.Error()}, err
		}
		svc.AddConsultant(req.User)
		return &models.AddConsultantResponse{}, nil
	}
}
