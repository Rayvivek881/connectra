package constants

var (
	InsertFileJobType = "insert_file"
	FileTypeCsv       = "csv"

	OpenJobStatus       = "open"
	InQueueJobStatus    = "in_queue"
	ProcessingJobStatus = "processing"
	CompletedJobStatus  = "completed"
	FailedJobStatus     = "failed"

	RetryInQueuedJobStatus = "retry_in_queued"
	RetryingJobStatus      = "retrying"
)
