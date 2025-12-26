<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"/>
  <img src="https://img.shields.io/badge/Elasticsearch-8.x-005571?style=for-the-badge&logo=elasticsearch&logoColor=white" alt="Elasticsearch"/>
  <img src="https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL"/>
  <img src="https://img.shields.io/badge/AWS_S3-232F3E?style=for-the-badge&logo=amazons3&logoColor=white" alt="AWS S3"/>
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker"/>
</p>

# Connectra API

> A high-performance B2B contact and company data platform with advanced search capabilities, real-time data ingestion, and a custom domain-specific query language.

## Overview

Connectra is a scalable backend service designed for B2B data enrichment platforms. It provides lightning-fast search across millions of contacts and companies using a hybrid PostgreSQL + Elasticsearch architecture, with AWS S3 integration for bulk data imports.

### Key Highlights

- **Custom Query Language (VQL)** - Domain-specific query language for complex search operations
- **Hybrid Search Architecture** - Elasticsearch for full-text search + PostgreSQL for relational data
- **Background Job Processing** - Async file processing with retry mechanisms and graceful shutdown
- **Production-Ready Security** - API key authentication with token bucket rate limiting
- **Concurrent Processing** - Parallel bulk upserts with worker pools for high throughput

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Connectra API                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐   │
│  │   Gin Web   │    │ Middleware  │    │ Controllers │    │  Services   │   │
│  │   Router    │───▶│ Auth + Rate │───▶│   Layer     │───▶│   Layer     │   │
│  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘   │
│                                                                  │          │
│                                              ┌───────────────────┼──────┐   │
│                                              ▼                   ▼      │   │
│                                     ┌─────────────┐    ┌─────────────┐  │   │
│                                     │ Elasticsearch│   │ PostgreSQL  │  │   │
│                                     │  Repository  │   │ Repository  │  │   │
│                                     └─────────────┘    └─────────────┘  │   │
│                                                                         │   │
├─────────────────────────────────────────────────────────────────────────┼───┤
│  Background Jobs (Worker Pool)                                          │   │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                  │   │
│  │ Job Queue   │───▶│  Workers    │───▶│ Batch Upsert│──────────────────┘   │
│  │ (Channel)   │    │ (Goroutines)│    │  Service    │                      │
│  └─────────────┘    └─────────────┘    └─────────────┘                      │
│         ▲                                    │                              │
│         │                                    ▼                              │
│  ┌─────────────┐                    ┌─────────────┐                         │
│  │   AWS S3    │◀───────────────────│ CSV Parser  │                         │
│  │ File Stream │                    │             │                         │
│  └─────────────┘                    └─────────────┘                         │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Web Framework** | Gin | High-performance HTTP router with middleware support |
| **Primary DB** | PostgreSQL + Bun ORM | Relational data storage with type-safe queries |
| **Search Engine** | Elasticsearch 8.x | Full-text search, fuzzy matching, and aggregations |
| **Object Storage** | AWS S3 | Bulk file uploads and streaming |
| **Configuration** | Viper | Environment-based configuration management |
| **CLI** | Cobra | Command-line interface for different run modes |
| **Logging** | Zerolog | Structured, zero-allocation JSON logging |

---

## Features

### 1. VQL - Vivek Query Language

A custom DSL for querying contacts and companies with support for:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "filter_key": "title",
          "text_value": "Software Engineer",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "seniority": ["senior", "lead"],
        "country": "united states"
      }
    },
    "range_query": {
      "must": {
        "company_employees_count": { "gte": 100, "lte": 1000 }
      }
    }
  },
  "order_by": [{ "order_by": "created_at", "order_direction": "desc" }],
  "page": 1,
  "limit": 25,
  "company_config": {
    "populate": true,
    "select_columns": ["name", "website", "industries"]
  }
}
```

**Search Types:**
- `exact` - Phrase matching with configurable slop
- `shuffle` - Multi-word matching with AND/OR operators
- `substring` - N-gram based partial matching

### 2. Hybrid Search Architecture

```
Query Flow:
┌──────────┐     ┌──────────────────┐     ┌────────────────┐
│  VQL     │────▶│  Elasticsearch   │────▶│   Get IDs      │
│  Query   │     │  (Fast Search)   │     │                │
└──────────┘     └──────────────────┘     └───────┬────────┘
                                                  │
                                                  ▼
┌──────────┐     ┌──────────────────┐     ┌────────────────┐
│  Result  │◀────│   PostgreSQL     │◀────│  Fetch Full    │
│          │     │   (Rich Data)    │     │  Records       │
└──────────┘     └──────────────────┘     └────────────────┘
```

- Elasticsearch handles text search, filtering, and pagination
- PostgreSQL returns complete records with all fields
- Concurrent fetching of contacts and companies using goroutines

### 3. Background Job Processing

```go
// Worker pool with configurable parallelism
for i := 0; i < conf.JobConfig.ParallelJobs; i++ {
    wg.Add(1)
    go insertJob.Run(&wg, 0, jobsChannel)
}
```

**Features:**
- Ticker-based job polling with configurable intervals
- Buffered channels for job queuing (1000 capacity)
- Automatic retry mechanism for failed jobs
- Graceful shutdown with context cancellation
- Job status tracking: `open` → `in_queue` → `processing` → `completed/failed`

### 4. Concurrent Batch Upserts

```go
// Parallel writes to 4 data stores
wg.Add(4)
go func() { s.companyRepo.BulkUpsert(pgCompanies) }()
go func() { s.contactRepo.BulkUpsert(pgContacts) }()
go func() { s.companyElasticRepo.BulkUpsert(esCompanies) }()
go func() { s.contactElasticRepo.BulkUpsert(esContacts) }()
wg.Wait()
```

### 5. Security & Rate Limiting

- **API Key Authentication** - Header-based (`X-API-Key`) authentication middleware
- **Token Bucket Rate Limiter** - Configurable requests/minute with automatic token refill

```go
// Token bucket implementation
func RateLimiter() gin.HandlerFunc {
    // Tokens refill every second
    // Burst protection with configurable limits
}
```

---

## API Endpoints

### Contacts
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/contacts/filter` | Query contacts with VQL |
| `POST` | `/contacts/count` | Get count of matching contacts |

### Companies
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/companies/filter` | Query companies with VQL |
| `POST` | `/companies/count` | Get count of matching companies |

### Common
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/common/filters/:service` | Get available filters for a service |
| `POST` | `/common/filters/:service/data` | Get filter options/values |
| `GET` | `/common/upload-url` | Generate S3 presigned upload URL |
| `POST` | `/common/jobs` | Create a new background job |
| `GET` | `/common/jobs/:id` | Get job status |

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check endpoint |

---

## Project Structure

```
connectra/
├── cmd/                      # CLI commands (Cobra)
│   ├── root.go              # Root command with graceful shutdown
│   ├── server.go            # API server command
│   └── s3_file_insertions.go # Background job runner
├── conf/                     # Configuration management
│   └── viper.go             # Viper-based env config
├── connections/              # Database & service connections
│   ├── database.go          # PostgreSQL connection pool
│   ├── s3_connection.go     # AWS S3 client
│   └── search_engine.go     # Elasticsearch client
├── clients/                  # Low-level client implementations
│   ├── pgsql.go
│   ├── elastic_search.go
│   ├── s3.go
│   └── mongo.go
├── middleware/               # HTTP middlewares
│   ├── authMiddleware.go    # API key authentication
│   └── rateMiddleware.go    # Token bucket rate limiter
├── models/                   # Data models & repositories
│   ├── contact.pgsql.go     # PostgreSQL contact model
│   ├── contact.elastic.go   # Elasticsearch contact model
│   ├── company.pgsql.go
│   ├── company.elastic.go
│   └── *.repo.go            # Repository implementations
├── modules/                  # Feature modules
│   ├── contacts/
│   │   ├── controller/      # HTTP handlers
│   │   ├── service/         # Business logic
│   │   ├── helper/          # Request/response helpers
│   │   └── routes.go
│   ├── companies/
│   └── common/
├── jobs/                     # Background job workers
│   ├── s3_file_insertions.go # S3 file processing job
│   └── insert_direct_file.go
├── utilities/                # Shared utilities
│   ├── query.go             # VQL to Elasticsearch converter
│   ├── structures.go        # Common data structures
│   └── common.go
├── constants/                # Application constants
├── Dockerfile               # Multi-stage Docker build
└── main.go                  # Application entry point
```

---

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 15+
- Elasticsearch 8.x
- AWS S3 (or compatible storage)

### Environment Variables

```env
# Application
APP_ENV=development
API_KEY=your-secret-api-key
MAX_REQUESTS_PER_MINUTE=1000

# PostgreSQL
PG_DB_HOST=localhost
PG_DB_PORT=5432
PG_DB_DATABASE=connectra
PG_DB_USERNAME=postgres
PG_DB_PASSWORD=password
PG_DB_SSL=false
PG_DB_DEBUG=true

# Elasticsearch
ELASTICSEARCH_HOST=localhost
ELASTICSEARCH_PORT=9200
ELASTICSEARCH_USERNAME=elastic
ELASTICSEARCH_PASSWORD=password
ELASTICSEARCH_SSL=false
ELASTICSEARCH_AUTH=true

# AWS S3
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
S3_REGION=us-east-1
S3_BUCKET=connectra-uploads
S3_ENDPOINT=https://s3.amazonaws.com

# Jobs
PARALLEL_JOBS=4
BATCH_SIZE_FOR_INSERTION=500
TICKER_INTERVAL_MINUTES=5
JOB_IN_QUEUE_SIZE=100
```

### Running Locally

```bash
# Clone the repository
git clone https://github.com/yourusername/connectra.git
cd connectra

# Install dependencies
go mod download

# Run the API server
go run main.go api-server

# Or run the background job processor
go run main.go insert-job
```

### Docker

```bash
# Build the image
docker build -t connectra:latest .

# Run API server
docker run -d \
  --env-file .env \
  -e RUN_COMMAND=api-server \
  -p 8000:8000 \
  connectra:latest

# Run job processor
docker run -d \
  --env-file .env \
  -e RUN_COMMAND=insert-job \
  connectra:latest
```

---

## Design Patterns & Best Practices

| Pattern | Implementation |
|---------|----------------|
| **Repository Pattern** | Decoupled data access layer for PostgreSQL and Elasticsearch |
| **Dependency Injection** | Services receive repository interfaces, not concrete implementations |
| **Worker Pool** | Configurable goroutine pool for parallel job processing |
| **Graceful Shutdown** | Context cancellation with signal handling (SIGTERM/SIGINT) |
| **Middleware Chain** | Composable middleware for auth, rate limiting, CORS, gzip |
| **Clean Architecture** | Separation of concerns across controller, service, repository layers |

---

## Performance Optimizations

- **Connection Pooling** - Reused database connections via Bun ORM
- **Gzip Compression** - Response compression middleware
- **Streaming CSV Parser** - Memory-efficient file processing
- **Bulk Operations** - Batch inserts for both PostgreSQL and Elasticsearch
- **Concurrent I/O** - Parallel writes to multiple data stores
- **Presigned URLs** - Direct S3 uploads without proxying through the API

---

## License

MIT License - feel free to use this project as a reference or starting point.

---

<p align="center">
  <sub>Built with ❤️ in Go</sub>
</p>
