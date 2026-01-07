package helper

import (
	"fmt"
	"net/http"
	"strings"
	"vivek-ray/constants"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	ErrorCode string `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
	Details   any    `json:"details,omitempty"`
}

// ErrorCode constants
const (
	ErrorCodeValidation      = "VALIDATION_ERROR"
	ErrorCodeNotFound        = "NOT_FOUND"
	ErrorCodeUnauthorized    = "UNAUTHORIZED"
	ErrorCodeForbidden       = "FORBIDDEN"
	ErrorCodeInternalError   = "INTERNAL_ERROR"
	ErrorCodeBadRequest      = "BAD_REQUEST"
	ErrorCodeConflict        = "CONFLICT"
	ErrorCodeRateLimit       = "RATE_LIMIT_EXCEEDED"
	ErrorCodeElasticsearch   = "ELASTICSEARCH_ERROR"
	ErrorCodeDatabase        = "DATABASE_ERROR"
)

// SendErrorResponse sends a standardized error response
func SendErrorResponse(c *gin.Context, statusCode int, errorCode, message string, details any) {
	response := ErrorResponse{
		Success:   false,
		Error:     message,
		ErrorCode: errorCode,
		Message:   message,
		Details:   details,
	}

	// Log error for debugging (except for client errors)
	if statusCode >= http.StatusInternalServerError {
		log.Error().
			Str("error_code", errorCode).
			Str("message", message).
			Interface("details", details).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Msg("Internal server error")
	}

	c.JSON(statusCode, response)
}

// SendValidationError sends a 400 Bad Request error for validation failures
func SendValidationError(c *gin.Context, message string, details any) {
	SendErrorResponse(c, http.StatusBadRequest, ErrorCodeValidation, message, details)
}

// SendNotFoundError sends a 404 Not Found error
func SendNotFoundError(c *gin.Context, resourceType, identifier string) {
	message := fmt.Sprintf("%s with identifier '%s' not found", resourceType, identifier)
	SendErrorResponse(c, http.StatusNotFound, ErrorCodeNotFound, message, nil)
}

// SendInternalError sends a 500 Internal Server Error
func SendInternalError(c *gin.Context, message string, err error) {
	details := map[string]string{}
	if err != nil {
		details["internal_error"] = err.Error()
	}
	SendErrorResponse(c, http.StatusInternalServerError, ErrorCodeInternalError, message, details)
}

// SendBadRequestError sends a 400 Bad Request error
func SendBadRequestError(c *gin.Context, message string) {
	SendErrorResponse(c, http.StatusBadRequest, ErrorCodeBadRequest, message, nil)
}

// SendUnauthorizedError sends a 401 Unauthorized error
func SendUnauthorizedError(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized: Invalid or missing API key"
	}
	SendErrorResponse(c, http.StatusUnauthorized, ErrorCodeUnauthorized, message, nil)
}

// SendForbiddenError sends a 403 Forbidden error
func SendForbiddenError(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden: Insufficient permissions"
	}
	SendErrorResponse(c, http.StatusForbidden, ErrorCodeForbidden, message, nil)
}

// SendConflictError sends a 409 Conflict error
func SendConflictError(c *gin.Context, message string, details any) {
	SendErrorResponse(c, http.StatusConflict, ErrorCodeConflict, message, details)
}

// SendRateLimitError sends a 429 Too Many Requests error
func SendRateLimitError(c *gin.Context, message string) {
	if message == "" {
		message = "Rate limit exceeded. Please try again later."
	}
	SendErrorResponse(c, http.StatusTooManyRequests, ErrorCodeRateLimit, message, nil)
}

// HandleServiceError handles service layer errors and converts them to appropriate HTTP responses
func HandleServiceError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Check for known error types
	switch err {
	case constants.CompanyNotFoundError:
		SendNotFoundError(c, "Company", "specified UUID")
		return
	case constants.ContactNotFoundError:
		SendNotFoundError(c, "Contact", "specified UUID")
		return
	case constants.PageSizeExceededError:
		SendBadRequestError(c, err.Error())
		return
	case constants.PageNumberExceededError:
		SendBadRequestError(c, err.Error())
		return
	case constants.FailedToFetchDataError:
		SendInternalError(c, "Failed to fetch data from database", err)
		return
	}

	// Check error message for common patterns
	errMsg := err.Error()
	if contains(errMsg, "validation") || contains(errMsg, "invalid") || contains(errMsg, "required") {
		SendValidationError(c, errMsg, nil)
		return
	}

	if contains(errMsg, "not found") {
		SendNotFoundError(c, "Resource", "specified identifier")
		return
	}

	// Default to internal server error
	SendInternalError(c, "An unexpected error occurred", err)
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
