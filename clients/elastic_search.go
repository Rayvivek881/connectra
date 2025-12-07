package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

var ElasticsearchClient *ElasticsearchConnection

type ElasticsearchConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Debug    bool
	Auth     bool
}

type ElasticsearchConnection struct {
	Config *ElasticsearchConfig
	Client *elasticsearch.Client
}

func NewElasticsearchConnection(config *ElasticsearchConfig) *ElasticsearchConnection {
	return &ElasticsearchConnection{
		Config: config,
	}
}

func (c *ElasticsearchConnection) Open() {
	var addresses []string
	if c.Config.Port != "" {
		addresses = []string{fmt.Sprintf("https://%s:%s", c.Config.Host, c.Config.Port)}
	} else {
		addresses = []string{fmt.Sprintf("https://%s", c.Config.Host)}
	}

	cfg := elasticsearch.Config{
		Addresses: addresses,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	if c.Config.Auth {
		cfg.Username = c.Config.User
		cfg.Password = c.Config.Password
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Error().Msgf("unable to create elasticsearch client with error %v", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ping the Elasticsearch cluster by making a GET request to root endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		log.Error().Msgf("unable to create request with error %v", err.Error())
		return
	}

	res, err := client.Perform(req)
	if err != nil {
		log.Error().Msgf("unable to connect elasticsearch with error %v", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		log.Error().Msgf("elasticsearch error response: status code %d", res.StatusCode)
		return
	}

	c.Client = client
	log.Info().Msgf("Elasticsearch Connected Successfully")
}
