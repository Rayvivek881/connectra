package connections

import (
	"sync"
	"vivek-ray/clients"
	"vivek-ray/conf"

	"github.com/rs/zerolog/log"
)

var (
	PgDBConnection *clients.PgsqlConnection
	pgOnce         sync.Once
)

// InitDB initializes PostgreSQL connection using singleton pattern
// In Lambda, this ensures connections are reused across invocations
func InitDB() {
	pgOnce.Do(func() {
		PgDBConnection = clients.NewPgsqlConnection(&clients.PgsqlConfig{
			User:     conf.DatabaseConfig.PgSQLUser,
			Password: conf.DatabaseConfig.PgSQLPassword,
			Host:     conf.DatabaseConfig.PgSQLHost,
			Port:     conf.DatabaseConfig.PgSQLPort,
			Database: conf.DatabaseConfig.PgSQLDatabase,
			Debug:    conf.DatabaseConfig.PgSQLDebug,
		})
		PgDBConnection.Open()
		log.Info().Msg("PostgreSQL connection initialized (singleton)")
	})
}

// GetDB returns the PostgreSQL connection (thread-safe)
func GetDB() *clients.PgsqlConnection {
	if PgDBConnection == nil {
		InitDB()
	}
	return PgDBConnection
}

func CloseDB() {
	if PgDBConnection != nil {
		PgDBConnection.Close()
		PgDBConnection = nil
		// Reset once to allow re-initialization if needed
		pgOnce = sync.Once{}
	}
}
