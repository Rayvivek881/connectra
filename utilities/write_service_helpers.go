package utilities

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/uptrace/bun"
)

// BulkUpsertResult represents the result of a bulk upsert operation
type BulkUpsertResult struct {
	Created              int64
	Updated              int64
	Total                int64
	BatchesProcessed     int64
	ElasticsearchIndexed int64
	ElasticsearchFailed  int64
	ProcessingTime       time.Duration
	Errors               []string
}

// TimestampSetter is a function that sets timestamps on an entity
type TimestampSetter[T any] func(entity *T, now time.Time, isNew bool)

// EntityUUIDGetter is a function that extracts UUID from an entity
type EntityUUIDGetter[T any] func(entity T) string

// EntityToSlicePtr converts a slice of entities to a slice of pointers
type EntityToSlicePtr[T any] func(entities []T) []*T

// BulkUpsertOptions contains configuration for bulk upsert operations
type BulkUpsertOptions[T any] struct {
	BatchSize        int
	TimestampSetter  TimestampSetter[T]
	UUIDGetter       EntityUUIDGetter[T]
	EntityToPtrSlice EntityToSlicePtr[T]
	EntityTypeName   string // e.g., "contact", "company"
}

// ExecuteBulkUpsertInTransaction executes a bulk upsert operation within a transaction.
// This helper extracts the common transaction pattern for bulk operations.
func ExecuteBulkUpsertInTransaction[T any](
	ctx context.Context,
	db *bun.DB,
	entities []*T,
	upsertFn func(tx bun.Tx, entities []*T) error,
) error {
	return WithTransactionVoid(ctx, db, func(tx bun.Tx) error {
		return upsertFn(tx, entities)
	})
}

// CalculateCreatedUpdatedCount calculates how many entities are new vs updated
// based on a map of existing UUIDs.
func CalculateCreatedUpdatedCount[T any](
	entities []T,
	existingMap map[string]bool,
	uuidGetter EntityUUIDGetter[T],
) (created int64, updated int64) {
	created = 0
	for _, entity := range entities {
		uuid := uuidGetter(entity)
		if !existingMap[uuid] {
			created++
		}
	}
	updated = int64(len(entities)) - created
	return created, updated
}

// BulkUpsertStats tracks statistics for bulk upsert operations
type BulkUpsertStats struct {
	BatchesProcessed     int64
	ElasticsearchIndexed int64
	ElasticsearchFailed  int64
	Errors               []string
	mu                   sync.Mutex
}

// NewBulkUpsertStats creates a new BulkUpsertStats instance
func NewBulkUpsertStats() *BulkUpsertStats {
	return &BulkUpsertStats{
		Errors: make([]string, 0),
	}
}

// RecordBatch records that a batch was processed
func (s *BulkUpsertStats) RecordBatch() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BatchesProcessed++
}

// RecordElasticsearchSuccess records successful Elasticsearch indexing
func (s *BulkUpsertStats) RecordElasticsearchSuccess(indexed int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ElasticsearchIndexed += indexed
}

// RecordElasticsearchFailure records failed Elasticsearch indexing
func (s *BulkUpsertStats) RecordElasticsearchFailure(count int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ElasticsearchFailed += count
	s.Errors = append(s.Errors, fmt.Sprintf("batch failed to index in Elasticsearch: %v", err))
}

// RecordError records a general error
func (s *BulkUpsertStats) RecordError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Errors = append(s.Errors, err.Error())
}

// SetTimestamps applies timestamps to a slice of entities
func SetTimestamps[T any](
	entities []T,
	now time.Time,
	setter TimestampSetter[T],
	existingMap map[string]bool,
	uuidGetter EntityUUIDGetter[T],
) {
	for i := range entities {
		uuid := uuidGetter(entities[i])
		isNew := !existingMap[uuid]
		setter(&entities[i], now, isNew)
	}
}

