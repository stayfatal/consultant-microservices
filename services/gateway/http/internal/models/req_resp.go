package models

import "cm/services/entities"

type RegistrationRequest struct {
	User entities.User
}

type RegistrationResponse struct {
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
