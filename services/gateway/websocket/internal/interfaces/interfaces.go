package interfaces

import (
	"cm/internal/entities"
	"net"
)

type Service interface {
	AddUser(conn net.Conn, user entities.User) error
	AddConsultant(conn net.Conn, user entities.User) error
	GratefulStop()
}
