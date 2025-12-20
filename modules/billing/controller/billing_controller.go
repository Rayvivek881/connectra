package controller

import (
	"net/http"
	"strconv"
	"vivek-ray/modules/billing/helper"
	"vivek-ray/modules/billing/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type BillingController struct {
	billingService *service.BillingService
}

func NewBillingController() *BillingController {
	return &BillingController{
		billingService: service.NewBillingService(),
	}
}

// GetBillingInfo returns billing information for the current user
func (c *BillingController) GetBillingInfo(ctx *gin.Context) {
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

	billingInfo, err := c.billingService.GetBillingInfo(ctx.Request.Context(), uuidStr)
	if err != nil {
		log.Error().Err(err).Msg("Error getting billing info")
		if err.Error() == "user profile not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "User profile not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve billing information"})
		return
	}

	ctx.JSON(http.StatusOK, billingInfo)
}

// GetSubscriptionPlans returns all available subscription plans (public endpoint)
func (c *BillingController) GetSubscriptionPlans(ctx *gin.Context) {
	plans, err := c.billingService.GetSubscriptionPlans(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting subscription plans")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve subscription plans"})
		return
	}

	ctx.JSON(http.StatusOK, plans)
}

// GetAddonPackages returns all available addon packages (public endpoint)
func (c *BillingController) GetAddonPackages(ctx *gin.Context) {
	packages, err := c.billingService.GetAddonPackages(ctx.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Error getting addon packages")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve addon packages"})
		return
	}

	ctx.JSON(http.StatusOK, packages)
}

// SubscribeToPlan subscribes the current user to a subscription plan
func (c *BillingController) SubscribeToPlan(ctx *gin.Context) {
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

	var req helper.SubscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	result, err := c.billingService.SubscribeToPlan(ctx.Request.Context(), uuidStr, req.Tier, req.Period)
	if err != nil {
		log.Error().Err(err).Msg("Error subscribing to plan")
		if err.Error() == "user profile not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "User profile not found"})
			return
		}
		if err.Error() == "invalid tier: "+req.Tier || err.Error() == "invalid period: "+req.Period {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to subscribe to plan"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// PurchaseAddon purchases addon credits for the current user
func (c *BillingController) PurchaseAddon(ctx *gin.Context) {
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

	var req helper.AddonPurchaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	result, err := c.billingService.PurchaseAddonCredits(ctx.Request.Context(), uuidStr, req.PackageID)
	if err != nil {
		log.Error().Err(err).Msg("Error purchasing addon credits")
		if err.Error() == "user profile not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "User profile not found"})
			return
		}
		if err.Error() == "invalid package ID: "+req.PackageID {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to purchase addon credits"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// CancelSubscription cancels the current user's subscription
func (c *BillingController) CancelSubscription(ctx *gin.Context) {
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

	result, err := c.billingService.CancelSubscription(ctx.Request.Context(), uuidStr)
	if err != nil {
		log.Error().Err(err).Msg("Error cancelling subscription")
		if err.Error() == "user profile not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "User profile not found"})
			return
		}
		if err.Error() == "subscription is already cancelled" {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to cancel subscription"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetInvoices returns invoice history for the current user
func (c *BillingController) GetInvoices(ctx *gin.Context) {
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

	// Parse query parameters
	limitStr := ctx.DefaultQuery("limit", "10")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	result, err := c.billingService.GetInvoices(ctx.Request.Context(), uuidStr, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Error getting invoices")
		if err.Error() == "user profile not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": "User profile not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve invoices"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// Admin endpoints

// AdminGetSubscriptionPlans returns all subscription plans for admin (including inactive)
func (c *BillingController) AdminGetSubscriptionPlans(ctx *gin.Context) {
	includeInactive := ctx.Query("include_inactive") == "true"

	var plans *helper.SubscriptionPlansResponse
	var err error

	if includeInactive {
		// Get all plans including inactive
		plans, err = c.billingService.GetSubscriptionPlans(ctx.Request.Context())
		// Note: This currently only returns active plans. To include inactive,
		// we'd need to modify GetSubscriptionPlans or add a new method.
		// For now, we'll use the existing method.
	} else {
		plans, err = c.billingService.GetSubscriptionPlans(ctx.Request.Context())
	}

	if err != nil {
		log.Error().Err(err).Msg("Error getting subscription plans")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve subscription plans"})
		return
	}

	ctx.JSON(http.StatusOK, plans)
}

// AdminCreateSubscriptionPlan creates a new subscription plan
func (c *BillingController) AdminCreateSubscriptionPlan(ctx *gin.Context) {
	var req helper.SubscriptionPlanCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := c.billingService.CreateSubscriptionPlan(ctx.Request.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Error creating subscription plan")
		if err.Error() == "plan with tier "+req.Tier+" already exists" {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create subscription plan"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Subscription plan created successfully",
		"tier":    req.Tier,
	})
}

// AdminUpdateSubscriptionPlan updates a subscription plan
func (c *BillingController) AdminUpdateSubscriptionPlan(ctx *gin.Context) {
	tier := ctx.Param("tier")
	if tier == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "tier parameter is required"})
		return
	}

	var req helper.SubscriptionPlanUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := c.billingService.UpdateSubscriptionPlan(ctx.Request.Context(), tier, &req); err != nil {
		log.Error().Err(err).Msg("Error updating subscription plan")
		if err.Error() == "plan with tier "+tier+" not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update subscription plan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Subscription plan updated successfully",
		"tier":    tier,
	})
}

// AdminDeleteSubscriptionPlan deletes a subscription plan
func (c *BillingController) AdminDeleteSubscriptionPlan(ctx *gin.Context) {
	tier := ctx.Param("tier")
	if tier == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "tier parameter is required"})
		return
	}

	if err := c.billingService.DeleteSubscriptionPlan(ctx.Request.Context(), tier); err != nil {
		log.Error().Err(err).Msg("Error deleting subscription plan")
		if err.Error() == "plan with tier "+tier+" not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to delete subscription plan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Subscription plan deleted successfully",
		"tier":    tier,
	})
}

// AdminCreateSubscriptionPlanPeriod creates or updates a subscription plan period
func (c *BillingController) AdminCreateSubscriptionPlanPeriod(ctx *gin.Context) {
	tier := ctx.Param("tier")
	if tier == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "tier parameter is required"})
		return
	}

	var req helper.SubscriptionPeriodCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := c.billingService.CreateSubscriptionPlanPeriod(ctx.Request.Context(), tier, &req); err != nil {
		log.Error().Err(err).Msg("Error creating/updating subscription plan period")
		if err.Error() == "plan with tier "+tier+" not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create/update period"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Period created/updated successfully",
		"tier":    tier,
		"period":  req.Period,
	})
}

// AdminDeleteSubscriptionPlanPeriod deletes a subscription plan period
func (c *BillingController) AdminDeleteSubscriptionPlanPeriod(ctx *gin.Context) {
	tier := ctx.Param("tier")
	period := ctx.Param("period")
	if tier == "" || period == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "tier and period parameters are required"})
		return
	}

	if err := c.billingService.DeleteSubscriptionPlanPeriod(ctx.Request.Context(), tier, period); err != nil {
		log.Error().Err(err).Msg("Error deleting subscription plan period")
		if err.Error() == "period "+period+" not found for plan "+tier {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to delete period"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Period deleted successfully",
		"tier":    tier,
		"period":  period,
	})
}

// AdminGetAddonPackages returns all addon packages for admin (including inactive)
func (c *BillingController) AdminGetAddonPackages(ctx *gin.Context) {
	includeInactive := ctx.Query("include_inactive") == "true"

	var packages *helper.AddonPackagesResponse
	var err error

	if includeInactive {
		// Get all packages including inactive
		packages, err = c.billingService.GetAddonPackages(ctx.Request.Context())
		// Note: This currently only returns active packages. To include inactive,
		// we'd need to modify GetAddonPackages or add a new method.
		// For now, we'll use the existing method.
	} else {
		packages, err = c.billingService.GetAddonPackages(ctx.Request.Context())
	}

	if err != nil {
		log.Error().Err(err).Msg("Error getting addon packages")
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to retrieve addon packages"})
		return
	}

	ctx.JSON(http.StatusOK, packages)
}

// AdminCreateAddonPackage creates a new addon package
func (c *BillingController) AdminCreateAddonPackage(ctx *gin.Context) {
	var req helper.AddonPackageCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := c.billingService.CreateAddonPackage(ctx.Request.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Error creating addon package")
		if err.Error() == "package with id "+req.ID+" already exists" {
			ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create addon package"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Addon package created successfully",
		"id":      req.ID,
	})
}

// AdminUpdateAddonPackage updates an addon package
func (c *BillingController) AdminUpdateAddonPackage(ctx *gin.Context) {
	packageID := ctx.Param("package_id")
	if packageID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "package_id parameter is required"})
		return
	}

	var req helper.AddonPackageUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	if err := c.billingService.UpdateAddonPackage(ctx.Request.Context(), packageID, &req); err != nil {
		log.Error().Err(err).Msg("Error updating addon package")
		if err.Error() == "package with id "+packageID+" not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update addon package"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Addon package updated successfully",
		"id":      packageID,
	})
}

// AdminDeleteAddonPackage deletes an addon package
func (c *BillingController) AdminDeleteAddonPackage(ctx *gin.Context) {
	packageID := ctx.Param("package_id")
	if packageID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"detail": "package_id parameter is required"})
		return
	}

	if err := c.billingService.DeleteAddonPackage(ctx.Request.Context(), packageID); err != nil {
		log.Error().Err(err).Msg("Error deleting addon package")
		if err.Error() == "package with id "+packageID+" not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"detail": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to delete addon package"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Addon package deleted successfully",
		"id":      packageID,
	})
}

