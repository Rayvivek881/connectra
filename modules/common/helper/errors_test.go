package helper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSendErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name       string
		statusCode int
		errorCode  string
		message    string
		details    any
	}{
		{
			name:       "validation error",
			statusCode: http.StatusBadRequest,
			errorCode:  ErrorCodeValidation,
			message:    "Validation failed",
			details:    map[string]string{"field": "error"},
		},
		{
			name:       "not found error",
			statusCode: http.StatusNotFound,
			errorCode:  ErrorCodeNotFound,
			message:    "Resource not found",
			details:    nil,
		},
		{
			name:       "internal error",
			statusCode: http.StatusInternalServerError,
			errorCode:  ErrorCodeInternalError,
			message:    "Internal error",
			details:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			SendErrorResponse(c, tt.statusCode, tt.errorCode, tt.message, tt.details)
			
			if w.Code != tt.statusCode {
				t.Errorf("SendErrorResponse() status code = %v, want %v", w.Code, tt.statusCode)
			}
		})
	}
}

func TestSendValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	SendValidationError(c, "Invalid input", map[string]string{"field": "error"})
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("SendValidationError() status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestSendNotFoundError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	SendNotFoundError(c, "Company", "test-uuid")
	
	if w.Code != http.StatusNotFound {
		t.Errorf("SendNotFoundError() status code = %v, want %v", w.Code, http.StatusNotFound)
	}
}

func TestSendBadRequestError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	SendBadRequestError(c, "Bad request")
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("SendBadRequestError() status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestSendUnauthorizedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	SendUnauthorizedError(c, "Invalid API key")
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("SendUnauthorizedError() status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestSendRateLimitError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	SendRateLimitError(c, "Rate limit exceeded")
	
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("SendRateLimitError() status code = %v, want %v", w.Code, http.StatusTooManyRequests)
	}
}
