# Filter Field Reference

## Table of Contents

1. [Overview](#overview)
2. [Company Fields](#company-fields)
3. [Contact Fields](#contact-fields)
4. [Denormalized Company Fields in Contact Index](#denormalized-company-fields-in-contact-index)
5. [Response-Only Fields](#response-only-fields)
6. [Field Type Reference](#field-type-reference)
7. [Filter Availability](#filter-availability)
8. [Sorting Capabilities](#sorting-capabilities)
9. [Authentication and Security](#authentication-and-security)
10. [VQL Syntax Reference](#vql-syntax-reference)
11. [Pagination Strategies](#pagination-strategies)
12. [Field Selection Optimization](#field-selection-optimization)
13. [Error Handling](#error-handling)

## Overview

This document provides a comprehensive reference for all fields in the Company and Contact APIs. Fields are categorized as:

- **Filterable**: Fields indexed in Elasticsearch that can be used in filter queries
- **Response-Only**: Fields stored in PostgreSQL but not indexed, only returned in API responses

Each field includes its data type, search capabilities, sorting support, ngram configuration, and usage examples.

> **Authentication Required**: All API requests require the `X-API-Key` header for authentication. See [06-api-reference.md](./06-api-reference.md#authentication) for details.

---

## Company Fields

### Filterable Company Fields (Elasticsearch Index)

These fields are indexed in Elasticsearch and can be used in filter queries.

#### Text Search Fields

| Field | Type | Search Type | Sortable | Filterable | Ngram Support | Description |
|-------|------|-------------|----------|------------|---------------|-------------|
| `name` | string | text | No | Yes | Yes (3-10) | Company name |
| `address` | string | text | No | Yes | No | Company address |
| `city` | string | text | No | Yes | No | City name |
| `state` | string | text | No | Yes | No | State/Province |
| `country` | string | text | No | Yes | No | Country name |
| `linkedin_url` | string | text | No | Yes | No | LinkedIn company URL |
| `website` | string | text | No | Yes | No | Company website URL |
| `normalized_domain` | string | text | No | Yes | No | Normalized domain name |

**Ngram Support**: Only the `name` field has ngram analysis enabled (min_gram: 3, max_gram: 10), allowing efficient substring/partial matching. Other text fields do not support substring search.

**Note**: `city`, `state`, and `country` are text fields in Elasticsearch. They can be used in `text_matches` for fuzzy/flexible search, but for exact matching, use `keyword_match` with the exact values.

#### Keyword Fields

| Field | Type | Search Type | Sortable | Filterable | Description |
|-------|------|-------------|----------|------------|-------------|
| `id` | integer | keyword | Yes | Yes | Company ID |
| `industries` | array[string] | keyword | Yes | Yes | Industries array |
| `keywords` | array[string] | keyword | Yes | Yes | Keywords array |
| `technologies` | array[string] | keyword | Yes | Yes | Technologies array |

#### Range Query Fields

| Field | Type | Search Type | Sortable | Filterable | Description |
|-------|------|-------------|----------|------------|-------------|
| `employees_count` | integer | range | Yes | Yes | Employee count |
| `annual_revenue` | integer | range | Yes | Yes | Annual revenue (in cents) |
| `total_funding` | integer | range | Yes | Yes | Total funding (in cents) |
| `created_at` | datetime | range | Yes | Yes | Creation date (ISO 8601) |

### Company Field Details

#### `name` (text)

- **Search Type**: Text search (exact, shuffle, substring)
- **Ngram Support**: Yes - has `name.ngram` sub-field for substring matching (min_gram: 3, max_gram: 10)
- **Use Cases**: Company name searches, brand searches, partial name matching
- **Substring Search**: Minimum 3 characters required
- **Example Query (Shuffle)**:

```json
{
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
}
```

- **Example Query (Substring)**:

```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "soft",
        "filter_key": "name",
        "search_type": "substring",
        "operator": "and"
      }
    ]
  }
}
```

#### `industries` (keyword array)

- **Search Type**: Keyword match only
- **Use Cases**: Industry filtering, market segmentation
- **Common Values**: "Software", "Technology", "SaaS", "Healthcare", "Finance", "Retail", etc.
- **Example Query**:

```json
{
  "keyword_match": {
    "must": {
      "industries": ["Software", "Technology"]
    }
  }
}
```

#### `technologies` (keyword array)

- **Search Type**: Keyword match only
- **Use Cases**: Technology stack filtering, technical targeting
- **Common Values**: "Python", "Go", "React", "JavaScript", "AWS", "Azure", "GCP", etc.
- **Example Query**:

```json
{
  "keyword_match": {
    "must": {
      "technologies": ["Python", "Go"]
    }
  }
}
```

#### `employees_count` (integer range)

- **Search Type**: Range query only
- **Use Cases**: Company size filtering, segmentation
- **Example Query**:

```json
{
  "range_query": {
    "must": {
      "employees_count": {
        "gte": 50,
        "lte": 500
      }
    }
  }
}
```

#### `annual_revenue` (integer range)

- **Search Type**: Range query only
- **Note**: Stored in cents (e.g., 5000000 = $50,000.00)
- **Use Cases**: Revenue-based targeting, high-value account identification
- **Example Query**:

```json
{
  "range_query": {
    "must": {
      "annual_revenue": {
        "gte": 1000000
      }
    }
  }
}
```

---

## Contact Fields

### Filterable Contact Fields (Elasticsearch Index)

These fields are indexed in Elasticsearch and can be used in filter queries.

#### Text Search Fields

| Field | Type | Search Type | Sortable | Filterable | Ngram Support | Description |
|-------|------|-------------|----------|------------|---------------|-------------|
| `first_name` | string | text | No | Yes | Yes (5-10) | First name |
| `last_name` | string | text | No | Yes | Yes (5-10) | Last name |
| `title` | string | text | No | Yes | Yes (5-10) | Job title |
| `city` | string | text | No | Yes | No | City name |
| `state` | string | text | No | Yes | No | State/Province |
| `country` | string | text | No | Yes | No | Country name |
| `linkedin_url` | string | text | No | Yes | No | LinkedIn profile URL |

**Ngram Support**: The `first_name`, `last_name`, and `title` fields have ngram analysis enabled (min_gram: 5, max_gram: 10), allowing efficient substring/partial matching. Other text fields do not support substring search.

**Note**: `city`, `state`, and `country` are text fields in Elasticsearch. They can be used in `text_matches` for fuzzy/flexible search, but for exact matching, use `keyword_match` with the exact values.

#### Keyword Fields

| Field | Type | Search Type | Sortable | Filterable | Description |
|-------|------|-------------|----------|------------|-------------|
| `id` | integer | keyword | Yes | Yes | Contact ID |
| `company_id` | string | keyword | Yes | Yes | Company UUID |
| `email` | string | keyword | Yes | Yes | Email address |
| `departments` | array[string] | keyword | Yes | Yes | Departments array |
| `mobile_phone` | string | keyword | Yes | Yes | Mobile phone number |
| `email_status` | string | keyword | Yes | Yes | Email verification status |
| `seniority` | string | keyword | Yes | Yes | Seniority level |

#### Range Query Fields

| Field | Type | Search Type | Sortable | Filterable | Description |
|-------|------|-------------|----------|------------|-------------|
| `created_at` | datetime | range | Yes | Yes | Creation date (ISO 8601) |

### Contact Field Details

#### `first_name` / `last_name` (text)

- **Search Type**: Text search (exact, shuffle, substring)
- **Ngram Support**: Yes - both have `.ngram` sub-fields (`first_name.ngram`, `last_name.ngram`) for substring matching (min_gram: 5, max_gram: 10)
- **Use Cases**: Name searches, people lookup, partial name matching
- **Substring Search**: Minimum 5 characters required
- **Best Practice**: Search both fields separately for better results
- **Example Query (Shuffle)**:

```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "John",
        "filter_key": "first_name",
        "search_type": "shuffle",
        "fuzzy": true
      },
      {
        "text_value": "Smith",
        "filter_key": "last_name",
        "search_type": "shuffle",
        "fuzzy": true
      }
    ]
  }
}
```

- **Example Query (Substring)**:

```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "johnn",
        "filter_key": "first_name",
        "search_type": "substring",
        "operator": "and"
      }
    ]
  }
}
```

#### `title` (text)

- **Search Type**: Text search (exact, shuffle, substring)
- **Ngram Support**: Yes - has `title.ngram` sub-field for substring matching (min_gram: 5, max_gram: 10)
- **Use Cases**: Job title searches, role-based targeting, partial title matching
- **Substring Search**: Minimum 5 characters required
- **Best Practice**: Use `exact` with `slop` for multi-word titles, use `substring` for partial matches
- **Example Query (Exact)**:

```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "software engineer",
        "filter_key": "title",
        "search_type": "exact",
        "slop": 2
      }
    ]
  }
}
```

- **Example Query (Substring)**:

```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "engin",
        "filter_key": "title",
        "search_type": "substring",
        "operator": "and"
      }
    ]
  }
}
```

#### `departments` (keyword array)

- **Search Type**: Keyword match only
- **Use Cases**: Department filtering, organizational targeting
- **Common Values**: "Engineering", "Sales", "Marketing", "Customer Success", "Support", "HR", "Finance", "Product", "Operations", "Legal", "Research"
- **Example Query**:

```json
{
  "keyword_match": {
    "must": {
      "departments": ["Engineering", "Sales"]
    }
  }
}
```

#### `email_status` (keyword)

- **Search Type**: Keyword match only
- **Use Cases**: Quality filtering, contactability assurance
- **Values**: `"verified"`, `"unverified"`, `"invalid"`, `"bounced"`
- **Best Practice**: Always filter by `"verified"` for high-quality contacts
- **Example Query**:

```json
{
  "keyword_match": {
    "must": {
      "email_status": "verified"
    }
  }
}
```

#### `seniority` (keyword)

- **Search Type**: Keyword match only
- **Use Cases**: Decision-maker targeting, role-based filtering
- **Values**: `"Junior"`, `"Mid"`, `"Senior"`, `"Lead"`, `"Principal"`, `"Executive"`
- **Best Practice**: Use `["Senior", "Lead", "Principal", "Executive"]` for decision-makers
- **Example Query**:

```json
{
  "keyword_match": {
    "must": {
      "seniority": ["Senior", "Lead", "Principal", "Executive"]
    }
  }
}
```

#### `company_id` (keyword)

- **Search Type**: Keyword match only
- **Use Cases**: Account-based filtering, company-contact relationships
- **Best Practice**: Use for finding all contacts at specific companies
- **Example Query**:

```json
{
  "keyword_match": {
    "must": {
      "company_id": [
        "c0a8012e-1111-2222-3333-444455556666",
        "c0a8012e-2222-3333-4444-555566667777"
      ]
    }
  }
}
```

---

## Denormalized Company Fields in Contact Index

The contact index includes denormalized company data with the `company_` prefix, allowing you to filter contacts directly by company attributes without a separate company query.

### Denormalized Company Fields

| Field | Type | Search Type | Sortable | Filterable | Ngram Support | Description |
|-------|------|-------------|----------|------------|---------------|-------------|
| `company_name` | string | text | No | Yes | Yes (5-10) | Company name |
| `company_employees_count` | integer | range | No | Yes | No | Company employee count |
| `company_industries` | array[string] | keyword | No | Yes | No | Company industries |
| `company_keywords` | array[string] | keyword | No | Yes | No | Company keywords |
| `company_address` | string | text | No | Yes | No | Company address |
| `company_annual_revenue` | integer | range | No | Yes | No | Company annual revenue (in cents) |
| `company_total_funding` | integer | range | No | Yes | No | Company total funding (in cents) |
| `company_technologies` | array[string] | keyword | No | Yes | No | Company technologies |
| `company_city` | string | text | No | Yes | No | Company city |
| `company_state` | string | text | No | Yes | No | Company state |
| `company_country` | string | text | No | Yes | No | Company country |
| `company_linkedin_url` | string | text | No | Yes | No | Company LinkedIn URL |
| `company_website` | string | text | No | Yes | Yes (5-10) | Company website |
| `company_normalized_domain` | string | text | No | Yes | Yes (5-10) | Company normalized domain |

**Ngram Support**: `company_name`, `company_website`, and `company_normalized_domain` support substring search via ngram (min_gram: 5, max_gram: 10), same as contact name fields.

### Example: Filtering Contacts by Company Attributes

**Example: Find contacts at high-revenue software companies**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS"],
        "seniority": ["Senior", "Lead"]
      }
    },
    "range_query": {
      "must": {
        "company_annual_revenue": {
          "gte": 5000000
        },
        "company_employees_count": {
          "gte": 100
        }
      }
    }
  }
}
```

**Example: Search contacts by company name (substring)**
```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "micros",
          "filter_key": "company_name",
          "search_type": "substring",
          "operator": "and"
        }
      ]
    }
  }
}
```

**Benefits**:

- Single query to filter contacts by company attributes
- No need to first query companies, then filter contacts by `company_id`
- Faster queries for account-based filtering scenarios
- Enables complex filtering combining contact and company criteria in one request

---

## Response-Only Fields

These fields are stored in PostgreSQL but are **NOT** indexed in Elasticsearch and **cannot be used in filters**. They are only returned in API responses and can be selected using `select_columns`.

### Company Response-Only Fields

| Field | Type | Description |
|-------|------|-------------|
| `uuid` | string | Company UUID |
| `facebook_url` | string | Facebook page URL |
| `twitter_url` | string | Twitter profile URL |
| `company_name_for_emails` | string | Company name formatted for email use |
| `phone_number` | string | Company phone number |
| `latest_funding` | string | Latest funding round (e.g., "Series B") |
| `latest_funding_amount` | integer | Latest funding amount (in cents) |
| `last_raised_at` | string | Date of last funding round |
| `updated_at` | datetime | Last update timestamp |
| `deleted_at` | datetime | Soft delete timestamp (null if active) |

### Contact Response-Only Fields

| Field | Type | Description |
|-------|------|-------------|
| `uuid` | string | Contact UUID |
| `facebook_url` | string | Facebook profile URL |
| `twitter_url` | string | Twitter profile URL |
| `website` | string | Personal website URL |
| `work_direct_phone` | string | Work direct phone number |
| `home_phone` | string | Home phone number |
| `other_phone` | string | Other phone number |
| `stage` | string | Contact stage (e.g., "Closed Won") |
| `updated_at` | datetime | Last update timestamp |
| `deleted_at` | datetime | Soft delete timestamp (null if active) |

**Important**: These fields cannot be used in `where` clauses. They are only available in response data and can be selected using `select_columns` if needed.

---

## Field Type Reference

### Text Fields

**Characteristics**:

- Support full-text search with `text_matches`
- Can use `search_type: "exact"` (phrase matching), `"shuffle"` (word matching), or `"substring"` (partial matching)
- Support fuzzy matching for typo tolerance (exact and shuffle only)
- **Cannot be used for sorting** (analyzed fields in Elasticsearch)
- Some fields have `.ngram` sub-fields for efficient substring matching

**Company Text Fields**: `name`, `address`, `city`, `state`, `country`, `linkedin_url`, `website`, `normalized_domain`

**Contact Text Fields**: `first_name`, `last_name`, `title`, `city`, `state`, `country`, `linkedin_url`

**Ngram Fields for Substring Search**:

- Certain text fields have `.ngram` sub-fields that enable efficient partial text matching
- **Company fields with ngram support**: `name.ngram` (min_gram: 3, max_gram: 10)
- **Contact fields with ngram support**: `first_name.ngram`, `last_name.ngram`, `title.ngram` (min_gram: 5, max_gram: 10)
- **Denormalized company fields with ngram**: `company_name.ngram`, `company_website.ngram`, `company_normalized_domain.ngram` (min_gram: 5, max_gram: 10)
- When using `search_type: "substring"`, the system automatically queries the `.ngram` field
- Use substring search for autocomplete-style queries and partial word matching
- **Important**: Minimum search text length must match the min_gram value (3 for company name, 5 for contact names/titles)

### Keyword Fields

**Characteristics**:

- Support exact matching with `keyword_match`
- Can match single values or arrays
- **Can be used for sorting** (not analyzed in Elasticsearch)
- Faster than text searches

**Company Keyword Fields**: `id`, `industries`, `keywords`, `technologies`

**Contact Keyword Fields**: `id`, `company_id`, `email`, `departments`, `mobile_phone`, `email_status`, `seniority`

**Denormalized Company Keyword Fields**: `company_industries`, `company_keywords`, `company_technologies`

### Range Query Fields

**Characteristics**:

- Support numeric or date range queries with `range_query`
- Operators: `gte`, `lte`, `gt`, `lt`
- **Can be used for sorting**
- Dates must be in ISO 8601 format (RFC3339)

**Company Range Fields**: `employees_count`, `annual_revenue`, `total_funding`, `created_at`

**Contact Range Fields**: `created_at`

**Denormalized Company Range Fields**: `company_employees_count`, `company_annual_revenue`, `company_total_funding`

---

## Filter Availability

### Direct-Derived Filters

**Definition**: Filter values are extracted directly from database records in real-time.

**Company Direct-Derived Filters**:

- `address`
- `annual_revenue`
- `employees_count`
- `linkedin_url`
- `normalized_domain`
- `total_funding`
- `website`

**Contact Direct-Derived Filters**:

- `company_id`
- `email`
- `first_name`
- `last_name`
- `linkedin_url`
- `mobile_phone`

**Characteristics**:

- Values come directly from PostgreSQL records
- Always up-to-date
- May be slower for filters with many distinct values

### Stored Filters

**Definition**: Filter values are pre-computed and stored in the `filters_data` table for faster access.

**Company Stored Filters**:

- `city`
- `country`
- `industries`
- `keywords`
- `state`
- `technologies`
- `uuid` (displayed as "Name" in filter UI, but not filterable in Elasticsearch)

**Contact Stored Filters**:

- `city`
- `country`
- `departments`
- `email_status`
- `seniority`
- `state`
- `title`

**Characteristics**:

- Values are pre-computed and cached
- Faster for frequently used filters
- May require periodic updates

---

## Sorting Capabilities

### Sortable Fields

**Company Sortable Fields**:

- `id` (keyword)
- `employees_count` (range)
- `annual_revenue` (range)
- `total_funding` (range)
- `created_at` (range)
- `industries` (keyword)
- `keywords` (keyword)
- `technologies` (keyword)

**Contact Sortable Fields**:

- `id` (keyword)
- `company_id` (keyword)
- `email` (keyword)
- `departments` (keyword)
- `mobile_phone` (keyword)
- `email_status` (keyword)
- `seniority` (keyword)
- `created_at` (range)

**Note**: Denormalized company fields in the contact index are **not sortable** because they are not keyword or range fields suitable for sorting.

### Non-Sortable Fields

**Company Non-Sortable Fields** (text fields):

- `name`
- `address`
- `city`
- `state`
- `country`
- `linkedin_url`
- `website`
- `normalized_domain`

**Contact Non-Sortable Fields** (text fields):

- `first_name`
- `last_name`
- `title`
- `city`
- `state`
- `country`
- `linkedin_url`

**Denormalized Company Non-Sortable Fields** (text fields):

- `company_name`
- `company_address`
- `company_city`
- `company_state`
- `company_country`
- `company_linkedin_url`
- `company_website`
- `company_normalized_domain`

**Why**: Text fields are analyzed in Elasticsearch (tokenized, lowercased, etc.), making them unsuitable for sorting. Use keyword fields or date fields for sorting operations.

---

## Field Usage Examples

### Company Field Examples

#### Example 1: Industry + Technology Filtering
```json
{
  "keyword_match": {
    "must": {
      "industries": ["Software", "SaaS"],
      "technologies": ["Python", "Go", "AWS"]
    }
  }
}
```

#### Example 2: Size + Revenue Filtering
```json
{
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
}
```

#### Example 3: Name + Location Text Search
```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "software",
        "filter_key": "name",
        "search_type": "shuffle",
        "fuzzy": true
      },
      {
        "text_value": "Silicon Valley",
        "filter_key": "address",
        "search_type": "shuffle",
        "fuzzy": true
      }
    ]
  }
}
```

### Contact Field Examples

#### Example 1: Department + Seniority Filtering
```json
{
  "keyword_match": {
    "must": {
      "departments": ["Engineering", "Product"],
      "seniority": ["Senior", "Lead", "Principal"],
      "email_status": "verified"
    }
  }
}
```

#### Example 2: Name + Title Text Search
```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "John",
        "filter_key": "first_name",
        "search_type": "shuffle",
        "fuzzy": true
      },
      {
        "text_value": "engineer",
        "filter_key": "title",
        "search_type": "shuffle",
        "fuzzy": true
      }
    ]
  }
}
```

#### Example 3: Company-Based Contact Filtering (using company_id)
```json
{
  "keyword_match": {
    "must": {
      "company_id": [
        "c0a8012e-1111-2222-3333-444455556666"
      ],
      "seniority": ["Senior", "Lead", "Principal", "Executive"],
      "email_status": "verified"
    }
  }
}
```

#### Example 4: Company-Based Contact Filtering (using denormalized fields)
```json
{
  "keyword_match": {
    "must": {
      "company_industries": ["Software", "SaaS"],
      "seniority": ["Senior", "Lead"],
      "email_status": "verified"
    },
    "range_query": {
      "must": {
        "company_annual_revenue": {
          "gte": 5000000
        }
      }
    }
  }
}
```

---

## Authentication and Security

### API Key Authentication

All filter endpoints require authentication using the `X-API-Key` header:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{...}'
```

### Security Best Practices

1. **Secure Storage**: Store API keys securely (environment variables, secret management systems)
2. **HTTPS Only**: Always use HTTPS in production to protect API keys in transit
3. **Key Rotation**: Rotate API keys periodically for enhanced security
4. **Separate Keys**: Use different API keys for different environments (dev, staging, production)
5. **Access Control**: Limit API key access to only necessary endpoints
6. **Security Recommendations**: 
   - Never commit API keys to version control
   - Use environment variables for configuration
   - Implement key rotation policies
   - Monitor API key usage for anomalies

**See**: [API Reference - Authentication](./06-api-reference.md#authentication) for complete authentication details.

---

## VQL Syntax Reference

VQL (Vivek Query Language) is the query language used for filtering in Connectra. This section provides a quick reference for VQL syntax.

### Filter Structure

```json
{
  "where": {
    "text_matches": {...},
    "keyword_match": {...},
    "range_query": {...}
  },
  "order_by": [...],
  "page": 1,
  "limit": 25,
  "search_after": [...],
  "select_columns": [...]
}
```

### Text Match Filters

**Structure**:
```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "search text",
        "filter_key": "field_name",
        "search_type": "exact|shuffle|substring",
        "operator": "and|or",
        "fuzzy": true|false,
        "slop": 2
      }
    ],
    "must_not": [...]
  }
}
```

**Search Types**:

- `exact`: Phrase matching with word order (supports `slop` for word distance)
- `shuffle`: Word matching (order doesn't matter)
- `substring`: Partial text matching using ngram analysis

**Ngram Fields**:

- Company `name`: min_gram: 3, max_gram: 10
- Contact `first_name`, `last_name`, `title`: min_gram: 5, max_gram: 10
- Denormalized `company_name`, `company_website`, `company_normalized_domain`: min_gram: 5, max_gram: 10

### Keyword Match Filters

**Structure**:
```json
{
  "keyword_match": {
    "must": {
      "field_name": ["value1", "value2"],
      "another_field": "single_value"
    },
    "must_not": {
      "field_name": ["excluded_value"]
    }
  }
}
```

### Range Query Filters

**Structure**:
```json
{
  "range_query": {
    "must": {
      "field_name": {
        "gte": 100,
        "lte": 500,
        "gt": 50,
        "lt": 1000
      }
    }
  }
}
```

**Operators**:

- `gte`: Greater than or equal
- `lte`: Less than or equal
- `gt`: Greater than
- `lt`: Less than

**Date Format**: ISO 8601 (RFC3339): `"YYYY-MM-DDTHH:MM:SSZ"`

**See**: [Company Filters Guide - VQL Syntax](./01-company-filters-complete-guide.md#vql-syntax-reference) for complete VQL documentation.

---

## Pagination Strategies

Connectra supports two pagination strategies: page-based and cursor-based.

### Page-Based Pagination

**Use When**: Small datasets (< 250 records), simple navigation

**Parameters**:

- `page`: Page number (1-10, maximum)
- `limit`: Results per page (1-100, maximum)

**Example**:
```json
{
  "where": {...},
  "page": 1,
  "limit": 25
}
```

**Limitations**:

- Maximum page number: 10
- Maximum page size: 100
- Not suitable for large datasets

### Cursor-Based Pagination (search_after)

**Use When**: Large datasets (> 250 records), deep pagination, stable sorting

**Parameters**:

- `search_after`: Array of sort values from last result
- `order_by`: Required for cursor-based pagination
- `limit`: Results per page (1-100)

**Example**:
```json
{
  "where": {...},
  "order_by": [
    {
      "order_by": "annual_revenue",
      "order_direction": "desc"
    },
    {
      "order_by": "id",
      "order_direction": "asc"
    }
  ],
  "search_after": [5000000, 123],
  "limit": 25
}
```

**Best Practices**:

- Always include `id` as the last sort field for stable pagination
- Use consistent `order_by` across all pages
- Store `search_after` value from last result for next page

**See**: [Company Filters Guide - Pagination](./01-company-filters-complete-guide.md#pagination-strategies) for detailed pagination documentation.

---

## Field Selection Optimization

### Using `select_columns`

The `select_columns` parameter allows you to specify which fields are returned in the API response, optimizing payload size and performance.

**Syntax**:
```json
{
  "where": {...},
  "select_columns": ["id", "name", "employees_count"]
}
```

### Benefits

1. **Reduced Payload Size**: Only return needed fields
2. **Improved Performance**: Less data to transfer and process
3. **PostgreSQL Optimization**: Fields are selected from PostgreSQL after Elasticsearch search

### Field Types

- **Filterable Fields**: Can be used in `where` clauses and `select_columns`
- **Response-Only Fields**: Cannot be used in `where` clauses but can be selected (e.g., `facebook_url`, `twitter_url`, `phone_number`, `stage`)
- **Denormalized Company Fields**: For contacts, `company_*` fields can be used in filters but NOT in `select_columns`. Use `company_config` to get company data.

### Contact-Specific: Company Config

For contacts, use `company_config` to populate full company objects:

```json
{
  "where": {...},
  "select_columns": ["id", "first_name", "last_name", "company_id"],
  "company_config": {
    "populate": true,
    "select_columns": ["uuid", "name", "employees_count", "industries"]
  }
}
```

**Note**: `company_config.select_columns` uses direct field names (e.g., `name`, `employees_count`) without the `company_` prefix.

**See**: [Contact Filters Guide - Field Selection](./02-contact-filters-complete-guide.md#field-selection-optimization) for complete field selection documentation.

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

### Common Error Codes

#### 400 Bad Request

- Invalid request body
- Missing required fields
- Invalid field types
- Page number/size exceeded

#### 401 Unauthorized

- Missing `X-API-Key` header
- Invalid API key

#### 429 Too Many Requests

- Rate limit exceeded
- Too many requests in time window

#### 500 Internal Server Error

- Elasticsearch errors
- Database errors
- Server-side issues

### Error Handling Best Practices

1. **Check `success` field** in all responses
2. **Handle rate limits** with exponential backoff
3. **Validate requests** before sending
4. **Log errors** for debugging
5. **Retry with backoff** for transient errors

**See**: [API Reference - Error Handling](./06-api-reference.md#error-handling) for complete error handling guide with examples and solutions.

---

## Fields in Write Operations

When creating or updating companies and contacts, you can use the following fields. Some fields are required, others are optional.

### Company Write Operations

#### Required Fields for Create

- `name` (string): Company name

#### Optional Fields for Create/Update

**Basic Information**:
- `uuid` (string): UUID (generated if not provided)
- `normalized_domain` (string): Company domain
- `address` (string): Company address
- `city` (string): City
- `state` (string): State/Province
- `country` (string): Country

**Company Metrics**:
- `employees_count` (integer): Number of employees (>= 0)
- `annual_revenue` (number): Annual revenue (>= 0)
- `total_funding` (number): Total funding amount (>= 0)
- `latest_funding_amount` (number): Latest funding amount (>= 0)

**Categorical Data**:
- `industries` (array of strings): Industry categories
- `keywords` (array of strings): Keywords/tags
- `technologies` (array of strings): Technologies used

**URLs**:
- `website` (string): Company website URL
- `linkedin_url` (string): LinkedIn company page URL
- `facebook_url` (string): Facebook page URL
- `twitter_url` (string): Twitter handle URL

**Other Fields**:
- `company_name_for_emails` (string): Name used in emails
- `phone_number` (string): Phone number
- `latest_funding` (string): Latest funding round
- `last_raised_at` (string): Date of last funding (ISO 8601)

**Validation Rules**:
- `normalized_domain`: Must be valid FQDN if provided
- `employees_count`, `annual_revenue`, `total_funding`, `latest_funding_amount`: Must be >= 0
- `linkedin_url`, `website`, `facebook_url`, `twitter_url`: Must be valid URL format if provided

### Contact Write Operations

#### Required Fields for Create

- `email` (string): Contact email (must be valid email format)
- `first_name` (string): First name
- `last_name` (string): Last name

#### Optional Fields for Create/Update

**Basic Information**:
- `uuid` (string): UUID (generated if not provided)
- `title` (string): Job title
- `company_id` (string): UUID of associated company
- `mobile_phone` (string): Mobile phone number
- `work_direct_phone` (string): Work direct phone number

**Location**:
- `city` (string): City
- `state` (string): State/Province
- `country` (string): Country

**Categorical Data**:
- `departments` (array of strings): Department names
- `seniority` (string): Seniority level (e.g., "Junior", "Senior", "Executive")
- `email_status` (string): Email verification status (e.g., "verified", "unverified")
- `technologies` (array of strings): Technologies/skills

**URLs**:
- `linkedin_url` (string): LinkedIn profile URL
- `website` (string): Personal website URL

**Validation Rules**:
- `email`: Must be valid email format
- `phone_number`: Must be valid phone format if provided
- `linkedin_url`, `website`: Must be valid URL format if provided

### Partial Updates

When updating records, you can provide only the fields you want to update. All other fields remain unchanged.

**Example - Update Company**:
```json
{
  "employees_count": 150,
  "annual_revenue": 7500000
}
```

Only `employees_count` and `annual_revenue` will be updated; all other fields remain unchanged.

**Example - Update Contact**:
```json
{
  "title": "Senior Software Engineer",
  "seniority": "Senior"
}
```

Only `title` and `seniority` will be updated.

### Upsert Operations

Upsert operations (create or update) can match records by:
- **Companies**: `uuid` or `normalized_domain`
- **Contacts**: `uuid` or `email`

If a matching record is found, it's updated; otherwise, a new record is created.

**See**: 
- [Company API - Write Operations](../company.md#write-operations)
- [Contact API - Write Operations](../contacts.md#write-operations)
- [CRUD Implementation Plan](../CRUD_IMPLEMENTATION_PLAN.md)

---

## Related Documentation

- [Company Filters Guide](./01-company-filters-complete-guide.md)
- [Contact Filters Guide](./02-contact-filters-complete-guide.md)
- [Combined Filters Guide](./03-combined-filters-guide.md)
- [Examples and Use Cases](./05-examples-use-cases.md)
- [API Reference](./06-api-reference.md)

---

**Last Updated**: 2025-01-XX  
**Version**: 1.2
