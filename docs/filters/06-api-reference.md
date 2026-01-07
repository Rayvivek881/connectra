# API Reference - Filter Endpoints

## Table of Contents

1. [Overview](#overview)
2. [Company Filter Endpoints](#company-filter-endpoints)
   - [Read Operations](#read-operations)
   - [Write Operations](#write-operations)
3. [Contact Filter Endpoints](#contact-filter-endpoints)
   - [Read Operations](#read-operations-1)
   - [Write Operations](#write-operations-1)
4. [Jobs API Endpoints](#jobs-api-endpoints)
5. [Request/Response Formats](#requestresponse-formats)
6. [Error Handling](#error-handling)
7. [Rate Limiting](#rate-limiting)
8. [Authentication](#authentication)

## Overview

This document provides a complete API reference for all filter-related endpoints in the Connectra API. All endpoints use JSON for request and response bodies.

### Base URL

**Lambda Deployment** (Production):
```
https://iarj32v8e1.execute-api.us-east-1.amazonaws.com
```

**Local Development**:
```
http://localhost:8000
```

**Note**: The Lambda URL above is the production deployment. For local development, use `http://localhost:8000`.

### Authentication

All endpoints (except `/health`) require an API Key via the `X-API-Key` header.

### Rate Limiting

Token bucket algorithm with configurable requests per minute (default: 60). In Lambda deployment, rate limiting may also be handled by API Gateway.

### Deployment

Connectra can be deployed as:
- **Lambda Function** (Recommended): Serverless deployment on AWS Lambda. See [Lambda Deployment Guide](../LAMBDA_DEPLOYMENT.md).
- **Traditional Server**: Run as HTTP server on port 8000. See [System Documentation](../system.md#api-server).

---

## Company Filter Endpoints

### 1. List Companies by Filters

Search and retrieve companies based on complex filter criteria using VQL (Vivek Query Language).

**Endpoint**: `POST /companies`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "where": {
    "text_matches": {
      "must": [ ... ],
      "must_not": [ ... ]
    },
    "keyword_match": {
      "must": { ... },
      "must_not": { ... }
    },
    "range_query": {
      "must": { ... }
    }
  },
  "order_by": [ ... ],
  "page": <integer>,
  "limit": <integer>,
  "search_after": [ ... ],
  "select_columns": [ ... ]
}
```

**Response** (200 OK):

```json
{
  "data": [
    {
      "id": 1,
      "uuid": "c0a8012e-1111-2222-3333-444455556666",
      "name": "Acme Software Corp",
      "employees_count": 120,
      "industries": ["Software", "Technology"],
      "keywords": ["AI", "Machine Learning"],
      "address": "123 Tech Street",
      "annual_revenue": 5000000,
      "total_funding": 10000000,
      "technologies": ["Python", "Go", "React"],
      "city": "New York",
      "state": "NY",
      "country": "USA",
      "linkedin_url": "https://linkedin.com/company/acme",
      "website": "https://acme.com",
      "normalized_domain": "acme.com",
      "facebook_url": "https://facebook.com/acme",
      "twitter_url": "https://twitter.com/acme",
      "company_name_for_emails": "Acme Corp",
      "phone_number": "+1-555-0123",
      "latest_funding": "Series B",
      "latest_funding_amount": 5000000,
      "last_raised_at": "2024-01-15",
      "created_at": "2025-12-10T11:20:05.605686Z",
      "updated_at": "2025-12-10T11:20:05.605689Z",
      "deleted_at": null
    }
  ],
  "success": true
}
```

**Pagination Limits**:

- `page`: Maximum 10 (1-indexed)
- `limit`: Maximum 100, default 25

### 2. Count Companies by Filters

Get the total count of companies matching the filter criteria without retrieving the actual records.

**Endpoint**: `POST /companies/count`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "where": {
    "text_matches": {
      "must": [ ... ]
    },
    "keyword_match": {
      "must": { ... }
    },
    "range_query": {
      "must": { ... }
    }
  },
  "limit": 25
}
```

**Response** (200 OK):

```json
{
  "count": 123,
  "success": true
}
```

**Note**: The `limit` parameter is optional for count queries but may be used for query validation.

### 3. Get Company Filters

Retrieve all available filters for the company service.

**Endpoint**: `GET /companies/filters`

**Request Headers**:

```
X-API-Key: your-secret-api-key
```

**Response** (200 OK):

```json
{
  "data": [
    {
      "id": 5,
      "key": "address",
      "service": "company",
      "display_name": "Address",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 6,
      "key": "annual_revenue",
      "service": "company",
      "display_name": "Annual Revenue",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 9,
      "key": "city",
      "service": "company",
      "display_name": "City",
      "direct_derived": false,
      "deleted_at": null
    }
  ],
  "success": true
}
```

**Filter Properties**:

- `id`: Unique identifier for the filter
- `key`: Filter identifier used in queries (matches field names in Elasticsearch)
- `service`: Service name (always `"company"` for company filters)
- `display_name`: Human-readable name for UI display
- `direct_derived`:
  - `true`: Filter values are extracted directly from company records in PostgreSQL
  - `false`: Filter values are stored in `filters_data` table for faster access
- `deleted_at`: Soft delete timestamp (null if active)

### 4. Get Company Filter Data

Retrieve available values for a specific filter with optional search and pagination.

**Endpoint**: `POST /companies/filters/data`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "service": "company",
  "filter_key": "country",
  "search_text": "Uni",
  "page": 1,
  "limit": 25
}
```

**Request Parameters**:

- `service` (required): Must be `"company"`
- `filter_key` (required): The filter key from the filters list
- `search_text` (optional): Text to filter results (case-insensitive, partial match)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Results per page (max: 100, default: 25)

**Response** (200 OK):

```json
{
  "data": [
    {
      "value": "USA",
      "display_value": "USA"
    },
    {
      "value": "United Kingdom",
      "display_value": "United Kingdom"
    },
    {
      "value": "United Arab Emirates",
      "display_value": "United Arab Emirates"
    }
  ],
  "success": true
}
```

**Response Fields**:

- `value`: The actual filter value (used in queries)
- `display_value`: The human-readable display value (may differ from value for stored filters)

**How Filter Data Works**:

1. **Direct-Derived Filters** (`direct_derived: true`):
   - Values are extracted directly from the `companies` table
   - Searches the actual field values
   - Examples: `address`, `annual_revenue`, `employees_count`, `linkedin_url`, `normalized_domain`, `total_funding`, `website`

2. **Stored Filters** (`direct_derived: false`):
   - Values are pre-computed and stored in `filters_data` table
   - Faster for frequently used filters with many distinct values
   - Examples: `city`, `country`, `industries`, `keywords`, `state`, `technologies`, `uuid` (displayed as "Name")

### 5. Create Company

Create a new company record with automatic Elasticsearch indexing.

**Endpoint**: `POST /companies/create`

**Request Headers**:
```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:
```json
{
  "name": "Acme Software Corp",
  "normalized_domain": "acme.com",
  "employees_count": 120,
  "industries": ["Software", "Technology"],
  "annual_revenue": 5000000
}
```

**Response** (201 Created):
```json
{
  "data": {
    "id": 1,
    "uuid": "c0a8012e-1111-2222-3333-444455556666",
    "name": "Acme Software Corp",
    ...
  },
  "success": true
}
```

**See**: [Company API - Create Company](../company.md#5-create-company) for complete documentation.

### 6. Update Company

Update an existing company by UUID.

**Endpoint**: `PUT /companies/:uuid`

**Request Headers**:
```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Response** (200 OK): Returns updated company data

**See**: [Company API - Update Company](../company.md#6-update-company) for complete documentation.

### 7. Delete Company

Soft delete a company by UUID.

**Endpoint**: `DELETE /companies/:uuid`

**Request Headers**:
```
X-API-Key: your-secret-api-key
```

**Response** (200 OK):
```json
{
  "message": "Company deleted successfully",
  "success": true
}
```

**See**: [Company API - Delete Company](../company.md#7-delete-company) for complete documentation.

### 8. Upsert Company

Create or update a company (identified by UUID or normalized_domain).

**Endpoint**: `POST /companies/upsert`

**See**: [Company API - Upsert Company](../company.md#8-upsert-company) for complete documentation.

### 9. Bulk Upsert Companies

Efficiently create or update multiple companies.

> **Status**: This endpoint is currently implemented as `POST /companies/batch-upsert`. See [CRUD Implementation Plan](../CRUD_IMPLEMENTATION_PLAN.md) for details.

**Endpoint**: `POST /companies/batch-upsert` (Currently Implemented)

**Request Body**:
```json
{
  "companies": [
    {
      "name": "Company 1",
      "employees_count": 100
    },
    {
      "name": "Company 2",
      "employees_count": 200
    }
  ]
}
```

**Response** (200 OK):
```json
{
  "count": 2,
  "success": true
}
```

**See**: [Company API - Bulk Upsert Companies](../company.md#9-bulk-upsert-companies) for complete documentation.

---

## Contact Filter Endpoints

### 1. List Contacts by Filters

Search and retrieve contacts based on complex filter criteria using VQL (Vivek Query Language).

**Endpoint**: `POST /contacts`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "where": {
    "text_matches": {
      "must": [ ... ],
      "must_not": [ ... ]
    },
    "keyword_match": {
      "must": { ... },
      "must_not": { ... }
    },
    "range_query": {
      "must": { ... }
    }
  },
  "order_by": [ ... ],
  "page": <integer>,
  "limit": <integer>,
  "search_after": [ ... ],
  "select_columns": [ ... ],
  "company_config": {
    "populate": <boolean>,
    "select_columns": [ ... ]
  }
}
```

**Request Parameters**:

- All VQL query parameters (see [Contact Filters Guide](./02-contact-filters-complete-guide.md#filter-structure))
- `select_columns` (optional): Array of contact field names to return. Only specified fields will be included in the response.
- `company_config` (optional): Configuration for populating company data in responses
  - `populate` (boolean, required when company_config is used): Set to `true` to include company objects in the response
  - `select_columns` (array of strings, optional): List of company fields to return. **Use direct field names** (e.g., `name`, `employees_count`), **NOT** `company_*` prefix

> **⚠️ Important**: 
> - Denormalized `company_*` fields (e.g., `company_name`, `company_industries`) are **ONLY for filtering** in `where` clauses. They are **NOT available** in `select_columns`.
> - To get company data in responses, use `company_config.populate: true` with `company_config.select_columns` containing direct field names (no prefix).
> - See [Select Columns Guide](./select_columns_filter.md) for complete documentation.

**Response** (200 OK):

```json
{
  "data": [
    {
      "id": 43171040,
      "uuid": "021c8c87-1a5b-55a7-86c8-8f6f4710924e",
      "first_name": "Shawna",
      "last_name": "Hegmann",
      "company_id": "0231e33b-acc7-5bd7-9290-421e73a41358",
      "email": "shawna.hegmann@multiconsulting.digital",
      "title": "Senior Software Engineer",
      "departments": ["Sales", "Customer Success", "Support"],
      "mobile_phone": "4706037761",
      "email_status": "verified",
      "seniority": "Junior",
      "city": "Adelaide",
      "state": "VIC",
      "country": "Australia",
      "linkedin_url": "https://linkedin.com/in/shawna-hegmann-575684",
      "facebook_url": "https://facebook.com/shawna.hegmann",
      "twitter_url": "https://twitter.com/shawnahegmann",
      "website": "https://shawnahegmann.com",
      "work_direct_phone": "5861275199",
      "home_phone": "5931318115",
      "other_phone": "2451009971",
      "stage": "Closed Won",
      "created_at": "2025-12-10T11:20:05.605686Z",
      "updated_at": "2025-12-10T11:20:05.605689Z",
      "deleted_at": null
    }
  ],
  "success": true
}
```

**Note**: 
- Response includes both filterable fields (indexed in Elasticsearch) and response-only fields (PostgreSQL only). Fields like `facebook_url`, `twitter_url`, `website`, `work_direct_phone`, `home_phone`, `other_phone`, and `stage` are response-only and cannot be used in filter queries.
- Use `select_columns` to limit which contact fields are returned.
- If `company_config.populate: true` is set, the response will include a nested `company` object with the selected company fields. Example response with company_config:

```json
{
  "data": [
    {
      "id": 123,
      "first_name": "John",
      "last_name": "Doe",
      "company_id": "company-uuid-here",
      "company": {
        "uuid": "company-uuid-here",
        "name": "Acme Software Corp",
        "employees_count": 150,
        "industries": ["Software", "SaaS"]
      }
    }
  ],
  "success": true
}
```

**Pagination Limits**:

- `page`: Maximum 10 (1-indexed)
- `limit`: Maximum 100, default 25

### 2. Count Contacts by Filters

Get the total count of contacts matching the filter criteria without retrieving the actual records.

**Endpoint**: `POST /contacts/count`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "where": {
    "text_matches": {
      "must": [ ... ]
    },
    "keyword_match": {
      "must": { ... }
    },
    "range_query": {
      "must": { ... }
    }
  },
  "limit": 25
}
```

**Response** (200 OK):

```json
{
  "count": 456,
  "success": true
}
```

**Note**: The `limit` parameter is optional for count queries but may be used for query validation.

### 3. Get Contact Filters

Retrieve all available filters for the contact service.

**Endpoint**: `GET /contacts/filters`

**Request Headers**:

```
X-API-Key: your-secret-api-key
```

**Response** (200 OK):

```json
{
  "data": [
    {
      "id": 24,
      "key": "city",
      "service": "contact",
      "display_name": "City",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 17,
      "key": "company_id",
      "service": "contact",
      "display_name": "Company ID",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 26,
      "key": "country",
      "service": "contact",
      "display_name": "Country",
      "direct_derived": false,
      "deleted_at": null
    }
  ],
  "success": true
}
```

**Filter Properties**:

- `id`: Unique identifier for the filter
- `key`: Filter identifier used in queries (matches field names in Elasticsearch)
- `service`: Service name (always `"contact"` for contact filters)
- `display_name`: Human-readable name for UI display
- `direct_derived`:
  - `true`: Filter values are extracted directly from contact records in PostgreSQL
  - `false`: Filter values are stored in `filters_data` table for faster access
- `deleted_at`: Soft delete timestamp (null if active)

### 4. Get Contact Filter Data

Retrieve available values for a specific filter with optional search and pagination.

**Endpoint**: `POST /contacts/filters/data`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "service": "contact",
  "filter_key": "departments",
  "search_text": "Eng",
  "page": 1,
  "limit": 25
}
```

**Request Parameters**:

- `service` (required): Must be `"contact"`
- `filter_key` (required): The filter key from the filters list
- `search_text` (optional): Text to filter results (case-insensitive, partial match)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Results per page (max: 100, default: 25)

**Response** (200 OK):

```json
{
  "data": [
    {
      "value": "Engineering",
      "display_value": "Engineering"
    },
    {
      "value": "Sales",
      "display_value": "Sales"
    },
    {
      "value": "Customer Success",
      "display_value": "Customer Success"
    },
    {
      "value": "Support",
      "display_value": "Support"
    },
    {
      "value": "HR",
      "display_value": "HR"
    },
    {
      "value": "Marketing",
      "display_value": "Marketing"
    },
    {
      "value": "Operations",
      "display_value": "Operations"
    },
    {
      "value": "Legal",
      "display_value": "Legal"
    },
    {
      "value": "Finance",
      "display_value": "Finance"
    },
    {
      "value": "Product",
      "display_value": "Product"
    }
  ],
  "success": true
}
```

**How Filter Data Works**:

1. **Direct-Derived Filters** (`direct_derived: true`):
   - Values are extracted directly from the `contacts` table
   - Searches the actual field values
   - Examples: `company_id`, `email`, `first_name`, `last_name`, `linkedin_url`, `mobile_phone`

2. **Stored Filters** (`direct_derived: false`):
   - Values are pre-computed and stored in `filters_data` table
   - Faster for frequently used filters
   - Examples: `city`, `country`, `departments`, `email_status`, `seniority`, `state`, `title`

**Note**: Contact index also includes denormalized company fields with `company_` prefix (e.g., `company_name`, `company_industries`, `company_employees_count`). These allow filtering contacts directly by company attributes. See [Contact Filters Guide](./02-contact-filters-complete-guide.md#denormalized-company-fields) for details.

### 5. Create Contact

Create a new contact record with automatic Elasticsearch indexing.

**Endpoint**: `POST /contacts/create`

**Request Headers**:
```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:
```json
{
  "first_name": "John",
  "last_name": "Smith",
  "email": "john.smith@example.com",
  "company_id": "c0a8012e-1111-2222-3333-444455556666",
  "title": "Senior Software Engineer",
  "departments": ["Engineering"],
  "seniority": "Senior"
}
```

**Response** (201 Created): Returns created contact data

**See**: [Contact API - Create Contact](../contacts.md#5-create-contact) for complete documentation.

### 6. Update Contact

Update an existing contact by UUID.

**Endpoint**: `PUT /contacts/:uuid`

**Request Headers**:
```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Response** (200 OK): Returns updated contact data

**See**: [Contact API - Update Contact](../contacts.md#6-update-contact) for complete documentation.

### 7. Delete Contact

Soft delete a contact by UUID.

**Endpoint**: `DELETE /contacts/:uuid`

**Request Headers**:
```
X-API-Key: your-secret-api-key
```

**Response** (200 OK):
```json
{
  "message": "Contact deleted successfully",
  "success": true
}
```

**See**: [Contact API - Delete Contact](../contacts.md#7-delete-contact) for complete documentation.

### 8. Upsert Contact

Create or update a contact (identified by UUID or email).

**Endpoint**: `POST /contacts/upsert`

**See**: [Contact API - Upsert Contact](../contacts.md#8-upsert-contact) for complete documentation.

### 9. Bulk Upsert Contacts

Efficiently create or update multiple contacts.

> **Status**: This endpoint is currently implemented as `POST /contacts/batch-upsert`. See [CRUD Implementation Plan](../CRUD_IMPLEMENTATION_PLAN.md) for details.

**Endpoint**: `POST /contacts/batch-upsert` (Currently Implemented)

**Request Body**:
```json
{
  "contacts": [
    {
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com"
    },
    {
      "first_name": "Jane",
      "last_name": "Doe",
      "email": "jane.doe@example.com"
    }
  ]
}
```

**Response** (200 OK):
```json
{
  "count": 2,
  "success": true
}
```

**See**: [Contact API - Bulk Upsert Contacts](../contacts.md#9-bulk-upsert-contacts) for complete documentation.

---

## Jobs API Endpoints

### 1. Create Job

Create a new background job for CSV import or export.

**Endpoint**: `POST /common/jobs/create`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "job_type": "insert_csv_file",
  "job_data": {
    "s3_key": "uploads/contacts.csv",
    "s3_bucket": "my-bucket"
  },
  "retry_count": 3
}
```

**Job Types**:

- `insert_csv_file`: Import CSV from S3
  - Required fields: `s3_key`, `s3_bucket` (optional, defaults to configured bucket)
- `export_csv_file`: Export data to S3 as CSV
  - Required fields: `service` ("contact" or "company"), `vql` (with `select_columns`), `s3_bucket` (optional)

**Response** (201 Created):

```json
{
  "message": "Job created successfully",
  "success": true
}
```

**Error Responses**:

- `400 Bad Request`: Missing required fields, invalid job type
- `401 Unauthorized`: Invalid API key
- `500 Internal Server Error`: Job creation failed

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

### 2. List Jobs

List jobs with optional filters.

**Endpoint**: `POST /common/jobs`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "job_type": "insert_csv_file",
  "status": ["open", "processing", "completed"],
  "limit": 25
}
```

**Request Parameters**:

- `job_type` (optional): Filter by job type
- `status` (optional): Array of statuses to filter by
- `limit` (optional): Max results (default: 25, max: 100)

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

**See**: [Jobs API Guide](./jobs.md) for complete documentation

---

## Request/Response Formats

### VQL Query Structure

The `where` clause in request bodies follows the VQL (Vivek Query Language) structure:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "search text",
          "filter_key": "field_name",
          "search_type": "exact" | "shuffle" | "substring",
          "slop": <integer>,
          "operator": "and" | "or",
          "fuzzy": true | false
        }
      ],
      "must_not": [ ... ]
    },
    "keyword_match": {
      "must": {
        "field_name": "value" | ["value1", "value2"],
        "array_field": ["value1", "value2"]
      },
      "must_not": { ... }
    },
    "range_query": {
      "must": {
        "numeric_field": {
          "gte": <number>,
          "lte": <number>,
          "gt": <number>,
          "lt": <number>
        },
        "date_field": {
          "gte": "2024-01-01T00:00:00Z",
          "lte": "2024-12-31T23:59:59Z"
        }
      }
    }
  },
  "order_by": [
    {
      "order_by": "field_name",
      "order_direction": "asc" | "desc"
    }
  ],
  "page": <integer>,
  "limit": <integer>,
  "search_after": ["value1", "value2"],
  "select_columns": ["field1", "field2", ...]
}
```

**Query Parameters**:

- `where`: Filter conditions (text_matches, keyword_match, range_query)
- `order_by`: Sorting configuration (array of order objects)
- `page`: Page number for pagination (1-indexed, max: 10)
- `limit`: Results per page (max: 100, default: 25)
- `search_after`: Cursor-based pagination values from previous response
- `select_columns`: Optional array of field names to return from PostgreSQL (limits which fields are fetched after Elasticsearch search)

### Sorting

**Sortable Fields**:

- **Companies**: `id`, `employees_count`, `annual_revenue`, `total_funding`, `created_at`, `industries`, `keywords`, `technologies`
- **Contacts**: `id`, `company_id`, `email`, `departments`, `mobile_phone`, `email_status`, `seniority`, `created_at`

**Non-Sortable Fields** (text fields):

- **Companies**: `name`, `address`, `city`, `state`, `country`, `linkedin_url`, `website`, `normalized_domain`
- **Contacts**: `first_name`, `last_name`, `title`, `city`, `state`, `country`, `linkedin_url`
- **Denormalized Company Fields in Contacts**: All `company_*` fields are not sortable (they are text fields or not indexed for sorting)

**Filterable vs Response-Only Fields**:

- **Filterable Fields**: Fields indexed in Elasticsearch that can be used in `where` clauses
- **Response-Only Fields**: Fields stored in PostgreSQL but not indexed (e.g., `facebook_url`, `twitter_url`, `phone_number`, `stage`, `work_direct_phone`, `home_phone`, `other_phone`, `updated_at`, `deleted_at`). These can be selected using `select_columns` but **cannot be used in filters**. Attempting to filter by these fields will result in an error.

### Pagination

**Page-Based Pagination**:

```json
{
  "page": 1,
  "limit": 25
}
```

**Cursor-Based Pagination**:

```json
{
  "order_by": [
    {
      "order_by": "annual_revenue",
      "order_direction": "desc"
    },
    {
      "order_by": "created_at",
      "order_direction": "asc"
    }
  ],
  "search_after": ["5000000", "2024-01-15T08:00:00Z"],
  "limit": 25
}
```

**Pagination Limits**:

- `page`: Maximum 10 (1-indexed)
- `limit`: Maximum 100, default 25
- `search_after`: Values from last document in previous response

### Field Selection (select_columns)

The `select_columns` parameter allows you to limit which fields are returned from PostgreSQL after the Elasticsearch search. This can improve performance by reducing the amount of data transferred.

**Important**: `select_columns` only affects PostgreSQL field retrieval after the Elasticsearch search completes. It does NOT affect:

- Elasticsearch search performance
- Which documents are matched by the search
- Filter query execution

**How it works**:

1. Elasticsearch executes the filter query and returns matching document IDs
2. PostgreSQL retrieves full records for those IDs
3. If `select_columns` is specified, only those fields are returned from PostgreSQL
4. If omitted, all fields (both filterable and response-only) are returned

**Example: Select specific filterable fields**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software"]
      }
    }
  },
  "select_columns": ["id", "name", "employees_count", "annual_revenue"],
  "page": 1,
  "limit": 25
}
```

**Example: Select response-only fields**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software"]
      }
    }
  },
  "select_columns": ["id", "name", "facebook_url", "twitter_url", "phone_number"],
  "page": 1,
  "limit": 25
}
```

**Example: Select denormalized company fields (contacts)**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software"],
        "seniority": ["Senior"]
      }
    }
  },
  "select_columns": ["id", "first_name", "last_name", "email", "company_name", "company_annual_revenue"],
  "page": 1,
  "limit": 25
}
```

**Notes**:

- `select_columns` is optional - if omitted, all fields are returned
- Only affects PostgreSQL field selection, not Elasticsearch search
- Useful for reducing response payload size
- Field names should match PostgreSQL column names
- Can include both filterable fields (from Elasticsearch) and response-only fields (PostgreSQL only)
- Response-only fields (e.g., `facebook_url`, `twitter_url`, `phone_number`, `stage`, `work_direct_phone`) cannot be used in `where` clauses but can be selected in responses
- Denormalized company fields (e.g., `company_name`, `company_annual_revenue`) can be selected but are not stored in PostgreSQL - they come from the company record via `company_id`

---

## Error Handling

### Error Response Format

All error responses follow this format:

```json
{
  "error": "error message",
  "success": false
}
```

**Note**: Some errors may include additional fields like `message` for more context.

### Complete Error Code Reference

#### 400 Bad Request

**Invalid Request Body**:

```json
{
  "error": "ERR_INVALID_REQUEST_BODY: the request body is invalid; check JSON syntax and required fields",
  "success": false
}
```

**Common Causes**:

- Malformed JSON in request body
- Missing required fields
- Invalid field types (e.g., string instead of integer)
- Invalid filter structure

**Example - Invalid JSON**:

```bash
# Request with invalid JSON
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-key" \
  -d '{"where": {invalid json}'
# Response: 400 Bad Request - "invalid request body"
```

**Example - Missing Required Field**:

```json
// Request missing 'where' clause
{
  "page": 1,
  "limit": 25
}
// Response: 400 Bad Request - "invalid request body"
```

**Page Size Exceeded**:

```json
{
  "error": "ERR_PAGE_SIZE_EXCEEDED: the requested page size surpasses the maximum allowed limit; consider using pagination with smaller batches",
  "success": false
}
```

**Causes**:

- `limit` parameter exceeds maximum value (100)
- `limit` is negative or zero

**Example**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software"]
      }
    }
  },
  "page": 1,
  "limit": 150  // ❌ ERROR: Maximum is 100
}
```

**Solution**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software"]
      }
    }
  },
  "page": 1,
  "limit": 100  // ✅ Correct: Maximum allowed
}
```

**Page Number Exceeded**:

```json
{
  "error": "ERR_PAGE_OUT_OF_RANGE: the requested page number is beyond the available range; verify total pages before requesting",
  "success": false
}
```

**Causes**:

- `page` parameter exceeds maximum value (10)
- `page` is zero or negative

**Example**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software"]
      }
    }
  },
  "page": 15,  // ❌ ERROR: Maximum is 10
  "limit": 25
}
```

**Solution**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software"]
      }
    }
  },
  "page": 10,  // ✅ Correct: Maximum allowed
  "limit": 25
}
```

**Note**: For pagination beyond page 10, use `search_after` cursor-based pagination instead.

**Invalid Filter Field**:

```json
{
  "error": "ERR_ELASTICSEARCH_FAILURE: search engine returned status 400; details: ...",
  "success": false
}
```

**Common Causes**:

- Using non-existent field names in filters
- Using response-only fields in `where` clauses
- Invalid field type for filter operation (e.g., using text field in range query)

**Example - Invalid Field Name**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "invalid_field": ["value"]  // ❌ ERROR: Field doesn't exist
      }
    }
  }
}
```

**Solution**: Use `/filters` endpoint to get valid field names:

```bash
GET /companies/filters
# Returns list of valid filter keys
```

**Example - Using Response-Only Field in Filter**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "facebook_url": "https://facebook.com/company"  // ❌ ERROR: Response-only field
      }
    }
  }
}
```

**Solution**: Filter by filterable fields, then select response-only fields:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software"]  // ✅ Filter by filterable field
      }
    }
  },
  "select_columns": ["id", "name", "facebook_url"]  // ✅ Select response-only field
}
```

#### 401 Unauthorized

**Missing or Invalid API Key**:

```json
{
  "error": "unauthorized",
  "message": "invalid API key"
}
```

**HTTP Status**: `401 Unauthorized`

**Causes**:

- Missing `X-API-Key` header
- Incorrect API key value
- API key doesn't match configured value
- Empty API key value

**Example - Missing Header**:

```bash
# Request without X-API-Key header
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -d '{"where": {}}'
# Response: 401 Unauthorized
```

**Solution**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{"where": {}}'
```

**Example - Invalid API Key**:

```bash
# Request with wrong API key
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: wrong-key" \
  -d '{"where": {}}'
# Response: 401 Unauthorized
```

**Solution**: Verify API key matches server configuration:

```env
# .env file on server
API_KEY=your-secret-api-key
```

#### 429 Too Many Requests

**Rate Limit Exceeded**:

```json
{
  "error": "rate limit exceeded",
  "message": "too many requests, please try again later"
}
```

**HTTP Status**: `429 Too Many Requests`

**Causes**:

- Exceeded `MAX_REQUESTS_PER_MINUTE` configured limit
- Token bucket depleted (wait for tokens to refill)

**Example - Rate Limit Exceeded**:

```bash
# Making too many requests in short time
for i in {1..100}; do
  curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
    -H "X-API-Key: your-key" \
    -d '{"where": {}}'
done
# After limit exceeded: 429 Too Many Requests
```

**Solution - Exponential Backoff**:

```javascript
async function makeRequestWithRetry(url, options, maxRetries = 5) {
  let delay = 1000; // Start with 1 second
  
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    const response = await fetch(url, options);
    
    if (response.status === 429) {
      // Rate limit exceeded - wait and retry
      await sleep(delay);
      delay = Math.min(delay * 2, 30000); // Exponential backoff, cap at 30s
      continue;
    }
    
    return response;
  }
  
  throw new Error('Max retries exceeded');
}
```

#### 500 Internal Server Error

**Elasticsearch Error**:

```json
{
  "error": "ERR_ELASTICSEARCH_FAILURE: search engine returned status 400; details: ...",
  "success": false
}
```

**Common Causes**:

- Elasticsearch cluster unavailable
- Invalid Elasticsearch query syntax
- Index not found or corrupted
- Field mapping mismatch
- Query too complex or resource-intensive

**Example - Elasticsearch Connection Error**:

```json
{
  "error": "ERR_ELASTICSEARCH_FAILURE: search engine connection error; details: connection refused",
  "success": false
}
```

**Troubleshooting Steps**:

1. Check Elasticsearch cluster health
2. Verify Elasticsearch connection settings in `.env`
3. Check Elasticsearch logs for detailed error
4. Verify index exists: `GET /companies_index`
5. Check field mappings match query structure

**Example - Invalid Query Syntax**:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "software",
          "filter_key": "name",
          "search_type": "invalid_type"  // ❌ ERROR: Invalid search_type
        }
      ]
    }
  }
}
```

**Solution**: Use valid `search_type` values: `"exact"`, `"shuffle"`, or `"substring"`

**Database Error**:

```json
{
  "error": "database error: ...",
  "success": false
}
```

**Common Causes**:

- PostgreSQL connection unavailable
- Database query timeout
- Invalid SQL query
- Database connection pool exhausted
- Transaction deadlock

**Example - Database Connection Error**:

```json
{
  "error": "database error: connection refused",
  "success": false
}
```

**Troubleshooting Steps**:

1. Check PostgreSQL server status
2. Verify database connection settings in `.env`
3. Check connection pool limits
4. Review PostgreSQL logs
5. Verify database exists and is accessible

**Example - Connection Pool Exhausted**:

```json
{
  "error": "database error: too many connections",
  "success": false
}
```

**Solution**: 

- Reduce concurrent requests
- Increase connection pool size in configuration
- Implement request queuing

### Validation Error Examples

**Invalid Filter Structure**:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "software",
          // ❌ ERROR: Missing required 'filter_key'
          "search_type": "shuffle"
        }
      ]
    }
  }
}
```

**Solution**:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "software",
          "filter_key": "name",  // ✅ Required field
          "search_type": "shuffle"
        }
      ]
    }
  }
}
```

**Invalid Range Query**:

```json
{
  "where": {
    "range_query": {
      "must": {
        "employees_count": {
          "gte": "fifty"  // ❌ ERROR: Must be integer, not string
        }
      }
    }
  }
}
```

**Solution**:

```json
{
  "where": {
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50  // ✅ Correct: Integer value
        }
      }
    }
  }
}
```

**Invalid Date Format**:

```json
{
  "where": {
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2024-01-01"  // ❌ ERROR: Missing time component
        }
      }
    }
  }
}
```

**Solution**:

```json
{
  "where": {
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2024-01-01T00:00:00Z"  // ✅ Correct: ISO 8601 format
        }
      }
    }
  }
}
```

### Error Handling Best Practices

1. **Always check the `success` field** in responses

   ```javascript
   const response = await fetch(url, options);
   const data = await response.json();
   
   if (!data.success) {
     console.error('API Error:', data.error);
     // Handle error appropriately
   }
   ```

2. **Handle authentication errors** by ensuring valid `X-API-Key` header is always sent

   ```javascript
   if (response.status === 401) {
     // Re-authenticate or show login prompt
     console.error('Authentication failed');
   }
   ```

3. **Implement retry logic** for rate limit errors (429) with exponential backoff

   ```javascript
   if (response.status === 429) {
     const retryAfter = response.headers.get('Retry-After') || 60;
     await sleep(retryAfter * 1000);
     return retryRequest();
   }
   ```

4. **Handle pagination errors** by validating `page` and `limit` before sending requests

   ```javascript
   function validatePagination(page, limit) {
     if (page < 1 || page > 10) {
       throw new Error('Page must be between 1 and 10');
     }
     if (limit < 1 || limit > 100) {
       throw new Error('Limit must be between 1 and 100');
     }
   }
   ```

5. **Validate filter keys** using the `/filters` endpoint before using them

   ```javascript
   async function getValidFilters(service) {
     const response = await fetch(`/${service}/filters`, {
       headers: { 'X-API-Key': apiKey }
     });
     const data = await response.json();
     return data.data.map(f => f.key);
   }
   ```

6. **Handle empty results** gracefully (empty `data` array with `success: true`)

   ```javascript
   if (data.success && data.data.length === 0) {
     console.log('No results found');
     // Show appropriate message to user
   }
   ```

7. **Log errors for debugging** while protecting sensitive information

   ```javascript
   if (!data.success) {
     logger.error('API Error', {
       error: data.error,
       endpoint: url,
       timestamp: new Date().toISOString()
       // Don't log API keys or sensitive data
     });
   }
   ```

8. **Implement circuit breaker pattern** for repeated failures

   ```javascript
   class CircuitBreaker {
     constructor(threshold = 5, timeout = 60000) {
       this.failures = 0;
       this.threshold = threshold;
       this.timeout = timeout;
       this.state = 'CLOSED'; // CLOSED, OPEN, HALF_OPEN
     }
     
     async call(fn) {
       if (this.state === 'OPEN') {
         throw new Error('Circuit breaker is OPEN');
       }
       
       try {
         const result = await fn();
         this.onSuccess();
         return result;
       } catch (error) {
         this.onFailure();
         throw error;
       }
     }
     
     onSuccess() {
       this.failures = 0;
       this.state = 'CLOSED';
     }
     
     onFailure() {
       this.failures++;
       if (this.failures >= this.threshold) {
         this.state = 'OPEN';
         setTimeout(() => {
           this.state = 'HALF_OPEN';
         }, this.timeout);
       }
     }
   }
   ```

### Troubleshooting Common Errors

**Problem**: Getting 500 errors for valid queries

- **Check**: Elasticsearch cluster health
- **Check**: Database connection status
- **Check**: Server logs for detailed error messages
- **Solution**: Verify all services are running and accessible

**Problem**: Getting 400 errors for filter queries

- **Check**: Field names match exactly (case-sensitive)
- **Check**: Field types match filter operation (text vs keyword vs range)
- **Check**: Required parameters are present
- **Solution**: Use `/filters` endpoint to verify valid field names

**Problem**: Rate limit errors in production

- **Check**: Current `MAX_REQUESTS_PER_MINUTE` setting
- **Check**: Request patterns and frequency
- **Solution**: Implement exponential backoff, reduce request frequency, or increase limit

**Problem**: Authentication errors after key rotation

- **Check**: API key matches server configuration
- **Check**: Header name is exactly `X-API-Key` (case-sensitive)
- **Solution**: Verify key in environment variables matches request header

---

## Rate Limiting

The API implements a token bucket rate limiter to prevent abuse and ensure fair usage.

### Rate Limiting Details

**Algorithm**: Token Bucket with continuous refill

- **Configuration**: `MAX_REQUESTS_PER_MINUTE` environment variable (default: 60)
- **Token Refill Rate**: `MAX_REQUESTS_PER_MINUTE / 60` tokens per second
- **Bucket Capacity**: Equal to `MAX_REQUESTS_PER_MINUTE`
- **Behavior**: When tokens are exhausted, requests receive `429 Too Many Requests`
- **Refill Mechanism**: Continuous refill at constant rate (not burst refill)

**Token Bucket Algorithm Formula**:

```
Token Refill Rate = MAX_REQUESTS_PER_MINUTE / 60 tokens/second
Bucket Capacity = MAX_REQUESTS_PER_MINUTE tokens
Available Tokens = min(Bucket Capacity, Previous Tokens + Refill Rate × Time Elapsed)
```

**Example Calculation**:

- If `MAX_REQUESTS_PER_MINUTE = 60`:
  - Refill rate: `60 / 60 = 1 token/second`
  - Bucket capacity: `60 tokens`
  - After 10 seconds of inactivity: `min(60, 0 + 1 × 10) = 10 tokens available`

**Configuration Example** (`.env` file):

```env
MAX_REQUESTS_PER_MINUTE=60
```

**Production Configuration Recommendations**:

```env
# Development/Testing
MAX_REQUESTS_PER_MINUTE=300

# Production (moderate traffic)
MAX_REQUESTS_PER_MINUTE=120

# Production (high traffic)
MAX_REQUESTS_PER_MINUTE=600

# Enterprise (very high traffic)
MAX_REQUESTS_PER_MINUTE=1200
```

### Rate Limit Response

When rate limit is exceeded:

```json
{
  "error": "rate limit exceeded",
  "message": "too many requests, please try again later"
}
```

**HTTP Status**: `429 Too Many Requests`

**Response Headers** (if available):

```
Retry-After: 60
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1640995200
```

### Best Practices for Rate Limiting

1. **Implement Exponential Backoff**: When receiving 429 errors, wait before retrying (e.g., 1s, 2s, 4s, 8s)

   ```javascript
   // Example exponential backoff implementation
   let delay = 1000; // Start with 1 second
   while (response.status === 429) {
     await sleep(delay);
     delay = Math.min(delay * 2, 30000); // Cap at 30 seconds
     response = await retryRequest();
   }
   ```

2. **Batch Requests**: Group multiple queries to reduce total request count
   - Use array filters (e.g., `industries: ["Software", "SaaS"]`) instead of multiple requests
   - Combine filters in single queries when possible

3. **Cache Results**: Store frequently accessed data to avoid repeated requests
   - Cache filter metadata (`/filters` endpoints)
   - Cache filter data values (`/filters/data` endpoints)
   - Cache frequently used query results

4. **Monitor Usage**: Track your API usage to stay within limits
   - Log rate limit errors
   - Track request patterns
   - Set up alerts for approaching limits

5. **Request Limit Increase**: Contact support if you need higher limits for production

**Recommendations**:

- Development: Use higher values (e.g., 300-600 requests/minute)
- Production: Set based on infrastructure capacity
- Enterprise: Consider per-API-key rate limiting

### Rate Limiting Bypass Scenarios

**Exempt Endpoints**:

- `/health` - Health check endpoint (no authentication or rate limiting)

**Note**: All other endpoints are subject to rate limiting. There is no way to bypass rate limits for authenticated endpoints.

### Exemptions

The `/health` endpoint is NOT rate limited or authenticated for monitoring purposes.

---

## Authentication

All API endpoints (except `/health`) require authentication using an API Key.

### Authentication Method

**Header Required**:

```
X-API-Key: your-secret-api-key
```

**Complete HTTP Request Example**:

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

**Configuration** (`.env` file):

```env
API_KEY=your-secret-api-key
```

**Note**: The API key must match exactly (case-sensitive) with the value configured in the server's environment.

### Authentication Response

**Unauthorized** (401):

```json
{
  "error": "unauthorized",
  "message": "invalid API key"
}
```

**HTTP Status**: `401 Unauthorized`

**Causes**:

- Missing `X-API-Key` header
- Incorrect API key value
- API key doesn't match configured value
- Empty API key value

**Example Error Scenarios**:

**Missing Header**:

```bash
# Request without X-API-Key header
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -d '{"where": {}}'
# Response: 401 Unauthorized
```

**Invalid API Key**:

```bash
# Request with wrong API key
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: wrong-key" \
  -d '{"where": {}}'
# Response: 401 Unauthorized
```

### API Key Security Best Practices

#### 1. Secure Storage

**✅ DO**:

- Store API keys in environment variables
- Use secret management services (AWS Secrets Manager, HashiCorp Vault, etc.)
- Use configuration files excluded from version control (`.env`, `.env.local`)
- Rotate keys regularly (every 90 days recommended)

**❌ DON'T**:

- Hardcode API keys in source code
- Commit API keys to version control (Git, SVN, etc.)
- Share API keys in chat/email/forums
- Store API keys in client-side code (browser JavaScript, mobile apps)

**Example - Secure Storage**:

```bash
# .env file (excluded from Git via .gitignore)
API_KEY=sk_live_abc123xyz789...

# .gitignore
.env
.env.local
.env.*.local
```

#### 2. HTTPS Only in Production

**✅ DO**:

- Always use HTTPS in production environments
- Enforce TLS 1.2 or higher
- Use valid SSL certificates
- Enable certificate pinning for mobile apps

**❌ DON'T**:

- Send API keys over unencrypted HTTP in production
- Use self-signed certificates in production
- Disable SSL verification in production code

**Example - HTTPS Request**:

```bash
curl -X POST https://api.example.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{"where": {}}'
```

#### 3. Key Rotation

**Best Practices**:

- Rotate API keys every 90 days
- Use key versioning to support gradual rotation
- Monitor for key expiration warnings
- Have a rollback plan if rotation fails

**Rotation Process**:

1. Generate new API key
2. Update configuration with new key
3. Deploy updated configuration
4. Verify new key works
5. Revoke old key after verification period

#### 4. Separate Keys for Different Environments

**✅ DO**:

- Use different API keys for development, staging, and production
- Use separate keys for different services/applications
- Track which key is used by which service

**Example Configuration**:

```env
# Development
API_KEY_DEV=sk_dev_abc123...

# Staging
API_KEY_STAGING=sk_staging_xyz789...

# Production
API_KEY_PROD=sk_prod_def456...
```

#### 5. Access Control and Monitoring

**✅ DO**:

- Monitor API key usage for suspicious activity
- Set up alerts for unusual access patterns
- Log all API requests with API key identifiers
- Implement IP whitelisting if possible
- Track request patterns per API key

**❌ DON'T**:

- Share API keys between multiple services without tracking
- Ignore unusual access patterns
- Use same key for multiple unrelated services

### Security Recommendations

**API Key Generation**:

- Use strong, randomly generated API keys
- Minimum length: 32 characters (recommended: 64+ characters)
- Use cryptographically secure random generators
- Include alphanumeric and special characters

**Example - Strong API Key**:

```
sk_live_.....
```

**Key Format Recommendations**:

- Prefix keys with environment identifier: `sk_dev_`, `sk_staging_`, `sk_prod_`
- Use base64 or hex encoding for readability
- Include checksum or validation characters

**Additional Security Layers**:

- **IP Whitelisting**: Restrict API access to specific IP addresses (if supported)
- **Request Signing**: Sign requests with HMAC for additional verification (future enhancement)
- **OAuth2**: Consider OAuth2 for enhanced security in future versions
- **API Key Scoping**: Limit API keys to specific endpoints/resources (future enhancement)
- **Rate Limiting per Key**: Track and limit requests per API key (future enhancement)

### Production Deployment Security Checklist

- [ ] API keys stored in environment variables (not in code)
- [ ] `.env` files excluded from version control
- [ ] HTTPS enabled for all API requests
- [ ] Valid SSL certificates configured
- [ ] API keys rotated regularly (every 90 days)
- [ ] Different keys for dev/staging/production
- [ ] API key usage monitoring enabled
- [ ] Alerts configured for suspicious activity
- [ ] Access logs reviewed regularly
- [ ] Key rotation process documented
- [ ] Incident response plan for key compromise

### CORS Configuration

**Current Configuration** (Development):

- **Allowed Origins**: All origins (`*`)
- **Allowed Methods**: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`, `PATCH`
- **Allowed Headers**: Common headers including `Content-Type`, `X-API-Key`
- **Credentials**: Not configured (no credentials support)

**Production Recommendations**:

- Restrict allowed origins to specific domains
- Limit allowed methods to required ones only
- Specify exact allowed headers
- Configure credentials if needed

**Example - Production CORS Configuration**:

```go
// Production CORS configuration (example)
corsConfig := cors.Config{
    AllowOrigins:     []string{"https://app.example.com", "https://admin.example.com"},
    AllowMethods:     []string{"GET", "POST"},
    AllowHeaders:     []string{"Content-Type", "X-API-Key"},
    ExposeHeaders:    []string{"X-RateLimit-Limit", "X-RateLimit-Remaining"},
    AllowCredentials: false,
    MaxAge:           3600,
}
```

**Security Considerations**:

- **Wildcard Origins**: Only use `*` in development. Never use in production.
- **Credentials**: If using credentials, specify exact origins (no wildcards)
- **Preflight Requests**: OPTIONS requests are handled automatically
- **Header Exposure**: Only expose necessary headers in `ExposeHeaders`

### Future Authentication Enhancements

**Planned Features**:

- **OAuth2 Support**: Token-based authentication with refresh tokens
- **JWT Tokens**: JSON Web Token support for stateless authentication
- **API Key Scoping**: Limit keys to specific endpoints or resources
- **Key Expiration**: Automatic expiration and renewal
- **Multi-Factor Authentication**: Additional security layer for sensitive operations
- **Role-Based Access Control (RBAC)**: Different permissions per API key
- **Audit Logging**: Comprehensive audit trail for all API operations

---

## Related Documentation

- [Company Filters Guide](./01-company-filters-complete-guide.md)
- [Contact Filters Guide](./02-contact-filters-complete-guide.md)
- [Combined Filters Guide](./03-combined-filters-guide.md)
- [Filter Field Reference](./04-filter-field-reference.md)
- [Examples and Use Cases](./05-examples-use-cases.md)

---

## API Endpoint Summary

### Company Endpoints

#### Read Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/companies` | List companies by filters |
| POST | `/companies/count` | Count companies by filters |
| GET | `/companies/filters` | Get available company filters |
| POST | `/companies/filters/data` | Get filter data values |

#### Write Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/companies/create` | Create a new company |
| PUT | `/companies/:uuid` | Update company by UUID |
| DELETE | `/companies/:uuid` | Delete company by UUID |
| POST | `/companies/upsert` | Create or update company |
| POST | `/companies/bulk` | Bulk upsert companies |

### Contact Endpoints

#### Read Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/contacts` | List contacts by filters |
| POST | `/contacts/count` | Count contacts by filters |
| GET | `/contacts/filters` | Get available contact filters |
| POST | `/contacts/filters/data` | Get filter data values |

#### Write Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/contacts/create` | Create a new contact |
| PUT | `/contacts/:uuid` | Update contact by UUID |
| DELETE | `/contacts/:uuid` | Delete contact by UUID |
| POST | `/contacts/upsert` | Create or update contact |
| POST | `/contacts/bulk` | Bulk upsert contacts |

### Jobs Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/common/jobs/create` | Create a new background job |
| POST | `/common/jobs` | List jobs with filters |

**See**: [Jobs API Guide](./jobs.md) for complete documentation

---

## Version Information

**API Version**: 1.2

**Last Updated**: 2025-12-24

**Base URL**: `https://iarj32v8e1.execute-api.us-east-1.amazonaws.com`

**Recent Updates**:
- ✅ **Write Operations** (2025-12-24): Added full CRUD operations for contacts and companies
- Enhanced authentication and security documentation
- Comprehensive error handling guide

**Content-Type**: `application/json`
