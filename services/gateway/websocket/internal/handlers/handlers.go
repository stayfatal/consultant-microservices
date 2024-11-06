package handlers

import (
	"cm/internal/entities"
	"cm/internal/log"
	"cm/services/gateway/websocket/internal/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/pkg/errors"
)

var (
	GinContextValueError = errors.New("no value in gin context")
	TypeAssertingError   = errors.New("can't assert types")
)

type HandlersManager struct {
	logger *log.Logger
	svc    interfaces.Service
}

func NewManager(logger *log.Logger, svc interfaces.Service) *HandlersManager {
	return &HandlersManager{logger, svc}
}

func (hm *HandlersManager) QuestionHandler(c *gin.Context) {
	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		hm.logger.Log(err)
		return
	}

	val, ok := c.Get("user")
	if !ok {
		c.String(http.StatusInternalServerError, GinContextValueError.Error())
		hm.logger.Log(GinContextValueError)
		return
	}

	user, ok := val.(entities.User)
	if !ok {
		c.String(http.StatusInternalServerError, TypeAssertingError.Error())
		hm.logger.Log(TypeAssertingError)
		return
	}

	if user.IsConsultant {
		hm.svc.StartAnsweringQuestions(conn, user)
	} else {
		hm.svc.StartAskingQuestions(conn, user)
	}

	c.Status(http.StatusOK)
}
