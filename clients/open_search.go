package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchtransport"
	"github.com/rs/zerolog/log"
)

var OpenSearchClient *OpenSearchConnection

type OpenSearchConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Debug    bool
	Auth     bool
	SSL      bool
}

type OpenSearchConnection struct {
	Config    *OpenSearchConfig
	Client    *opensearch.Client
	transport *http.Transport
}

func NewOpenSearchConnection(config *OpenSearchConfig) *OpenSearchConnection {
	return &OpenSearchConnection{
		Config: config,
	}
}

func (c *OpenSearchConnection) Open() {
	var addresses []string
	if c.Config.SSL {
		addresses = []string{fmt.Sprintf("https://%s:%s", c.Config.Host, c.Config.Port)}
	} else {
		addresses = []string{fmt.Sprintf("http://%s:%s", c.Config.Host, c.Config.Port)}
	}

	c.transport = &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		DisableCompression:    false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	if c.Config.SSL {
		c.transport.TLSClientConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		}
	}

	cfg := opensearch.Config{
		Addresses:            addresses,
		Transport:            c.transport,
		CompressRequestBody:  !c.Config.Debug,
		MaxRetries:           3,
		RetryOnStatus:        []int{502, 503, 504},
		EnableRetryOnTimeout: true,
	}

	if c.Config.Auth {
		cfg.Username = c.Config.User
		cfg.Password = c.Config.Password
	}

	if c.Config.Debug {
		cfg.Logger = &opensearchtransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to create OpenSearch client")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to OpenSearch cluster")
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Error().Msgf("OpenSearch cluster error: status %d", res.StatusCode)
		return
	}

	c.Client = client
	log.Info().Msg("OpenSearch connected successfully")
}

func (c *OpenSearchConnection) Close() {
	if c.transport != nil {
		c.transport.CloseIdleConnections()
	}
	c.Client = nil
	log.Info().Msg("OpenSearch connection closed")
}

func (c *OpenSearchConnection) IsConnected() bool {
	return c != nil && c.Client != nil
}
