# Testing Summary

This document summarizes the testing infrastructure and coverage for the Connectra API.

## Test Structure

### Unit Tests (`utilities/`, `modules/common/helper/`)

Unit tests are located alongside the code they test and can be run without external dependencies:

- **`utilities/common_test.go`** - Validation helper tests
  - Email validation
  - URL validation (HTTP/HTTPS, LinkedIn)
  - UUID validation
  - Required field validation
  - Non-negative number validation
  - String length validation
  - Phone number validation and cleaning

- **`utilities/cache_test.go`** - TTL cache tests
  - Get/Set operations
  - TTL expiration
  - Delete and clear operations
  - Prefix-based invalidation
  - Max size and eviction
  - Custom TTL support

- **`utilities/query_test.go`** - VQL query compilation tests
  - Empty query detection
  - Text match queries
  - Keyword match queries
  - Range queries
  - Combined queries
  - Count queries

- **`modules/common/helper/errors_test.go`** - Error response helper tests
  - Standardized error response format
  - Validation error responses
  - Not found error responses
  - Bad request error responses
  - Unauthorized error responses
  - Rate limit error responses

### Integration Tests (`tests/`)

Integration tests require database and Elasticsearch connections and use the `integration` build tag:

- **`tests/integration_test.go`** - Test infrastructure
  - Test server setup
  - Database and Elasticsearch initialization
  - Helper functions for making API requests

- **`tests/api_endpoints_test.go`** - API endpoint tests
  - Health check endpoint
  - Authentication validation
  - Company CRUD operations
  - Contact CRUD operations
  - Input validation (email, UUID)
  - Error response format validation
  - Filtering operations

- **`tests/database_test.go`** - Database operation tests
  - Database connectivity
  - Filter repository operations
  - Service-specific filter retrieval

- **`tests/elasticsearch_test.go`** - Elasticsearch sync tests
  - Elasticsearch connectivity
  - Company repository queries
  - Contact repository queries
  - Count operations

## Running Tests

### Unit Tests

```bash
# Run all unit tests
go test ./utilities ./modules/common/helper -v

# Run with coverage
go test ./utilities ./modules/common/helper -v -cover

# Run specific test
go test ./utilities -v -run TestValidateEmail
```

### Integration Tests

```bash
# Run all integration tests
go test -tags=integration ./tests -v

# Run specific test file
go test -tags=integration ./tests -v -run TestHealthCheck

# Run with coverage
go test -tags=integration ./tests -v -cover
```

## Test Coverage

### Unit Test Coverage

- ✅ **Validation Helpers**: 100% coverage of validation functions
- ✅ **Cache Operations**: All cache operations tested
- ✅ **Query Compilation**: VQL to Elasticsearch query conversion
- ✅ **Error Handling**: Standardized error response format

### Integration Test Coverage

- ✅ **API Endpoints**: All CRUD endpoints for companies and contacts
- ✅ **Authentication**: API key validation
- ✅ **Input Validation**: Email, UUID, required fields
- ✅ **Database Operations**: Filter repository operations
- ✅ **Elasticsearch Operations**: Query and count operations
- ✅ **Error Handling**: Standardized error response format

## Test Features

### Build Tag Separation

Integration tests use the `//go:build integration` and `// +build integration` build tags to:
- Separate unit tests (fast, no dependencies) from integration tests (slower, require DB/ES)
- Allow CI/CD pipelines to run tests separately
- Enable developers to run quick unit tests during development

### Graceful Skipping

Integration tests gracefully skip when dependencies are unavailable:
- Database connection tests skip if PostgreSQL is not available
- Elasticsearch tests skip if Elasticsearch is not available
- Tests log warnings instead of failing when dependencies are missing

### Test Helpers

The `makeRequest()` helper function simplifies API testing:
- Automatically sets Content-Type header
- Adds API key authentication
- Handles JSON marshaling
- Returns response recorder for assertions

## Best Practices

1. **Unit tests** should be fast and not require external dependencies
2. **Integration tests** should be idempotent (can run multiple times)
3. **Test data** should be cleaned up after tests when possible
4. **Error cases** should be tested alongside success cases
5. **Build tags** should be used to separate test types

## Future Improvements

Potential areas for test expansion:

- [ ] Mock database for faster integration tests
- [ ] Mock Elasticsearch for faster integration tests
- [ ] End-to-end tests for complete workflows
- [ ] Performance/load tests
- [ ] Security tests (SQL injection, XSS, etc.)
- [ ] Contract tests for API compatibility
