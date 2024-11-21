package service

import (
	"cm/internal/entities"
	"cm/internal/log"
	"cm/services/chat/internal/interfaces"

	"sync"
)

type service struct {
	repo        interfaces.Repository
	consultants map[int]entities.User
	mu          sync.RWMutex
}

func New(repo interfaces.Repository) interfaces.Service {
	return &service{repo: repo, consultants: make(map[int]entities.User)}
}

func (svc *service) AddConsultant(user entities.User) {
	svc.mu.Lock()
	svc.consultants[user.Id] = user
	svc.mu.Unlock()
	log.New().Println(svc.consultants)
}
