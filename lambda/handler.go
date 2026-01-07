package lambda

import (
	"context"
	"sync"

	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

var (
	ginLambda *ginadapter.GinLambda
	initOnce  sync.Once
)

// Handler is the main Lambda handler function
// It processes API Gateway HTTP API events and returns responses
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize router adapter on first invocation (singleton pattern)
	initOnce.Do(func() {
		log.Info().Msg("Initializing Lambda handler (first invocation)")
		
		// Initialize connections if not already done
		if err := InitConnections(); err != nil {
			log.Error().Err(err).Msg("Failed to initialize connections")
			// Continue anyway - connections might be initialized elsewhere
		}

		// Initialize router
		router := InitRouter()
		ginLambda = ginadapter.New(router)
		
		log.Info().Msg("Lambda handler initialized successfully")
	})

	// Proxy the request to Gin router
	resp, err := ginLambda.ProxyWithContext(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Error processing request")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"internal server error"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, err
	}

	return resp, nil
}
