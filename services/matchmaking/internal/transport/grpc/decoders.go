package transport

import (
	"cm/gen/chatpb"
	"cm/libs/entities"
	"cm/services/matchmaking/internal/models"
	"context"

	"github.com/pkg/errors"
)

func decodeAddConsultantRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*chatpb.AddConsultantRequest)
	if !ok {
		return nil, errors.New("type assertion error")
	}
	return &models.AddConsultantRequest{
		User: entities.User{
			Id:    int(req.Id),
			Email: req.Email,
		},
	}, nil
}
