package lambda

import (
	"os"
	"vivek-ray/conf"
	"vivek-ray/connections"

	"github.com/rs/zerolog/log"
)

// InitConnections initializes all database and service connections for Lambda
// This function is called once during Lambda cold start and connections are reused
// across invocations within the same container
func InitConnections() error {
	// Set Lambda mode before initialization
	os.Setenv("LAMBDA_MODE", "true")

	// Initialize configuration
	v := conf.Viper{}
	v.Init()

	// Validate configuration
	if err := conf.ValidateConfig(); err != nil {
		log.Error().Err(err).Msg("Configuration validation failed")
		return err
	}

	// Set Lambda-optimized defaults
	conf.SetLambdaDefaults()

	// Initialize all connections
	log.Info().Msg("Initializing connections for Lambda...")
	connections.InitDB()
	connections.InitSearchEngine()
	connections.InitS3()

	// Verify connections are healthy
	if err := verifyConnections(); err != nil {
		log.Error().Err(err).Msg("Connection verification failed")
		return err
	}

	log.Info().Msg("All connections initialized successfully for Lambda")
	return nil
}

// verifyConnections checks that all connections are healthy
func verifyConnections() error {
	// Verify PostgreSQL connection
	if connections.PgDBConnection == nil || connections.PgDBConnection.Client == nil {
		return &ConnectionError{Service: "PostgreSQL", Message: "connection is nil"}
	}
	if err := connections.PgDBConnection.Client.Ping(); err != nil {
		return &ConnectionError{Service: "PostgreSQL", Message: err.Error()}
	}

	// Verify Elasticsearch connection
	if connections.ElasticsearchConnection == nil || connections.ElasticsearchConnection.Client == nil {
		return &ConnectionError{Service: "Elasticsearch", Message: "connection is nil"}
	}
	// Elasticsearch ping is done during Open(), so we just check if client exists

	// Verify S3 connection
	if connections.S3Connection == nil || connections.S3Connection.Client == nil {
		return &ConnectionError{Service: "S3", Message: "connection is nil"}
	}

	return nil
}

// ConnectionError represents a connection initialization error
type ConnectionError struct {
	Service string
	Message string
}

func (e *ConnectionError) Error() string {
	return "failed to initialize " + e.Service + ": " + e.Message
}
