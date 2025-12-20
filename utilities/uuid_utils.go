package utilities

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// IsValidUUID checks if a string is a valid UUID format
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// ParseUUID parses a UUID string and returns the UUID object
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

