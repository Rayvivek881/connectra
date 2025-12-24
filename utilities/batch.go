package utilities

import (
	"fmt"
	"sync/atomic"
	"time"
)

// BatchProcessor handles batch processing with statistics
type BatchProcessor struct {
	BatchSize int
	Stats     *BatchStats
}

// BatchStats tracks statistics for batch processing operations
type BatchStats struct {
	TotalProcessed   int64
	TotalSuccess     int64
	TotalFailed      int64
	BatchesProcessed int64
	StartTime        time.Time
	EndTime          time.Time
}

// NewBatchProcessor creates a new batch processor with the specified batch size
func NewBatchProcessor(batchSize int) *BatchProcessor {
	if batchSize <= 0 {
		batchSize = 1000 // Default batch size
	}
	return &BatchProcessor{
		BatchSize: batchSize,
		Stats: &BatchStats{
			StartTime: time.Now(),
		},
	}
}

// ProcessInBatches processes items in batches using the provided processor function
// T is a generic type parameter, allowing this function to work with any slice type
func ProcessInBatches[T any](
	batchSize int,
	items []T,
	processor func(batch []T) (success int, failed int, err error),
) (*BatchStats, error) {
	if batchSize <= 0 {
		batchSize = 1000
	}
	if len(items) == 0 {
		return &BatchStats{}, nil
	}

	stats := &BatchStats{
		StartTime: time.Now(),
	}

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		success, failed, err := processor(batch)

		atomic.AddInt64(&stats.TotalProcessed, int64(len(batch)))
		atomic.AddInt64(&stats.TotalSuccess, int64(success))
		atomic.AddInt64(&stats.TotalFailed, int64(failed))
		atomic.AddInt64(&stats.BatchesProcessed, 1)

		if err != nil {
			stats.EndTime = time.Now()
			return stats, fmt.Errorf("batch %d failed: %w", stats.BatchesProcessed, err)
		}
	}

	stats.EndTime = time.Now()
	return stats, nil
}

// GetProcessingTime returns the duration of the batch processing
func (bs *BatchStats) GetProcessingTime() time.Duration {
	if bs.EndTime.IsZero() {
		return time.Since(bs.StartTime)
	}
	return bs.EndTime.Sub(bs.StartTime)
}
