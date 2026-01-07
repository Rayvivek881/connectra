package connections

import (
	"sync"
	"vivek-ray/clients"
	"vivek-ray/conf"

	"github.com/rs/zerolog/log"
)

var (
	S3Connection *clients.S3Connection
	s3Once       sync.Once
)

// InitS3 initializes S3 connection using singleton pattern
// In Lambda, this ensures connections are reused across invocations
func InitS3() {
	s3Once.Do(func() {
		S3Connection = clients.NewS3Connection(&clients.S3Config{
			AccessKey: conf.S3StorageConfig.S3AccessKey,
			SecretKey: conf.S3StorageConfig.S3SecretKey,
			Region:    conf.S3StorageConfig.S3Region,
			Bucket:    conf.S3StorageConfig.S3Bucket,
			Endpoint:  conf.S3StorageConfig.S3Endpoint,
			UseSSL:    conf.S3StorageConfig.S3SSL,
			Debug:     conf.S3StorageConfig.S3Debug,
		})
		S3Connection.Open()
		log.Info().Msg("S3 connection initialized (singleton)")
	})
}

// GetS3 returns the S3 connection (thread-safe)
func GetS3() *clients.S3Connection {
	if S3Connection == nil {
		InitS3()
	}
	return S3Connection
}
