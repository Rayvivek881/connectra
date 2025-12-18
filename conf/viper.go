package conf

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Viper struct{}

type app struct {
	ENV                  string `mapstructure:"APP_ENV"`
	ParallelJobs         int    `mapstructure:"PARALLEL_JOBS"`
	InQueueSize          int    `mapstructure:"IN_QUEUE_SIZE"`
	TickerInterval       int    `mapstructure:"TICKER_INTERVAL_MINUTES"`
	APIKey               string `mapstructure:"API_KEY"`
	MaxRequestsPerMinute int    `mapstructure:"MAX_REQUESTS_PER_MINUTE"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
}

type database struct {
	PgSQLConnection string `mapstructure:"PG_DB_CONNECTION"`
	PgSQLHost       string `mapstructure:"PG_DB_HOST"`
	PgSQLPort       string `mapstructure:"PG_DB_PORT"`
	PgSQLDatabase   string `mapstructure:"PG_DB_DATABASE"`
	PgSQLUser       string `mapstructure:"PG_DB_USERNAME"`
	PgSQLPassword   string `mapstructure:"PG_DB_PASSWORD"`
	PgSQLDebug      bool   `mapstructure:"PG_DB_DEBUG"`
	PgSQLSSL        bool   `mapstructure:"PG_DB_SSL"`
}

type searchEngine struct {
	ElasticsearchConnection string `mapstructure:"ELASTICSEARCH_CONNECTION"`
	ElasticsearchHost       string `mapstructure:"ELASTICSEARCH_HOST"`
	ElasticsearchPort       string `mapstructure:"ELASTICSEARCH_PORT"`
	ElasticsearchUser       string `mapstructure:"ELASTICSEARCH_USERNAME"`
	ElasticsearchPassword   string `mapstructure:"ELASTICSEARCH_PASSWORD"`
	ElasticsearchDebug      bool   `mapstructure:"ELASTICSEARCH_DEBUG"`
	ElasticsearchSSL        bool   `mapstructure:"ELASTICSEARCH_SSL"`
	ElasticsearchAuth       bool   `mapstructure:"ELASTICSEARCH_AUTH"`
}

type s3Storage struct {
	S3AccessKey string `mapstructure:"S3_ACCESS_KEY"`
	S3SecretKey string `mapstructure:"S3_SECRET_KEY"`
	S3Region    string `mapstructure:"S3_REGION"`
	S3Bucket    string `mapstructure:"S3_BUCKET"`
	S3Endpoint  string `mapstructure:"S3_ENDPOINT"`
	S3SSL       bool   `mapstructure:"S3_SSL"`
	S3Debug     bool   `mapstructure:"S3_DEBUG"`
}

var AppConfig = &app{}
var DatabaseConfig = &database{}
var SearchEngineConfig = &searchEngine{}
var S3StorageConfig = &s3Storage{}

func (v *Viper) Init() {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	if err := viper.MergeInConfig(); err != nil {
		log.Fatal().Msgf("Viper merge failed: %v", err)
	}

	v.setDefaults()
	v.unmarshal(&AppConfig)
	v.unmarshal(&DatabaseConfig)
	v.unmarshal(&SearchEngineConfig)
	v.unmarshal(&S3StorageConfig)

	log.Info().Msgf("Viper initialized successfully")
}

func (v *Viper) setDefaults() {
	defer func() {
		if err := recover(); err != nil {
			log.Info().Msgf("Panic occurred: %v", err)
			return
		}
	}()
	structFields := [][]reflect.StructField{
		reflect.VisibleFields(reflect.TypeOf(struct{ app }{})),
		reflect.VisibleFields(reflect.TypeOf(struct{ database }{})),
		reflect.VisibleFields(reflect.TypeOf(struct{ searchEngine }{})),
		reflect.VisibleFields(reflect.TypeOf(struct{ s3Storage }{})),
	}
	v.setFields(structFields)
	log.Info().Msgf("Setting defaults for viper, completed")
}

func (v *Viper) setFields(fields [][]reflect.StructField) {
	for _, fs := range fields {
		for _, field := range fs {
			viper.SetDefault(field.Tag.Get("mapstructure"), "")
		}
	}
}

func (v *Viper) unmarshal(conf interface{}) {
	err := viper.Unmarshal(conf)
	if err != nil {
		log.Fatal().Msgf("Viper app unmarshal error: %s", err)
	}
}
