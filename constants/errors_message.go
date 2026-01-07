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

	CompanyNotFoundError = errors.New("ERR_COMPANY_NOT_FOUND: the specified company UUID does not exist")
	ContactNotFoundError = errors.New("ERR_CONTACT_NOT_FOUND: the specified contact UUID does not exist")
)

func InvalidJobTypeError(jobType string) error {
	return fmt.Errorf("ERR_INVALID_JOB_TYPE: job type '%s' is not recognized; supported types are 'insert_csv_file' and 'export_csv_file'", jobType)
}

func ElasticsearchError(statusCode int, body string) error {
	return fmt.Errorf("ERR_ELASTICSEARCH_FAILURE: search engine returned status %d; details: %s", statusCode, body)
}

func ElasticsearchBulkError(statusCode int, body string) error {
	return fmt.Errorf("ERR_ELASTICSEARCH_BULK_FAILURE: bulk indexing operation returned status %d; details: %s", statusCode, body)
}
