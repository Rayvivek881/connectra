package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vivek-ray/middleware"
	"vivek-ray/modules/common"
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
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
	}))

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.Use(middleware.RateLimiter())
	router.Use(middleware.APIKeyAuth())

	router.SetTrustedProxies(nil)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	common.Routes(router.Group("/common"))
	companies.Routes(router.Group("/companies"))
	contacts.Routes(router.Group("/contacts"))

	srv := &http.Server{Addr: ":8000", Handler: router}

	go func() {
		log.Info().Msg("Starting server on :8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Error starting server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
