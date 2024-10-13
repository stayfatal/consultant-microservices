package models

type RegisterRequest struct {
	User User
}

type RegisterResponse struct {
	Token string
	Error string
}

type LoginRequest struct {
	User User
}

type LoginResponse struct {
	Token string
	Error string
}
