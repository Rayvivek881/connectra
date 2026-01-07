package main

import (
	"os"
	"vivek-ray/lambda"

	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set Lambda mode environment variable
	os.Setenv("LAMBDA_MODE", "true")

	// Initialize connections once during cold start
	// These connections will be reused across invocations
	if err := lambda.InitConnections(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize connections during Lambda startup")
		// Continue anyway - connections might be initialized in handler
	}

	// Start Lambda runtime
	// The handler function will be called for each API Gateway event
	awslambda.Start(lambda.Handler)
}
