package cmd

import (
	"vivek-ray/modules/companies"
	"vivek-ray/modules/contacts"

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
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.SetTrustedProxies(nil)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	companies.Routes(router.Group("/companies"))
	contacts.Routes(router.Group("/contacts"))

	log.Info().Msg("Starting server on :8000")
	if err := router.Run(":8000"); err != nil {
		log.Error().Err(err).Msgf("Error starting server: %v", err.Error())
		return
	}
}
