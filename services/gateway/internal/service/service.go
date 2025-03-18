package service

import (
	"cm/gen/authpb"
	"cm/libs/config"
	"cm/libs/entities"
	"cm/services/gateway/internal/interfaces"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type service struct {
	clientConn *grpc.ClientConn
	client     authpb.AuthenticationClient
}

func New(cfg *config.ServicesConfig) (interfaces.Service, error) {
	clientConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Auth.Host, cfg.Auth.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &service{client: authpb.NewAuthenticationClient(clientConn), clientConn: clientConn}, nil
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

func (s *service) GratefulStop() {
	s.clientConn.Close()
}
