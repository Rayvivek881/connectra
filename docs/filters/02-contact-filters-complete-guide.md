# Contact Filters - Complete Guide

**Version**: 1.2  
**Last Updated**: 2025-01-XX

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Filter Structure](#filter-structure)
   - [Base VQL Query Structure](#base-vql-query-structure)
   - [Contact Filter Field Categories](#contact-filter-field-categories)
   - [Denormalized Company Fields](#denormalized-company-fields)
4. [Text Match Filters](#text-match-filters)
   - [Substring Text Search](#substring-text-search)
   - [Text Match Parameters Reference](#text-match-parameters-reference)
5. [Denormalized Company Fields](#denormalized-company-fields)
   - [Filtering Contacts by Company Attributes](#filtering-contacts-by-company-attributes)
   - [Available Denormalized Company Fields](#available-denormalized-company-fields)
6. [Keyword Match Filters](#keyword-match-filters)
7. [Range Query Filters](#range-query-filters)
8. [Combined Filter Patterns](#combined-filter-patterns)
9. [Sorting and Pagination](#sorting-and-pagination)
   - [Sortable Fields](#sortable-fields)
   - [Page-Based Pagination](#page-based-pagination)
   - [Cursor-Based Pagination](#cursor-based-pagination)
10. [Real-World Use Cases](#real-world-use-cases)
11. [Field Reference](#field-reference)

- [Filterable Contact Fields](#filterable-contact-fields-elasticsearch-index)
- [Denormalized Company Fields](#denormalized-company-fields-in-contact-index)
- [Response-Only Fields](#response-only-fields-postgresql-only)
- [Filter Availability](#filter-availability)

12. [Best Practices](#best-practices)
13. [Error Handling](#error-handling)
14. [Relationship with Companies](#relationship-with-companies)
15. [Related Documentation](#related-documentation)

## Overview

The Contact API supports comprehensive filtering using VQL (Vivek Query Language). Filters are organized into three main categories:

- **Text Matches**: Full-text search on text fields (supports exact, shuffle, and substring search types)
- **Keyword Matches**: Exact matching on keyword/array fields
- **Range Queries**: Date range filtering

All filters can be combined using `must` (AND logic) and `must_not` (NOT logic) conditions.

> **Authentication Required**: All examples in this guide show only the JSON request body. When making actual API calls, you must include the `X-API-Key` header for authentication. See [06-api-reference.md](./06-api-reference.md#authentication) for complete HTTP request examples with headers.

### Key Features

- **Denormalized Company Fields**: Filter contacts directly by company attributes using `company_*` prefix fields
- **Company Config**: Populate full company objects (27 fields) in responses using `company_config.populate`
- **Flexible Text Search**: Supports exact, shuffle, and substring search types
- **Ngram Support**: Efficient partial text matching for names and titles
- **Account-Based Filtering**: Single-query filtering combining contact and company criteria

---

## Prerequisites

Before using contact filters, ensure you have:

- ✅ Valid API key configured
- ✅ Understanding of VQL query structure
- ✅ Knowledge of available filterable fields
- ✅ Understanding of denormalized company fields vs company_config
- ✅ Access to the Contact API endpoint

**See**: [API Reference](./06-api-reference.md) for authentication and endpoint details.

## Filter Structure

### Base VQL Query Structure

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

### Contact Filter Field Categories

**Text Search Fields** (use in `text_matches`):

- `first_name` - First name
- `last_name` - Last name
- `title` - Job title
- `city` - City name
- `state` - State/Province
- `country` - Country name
- `linkedin_url` - LinkedIn URL

**Keyword Fields** (use in `keyword_match`):

- `id` - Contact ID
- `company_id` - Company UUID
- `email` - Email address
- `departments` - Departments array
- `mobile_phone` - Mobile phone number
- `email_status` - Email verification status
- `seniority` - Seniority level
- `country` - Country (can be text or keyword)
- `city` - City (can be text or keyword)
- `state` - State (can be text or keyword)

**Range Query Fields**:

- `created_at` - Creation date (ISO 8601 string)

**Denormalized Company Fields** (available in contact index with `company_` prefix):

- `company_name` - Company name (text, with ngram support)
- `company_employees_count` - Company employee count (range)
- `company_industries` - Company industries (keyword array)
- `company_keywords` - Company keywords (keyword array)
- `company_address` - Company address (text)
- `company_annual_revenue` - Company annual revenue (range)
- `company_total_funding` - Company total funding (range)
- `company_technologies` - Company technologies (keyword array)
- `company_city` - Company city (text)
- `company_state` - Company state (text)
- `company_country` - Company country (text)
- `company_linkedin_url` - Company LinkedIn URL (text)
- `company_website` - Company website (text, with ngram support)
- `company_normalized_domain` - Company normalized domain (text, with ngram support)

These denormalized fields allow you to filter contacts directly by their company attributes without needing a separate company query. See [Denormalized Company Fields](#denormalized-company-fields) section for details.

---

## Text Match Filters

### Single Text Match - First Name

**Example: Search by first name**

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

### Single Text Match - Last Name

**Example: Search by last name**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "Smith",
          "filter_key": "last_name",
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

### Single Text Match - Title (Exact)

**Example: Exact job title search**

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
      ]
    }
  },
  "page": 1,
  "limit": 25
}
```

### Single Text Match - Title (Shuffle)

**Example: Flexible job title search**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "engineer developer",
          "filter_key": "title",
          "search_type": "shuffle",
          "operator": "or",
          "fuzzy": true
        }
      ]
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Text Matches - Name Fields

**Example: Search first and last name**

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
        },
        {
          "text_value": "Smith",
          "filter_key": "last_name",
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

### Multiple Text Matches - Title and Name

**Example: Title and name search**

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
        },
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

### Text Match with must_not

**Example: Include "engineer", exclude "intern"**

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
      ],
      "must_not": [
        {
          "text_value": "intern",
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

### Location Text Search - City

**Example: Search by city**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "San Francisco",
          "filter_key": "city",
          "search_type": "exact",
          "slop": 0
        }
      ]
    }
  },
  "page": 1,
  "limit": 25
}
```

### Location Text Search - State

**Example: Search by state**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "California",
          "filter_key": "state",
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

### Location Text Search - Country

**Example: Search by country**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "United States",
          "filter_key": "country",
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

### LinkedIn URL Text Search

**Example: Search by LinkedIn URL**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "linkedin.com/in/john-smith",
          "filter_key": "linkedin_url",
          "search_type": "shuffle",
          "fuzzy": false
        }
      ]
    }
  },
  "page": 1,
  "limit": 25
}
```

### Substring Text Search

**Example: Partial name matching using substring search**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "john",
          "filter_key": "first_name",
          "search_type": "substring",
          "operator": "and"
        }
      ]
    }
  },
  "page": 1,
  "limit": 25
}
```

**Example: Partial title matching**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "eng",
          "filter_key": "title",
          "search_type": "substring",
          "operator": "and"
        }
      ]
    }
  },
  "page": 1,
  "limit": 25
}
```

**Use Cases for Substring Search**:

- Finding contacts with partial name matches (e.g., "johnn" matches "Johnny", "Johnson")
- Searching for partial words within job titles (e.g., "engin" matches "Engineer", "Engineering")
- Autocomplete-style searches where users type partial text
- Finding name variations and abbreviations

**Important**:

- Only `first_name`, `last_name`, and `title` fields support substring search (via ngram analysis)
- Minimum search text length is 5 characters (min_gram: 5)
- Maximum substring match is 10 characters (max_gram: 10)
- Other text fields (`city`, `state`, `country`, `linkedin_url`) do not support substring search

**Note**: Substring search uses ngram analysis and automatically queries the `.ngram` field. For contact filters, the `first_name`, `last_name`, and `title` fields support ngram matching (min_gram: 5, max_gram: 10). This enables efficient partial text matching. The search text must be at least 5 characters long for substring search on these fields.

### Text Match Parameters Reference

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `text_value` | string | Yes | Text to search for |
| `filter_key` | string | Yes | Field name (first_name, last_name, title, etc.) |
| `search_type` | string | Yes | `"exact"`, `"shuffle"`, or `"substring"` |
| `slop` | integer | No | Word distance for exact search (default: 0) |
| `operator` | string | No | `"and"` or `"or"` for shuffle/substring search |
| `fuzzy` | boolean | No | Enable fuzzy matching for exact/shuffle (default: false) |

**Search Type Guide**:

- **`"exact"`**: Phrase matching with word order. Use for multi-word phrases where order matters (e.g., "software engineer")
- **`"shuffle"`**: Word matching where order doesn't matter. Use for flexible text search (e.g., "engineer software" matches "software engineer")
- **`"substring"`**: Partial text matching using ngram analysis. Use for finding partial matches within words (e.g., "engin" matches "Engineer", "Engineering"). For contact filters, `first_name`, `last_name`, and `title` support substring search via ngram (min_gram: 5, max_gram: 10). The search text must be at least 5 characters long.

---

## Denormalized Company Fields

> **⚠️ CRITICAL DISTINCTION**: Denormalized `company_*` fields are **ONLY for filtering** in `where` clauses. They are **NOT available** in `select_columns`. To get company data in API responses, use `company_config.select_columns` with direct field names (e.g., `name`, `employees_count`). See [Relationship with Companies](#relationship-with-companies) section for details on `company_config`.

The contact index includes denormalized company data, allowing you to filter contacts directly by their company attributes without performing a separate company query. All denormalized fields use the `company_` prefix.

**Key Points**:
- ✅ **Use in `where` clauses**: Filter contacts by company attributes using `company_*` fields
- ❌ **NOT in `select_columns`**: Denormalized fields cannot be selected in responses
- ✅ **Use `company_config`**: To get company data in responses, use `company_config.populate` with `company_config.select_columns`

### Filtering Contacts by Company Attributes

**Example: Find contacts at companies with specific employee count**

```json
{
  "where": {
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
  "limit": 25
}
```

**Example: Find contacts at companies in specific industries**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS"],
        "seniority": ["Senior", "Lead"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

**Example: Find contacts at companies using specific technologies**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_technologies": ["Python", "Go"],
        "departments": ["Engineering"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

**Example: Find contacts at high-revenue companies**

```json
{
  "where": {
    "range_query": {
      "must": {
        "company_annual_revenue": {
          "gte": 5000000
        }
      }
    },
    "keyword_match": {
      "must": {
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified"
      }
    }
  },
  "page": 1,
  "limit": 25
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
  },
  "page": 1,
  "limit": 25
}
```

**Note**: The `company_name`, `company_website`, and `company_normalized_domain` fields support substring search via ngram (min_gram: 5, max_gram: 10), same as contact name fields.

### Available Denormalized Company Fields

| Field | Type | Search Type | Description |
|-------|------|-------------|-------------|
| `company_name` | text | text (ngram: 5-10) | Company name |
| `company_employees_count` | integer | range | Company employee count |
| `company_industries` | array[string] | keyword | Company industries |
| `company_keywords` | array[string] | keyword | Company keywords |
| `company_address` | string | text | Company address |
| `company_annual_revenue` | integer | range | Company annual revenue (in cents) |
| `company_total_funding` | integer | range | Company total funding (in cents) |
| `company_technologies` | array[string] | keyword | Company technologies |
| `company_city` | string | text | Company city |
| `company_state` | string | text | Company state |
| `company_country` | string | text | Company country |
| `company_linkedin_url` | string | text | Company LinkedIn URL |
| `company_website` | string | text (ngram: 5-10) | Company website |
| `company_normalized_domain` | string | text (ngram: 5-10) | Company normalized domain |

**Benefits of Denormalized Fields**:

- Single query to filter contacts by company attributes
- No need to first query companies, then filter contacts by `company_id`
- Faster queries for account-based filtering scenarios
- Enables complex filtering combining contact and company criteria in one request

---

## Keyword Match Filters

### Single Keyword - Departments

**Example: Single department**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Departments

**Example: Multiple departments**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering", "Sales", "Marketing"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Single Keyword - Email Status

**Example: Verified email status**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Email Status Values

**Available values**:

- `"verified"` - Email address has been verified
- `"unverified"` - Email address not yet verified
- `"invalid"` - Email address is invalid
- `"bounced"` - Email has bounced

### Single Keyword - Seniority

**Example: Senior level**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": "Senior"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Seniority Values

**Available values**:

- `"Junior"` - Junior level
- `"Mid"` - Mid-level
- `"Senior"` - Senior level
- `"Lead"` - Lead/Team Lead
- `"Principal"` - Principal level
- `"Executive"` - Executive level

### Multiple Seniority Levels

**Example: Senior and Lead**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": ["Senior", "Lead", "Principal"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Country Keyword Filter

**Example: Single country**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "country": "USA"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Countries

**Example: Multiple countries**

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

### Company ID Filter

**Example: Single company**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": "c0a8012e-1111-2222-3333-444455556666"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Company IDs

**Example: Multiple companies**

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
  "page": 1,
  "limit": 25
}
```

### Email Filter

**Example: Specific email**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email": "john.smith@example.com"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Mobile Phone Filter

**Example: Specific phone number**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "mobile_phone": "4706037761"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Keyword Filters Combined

**Example: Departments, seniority, and country**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering", "Sales"],
        "seniority": "Senior",
        "country": ["USA", "Canada"],
        "email_status": "verified"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Keyword Match with must_not

**Example: Include departments, exclude seniority**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering"],
        "country": ["USA"]
      },
      "must_not": {
        "seniority": ["Junior", "Intern"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Contact ID Filter

**Example: Specific contact ID**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": 43171040
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Contact IDs

**Example: Multiple contact IDs**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [43171040, 43171041, 43171042]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

---

## Range Query Filters

### Date Range - Created At

**Example: Contacts created in date range**

```json
{
  "where": {
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2023-01-01T00:00:00Z",
          "lte": "2024-12-31T23:59:59Z"
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Date After

**Example: Contacts created after date**

```json
{
  "where": {
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2024-01-01T00:00:00Z"
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Date Before

**Example: Contacts created before date**

```json
{
  "where": {
    "range_query": {
      "must": {
        "created_at": {
          "lte": "2023-12-31T23:59:59Z"
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Date Range Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `gte` | Greater than or equal | `{"gte": "2024-01-01T00:00:00Z"}` |
| `lte` | Less than or equal | `{"lte": "2024-12-31T23:59:59Z"}` |
| `gt` | Greater than | `{"gt": "2024-01-01T00:00:00Z"}` |
| `lt` | Less than | `{"lt": "2024-12-31T23:59:59Z"}` |

**Date Format**: ISO 8601 (RFC3339) - `"2024-01-01T00:00:00Z"`

---

## Combined Filter Patterns

### Text Match + Keyword Match

**Example: Name search + department filter**

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
    },
    "keyword_match": {
      "must": {
        "departments": ["Engineering"],
        "email_status": "verified"
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Text Match + Range Query

**Example: Title search + date filter**

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
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2023-01-01T00:00:00Z"
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Keyword Match + Range Query

**Example: Department + date filter**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering"],
        "seniority": "Senior",
        "country": ["USA"]
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
  "page": 1,
  "limit": 25
}
```

### All Three Filter Types Combined

**Example: Complete filter combination**

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
        "country": ["USA", "Canada"],
        "email_status": "verified"
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

---

## Sorting and Pagination

### Single Sort

**Example: Sort by creation date**

```json
{
  "where": {
    "keyword_match": {
      "must": {
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
  "limit": 25
}
```

### Multiple Sorts

**Example: Sort by email then creation date**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "country": ["USA"]
      }
    }
  },
  "order_by": [
    {
      "order_by": "email",
      "order_direction": "asc"
    },
    {
      "order_by": "created_at",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 25
}
```

### Sortable Fields

**Can sort by**:

- `id`
- `company_id`
- `email`
- `departments` (keyword field)
- `mobile_phone`
- `email_status`
- `seniority`
- `created_at`

**Cannot sort by** (text fields):

- `first_name`
- `last_name`
- `title`
- `city`
- `state`
- `country`
- `linkedin_url`

**Important**: Text fields (`first_name`, `last_name`, `title`, `city`, `state`, `country`, `linkedin_url`) cannot be used for sorting in Elasticsearch because they are analyzed fields. Use keyword fields (`email`, `company_id`, `seniority`, `email_status`, etc.) or date fields (`created_at`) for sorting operations.

### Page-Based Pagination

**Example: First page**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

**Pagination Limits**:

- `page`: Maximum 10
- `limit`: Maximum 100, default 25

### Cursor-Based Pagination

**Example: Using search_after**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering"]
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
  "search_after": ["2024-01-15T08:00:00Z", "john.smith@example.com"],
  "limit": 25
}
```

**Note**: `search_after` values come from the last document in the previous response. Use the sort field values in the same order as `order_by`.

---

## Real-World Use Cases

### Email Campaign Targeting

**Find verified, active contacts in specific countries**

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

### Recruiting Search

**Find senior engineers in specific technologies**

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

### Account-Based Marketing

**Find contacts at target companies**

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

### Lead Qualification

**Find high-quality leads**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "director manager",
          "filter_key": "title",
          "search_type": "shuffle",
          "operator": "or",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified",
        "country": ["USA", "Canada", "UK"]
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
  "order_by": [
    {
      "order_by": "seniority",
      "order_direction": "desc"
    },
    {
      "order_by": "created_at",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 50
}
```

### Sales Outreach

**Find decision-makers in target departments**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Sales", "Marketing", "Business Development"],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified"
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
      "order_by": "seniority",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 100
}
```

---

## Field Reference

### Filterable Contact Fields (Elasticsearch Index)

These fields are indexed in Elasticsearch and can be used in filter queries:

| Field | Type | Search Type | Sortable | Ngram Support | Description |
|-------|------|-------------|----------|---------------|-------------|
| `id` | integer | keyword | Yes | No | Contact ID |
| `first_name` | string | text | No | Yes (5-10) | First name |
| `last_name` | string | text | No | Yes (5-10) | Last name |
| `company_id` | string | keyword | Yes | No | Company UUID |
| `email` | string | keyword | Yes | No | Email address |
| `title` | string | text | No | Yes (5-10) | Job title |
| `departments` | array[string] | keyword | Yes | No | Departments array |
| `mobile_phone` | string | keyword | Yes | No | Mobile phone number |
| `email_status` | string | keyword | Yes | No | Email verification status |
| `seniority` | string | keyword | Yes | No | Seniority level |
| `city` | string | text | No | No | City name |
| `state` | string | text | No | No | State/Province |
| `country` | string | text | No | No | Country name |
| `linkedin_url` | string | text | No | No | LinkedIn URL |
| `created_at` | datetime | range | Yes | No | Creation date (ISO 8601) |

**Ngram Support**: The `first_name`, `last_name`, and `title` fields have ngram analysis enabled (min_gram: 5, max_gram: 10), allowing efficient substring/partial matching. Other text fields do not support substring search.

**Note**: `city`, `state`, and `country` are text fields in Elasticsearch. They can be used in `text_matches` for fuzzy/flexible search, but for exact matching, use `keyword_match` with the exact values.

### Denormalized Company Fields (in Contact Index)

These company fields are denormalized into the contact index with the `company_` prefix:

| Field | Type | Search Type | Sortable | Ngram Support | Description |
|-------|------|-------------|----------|---------------|-------------|
| `company_name` | string | text | No | Yes (5-10) | Company name |
| `company_employees_count` | integer | range | No | No | Company employee count |
| `company_industries` | array[string] | keyword | No | No | Company industries |
| `company_keywords` | array[string] | keyword | No | No | Company keywords |
| `company_address` | string | text | No | No | Company address |
| `company_annual_revenue` | integer | range | No | No | Company annual revenue (in cents) |
| `company_total_funding` | integer | range | No | No | Company total funding (in cents) |
| `company_technologies` | array[string] | keyword | No | No | Company technologies |
| `company_city` | string | text | No | No | Company city |
| `company_state` | string | text | No | No | Company state |
| `company_country` | string | text | No | No | Company country |
| `company_linkedin_url` | string | text | No | No | Company LinkedIn URL |
| `company_website` | string | text | No | Yes (5-10) | Company website |
| `company_normalized_domain` | string | text | No | Yes (5-10) | Company normalized domain |

**Ngram Support**: `company_name`, `company_website`, and `company_normalized_domain` support substring search via ngram (min_gram: 5, max_gram: 10).

### Response-Only Fields (PostgreSQL Only)

> **⚠️ Important**: These fields are stored in PostgreSQL but are **NOT** indexed in Elasticsearch and **cannot be used in filters**. They are only returned in API responses:

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

### Filter Availability

**Direct-Derived Filters** (`direct_derived: true`):

- `company_id`
- `email`
- `first_name`
- `last_name`
- `linkedin_url`
- `mobile_phone`

**Stored Filters** (`direct_derived: false`):

- `city`
- `country`
- `departments`
- `email_status`
- `seniority`
- `state`
- `title`

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

**departments** (keyword array field):
Common values include: `"Engineering"`, `"Sales"`, `"Customer Success"`, `"Support"`, `"HR"`, `"Marketing"`, `"Operations"`, `"Legal"`, `"Finance"`, `"Product"`, `"Research"`, etc.

---

## Best Practices

### 1. Name and Title Searches

**✅ DO**:

- Use `shuffle` search type for flexible name matching
- Enable `fuzzy` for typo tolerance
- Search both `first_name` and `last_name` separately for better results
- Use `substring` for partial title matching (e.g., "eng" matches "Engineer", "Engineering")
- Use `exact` with appropriate `slop` for job titles with multiple words

**❌ DON'T**:

- Use `substring` for full word matching (use `shuffle` instead)
- Forget to enable `fuzzy` for user-generated search terms
- Use `exact` without `slop` for flexible multi-word searches

### 2. Email Status and Data Quality

**✅ DO**:

- Always filter by `email_status: "verified"` for high-quality contacts
- Use `seniority` filter to target decision-makers
- Combine with company filters for account-based approaches

**❌ DON'T**:

- Ignore email verification status
- Skip seniority filters when targeting decision-makers

### 3. Pagination Strategy

**✅ DO**:

- Use `search_after` for large result sets (> 250 records)
- Use page-based pagination for small datasets (< 250 records)
- Keep `limit` reasonable (25-50 for best performance)
- Always include `id` as last sort field for stable pagination

**❌ DON'T**:

- Use page-based pagination beyond page 10
- Use `search_after` without `order_by`

### 4. Performance Optimization

**✅ DO**:

- Use count endpoint when you only need the total
- Combine multiple keyword filters in a single `must` object
- Use stored filters (`direct_derived: false`) when available
- Use `select_columns` to limit returned fields and reduce payload size
- Use denormalized company fields for single-query filtering

**❌ DON'T**:

- Request all fields when you only need a few
- Make separate company queries when denormalized fields would work
- Ignore filter performance characteristics

### 5. Denormalized Company Fields

**✅ DO**:

- Use `company_*` fields for filtering contacts by company attributes
- Filter by denormalized fields in `where` clauses
- Use `company_config.populate` to get company data in responses

**❌ DON'T**:

- Use `company_*` fields in `select_columns` (they're for filtering only)
- Use `company_*` prefix in `company_config.select_columns` (use direct names)
- Forget that denormalized fields are only for filtering

### 6. Filter Combination

**✅ DO**:

- Combine text, keyword, and range filters for precise results
- Use `must` for AND logic (all conditions must match)
- Use `must_not` to exclude matching records
- Combine contact and company filters using denormalized fields

**❌ DON'T**:

- Create overly complex nested filters
- Mix incompatible filter types incorrectly

---

## Error Handling

### Common Errors and Solutions

#### 400 Bad Request - Invalid Request Body

**Error Response**:

```json
{
  "error": "ERR_INVALID_REQUEST_BODY: the request body is invalid; check JSON syntax and required fields",
  "success": false
}
```

**Common Causes**:

- Malformed JSON in request body
- Missing required fields
- Invalid field types
- Using response-only fields in filters

**Solution**:

- Validate JSON syntax
- Ensure all required fields are present
- Check field types match expected values
- Verify you're not using response-only fields in `where` clauses

#### 400 Bad Request - Page Size Exceeded

**Error Response**:

```json
{
  "error": "ERR_PAGE_SIZE_EXCEEDED: the requested page size surpasses the maximum allowed limit; consider using pagination with smaller batches",
  "success": false
}
```

**Cause**: `limit` parameter exceeds maximum value (100)

**Solution**: Reduce `limit` to maximum 100

#### 400 Bad Request - Page Number Exceeded

**Error Response**:

```json
{
  "error": "ERR_PAGE_OUT_OF_RANGE: the requested page number is beyond the available range; verify total pages before requesting",
  "success": false
}
```

**Cause**: `page` parameter exceeds maximum value (10)

**Solution**: Use `search_after` for pagination beyond page 10

#### 500 Internal Server Error - Elasticsearch Error

**Error Response**:

```json
{
  "error": "ERR_ELASTICSEARCH_FAILURE: search engine returned status 400; details: ...",
  "success": false
}
```

**Common Causes**:

- Invalid field names in filters
- Using denormalized fields incorrectly
- Elasticsearch cluster unavailable

**Solution**:

- Verify field names using `/filters` endpoint
- Check that denormalized fields use `company_*` prefix in `where` clauses
- Verify Elasticsearch cluster health

**See**: [Error Handling Guide](./06-api-reference.md#error-handling) for complete error reference and troubleshooting.

---

## Relationship with Companies

Contacts are linked to companies via the `company_id` field. You can:

1. **Filter contacts by company ID**: Use `company_id` in `keyword_match` to get all contacts for specific companies
2. **Filter contacts by company attributes**: Use denormalized `company_*` fields to filter contacts directly by company attributes (see [Denormalized Company Fields](#denormalized-company-fields) section)
3. **Account-based searches**: Use company filters in combination with contact filters for account-based approaches
4. **Join company data**: Join company data on the application side using the `company_id` field

**Two Approaches for Company-Based Contact Filtering**:

**Approach 1: Filter by company_id (requires separate company query)**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": ["company-uuid-1", "company-uuid-2"]
      }
    }
  }
}
```

**Approach 2: Filter by denormalized company fields (single query)**

```json
{
  "where": {
    "range_query": {
      "must": {
        "company_employees_count": {"gte": 100},
        "company_annual_revenue": {"gte": 5000000}
      }
    },
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS"]
      }
    }
  }
}
```

**Example: Find contacts at high-value companies**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "company-uuid-1",
          "company-uuid-2"
        ],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified"
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

### Populating Company Data in Responses (company_config)

> **Key Concept**: While denormalized `company_*` fields are used for **filtering** contacts by company attributes, `company_config` is used to **populate full company objects** in API responses.

The `company_config` parameter allows you to include full company objects alongside contact data in responses. This is different from denormalized fields, which are only for filtering.

**Structure**:
```json
{
  "where": {...},
  "company_config": {
    "populate": true,
    "select_columns": ["name", "employees_count", "industries", "annual_revenue"]
  }
}
```

**Parameters**:
- `populate` (boolean, required): Set to `true` to enable company population
- `select_columns` (array of strings, optional): List of company fields to return (use direct field names, NOT `company_*` prefix)

**Example: Filter contacts and populate company data**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS"],
        "seniority": ["Senior", "Lead"],
        "email_status": "verified"
      }
    }
  },
  "select_columns": ["id", "first_name", "last_name", "title", "company_id"],
  "company_config": {
    "populate": true,
    "select_columns": ["uuid", "name", "employees_count", "industries", "annual_revenue", "website"]
  },
  "page": 1,
  "limit": 25
}
```

**Response Structure**:
```json
{
  "data": [
    {
      "id": 123,
      "first_name": "John",
      "last_name": "Doe",
      "title": "Software Engineer",
      "company_id": "company-uuid-here",
      "company": {
        "uuid": "company-uuid-here",
        "name": "Acme Software Corp",
        "employees_count": 150,
        "industries": ["Software", "SaaS"],
        "annual_revenue": 5000000,
        "website": "https://acme.com"
      }
    }
  ],
  "success": true
}
```

**Key Points**:
- ✅ Use `company_*` fields in `where` clauses for filtering (denormalized fields)
- ✅ Use `company_config.select_columns` with direct field names (e.g., `name`, `employees_count`) to select company data
- ❌ Do NOT use `company_*` prefix in `company_config.select_columns`
- ❌ Denormalized fields are NOT returned in responses (they're filter-only)

**Available Company Fields** (27 total): All company fields from PostgreSQL are available in `company_config.select_columns`. See [Select Columns Guide](./select_columns_filter.md#company-config---populate-company-objects-27-fields) for the complete list.

**See**: [Select Columns Guide](./select_columns_filter.md) for comprehensive documentation on `company_config` and field selection.

---

## Write Operations (CRUD)

> **Status**: The following CRUD operations are documented but currently only `batch-upsert` is implemented. See [CRUD Implementation Plan](../CRUD_IMPLEMENTATION_PLAN.md) for implementation details.

### Create Contact

Create a new contact record with automatic Elasticsearch indexing.

**Endpoint**: `POST /contacts/create`

**Request**:
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "email": "john.doe@acme.com",
    "first_name": "John",
    "last_name": "Doe",
    "title": "Software Engineer",
    "company_id": "company-uuid-here",
    "seniority": "Senior",
    "email_status": "verified"
  }'
```

**Response** (201 Created):
```json
{
  "data": {
    "uuid": "contact-uuid-here",
    "email": "john.doe@acme.com",
    "first_name": "John",
    "last_name": "Doe",
    "title": "Software Engineer",
    "company_id": "company-uuid-here",
    "seniority": "Senior",
    "email_status": "verified",
    "created_at": "2025-12-24T10:30:00Z"
  },
  "success": true
}
```

### Update Contact

Update an existing contact by UUID.

**Endpoint**: `PUT /contacts/:uuid`

**Request**:
```bash
curl -X PUT https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/contact-uuid-here \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "title": "Senior Software Engineer",
    "seniority": "Senior"
  }'
```

### Delete Contact

Soft delete a contact by UUID.

**Endpoint**: `DELETE /contacts/:uuid`

**Request**:
```bash
curl -X DELETE https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/contact-uuid-here \
  -H "X-API-Key: your-secret-api-key"
```

### Upsert Contact

Create or update a contact (identified by UUID or email).

**Endpoint**: `POST /contacts/upsert`

**Request**:
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "email": "john.doe@acme.com",
    "first_name": "John",
    "last_name": "Doe",
    "title": "Software Engineer"
  }'
```

### Bulk Upsert Contacts

Efficiently create or update multiple contacts.

**Endpoint**: `POST /contacts/batch-upsert` (Currently Implemented)

**Request**:
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/batch-upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "pg_contacts": [
      {
        "uuid": "contact-uuid-here",
        "email": "john.doe@acme.com",
        "first_name": "John",
        "last_name": "Doe"
      }
    ],
    "es_contacts": [
      {
        "uuid": "contact-uuid-here",
        "email": "john.doe@acme.com",
        "first_name": "John",
        "last_name": "Doe"
      }
    ]
  }'
```

**See**: [Contact API - Write Operations](../contacts.md#write-operations) for complete CRUD documentation.

---

## Related Documentation

### Filter Documentation

- [Company Filters Guide](./01-company-filters-complete-guide.md) - Complete company filtering guide
- [Combined Filters Guide](./03-combined-filters-guide.md) - Account-based filtering strategies
- [Filter Field Reference](./04-filter-field-reference.md) - Complete field reference
- [Examples and Use Cases](./05-examples-use-cases.md) - Real-world examples
- [API Reference](./06-api-reference.md) - Complete API endpoint reference
- [Select Columns Guide](./select_columns_filter.md) - Field selection and company_config

### Main Documentation

- [System Documentation](../system.md) - System architecture and setup
- [Company API](../company.md) - Company API documentation
- [Contact API](../contacts.md) - Contact API documentation
- [Main README](../README.md) - Documentation index

---

**Last Updated**: 2025-01-XX  
**Version**: 1.2
