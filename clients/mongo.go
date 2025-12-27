package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var MongoDBClient *MongoConnection

type MongoConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Debug    bool
	Auth     bool
}

type MongoConnection struct {
	Config *MongoConfig
	Client *mongo.Client
	DBName string
}

func NewMongoConnection(config *MongoConfig) *MongoConnection {
	return &MongoConnection{
		Config: config,
		DBName: "mongo_message_database",
	}
}

func (c *MongoConnection) Open() {
	var url = fmt.Sprintf("mongodb://%s:%s", c.Config.Host, c.Config.Port)
	var clientOpts = options.Client().ApplyURI(url)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if c.Config.Auth {
		clientOpts = clientOpts.SetAuth(options.Credential{
			Username: c.Config.User,
			Password: c.Config.Password,
		})
	}

	if c.Config.Debug {
		cmdMonitor := &event.CommandMonitor{
			Started: func(_ context.Context, cse *event.CommandStartedEvent) {
				fmt.Println(cse.Command)
			},
		}
		clientOpts.SetMonitor(cmdMonitor)
	}
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Error().Msgf("unable to connect mongodb with error %v", err.Error())
		return
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error().Msgf("unable to ping mongodb with error %v", err.Error())
		return
	}
	c.Client = client
	log.Info().Msgf("Mongo DB Connected Successfully")
}

func (c *MongoConnection) CloseConnection() {
	defer func() {
		if err := c.Client.Disconnect(context.TODO()); err != nil {
			log.Error().Msgf("unable to discunnect mongo with %v", err.Error())
			return
		}
	}()
}
