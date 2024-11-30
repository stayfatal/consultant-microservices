package service

import (
	"cm/internal/entities"
	"cm/internal/log"
	"cm/services/chat/internal/interfaces"
	"encoding/json"
	"fmt"
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

				s.startChat(entities.Chat{
					ConsultantId: s.consultants[0].Id,
					UserId:       user.Id,
				})

				s.consultants = s.consultants[1:]

				fmt.Println(user)
			}
		}
	}()
}

func (s *service) startChat(chat entities.Chat) {

}

func (svc *service) AddConsultant(user entities.User) {
	svc.mu.Lock()
	svc.consultants = append(svc.consultants, user)
	svc.mu.Unlock()
	log.New().Println(svc.consultants)
}
