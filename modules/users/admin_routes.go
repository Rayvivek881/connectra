package users

import (
	"vivek-ray/middleware"
	"vivek-ray/modules/users/controller"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.RouterGroup) {
	adminController := controller.NewAdminController()

	// All admin routes require authentication
	r.Use(middleware.JWTAuth())

	// SuperAdmin only endpoints
	r.GET("", middleware.RequireSuperAdmin(), adminController.ListAllUsers)
	r.PUT("/:user_id/role", middleware.RequireSuperAdmin(), adminController.UpdateUserRole)
	r.PUT("/:user_id/credits", middleware.RequireSuperAdmin(), adminController.UpdateUserCredits)
	r.DELETE("/:user_id", middleware.RequireSuperAdmin(), adminController.DeleteUser)
	r.GET("/history", middleware.RequireSuperAdmin(), adminController.GetUserHistory)

	// Admin+ endpoints
	r.GET("/stats", middleware.RequireAdmin(), adminController.GetUserStats)
}

