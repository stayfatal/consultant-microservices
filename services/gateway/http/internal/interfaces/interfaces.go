package interfaces

import (
	"cm/services/entities"
	"cm/services/gen/authpb"
)

type Service interface {
	Register(user entities.User) (*authpb.RegisterResponse, error)
	Login(user entities.User) (*authpb.LoginResponse, error)
}
