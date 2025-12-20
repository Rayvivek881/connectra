# Companies API

## Overview

The Companies API provides endpoints for searching, filtering, and retrieving company data. It uses a hybrid search approach combining Elasticsearch for fast filtering and PostgreSQL for detailed data retrieval.

**Related Documentation:**
- [Auth API](./auth.md) - For authentication
- [Contacts API](./contacts.md) - For contact management
- [User API](./user.md) - For user profile management

**Base URL**: `/companies`

**Note**: This is a legacy route (not under `/api/v2/`). For consistency, consider using `/api/v2/companies/` in future versions.

**Authentication**: All endpoints require JWT Bearer token authentication.

---

## Endpoints

### 1. Get Companies by Filter

Search and retrieve companies based on VQL (Vivek Query Language) filter criteria.

**Endpoint**: `POST /companies/`

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
          "text_value": "technology",
          "filter_key": "name",
          "search_type": "shuffle",
          "operator": "and",
          "fuzzy": true
        }
      ],
      "must_not": []
    },
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology"],
        "country": "USA"
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
  "limit": 25,
  "select_columns": ["name", "employees_count", "industries", "city", "country"]
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
| `select_columns` | array | No | Specific fields to return from PostgreSQL |

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
      "id": 123,
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Tech Solutions Inc",
      "employees_count": 250,
      "industries": ["Software", "Technology"],
      "keywords": ["AI", "Machine Learning"],
      "address": "123 Tech Street",
      "annual_revenue": 50000000,
      "total_funding": 10000000,
      "technologies": ["Python", "JavaScript"],
      "city": "San Francisco",
      "state": "California",
      "country": "USA",
      "linkedin_url": "https://linkedin.com/company/tech-solutions",
      "website": "https://techsolutions.com",
      "normalized_domain": "techsolutions.com",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Error Responses**:

- `400 Bad Request`: Invalid VQL query or validation error
- `401 Unauthorized`: Missing or invalid authentication token
- `500 Internal Server Error`: Server error

---

### 2. Get Companies Count

Get the total count of companies matching the filter criteria.

**Endpoint**: `POST /companies/count`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**: Same as Get Companies by Filter (VQL query structure)

**Response** (200 OK):

```json
{
  "success": true,
  "count": 1250
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `count` | integer | Total number of companies matching the filter |
| `success` | boolean | Request success status |

**Error Responses**:

- `400 Bad Request`: Invalid VQL query
- `401 Unauthorized`: Missing or invalid authentication token
- `500 Internal Server Error`: Server error

---

### 3. Get Filters

Retrieve all available filter metadata for companies.

**Endpoint**: `GET /companies/filters`

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
      "filter_key": "name",
      "filter_type": "text_matches",
      "display_name": "Company Name",
      "description": "Search by company name"
    },
    {
      "filter_key": "industries",
      "filter_type": "keyword_match",
      "display_name": "Industries",
      "description": "Filter by industries"
    },
    {
      "filter_key": "employees_count",
      "filter_type": "range_query",
      "display_name": "Employee Count",
      "description": "Filter by number of employees"
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

**Endpoint**: `POST /companies/filters/data`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "filter_key": "industries",
  "limit": 100,
  "page": 1
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `filter_key` | string | Yes | The filter field key (e.g., "industries", "country") |
| `limit` | integer | No | Maximum number of results (default: 100) |
| `page` | integer | No | Page number for pagination |

**Response** (200 OK):

```json
{
  "success": true,
  "data": [
    {
      "value": "Software",
      "display_value": "Software"
    },
    {
      "value": "Technology",
      "display_value": "Technology"
    },
    {
      "value": "Healthcare",
      "display_value": "Healthcare"
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
| `name` | exact, shuffle, substring | Company name (supports ngram for substring) |
| `address` | exact, shuffle | Company address |
| `city` | exact, shuffle | City name |
| `state` | exact, shuffle | State/Province |
| `country` | exact, shuffle | Country name |
| `linkedin_url` | exact, shuffle | LinkedIn company URL |
| `website` | exact, shuffle | Company website URL |
| `normalized_domain` | exact, shuffle | Normalized domain name |

**Search Types**:
- `exact`: Exact phrase matching (match_phrase)
- `shuffle`: Flexible word matching (match with fuzziness)
- `substring`: Partial string matching (ngram, only for `name` field)

### Keyword Match Fields

These fields support exact matching on keywords or arrays:

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Company ID |
| `industries` | array[string] | Industries array |
| `keywords` | array[string] | Keywords array |
| `technologies` | array[string] | Technologies array |
| `city` | string | City (exact match) |
| `state` | string | State (exact match) |
| `country` | string | Country (exact match) |

### Range Query Fields

These fields support numeric/date range filtering:

| Field | Type | Operators | Description |
|-------|------|-----------|-------------|
| `employees_count` | integer | gte, lte, gt, lt | Employee count |
| `annual_revenue` | integer | gte, lte, gt, lt | Annual revenue (in cents) |
| `total_funding` | integer | gte, lte, gt, lt | Total funding (in cents) |
| `created_at` | datetime | gte, lte, gt, lt | Creation date (ISO 8601) |

**Range Operators**:
- `gte`: Greater than or equal
- `lte`: Less than or equal
- `gt`: Greater than
- `lt`: Less than

---

## VQL Query Examples

### Simple Text Search

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "tech",
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

### Keyword Filter

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology"],
        "country": "USA"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Range Query

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
  "page": 1,
  "limit": 25
}
```

### Combined Filters

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "AI",
          "filter_key": "name",
          "search_type": "shuffle"
        }
      ]
    },
    "keyword_match": {
      "must": {
        "industries": ["Technology"],
        "country": "USA"
      }
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 100,
          "lte": 1000
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
  "limit": 25,
  "select_columns": ["name", "employees_count", "industries", "city"]
}
```

---

## Example cURL Requests

### Get Companies

```bash
curl -X POST http://localhost:8000/companies/ \
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

### Get Count

```bash
curl -X POST http://localhost:8000/companies/count \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "industries": ["Technology"]
        }
      }
    }
  }'
```

### Get Filters

```bash
curl -X GET http://localhost:8000/companies/filters \
  -H "Authorization: Bearer <access_token>"
```

### Get Filter Data

```bash
curl -X POST http://localhost:8000/companies/filters/data \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "filter_key": "industries",
    "limit": 50
  }'
```

---

## Response-Only Fields

The following fields are stored in PostgreSQL but are **NOT** indexed in Elasticsearch and **cannot be used in filters**. They are only returned in API responses:

- `facebook_url`
- `twitter_url`
- `company_name_for_emails`
- `phone_number`
- `latest_funding`
- `latest_funding_amount`
- `last_raised_at`

---

## Notes

- The API uses a hybrid approach: Elasticsearch for fast filtering, PostgreSQL for detailed data
- Only fields specified in `select_columns` are fetched from PostgreSQL (improves performance)
- If `select_columns` is empty, all available fields are returned
- Pagination supports both page-based (`page`, `limit`) and cursor-based (`search_after`) approaches
- Maximum page number is 10
- Maximum limit per page is 100
- Default limit is 25 if not specified

