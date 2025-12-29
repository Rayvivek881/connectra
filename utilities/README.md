# Connectra Utilities

This directory contains reusable utility functions and helpers for common patterns across the Connectra (Go) backend.

## Transaction Helpers

**File**: `transaction_helpers.go`

Generic transaction wrappers with automatic rollback handling.

### WithTransaction[T]
Execute a function within a database transaction with automatic rollback on error.

```go
result, err := utilities.WithTransaction(ctx, db, func(tx bun.Tx) (*MyModel, error) {
    // Perform operations with tx
    return &MyModel{}, nil
})
```

### WithTransactionVoid
Transaction wrapper for functions that return no value.

```go
err := utilities.WithTransactionVoid(ctx, db, func(tx bun.Tx) error {
    // Perform operations with tx
    return nil
})
```

## Elasticsearch Helpers

**File**: `elasticsearch_helpers.go`

Standardized async Elasticsearch indexing patterns with error handling.

### AsyncIndex
Execute indexing function asynchronously with standardized error handling.

```go
utilities.AsyncIndex("contact", contact.UUID, func() error {
    return elasticRepo.IndexContact(elasticContact)
})
```

### BulkIndexWithResult[T]
Process bulk indexing operations with result tracking.

```go
result := utilities.BulkIndexWithResult(contacts, 100, func(batch []Contact) error {
    // Index batch
    return nil
})
```

## Write Service Helpers

**File**: `write_service_helpers.go`

Common patterns for write service operations including bulk upserts.

### Key Types
- `BulkUpsertResult` - Result structure for bulk operations
- `BulkUpsertStats` - Statistics tracking for bulk operations

### Helper Functions
- `ExecuteBulkUpsertInTransaction[T]` - Execute bulk upsert within transaction
- `CalculateCreatedUpdatedCount[T]` - Calculate new vs updated entity counts
- `SetTimestamps[T]` - Apply timestamps to entities

## Usage Guidelines

1. **Use transaction helpers** for all database operations requiring transactions
2. **Use Elasticsearch helpers** for async indexing operations
3. **Use write service helpers** for bulk operations and statistics tracking
4. **Follow error handling patterns** established in the helpers
5. **Add new helpers** when patterns are repeated across multiple services

## Migration Guide

When migrating existing code to use these helpers:

1. Replace manual transaction handling with `WithTransaction` or `WithTransactionVoid`
2. Replace async indexing goroutines with `AsyncIndex`
3. Use bulk operation helpers for batch processing
4. Update error handling to use helper error handlers

