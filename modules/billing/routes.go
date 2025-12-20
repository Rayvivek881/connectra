package billing

import (
	"vivek-ray/middleware"
	"vivek-ray/modules/billing/controller"

	"github.com/gin-gonic/gin"
)

// Routes sets up billing routes
func Routes(r *gin.RouterGroup) {
	billingController := controller.NewBillingController()

	// Public endpoints (no authentication required)
	r.GET("/plans/", billingController.GetSubscriptionPlans)
	r.GET("/addons/", billingController.GetAddonPackages)

	// Authenticated endpoints (require JWT)
	authenticated := r.Group("")
	authenticated.Use(middleware.JWTAuth())
	{
		authenticated.GET("/", billingController.GetBillingInfo)
		authenticated.POST("/subscribe/", billingController.SubscribeToPlan)
		authenticated.POST("/addon/", billingController.PurchaseAddon)
		authenticated.POST("/cancel/", billingController.CancelSubscription)
		authenticated.GET("/invoices/", billingController.GetInvoices)
	}

	// Admin endpoints (require SuperAdmin role)
	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.RequireSuperAdmin())
	{
		// Subscription plan management
		admin.GET("/plans/", billingController.AdminGetSubscriptionPlans)
		admin.POST("/plans/", billingController.AdminCreateSubscriptionPlan)
		admin.PUT("/plans/:tier/", billingController.AdminUpdateSubscriptionPlan)
		admin.DELETE("/plans/:tier/", billingController.AdminDeleteSubscriptionPlan)

		// Subscription plan period management
		admin.POST("/plans/:tier/periods/", billingController.AdminCreateSubscriptionPlanPeriod)
		admin.DELETE("/plans/:tier/periods/:period/", billingController.AdminDeleteSubscriptionPlanPeriod)

		// Addon package management
		admin.GET("/addons/", billingController.AdminGetAddonPackages)
		admin.POST("/addons/", billingController.AdminCreateAddonPackage)
		admin.PUT("/addons/:package_id/", billingController.AdminUpdateAddonPackage)
		admin.DELETE("/addons/:package_id/", billingController.AdminDeleteAddonPackage)
	}
}

