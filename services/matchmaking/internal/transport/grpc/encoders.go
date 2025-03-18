package transport

import (
	"cm/gen/chatpb"
	"cm/services/matchmaking/internal/models"
	"context"

	"github.com/pkg/errors"
)

func encodeAddConsultantResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.AddConsultantResponse)
	if !ok {
		err := errors.New("type assertion error")
		return &chatpb.AddConsultantResponse{Error: err.Error()}, err
	}

	return &chatpb.AddConsultantResponse{Error: resp.Error}, nil
}
