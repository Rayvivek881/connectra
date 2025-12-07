package connections

import (
	"vivek-ray/clients"
	"vivek-ray/conf"
)

var PgDBConnection *clients.PgsqlConnection

func InitDB() {
	PgDBConnection = clients.NewPgsqlConnection(&clients.PgsqlConfig{
		User:     conf.DatabaseConfig.PgSQLUser,
		Password: conf.DatabaseConfig.PgSQLPassword,
		Host:     conf.DatabaseConfig.PgSQLHost,
		Port:     conf.DatabaseConfig.PgSQLPort,
		Database: conf.DatabaseConfig.PgSQLDatabase,
		Debug:    conf.DatabaseConfig.PgSQLDebug,
	})
	PgDBConnection.Open()
}

func CloseDB() {
	PgDBConnection.Close()
}
