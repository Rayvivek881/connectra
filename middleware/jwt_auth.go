package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"vivek-ray/conf"
	authrepo "vivek-ray/modules/auth/repository"
	userrepo "vivek-ray/modules/users/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth validates JWT tokens and extracts user information
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Authentication credentials were not provided."})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
			c.Abort()
			return
		}

		// Check blacklist
		repo := authrepo.NewAuthRepository()
		isBlacklisted, err := repo.IsBlacklisted(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Error checking token blacklist"})
			c.Abort()
			return
		}
		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(conf.AppConfig.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user UUID
			userUUID, ok := claims["sub"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
				c.Abort()
				return
			}

			// Get user profile for role and credits
			userRepo := userrepo.NewUserRepository()
			profileRepo := userrepo.NewUserProfileRepository()
			
			user, err := userRepo.GetByUUID(c.Request.Context(), userUUID)
			if err != nil || user == nil || !user.IsActive {
				c.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
				c.Abort()
				return
			}

			profile, _ := profileRepo.GetByUserID(c.Request.Context(), userUUID)
			
			// Set context values
			c.Set("user_uuid", userUUID)
			c.Set("user_id", user.ID)
			c.Set("email", user.Email)
			
			if profile != nil && profile.Role != nil {
				c.Set("role", *profile.Role)
				c.Set("credits", profile.Credits)
			} else {
				c.Set("role", "FreeUser")
				c.Set("credits", 0)
			}
			
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Given token not valid for any token type"})
			c.Abort()
		}
	}
}

// RoleAuth checks if user has one of the required roles
func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"detail": "You do not have permission to perform this action. Role required."})
			c.Abort()
			return
		}

		roleString, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"detail": "You do not have permission to perform this action. Role required."})
			c.Abort()
			return
		}

		for _, role := range roles {
			if role == roleString {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"detail": "You do not have permission to perform this action. " + strings.Join(roles, " or ") + " role required."})
		c.Abort()
	}
}

// RequireSuperAdmin requires SuperAdmin role
func RequireSuperAdmin() gin.HandlerFunc {
	return RoleAuth("SuperAdmin")
}

// RequireAdmin requires Admin or SuperAdmin role
func RequireAdmin() gin.HandlerFunc {
	return RoleAuth("Admin", "SuperAdmin")
}

// RequireProUser requires ProUser, Admin, or SuperAdmin role
func RequireProUser() gin.HandlerFunc {
	return RoleAuth("ProUser", "Admin", "SuperAdmin")
}
