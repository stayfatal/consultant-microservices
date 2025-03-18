package models

import "cm/libs/entities"

type RegisterRequest struct {
	User entities.User
}

type RegisterResponse struct {
	Token string
	Error string
}

type LoginRequest struct {
	User entities.User
}

type LoginResponse struct {
	Token string
	Error string
}
