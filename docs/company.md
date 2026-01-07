# Company API Documentation

## Overview

The Company API provides endpoints for searching, filtering, and managing company data. It uses a dual-storage architecture with PostgreSQL for primary data and Elasticsearch for fast search capabilities.

**All endpoints require authentication using an API Key via the `X-API-Key` header and are subject to rate limiting.**

## Base URL

**Lambda Deployment** (Production):

```
https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies
```

**Local Development**:

```
http://localhost:8000/companies
```

**Note**: The Lambda URL above is the production deployment. For local development, use `http://localhost:8000`.

## Endpoints

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
      "must": [
        {
          "text_value": "software",
          "filter_key": "name",
          "search_type": "exact",
          "slop": 3
        },
        {
          "text_value": "technology",
          "filter_key": "address",
          "search_type": "shuffle",
          "fuzzy": false
        }
      ],
      "must_not": [
        {
          "text_value": "consulting",
          "filter_key": "name",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "country": ["USA", "Canada"],
        "industries": ["Software", "Technology"]
      },
      "must_not": {}
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50,
          "lte": 1000
        },
        "annual_revenue": {
          "gte": 1000000
        },
        "created_at": {
          "gte": "2023-01-01T00:00:00Z",
          "lte": "2024-12-31T23:59:59Z"
        }
      }
    }
  },
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
  "page": 1,
  "limit": 25
}
```

**Query Parameters Explained**:

#### Text Matches

- **text_value**: The text to search for
- **filter_key**: Field name to search in (e.g., "name", "address", "city")
- **search_type**:
  - `"exact"`: Phrase matching with slop (word order matters)
  - `"shuffle"`: Word matching (order doesn't matter)
- **slop**: Number of words that can be between terms (for exact search)
- **fuzzy**: Enable fuzzy matching for typos (true/false)
- **operator**: "and" or "or" for multiple terms (for shuffle search)

#### Keyword Match

- **must**: Fields that must match exactly (supports arrays for multiple values)
- **must_not**: Fields that must not match

#### Range Query

- **gte**: Greater than or equal
- **lte**: Less than or equal
- **gt**: Greater than
- **lt**: Less than

#### Pagination & Sorting

- **page**: Page number (1-indexed, max: 10)
- **limit**: Results per page (max: 100, default: 25)
- **search_after**: Cursor-based pagination values (from previous response)
- **order_by**: Array of sort criteria

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

**Error Responses**:

401 Unauthorized:

```json
{
  "success": false,
  "error": "Unauthorized: Invalid or missing API key",
  "error_code": "UNAUTHORIZED",
  "message": "Unauthorized: Invalid or missing API key"
}
```

429 Too Many Requests:

```json
{
  "success": false,
  "error": "Rate limit exceeded. Please try again later.",
  "error_code": "RATE_LIMIT_EXCEEDED",
  "message": "Rate limit exceeded. Please try again later."
}
```

400 Bad Request:

```json
{
  "success": false,
  "error": "Validation error message",
  "error_code": "VALIDATION_ERROR",
  "message": "Validation error message",
  "details": {
    "field_name": "error details"
  }
}
```

400 Bad Request (Pagination):

```json
{
  "success": false,
  "error": "ERR_PAGE_SIZE_EXCEEDED: the requested page size surpasses the maximum allowed limit; consider using pagination with smaller batches",
  "error_code": "VALIDATION_ERROR",
  "message": "ERR_PAGE_SIZE_EXCEEDED: the requested page size surpasses the maximum allowed limit; consider using pagination with smaller batches"
}
```

500 Internal Server Error:

```json
{
  "success": false,
  "error": "An unexpected error occurred",
  "error_code": "INTERNAL_ERROR",
  "message": "An unexpected error occurred",
  "details": {
    "internal_error": "detailed error message"
  }
}
```

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
      "must": [
        {
          "text_value": "software",
          "filter_key": "name",
          "search_type": "exact",
          "slop": 3
        }
      ]
    },
    "keyword_match": {
      "must": {
        "country": ["USA"]
      }
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50
        }
      }
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
    },
    {
      "id": 11,
      "key": "country",
      "service": "company",
      "display_name": "Country",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 2,
      "key": "employees_count",
      "service": "company",
      "display_name": "Employees Count",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 3,
      "key": "industries",
      "service": "company",
      "display_name": "Industries",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 4,
      "key": "keywords",
      "service": "company",
      "display_name": "Keywords",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 12,
      "key": "linkedin_url",
      "service": "company",
      "display_name": "LinkedIn URL",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 14,
      "key": "normalized_domain",
      "service": "company",
      "display_name": "Normalized Domain",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 10,
      "key": "state",
      "service": "company",
      "display_name": "State",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 8,
      "key": "technologies",
      "service": "company",
      "display_name": "Technologies",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 7,
      "key": "total_funding",
      "service": "company",
      "display_name": "Total Funding",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 1,
      "key": "uuid",
      "service": "company",
      "display_name": "Name",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 13,
      "key": "website",
      "service": "company",
      "display_name": "Website",
      "direct_derived": true,
      "deleted_at": null
    }
  ],
  "success": true
}
```

**Filter Properties**:

- **id**: Unique identifier for the filter
- **key**: Filter identifier used in queries (matches field names in Elasticsearch)
- **service**: Service name (always `"company"` for company filters)
- **display_name**: Human-readable name for UI display
- **direct_derived**:
  - `true`: Filter values are extracted directly from company records in PostgreSQL
  - `false`: Filter values are stored in `filters_data` table for faster access
- **deleted_at**: Soft delete timestamp (null if active)

**Available Company Filters** (14 total):

- **Direct-Derived** (`direct_derived: true`): `address`, `annual_revenue`, `employees_count`, `linkedin_url`, `normalized_domain`, `total_funding`, `website`
- **Stored** (`direct_derived: false`): `city`, `country`, `industries`, `keywords`, `state`, `technologies`, `uuid` (displayed as "Name")

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

- **service**: Must be `"company"`
- **filter_key**: The filter key from the filters list
- **search_text**: (Optional) Text to filter results (case-insensitive, partial match)
- **page**: (Optional) Page number (default: 1)
- **limit**: (Optional) Results per page (max: 100, default: 25)

**Response** (200 OK):

```json
{
  "data": [
    "USA",
    "United Kingdom",
    "United Arab Emirates"
  ],
  "success": true
}
```

**How Filter Data Works**:

1. **Direct-Derived Filters** (`direct_derived: true`):
   - Values are extracted directly from the `companies` table
   - Searches the actual field values
   - Examples: `address`, `annual_revenue`, `employees_count`, `linkedin_url`, `normalized_domain`, `total_funding`, `website`

2. **Stored Filters** (`direct_derived: false`):
   - Values are pre-computed and stored in `filters_data` table
   - Faster for frequently used filters with many distinct values
   - Examples: `city`, `country`, `industries`, `keywords`, `state`, `technologies`, `uuid` (displayed as "Name")

---

## Write Operations

> **New Feature**: Write operations (create, update, delete, upsert, bulk) are now available for company data management.

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
  "uuid": "c0a8012e-1111-2222-3333-444455556666",
  "name": "Acme Software Corp",
  "normalized_domain": "acme.com",
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
  "facebook_url": "https://facebook.com/acme",
  "twitter_url": "https://twitter.com/acme",
  "company_name_for_emails": "Acme Corp",
  "phone_number": "+1-555-0123",
  "latest_funding": "Series B",
  "latest_funding_amount": 5000000,
  "last_raised_at": "2024-01-15"
}
```

**Validation Rules**:

- `name`: Required
- `uuid`: Optional (generated if not provided), must be valid UUID format
- `normalized_domain`: Optional, must be valid FQDN
- `employees_count`, `annual_revenue`, `total_funding`, `latest_funding_amount`: Must be >= 0
- `linkedin_url`, `website`, `facebook_url`, `twitter_url`: Must be valid URL format

**Response** (201 Created):

```json
{
  "data": {
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
    "created_at": "2025-12-24T10:30:00Z",
    "updated_at": "2025-12-24T10:30:00Z",
    "deleted_at": null
  },
  "success": true
}
```

**Automatic Indexing**: The company is automatically indexed in Elasticsearch for search capabilities.

### 6. Update Company

Update an existing company record by UUID with automatic Elasticsearch reindexing.

**Endpoint**: `PUT /companies/:uuid`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body** (all fields optional, provide only fields to update):

```json
{
  "name": "Acme Software Corporation",
  "employees_count": 150,
  "industries": ["Software", "Technology", "SaaS"],
  "annual_revenue": 7500000,
  "city": "San Francisco",
  "state": "CA"
}
```

**Response** (200 OK):

```json
{
  "data": {
    "id": 1,
    "uuid": "c0a8012e-1111-2222-3333-444455556666",
    "name": "Acme Software Corporation",
    "employees_count": 150,
    "industries": ["Software", "Technology", "SaaS"],
    "annual_revenue": 7500000,
    "city": "San Francisco",
    "state": "CA",
    "updated_at": "2025-12-24T11:45:00Z",
    ...
  },
  "success": true
}
```

**Error Response** (404 Not Found):

```json
{
  "success": false,
  "error": "Company with identifier 'specified UUID' not found",
  "error_code": "NOT_FOUND",
  "message": "Company with identifier 'specified UUID' not found"
}
```

### 7. Get Company by UUID

Retrieve a single company by its UUID with optional field selection.

**Endpoint**: `GET /companies/:uuid`

**Request Headers**:

```
X-API-Key: your-secret-api-key
```

**Query Parameters**:

- `select_columns` (Optional): Comma-separated list of field names to return. If not provided, all fields are returned.

**Example Request**:

```
GET /companies/c0a8012e-1111-2222-3333-444455556666?select_columns=uuid,name,employees_count,industries,website
```

**Response** (200 OK):

```json
{
  "data": {
    "id": 1,
    "uuid": "c0a8012e-1111-2222-3333-444455556666",
    "name": "Acme Software Corp",
    "employees_count": 120,
    "industries": ["Software", "Technology"],
    "website": "https://acme.com",
    "created_at": "2025-12-24T10:30:00Z",
    "updated_at": "2025-12-24T10:30:00Z",
    "deleted_at": null
  },
  "success": true
}
```

**Error Response** (404 Not Found):

```json
{
  "success": false,
  "error": "Company with identifier 'specified UUID' not found",
  "error_code": "NOT_FOUND",
  "message": "Company with identifier 'specified UUID' not found"
}
```

**Error Response** (400 Bad Request - Invalid UUID):

```json
{
  "success": false,
  "error": "Invalid UUID format",
  "error_code": "VALIDATION_ERROR",
  "message": "Invalid UUID format"
}
```

### 8. Delete Company

Soft delete a company record by UUID (sets `deleted_at` timestamp).

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

**Error Response** (404 Not Found):

```json
{
  "success": false,
  "error": "Company with identifier 'specified UUID' not found",
  "error_code": "NOT_FOUND",
  "message": "Company with identifier 'specified UUID' not found"
}
```

**Note**: Deleted companies are removed from Elasticsearch index but retained in PostgreSQL with `deleted_at` timestamp for audit purposes.

### 9. Upsert Company

Create a new company or update an existing one (identified by UUID or normalized_domain).

**Endpoint**: `POST /companies/upsert`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "uuid": "c0a8012e-1111-2222-3333-444455556666",
  "name": "Acme Software Corp",
  "normalized_domain": "acme.com",
  "employees_count": 120,
  "industries": ["Software", "Technology"],
  "annual_revenue": 5000000
}
```

**Response** (200 OK):

```json
{
  "data": {
    "id": 1,
    "uuid": "c0a8012e-1111-2222-3333-444455556666",
    "name": "Acme Software Corp",
    "employees_count": 120,
    ...
  },
  "success": true
}
```

**Upsert Logic**:

1. If company with matching UUID exists → update
2. If company with matching normalized_domain exists → update
3. Otherwise → create new company

**Response** (201 Created - New Company):

```json
{
  "data": {
    "id": 1,
    "uuid": "c0a8012e-1111-2222-3333-444455556666",
    "name": "Acme Software Corp",
    ...
  },
  "is_new": true,
  "success": true
}
```

**Response** (200 OK - Updated Company):

```json
{
  "data": {
    "id": 1,
    "uuid": "c0a8012e-1111-2222-3333-444455556666",
    "name": "Acme Software Corp",
    ...
  },
  "is_new": false,
  "success": true
}
```

### 10. Bulk Upsert Companies

Efficiently create or update multiple companies in a single request.

**Endpoint**: `POST /companies/bulk`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "companies": [
    {
      "uuid": "c0a8012e-1111-2222-3333-444455556666",
      "name": "Acme Software Corp",
      "employees_count": 120,
      "industries": ["Software"]
    },
    {
      "uuid": "c0a8012e-2222-3333-4444-555566667777",
      "name": "TechCorp Industries",
      "employees_count": 250,
      "industries": ["Technology"]
    }
  ]
}
```

**Validation**:

- `companies`: Required array with minimum 1 item
- Each company must have `name` field
- UUIDs generated automatically if not provided

**Response** (200 OK):

```json
{
  "data": {
    "total_count": 2,
    "success_count": 2,
    "error_count": 0,
    "created": 1,
    "updated": 1,
    "errors": []
  },
  "success": true
}
```

**Response** (207 Partial Content - Some Errors):

```json
{
  "data": {
    "total_count": 2,
    "success_count": 1,
    "error_count": 1,
    "created": 0,
    "updated": 1,
    "errors": [
      {
        "index": 1,
        "error": "Validation error: name is required"
      }
    ]
  },
  "success": false
}
```

**Performance**:

- Uses PostgreSQL `ON CONFLICT` for atomic upsert
- Batch Elasticsearch indexing for efficiency
- Can handle 1000+ records per request

**Best Practices for Bulk Operations**:

- Batch size: 100-500 records recommended
- Validate data before sending
- Monitor response times for large batches
- Use for data imports and synchronization
- For very large datasets (multi-GB CSV files), consider using [Jobs API](../filters/jobs.md) for asynchronous processing
- Jobs API supports streaming CSV import/export with automatic retry mechanisms

---

## Company Data Model

### PostgreSQL Schema

```go
type PgCompany struct {
    ID              uint64    `json:"id"`
    UUID            string    `json:"uuid"`
    Name            string    `json:"name"`
    EmployeesCount  int64     `json:"employees_count"`
    Industries      []string  `json:"industries"`
    Keywords        []string  `json:"keywords"`
    Address         string    `json:"address"`
    AnnualRevenue   int64     `json:"annual_revenue"`
    TotalFunding    int64     `json:"total_funding"`
    Technologies    []string  `json:"technologies"`
    City            string    `json:"city"`
    State           string    `json:"state"`
    Country         string    `json:"country"`
    LinkedinURL     string    `json:"linkedin_url"`
    Website         string    `json:"website"`
    NormalizedDomain string   `json:"normalized_domain"`
    FacebookURL     string    `json:"facebook_url"`
    TwitterURL      string    `json:"twitter_url"`
    CompanyNameForEmails string `json:"company_name_for_emails"`
    PhoneNumber     string    `json:"phone_number"`
    LatestFunding   string    `json:"latest_funding"`
    LatestFundingAmount int64 `json:"latest_funding_amount"`
    LastRaisedAt    string    `json:"last_raised_at"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    DeletedAt       *time.Time `json:"deleted_at"`
}
```

### Searchable Fields

The following fields are indexed in Elasticsearch and can be used in queries:

**Text Fields** (supports n-gram search):

- `name`
- `address`
- `city`
- `state`
- `country`
- `linkedin_url`
- `website`
- `normalized_domain`

**Keyword Fields** (exact matching):

- `id`
- `industries`
- `keywords`
- `technologies`

**Numeric Fields** (range queries):

- `employees_count`
- `annual_revenue`
- `total_funding`

**Date Fields** (range queries):

- `created_at`

## Query Examples

### Example 1: Simple Text Search

Find companies with "software" in the name:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "software",
          "filter_key": "name",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    }
  },
  "page": 1,
  "limit": 25
}
```

### Example 2: Multiple Countries

Find companies in specific countries:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "country": ["USA", "Canada", "UK"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Example 3: Size and Revenue Filter

Find mid-size companies with good revenue:

```json
{
  "where": {
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50,
          "lte": 500
        },
        "annual_revenue": {
          "gte": 1000000
        }
      }
    }
  },
  "order_by": [
    {
      "order_by": "annual_revenue",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 25
}
```

### Example 4: Complex Multi-Criteria Search

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "artificial intelligence",
          "filter_key": "name",
          "search_type": "shuffle",
          "fuzzy": true,
          "operator": "and"
        }
      ]
    },
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology"],
        "country": ["USA"]
      }
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 100
        },
        "created_at": {
          "gte": "2020-01-01T00:00:00Z"
        }
      }
    }
  },
  "order_by": [
    {
      "order_by": "employees_count",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 50
}
```

## Best Practices

1. **Use Appropriate Search Types**:
   - Use `exact` for phrases where word order matters
   - Use `shuffle` for general text search
   - Enable `fuzzy` for user-generated search terms

2. **Pagination**:
   - Use `search_after` for large result sets (more efficient than `page`)
   - Keep `limit` reasonable (25-50 for best performance)

3. **Filter Order**:
   - Place most selective filters first
   - Use `keyword_match` for exact values (faster)
   - Use `range_query` for numeric/date filters

4. **Performance**:
   - Use count endpoint when you only need the total
   - Avoid overly complex text searches
   - Use stored filters (`direct_derived: false`) when available

## Error Codes

The API uses standardized error codes for consistent error handling:

- `VALIDATION_ERROR` (400): Request validation failed (invalid format, missing required fields, etc.)
- `NOT_FOUND` (404): Resource not found (company, contact, etc.)
- `UNAUTHORIZED` (401): Missing or invalid API key
- `FORBIDDEN` (403): Insufficient permissions
- `RATE_LIMIT_EXCEEDED` (429): Too many requests, rate limit exceeded
- `BAD_REQUEST` (400): General bad request error
- `CONFLICT` (409): Resource conflict (e.g., duplicate entry)
- `INTERNAL_ERROR` (500): Internal server error (database, Elasticsearch, etc.)
- `ELASTICSEARCH_ERROR` (500): Elasticsearch operation failed
- `DATABASE_ERROR` (500): Database operation failed

All error responses follow this format:

```json
{
  "success": false,
  "error": "Error message",
  "error_code": "ERROR_CODE",
  "message": "Error message",
  "details": {}
}
```

## Authentication

All Company API endpoints require authentication using an API Key:

**Header Required**:

```
X-API-Key: your-secret-api-key
```

Configure the API key in your `.env` file:

```env
API_KEY=your-secret-api-key
```

## Rate Limiting

The API implements a token bucket rate limiter:

- **Algorithm**: Token bucket with continuous refill
- **Configuration**: `MAX_REQUESTS_PER_MINUTE` environment variable (default: 60)
- **Behavior**: Requests exceeding the limit receive a `429 Too Many Requests` response
- **Tokens**: Refill continuously at rate of `MAX_REQUESTS_PER_MINUTE / 60` per second

**Rate Limit Configuration**:

```env
MAX_REQUESTS_PER_MINUTE=60
```

**Recommendations**:

- Production: Set to appropriate value based on your infrastructure capacity
- Development: Use higher values for testing (e.g., 300-600)
- Enterprise: Consider per-client rate limiting with API key tracking
