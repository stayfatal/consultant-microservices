package service

import "cm/services/user/internal/interfaces"

type service struct {
	repo *interfaces.Repository
}

func New(repo *interfaces.Repository) interfaces.Service {
	return &service{repo: repo}
}

func (svc *service) Register() {

}

func (svc *service) Login() {

}
