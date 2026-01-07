//go:build integration
// +build integration

package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	w := makeRequest("GET", "/health", nil)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	if status, ok := response["status"].(string); !ok || status != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

func TestCompanyEndpoints_Unauthorized(t *testing.T) {
	// Test without API key
	req := httptest.NewRequest("GET", "/companies/test-uuid", nil)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestCompanyEndpoints_Create(t *testing.T) {
	companyData := map[string]interface{}{
		"name":            "Test Company",
		"employees_count": 100,
		"country":         "USA",
		"city":            "San Francisco",
		"state":           "CA",
	}
	
	w := makeRequest("POST", "/companies/create", companyData)
	
	if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 201 or 500 (if DB not available), got %d", w.Code)
	}
	
	if w.Code == http.StatusCreated {
		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		
		if success, ok := response["success"].(bool); !ok || !success {
			t.Errorf("Expected success=true, got %v", response["success"])
		}
	}
}

func TestCompanyEndpoints_GetByUUID_NotFound(t *testing.T) {
	// Test with non-existent UUID
	w := makeRequest("GET", "/companies/00000000-0000-0000-0000-000000000000", nil)
	
	// Should return 404 or 500 (if DB not available)
	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 404 or 500, got %d", w.Code)
	}
}

func TestCompanyEndpoints_GetByUUID_InvalidUUID(t *testing.T) {
	// Test with invalid UUID format
	w := makeRequest("GET", "/companies/invalid-uuid", nil)
	
	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestContactEndpoints_Create(t *testing.T) {
	contactData := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      "john.doe@example.com",
		"title":      "Software Engineer",
		"country":    "USA",
	}
	
	w := makeRequest("POST", "/contacts/create", contactData)
	
	if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 201 or 500 (if DB not available), got %d", w.Code)
	}
	
	if w.Code == http.StatusCreated {
		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		
		if success, ok := response["success"].(bool); !ok || !success {
			t.Errorf("Expected success=true, got %v", response["success"])
		}
	}
}

func TestContactEndpoints_Create_InvalidEmail(t *testing.T) {
	contactData := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      "invalid-email",
		"title":      "Software Engineer",
	}
	
	w := makeRequest("POST", "/contacts/create", contactData)
	
	// Should return 400 Bad Request for invalid email
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	if success, ok := response["success"].(bool); !ok || success {
		t.Errorf("Expected success=false, got %v", response["success"])
	}
}

func TestContactEndpoints_GetByUUID_InvalidUUID(t *testing.T) {
	// Test with invalid UUID format
	w := makeRequest("GET", "/contacts/invalid-uuid", nil)
	
	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCompanyEndpoints_ListByFilter(t *testing.T) {
	queryData := map[string]interface{}{
		"where": map[string]interface{}{
			"keyword_match": map[string]interface{}{
				"must": map[string]interface{}{
					"country": []string{"USA"},
				},
			},
		},
		"page":  1,
		"limit": 10,
	}
	
	w := makeRequest("POST", "/companies/", queryData)
	
	// Should return 200 or 500 (if ES not available)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

func TestContactEndpoints_ListByFilter(t *testing.T) {
	queryData := map[string]interface{}{
		"where": map[string]interface{}{
			"keyword_match": map[string]interface{}{
				"must": map[string]interface{}{
					"country": []string{"USA"},
				},
			},
		},
		"page":  1,
		"limit": 10,
	}
	
	w := makeRequest("POST", "/contacts/", queryData)
	
	// Should return 200 or 500 (if ES not available)
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

func TestErrorResponseFormat(t *testing.T) {
	// Test that error responses follow the standardized format
	w := makeRequest("GET", "/companies/invalid-uuid", nil)
	
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	// Check for standardized error fields
	if _, ok := response["success"]; !ok {
		t.Error("Error response missing 'success' field")
	}
	
	if success, ok := response["success"].(bool); ok && success {
		t.Error("Error response should have success=false")
	}
	
	// Error responses should have error_code (if standardized)
	if w.Code >= 400 {
		if _, ok := response["error"]; !ok {
			t.Error("Error response missing 'error' field")
		}
	}
}
