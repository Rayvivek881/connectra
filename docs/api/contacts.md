# Contacts API

## Overview

The Contacts API provides endpoints for searching, filtering, and retrieving contact data. It uses a hybrid search approach combining Elasticsearch for fast filtering and PostgreSQL for detailed data retrieval. Contacts can optionally include associated company information.

**Related Documentation:**
- [Auth API](./auth.md) - For authentication
- [Companies API](./companies.md) - For company management
- [User API](./user.md) - For user profile management

**Base URL**: `/contacts`

**Note**: This is a legacy route (not under `/api/v2/`). For consistency, consider using `/api/v2/contacts/` in future versions.

**Authentication**: All endpoints require JWT Bearer token authentication.

---

## Endpoints

### 1. Get Contacts by Filter

Search and retrieve contacts based on VQL (Vivek Query Language) filter criteria.

**Endpoint**: `POST /contacts/`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "John",
          "filter_key": "first_name",
          "search_type": "shuffle",
          "operator": "and",
          "fuzzy": true
        },
        {
          "text_value": "engineer",
          "filter_key": "title",
          "search_type": "shuffle"
        }
      ],
      "must_not": []
    },
    "keyword_match": {
      "must": {
        "country": "USA",
        "departments": ["Engineering", "Product"]
      },
      "must_not": {
        "email_status": "invalid"
      }
    },
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2024-01-01T00:00:00Z",
          "lte": "2024-12-31T23:59:59Z"
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
  "limit": 25,
  "select_columns": ["first_name", "last_name", "title", "email", "company_id"],
  "company_config": {
    "populate": true,
    "select_columns": ["name", "industries", "employees_count"]
  }
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `where` | object | No | Filter conditions (VQL query structure) |
| `where.text_matches` | object | No | Text search filters |
| `where.text_matches.must` | array | No | Text match conditions (AND logic) |
| `where.text_matches.must_not` | array | No | Text match exclusions (NOT logic) |
| `where.keyword_match` | object | No | Exact keyword/array filters |
| `where.keyword_match.must` | object | No | Keyword match conditions |
| `where.keyword_match.must_not` | object | No | Keyword exclusions |
| `where.range_query` | object | No | Numeric/date range filters |
| `where.range_query.must` | object | No | Range conditions |
| `order_by` | array | No | Sorting configuration |
| `page` | integer | No | Page number (1-indexed, max: 10) |
| `limit` | integer | No | Results per page (max: 100, default: 25) |
| `search_after` | array | No | Cursor-based pagination values |
| `select_columns` | array | No | Specific contact fields to return |
| `company_config` | object | No | Company data population configuration |
| `company_config.populate` | boolean | No | Whether to include company data (default: false) |
| `company_config.select_columns` | array | No | Specific company fields to return |

**Text Match Structure**:

```json
{
  "text_value": "search text",
  "filter_key": "field_name",
  "search_type": "exact" | "shuffle" | "substring",
  "slop": 0,
  "operator": "and" | "or",
  "fuzzy": true | false
}
```

**Response** (200 OK):

```json
{
  "success": true,
  "data": [
    {
      "id": 456,
      "uuid": "660e8400-e29b-41d4-a716-446655440001",
      "first_name": "John",
      "last_name": "Doe",
      "company_id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john.doe@example.com",
      "title": "Senior Software Engineer",
      "departments": ["Engineering", "Product"],
      "mobile_phone": "+1-555-0123",
      "email_status": "verified",
      "seniority": "Senior",
      "city": "San Francisco",
      "state": "California",
      "country": "USA",
      "linkedin_url": "https://linkedin.com/in/johndoe",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "company": {
        "uuid": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Tech Solutions Inc",
        "industries": ["Software", "Technology"],
        "employees_count": 250
      }
    }
  ]
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Request success status |
| `data` | array | Array of contact objects |
| `data[].company` | object | Associated company data (if `company_config.populate` is true) |

**Error Responses**:

- `400 Bad Request`: Invalid VQL query or validation error
- `401 Unauthorized`: Missing or invalid authentication token
- `500 Internal Server Error`: Server error

---

### 2. Get Contacts Count

Get the total count of contacts matching the filter criteria.

**Endpoint**: `POST /contacts/count`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**: Same as Get Contacts by Filter (VQL query structure)

**Response** (200 OK):

```json
{
  "success": true,
  "count": 5432
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `count` | integer | Total number of contacts matching the filter |
| `success` | boolean | Request success status |

**Error Responses**:

- `400 Bad Request`: Invalid VQL query
- `401 Unauthorized`: Missing or invalid authentication token
- `500 Internal Server Error`: Server error

---

### 3. Get Filters

Retrieve all available filter metadata for contacts.

**Endpoint**: `GET /contacts/filters`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```
Authorization: Bearer <access_token>
```

**Request Body**: None

**Response** (200 OK):

```json
{
  "success": true,
  "data": [
    {
      "filter_key": "first_name",
      "filter_type": "text_matches",
      "display_name": "First Name",
      "description": "Search by first name"
    },
    {
      "filter_key": "departments",
      "filter_type": "keyword_match",
      "display_name": "Departments",
      "description": "Filter by departments"
    },
    {
      "filter_key": "created_at",
      "filter_type": "range_query",
      "display_name": "Created Date",
      "description": "Filter by creation date"
    }
  ]
}
```

**Error Responses**:

- `401 Unauthorized`: Missing or invalid authentication token
- `500 Internal Server Error`: Server error

---

### 4. Get Filter Data

Retrieve available values for a specific filter (useful for dropdowns/autocomplete).

**Endpoint**: `POST /contacts/filters/data`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "filter_key": "departments",
  "limit": 100,
  "page": 1
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `filter_key` | string | Yes | The filter field key (e.g., "departments", "country") |
| `limit` | integer | No | Maximum number of results (default: 100) |
| `page` | integer | No | Page number for pagination |

**Response** (200 OK):

```json
{
  "success": true,
  "data": [
    {
      "value": "Engineering",
      "display_value": "Engineering"
    },
    {
      "value": "Product",
      "display_value": "Product"
    },
    {
      "value": "Sales",
      "display_value": "Sales"
    }
  ]
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `value` | string | The actual filter value |
| `display_value` | string | Human-readable display value |

**Error Responses**:

- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Missing or invalid authentication token
- `500 Internal Server Error`: Server error

---

## Filter Field Reference

### Text Match Fields

These fields support full-text search with different search types:

| Field | Search Types | Description |
|-------|--------------|-------------|
| `first_name` | exact, shuffle | Contact's first name |
| `last_name` | exact, shuffle | Contact's last name |
| `title` | exact, shuffle | Job title |
| `city` | exact, shuffle | City name |
| `state` | exact, shuffle | State/Province |
| `country` | exact, shuffle | Country name |
| `linkedin_url` | exact, shuffle | LinkedIn profile URL |

**Search Types**:
- `exact`: Exact phrase matching (match_phrase)
- `shuffle`: Flexible word matching (match with fuzziness)
- `substring`: Partial string matching (not commonly used for contacts)

### Keyword Match Fields

These fields support exact matching on keywords or arrays:

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Contact ID |
| `company_id` | string | Associated company UUID |
| `departments` | array[string] | Departments array |
| `seniority` | string | Seniority level |
| `email_status` | string | Email verification status |
| `stage` | string | Contact stage |
| `city` | string | City (exact match) |
| `state` | string | State (exact match) |
| `country` | string | Country (exact match) |

### Range Query Fields

These fields support numeric/date range filtering:

| Field | Type | Operators | Description |
|-------|------|-----------|-------------|
| `created_at` | datetime | gte, lte, gt, lt | Creation date (ISO 8601) |

**Range Operators**:
- `gte`: Greater than or equal
- `lte`: Less than or equal
- `gt`: Greater than
- `lt`: Less than

### Denormalized Company Fields

The contact index also includes denormalized company fields with `company_` prefix. These allow filtering contacts directly by company attributes:

- `company_name` - Company name (text search)
- `company_industries` - Company industries (keyword match)
- `company_employees_count` - Company employee count (range query)
- `company_country` - Company country (keyword match)
- `company_state` - Company state (keyword match)
- `company_city` - Company city (keyword match)

**Note**: These fields are indexed in Elasticsearch for efficient filtering but are not returned in the response. Use `company_config.populate` to get actual company data.

---

## VQL Query Examples

### Simple Text Search

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
    }
  },
  "page": 1,
  "limit": 25
}
```

### Keyword Filter

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering", "Product"],
        "country": "USA"
      },
      "must_not": {
        "email_status": "invalid"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Filter by Company Attributes

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Technology"],
        "company_country": "USA"
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
  "limit": 25,
  "company_config": {
    "populate": true,
    "select_columns": ["name", "industries", "employees_count"]
  }
}
```

### Combined Filters with Company Data

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "John",
          "filter_key": "first_name",
          "search_type": "shuffle"
        }
      ]
    },
    "keyword_match": {
      "must": {
        "departments": ["Engineering"],
        "country": "USA"
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
  "limit": 25,
  "select_columns": ["first_name", "last_name", "title", "email", "company_id"],
  "company_config": {
    "populate": true,
    "select_columns": ["name", "industries", "employees_count", "city"]
  }
}
```

---

## Example cURL Requests

### Get Contacts

```bash
curl -X POST http://localhost:8000/contacts/ \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "country": "USA"
        }
      }
    },
    "page": 1,
    "limit": 25
  }'
```

### Get Contacts with Company Data

```bash
curl -X POST http://localhost:8000/contacts/ \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "departments": ["Engineering"]
        }
      }
    },
    "page": 1,
    "limit": 25,
    "company_config": {
      "populate": true,
      "select_columns": ["name", "industries"]
    }
  }'
```

### Get Count

```bash
curl -X POST http://localhost:8000/contacts/count \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "departments": ["Engineering"]
        }
      }
    }
  }'
```

### Get Filters

```bash
curl -X GET http://localhost:8000/contacts/filters \
  -H "Authorization: Bearer <access_token>"
```

### Get Filter Data

```bash
curl -X POST http://localhost:8000/contacts/filters/data \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "filter_key": "departments",
    "limit": 50
  }'
```

---

## Response-Only Fields

The following fields are stored in PostgreSQL but are **NOT** indexed in Elasticsearch and **cannot be used in filters**. They are only returned in API responses:

- `email` - Email address
- `mobile_phone` - Mobile phone number
- `work_direct_phone` - Work direct phone
- `home_phone` - Home phone number
- `other_phone` - Other phone number
- `facebook_url` - Facebook profile URL
- `twitter_url` - Twitter profile URL
- `website` - Personal website URL

---

## Company Data Population

When `company_config.populate` is set to `true`, the API will:

1. Fetch contacts matching the filter criteria
2. Extract unique `company_id` values from the results
3. Fetch company data from PostgreSQL for those companies
4. Merge company data into each contact response

**Performance Note**: Populating company data adds a database query but provides complete contact information in a single API call. Use `company_config.select_columns` to limit which company fields are fetched.

**Example**:

```json
{
  "company_config": {
    "populate": true,
    "select_columns": ["name", "industries", "employees_count", "city", "country"]
  }
}
```

If `populate` is `false` or omitted, the `company` field will be `null` in the response.

---

## Notes

- The API uses a hybrid approach: Elasticsearch for fast filtering, PostgreSQL for detailed data
- Only fields specified in `select_columns` are fetched from PostgreSQL (improves performance)
- If `select_columns` is empty, all available contact fields are returned
- Company data is fetched separately and merged when `company_config.populate` is true
- Pagination supports both page-based (`page`, `limit`) and cursor-based (`search_after`) approaches
- Maximum page number is 10
- Maximum limit per page is 100
- Default limit is 25 if not specified
- Denormalized company fields (`company_*`) allow efficient filtering by company attributes without joins

