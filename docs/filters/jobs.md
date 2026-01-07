# Jobs API - Complete Guide

**Version**: 1.0  
**Last Updated**: 2025-01-XX

## Table of Contents

1. [Overview](#overview)
2. [Job System Architecture](#job-system-architecture)
3. [Job Types](#job-types)
4. [Job States and Lifecycle](#job-states-and-lifecycle)
5. [API Reference](#api-reference)
6. [Job Runners](#job-runners)
7. [Configuration](#configuration)
8. [Examples](#examples)
9. [Error Handling](#error-handling)
10. [Best Practices](#best-practices)
11. [Related Documentation](#related-documentation)

---

## Overview

The Jobs API provides a distributed, asynchronous job processing system for handling large-scale data operations. Jobs are processed in the background using a worker pool pattern with support for retries, backpressure management, and graceful shutdown.

### Key Features

- **Asynchronous Processing**: Jobs are queued and processed in the background
- **Distributed Processing**: Multiple workers can process jobs concurrently
- **Automatic Retries**: Failed jobs can be automatically retried with configurable intervals
- **Memory Efficient**: Streaming CSV processing for large files (multi-GB support)
- **State Management**: Complete job lifecycle tracking with status updates
- **Backpressure Control**: Prevents system overload with configurable queue limits

### Use Cases

- **CSV Import**: Import large CSV files from S3 into PostgreSQL and Elasticsearch
- **CSV Export**: Export filtered contact/company data to S3 as CSV files
- **Batch Operations**: Process large datasets without blocking API requests
- **Data Migration**: Move data between systems efficiently

---

## Job System Architecture

### High-Level Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  API REQUEST: POST /common/jobs/create                                       │
│  Creates job with status: "open"                                            │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  JOB RUNNER (first_time or retry)                                           │
│                                                                             │
│  Ticker (poll every N minutes)                                              │
│    ↓                                                                         │
│  Fetch jobs with status: "open" or "failed"                                │
│    ↓                                                                         │
│  Update status: "in_queue" or "retry_in_queued"                             │
│    ↓                                                                         │
│  Push to buffered channel (capacity: 1000)                                  │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  WORKER POOL (N workers for first_time, 1 worker for retry)                │
│                                                                             │
│  Worker 1  │  Worker 2  │  Worker 3  │  ...  │  Worker N                   │
│     ↓            ↓            ↓              ↓                              │
│  Update status: "processing"                                                │
│     ↓                                                                        │
│  Execute job (ProcessInsertCsvFile or ProcessExportCsvFile)                 │
│     ↓                                                                        │
│  Update status: "completed" or "failed"                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Components

| Component | Description | Location |
|-----------|-------------|----------|
| **Job Repository** | PostgreSQL storage for job state | `models/jobs.repo.go` |
| **Job Service** | Business logic for job creation/listing | `modules/common/service/jobService.go` |
| **Job Controller** | HTTP handlers for job API | `modules/common/controller/jobController.go` |
| **Job Runner** | Background worker that polls and processes jobs | `jobs/jobs.go` |
| **Job Consumer** | Worker that executes individual jobs | `jobs/jobs.go` |
| **CSV Processors** | Handlers for insert/export operations | `jobs/s3_files.go` |

---

## Job Types

### 1. Insert CSV File (`insert_csv_file`)

**Purpose**: Import CSV data from S3 into PostgreSQL and Elasticsearch.

**Data Flow**:
```
S3 CSV File → Streaming Reader → Batch Processing (500 records/batch) → 
Parallel Writes (PostgreSQL + Elasticsearch + Filters) → Complete
```

**Job Data Structure**:
```json
{
  "s3_key": "uploads/example.csv",
  "s3_bucket": "my-bucket-name"  // Optional, defaults to configured bucket
}
```

**Processing Details**:
- Streams CSV file from S3 (never loads full file into memory)
- Processes in configurable batches (default: 500 records)
- Performs parallel writes to:
  - PostgreSQL (contacts/companies tables)
  - Elasticsearch (contacts/companies indices)
  - Filters data table
- Handles large files (multi-GB) efficiently

### 2. Export CSV File (`export_csv_file`)

**Purpose**: Export filtered contact/company data to S3 as CSV.

**Data Flow**:
```
VQL Query → Cursor-based Pagination → Streaming Writer → S3 Upload → Complete
```

**Job Data Structure**:
```json
{
  "s3_bucket": "my-bucket-name",  // Optional, defaults to configured bucket
  "service": "contact",  // or "company"
  "vql": {
    "where": { ... },
    "select_columns": ["field1", "field2", ...],
    "order_by": [{ "order_by": "uuid", "order_direction": "desc" }],
    "limit": 500
  }
}
```

**Processing Details**:
- Uses cursor-based pagination (`search_after`) for efficient large dataset export
- Streams data directly to S3 using `io.Pipe()` (no memory buffering)
- Requires `select_columns` in VQL query
- Automatically sorts by `uuid` for stable pagination
- Output file: `{S3_UPLOAD_FILE_PATH}/{job_uuid}.csv`

**Important**: The `select_columns` field is **required** for export jobs. See [Select Columns Guide](./select_columns_filter.md) for details.

---

## Job States and Lifecycle

### State Machine

```
┌─────────┐
│  OPEN   │  ← Initial state when job is created
└────┬────┘
     │ (first_time runner picks up)
     ▼
┌──────────────┐
│  IN_QUEUE    │  ← Job is queued for processing
└──────┬───────┘
       │ (worker picks up)
       ▼
┌──────────────┐
│  PROCESSING  │  ← Job is being executed
└──────┬───────┘
       │
       ├─── Success ───▶ ┌─────────────┐
       │                 │  COMPLETED  │  ← Job finished successfully
       │                 └─────────────┘
       │
       └─── Failure ───▶ ┌─────────────┐
                         │   FAILED    │  ← Job failed, can be retried
                         └──────┬──────┘
                                │ (retry runner picks up after run_after time)
                                ▼
                         ┌──────────────────┐
                         │  RETRY_IN_QUEUED  │  ← Failed job queued for retry
                         └────────┬─────────┘
                                  │ (worker picks up)
                                  ▼
                         ┌──────────────┐
                         │  PROCESSING  │  ← Retry attempt
                         └──────────────┘
```

### Job States

| State | Description | Transition Trigger |
|-------|-------------|-------------------|
| `open` | Job created, waiting to be picked up | Job creation via API |
| `in_queue` | Job queued in channel, waiting for worker | First-time runner picks up job |
| `processing` | Job is currently being executed | Worker starts processing |
| `completed` | Job finished successfully | Job execution succeeds |
| `failed` | Job execution failed | Job execution fails |
| `retry_in_queued` | Failed job queued for retry | Retry runner picks up failed job |
| `retrying` | Job is being retried | Worker processes retry |

### Job Fields

| Field | Type | Description |
|-------|------|-------------|
| `uuid` | string | Unique job identifier (auto-generated) |
| `job_type` | string | Job type: `insert_csv_file` or `export_csv_file` |
| `data` | jsonb | Job-specific data (see Job Types section) |
| `status` | string | Current job state (see Job States) |
| `job_response` | jsonb | Response data (errors, S3 key, messages) |
| `retry_count` | integer | Number of retry attempts (default: 0) |
| `retry_interval` | integer | Seconds to wait before retry (default: 30) |
| `run_after` | timestamp | When to retry failed job (auto-calculated) |
| `created_at` | timestamp | Job creation time |
| `updated_at` | timestamp | Last update time |

### Job Response Structure

The `job_response` field contains:

```json
{
  "runtime_errors": ["error message 1", "error message 2"],
  "messages": "Success message",
  "s3_key": "uploads/abc123.csv"  // For export jobs
}
```

---

## API Reference

### Create Job

**Endpoint**: `POST /common/jobs/create`

**Authentication**: Required (`X-API-Key` header)

**Request Body**:
```json
{
  "job_type": "insert_csv_file",  // or "export_csv_file"
  "job_data": {
    // Job-specific data (see Job Types section)
  },
  "retry_count": 0  // Optional, default: 0
}
```

**Response** (201 Created):
```json
{
  "message": "Job created successfully",
  "success": true
}
```

**Error Responses**:

- **400 Bad Request**: Invalid request body, missing required fields
- **401 Unauthorized**: Invalid or missing API key
- **500 Internal Server Error**: Job creation failed

**Example - Create Insert Job**:
```bash
curl -X POST https://api.example.com/common/jobs/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "job_type": "insert_csv_file",
    "job_data": {
      "s3_key": "uploads/contacts.csv",
      "s3_bucket": "my-bucket"
    },
    "retry_count": 3
  }'
```

**Example - Create Export Job**:
```bash
curl -X POST https://api.example.com/common/jobs/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "job_type": "export_csv_file",
    "job_data": {
      "s3_bucket": "my-bucket",
      "service": "contact",
      "vql": {
        "where": {
          "keyword_match": {
            "must": {
              "country": ["USA"]
            }
          }
        },
        "select_columns": ["first_name", "last_name", "email", "title"],
        "order_by": [
          {
            "order_by": "uuid",
            "order_direction": "desc"
          }
        ],
        "limit": 500
      }
    }
  }'
```

### List Jobs

**Endpoint**: `POST /common/jobs`

**Authentication**: Required (`X-API-Key` header)

**Request Body**:
```json
{
  "job_type": "insert_csv_file",  // Optional: filter by job type
  "status": ["open", "processing"],  // Optional: filter by status(es)
  "limit": 25  // Optional: max results (default: 25, max: 100)
}
```

**Response** (200 OK):
```json
{
  "data": [
    {
      "id": 1,
      "uuid": "abc123-def456-ghi789",
      "job_type": "insert_csv_file",
      "data": {
        "s3_key": "uploads/contacts.csv",
        "s3_bucket": "my-bucket"
      },
      "status": "completed",
      "job_response": {
        "messages": "Import completed successfully"
      },
      "retry_count": 0,
      "retry_interval": 30,
      "run_after": null,
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:35:00Z"
    }
  ],
  "success": true
}
```

**Example - List All Open Jobs**:
```bash
curl -X POST https://api.example.com/common/jobs \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "status": ["open"],
    "limit": 50
  }'
```

**Example - List Failed Jobs**:
```bash
curl -X POST https://api.example.com/common/jobs \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "status": ["failed"],
    "limit": 100
  }'
```

---

## Job Runners

Jobs are processed by background workers that run as separate processes. There are two runner modes:

### First-Time Runner (`first_time`)

**Purpose**: Process new jobs (status: `open`)

**Command**: `./connectra jobs first_time`

**Characteristics**:
- Multiple workers (configurable via `PARALLEL_JOBS`, default: 4)
- Polls database every N minutes (configurable via `TICKER_INTERVAL_MINUTES`)
- Processes jobs with status: `open`
- Updates status to `in_queue` before processing
- Handles graceful shutdown (dequeues remaining jobs back to `open`)

**When to Use**: 
- Primary job processing
- High-throughput scenarios
- Production deployments

**Example**:
```bash
# Run first-time job processor
./connectra jobs first_time

# Or with Docker
docker run -d \
  --env-file .env \
  -e RUN_COMMAND="jobs first_time" \
  --name connectra-jobs \
  connectra:latest
```

### Retry Runner (`retry`)

**Purpose**: Retry failed jobs (status: `failed`)

**Command**: `./connectra jobs retry`

**Characteristics**:
- Single worker (controlled retry rate)
- Polls database every N minutes (configurable via `TICKER_INTERVAL_MINUTES`)
- Only processes jobs where `run_after` time has passed
- Updates status to `retry_in_queued` before processing
- Respects `retry_interval` for failed jobs

**When to Use**:
- Dedicated retry processing
- Controlled retry rate
- Separate from primary processing

**Example**:
```bash
# Run retry job processor
./connectra jobs retry

# Or with Docker
docker run -d \
  --env-file .env \
  -e RUN_COMMAND="jobs retry" \
  --name connectra-jobs-retry \
  connectra:latest
```

### Running Both Runners

For production deployments, run both runners:

```bash
# Terminal 1: First-time processor
./connectra jobs first_time

# Terminal 2: Retry processor
./connectra jobs retry
```

Or with Docker Compose:
```yaml
services:
  jobs-first-time:
    image: connectra:latest
    env_file: .env
    environment:
      - RUN_COMMAND=jobs first_time
    restart: unless-stopped

  jobs-retry:
    image: connectra:latest
    env_file: .env
    environment:
      - RUN_COMMAND=jobs retry
    restart: unless-stopped
```

---

## Configuration

### Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PARALLEL_JOBS` | Number of concurrent workers (first_time) | 4 | `8` |
| `BATCH_SIZE_FOR_INSERTION` | Records per batch for CSV processing | 500 | `1000` |
| `TICKER_INTERVAL_MINUTES` | Polling interval in minutes | 5 | `2` |
| `JOB_IN_QUEUE_SIZE` | Backpressure threshold (channel capacity) | 100 | `200` |
| `S3_BUCKET` | Default S3 bucket for jobs | - | `my-bucket` |
| `S3_UPLOAD_FILE_PATH_PRIFIX` | S3 path prefix for exports | `uploads` | `exports` |

### Configuration Example

```env
# Job Processing Configuration
PARALLEL_JOBS=8                    # 8 concurrent workers
BATCH_SIZE_FOR_INSERTION=1000      # Process 1000 records per batch
TICKER_INTERVAL_MINUTES=2           # Poll every 2 minutes
JOB_IN_QUEUE_SIZE=200               # Backpressure at 200 queued jobs

# S3 Configuration
S3_BUCKET=my-connectra-bucket
S3_UPLOAD_FILE_PATH_PRIFIX=exports
```

### Performance Tuning

**For High Throughput**:
- Increase `PARALLEL_JOBS` (e.g., 8-16 workers)
- Increase `BATCH_SIZE_FOR_INSERTION` (e.g., 1000-2000)
- Decrease `TICKER_INTERVAL_MINUTES` (e.g., 1-2 minutes)

**For Resource-Constrained Environments**:
- Decrease `PARALLEL_JOBS` (e.g., 2-4 workers)
- Decrease `BATCH_SIZE_FOR_INSERTION` (e.g., 250-500)
- Increase `TICKER_INTERVAL_MINUTES` (e.g., 5-10 minutes)

---

## Examples

### Example 1: Import Contacts CSV

**Step 1: Upload CSV to S3**

First, upload your CSV file to S3. You can use the upload URL endpoint:

```bash
# Get presigned upload URL
curl -X GET "https://api.example.com/common/upload-url?filename=contacts.csv" \
  -H "X-API-Key: your-secret-api-key"

# Response:
# {
#   "upload_url": "https://s3.amazonaws.com/...",
#   "s3_key": "uploads/contacts.csv"
# }

# Upload file using presigned URL
curl -X PUT "https://s3.amazonaws.com/..." \
  -H "Content-Type: text/csv" \
  --upload-file contacts.csv
```

**Step 2: Create Import Job**

```bash
curl -X POST https://api.example.com/common/jobs/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "job_type": "insert_csv_file",
    "job_data": {
      "s3_key": "uploads/contacts.csv",
      "s3_bucket": "my-bucket"
    },
    "retry_count": 3
  }'
```

**Step 3: Monitor Job Status**

```bash
# Check job status
curl -X POST https://api.example.com/common/jobs \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "status": ["processing", "completed", "failed"],
    "limit": 10
  }'
```

### Example 2: Export Companies to CSV

**Step 1: Create Export Job**

```bash
curl -X POST https://api.example.com/common/jobs/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "job_type": "export_csv_file",
    "job_data": {
      "s3_bucket": "my-bucket",
      "service": "company",
      "vql": {
        "where": {
          "keyword_match": {
            "must": {
              "industries": ["Software", "SaaS"],
              "country": ["USA"]
            }
          },
          "range_query": {
            "must": {
              "employees_count": {
                "gte": 50,
                "lte": 1000
              },
              "annual_revenue": {
                "gte": 1000000
              }
            }
          }
        },
        "select_columns": [
          "id",
          "name",
          "employees_count",
          "annual_revenue",
          "industries",
          "country",
          "website"
        ],
        "order_by": [
          {
            "order_by": "uuid",
            "order_direction": "desc"
          }
        ],
        "limit": 500
      }
    }
  }'
```

**Step 2: Get Export File Location**

Once the job completes, check the `job_response.s3_key` field:

```bash
curl -X POST https://api.example.com/common/jobs \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "status": ["completed"],
    "job_type": "export_csv_file",
    "limit": 1
  }'

# Response includes:
# {
#   "job_response": {
#     "s3_key": "uploads/abc123-def456.csv"
#   }
# }
```

**Step 3: Download CSV from S3**

Use your S3 client or generate a presigned download URL to retrieve the file.

### Example 3: Large Dataset Export

For exporting very large datasets (millions of records), the system automatically uses cursor-based pagination:

```json
{
  "job_type": "export_csv_file",
  "job_data": {
    "service": "contact",
    "vql": {
      "where": {
        "keyword_match": {
          "must": {
            "country": ["USA"]
          }
        }
      },
      "select_columns": [
        "first_name",
        "last_name",
        "email",
        "title",
        "company_id"
      ],
      "order_by": [
        {
          "order_by": "uuid",
          "order_direction": "desc"
        }
      ],
      "limit": 500
    }
  }
}
```

The system will:
1. Use `search_after` pagination automatically
2. Stream data directly to S3 (no memory buffering)
3. Handle millions of records efficiently

---

## Error Handling

### Common Errors

#### Job Creation Errors

**Error**: `ERR_MISSING_JOB_TYPE`
```json
{
  "error": "ERR_MISSING_JOB_TYPE: the 'job_type' field is required; specify a valid job type such as 'insert_csv_file' or 'export_csv_file'",
  "success": false
}
```
**Solution**: Include `job_type` field in request body.

**Error**: `ERR_MISSING_JOB_DATA`
```json
{
  "error": "ERR_MISSING_JOB_DATA: the 'job_data' field is required; include the necessary payload for job execution",
  "success": false
}
```
**Solution**: Include `job_data` field with appropriate structure.

**Error**: `ERR_INVALID_RETRY_COUNT`
```json
{
  "error": "ERR_INVALID_RETRY_COUNT: 'retry_count' must be a non-negative integer; use 0 for no retries or a positive number for retry attempts",
  "success": false
}
```
**Solution**: Use non-negative integer for `retry_count`.

#### Job Execution Errors

**Error**: `ERR_MISSING_SELECT_COLUMNS` (Export jobs)
```json
{
  "error": "ERR_MISSING_SELECT_COLUMNS: 'select_columns' is required for export operations; specify at least one column to include in the output",
  "success": false
}
```
**Solution**: Include `select_columns` array in VQL query for export jobs.

**Error**: S3 File Not Found (Insert jobs)
- Check `s3_key` and `s3_bucket` are correct
- Verify file exists in S3
- Check S3 permissions

**Error**: Invalid VQL Query (Export jobs)
- Validate VQL query structure
- Ensure `select_columns` is not empty
- Check field names match available columns

### Retry Logic

**Automatic Retries**:
- Failed jobs are automatically retried by the retry runner
- Retry interval: `retry_interval` seconds (default: 30)
- Retry count: Tracked in `retry_count` field
- Maximum retries: Controlled by `retry_count` set at job creation

**Retry Behavior**:
1. Job fails with error
2. Status updated to `failed`
3. `run_after` set to `current_time + retry_interval`
4. Retry runner picks up job after `run_after` time
5. Status updated to `retry_in_queued`, then `processing`
6. If successful: `completed`; if failed: back to `failed` with updated `run_after`

**Manual Retry**:
You can manually trigger a retry by updating a failed job's `run_after` to a past time, or by creating a new job with the same data.

---

## Best Practices

### Job Creation

1. **Set Appropriate Retry Count**:
   - Use `retry_count: 3` for critical jobs
   - Use `retry_count: 0` for jobs that should not retry

2. **Validate Data Before Creating Job**:
   - Verify S3 file exists before creating insert job
   - Test VQL query before creating export job
   - Ensure `select_columns` includes required fields

3. **Use Descriptive S3 Keys**:
   - Include timestamp: `uploads/contacts-2025-01-15.csv`
   - Include job type: `exports/companies-usa-2025-01-15.csv`

### Export Jobs

1. **Always Include `select_columns`**:
   - Required for export jobs
   - Only select fields you need (reduces processing time)
   - Include `uuid` for reference

2. **Use Cursor-Based Pagination**:
   - System automatically uses `search_after` for large datasets
   - Always include `uuid` in `order_by` for stable pagination
   - Use `limit: 500` for optimal batch size

3. **Optimize VQL Queries**:
   - Use specific filters to reduce dataset size
   - Combine filters efficiently (see [Filter Guides](./README.md))
   - Test query performance before creating export job

### Import Jobs

1. **CSV Format Requirements**:
   - First row must be headers
   - Headers should match database field names
   - Use UTF-8 encoding
   - Handle special characters properly

2. **File Size Considerations**:
   - System handles multi-GB files efficiently
   - Processing time depends on file size and batch size
   - Monitor job status for large imports

3. **Data Validation**:
   - Validate CSV structure before upload
   - Check required fields are present
   - Handle missing/null values appropriately

### Monitoring

1. **Regular Status Checks**:
   - Poll job status every few minutes
   - Monitor for stuck jobs (long `processing` status)
   - Check failed jobs regularly

2. **Error Review**:
   - Review `runtime_errors` in `job_response`
   - Fix data issues before retrying
   - Update job configuration if needed

3. **Performance Monitoring**:
   - Track job completion times
   - Monitor worker utilization
   - Adjust configuration based on load

### Production Deployment

1. **Run Both Runners**:
   - Deploy `first_time` runner for new jobs
   - Deploy `retry` runner for failed jobs
   - Use separate containers/processes

2. **Configure Resource Limits**:
   - Set appropriate `PARALLEL_JOBS` based on CPU/memory
   - Adjust `BATCH_SIZE_FOR_INSERTION` for optimal throughput
   - Monitor system resources

3. **Implement Monitoring**:
   - Set up alerts for failed jobs
   - Monitor job queue size
   - Track job processing rates

---

## Related Documentation

### Filter Documentation

- [Company Filters Guide](./01-company-filters-complete-guide.md) - Company filtering with VQL
- [Contact Filters Guide](./02-contact-filters-complete-guide.md) - Contact filtering with VQL
- [Select Columns Guide](./select_columns_filter.md) - Field selection optimization
- [API Reference](./06-api-reference.md) - Complete API endpoint reference

### Main Documentation

- [System Documentation](../system.md) - System architecture and setup
- [Company API](../company.md) - Company API documentation
- [Contact API](../contacts.md) - Contact API documentation
- [Main README](../README.md) - Documentation index

---

**Last Updated**: 2025-01-XX  
**Version**: 1.0
