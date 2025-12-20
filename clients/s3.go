package clients

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

type S3Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
	Endpoint  string
	UseSSL    bool
	Debug     bool
}

type S3Connection struct {
	Config *S3Config
	Client *s3.Client
}

func NewS3Connection(cfg *S3Config) *S3Connection {
	return &S3Connection{
		Config: cfg,
	}
}

func (c *S3Connection) Open() {
	if c.Client != nil {
		log.Error().Msg("Already open S3 connection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if c.Config.Endpoint != "" {
				protocol := "http"
				if c.Config.UseSSL {
					protocol = "https"
				}
				return aws.Endpoint{
					URL:               protocol + "://" + c.Config.Endpoint,
					SigningRegion:     c.Config.Region,
					HostnameImmutable: true,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		},
	)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(c.Config.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				c.Config.AccessKey,
				c.Config.SecretKey,
				"",
			),
		),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Error().Err(err).Msg("Error loading AWS config")
		return
	}

	c.Client = s3.NewFromConfig(cfg)

	_, err = c.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to S3")
		return
	}

	log.Info().Msgf("S3 Connected Successfully")
}

func (c *S3Connection) Close() {
	c.Client = nil
	log.Info().Msg("S3 connection closed")
}

func (c *S3Connection) ListFiles(ctx context.Context, bucket, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	result, err := c.Client.ListObjectsV2(ctx, input)
	if err != nil {
		log.Error().Err(err).Msgf("Error listing files from S3 with prefix: %s", prefix)
		return nil, err
	}

	var keys []string
	for _, obj := range result.Contents {
		keys = append(keys, *obj.Key)
	}

	log.Debug().Msgf("Successfully listed %d files from S3 with prefix: %s", len(keys), prefix)
	return keys, nil
}

func (c *S3Connection) FileExists(ctx context.Context, bucket, key string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := c.Client.HeadObject(ctx, input)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (c *S3Connection) GetPresignedURL(ctx context.Context, bucket, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.Client)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	presignedURL, err := presignClient.PresignGetObject(ctx, input, s3.WithPresignExpires(expiry))
	if err != nil {
		log.Error().Err(err).Msgf("Error generating presigned URL for: %s", key)
		return "", err
	}

	log.Debug().Msgf("Successfully generated presigned URL for: %s", key)
	return presignedURL.URL, nil
}

func (c *S3Connection) ReadFileStream(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := c.Client.GetObject(ctx, input)
	if err != nil {
		log.Error().Err(err).Msgf("Error reading file stream from S3: %s", key)
		return nil, err
	}

	log.Debug().Msgf("Successfully opened file stream from S3: %s", key)
	return result.Body, nil
}
