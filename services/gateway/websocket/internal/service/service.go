package service

import (
	"cm/gen/chatpb"
	"cm/internal/entities"
	"cm/internal/log"
	"cm/services/gateway/websocket/config"
	"cm/services/gateway/websocket/internal/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/gobwas/ws/wsutil"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type connectedUser struct {
	conn net.Conn
	user entities.User
}

type service struct {
	rabbit     *rabbitConfig
	clientConn *grpc.ClientConn
	client     chatpb.ChatClient
	mu         sync.RWMutex
	users      map[int]connectedUser
	logger     *log.Logger
}

func New(cfg *config.ServiceConfig, logger *log.Logger) (interfaces.Service, error) {
	clientConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.CHAT_HOST, cfg.CHAT_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	rabbit, err := dialRabbit()
	if err != nil {
		return nil, err
	}

	return &service{users: make(map[int]connectedUser), client: chatpb.NewChatClient(clientConn), clientConn: clientConn, rabbit: rabbit, logger: logger}, nil
}

func (s *service) AddUser(conn net.Conn, user entities.User) error {
	s.addUserWithConn(conn, user)

	bUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = s.rabbit.ch.Publish("", s.rabbit.q.Name, false, false, amqp.Publishing{ContentType: "application/json", Body: bUser})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) AddConsultant(conn net.Conn, user entities.User) error {
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

func (s *service) StartChat(chat entities.Chat) error {
	user, ok := s.getUserWithConn(chat.UserId)
	if !ok {
		return errors.Errorf("user %v is not connected", user.user)
	}
	consultant, ok := s.getUserWithConn(chat.ConsultantId)
	if !ok {
		return errors.Errorf("consultant %v is not connected", consultant.user)
	}

	userToConsultant := make(chan []byte)
	consultantToUser := make(chan []byte)
	s.startRW(user, consultantToUser, userToConsultant)
	s.startRW(consultant, userToConsultant, consultantToUser)
	return nil
}

func (s *service) startRW(user connectedUser, in chan []byte, out chan []byte) {
	go func() {
		for {
			msg, err := wsutil.ReadClientText(user.conn)
			if err != nil {
				//delete chat invoke
				s.logger.Log(err)
				return
			}

			out <- msg
		}
	}()

	go func() {
		for {
			msg := <-in

			err := wsutil.WriteServerText(user.conn, msg)
			if err != nil {
				//delete chat invoke
				s.logger.Log(err)
				return
			}
		}
	}()
}

func (s *service) GratefulStop() {
	for _, user := range s.users {
		user.conn.Close()
	}

	s.rabbit.conn.Close()
	s.rabbit.ch.Close()

	s.clientConn.Close()
}

func (s *service) addUserWithConn(conn net.Conn, user entities.User) {
	s.mu.Lock()
	s.users[user.Id] = connectedUser{conn, user}
	s.mu.Unlock()
}

func (s *service) getUserWithConn(id int) (connectedUser, bool) {
	s.mu.RLock()
	userWithConn, ok := s.users[id]
	s.mu.RUnlock()
	return userWithConn, ok
}
