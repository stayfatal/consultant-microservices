package service

import (
	"cm/gen/chatpb"
	"cm/internal/entities"
	"cm/services/gateway/websocket/config"
	"cm/services/gateway/websocket/internal/interfaces"
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type connectedUser struct {
	conn net.Conn
	user entities.User
}

type service struct {
	clientConn *grpc.ClientConn
	client     chatpb.ChatClient
	mu         sync.RWMutex
	users      map[int]connectedUser
}

func New(cfg *config.ServiceConfig) (interfaces.Service, error) {
	clientConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.CHAT_HOST, cfg.CHAT_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &service{users: make(map[int]connectedUser), client: chatpb.NewChatClient(clientConn), clientConn: clientConn}, nil
}

func (s *service) addUserWithConn(conn net.Conn, user entities.User) {
	s.mu.Lock()
	s.users[user.Id] = connectedUser{conn, user}
	s.mu.Unlock()
}

func (s *service) StartAskingQuestions(conn net.Conn, user entities.User) {
	s.addUserWithConn(conn, user)
	// invoke chat service
}

func (s *service) StartAnsweringQuestions(conn net.Conn, user entities.User) error {
	s.addUserWithConn(conn, user)

	resp, err := s.client.AddConsultant(context.Background(), &chatpb.AddConsultantRequest{
		Id:    int32(user.Id),
		Email: user.Email,
	})
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	return nil
}

func (s *service) GratefulStop() {
	s.clientConn.Close()
}
