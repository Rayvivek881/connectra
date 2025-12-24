package utilities

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ProducerConsumerConfig holds configuration for producer-consumer pattern
type ProducerConsumerConfig struct {
	NumProducers      int           // Number of producer goroutines
	NumConsumers      int           // Number of consumer goroutines
	BatchSize         int           // Size of each batch
	ChannelBufferSize int           // Buffer size for channels
	ProgressInterval  time.Duration // Interval for progress reporting (0 to disable)
}

// DefaultProducerConsumerConfig returns default configuration
func DefaultProducerConsumerConfig() ProducerConsumerConfig {
	return ProducerConsumerConfig{
		NumProducers:      2,
		NumConsumers:      4,
		BatchSize:         1000,
		ChannelBufferSize: 10000,
		ProgressInterval:  10 * time.Second,
	}
}

// ProducerConsumerStats tracks statistics for producer-consumer operations
type ProducerConsumerStats struct {
	TotalProcessed   int64
	TotalSuccess     int64
	TotalFailed      int64
	BatchesProcessed int64
	StartTime        time.Time
	EndTime          time.Time
	mu               sync.RWMutex
	errors           []string
}

// AddError adds an error to the statistics
func (pcs *ProducerConsumerStats) AddError(err string) {
	pcs.mu.Lock()
	defer pcs.mu.Unlock()
	pcs.errors = append(pcs.errors, err)
}

// GetErrors returns all collected errors
func (pcs *ProducerConsumerStats) GetErrors() []string {
	pcs.mu.RLock()
	defer pcs.mu.RUnlock()
	return pcs.errors
}

// GetProcessingTime returns the duration of the processing
func (pcs *ProducerConsumerStats) GetProcessingTime() time.Duration {
	if pcs.EndTime.IsZero() {
		return time.Since(pcs.StartTime)
	}
	return pcs.EndTime.Sub(pcs.StartTime)
}

// SimpleProcessWithProducerConsumer processes items using producer-consumer pattern
// This is useful for very large datasets (millions of records)
// T is the item type
func SimpleProcessWithProducerConsumer[T any](
	ctx context.Context,
	config ProducerConsumerConfig,
	items []T,
	consumer func(ctx context.Context, batch []T, batchNum int) (success int, failed int, err error),
	progressReporter func(stats *ProducerConsumerStats),
) (*ProducerConsumerStats, error) {
	if len(items) == 0 {
		return &ProducerConsumerStats{}, nil
	}

	if config.BatchSize <= 0 {
		config.BatchSize = 1000
	}
	if config.NumProducers <= 0 {
		config.NumProducers = 2
	}
	if config.NumConsumers <= 0 {
		config.NumConsumers = 4
	}
	if config.ChannelBufferSize <= 0 {
		config.ChannelBufferSize = 10000
	}

	stats := &ProducerConsumerStats{
		StartTime: time.Now(),
	}

	// Calculate total batches
	totalBatches := (len(items) + config.BatchSize - 1) / config.BatchSize

	// Create channel for batches
	type batchData struct {
		batch    []T
		batchNum int
	}
	batchChan := make(chan batchData, config.ChannelBufferSize)

	// Wait groups
	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	// Start progress reporter if enabled
	var progressTicker *time.Ticker
	var progressDone chan struct{}
	if config.ProgressInterval > 0 && progressReporter != nil {
		progressTicker = time.NewTicker(config.ProgressInterval)
		progressDone = make(chan struct{})
		go func() {
			for {
				select {
				case <-progressTicker.C:
					progressReporter(stats)
				case <-progressDone:
					return
				}
			}
		}()
	}

	// Start producers
	batchesPerProducer := totalBatches / config.NumProducers
	for i := 0; i < config.NumProducers; i++ {
		producerWg.Add(1)
		startBatch := i * batchesPerProducer
		endBatch := startBatch + batchesPerProducer
		if i == config.NumProducers-1 {
			endBatch = totalBatches // Last producer handles remaining batches
		}

		go func(start, end int, producerID int) {
			defer producerWg.Done()
			for batchNum := start; batchNum < end; batchNum++ {
				startIdx := batchNum * config.BatchSize
				endIdx := startIdx + config.BatchSize
				if endIdx > len(items) {
					endIdx = len(items)
				}

				// Create a copy of the batch slice to avoid race conditions
				batch := make([]T, endIdx-startIdx)
				copy(batch, items[startIdx:endIdx])

				select {
				case batchChan <- batchData{batch: batch, batchNum: batchNum}:
				case <-ctx.Done():
					return
				}
			}
		}(startBatch, endBatch, i+1)
	}

	// Close channel when all producers are done
	go func() {
		producerWg.Wait()
		close(batchChan)
	}()

	// Start consumers
	for i := 0; i < config.NumConsumers; i++ {
		consumerWg.Add(1)
		go func(consumerID int) {
			defer consumerWg.Done()
			for {
				select {
				case batchData, ok := <-batchChan:
					if !ok {
						return // Channel closed, no more batches
					}

					success, failed, err := consumer(ctx, batchData.batch, batchData.batchNum)
					atomic.AddInt64(&stats.TotalProcessed, int64(success+failed))
					atomic.AddInt64(&stats.TotalSuccess, int64(success))
					atomic.AddInt64(&stats.TotalFailed, int64(failed))
					atomic.AddInt64(&stats.BatchesProcessed, 1)

					if err != nil {
						stats.AddError(fmt.Sprintf("consumer %d: batch %d failed: %v", consumerID, batchData.batchNum, err))
					}
				case <-ctx.Done():
					return
				}
			}
		}(i + 1)
	}

	// Wait for all consumers to finish
	consumerWg.Wait()

	// Stop progress reporter
	if progressTicker != nil {
		progressTicker.Stop()
		close(progressDone)
	}

	stats.EndTime = time.Now()

	return stats, nil
}

