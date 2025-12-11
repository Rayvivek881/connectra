package connections

import (
	"vivek-ray/clients"
	"vivek-ray/conf"
)

var ElasticsearchConnection *clients.ElasticsearchConnection

func InitSearchEngine() {
	ElasticsearchConnection = clients.NewElasticsearchConnection(&clients.ElasticsearchConfig{
		User:     conf.SearchEngineConfig.ElasticsearchUser,
		Password: conf.SearchEngineConfig.ElasticsearchPassword,
		Host:     conf.SearchEngineConfig.ElasticsearchHost,
		Port:     conf.SearchEngineConfig.ElasticsearchPort,
		Debug:    conf.SearchEngineConfig.ElasticsearchDebug,
		Auth:     conf.SearchEngineConfig.ElasticsearchAuth,
		SSL:      conf.SearchEngineConfig.ElasticsearchSSL,
	})
	ElasticsearchConnection.Open()
}
