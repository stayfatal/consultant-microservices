package interfaces

import (
	"cm/libs/entities"
	"net"
)

type Service interface {
	AddUser(conn net.Conn, user entities.User) error
	AddConsultant(conn net.Conn, user entities.User) error
	StartChat(chat entities.Chat) error
	GratefulStop()
}
