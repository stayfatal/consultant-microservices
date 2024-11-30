package middlewares

import (
	"cm/internal/entities"
	"cm/internal/publicauth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Authenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.AbortWithError(http.StatusBadRequest, publicauth.InvalidTokenError)
			return
		}

		claims, err := publicauth.ValidateToken(token)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, publicauth.InvalidTokenError)
			return
		}

		c.Set("user", entities.User{Id: claims.Id, Email: claims.Email, IsConsultant: claims.IsConsultant})

		c.Next()
	}
}

func IpWhiteList(whiteList map[string]struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := whiteList[c.ClientIP()]; !ok {
			log.Warn().Msgf("Blocked IP: %s", c.ClientIP())
			c.AbortWithStatus(http.StatusForbidden)
		}
		c.Next()
	}
}
