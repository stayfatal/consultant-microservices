package transport

import (
	"cm/services/sso/internal/models"
	"cm/services/sso/internal/transport/pb"
	"context"
)

func decodeRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.RegisterRequest)
	return models.RegisterRequest{
		User: models.User{
			Name:         req.Name,
			Email:        req.Email,
			Password:     req.Password,
			IsConsultant: req.IsConsultant,
		},
	}, nil
}

func decodeLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LoginRequest)
	return models.LoginRequest{
		User: models.User{
			Email:    req.Email,
			Password: req.Password,
		},
	}, nil
}
