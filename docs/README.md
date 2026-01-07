# Connectra API Documentation

## Overview

Welcome to the Connectra API documentation! Connectra is a high-performance contact and company management API that provides advanced search capabilities, comprehensive filtering, and efficient data management.

**Version**: 1.2  
**Last Updated**: 2025-01-XX

## Quick Start

### Deployment Options

Connectra can be deployed in two ways:

1. **Lambda Deployment** (Recommended for Production)
   - Serverless deployment on AWS Lambda
   - Auto-scaling and cost-effective
   - See: [Lambda Deployment Guide](./LAMBDA_DEPLOYMENT.md)

2. **Server Deployment** (Traditional)
   - Run as a traditional HTTP server
   - See: [System Documentation](./system.md#api-server)

### For New Users

1. **Start Here**: [System Documentation](./system.md) - Understand the architecture and setup
2. **Deployment**: [Lambda Deployment Guide](./LAMBDA_DEPLOYMENT.md) - Deploy to AWS Lambda
3. **API Basics**: [Company API](./company.md) or [Contact API](./contacts.md) - Learn the core APIs
4. **Filtering**: [Filter Documentation Index](./filters/README.md) - Master filtering capabilities
5. **Examples**: [Examples and Use Cases](./filters/05-examples-use-cases.md) - See practical examples

### For Developers

1. **Setup**: [System Documentation](./system.md#development) - Get started with development
2. **API Reference**: [Filter API Reference](./filters/06-api-reference.md) - Complete endpoint reference
3. **Field Reference**: [Filter Field Reference](./filters/04-filter-field-reference.md) - All available fields
4. **Integration**: See examples below for common integration patterns

### For Business Users

1. **Use Cases**: [Examples and Use Cases](./filters/05-examples-use-cases.md) - Real-world scenarios
2. **Filtering Strategies**: [Combined Filters Guide](./filters/03-combined-filters-guide.md) - Account-based filtering
3. **Field Guide**: [Filter Field Reference](./filters/04-filter-field-reference.md) - What you can filter by

## Documentation Structure

```
docs/
├── README.md (this file)
├── system.md                    # System architecture and technical details
├── company.md                   # Company API documentation
├── contacts.md                  # Contact API documentation
└── filters/                     # Comprehensive filter documentation
    ├── README.md               # Filter documentation index
    ├── 01-company-filters-complete-guide.md
    ├── 02-contact-filters-complete-guide.md
    ├── 03-combined-filters-guide.md
    ├── 04-filter-field-reference.md
    ├── 05-examples-use-cases.md
    ├── 06-api-reference.md
    └── select_columns_filter.md
```

## Core Documentation Files

### 1. [System Documentation](./system.md)

**Complete system architecture and technical reference**

- Technology stack (Go, Gin, PostgreSQL, Elasticsearch)
- Dual storage architecture
- System components and data flow
- Database schema
- Elasticsearch indices
- Performance considerations
- Security features
- Development setup
- Troubleshooting guide

**Best for**: Understanding the system architecture, setup, and technical details

---

### 2. [Company API](./company.md)

**Company data management and filtering**

- Company endpoints
- Company data model
- Filter capabilities
- Field reference
- Common operations

**Best for**: Working with company data

---

### 3. [Contact API](./contacts.md)

**Contact data management and filtering**

- Contact endpoints
- Contact data model
- Filter capabilities
- Denormalized company fields
- Field reference
- Common operations

**Best for**: Working with contact data

---

### 4. [Filter Documentation](./filters/README.md)

**Comprehensive filtering guide**

- Company filters complete guide
- Contact filters complete guide
- Combined filtering strategies
- Field reference
- Examples and use cases
- API reference

**Best for**: Mastering all filtering capabilities

---

### 5. [Jobs API](./filters/jobs.md)

**Background job processing for large operations**

- CSV import/export jobs
- Asynchronous processing
- Job status tracking
- Retry mechanisms
- Job runners and configuration

**Best for**: Bulk data operations, CSV imports/exports

---

## Key Features

### 🔍 Advanced Search

- **Text Search**: Full-text search with exact, shuffle, and substring matching
- **Fuzzy Matching**: Typo-tolerant search
- **Ngram Analysis**: Efficient partial text matching
- **Keyword Matching**: Exact value matching for categorical data
- **Range Queries**: Numeric and date range filtering

### 🏢 Company Management

- Complete company profiles
- Industry and technology tracking
- Financial metrics (revenue, funding, employees)
- Geographic data
- Social media links

### 👥 Contact Management

- Contact profiles with job titles and departments
- Email verification status
- Seniority levels
- Company relationships
- Denormalized company data for efficient filtering

### 🔗 Account-Based Filtering

- Filter contacts by company attributes in a single query using denormalized fields (`company_*` prefix)
- Populate full company objects in contact responses using `company_config`
- Efficient account-based marketing workflows
- No need for separate company queries when filtering contacts

> **Note**: Denormalized `company_*` fields are **ONLY for filtering** in `where` clauses. To get company data in responses, use `company_config.populate: true` with `company_config.select_columns`. See [Select Columns Guide](./filters/select_columns_filter.md) for details.

### ✏️ Write Operations

- **Full CRUD Support**: Create, update, delete operations for contacts and companies
- **Bulk Operations**: Efficient bulk upsert for data imports
- **Automatic Indexing**: All writes automatically sync to Elasticsearch
- **Validation**: Comprehensive request validation with clear error messages
- **Upsert Logic**: Smart create-or-update based on UUID or unique fields

### 📊 Background Jobs

- **CSV Import**: Bulk import from S3 to database
- **CSV Export**: Export filtered data to S3
- **Asynchronous Processing**: Non-blocking job execution
- **Retry Logic**: Automatic retry for failed jobs
- **Status Tracking**: Complete job lifecycle monitoring

### ⚡ Performance

- Dual storage architecture (PostgreSQL + Elasticsearch)
- Fast Elasticsearch queries
- Connection pooling
- Efficient pagination (page-based and cursor-based)
- Field selection optimization

## Common Workflows

### 1. Lead Generation

**Goal**: Find high-value companies and decision-makers

**Steps**:
1. Filter companies by industry, size, revenue → [Company Filters Guide](./filters/01-company-filters-complete-guide.md)
2. Filter contacts at those companies by seniority, email status → [Contact Filters Guide](./filters/02-contact-filters-complete-guide.md)
3. Or use single query with denormalized fields → [Combined Filters Guide](./filters/03-combined-filters-guide.md)

**Example**: [Lead Generation Examples](./filters/05-examples-use-cases.md#lead-generation)

---

### 2. Account-Based Marketing

**Goal**: Target specific accounts with multiple contacts

**Steps**:
1. Identify target companies by criteria
2. Find all contacts at those companies
3. Filter contacts by role, seniority, department
4. Use denormalized company fields for single-query efficiency

**Example**: [Account-Based Marketing Examples](./filters/03-combined-filters-guide.md#account-based-filtering-patterns)

---

### 3. Competitive Intelligence

**Goal**: Analyze competitors and their key personnel

**Steps**:
1. Find competitors by industry, technology stack
2. Identify key contacts at competitor companies
3. Track technology usage and company metrics

**Example**: [Competitive Analysis Examples](./filters/05-examples-use-cases.md#competitive-analysis)

---

### 4. Email Campaign Targeting

**Goal**: Build targeted email lists

**Steps**:
1. Filter contacts by email status (`verified`)
2. Filter by company attributes (industry, size, location)
3. Filter by contact attributes (seniority, department)
4. Export for email campaigns

**Example**: [Email Campaign Examples](./filters/05-examples-use-cases.md#email-campaign-targeting)

---

## Authentication

All API endpoints (except `/health`) require authentication using an API key.

**Header Required**:
```
X-API-Key: your-secret-api-key
```

**Complete Example** (Lambda Deployment - Production):
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "industries": ["Software"]
        }
      }
    },
    "page": 1,
    "limit": 25
  }'
```

**Note**: For local development, use `http://localhost:8000` instead of the Lambda URL.

**See**: [Authentication Documentation](./filters/06-api-reference.md#authentication) for complete details

---

## Rate Limiting

The API implements token bucket rate limiting to prevent abuse.

**Default**: 60 requests per minute  
**Configurable**: Via `MAX_REQUESTS_PER_MINUTE` environment variable

**Response**: `429 Too Many Requests` when limit exceeded

**See**: [Rate Limiting Documentation](./filters/06-api-reference.md#rate-limiting) for details and best practices

---

## Error Handling

All errors follow a consistent format:

```json
{
  "error": "error message",
  "success": false
}
```

**Common Error Codes**:
- `400 Bad Request` - Invalid request body, pagination errors
- `401 Unauthorized` - Missing or invalid API key
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server or database errors

**See**: [Error Handling Documentation](./filters/06-api-reference.md#error-handling) for complete error reference

---

## Best Practices

### 1. Authentication
- ✅ Store API keys in environment variables
- ✅ Use HTTPS in production
- ✅ Rotate keys regularly
- ❌ Never commit API keys to version control

### 2. Filtering
- ✅ Use `email_status: "verified"` for contacts
- ✅ Use `seniority` filters to target decision-makers
- ✅ Use denormalized company fields (`company_*`) for single-query efficiency
- ✅ Combine multiple filter types for precise targeting

### 3. Performance
- ✅ Use `search_after` for large result sets
- ✅ Use `select_columns` to limit returned fields
- ✅ Keep `limit` reasonable (25-50 for best performance)
- ✅ Use stored filters when available

### 4. Pagination
- ✅ Use `search_after` for pagination beyond page 10
- ✅ Validate `page` (max 10) and `limit` (max 100) before requests
- ✅ Use count endpoints when you only need totals

### 5. Error Handling
- ✅ Always check `success` field in responses
- ✅ Implement exponential backoff for rate limit errors
- ✅ Validate filter keys using `/filters` endpoint
- ✅ Handle empty results gracefully

**See**: [Best Practices Guide](./filters/README.md#best-practices) for more details

---

## API Endpoints

### Company Endpoints

#### Read Operations
- `POST /companies` - Search/filter companies
- `POST /companies/count` - Get count of matching companies
- `GET /companies/filters` - Get available filter fields
- `POST /companies/filters/data` - Get filter data values

#### Write Operations
- `POST /companies/create` - Create a new company
- `PUT /companies/:uuid` - Update company by UUID
- `DELETE /companies/:uuid` - Delete company by UUID
- `POST /companies/upsert` - Create or update company
- `POST /companies/bulk` - Bulk upsert companies

### Contact Endpoints

#### Read Operations
- `POST /contacts` - Search/filter contacts
- `POST /contacts/count` - Get count of matching contacts
- `GET /contacts/filters` - Get available filter fields
- `POST /contacts/filters/data` - Get filter data values

#### Write Operations
- `POST /contacts/create` - Create a new contact
- `PUT /contacts/:uuid` - Update contact by UUID
- `DELETE /contacts/:uuid` - Delete contact by UUID
- `POST /contacts/upsert` - Create or update contact
- `POST /contacts/bulk` - Bulk upsert contacts

### System Endpoints

- `GET /health` - Health check (no authentication required)

### Jobs Endpoints

- `POST /common/jobs/create` - Create a new background job
- `POST /common/jobs` - List jobs with filters

**See**: [Jobs API Guide](./filters/jobs.md) for complete documentation

**See**: [API Reference](./filters/06-api-reference.md) for complete endpoint documentation

---

## Filter Types

### 1. Text Matches

Full-text search on text fields with three search types:

- **`exact`**: Phrase matching with word order
- **`shuffle`**: Word matching (order doesn't matter)
- **`substring`**: Partial text matching using ngram analysis

**Supports**: Fuzzy matching for typo tolerance

**See**: [VQL Syntax Reference](./filters/01-company-filters-complete-guide.md#comprehensive-vql-syntax-reference)

### 2. Keyword Matches

Exact matching on keyword/array fields:

- Single values or arrays
- Faster than text searches
- Use for categorical data

**See**: [Keyword Match Filters](./filters/01-company-filters-complete-guide.md#keyword-match-filters)

### 3. Range Queries

Numeric and date range filtering:

- Operators: `gte`, `lte`, `gt`, `lt`
- Dates in ISO 8601 format
- Efficient for numeric filtering

**See**: [Range Query Filters](./filters/01-company-filters-complete-guide.md#range-query-filters)

---

## Field Types

### Filterable Fields

Fields indexed in Elasticsearch that can be used in filter queries:

- **Text Fields**: Full-text searchable (name, address, title, etc.)
- **Keyword Fields**: Exact matching (industries, technologies, email_status, etc.)
- **Range Fields**: Numeric/date ranges (employees_count, annual_revenue, created_at, etc.)

### Response-Only Fields

Fields stored in PostgreSQL but not indexed in Elasticsearch:

- Cannot be used in `where` clauses
- Available in API responses
- Can be selected using `select_columns`

**Examples**: `facebook_url`, `twitter_url`, `phone_number`, `stage`, `work_direct_phone`

**See**: [Field Reference](./filters/04-filter-field-reference.md) for complete field list

---

## Denormalized Company Fields & Company Config

The contact index includes denormalized company data with the `company_*` prefix, enabling efficient account-based filtering:

- **For Filtering**: Use denormalized `company_*` fields (e.g., `company_name`, `company_industries`) in `where` clauses
- **For Responses**: Use `company_config.populate: true` with `company_config.select_columns` to get full company objects in responses

> **⚠️ Important**: Denormalized `company_*` fields are **ONLY for filtering**. They are **NOT available** in `select_columns`. To get company data in responses, use `company_config.select_columns` with direct field names (no prefix).

**Available Denormalized Filter Fields**:
- `company_name`, `company_industries`, `company_technologies`
- `company_employees_count`, `company_annual_revenue`, `company_total_funding`
- `company_city`, `company_state`, `company_country`
- And more (13 total fields for filtering)

**Company Config**: 27 company fields available via `company_config.select_columns`

**See**: 
- [Denormalized Company Fields](./filters/02-contact-filters-complete-guide.md#denormalized-company-fields)
- [Company Config](./filters/02-contact-filters-complete-guide.md#populating-company-data-in-responses-company_config)
- [Select Columns Guide](./filters/select_columns_filter.md) for complete documentation

---

## Examples

### Basic Company Search

**Lambda Deployment** (Production):
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "industries": ["Software"],
          "country": ["USA"]
        }
      },
      "range_query": {
        "must": {
          "employees_count": {
            "gte": 50,
            "lte": 500
          }
        }
      }
    },
    "page": 1,
    "limit": 25
  }'
```

**Local Development**:
```bash
curl -X POST http://localhost:8000/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "industries": ["Software"],
          "country": ["USA"]
        }
      },
      "range_query": {
        "must": {
          "employees_count": {
            "gte": 50,
            "lte": 500
          }
        }
      }
    },
    "page": 1,
    "limit": 25
  }'
```

### Contact Search with Company Filters

**Lambda Deployment** (Production):
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "company_industries": ["Software", "SaaS"],
          "seniority": ["Senior", "Lead", "Executive"],
          "email_status": "verified"
        }
      },
      "range_query": {
        "must": {
          "company_employees_count": {
            "gte": 100,
            "lte": 1000
          }
        }
      }
    },
    "page": 1,
    "limit": 50
  }'
```

**Local Development**:
```bash
curl -X POST http://localhost:8000/contacts \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "company_industries": ["Software", "SaaS"],
          "seniority": ["Senior", "Lead", "Executive"],
          "email_status": "verified"
        }
      },
      "range_query": {
        "must": {
          "company_employees_count": {
            "gte": 100,
            "lte": 1000
          }
        }
      }
    },
    "page": 1,
    "limit": 50
  }'
```

**See**: [Examples and Use Cases](./filters/05-examples-use-cases.md) for more examples

---

## Support and Resources

### Documentation

- **System**: [System Documentation](./system.md)
- **Filters**: [Filter Documentation Index](./filters/README.md)
- **API Reference**: [API Reference](./filters/06-api-reference.md)
- **Examples**: [Examples and Use Cases](./filters/05-examples-use-cases.md)

### Troubleshooting

- **Common Errors**: [Error Handling Guide](./filters/06-api-reference.md#error-handling)
- **Performance**: [Performance Considerations](./system.md#performance-considerations)
- **Connection Issues**: [Troubleshooting Guide](./system.md#troubleshooting)

### Getting Help

1. Review the relevant documentation file
2. Check the [Examples](./filters/05-examples-use-cases.md) for similar use cases
3. Review the [API Reference](./filters/06-api-reference.md) for endpoint details
4. Check [Error Handling](./filters/06-api-reference.md#error-handling) for error solutions

---

## Version History

### Version 1.2 (2025-12-24)

**Recent Updates**:
- ✅ **Write Operations**: Added full CRUD operations (create, update, delete, upsert, bulk) for contacts and companies
- ✅ **Automatic Elasticsearch Indexing**: All write operations automatically sync to search index
- ✅ **Validation Layer**: Comprehensive request validation with clear error messages
- ✅ **Bulk Operations**: Efficient bulk upsert endpoints for data imports
- Enhanced authentication and security documentation
- Comprehensive VQL syntax reference
- Complete error handling guide with solutions
- Enhanced HTTP request examples with curl commands
- Improved cross-references and navigation
- Standardized documentation structure

**Previous Updates**:
- Fixed `direct_derived` spelling throughout codebase
- Enhanced `select_columns` documentation
- Added comprehensive denormalized company fields documentation
- Documented all response-only fields
- Updated ngram configuration with accurate values

---

## Related Resources

- **Filter Documentation**: [Filter Index](./filters/README.md)
- **System Architecture**: [System Documentation](./system.md)
- **Company API**: [Company Documentation](./company.md)
- **Contact API**: [Contact Documentation](./contacts.md)

---

**Last Updated**: 2025-01-XX  
**Documentation Version**: 1.2

