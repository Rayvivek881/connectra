package connections

import (
	"sync"
	"vivek-ray/clients"
	"vivek-ray/conf"
	"vivek-ray/utilities"

	"github.com/rs/zerolog/log"
)

var (
	ElasticsearchConnection *clients.ElasticsearchConnection
	esOnce                  sync.Once
)

// InitSearchEngine initializes Elasticsearch connection using singleton pattern
// In Lambda, this ensures connections are reused across invocations
func InitSearchEngine() {
	esOnce.Do(func() {
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

		// Initialize async indexing queues after Elasticsearch connection is established
		if ElasticsearchConnection.Client != nil {
			utilities.InitializeElasticsearchQueues(ElasticsearchConnection.Client)
		}
		log.Info().Msg("Elasticsearch connection initialized (singleton)")
	})
}

// GetSearchEngine returns the Elasticsearch connection (thread-safe)
func GetSearchEngine() *clients.ElasticsearchConnection {
	if ElasticsearchConnection == nil {
		InitSearchEngine()
	}
	return ElasticsearchConnection
}

func CloseSearchEngine() {
	// Stop queues gracefully
	if utilities.GetCompanyQueue() != nil {
		utilities.GetCompanyQueue().Stop()
	}
	if utilities.GetContactQueue() != nil {
		utilities.GetContactQueue().Stop()
	}

	if ElasticsearchConnection != nil {
		ElasticsearchConnection.Close()
		ElasticsearchConnection = nil
		// Reset once to allow re-initialization if needed
		esOnce = sync.Once{}
	}
}
