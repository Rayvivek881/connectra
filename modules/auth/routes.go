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

	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	r.Use(middleware.JWTAuth())
	r.POST("/logout", authController.Logout)
}

func EnsureTables() error {
	repo := repository.NewAuthRepository()
	return repo.EnsureTables(context.Background())
}
