package usage

import (
	"vivek-ray/middleware"
	"vivek-ray/modules/usage/controller"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	usageController := controller.NewUsageController()

	// All routes require JWT authentication
	r.Use(middleware.JWTAuth())

	r.GET("/current/", usageController.GetCurrentUsage)
	r.POST("/track/", usageController.TrackUsage)
}

