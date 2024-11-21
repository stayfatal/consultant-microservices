package interfaces

import "cm/internal/entities"

type Service interface {
	AddConsultant(user entities.User)
}

type Repository interface {
	CreateChat(chat entities.Chat) (id int, err error)
	SaveMessage(message entities.Message) error
}
