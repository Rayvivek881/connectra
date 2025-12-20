package utilities

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
	ErrFileTooLarge    = errors.New("file too large")
	ErrInvalidImage    = errors.New("file does not appear to be a valid image file")
)

// Allowed image extensions
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// ValidateFileExtension checks if the file extension is allowed
func ValidateFileExtension(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return ErrInvalidFileType
	}
	return nil
}

// ValidateImageMagicBytes validates image file by checking magic bytes (file signature)
func ValidateImageMagicBytes(data []byte) error {
	if len(data) < 12 {
		return ErrInvalidImage
	}

	// JPEG: starts with FF D8 FF
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return nil
	}

	// PNG: starts with 89 50 4E 47
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return nil
	}

	// GIF87a or GIF89a
	if bytes.HasPrefix(data, []byte("GIF87a")) || bytes.HasPrefix(data, []byte("GIF89a")) {
		return nil
	}

	// WebP: starts with RIFF and contains WEBP
	if bytes.HasPrefix(data, []byte("RIFF")) && bytes.Contains(data[:12], []byte("WEBP")) {
		return nil
	}

	return ErrInvalidImage
}

// ValidateFileSize checks if file size is within limit (5MB)
func ValidateFileSize(size int64, maxSize int64) error {
	if size > maxSize {
		return ErrFileTooLarge
	}
	return nil
}

// GenerateUniqueFilename generates a unique filename with timestamp
func GenerateUniqueFilename(userID, originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().UTC().Format("20060102T150405000000")
	return userID + "_" + timestamp + ext
}

// ValidateAvatarFile validates avatar file (extension, magic bytes, and size)
func ValidateAvatarFile(filename string, data []byte, maxSize int64) error {
	// Validate extension
	if err := ValidateFileExtension(filename); err != nil {
		return err
	}

	// Validate magic bytes
	if err := ValidateImageMagicBytes(data); err != nil {
		return err
	}

	// Validate size
	if err := ValidateFileSize(int64(len(data)), maxSize); err != nil {
		return err
	}

	return nil
}

