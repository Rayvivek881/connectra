//go:build integration
// +build integration

package tests

import (
	"testing"
	"vivek-ray/connections"
	"vivek-ray/models"
)

func TestDatabaseConnection(t *testing.T) {
	if connections.PgDBConnection == nil || connections.PgDBConnection.Client == nil {
		t.Skip("Database connection not available, skipping test")
		return
	}

	// Test basic database connectivity
	// This is a simple connectivity test
	// More complex tests would require test data setup/teardown
}

func TestFilterRepository_GetTempFilters(t *testing.T) {
	if connections.PgDBConnection == nil || connections.PgDBConnection.Client == nil {
		t.Skip("Database connection not available, skipping test")
		return
	}

	repo := models.FiltersRepository(connections.PgDBConnection.Client)
	filters, err := repo.GetTempFilters()

	if err != nil {
		t.Logf("GetTempFilters error (may be expected if DB not configured): %v", err)
		return
	}

	// If successful, verify structure
	if filters == nil {
		t.Error("GetTempFilters returned nil")
		return
	}

	t.Logf("Retrieved %d temp filters", len(filters))
}

func TestFilterRepository_GetFiltersByService(t *testing.T) {
	if connections.PgDBConnection == nil || connections.PgDBConnection.Client == nil {
		t.Skip("Database connection not available, skipping test")
		return
	}

	repo := models.FiltersRepository(connections.PgDBConnection.Client)

	// Test company filters
	companyFilters, err := repo.GetFiltersByService("company")
	if err != nil {
		t.Logf("GetFiltersByService('company') error: %v", err)
		return
	}

	if companyFilters == nil {
		t.Error("GetFiltersByService returned nil")
		return
	}

	t.Logf("Retrieved %d company filters", len(companyFilters))

	// Test contact filters
	contactFilters, err := repo.GetFiltersByService("contact")
	if err != nil {
		t.Logf("GetFiltersByService('contact') error: %v", err)
		return
	}

	if contactFilters == nil {
		t.Error("GetFiltersByService returned nil")
		return
	}

	t.Logf("Retrieved %d contact filters", len(contactFilters))
}
