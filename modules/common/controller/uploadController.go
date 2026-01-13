package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vivek-ray/conf"
	"vivek-ray/connections"
	"vivek-ray/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var URLTTL = time.Duration(conf.S3StorageConfig.S3URLTTL) * time.Hour

func GetUploadURL(c *gin.Context) {
	fileName := c.Query("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.FilenameRequiredError.Error(), "success": false})
		return
	}
	bucket := c.Query("bucket")
	if bucket == "" {
		bucket = conf.S3StorageConfig.S3Bucket
	}

	s3Key := fmt.Sprintf("%s/%s_%s", conf.S3StorageConfig.S3UploadFilePath, uuid.New().String(), fileName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	presignedURL, err := connections.S3Connection.GetUploadPresignedURL(
		ctx,
		bucket,
		s3Key,
		URLTTL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.FailedToGenerateUploadURLError.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"upload_url": presignedURL,
		"s3_key":     s3Key,
		"expires_in": URLTTL.String(),
	})
}

func GetDownloadURL(c *gin.Context) {
	s3Key := c.Query("s3_key")
	if s3Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.S3KeyRequiredError.Error(), "success": false})
		return
	}
	bucket := c.Query("bucket")
	if bucket == "" {
		bucket = conf.S3StorageConfig.S3Bucket
	}

	presignedURL, err := connections.S3Connection.GetDownloadPresignedURL(
		context.Background(),
		bucket,
		s3Key,
		URLTTL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.FailedToGenerateDownloadURLError.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"download_url": presignedURL,
		"s3_key":       s3Key,
		"expires_in":   URLTTL.String(),
	})
}
