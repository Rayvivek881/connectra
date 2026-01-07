# Integration Tests

This directory contains integration tests for the Connectra API.

## Prerequisites

Integration tests require:

- PostgreSQL database (can use test database)
- Elasticsearch instance (can use test instance or mock)
- Environment variables configured (see `.env.example`)

## Running Integration Tests

### Run All Integration Tests

```bash
go test -tags=integration ./tests -v
```

### Run Specific Test File

```bash
go test -tags=integration ./tests -v -run TestHealthCheck
```

### Run with Coverage

```bash
go test -tags=integration ./tests -v -cover
```

## Test Structure

- `integration_test.go` - Test setup, teardown, and helper functions
- `api_endpoints_test.go` - API endpoint integration tests
- `database_test.go` - Database operation tests (if needed)
- `elasticsearch_test.go` - Elasticsearch sync tests (if needed)

## Test Environment

Integration tests use the `integration` build tag to separate them from unit tests. This allows:
- Unit tests to run quickly without external dependencies
- Integration tests to run only when explicitly requested
- CI/CD pipelines to run both separately

## Configuration

Set the following environment variables for integration tests:

```bash
# API Configuration
API_KEY=test-api-key-for-integration-tests
MAX_REQUESTS_PER_MINUTE=1000

# Database (use test database)
PG_DB_HOST=localhost
PG_DB_PORT=5432
PG_DB_DATABASE=connectra_test
PG_DB_USERNAME=test_user
PG_DB_PASSWORD=test_password

# Elasticsearch (use test instance)
ELASTICSEARCH_HOST=localhost
ELASTICSEARCH_PORT=9200
```


## Notes

- Integration tests may require actual database and Elasticsearch connections
- Tests are designed to be idempotent (can run multiple times)
- Tests clean up after themselves when possible
- Some tests may be skipped if dependencies are not available
