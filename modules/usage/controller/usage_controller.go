package controller

import (
	"net/http"
	"strings"
	"vivek-ray/modules/usage/helper"
	"vivek-ray/modules/usage/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UsageController struct {
	usageService *service.UsageService
}

func NewUsageController() *UsageController {
	return &UsageController{
		usageService: service.NewUsageService(),
	}
}

// GetCurrentUsage returns current feature usage for the authenticated user
// GET /api/v2/usage/current/
func (c *UsageController) GetCurrentUsage(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Not authenticated"})
		return
	}

	uuidStr, ok := userUUID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid user context"})
		return
	}

	usageData, err := c.usageService.GetCurrentUsage(ctx.Request.Context(), uuidStr)
	if err != nil {
		log.Error().Err(err).Msg("Error getting current usage")
		if strings.Contains(err.Error(), "User profile not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve feature usage"})
		return
	}

	ctx.JSON(http.StatusOK, usageData)
}

// TrackUsage tracks feature usage for the authenticated user
// POST /api/v2/usage/track/
func (c *UsageController) TrackUsage(ctx *gin.Context) {
	userUUID, exists := ctx.Get("user_uuid")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Not authenticated"})
		return
	}

	uuidStr, ok := userUUID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid user context"})
		return
	}

	var request helper.TrackUsageRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	// Default amount to 1 if not provided
	amount := request.Amount
	if amount == 0 {
		amount = 1
	}

	// Validate amount
	if amount < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": []map[string]interface{}{
				{
					"type":   "greater_than_equal",
					"loc":    []string{"body", "amount"},
					"msg":    "Input should be greater than or equal to 1",
					"input":  amount,
				},
			},
		})
		return
	}

	result, err := c.usageService.TrackUsage(ctx.Request.Context(), uuidStr, request.Feature, amount)
	if err != nil {
		log.Error().Err(err).Msg("Error tracking usage")
		if strings.Contains(err.Error(), "Invalid feature") {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "User profile not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to track feature usage"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

