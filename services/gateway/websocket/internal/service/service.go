package service

import (
	"cm/internal/entities"
	"cm/services/gateway/websocket/internal/interfaces"
	"net"
	"sync"
)

type connectedUser struct {
	conn net.Conn
	user entities.User
}

type service struct {
	mu    sync.RWMutex
	users map[int]connectedUser
}

func New() interfaces.Service {
	return &service{users: make(map[int]connectedUser)}
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

func (s *service) StartAnsweringQuestions(conn net.Conn, user entities.User) {
	s.addUserWithConn(conn, user)
	// invoke chat service
}
