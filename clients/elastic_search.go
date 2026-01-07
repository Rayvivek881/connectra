package clients

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
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
	SSL      bool
}

type ElasticsearchConnection struct {
	Config    *ElasticsearchConfig
	Client    *elasticsearch.Client
	transport *http.Transport
}

func NewElasticsearchConnection(config *ElasticsearchConfig) *ElasticsearchConnection {
	return &ElasticsearchConnection{
		Config: config,
	}
}

func (c *ElasticsearchConnection) Open() {
	var addresses []string
	if c.Config.SSL {
		addresses = []string{fmt.Sprintf("https://%s:%s", c.Config.Host, c.Config.Port)}
	} else {
		addresses = []string{fmt.Sprintf("http://%s:%s", c.Config.Host, c.Config.Port)}
	}

	// Optimize connection pool based on execution environment
	isLambda := os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" || os.Getenv("LAMBDA_MODE") == "true"
	
	var maxIdleConns, maxIdleConnsPerHost int
	var idleConnTimeout time.Duration
	
	if isLambda {
		// Lambda-optimized pool sizes (smaller for serverless)
		maxIdleConns = 3        // Reduced from 20
		maxIdleConnsPerHost = 1 // Reduced from 5
		idleConnTimeout = 15 * time.Minute // Reduced from 30
		log.Info().Msg("Elasticsearch connection pool configured for Lambda (3 max idle, 1 per host)")
	} else {
		// Server mode - use larger pools
		maxIdleConns = 20
		maxIdleConnsPerHost = 5
		idleConnTimeout = 30 * time.Minute
		log.Info().Msg("Elasticsearch connection pool configured for server (20 max idle, 5 per host)")
	}

	c.transport = &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		IdleConnTimeout:     idleConnTimeout,
		DisableCompression:  false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	cfg := elasticsearch.Config{
		Addresses:           addresses,
		Transport:           c.transport,
		CompressRequestBody: !c.Config.Debug,
	}

	if c.Config.Auth {
		cfg.Username = c.Config.User
		cfg.Password = c.Config.Password
	}

	if c.Config.Debug {
		cfg.Logger = &elastictransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: false,
		}
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

func (c *ElasticsearchConnection) Close() {
	if c.transport != nil {
		c.transport.CloseIdleConnections()
	}
	c.Client = nil
	log.Info().Msg("Elasticsearch connection closed")
}
