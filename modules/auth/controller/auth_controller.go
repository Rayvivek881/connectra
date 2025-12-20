package controller

import (
	"net/http"
	"strings"
	"vivek-ray/modules/auth/service"
	"vivek-ray/modules/users/helper"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		service: service.NewAuthService(),
	}
}

// Register handles user registration
func (c *AuthController) Register(ctx *gin.Context) {
	var req helper.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"detail": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := c.service.Register(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "email already exists" {
			ctx.JSON(http.StatusBadRequest, gin.H{"email": []string{"Email already exists"}})
			return
		}
		if err.Error() == "password must be at least 8 characters" {
			ctx.JSON(http.StatusBadRequest, gin.H{"password": []string{"This password is too short. It must contain at least 8 characters."}})
			return
		}
		if err.Error() == "password must be at most 72 characters" {
			ctx.JSON(http.StatusBadRequest, gin.H{"password": []string{"String should have at most 72 characters"}})
			return
		}
		log.Error().Err(err).Msg("Error registering user")
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "Registration failed"})
		return
	}

	ctx.JSON(http.StatusCreated, helper.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: helper.UserResponse{
			UUID:  user.UUID,
			Email: user.Email,
		},
		Message: "Registration successful! Please check your email to verify your account.",
	})
}

// Login handles user login
func (c *AuthController) Login(ctx *gin.Context) {
	var req helper.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"detail": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := c.service.Login(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid email or password"})
			return
		}
		if err.Error() == "user account is disabled" {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": map[string][]string{"non_field_errors": {"User account is disabled"}}})
			return
		}
		log.Error().Err(err).Msg("Error logging in user")
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "Login failed"})
		return
	}

	ctx.JSON(http.StatusOK, helper.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: helper.UserResponse{
			UUID:  user.UUID,
			Email: user.Email,
		},
	})
}

// Logout handles user logout
func (c *AuthController) Logout(ctx *gin.Context) {
	var req helper.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Logout can succeed even without refresh token
		req = helper.LogoutRequest{}
	}

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Authentication credentials were not provided."})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
		return
	}

	var userID *string
	if userUUID, exists := ctx.Get("user_uuid"); exists {
		if uuidStr, ok := userUUID.(string); ok {
			userID = &uuidStr
		}
	}

	if err := c.service.Logout(ctx.Request.Context(), tokenString, userID); err != nil {
		log.Warn().Err(err).Msg("Error during logout (non-critical)")
		// Logout still succeeds even if blacklist fails
	}

	if req.RefreshToken != nil && *req.RefreshToken != "" {
		if err := c.service.Logout(ctx.Request.Context(), *req.RefreshToken, userID); err != nil {
			log.Warn().Err(err).Msg("Error blacklisting refresh token (non-critical)")
		}
	}

	ctx.JSON(http.StatusOK, helper.LogoutResponse{
		Message: "Logout successful",
	})
}

// GetSession returns current session information
func (c *AuthController) GetSession(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Authentication credentials were not provided."})
		return
	}

	uuidStr, ok := userUUID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
		return
	}

	session, err := c.service.GetSession(ctx.Request.Context(), uuidStr)
	if err != nil {
		log.Error().Err(err).Msg("Error getting session")
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
		return
	}

	ctx.JSON(http.StatusOK, session)
}

// RefreshToken handles token refresh
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req helper.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"refresh_token": []string{"This field is required."}})
		return
	}

	accessToken, refreshToken, err := c.service.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		if err.Error() == "invalid refresh token" {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid refresh token"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "Token is invalid or expired"})
		return
	}

	ctx.JSON(http.StatusOK, helper.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
