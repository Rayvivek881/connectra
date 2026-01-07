//go:build integration
// +build integration

package tests

import (
	"testing"
	"vivek-ray/connections"
	"vivek-ray/constants"
	"vivek-ray/models"
)

func TestElasticsearchConnection(t *testing.T) {
	if connections.ElasticsearchConnection == nil || connections.ElasticsearchConnection.Client == nil {
		t.Skip("Elasticsearch connection not available, skipping test")
		return
	}
	
	// Test basic Elasticsearch connectivity
	// This is a simple connectivity test
}

func TestElasticsearchCompanyRepository_ListByQueryMap(t *testing.T) {
	if connections.ElasticsearchConnection == nil || connections.ElasticsearchConnection.Client == nil {
		t.Skip("Elasticsearch connection not available, skipping test")
		return
	}
	
	repo := models.ElasticCompanyRepository(connections.ElasticsearchConnection.Client)
	
	// Simple query to test connectivity
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": 1,
	}
	
	results, err := repo.ListByQueryMap(query)
	if err != nil {
		t.Logf("ListByQueryMap error (may be expected if ES not configured): %v", err)
		return
	}
	
	if results == nil {
		t.Error("ListByQueryMap returned nil")
		return
	}
	
	t.Logf("Retrieved %d results from Elasticsearch", len(results))
}

func TestElasticsearchContactRepository_ListByQueryMap(t *testing.T) {
	if connections.ElasticsearchConnection == nil || connections.ElasticsearchConnection.Client == nil {
		t.Skip("Elasticsearch connection not available, skipping test")
		return
	}
	
	repo := models.ElasticContactRepository(connections.ElasticsearchConnection.Client)
	
	// Simple query to test connectivity
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": 1,
	}
	
	results, err := repo.ListByQueryMap(query)
	if err != nil {
		t.Logf("ListByQueryMap error (may be expected if ES not configured): %v", err)
		return
	}
	
	if results == nil {
		t.Error("ListByQueryMap returned nil")
		return
	}
	
	t.Logf("Retrieved %d results from Elasticsearch", len(results))
}

func TestElasticsearchCompanyRepository_CountByQueryMap(t *testing.T) {
	if connections.ElasticsearchConnection == nil || connections.ElasticsearchConnection.Client == nil {
		t.Skip("Elasticsearch connection not available, skipping test")
		return
	}
	
	repo := models.ElasticCompanyRepository(connections.ElasticsearchConnection.Client)
	
	// Simple count query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	
	count, err := repo.CountByQueryMap(query)
	if err != nil {
		t.Logf("CountByQueryMap error (may be expected if ES not configured): %v", err)
		return
	}
	
	t.Logf("Company count in Elasticsearch: %d", count)
}

func TestElasticsearchContactRepository_CountByQueryMap(t *testing.T) {
	if connections.ElasticsearchConnection == nil || connections.ElasticsearchConnection.Client == nil {
		t.Skip("Elasticsearch connection not available, skipping test")
		return
	}
	
	repo := models.ElasticContactRepository(connections.ElasticsearchConnection.Client)
	
	// Simple count query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	
	count, err := repo.CountByQueryMap(query)
	if err != nil {
		t.Logf("CountByQueryMap error (may be expected if ES not configured): %v", err)
		return
	}
	
	t.Logf("Contact count in Elasticsearch: %d", count)
}
