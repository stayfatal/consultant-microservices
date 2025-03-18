package models

import "cm/libs/entities"

type AddConsultantRequest struct {
	User entities.User
}

type AddConsultantResponse struct {
	Error string
}
