package service

import (
	"bytes"
	"cm/libs/entities"
	"cm/libs/log"
	"cm/services/matchmaking/internal/interfaces"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"sync"
)

type service struct {
	rabbit      *rabbitConfig
	repo        interfaces.Repository
	consultants []entities.User
	mu          sync.RWMutex
	logger      *log.Logger
}

func New(repo interfaces.Repository, logger *log.Logger) (interfaces.Service, error) {
	rabbit, err := dialRabbit()
	if err != nil {
		return nil, err
	}

	svc := &service{repo: repo, logger: logger, rabbit: rabbit}
	svc.startService()
	return svc, nil
}

func (s *service) AddConsultant(user entities.User) {
	s.mu.Lock()
	s.consultants = append(s.consultants, user)
	s.mu.Unlock()
}

func (s *service) GratefulStop() {
	s.rabbit.ch.Close()
	s.rabbit.conn.Close()
}

func (s *service) startService() {
	msgs, err := s.rabbit.ch.Consume(
		s.rabbit.q.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			if len(s.consultants) > 0 {
				msg := <-msgs
				user := entities.User{}
				err = json.Unmarshal(msg.Body, &user)
				if err != nil {
					s.logger.Log(err)
					continue
				}

				err := s.startChat(entities.Chat{
					ConsultantId: s.consultants[0].Id,
					UserId:       user.Id,
				})
				if err != nil {
					s.logger.Log(err)
					continue
				}

				s.popConsultantFromQueue()
			}
		}
	}()
}

func (s *service) startChat(chat entities.Chat) error {
	bChat, err := json.Marshal(chat)
	if err != nil {
		s.logger.Log(err)
		return err
	}

	buf := bytes.NewBuffer(bChat)
	resp, err := http.Post("http://chat:2955/chat", "application/json", buf)
	if err != nil {
		s.logger.Log(err)
		return err
	}

	_, err = io.Copy(os.Stderr, resp.Body)
	if err != nil {
		s.logger.Log(err)
	}
	return nil
}

func (s *service) popConsultantFromQueue() {
	s.mu.Lock()
	s.consultants = s.consultants[1:]
	s.mu.Unlock()
}
