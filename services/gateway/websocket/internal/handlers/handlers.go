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

func (hm *HandlersManager) WebsocketChatHandler(c *gin.Context) {
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
		err := hm.svc.AddConsultant(conn, user)
		if err != nil {
			c.String(http.StatusInternalServerError, TypeAssertingError.Error())
			hm.logger.Log(err)
			return
		}
	} else {
		err := hm.svc.AddUser(conn, user)
		if err != nil {
			c.String(http.StatusInternalServerError, TypeAssertingError.Error())
			hm.logger.Log(err)
			return
		}
	}

	c.Status(http.StatusOK)
}

func (hm *HandlersManager) ChatHandler(c *gin.Context) {
	chat := entities.Chat{}
	err := c.BindJSON(&chat)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		hm.logger.Log(err)
		return
	}

	if err := hm.svc.StartChat(chat); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		hm.logger.Log(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
