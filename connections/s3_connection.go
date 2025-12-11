package connections

import (
	"vivek-ray/clients"
	"vivek-ray/conf"
)

var S3Connection *clients.S3Connection

func InitS3() {
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
}
