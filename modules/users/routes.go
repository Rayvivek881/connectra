package users

import (
	"vivek-ray/middleware"
	"vivek-ray/modules/users/controller"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	userController := controller.NewUserController()

	// All user routes require authentication
	r.Use(middleware.JWTAuth())

	// Profile endpoints
	r.GET("/profile", userController.GetProfile)
	r.PUT("/profile", userController.UpdateProfile)
	r.POST("/profile/avatar", userController.UploadAvatar)

	// Promotion endpoints
	r.POST("/promote-to-admin", userController.PromoteToAdmin)
	r.POST("/promote-to-super-admin", middleware.RequireSuperAdmin(), userController.PromoteToSuperAdmin)
}

