package endpoints

import (
	"cm/services/gateway/http/internal/interfaces"
	"cm/services/gateway/http/internal/models"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Register endpoint.Endpoint
	Login    endpoint.Endpoint
}

func MakeEndpoints(svc interfaces.Service) *Endpoints {
	return &Endpoints{
		Register: makeRegisterEndpoint(svc),
		Login:    makeLoginEndpoint(svc),
	}
}

func makeRegisterEndpoint(svc interfaces.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.RegistrationRequest)
		return svc.Register(req.User)
	}
}

func makeLoginEndpoint(svc interfaces.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(models.LoginRequest)
		return svc.Login(req.User)
	}
}
