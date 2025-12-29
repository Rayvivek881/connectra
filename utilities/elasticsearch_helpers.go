package utilities

import (
	"fmt"
	"sync"
)

// AsyncIndexErrorHandler is called when async indexing fails.
// If nil, errors are silently ignored (existing behavior).
type AsyncIndexErrorHandler func(err error, entityType string, identifier string)

// DefaultAsyncIndexErrorHandler logs errors using fmt.Printf (existing behavior).
func DefaultAsyncIndexErrorHandler(err error, entityType string, identifier string) {
	fmt.Printf("Warning: Failed to index %s in Elasticsearch (id: %s): %v\n", entityType, identifier, err)
}

var (
	// Global error handler for async indexing
	asyncIndexErrorHandler AsyncIndexErrorHandler = DefaultAsyncIndexErrorHandler
	errorHandlerMutex      sync.RWMutex
)

// SetAsyncIndexErrorHandler sets a custom error handler for async indexing operations.
func SetAsyncIndexErrorHandler(handler AsyncIndexErrorHandler) {
	errorHandlerMutex.Lock()
	defer errorHandlerMutex.Unlock()
	asyncIndexErrorHandler = handler
}

// GetAsyncIndexErrorHandler returns the current async index error handler.
func GetAsyncIndexErrorHandler() AsyncIndexErrorHandler {
	errorHandlerMutex.RLock()
	defer errorHandlerMutex.RUnlock()
	return asyncIndexErrorHandler
}

// AsyncIndex executes an indexing function asynchronously with error handling.
// This is a helper to standardize async Elasticsearch indexing patterns.
//
// Example:
//   AsyncIndex("contact", contact.UUID, func() error {
//       return elasticRepo.IndexContact(elasticContact)
//   })
func AsyncIndex(entityType string, identifier string, indexFn func() error) {
	go func() {
		if err := indexFn(); err != nil {
			handler := GetAsyncIndexErrorHandler()
			if handler != nil {
				handler(err, entityType, identifier)
			}
		}
	}()
}

// BulkIndexResult represents the result of a bulk indexing operation.
type BulkIndexResult struct {
	Indexed int64
	Failed  int64
	Errors  []string
}

// BulkIndexWithResult processes items in batches and indexes them, collecting results.
// This is a helper pattern for bulk Elasticsearch indexing operations.
//
// Example:
//   result := BulkIndexWithResult(items, batchSize, func(batch []MyType) (int64, error) {
//       return elasticRepo.BulkUpsert(batch)
//   })
func BulkIndexWithResult[T any](
	items []T,
	batchSize int,
	indexFn func([]T) (int64, error),
) BulkIndexResult {
	if len(items) == 0 {
		return BulkIndexResult{}
	}

	result := BulkIndexResult{
		Errors: make([]string, 0),
	}

	// Process in batches
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		indexed, err := indexFn(batch)
		if err != nil {
			result.Failed += int64(len(batch))
			result.Errors = append(result.Errors, fmt.Sprintf("batch %d-%d failed: %v", i, end-1, err))
		} else {
			result.Indexed += indexed
		}
	}

	return result
}
