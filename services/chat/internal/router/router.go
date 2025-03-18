package router

import (
	"cm/services/chat/internal/handlers"
	"cm/services/chat/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func NewRouter(hm *handlers.HandlersManager) *gin.Engine {
	r := gin.Default()

	withAuth := r.Group("/").Use(middlewares.Authenticator())
	withAuth.GET("/ws/chat", hm.WebsocketChatHandler)

	// whiteIp :=
	r.Group("/").Use(middlewares.IpWhiteList(map[string]struct{}{
		"chat": {},
	}))
	r.POST("/chat", hm.ChatHandler)

	return r
}
