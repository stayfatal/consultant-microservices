package interfaces

import (
	"cm/internal/entities"
	"net"
)

type Service interface {
	StartAskingQuestions(conn net.Conn, user entities.User)
	StartAnsweringQuestions(conn net.Conn, user entities.User) error
	GratefulStop()
}
