package cmd

import (
	"vivek-ray/middleware"
	"vivek-ray/modules/auth"
	"vivek-ray/modules/companies"
	"vivek-ray/modules/contacts"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Start the API server",
	Long:  `Start the API server for the Connectra API`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
	}))

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.Use(middleware.RateLimiter())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth.Routes(router.Group("/auth"))

	router.SetTrustedProxies(nil)

	router.Use(middleware.JWTAuth())

	companies.Routes(router.Group("/companies"))
	contacts.Routes(router.Group("/contacts"))

	log.Info().Msg("Starting server on :8000")
	if err := router.Run(":8000"); err != nil {
		log.Error().Err(err).Msgf("Error starting server: %v", err.Error())
		return
	}
}
