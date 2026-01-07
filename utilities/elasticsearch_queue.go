package utilities

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// IndexOperation represents a single Elasticsearch indexing operation
type IndexOperation struct {
	Type      string      // "create", "update", "delete"
	Index     string      // Elasticsearch index name
	DocumentID string     // Document UUID
	Document  interface{} // Document data (nil for delete operations)
	Retries   int         // Current retry count
	MaxRetries int        // Maximum retry attempts
}

// ElasticsearchQueue manages async indexing operations with retry logic
type ElasticsearchQueue struct {
	queue      chan *IndexOperation
	workers    int
	maxRetries int
	processor  func(*IndexOperation) error
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewElasticsearchQueue creates a new async indexing queue
func NewElasticsearchQueue(workers int, queueSize int, maxRetries int, processor func(*IndexOperation) error) *ElasticsearchQueue {
	ctx, cancel := context.WithCancel(context.Background())
	
	queue := &ElasticsearchQueue{
		queue:      make(chan *IndexOperation, queueSize),
		workers:    workers,
		maxRetries: maxRetries,
		processor:  processor,
		ctx:        ctx,
		cancel:     cancel,
	}
	
	// Start worker goroutines
	for i := 0; i < workers; i++ {
		queue.wg.Add(1)
		go queue.worker(i)
	}
	
	return queue
}

// Enqueue adds an operation to the queue
func (q *ElasticsearchQueue) Enqueue(op *IndexOperation) bool {
	op.MaxRetries = q.maxRetries
	select {
	case q.queue <- op:
		return true
	case <-q.ctx.Done():
		return false
	default:
		// Queue is full, log warning but don't block
		log.Warn().
			Str("type", op.Type).
			Str("index", op.Index).
			Str("document_id", op.DocumentID).
			Msg("Elasticsearch queue is full, operation dropped")
		return false
	}
}

// worker processes operations from the queue
func (q *ElasticsearchQueue) worker(id int) {
	defer q.wg.Done()
	
	for {
		select {
		case op := <-q.queue:
			q.processOperation(op)
		case <-q.ctx.Done():
			return
		}
	}
}

// processOperation processes a single operation with retry logic
func (q *ElasticsearchQueue) processOperation(op *IndexOperation) {
	for op.Retries <= op.MaxRetries {
		err := q.processor(op)
		if err == nil {
			// Success
			return
		}
		
		// Check if we should retry
		if op.Retries >= op.MaxRetries {
			log.Error().
				Err(err).
				Str("type", op.Type).
				Str("index", op.Index).
				Str("document_id", op.DocumentID).
				Int("retries", op.Retries).
				Msg("Elasticsearch operation failed after max retries")
			return
		}
		
		// Exponential backoff: 1s, 2s, 4s, 8s, etc.
		backoff := time.Duration(1<<uint(op.Retries)) * time.Second
		if backoff > 30*time.Second {
			backoff = 30 * time.Second // Cap at 30 seconds
		}
		
		log.Warn().
			Err(err).
			Str("type", op.Type).
			Str("index", op.Index).
			Str("document_id", op.DocumentID).
			Int("retry", op.Retries+1).
			Dur("backoff", backoff).
			Msg("Elasticsearch operation failed, retrying")
		
		time.Sleep(backoff)
		op.Retries++
	}
}

// Stop gracefully stops the queue
func (q *ElasticsearchQueue) Stop() {
	close(q.queue)
	q.cancel()
	q.wg.Wait()
}

// Size returns the current queue size
func (q *ElasticsearchQueue) Size() int {
	return len(q.queue)
}

// Global queue instances (initialized on first use)
var (
	companyIndexQueue  *ElasticsearchQueue
	contactIndexQueue  *ElasticsearchQueue
	queueOnce          sync.Once
)

// InitializeQueues initializes the global Elasticsearch indexing queues
func InitializeQueues(
	companyProcessor func(*IndexOperation) error,
	contactProcessor func(*IndexOperation) error,
) {
	queueOnce.Do(func() {
		// Queue size: 1000 operations, 5 workers, max 3 retries
		companyIndexQueue = NewElasticsearchQueue(5, 1000, 3, companyProcessor)
		contactIndexQueue = NewElasticsearchQueue(5, 1000, 3, contactProcessor)
		
		log.Info().
			Int("workers", 5).
			Int("queue_size", 1000).
			Int("max_retries", 3).
			Msg("Elasticsearch indexing queues initialized")
	})
}

// GetCompanyQueue returns the company indexing queue
func GetCompanyQueue() *ElasticsearchQueue {
	return companyIndexQueue
}

// GetContactQueue returns the contact indexing queue
func GetContactQueue() *ElasticsearchQueue {
	return contactIndexQueue
}

// EnqueueCompanyOperation enqueues a company indexing operation
func EnqueueCompanyOperation(opType, documentID string, document interface{}) {
	if companyIndexQueue == nil {
		// Queue not initialized, skip
		return
	}
	
	companyIndexQueue.Enqueue(&IndexOperation{
		Type:       opType,
		Index:      "companies", // Will be resolved from constants
		DocumentID: documentID,
		Document:   document,
		Retries:    0,
	})
}

// EnqueueContactOperation enqueues a contact indexing operation
func EnqueueContactOperation(opType, documentID string, document interface{}) {
	if contactIndexQueue == nil {
		// Queue not initialized, skip
		return
	}
	
	contactIndexQueue.Enqueue(&IndexOperation{
		Type:       opType,
		Index:      "contacts", // Will be resolved from constants
		DocumentID: documentID,
		Document:   document,
		Retries:    0,
	})
}

// SerializeDocument converts a document to JSON bytes for indexing
func SerializeDocument(doc interface{}) ([]byte, error) {
	return json.Marshal(doc)
}
