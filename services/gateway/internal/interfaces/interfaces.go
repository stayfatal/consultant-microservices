package interfaces

import (
	"cm/gen/authpb"
	"cm/libs/entities"
)

type Service interface {
	Register(user entities.User) (*authpb.RegisterResponse, error)
	Login(user entities.User) (*authpb.LoginResponse, error)
	GratefulStop()
}
