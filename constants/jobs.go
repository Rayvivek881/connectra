package constants

var (
	InsertFileJob = "insert_file"

	OpenJobStatus       = "open"
	InQueueJobStatus    = "in_queue"
	ProcessingJobStatus = "processing"
	CompletedJobStatus  = "completed"
	FailedJobStatus     = "failed"

	RetryInQueuedJobStatus = "retry_in_queued"
	RetryingJobStatus      = "retrying"
)
