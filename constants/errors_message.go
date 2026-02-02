package constants

import (
	"errors"
	"fmt"
)

var (
	PageSizeExceededError   = errors.New("ERR_PAGE_SIZE_EXCEEDED: the requested page size surpasses the maximum allowed limit; consider using pagination with smaller batches")
	PageNumberExceededError = errors.New("ERR_PAGE_OUT_OF_RANGE: the requested page number is beyond the available range; verify total pages before requesting")

	FailedToFetchDataError = errors.New("ERR_DATA_FETCH_FAILED: an unexpected error occurred while retrieving records from the data store; please retry or contact support if the issue persists")

	DataArrayEmptyError        = errors.New("ERR_EMPTY_PAYLOAD: the 'data' array in the request body is empty; provide at least one record to process")
	JobTypeRequiredError       = errors.New("ERR_MISSING_JOB_TYPE: the 'job_type' field is required; specify a valid job type such as 'insert_csv_file' or 'export_csv_file'")
	JobDataRequiredError       = errors.New("ERR_MISSING_JOB_DATA: the 'job_data' field is required; include the necessary payload for job execution")
	RetryCountNegativeError    = errors.New("ERR_INVALID_RETRY_COUNT: 'retry_count' must be a non-negative integer; use 0 for no retries or a positive number for retry attempts")
	LimitNegativeError         = errors.New("ERR_INVALID_LIMIT: 'limit' must be a non-negative integer; use 0 for default or specify a positive value")
	LimitExceededError         = errors.New("ERR_LIMIT_TOO_HIGH: 'limit' exceeds the maximum of 100 records per request; reduce the value or use pagination")
	BatchSizeExceededError     = errors.New("ERR_BATCH_TOO_LARGE: the number of records in the batch exceeds the allowed maximum; split the data into smaller chunks")
	SelectColumnsRequiredError = errors.New("ERR_MISSING_SELECT_COLUMNS: 'select_columns' is required for export operations; specify at least one column to include in the output")

	InvalidServiceError     = errors.New("ERR_UNKNOWN_SERVICE: the provided service identifier is not recognized; use 'contacts' or 'companies'")
	InvalidServiceTypeError = errors.New("ERR_UNSUPPORTED_SERVICE: the specified service type is not supported for this operation; verify the endpoint and try again")

	JobNotFoundError     = errors.New("ERR_JOB_NOT_FOUND: the requested job does not exist; verify the job UUID and try again")
	JobUuidRequiredError = errors.New("ERR_JOB_UUID_REQUIRED: 'job_uuid' path parameter is required; provide a valid job UUID")

	FilenameRequiredError            = errors.New("ERR_FILENAME_REQUIRED: 'filename' query parameter is required; provide a valid filename")
	S3KeyRequiredError               = errors.New("ERR_S3_KEY_REQUIRED: 's3_key' query parameter is required; provide a valid S3 key")
	FailedToGenerateUploadURLError   = errors.New("ERR_UPLOAD_URL_GENERATION_FAILED: failed to generate presigned upload URL; please retry or contact support")
	FailedToGenerateDownloadURLError = errors.New("ERR_DOWNLOAD_URL_GENERATION_FAILED: failed to generate presigned download URL; please retry or contact support")

	FailedToInitBatchServiceError = errors.New("ERR_BATCH_SERVICE_INIT_FAILED: failed to initialize batch service; please retry or contact support")

	RateLimitExceededError = errors.New("ERR_RATE_LIMIT_EXCEEDED: too many requests; please try again later")
	UnauthorizedError      = errors.New("ERR_UNAUTHORIZED: invalid or missing API key; provide a valid 'X-API-Key' header")

	OpenSearchNotConnectedError = errors.New("ERR_OPENSEARCH_NOT_CONNECTED: search engine is not connected; check OpenSearch is running and configuration")
)

func InvalidJobTypeError(jobType string) error {
	return fmt.Errorf("ERR_INVALID_JOB_TYPE: job type '%s' is not recognized; supported types are 'insert_csv_file' and 'export_csv_file'", jobType)
}

func InvalidUUIDError(uuid string) error {
	return fmt.Errorf("ERR_INVALID_UUID: '%s' is not a valid UUID format; use a valid UUID like 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'", uuid)
}

func OpenSearchError(statusCode int, body string) error {
	return fmt.Errorf("ERR_OPENSEARCH_FAILURE: search engine returned status %d; details: %s", statusCode, body)
}

func OpenSearchBulkError(statusCode int, body string) error {
	return fmt.Errorf("ERR_OPENSEARCH_BULK_FAILURE: bulk indexing operation returned status %d; details: %s", statusCode, body)
}
