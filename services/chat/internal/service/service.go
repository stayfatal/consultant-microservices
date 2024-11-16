package service

import "cm/services/chat/internal/interfaces"

type service struct {
}

func New() interfaces.Service {
	return &service{}
}

func (svc *service) AddConsultant() {

}
