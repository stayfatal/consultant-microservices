package endpoints

import (
	"cm/services/sso/internal/interfaces"
	"cm/services/sso/internal/models"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
}

func MakeEndpoints(svc interfaces.Service) *Endpoints {
	return &Endpoints{
		RegisterEndpoint: makeRegisterEndpoint(svc),
		LoginEndpoint:    makeLoginEndpoint(svc),
	}
}

func makeRegisterEndpoint(svc interfaces.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.RegisterRequest)
		token, err := svc.Register(req.User)
		return models.RegisterResponse{Token: token, Error: err.Error()}, err
	}
}

func makeLoginEndpoint(svc interfaces.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.LoginRequest)
		token, err := svc.Login(req.User)
		return models.LoginResponse{Token: token, Error: err.Error()}, err
	}
}
