package middleware

import (
	"net/http"
	"vivek-ray/conf"
	"vivek-ray/constants"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" || apiKey != conf.AppConfig.APIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   constants.UnauthorizedError.Error(),
				"success": false,
			})
			return
		}
		c.Next()
	}
}
