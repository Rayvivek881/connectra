package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"vivek-ray/conf"
	"vivek-ray/connections"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const uploadURLTTL = 24 * time.Hour

func GetUploadURL(c *gin.Context) {
	fileName := c.Query("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename query parameter is required", "success": false})
		return
	}

	s3Key := fmt.Sprintf("%s/%s_%s", conf.S3StorageConfig.S3UploadFilePath, uuid.New().String(), fileName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	presignedURL, err := connections.S3Connection.GetUploadPresignedURL(
		ctx,
		conf.S3StorageConfig.S3Bucket,
		s3Key,
		uploadURLTTL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate upload URL", "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"upload_url": presignedURL,
		"s3_key":     s3Key,
		"expires_in": uploadURLTTL.String(),
	})
}
