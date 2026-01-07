//go:build integration
// +build integration

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"vivek-ray/conf"
	"vivek-ray/connections"
	"vivek-ray/middleware"
	"vivek-ray/modules/common"
	"vivek-ray/modules/companies"
	"vivek-ray/modules/contacts"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var testRouter *gin.Engine
var testAPIKey = "test-api-key-for-integration-tests"

func setupTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	
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
	
	// Setup routes
	common.Routes(router.Group("/common"))
	companies.Routes(router.Group("/companies"))
	contacts.Routes(router.Group("/contacts"))
	
	return router
}

func TestMain(m *testing.M) {
	// Initialize configuration
	v := conf.Viper{}
	v.Init()
	
	// Set test API key
	conf.AppConfig.APIKey = testAPIKey
	
	// Initialize connections (use test database if available)
	connections.InitDB()
	connections.InitSearchEngine()
	
	// Setup test router
	testRouter = setupTestServer()
	
	// Run tests
	code := m.Run()
	
	// Cleanup
	connections.CloseDB()
	connections.CloseSearchEngine()
	
	// Exit
	os.Exit(code)
}

func makeRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	
	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", testAPIKey)
	
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	
	return w
}
