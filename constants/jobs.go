package constants

var (
	OpenJobStatus       = "open"
	InQueueJobStatus    = "in_queue"
	ProcessingJobStatus = "processing"
	CompletedJobStatus  = "completed"
	FailedJobStatus     = "failed"

	RetryInQueuedJobStatus = "retry_in_queued"
	RetryingJobStatus      = "retrying"

	FirstTimeJobType = "first_time"
	RetryJobType     = "retry"
	InsertCsvFile    = "insert_csv_file"
	ExportCsvFile    = "export_csv_file"
)
