package models

import "cm/internal/entities"

type AddConsultantRequest struct {
	User entities.User
}

type AddConsultantResponse struct {
	Error string
}
