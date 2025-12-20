package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"vivek-ray/conf"
	"vivek-ray/modules/users/repository"
	"vivek-ray/utilities"

	"github.com/rs/zerolog/log"
)

const (
	maxAvatarSize = 5 * 1024 * 1024 // 5MB
)

type AvatarService struct {
	profileRepo *repository.UserProfileRepository
	s3Service   *utilities.S3Service
}

func NewAvatarService() *AvatarService {
	return &AvatarService{
		profileRepo: repository.NewUserProfileRepository(),
		s3Service:   utilities.NewS3Service(),
	}
}

// UploadAvatar handles avatar file upload (S3 or local storage)
func (s *AvatarService) UploadAvatar(ctx context.Context, userID string, filename string, fileContent []byte, contentType string) (string, error) {
	// Validate file
	if err := utilities.ValidateAvatarFile(filename, fileContent, maxAvatarSize); err != nil {
		return "", err
	}

	// Generate unique filename
	uniqueFilename := utilities.GenerateUniqueFilename(userID, filename)

	var avatarURL string
	var fullURL string

	// Try S3 first if configured
	if s.s3Service.IsS3Configured() {
		s3Key := s.s3Service.GetAvatarsPrefix() + uniqueFilename
		if err := s.s3Service.UploadFile(ctx, fileContent, s3Key, contentType); err != nil {
			log.Warn().Err(err).Msg("Failed to upload to S3, falling back to local storage")
			} else {
			avatarURL = s3Key
			// Generate presigned URL or public URL
			if presignedURL, err := s.s3Service.GeneratePresignedURL(ctx, s3Key, 7*24*time.Hour); err == nil {
				fullURL = presignedURL
			} else {
				fullURL = s.s3Service.GetPublicURL(s3Key)
			}
			// Successfully uploaded to S3, return
			return fullURL, nil
		}
	}

	// Fallback to local storage
	if avatarURL == "" {
		uploadDir := conf.AppConfig.UploadDir
		if uploadDir == "" {
			uploadDir = "uploads"
		}
		avatarsDir := filepath.Join(uploadDir, "avatars")
		if err := os.MkdirAll(avatarsDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create avatars directory: %w", err)
		}

		filePath := filepath.Join(avatarsDir, uniqueFilename)
		if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
			return "", fmt.Errorf("failed to save file: %w", err)
		}

		avatarURL = "/media/avatars/" + uniqueFilename
		baseURL := conf.AppConfig.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:8000"
		}
		fullURL = baseURL + avatarURL
	}

	// Update profile with avatar URL
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	if profile == nil {
		return "", errors.New("user profile not found")
	}

	// Delete old avatar if exists
	if profile.AvatarURL != nil && *profile.AvatarURL != "" {
		s.deleteOldAvatar(ctx, *profile.AvatarURL)
	}

	now := time.Now()
	profile.AvatarURL = &avatarURL
	profile.UpdatedAt = &now

	if err := s.profileRepo.UpdateProfile(ctx, profile); err != nil {
		return "", err
	}

	return fullURL, nil
}

// deleteOldAvatar deletes old avatar file (S3 or local)
func (s *AvatarService) deleteOldAvatar(ctx context.Context, oldAvatarURL string) {
	if s.s3Service.IsS3Key(oldAvatarURL) {
		if err := s.s3Service.DeleteFile(ctx, oldAvatarURL); err != nil {
			log.Warn().Err(err).Msg("Failed to delete old avatar from S3")
		}
	} else if !isExternalURL(oldAvatarURL) {
		// Local file
		uploadDir := conf.AppConfig.UploadDir
		if uploadDir == "" {
			uploadDir = "uploads"
		}
		// Extract filename from URL
		filename := filepath.Base(oldAvatarURL)
		filePath := filepath.Join(uploadDir, "avatars", filename)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			log.Warn().Err(err).Msg("Failed to delete old avatar file")
		}
	}
}

func isExternalURL(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}

// ReadAvatarFile reads avatar file content (for serving)
func (s *AvatarService) ReadAvatarFile(ctx context.Context, avatarURL string) (io.ReadCloser, error) {
	if s.s3Service.IsS3Key(avatarURL) {
		// Read from S3
		return s.s3Service.ReadFileStream(ctx, avatarURL)
	}

	// Read from local storage
	uploadDir := conf.AppConfig.UploadDir
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	filename := filepath.Base(avatarURL)
	filePath := filepath.Join(uploadDir, "avatars", filename)
	return os.Open(filePath)
}

