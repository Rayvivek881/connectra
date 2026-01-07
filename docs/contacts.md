# Contact API Documentation

## Overview

The Contact API provides endpoints for searching, filtering, and managing contact data. It uses a dual-storage architecture with PostgreSQL for primary data and Elasticsearch for fast search capabilities.

**All endpoints require authentication using an API Key via the `X-API-Key` header and are subject to rate limiting.**

## Base URL

**Lambda Deployment** (Production):

```
https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts
```

**Local Development**:

```
http://localhost:8000/contacts
```

**Note**: The Lambda URL above is the production deployment. For local development, use `http://localhost:8000`.

## Endpoints

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
      "must": [
        {
          "text_value": "software engineer",
          "filter_key": "title",
          "search_type": "exact",
          "slop": 3
        },
        {
          "text_value": "John",
          "filter_key": "first_name",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ],
      "must_not": [
        {
          "text_value": "intern",
          "filter_key": "title",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "country": ["USA", "Canada"],
        "departments": ["Engineering", "Research"],
        "email_status": "verified",
        "seniority": "Senior"
      },
      "must_not": {}
    },
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2023-01-01T00:00:00Z",
          "lte": "2024-12-31T23:59:59Z"
        }
      }
    }
  },
  "order_by": [
    {
      "order_by": "created_at",
      "order_direction": "desc"
    },
    {
      "order_by": "email",
      "order_direction": "asc"
    }
  ],
  "page": 1,
  "limit": 25,
  "select_columns": ["id", "first_name", "last_name", "email"],
  "company_config": {
    "populate": true,
    "select_columns": ["name", "employees_count", "industries"]
  }
}
```

**Query Parameters Explained**:

#### Text Matches

- **text_value**: The text to search for
- **filter_key**: Field name to search in (e.g., "first_name", "last_name", "title", "city")
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

#### Field Selection

- **select_columns** (optional): Array of contact field names to return. Only specified fields will be included in the response. Can include both filterable and response-only fields (e.g., `facebook_url`, `twitter_url`, `stage`).

#### Company Config (Contacts Only)

- **company_config** (optional): Configuration for populating company data in responses
  - **populate** (boolean, required when company_config is used): Set to `true` to include company objects in the response
  - **select_columns** (array of strings, optional): List of company fields to return. **Use direct field names** (e.g., `name`, `employees_count`), **NOT** `company_*` prefix

> **⚠️ Important**: 
> - Denormalized `company_*` fields (e.g., `company_name`, `company_industries`) are **ONLY for filtering** in `where` clauses. They are **NOT available** in `select_columns`.
> - To get company data in responses, use `company_config.populate: true` with `company_config.select_columns` containing direct field names (no prefix).
> - See [Populating Company Data](#populating-company-data-company_config) section below for complete documentation.

**Important - Sortable Fields**: Only certain fields can be used for sorting in Elasticsearch:

- **Sortable fields** (keyword or date): `id`, `company_id`, `email`, `departments`, `mobile_phone`, `email_status`, `seniority`, `created_at`
- **Non-sortable fields** (text fields): `first_name`, `last_name`, `title`, `city`, `state`, `country`, `linkedin_url` - These are analyzed text fields and cannot be used for sorting. Use keyword fields or date fields for sorting instead.

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
      "must": [
        {
          "text_value": "engineer",
          "filter_key": "title",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "country": ["USA"],
        "departments": ["Engineering"]
      }
    },
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2023-01-01T00:00:00Z"
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
    },
    {
      "id": 20,
      "key": "departments",
      "service": "contact",
      "display_name": "Departments",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 18,
      "key": "email",
      "service": "contact",
      "display_name": "Email",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 22,
      "key": "email_status",
      "service": "contact",
      "display_name": "Email Status",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 15,
      "key": "first_name",
      "service": "contact",
      "display_name": "First Name",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 16,
      "key": "last_name",
      "service": "contact",
      "display_name": "Last Name",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 27,
      "key": "linkedin_url",
      "service": "contact",
      "display_name": "LinkedIn URL",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 21,
      "key": "mobile_phone",
      "service": "contact",
      "display_name": "Mobile Phone",
      "direct_derived": true,
      "deleted_at": null
    },
    {
      "id": 23,
      "key": "seniority",
      "service": "contact",
      "display_name": "Seniority",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 25,
      "key": "state",
      "service": "contact",
      "display_name": "State",
      "direct_derived": false,
      "deleted_at": null
    },
    {
      "id": 19,
      "key": "title",
      "service": "contact",
      "display_name": "Title",
      "direct_derived": false,
      "deleted_at": null
    }
  ],
  "success": true
}
```

**Filter Properties**:

- **id**: Unique identifier for the filter
- **key**: Filter identifier used in queries (matches field names in Elasticsearch)
- **service**: Service name (always `"contact"` for contact filters)
- **display_name**: Human-readable name for UI display
- **direct_derived**:
  - `true`: Filter values are extracted directly from contact records in PostgreSQL
  - `false`: Filter values are stored in `filters_data` table for faster access
- **deleted_at**: Soft delete timestamp (null if active)

**Available Contact Filters** (13 total):

- **Direct-Derived** (`direct_derived: true`): `company_id`, `email`, `first_name`, `last_name`, `linkedin_url`, `mobile_phone`
- **Stored** (`direct_derived: false`): `city`, `country`, `departments`, `email_status`, `seniority`, `state`, `title`

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

- **service**: Must be `"contact"`
- **filter_key**: The filter key from the filters list
- **search_text**: (Optional) Text to filter results (case-insensitive, partial match)
- **page**: (Optional) Page number (default: 1)
- **limit**: (Optional) Results per page (max: 100, default: 25)

**Response** (200 OK):

```json
{
  "data": [
    "Engineering",
    "Sales",
    "Customer Success",
    "Support",
    "HR",
    "Marketing",
    "Operations",
    "Legal",
    "Finance",
    "Product"
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

---

## Write Operations

> **New Feature**: Write operations (create, update, delete, upsert, bulk) are now available for contact data management.

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
  "mobile_phone": "4706037761",
  "email_status": "verified",
  "seniority": "Senior",
  "city": "San Francisco",
  "state": "CA",
  "country": "USA",
  "linkedin_url": "https://linkedin.com/in/john-smith",
  "facebook_url": "https://facebook.com/johnsmith",
  "twitter_url": "https://twitter.com/johnsmith",
  "website": "https://johnsmith.dev"
}
```

**Validation Rules**:

- `first_name`: Required
- `last_name`: Required
- `email`: Required, must be valid email format
- `company_id`: Optional, must be valid UUID format
- `linkedin_url`, `facebook_url`, `twitter_url`, `website`: Must be valid URL format
- `uuid`: Optional (generated if not provided)

**Response** (201 Created):

```json
{
  "data": {
    "id": 43171040,
    "uuid": "d1e2f3a4-5678-90ab-cdef-1234567890ab",
    "first_name": "John",
    "last_name": "Smith",
    "email": "john.smith@example.com",
    "company_id": "c0a8012e-1111-2222-3333-444455556666",
    "title": "Senior Software Engineer",
    "departments": ["Engineering"],
    "mobile_phone": "4706037761",
    "email_status": "verified",
    "seniority": "Senior",
    "city": "San Francisco",
    "state": "CA",
    "country": "USA",
    "linkedin_url": "https://linkedin.com/in/john-smith",
    "facebook_url": "https://facebook.com/johnsmith",
    "twitter_url": "https://twitter.com/johnsmith",
    "website": "https://johnsmith.dev",
    "created_at": "2025-12-24T10:30:00Z",
    "updated_at": "2025-12-24T10:30:00Z",
    "deleted_at": null
  },
  "success": true
}
```

**Automatic Indexing**: The contact is automatically indexed in Elasticsearch for search capabilities.

### 6. Update Contact

Update an existing contact record by UUID with automatic Elasticsearch reindexing.

**Endpoint**: `PUT /contacts/:uuid`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body** (all fields optional, provide only fields to update):

```json
{
  "title": "Staff Software Engineer",
  "seniority": "Principal",
  "departments": ["Engineering", "Architecture"],
  "city": "New York",
  "state": "NY"
}
```

**Response** (200 OK):

```json
{
  "data": {
    "id": 43171040,
    "uuid": "d1e2f3a4-5678-90ab-cdef-1234567890ab",
    "first_name": "John",
    "last_name": "Smith",
    "title": "Staff Software Engineer",
    "seniority": "Principal",
    "departments": ["Engineering", "Architecture"],
    "city": "New York",
    "state": "NY",
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
  "error": "Contact with identifier 'specified UUID' not found",
  "error_code": "NOT_FOUND",
  "message": "Contact with identifier 'specified UUID' not found"
}
```

### 7. Get Contact by UUID

Retrieve a single contact by its UUID with optional field selection.

**Endpoint**: `GET /contacts/:uuid`

**Request Headers**:

```
X-API-Key: your-secret-api-key
```

**Query Parameters**:

- `select_columns` (Optional): Comma-separated list of field names to return. If not provided, all fields are returned.

**Example Request**:

```
GET /contacts/d1e2f3a4-5678-90ab-cdef-1234567890ab?select_columns=uuid,first_name,last_name,email,title,company_id
```

**Response** (200 OK):

```json
{
  "data": {
    "id": 43171040,
    "uuid": "d1e2f3a4-5678-90ab-cdef-1234567890ab",
    "first_name": "John",
    "last_name": "Smith",
    "email": "john.smith@example.com",
    "title": "Senior Software Engineer",
    "company_id": "c0a8012e-1111-2222-3333-444455556666",
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
  "error": "Contact with identifier 'specified UUID' not found",
  "error_code": "NOT_FOUND",
  "message": "Contact with identifier 'specified UUID' not found"
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

### 8. Delete Contact

Soft delete a contact record by UUID (sets `deleted_at` timestamp).

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

**Error Response** (404 Not Found):

```json
{
  "success": false,
  "error": "Contact with identifier 'specified UUID' not found",
  "error_code": "NOT_FOUND",
  "message": "Contact with identifier 'specified UUID' not found"
}
```

**Note**: Deleted contacts are removed from Elasticsearch index but retained in PostgreSQL with `deleted_at` timestamp for audit purposes.

### 9. Upsert Contact

Create a new contact or update an existing one (identified by UUID or email).

**Endpoint**: `POST /contacts/upsert`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "uuid": "d1e2f3a4-5678-90ab-cdef-1234567890ab",
  "first_name": "John",
  "last_name": "Smith",
  "email": "john.smith@example.com",
  "company_id": "c0a8012e-1111-2222-3333-444455556666",
  "title": "Senior Software Engineer",
  "departments": ["Engineering"],
  "seniority": "Senior"
}
```

**Response** (200 OK):

```json
{
  "data": {
    "id": 43171040,
    "uuid": "d1e2f3a4-5678-90ab-cdef-1234567890ab",
    "first_name": "John",
    "last_name": "Smith",
    "email": "john.smith@example.com",
    ...
  },
  "success": true
}
```

**Upsert Logic**:

1. If contact with matching UUID exists → update
2. If contact with matching email exists → update
3. Otherwise → create new contact

**Response** (201 Created - New Contact):

```json
{
  "data": {
    "id": 43171040,
    "uuid": "d1e2f3a4-5678-90ab-cdef-1234567890ab",
    "first_name": "John",
    "last_name": "Smith",
    ...
  },
  "is_new": true,
  "success": true
}
```

**Response** (200 OK - Updated Contact):

```json
{
  "data": {
    "id": 43171040,
    "uuid": "d1e2f3a4-5678-90ab-cdef-1234567890ab",
    "first_name": "John",
    "last_name": "Smith",
    ...
  },
  "is_new": false,
  "success": true
}
```

### 10. Bulk Upsert Contacts

Efficiently create or update multiple contacts in a single request.

**Endpoint**: `POST /contacts/bulk`

**Request Headers**:

```
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body**:

```json
{
  "contacts": [
    {
      "uuid": "d1e2f3a4-5678-90ab-cdef-1234567890ab",
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "company_id": "c0a8012e-1111-2222-3333-444455556666",
      "title": "Senior Software Engineer",
      "departments": ["Engineering"],
      "seniority": "Senior"
    },
    {
      "uuid": "e2f3g4h5-6789-01bc-def0-234567890bcd",
      "first_name": "Jane",
      "last_name": "Doe",
      "email": "jane.doe@example.com",
      "company_id": "c0a8012e-2222-3333-4444-555566667777",
      "title": "Engineering Manager",
      "departments": ["Engineering"],
      "seniority": "Lead"
    }
  ]
}
```

**Validation**:

- `contacts`: Required array with minimum 1 item
- Each contact must have `first_name`, `last_name`, and `email` fields
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
        "error": "Validation error: email is required"
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
- Validate data before sending (especially email formats)
- Monitor response times for large batches
- Use for data imports and synchronization
- Ensure company_id references exist for relational integrity
- For very large datasets (multi-GB CSV files), consider using [Jobs API](../filters/jobs.md) for asynchronous processing
- Jobs API supports streaming CSV import/export with automatic retry mechanisms

---

## Contact Data Model

### PostgreSQL Schema

```go
type PgContact struct {
    ID            uint64    `json:"id"`
    UUID          string    `json:"uuid"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    CompanyID     string    `json:"company_id"`
    Email         string    `json:"email"`
    Title         string    `json:"title"`
    Departments   []string  `json:"departments"`
    MobilePhone   string    `json:"mobile_phone"`
    EmailStatus   string    `json:"email_status"`
    Seniority     string    `json:"seniority"`
    City          string    `json:"city"`
    State         string    `json:"state"`
    Country       string    `json:"country"`
    LinkedinURL   string    `json:"linkedin_url"`
    FacebookURL   string    `json:"facebook_url"`
    TwitterURL    string    `json:"twitter_url"`
    Website       string    `json:"website"`
    WorkDirectPhone string  `json:"work_direct_phone"`
    HomePhone     string    `json:"home_phone"`
    OtherPhone    string    `json:"other_phone"`
    Stage         string    `json:"stage"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    DeletedAt     *time.Time `json:"deleted_at"`
}
```

### Searchable Fields

The following fields are indexed in Elasticsearch and can be used in queries:

**Text Fields** (supports n-gram search):

- `first_name`
- `last_name`
- `title`
- `city`
- `state`
- `country`
- `linkedin_url`

**Keyword Fields** (exact matching):

- `id`
- `company_id`
- `email`
- `departments`
- `mobile_phone`
- `email_status`
- `seniority`

**Date Fields** (range queries and sorting):

- `created_at`

**Note on Sorting**: Text fields (`first_name`, `last_name`, `title`, `city`, `state`, `country`, `linkedin_url`) cannot be used for sorting in Elasticsearch because they are analyzed fields. Use keyword fields (`email`, `company_id`, `seniority`, `email_status`, etc.) or date fields (`created_at`) for sorting operations.

### Field Value Enumerations

**email_status** (keyword field):

- `"verified"` - Email address has been verified
- `"unverified"` - Email address not yet verified
- `"invalid"` - Email address is invalid
- `"bounced"` - Email has bounced

**seniority** (keyword field):

- `"Junior"` - Junior level
- `"Mid"` - Mid-level
- `"Senior"` - Senior level
- `"Lead"` - Lead/Team Lead
- `"Principal"` - Principal level
- `"Executive"` - Executive level

**stage** (PostgreSQL field, not in Elasticsearch):

- `"Contacted"` - Initial contact made
- `"Qualified"` - Contact qualified
- `"Lead"` - Lead status
- `"Proposal"` - Proposal sent
- `"Negotiation"` - In negotiation
- `"Closed Won"` - Deal won
- `"Closed Lost"` - Deal lost

**departments** (keyword array field):
Common values include: `"Engineering"`, `"Sales"`, `"Customer Success"`, `"Support"`, `"HR"`, `"Marketing"`, `"Operations"`, `"Legal"`, `"Finance"`, `"Product"`, `"Research"`, etc.

## Query Examples

### Example 1: Search by Name

Find contacts with "John" in first or last name:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "John",
          "filter_key": "first_name",
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

### Example 2: Find Engineers in Specific Departments

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "engineer",
          "filter_key": "title",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "departments": ["Engineering", "Backend", "Frontend"]
      }
    }
  },
  "order_by": [
    {
      "order_by": "seniority",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 25
}
```

### Example 3: Verified Contacts in USA

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "country": ["USA"],
        "email_status": "verified"
      },
      "must_not": {}
    }
  },
  "page": 1,
  "limit": 50
}
```

### Example 4: Recent Senior Contacts

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": "Senior"
      }
    },
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2024-01-01T00:00:00Z"
        }
      }
    }
  },
  "order_by": [
    {
      "order_by": "created_at",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 25
}
```

### Example 5: Complex Multi-Criteria Search

Find senior software engineers in tech companies:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "software engineer",
          "filter_key": "title",
          "search_type": "exact",
          "slop": 2
        }
      ],
      "must_not": [
        {
          "text_value": "intern",
          "filter_key": "title",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "seniority": "Senior",
        "departments": ["Engineering"],
        "country": ["USA", "Canada"]
      },
      "must_not": {}
    },
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2022-01-01T00:00:00Z"
        }
      }
    }
  },
  "order_by": [
    {
      "order_by": "created_at",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 50
}
```

### Example 6: Search by Company

Find all contacts for specific companies:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "c0a8012e-1111-2222-3333-444455556666",
          "c0a8012e-2222-3333-4444-555566667777"
        ]
      }
    }
  },
  "order_by": [
    {
      "order_by": "email",
      "order_direction": "asc"
    }
  ],
  "page": 1,
  "limit": 100
}
```

## Best Practices

1. **Name Searches**:
   - Use `shuffle` search type for flexible name matching
   - Enable `fuzzy` for typo tolerance
   - Search both `first_name` and `last_name` separately for better results

2. **Title/Job Searches**:
   - Use `exact` with appropriate `slop` for job titles with multiple words
   - Example: "software engineer" with slop 2 allows "senior software engineer"

3. **Email Status Filtering**:
   - Use `email_status: "verified"` for high-quality contacts
   - Filter by `email_status` to ensure contactability

4. **Pagination**:
   - Use `search_after` for large result sets (more efficient than `page`)
   - Keep `limit` reasonable (25-50 for best performance)

5. **Performance**:
   - Use count endpoint when you only need the total
   - Combine multiple keyword filters in a single `must` object
   - Use stored filters (`direct_derived: false`) when available

6. **Data Quality**:
   - Filter by `email_status` to ensure contactability
   - Use `seniority` filter to target decision-makers
   - Combine with company filters for account-based approaches

## Common Use Cases

### Use Case 1: Email Campaign Targeting

Find verified, active contacts in specific countries:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "country": ["USA", "UK", "Canada"],
        "email_status": "verified"
      },
      "must_not": {}
    }
  },
  "page": 1,
  "limit": 100
}
```

### Use Case 2: Recruiting Search

Find senior engineers in specific technologies:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "python go javascript",
          "filter_key": "title",
          "search_type": "shuffle",
          "operator": "or"
        }
      ]
    },
    "keyword_match": {
      "must": {
        "seniority": "Senior",
        "departments": ["Engineering"]
      }
    }
  },
  "order_by": [
    {
      "order_by": "created_at",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 50
}
```

### Use Case 3: Account-Based Marketing

Find contacts at target companies:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "company-uuid-1",
          "company-uuid-2",
          "company-uuid-3"
        ]
      }
    }
  },
  "order_by": [
    {
      "order_by": "seniority",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 100
}
```

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

All Contact API endpoints require authentication using an API Key:

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

## Company Config - Populating Company Objects

The `company_config` feature allows you to populate **full company objects** from PostgreSQL alongside contact data in a single query. This is different from denormalized `company_*` fields:

- **Denormalized fields** (`company_*`): Already in contact index, can filter by them, limited to 13 fields
- **Company Config** (`company_config.populate`): Fetches full company objects separately, all 27 fields available

### How company_config Works

1. **Elasticsearch Search**: Filter contacts (can use denormalized `company_*` fields)
2. **Extract Company IDs**: Get `company_id` values from matched contacts
3. **Parallel Fetch**: Fetch full company records from PostgreSQL (in parallel with contacts)
4. **Attach to Response**: Company objects are attached to each contact in the response

### Company Config Structure

```json
{
  "company_config": {
    "populate": true,
    "select_columns": ["uuid", "name", "website", ...]
  }
}
```

- `populate`: `true` to enable company object population
- `select_columns`: Array of company field names to return (all 27 fields available)

### Complete Company Field Reference (27 fields)

#### Core Company Fields (17 fields)

1. `id` - Company ID (bigint, primary key)
2. `uuid` - Company UUID (text, unique)
3. `name` - Company name (text, ngram support 3-10)
4. `employees_count` - Employee count (bigint, range filterable)
5. `industries` - Industries array (text[], keyword filterable)
6. `keywords` - Keywords array (text[], keyword filterable)
7. `address` - Company address (text)
8. `annual_revenue` - Annual revenue in cents (bigint, range filterable)
9. `total_funding` - Total funding in cents (bigint, range filterable)
10. `technologies` - Technologies array (text[], keyword filterable)
11. `city` - City name (text)
12. `state` - State/Province (text)
13. `country` - Country name (text)
14. `linkedin_url` - LinkedIn URL (text)
15. `website` - Website URL (text)
16. `normalized_domain` - Normalized domain (text)
17. `created_at` - Creation date (timestamp, range filterable)

#### Company Metadata Fields (10 fields)

1. `facebook_url` - Facebook page URL (text)
2. `twitter_url` - Twitter profile URL (text)
3. `company_name_for_emails` - Company name formatted for emails (text)
4. `phone_number` - Company phone number (text)
5. `latest_funding` - Latest funding round (text, e.g., "Series B")
6. `latest_funding_amount` - Latest funding amount in cents (bigint)
7. `last_raised_at` - Date of last funding round (text)
8. `updated_at` - Last update timestamp (timestamp)
9. `deleted_at` - Soft delete timestamp (timestamp, null if active)
10. `linkedin_sales_url` - LinkedIn Sales Navigator URL (text)

### Example: Basic Company Population

**Request**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified",
        "seniority": ["Senior", "Lead"]
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",
      "employees_count",
      "industries",
      "annual_revenue",
      "website",
      "phone_number",
      "latest_funding"
    ]
  },
  "page": 1,
  "limit": 25
}
```

**Response**:

```json
{
  "data": [
    {
      "id": 43171040,
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "title": "Senior Software Engineer",
      "company_id": "c0a8012e-1111-2222-3333-444455556666",
      "company": {
        "uuid": "c0a8012e-1111-2222-3333-444455556666",
        "name": "Altitude Software",
        "employees_count": 500,
        "industries": ["Software", "SaaS"],
        "annual_revenue": 10000000,
        "website": "https://altitude.com",
        "phone_number": "+1-555-123-4567",
        "latest_funding": "Series B"
      }
    }
  ],
  "success": true
}
```

### Denormalized Fields vs Company Config

**Key Distinction**:

| Aspect | Denormalized Fields (`company_*`) | Company Config (`company_config.populate`) |
|--------|-----------------------------------|--------------------------------------------|
| **Usage** | **ONLY for filtering** in `where` clauses | **ONLY for selecting** in `company_config.select_columns` |
| **Field Names** | Use `company_*` prefix | Use direct names (NO prefix) |
| **Examples** | `company_name`, `company_employees_count` (in `where`) | `name`, `employees_count` (in `company_config.select_columns`) |
| **Response Location** | ❌ NOT returned in response | ✅ In nested `company` object |
| **Performance** | ⚡ Fast (already in index) | 🐢 Slower (separate query) |
| **Fields Available** | 13 fields (for filtering) | 27 fields (for selection) |

> **⚠️ IMPORTANT**: Denormalized `company_*` fields are **ONLY for filtering** in `where` clauses. They are **NOT available** in `select_columns`. To get company data in responses, you **MUST** use `company_config.select_columns`.

**Example - Filter by Denormalized, Populate with Company Config**:

```json
{
  "where": {
    "range_query": {
      "must": {
        "company_employees_count": {"gte": 100}  // ✅ Filter by denormalized field
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",  // ✅ Direct name, NO company_ prefix!
      "employees_count",  // ✅ Direct name, NO company_ prefix!
      "phone_number",  // ✅ Metadata field not in denormalized
      "latest_funding"  // ✅ Metadata field not in denormalized
    ]
  }
}
```

**See**: [Select Columns Guide](./filters/select_columns_filter.md) for complete documentation on `company_config` and field selection.

## Relationship with Companies

Contacts are linked to companies via the `company_id` field. You can:

1. Filter contacts by `company_id` to get all contacts for specific companies
2. Use denormalized `company_*` fields to filter contacts by company attributes in a single query
3. Use `company_config.populate` to get full company objects in responses
4. Use company filters in combination with contact filters for account-based searches
5. Join company data on the application side using the `company_id` field
