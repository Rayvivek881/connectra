package conf

import (
	"os"

	"github.com/rs/zerolog/log"
)

// SetLambdaDefaults sets Lambda-optimized default values
func SetLambdaDefaults() {
	if !IsLambdaMode() {
		return
	}

	// Set Lambda-optimized connection pool sizes via environment variables
	// These will be picked up by the connection initialization code
	if os.Getenv("PG_DB_MAX_OPEN_CONNS") == "" {
		os.Setenv("PG_DB_MAX_OPEN_CONNS", "5")
	}
	if os.Getenv("PG_DB_MAX_IDLE_CONNS") == "" {
		os.Setenv("PG_DB_MAX_IDLE_CONNS", "2")
	}
	if os.Getenv("ES_MAX_IDLE_CONNS") == "" {
		os.Setenv("ES_MAX_IDLE_CONNS", "3")
	}
	if os.Getenv("ES_MAX_IDLE_CONNS_PER_HOST") == "" {
		os.Setenv("ES_MAX_IDLE_CONNS_PER_HOST", "1")
	}

	log.Info().Msg("Lambda-optimized defaults applied")
}

// GetLambdaFunctionName returns the Lambda function name if running in Lambda
func GetLambdaFunctionName() string {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
}

// GetLambdaRegion returns the AWS region if running in Lambda
func GetLambdaRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	return region
}
