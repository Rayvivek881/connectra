package clients

import (
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type PgsqlConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Debug    bool
}

type PgsqlConnection struct {
	config *PgsqlConfig
	Client *bun.DB
}

func NewPgsqlConnection(config *PgsqlConfig) *PgsqlConnection {
	return &PgsqlConnection{
		config: config,
	}
}

func (c *PgsqlConnection) Open() {
	if c.Client != nil {
		log.Error().Msg("Already open pgsql connection")
		return
	}
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		c.config.User, c.config.Password, c.config.Host, c.config.Port, c.config.Database)

	config, err := pgx.ParseConfig(uri)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing connection url")
		return
	}
	config.PreferSimpleProtocol = false
	sqldb := stdlib.OpenDB(*config)

	c.Client = bun.NewDB(sqldb, pgdialect.New())
	if c.config.Debug {
		c.Client.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(c.config.Debug)))
	}

	// Optimize connection pool based on execution environment
	isLambda := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" || os.Getenv("LAMBDA_MODE") == "true"
	
	if isLambda {
		// Lambda-optimized pool sizes (smaller for serverless)
		sqldb.SetMaxOpenConns(5)   // Reduced from 40
		sqldb.SetMaxIdleConns(2)  // Reduced from 20
		sqldb.SetConnMaxLifetime(15 * time.Minute) // Reduced from 30
		sqldb.SetConnMaxIdleTime(10 * time.Minute) // Reduced from 30
		log.Info().Msg("PostgreSQL connection pool configured for Lambda (5 max, 2 idle)")
	} else {
		// Server mode - use larger pools
		sqldb.SetMaxOpenConns(40)
		sqldb.SetMaxIdleConns(20)
		sqldb.SetConnMaxLifetime(30 * time.Minute)
		sqldb.SetConnMaxIdleTime(30 * time.Minute)
		log.Info().Msg("PostgreSQL connection pool configured for server (40 max, 20 idle)")
	}

	if err := c.Client.Ping(); err != nil {
		log.Error().Err(err).Msg("Error connecting to pgsql")
		return
	}
	log.Info().Msgf("PostgreSQL Connected Successfully")
}

func (c *PgsqlConnection) Close() {
	if err := c.Client.Close(); err != nil {
		log.Error().Msg("Already close pgsql connection")
	}
	c.Client = nil
	log.Info().Msg("PostgreSQL connection closed")
}
