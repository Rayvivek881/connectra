package auth

import (
	"context"
	"vivek-ray/middleware"
	"vivek-ray/modules/auth/controller"
	"vivek-ray/modules/auth/repository"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	authController := controller.NewAuthController()

	// Public endpoints
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)
	r.POST("/refresh", authController.RefreshToken)

	// Protected endpoints
	r.Use(middleware.JWTAuth())
	r.POST("/logout", authController.Logout)
	r.GET("/session", authController.GetSession)
}

func EnsureTables() error {
	repo := repository.NewAuthRepository()
	return repo.EnsureTables(context.Background())
}
