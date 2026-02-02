package connections

import (
	"vivek-ray/clients"
	"vivek-ray/conf"
)

var OpenSearchConnection *clients.OpenSearchConnection

func InitSearchEngine() {
	OpenSearchConnection = clients.NewOpenSearchConnection(&clients.OpenSearchConfig{
		User:     conf.SearchEngineConfig.OpenSearchUser,
		Password: conf.SearchEngineConfig.OpenSearchPassword,
		Host:     conf.SearchEngineConfig.OpenSearchHost,
		Port:     conf.SearchEngineConfig.OpenSearchPort,
		Debug:    conf.SearchEngineConfig.OpenSearchDebug,
		Auth:     conf.SearchEngineConfig.OpenSearchAuth,
		SSL:      conf.SearchEngineConfig.OpenSearchSSL,
	})
	OpenSearchConnection.Open()
}

func CloseSearchEngine() {
	if OpenSearchConnection != nil {
		OpenSearchConnection.Close()
	}
}
