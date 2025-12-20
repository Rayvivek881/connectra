package middleware

import (
	"net/http"
	userrepo "vivek-ray/modules/users/repository"

	"github.com/gin-gonic/gin"
)

// CreditCheck checks if user has sufficient credits (skips for Admin/SuperAdmin)
func CreditCheck(requiredCredits int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Authentication required"})
			c.Abort()
			return
		}

		roleString, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Invalid role"})
			c.Abort()
			return
		}

		// Admin and SuperAdmin have unlimited credits
		if roleString == "Admin" || roleString == "SuperAdmin" {
			c.Next()
			return
		}

		// Check credits
		credits, exists := c.Get("credits")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Unable to check credits"})
			c.Abort()
			return
		}

		creditsInt, ok := credits.(int)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Invalid credits value"})
			c.Abort()
			return
		}

		if creditsInt < requiredCredits {
			c.JSON(http.StatusPaymentRequired, gin.H{"detail": "Insufficient credits"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// DeductCredits deducts credits after operation (should be called after successful operation)
func DeductCredits(amount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware should be used after the operation succeeds
		// For now, we'll just pass through - actual deduction happens in service layer
		userUUID, exists := c.Get("user_uuid")
		if !exists {
			c.Next()
			return
		}

		uuidStr, ok := userUUID.(string)
		if !ok {
			c.Next()
			return
		}

		userRole, exists := c.Get("role")
		if exists {
			roleString, ok := userRole.(string)
			if ok && (roleString == "Admin" || roleString == "SuperAdmin") {
				// Skip deduction for Admin/SuperAdmin
				c.Next()
				return
			}
		}

		// Deduct credits in background (non-blocking)
		go func() {
			profileRepo := userrepo.NewUserProfileRepository()
			profile, err := profileRepo.GetByUserID(c.Request.Context(), uuidStr)
			if err == nil && profile != nil {
				if profile.Credits >= amount {
					profile.Credits -= amount
					profileRepo.UpdateProfile(c.Request.Context(), profile)
				}
			}
		}()

		c.Next()
	}
}

