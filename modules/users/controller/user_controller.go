package controller

import (
	"io"
	"net/http"
	"vivek-ray/modules/users/helper"
	"vivek-ray/modules/users/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserController struct {
	userService   *service.UserService
	avatarService *service.AvatarService
}

func NewUserController() *UserController {
	return &UserController{
		userService:   service.NewUserService(),
		avatarService: service.NewAvatarService(),
	}
}

// GetProfile returns the current user's profile
func (c *UserController) GetProfile(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Authentication required"})
		return
	}

	uuidStr, ok := userUUID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid user context"})
		return
	}

	profile, err := c.userService.GetUserProfile(ctx.Request.Context(), uuidStr)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user profile")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve profile"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// UpdateProfile updates the current user's profile
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Authentication required"})
		return
	}

	uuidStr, ok := userUUID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid user context"})
		return
	}

	var update helper.ProfileUpdateRequest
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	profile, err := c.userService.UpdateUserProfile(ctx.Request.Context(), uuidStr, &update)
	if err != nil {
		log.Error().Err(err).Msg("Error updating user profile")
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid data provided"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// UploadAvatar handles avatar file upload
func (c *UserController) UploadAvatar(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Authentication required"})
		return
	}

	uuidStr, ok := userUUID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid user context"})
		return
	}

	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"avatar": []string{"This field is required."}})
		return
	}

	// Read file content
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"avatar": []string{"Error reading file"}})
		return
	}
	defer src.Close()

	fileContent, err := io.ReadAll(src)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"avatar": []string{"Error reading file"}})
		return
	}

	// Upload avatar
	avatarURL, err := c.avatarService.UploadAvatar(
		ctx.Request.Context(),
		uuidStr,
		file.Filename,
		fileContent,
		file.Header.Get("Content-Type"),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error uploading avatar")
		if err.Error() == "invalid file type" || err.Error() == "file too large" || err.Error() == "file does not appear to be a valid image file" {
			ctx.JSON(http.StatusBadRequest, gin.H{"avatar": []string{err.Error()}})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Error saving file: " + err.Error()})
		return
	}

	// Get updated profile
	profile, err := c.userService.GetUserProfile(ctx.Request.Context(), uuidStr)
	if err != nil {
		log.Error().Err(err).Msg("Error getting updated profile")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Error retrieving profile"})
		return
	}

	ctx.JSON(http.StatusOK, helper.AvatarUploadResponse{
		AvatarURL: avatarURL,
		Profile:   *profile,
		Message:   "Avatar uploaded successfully",
	})
}

// PromoteToAdmin promotes the current user to admin
func (c *UserController) PromoteToAdmin(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Authentication required"})
		return
	}

	uuidStr, ok := userUUID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid user context"})
		return
	}

	profile, err := c.userService.PromoteUserToAdmin(ctx.Request.Context(), uuidStr)
	if err != nil {
		log.Error().Err(err).Msg("Error promoting user to admin")
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
			return
		}
		if err.Error() == "user account is disabled" {
			ctx.JSON(http.StatusBadRequest, gin.H{"non_field_errors": []string{"User account is disabled"}})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to promote user to admin"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

// PromoteToSuperAdmin promotes a user to super admin (SuperAdmin only)
func (c *UserController) PromoteToSuperAdmin(ctx *gin.Context) {
	userRole, exists := ctx.Get("role")
	if !exists {
		ctx.JSON(http.StatusForbidden, gin.H{"detail": "You do not have permission to perform this action. SuperAdmin role required."})
		return
	}

	roleStr, ok := userRole.(string)
	if !ok || roleStr != "SuperAdmin" {
		ctx.JSON(http.StatusForbidden, gin.H{"detail": "You do not have permission to perform this action. SuperAdmin role required."})
		return
	}

	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "user_id query parameter is required"})
		return
	}

	profile, err := c.userService.PromoteUserToSuperAdmin(ctx.Request.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Error promoting user to super admin")
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
			return
		}
		if err.Error() == "user account is disabled" {
			ctx.JSON(http.StatusBadRequest, gin.H{"non_field_errors": []string{"User account is disabled"}})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to promote user to super admin"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

