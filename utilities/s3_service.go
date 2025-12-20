package utilities

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"vivek-ray/connections"
	"vivek-ray/conf"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

var (
	ErrS3NotConfigured = errors.New("S3 not configured")
	ErrS3UploadFailed    = errors.New("S3 upload failed")
)

type S3Service struct {
	bucketName string
	prefix     string
}

func NewS3Service() *S3Service {
	return &S3Service{
		bucketName: conf.S3StorageConfig.S3Bucket,
		prefix:     "avatars/",
	}
}

// IsS3Configured checks if S3 is properly configured
func (s *S3Service) IsS3Configured() bool {
	return s.bucketName != "" && connections.S3Connection != nil && connections.S3Connection.Client != nil
}

// IsS3Key checks if a path is an S3 key
func (s *S3Service) IsS3Key(path string) bool {
	return IsS3Key(path)
}

// UploadFile uploads a file to S3
func (s *S3Service) UploadFile(ctx context.Context, fileContent []byte, s3Key string, contentType string) error {
	if !s.IsS3Configured() {
		return ErrS3NotConfigured
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String(contentType),
	}

	_, err := connections.S3Connection.Client.PutObject(ctx, input)
	if err != nil {
		log.Error().Err(err).Msgf("Error uploading file to S3: %s", s3Key)
		return ErrS3UploadFailed
	}

	log.Debug().Msgf("Successfully uploaded file to S3: %s", s3Key)
	return nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(ctx context.Context, s3Key string) error {
	if !s.IsS3Configured() {
		return ErrS3NotConfigured
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s3Key),
	}

	_, err := connections.S3Connection.Client.DeleteObject(ctx, input)
	if err != nil {
		log.Error().Err(err).Msgf("Error deleting file from S3: %s", s3Key)
		return err
	}

	log.Debug().Msgf("Successfully deleted file from S3: %s", s3Key)
	return nil
}

// GeneratePresignedURL generates a presigned URL for an S3 object
func (s *S3Service) GeneratePresignedURL(ctx context.Context, s3Key string, expiry time.Duration) (string, error) {
	if !s.IsS3Configured() {
		return "", ErrS3NotConfigured
	}

	return connections.S3Connection.GetPresignedURL(ctx, s.bucketName, s3Key, expiry)
}

// GetPublicURL generates a public URL for an S3 object
func (s *S3Service) GetPublicURL(s3Key string) string {
	if !s.IsS3Configured() {
		return ""
	}

	// Construct public URL based on endpoint
	endpoint := conf.S3StorageConfig.S3Endpoint
	protocol := "https"
	if !conf.S3StorageConfig.S3SSL {
		protocol = "http"
	}

	if endpoint != "" {
		return protocol + "://" + endpoint + "/" + s.bucketName + "/" + s3Key
	}

	// Default AWS S3 URL format
	region := conf.S3StorageConfig.S3Region
	if region == "" {
		region = "us-east-1"
	}
	return "https://" + s.bucketName + ".s3." + region + ".amazonaws.com/" + s3Key
}

// GetAvatarsPrefix returns the prefix for avatar files
func (s *S3Service) GetAvatarsPrefix() string {
	return s.prefix
}

// ExtractS3KeyFromURL extracts S3 key from a full URL
func (s *S3Service) ExtractS3KeyFromURL(url string) string {
	// Remove protocol
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	// Remove domain and bucket
	parts := strings.Split(url, "/")
	if len(parts) > 1 {
		// Skip bucket name (first part after domain)
		return strings.Join(parts[1:], "/")
	}

	return url
}

// ReadFileStream reads a file stream from S3
func (s *S3Service) ReadFileStream(ctx context.Context, s3Key string) (io.ReadCloser, error) {
	if !s.IsS3Configured() {
		return nil, ErrS3NotConfigured
	}

	return connections.S3Connection.ReadFileStream(ctx, s.bucketName, s3Key)
}

