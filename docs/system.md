# Connectra System Documentation

## Overview

Connectra is a high-performance contact and company management API built with Go. It provides advanced search capabilities using Elasticsearch, primary data storage with PostgreSQL, and file processing through S3 integration.

## Architecture

### Technology Stack

- **Language**: Go 1.24
- **Web Framework**: Gin
- **Primary Database**: PostgreSQL (using Bun ORM)
- **Search Engine**: Elasticsearch 8.x
- **File Storage**: AWS S3 (or S3-compatible)
- **CLI Framework**: Cobra
- **Configuration**: Viper

### System Components

```
┌─────────────┐
│   API Layer │  (Gin Router)
└──────┬──────┘
       │
┌──────▼──────────────────┐
│   Service Layer          │
│  - Company Service      │
│  - Contact Service      │
│  - Filter Service       │
└──────┬──────────────────┘
       │
┌──────▼─────────────────────────────┐
│   Repository Layer                  │
│  - PostgreSQL Repositories          │
│  - Elasticsearch Repositories       │
└──────┬──────────────────────────────┘
       │
┌──────▼──────────────────────────────┐
│   Data Layer                         │
│  - PostgreSQL (Primary Storage)      │
│  - Elasticsearch (Search Index)      │
│  - S3 (File Storage)                 │
└──────────────────────────────────────┘
```

## Core Features

### 1. Dual Storage Architecture

- **PostgreSQL**: Stores complete, normalized data with relationships
- **Elasticsearch**: Provides fast, flexible search capabilities
- Data is synchronized between both systems

### 2. VQL (Vivek Query Language)

A powerful query system that converts user-friendly queries into Elasticsearch queries:

- **Text Matching**: Supports exact, shuffle, and substring search types with fuzzy matching
- **Keyword Matching**: Exact value matching for categorical data
- **Range Queries**: Numeric and date range filtering
- **Complex Boolean Logic**: Must, Must Not, and Filter conditions
- **Field Selection**: `select_columns` parameter to limit returned fields

### 3. Filter System

- **Direct-Derived Filters** (`direct_derived: true`): Dynamically extracted from actual data records
  - Values are queried directly from the main data tables (companies/contacts)
  - **Company examples**: `address`, `annual_revenue`, `employees_count`, `linkedin_url`, `normalized_domain`, `total_funding`, `website`
  - **Contact examples**: `company_id`, `email`, `first_name`, `last_name`, `linkedin_url`, `mobile_phone`
- **Stored Filters** (`direct_derived: false`): Pre-computed filter values stored in `filters_data` table
  - Faster for frequently used filters with many distinct values
  - **Company examples**: `city`, `country`, `industries`, `keywords`, `state`, `technologies`, `uuid` (displayed as "Name")
  - **Contact examples**: `city`, `country`, `departments`, `email_status`, `seniority`, `state`, `title`
- Supports searchable filter dropdowns with pagination
- Each filter includes metadata: `key`, `display_name`, `service`, and `direct_derived` flag

### 4. Background Jobs

- S3 file processing jobs for bulk data import
- Retry mechanism with configurable intervals
- Job status tracking in PostgreSQL
- Parallel job processing with configurable workers
- Job queue management with configurable queue size

### 5. Bulk Operations and S3 File Processing

Connectra supports bulk data import through S3 file processing, enabling efficient import of large datasets.

#### S3 File Processing Overview

**Purpose**: Process large data files stored in S3 for bulk import of companies and contacts

**Workflow**:

1. Upload data file to S3 bucket
2. Create job record in `jobs_data` table
3. S3 job processor picks up job
4. File is downloaded and processed
5. Data is inserted into PostgreSQL
6. Data is indexed in Elasticsearch
7. Job status is updated

#### Job Configuration

**Environment Variables**:

```env
# Job Processing Configuration
PARALLEL_JOBS=5              # Number of parallel job workers
IN_QUEUE_SIZE=100            # Job queue buffer size
TICKER_INTERVAL_MINUTES=5    # Interval for checking new jobs

# S3 Configuration
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
S3_REGION=us-east-1
S3_BUCKET=your-bucket-name
S3_ENDPOINT=s3.amazonaws.com
S3_SSL=true
S3_DEBUG=false
```

#### Job Data Model

**Table**: `jobs_data`

**Fields**:

- `id` - Job ID (primary key)
- `uuid` - Job UUID (unique identifier)
- `job_type` - Type of job (e.g., "s3_file_insert")
- `data` - Job data (JSONB, contains file path, metadata, etc.)
- `status` - Job status (e.g., "pending", "processing", "completed", "failed")
- `retry_after` - Timestamp for retry (if job failed)
- `retry_interval` - Retry interval in seconds (default: 60)
- `remaining_retries` - Number of retries remaining (default: 3)
- `runtime_errors` - Array of error messages
- `created_at` - Job creation timestamp
- `updated_at` - Last update timestamp

#### Job Status Values

- `pending` - Job created, waiting to be processed
- `processing` - Job is currently being processed
- `completed` - Job completed successfully
- `failed` - Job failed (may retry if retries remaining)
- `cancelled` - Job was cancelled

#### Starting the S3 Job Processor

**Command**:

```bash
./connectra s3-job
```

**What It Does**:

- Starts background job processor
- Polls `jobs_data` table for pending jobs
- Processes jobs in parallel (configurable via `PARALLEL_JOBS`)
- Implements retry mechanism for failed jobs
- Updates job status in database

**Configuration**:

- **Parallel Workers**: `PARALLEL_JOBS` (default: 5)
- **Queue Size**: `IN_QUEUE_SIZE` (default: 100)
- **Polling Interval**: `TICKER_INTERVAL_MINUTES` (default: 5 minutes)

#### Retry Mechanism

**How It Works**:

1. Job fails during processing
2. System checks `remaining_retries`
3. If retries remaining:
   - Job status set to `pending`
   - `retry_after` set to current time + `retry_interval`
   - `remaining_retries` decremented
   - Error added to `runtime_errors`
4. Retry worker picks up job after `retry_after` time
5. Process repeats until success or retries exhausted

**Retry Configuration**:

- **Default Retries**: 3 attempts
- **Default Interval**: 60 seconds
- **Configurable**: Via `retry_interval` and `remaining_retries` in job data

#### Bulk Import Best Practices

1. **File Format**:
   - Use JSON or CSV format
   - Ensure proper encoding (UTF-8)
   - Validate data structure before upload

2. **File Size**:
   - Recommended: < 100MB per file
   - For larger datasets: Split into multiple files
   - Process files in batches

3. **Job Management**:
   - Monitor job status via `jobs_data` table
   - Check `runtime_errors` for processing issues
   - Implement job status polling in your application

4. **Error Handling**:
   - Review `runtime_errors` for failed jobs
   - Fix data issues and retry
   - Monitor `remaining_retries` to prevent job loss

5. **Performance**:
   - Adjust `PARALLEL_JOBS` based on system capacity
   - Use appropriate `IN_QUEUE_SIZE` for your workload
   - Monitor database and Elasticsearch performance during bulk imports

#### Job Tracking Example

**Query Job Status**:

```sql
SELECT 
  uuid,
  job_type,
  status,
  remaining_retries,
  runtime_errors,
  created_at,
  updated_at
FROM jobs_data
WHERE status IN ('pending', 'processing', 'failed')
ORDER BY created_at DESC;
```

**Monitor Job Progress**:

```sql
SELECT 
  status,
  COUNT(*) as count
FROM jobs_data
GROUP BY status;
```

#### S3 File Processing Flow

```
1. Upload File to S3
   ↓
2. Create Job Record (status: 'pending')
   ↓
3. S3 Job Processor Polls for Jobs
   ↓
4. Job Picked Up (status: 'processing')
   ↓
5. Download File from S3
   ↓
6. Parse and Validate Data
   ↓
7. Insert into PostgreSQL
   ↓
8. Index in Elasticsearch
   ↓
9. Update Job Status (status: 'completed')
```

**Error Handling Flow**:

```
Job Fails
   ↓
Check remaining_retries
   ↓
If > 0:
   - Set status: 'pending'
   - Set retry_after: now + retry_interval
   - Decrement remaining_retries
   - Add error to runtime_errors
   ↓
Retry Worker Picks Up After retry_after
   ↓
If remaining_retries = 0:
   - Set status: 'failed'
   - Job will not retry automatically
```

### 6. Write Operations and Data Management

**New Feature** (December 2024): Connectra now supports full CRUD operations for both contacts and companies.

#### Write Service Architecture

The write service layer handles all data modifications with automatic dual-storage synchronization:

**Write Flow**:

```
API Request
   ↓
Controller Layer (validation)
   ↓
Write Service Layer
   ↓
PostgreSQL Repository (transaction)
   ↓
Commit Transaction
   ↓
Elasticsearch Repository (async indexing)
```

#### Supported Operations

**Contact Write Operations**:

- `POST /contacts/create` - Create single contact
- `PUT /contacts/:uuid` - Update contact by UUID
- `DELETE /contacts/:uuid` - Soft delete contact
- `POST /contacts/upsert` - Create or update contact
- `POST /contacts/bulk` - Bulk upsert contacts

**Company Write Operations**:

- `POST /companies/create` - Create single company
- `PUT /companies/:uuid` - Update company by UUID
- `DELETE /companies/:uuid` - Soft delete company
- `POST /companies/upsert` - Create or update company
- `POST /companies/bulk` - Bulk upsert companies

#### Key Features

**1. Automatic Elasticsearch Indexing**

- All write operations automatically index to Elasticsearch
- Async/fire-and-forget pattern (doesn't block writes)
- Resilient to ES failures (logs warning, doesn't fail write)
- Full document reindexing on updates

**2. Validation Layer**

- Request validation using Gin binding tags
- Email format validation
- LinkedIn URL format validation
- Required field validation (e.g., name, email)
- UUID format validation
- Numeric field validation (non-negative values)

**3. Upsert Logic**

- Contacts: Match by UUID or email
- Companies: Match by UUID or normalized_domain
- Uses PostgreSQL `ON CONFLICT` for atomic operations
- Preserves UUID on updates

**4. Bulk Operations**

- Efficient batch processing using PostgreSQL bulk insert
- Batch Elasticsearch indexing
- Supports 1000+ records per request
- Atomic transaction handling

**5. Soft Delete**

- Sets `deleted_at` timestamp instead of hard delete
- Removes from Elasticsearch index
- Retains in PostgreSQL for audit purposes
- Preserves data integrity

#### Error Handling

Write operations return clear error messages:

- **400 Bad Request**: Validation errors (invalid email, missing fields, etc.)
- **404 Not Found**: Resource not found (update/delete non-existent record)
- **500 Internal Server Error**: Database or Elasticsearch errors

#### Performance Characteristics

- **Single Operations**: < 100ms (PostgreSQL write + ES index)
- **Bulk Operations**: ~1-2 seconds for 500 records
- **Elasticsearch Lag**: < 1 second for search availability
- **Transaction Safety**: All PostgreSQL writes are transactional

#### Migration Impact

The write service enables:

- Direct data management through API (no direct database access needed)
- Automatic search index synchronization
- Simplified backend integrations (removed repository layer)
- Centralized data validation and business logic

See `backend/WRITE_MIGRATION_COMPLETE.md` for details on the migration from direct database access to Connectra write API.

## Configuration

### Environment Variables

The system uses Viper for configuration management. Create a `.env` file based on `.example.env`:

#### Application Configuration

```env
APP_ENV=development
PARALLEL_JOBS=5
IN_QUEUE_SIZE=100
TICKER_INTERVAL_MINUTES=5
API_KEY=your-secret-api-key
MAX_REQUESTS_PER_MINUTE=60
```

#### PostgreSQL Configuration

```env
PG_DB_CONNECTION=postgres
PG_DB_HOST=localhost
PG_DB_PORT=5432
PG_DB_DATABASE=connectra
PG_DB_USERNAME=postgres
PG_DB_PASSWORD=password
PG_DB_DEBUG=false
PG_DB_SSL=false
```

#### Elasticsearch Configuration

```env
ELASTICSEARCH_CONNECTION=http
ELASTICSEARCH_HOST=localhost
ELASTICSEARCH_PORT=9200
ELASTICSEARCH_USERNAME=elastic
ELASTICSEARCH_PASSWORD=password
ELASTICSEARCH_DEBUG=false
ELASTICSEARCH_SSL=false
ELASTICSEARCH_AUTH=false
```

#### S3 Configuration

```env
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
S3_REGION=us-east-1
S3_BUCKET=your-bucket-name
S3_ENDPOINT=s3.amazonaws.com
S3_SSL=true
S3_DEBUG=false
```

## Setup and Installation

### Prerequisites

1. Go 1.24 or higher
2. PostgreSQL 12+
3. Elasticsearch 8.x
4. (Optional) S3-compatible storage

### Installation Steps

1. **Clone and Install Dependencies**

   ```bash
   cd connectra
   go mod download
   ```

2. **Configure Environment**

   ```bash
   cp .example.env .env
   # Edit .env with your configuration
   ```

3. **Set Up Database**
   - Create PostgreSQL database
   - Run migrations (if available)

4. **Set Up Elasticsearch**
   - Start Elasticsearch instance
   - Create indices using examples in `examples/` directory:
     - `company_index_create.json`
     - `contact_index_create.json`

5. **Build and Run**

   ```bash
   # Build
   go build -o connectra .
   
   # Run API Server
   ./connectra api-server
   
   # Run S3 Job Processor
   ./connectra s3-job
   ```

## API Server

### Starting the Server

```bash
./connectra api-server
```

The server starts on port `8000` by default.

### Middleware

The API server uses the following middleware in order:

1. **Logger**: Gin's built-in request logger
2. **Recovery**: Panic recovery middleware
3. **CORS**: Configured to allow all origins with common HTTP methods and headers
4. **GZIP**: Response compression using default compression level
5. **Rate Limiting**: Token bucket algorithm limiting requests per minute (configured via `MAX_REQUESTS_PER_MINUTE`)
6. **API Key Authentication**: Requires `X-API-Key` header matching configured `API_KEY`

### Authentication

All API endpoints require authentication using an API Key:

**Request Header**:

```
X-API-Key: your-secret-api-key
```

**Unauthorized Response** (401):

```json
{
  "error": "unauthorized",
  "message": "invalid API key"
}
```

### Rate Limiting

The API implements a token bucket rate limiter:

- **Algorithm**: Token bucket with per-second refill
- **Configuration**: `MAX_REQUESTS_PER_MINUTE` environment variable
- **Default**: 60 requests per minute
- **Tokens refill**: Continuously at rate of `MAX_REQUESTS_PER_MINUTE / 60` per second

**Rate Limit Exceeded Response** (429):

```json
{
  "error": "rate limit exceeded",
  "message": "too many requests, please try again later"
}
```

### Health Check

```bash
GET /health
```

**Note**: Health check endpoint does NOT require authentication or count against rate limits.

Returns:

```json
{
  "status": "ok"
}
```

## CLI Commands

### Available Commands

1. **api-server**: Start the REST API server

   ```bash
   ./connectra api-server
   ```

2. **s3-job**: Start the S3 file processing job worker

   ```bash
   ./connectra s3-job
   ```

## Data Flow

### Search Flow

1. Client sends VQL query to API with authentication header
2. Middleware validates API key and checks rate limits
3. Service layer converts VQL to Elasticsearch query
4. Elasticsearch returns matching document IDs
5. PostgreSQL repository fetches full records by IDs (filtered by `select_columns` if specified)
6. Results returned to client

### Filter Data Flow

1. Client requests filter data
2. System checks if filter is `direct_derived`:
   - **False**: Query `filters_data` table
   - **True**: Query actual data table and extract field values
3. Results filtered by search text (if provided)
4. Paginated results returned

## Database Schema

### Core Tables

- **companies**: Company master data
- **contacts**: Contact master data
- **filters**: Filter definitions
- **filters_data**: Pre-computed filter values
- **jobs_data**: Background job tracking

## Elasticsearch Indices

- **companies_index**: Company search index
- **contacts_index**: Contact search index

Both indices use:

- N-gram analyzers for partial text matching
- Keyword fields for exact matching
- Date fields for temporal queries

## Performance Considerations

### Pagination Limits

- Maximum page size: 100 records
- Maximum page number: 10 (for Elasticsearch queries)
- Default page size: 25 records

### Connection Pooling

- PostgreSQL: Max 100 open connections, 50 idle
- Elasticsearch: Max 100 idle connections per host

### Caching Strategy

- Filter data can be cached at application level
- Elasticsearch query results are not cached (real-time search)

## Error Handling

Common error constants are defined in `constants/errors_message.go`:

- Invalid request body
- Page size/number exceeded
- Database connection errors
- Elasticsearch query errors

## Development

### Project Structure

```
connectra/
├── cmd/              # CLI commands
├── clients/          # Database/storage clients
├── conf/             # Configuration management
├── connections/      # Connection initialization
├── constants/        # Application constants
├── jobs/             # Background job processors
├── models/           # Data models and repositories
├── modules/          # API modules (companies, contacts)
│   ├── companies/
│   │   ├── controller/
│   │   ├── service/
│   │   └── helper/
│   └── contacts/
├── utilities/        # Common utilities
└── docs/             # Documentation
```

### Adding New Features

1. **New Entity**: Create models in `models/`, add repositories, services, and controllers
2. **New Endpoint**: Add route in module's `routes.go`, implement controller
3. **New Filter**: Add entry to `filters` table, implement in filter service

## Monitoring and Logging

- Uses `zerolog` for structured logging
- Log levels: Debug, Info, Error
- Database query logging available when `PG_DB_DEBUG=true`
- Elasticsearch request logging when `ELASTICSEARCH_DEBUG=true`

## Security Considerations

- **API Key Authentication**: All endpoints (except `/health`) require valid `X-API-Key` header
- **Rate Limiting**: Token bucket algorithm prevents API abuse (configurable via `MAX_REQUESTS_PER_MINUTE`)
- **Environment Variables**: Sensitive data stored in environment variables, never in code
- **SQL Injection Prevention**: Parameterized queries with Bun ORM
- **Elasticsearch Query Validation**: VQL queries validated before conversion
- **Connection Pooling Limits**: Prevent resource exhaustion
- **CORS Configuration**: Configurable origins (currently allows all for development)
- **SSL/TLS**: Configurable for PostgreSQL, Elasticsearch, and S3 connections

## Troubleshooting

### Common Issues

1. **Connection Errors**
   - Verify database/Elasticsearch credentials
   - Check network connectivity
   - Verify ports are accessible

2. **Query Errors**
   - Check Elasticsearch index exists
   - Verify field names match index mappings
   - Check pagination limits

3. **Performance Issues**
   - Monitor connection pool usage
   - Check Elasticsearch cluster health
   - Review query complexity

## Future Enhancements

- **Enhanced Authentication**: JWT-based authentication, OAuth2 integration, role-based access control (RBAC)
- **API Versioning**: Version-based routing (e.g., `/v1`, `/v2`)
- **WebSocket Support**: Real-time updates and notifications
- **GraphQL API**: Alternative query interface
- **Advanced Analytics**: Comprehensive reporting and metrics dashboard
- **Caching Layer**: Redis integration for frequently accessed data
- **Audit Logging**: Comprehensive audit trail for all API operations
- **Multi-tenancy**: Support for multiple organizations/workspaces
