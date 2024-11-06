package router

import (
	"cm/services/gateway/websocket/internal/handlers"
	"cm/services/gateway/websocket/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func NewRouter(hm *handlers.HandlersManager) *gin.Engine {
	r := gin.Default()
	withAuth := r.Group("/").Use(middlewares.Authenticator())
	withAuth.GET("/ws/question", hm.QuestionHandler)
	return r
}
