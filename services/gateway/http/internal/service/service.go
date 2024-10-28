package service

import (
	"cm/services/entities"
	"cm/services/gateway/http/internal/interfaces"
	"cm/services/gen/authpb"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type service struct {
	client authpb.AuthenticationClient
}

func New() (interfaces.Service, error) {
	client, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &service{client: authpb.NewAuthenticationClient(client)}, nil
}

func (s *service) Register(user entities.User) (*authpb.RegisterResponse, error) {
	return s.client.Register(context.Background(), &authpb.RegisterRequest{
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		IsConsultant: user.IsConsultant,
	})
}

func (s *service) Login(user entities.User) (*authpb.LoginResponse, error) {
	return s.client.Login(context.Background(), &authpb.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
}
