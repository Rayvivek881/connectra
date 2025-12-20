package controller

import (
	"net/http"
	"strconv"
	"vivek-ray/modules/users/helper"
	"vivek-ray/modules/users/repository"
	"vivek-ray/modules/users/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AdminController struct {
	userService    *service.UserService
	historyService *service.HistoryService
	userRepo       *repository.UserRepository
	profileRepo    *repository.UserProfileRepository
}

func NewAdminController() *AdminController {
	return &AdminController{
		userService:    service.NewUserService(),
		historyService: service.NewHistoryService(),
		userRepo:       repository.NewUserRepository(),
		profileRepo:    repository.NewUserProfileRepository(),
	}
}

// ListAllUsers returns a paginated list of all users (SuperAdmin only)
func (c *AdminController) ListAllUsers(ctx *gin.Context) {
	limit := 100
	offset := 0

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	if offsetStr := ctx.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, total, err := c.userRepo.ListAllUsers(ctx.Request.Context(), limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Error listing users")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to list users"})
		return
	}

	userList := make([]helper.UserListItem, 0, len(users))
	for _, user := range users {
		profile, _ := c.profileRepo.GetByUserID(ctx.Request.Context(), user.UUID)
		
		item := helper.UserListItem{
			UUID:             user.UUID,
			Email:            user.Email,
			Name:             user.Name,
			IsActive:         user.IsActive,
			CreatedAt:        user.CreatedAt,
			LastSignInAt:     user.LastSignInAt,
		}

		if profile != nil {
			item.Role = profile.Role
			item.Credits = profile.Credits
			item.SubscriptionPlan = profile.SubscriptionPlan
			item.SubscriptionPeriod = profile.SubscriptionPeriod
			item.SubscriptionStatus = profile.SubscriptionStatus
		}

		userList = append(userList, item)
	}

	ctx.JSON(http.StatusOK, helper.UserListResponse{
		Users: userList,
		Total: total,
	})
}

// UpdateUserRole updates a user's role (SuperAdmin only)
func (c *AdminController) UpdateUserRole(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "user_id is required"})
		return
	}

	var req helper.UpdateUserRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	// Validate role
	validRoles := map[string]bool{
		"SuperAdmin": true,
		"Admin":      true,
		"ProUser":    true,
		"FreeUser":   true,
	}
	if !validRoles[req.Role] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "Invalid role: " + req.Role + ". Valid roles: SuperAdmin, Admin, FreeUser, ProUser",
		})
		return
	}

	profile, err := c.profileRepo.GetByUserID(ctx.Request.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting profile")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update user role"})
		return
	}
	if profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
		return
	}

	profile.Role = &req.Role
	if err := c.profileRepo.UpdateProfile(ctx.Request.Context(), profile); err != nil {
		log.Error().Err(err).Msg("Error updating profile")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update user role"})
		return
	}

	// Get user for response
	user, _ := c.userRepo.GetByUUID(ctx.Request.Context(), userID)
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
		return
	}

	profileResponse, _ := c.userService.GetUserProfile(ctx.Request.Context(), userID)
	ctx.JSON(http.StatusOK, profileResponse)
}

// UpdateUserCredits updates a user's credits (SuperAdmin only)
func (c *AdminController) UpdateUserCredits(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "user_id is required"})
		return
	}

	var req helper.UpdateUserCreditsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if req.Credits < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "Credits must be non-negative"})
		return
	}

	profile, err := c.profileRepo.GetByUserID(ctx.Request.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Error getting profile")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update user credits"})
		return
	}
	if profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
		return
	}

	profile.Credits = req.Credits
	if err := c.profileRepo.UpdateProfile(ctx.Request.Context(), profile); err != nil {
		log.Error().Err(err).Msg("Error updating profile")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update user credits"})
		return
	}

	profileResponse, _ := c.userService.GetUserProfile(ctx.Request.Context(), userID)
	ctx.JSON(http.StatusOK, profileResponse)
}

// DeleteUser deletes a user (SuperAdmin only)
func (c *AdminController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "user_id is required"})
		return
	}

	currentUserUUID, exists := ctx.Get("user_uuid")
	if exists {
		if uuidStr, ok := currentUserUUID.(string); ok && uuidStr == userID {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": "Cannot delete your own account"})
			return
		}
	}

	if err := c.userRepo.DeleteUser(ctx.Request.Context(), userID); err != nil {
		log.Error().Err(err).Msg("Error deleting user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to delete user"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetUserStats returns user statistics (Admin+)
func (c *AdminController) GetUserStats(ctx *gin.Context) {
	stats, err := c.userService.GetUserStats(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting user stats")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to get user statistics"})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetUserHistory returns user history records (SuperAdmin only)
func (c *AdminController) GetUserHistory(ctx *gin.Context) {
	var userID *string
	if userIDStr := ctx.Query("user_uuid"); userIDStr != "" {
		userID = &userIDStr
	}

	var eventType *string
	if eventTypeStr := ctx.Query("event_type"); eventTypeStr != "" {
		eventType = &eventTypeStr
	}

	limit := 100
	offset := 0

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	if offsetStr := ctx.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	history, err := c.historyService.GetUserHistory(ctx.Request.Context(), userID, eventType, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Error getting user history")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to get user history"})
		return
	}

	ctx.JSON(http.StatusOK, history)
}

