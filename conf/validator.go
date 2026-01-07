package conf

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	MissingVars []string
	Message     string
}

func (e *ValidationError) Error() string {
	if len(e.MissingVars) > 0 {
		return fmt.Sprintf("%s. Missing variables: %s", e.Message, strings.Join(e.MissingVars, ", "))
	}
	return e.Message
}

// ValidateConfig validates that all required configuration is present
func ValidateConfig() error {
	var missingVars []string

	// Validate application config
	if AppConfig.APIKey == "" {
		missingVars = append(missingVars, "API_KEY")
	}

	// Validate database config
	// Either full connection string OR individual components must be provided
	hasConnectionString := DatabaseConfig.PgSQLConnection != ""
	hasComponents := DatabaseConfig.PgSQLHost != "" && 
					 DatabaseConfig.PgSQLDatabase != "" && 
					 DatabaseConfig.PgSQLUser != "" && 
					 DatabaseConfig.PgSQLPassword != ""

	if !hasConnectionString && !hasComponents {
		missingVars = append(missingVars, "PG_DB_CONNECTION or (PG_DB_HOST, PG_DB_DATABASE, PG_DB_USERNAME, PG_DB_PASSWORD)")
	}

	// Validate Elasticsearch config
	// Either full connection string OR individual components must be provided
	hasESConnectionString := SearchEngineConfig.ElasticsearchConnection != ""
	hasESComponents := SearchEngineConfig.ElasticsearchHost != "" && 
					   SearchEngineConfig.ElasticsearchPort != ""

	if !hasESConnectionString && !hasESComponents {
		missingVars = append(missingVars, "ELASTICSEARCH_CONNECTION or (ELASTICSEARCH_HOST, ELASTICSEARCH_PORT)")
	}

	// If Elasticsearch auth is enabled, username and password are required
	if SearchEngineConfig.ElasticsearchAuth {
		if SearchEngineConfig.ElasticsearchUser == "" {
			missingVars = append(missingVars, "ELASTICSEARCH_USERNAME")
		}
		if SearchEngineConfig.ElasticsearchPassword == "" {
			missingVars = append(missingVars, "ELASTICSEARCH_PASSWORD")
		}
	}

	// Validate S3 config
	if S3StorageConfig.S3AccessKey == "" {
		missingVars = append(missingVars, "S3_ACCESS_KEY")
	}
	if S3StorageConfig.S3SecretKey == "" {
		missingVars = append(missingVars, "S3_SECRET_KEY")
	}
	if S3StorageConfig.S3Region == "" {
		missingVars = append(missingVars, "S3_REGION")
	}
	if S3StorageConfig.S3Bucket == "" {
		missingVars = append(missingVars, "S3_BUCKET")
	}

	if len(missingVars) > 0 {
		err := &ValidationError{
			MissingVars: missingVars,
			Message:     "Required configuration variables are missing",
		}
		log.Error().Msgf("Configuration validation failed: %s", err.Error())
		return err
	}

	log.Info().Msg("Configuration validation passed")
	return nil
}
