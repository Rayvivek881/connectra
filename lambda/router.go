package lambda

import (
	"vivek-ray/middleware"
	"vivek-ray/modules/common"
	"vivek-ray/modules/companies"
	"vivek-ray/modules/contacts"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// InitRouter initializes and returns a fully configured Gin router
// This function extracts the router setup from cmd/server.go for reuse in Lambda
func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
	}))

	// Gzip compression
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Rate limiting middleware
	router.Use(middleware.RateLimiter())

	// API key authentication middleware
	router.Use(middleware.APIKeyAuth())

	router.SetTrustedProxies(nil)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register all routes
	common.Routes(router.Group("/common"))
	companies.Routes(router.Group("/companies"))
	contacts.Routes(router.Group("/contacts"))

	return router
}
