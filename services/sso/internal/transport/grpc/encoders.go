package transport

import (
	"cm/services/sso/internal/models"
	"cm/services/sso/internal/transport/pb"
	"context"
)

func encodeRegisterResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(models.RegisterResponse)
	return &pb.RegisterResponse{
		Token: resp.Token,
		Error: resp.Error,
	}, nil
}

func encodeLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(models.LoginResponse)
	return &pb.LoginResponse{
		Token: resp.Token,
		Error: resp.Error,
	}, nil
}
