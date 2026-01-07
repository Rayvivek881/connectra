package middleware

import (
	"vivek-ray/conf"
	commonHelper "vivek-ray/modules/common/helper"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" || apiKey != conf.AppConfig.APIKey {
			commonHelper.SendUnauthorizedError(c, "Invalid or missing API key")
			c.Abort()
			return
		}
		c.Next()
	}
}
