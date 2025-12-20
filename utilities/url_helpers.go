package utilities

import (
	"strings"
)

// GetFullAvatarURL converts a relative avatar URL to a full URL
func GetFullAvatarURL(avatarURL *string) *string {
	if avatarURL == nil || *avatarURL == "" {
		return nil
	}

	url := *avatarURL

	// If already a full URL, return as-is
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return &url
	}

	// Check if it's an S3 key (starts with avatars/ or similar)
	// For now, we'll handle S3 URLs in the service layer
	// Here we just handle relative paths

	// Remove leading slash if present
	path := strings.TrimPrefix(url, "/")
	
	// Get base URL from config (would need to add this to config)
	// For now, return relative path - service layer will handle full URL generation
	baseURL := "" // TODO: Get from config
	if baseURL != "" {
		baseURL = strings.TrimSuffix(baseURL, "/")
		fullURL := baseURL + "/" + path
		return &fullURL
	}

	return &url
}

// IsS3Key checks if a path is an S3 key
func IsS3Key(path string) bool {
	// S3 keys typically don't start with http:// or https://
	// and don't start with /
	return !strings.HasPrefix(path, "http://") &&
		!strings.HasPrefix(path, "https://") &&
		!strings.HasPrefix(path, "/")
}

