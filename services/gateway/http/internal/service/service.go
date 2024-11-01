package service

import (
	"cm/gen/authpb"
	"cm/internal/entities"
	"cm/services/gateway/http/config"
	"cm/services/gateway/http/internal/interfaces"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type service struct {
	client authpb.AuthenticationClient
}

func New(cfg *config.ServiceConfig) (interfaces.Service, error) {
	client, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.SSO_HOST, cfg.SSO_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
