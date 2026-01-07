# Select Columns - Complete Guide for Contact Filters

## Table of Contents

1. [Overview](#overview)
2. [How select_columns Works](#how-select_columns-works)
3. [Complete Field Reference](#complete-field-reference)
4. [Contact Core Fields (16 fields)](#contact-core-fields-16-fields)
5. [Contact Metadata Fields (9 fields)](#contact-metadata-fields-9-fields)
6. [Denormalized Company Fields (13 fields)](#denormalized-company-fields-13-fields)
7. [Company Config - Populate Company Objects (27 fields)](#company-config---populate-company-objects-27-fields)
8. [Denormalized vs Company Config Comparison](#denormalized-vs-company-config-comparison)
9. [Usage Patterns](#usage-patterns)
10. [Complete Examples](#complete-examples)
11. [Best Practices](#best-practices)
12. [Common Errors](#common-errors)
13. [Authentication](#authentication)

---

## Overview

The `select_columns` parameter allows you to specify which fields should be returned in the API response. This is a powerful optimization tool that:

- **Reduces response payload size** - Only fetch the data you need
- **Improves performance** - Less data transfer means faster responses
- **Enables flexible data selection** - Mix filterable and response-only fields
- **Supports different view modes** - Simple list views, full detail views, or exports

**Total Available Fields**: 25 contact fields + 27 company fields (via `company_config`)
- **16 Contact Core Fields** - Filterable in Elasticsearch, selectable in `select_columns`
- **9 Contact Metadata Fields** - Response-only (PostgreSQL only), selectable in `select_columns`
- **13 Denormalized Company Fields** - **ONLY for filtering** in `where` clauses (NOT in `select_columns`)
- **27 Company Fields** - Available via `company_config.select_columns` (full company objects)

> **⚠️ Important**: `select_columns` only affects PostgreSQL field retrieval **after** the Elasticsearch search completes. It does NOT affect:
> - Elasticsearch search performance
> - Which documents are matched by the search
> - Filter query execution

### ⚠️ CRITICAL: Denormalized Fields vs Company Config Fields

**Key Distinction:**

| Aspect | Denormalized Fields | Company Config Fields |
|--------|---------------------|----------------------|
| **Usage** | **ONLY for filtering** in `where` clauses | **ONLY for selecting** in `company_config.select_columns` |
| **Location** | In `where` clauses (filtering) | In `company_config.select_columns` (selection) |
| **Field Names** | Use `company_*` prefix | Use direct names (NO prefix) |
| **Examples** | `company_name`, `company_employees_count` (in `where`) | `name`, `employees_count` (in `company_config.select_columns`) |
| **Response Location** | ❌ NOT returned in response | ✅ In nested `company` object |
| **Performance** | ⚡ Fast (already in index) | 🐢 Slower (separate query) |
| **Fields Available** | 13 fields (for filtering) | 27 fields (for selection) |

> **⚠️ IMPORTANT**: Denormalized `company_*` fields are **ONLY for filtering** in `where` clauses. They are **NOT available** in `select_columns`. To get company data in responses, you **MUST** use `company_config.select_columns`.

**❌ WRONG - Using denormalized fields in select_columns:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software"]  // ✅ OK - filtering
      }
    }
  },
  "select_columns": [
    "id",
    "company_name"  // ❌ ERROR: Denormalized fields NOT in select_columns!
  ]
}
```

**✅ CORRECT - Using company_config for selection:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software"]  // ✅ OK - filtering with denormalized
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name"
  ],
  "company_config": {
    "populate": true,
    "select_columns": ["name", "employees_count"]  // ✅ Correct!
  }
}
```

---

## How select_columns Works

### Execution Flow

1. **Elasticsearch Search**: The filter query executes in Elasticsearch and returns matching document IDs
2. **PostgreSQL Retrieval**: PostgreSQL retrieves full records for those document IDs
3. **Field Selection**: If `select_columns` is specified, only those fields are returned from PostgreSQL
4. **Response**: The API returns only the selected fields in the response

### When to Use select_columns

- ✅ **Simple List Views**: Select only essential fields (8-10 fields)
- ✅ **Detail Views**: Select comprehensive field set (20-30 fields)
- ✅ **Export Mode**: Omit `select_columns` to get all fields
- ✅ **Performance Optimization**: Reduce payload size for large result sets
- ✅ **Mobile/API Clients**: Minimize data transfer

### When NOT to Use select_columns

- ❌ **When you need all fields**: Omit `select_columns` parameter
- ❌ **For filtering**: `select_columns` doesn't affect which documents match
- ❌ **For search performance**: Doesn't improve Elasticsearch query speed

---

## Complete Field Reference

### Summary Statistics

| Category | Count | Filterable | Selectable | Sortable |
|----------|-------|------------|------------|----------|
| Contact Core Fields | 16 | ✅ Yes | ✅ Yes | 8 fields |
| Contact Metadata Fields | 9 | ❌ No | ✅ Yes | 2 fields |
| Denormalized Company Fields | 13 | ✅ Yes | ❌ No* | ❌ No |
| Company Config Fields | 27 | ❌ No** | ✅ Yes | ❌ No |
| **Total Contact Fields** | **25** | **16** | **25** | **10** |
| **Total Company Fields** | **27** | **17** | **27** | **❌ No** |

*Denormalized fields are **ONLY for filtering** in `where` clauses, NOT in `select_columns`
**Company Config fields cannot be used in contact filters, but can be used in company filters when querying companies directly

### Quick Reference Table

| Field | Category | Filterable | Sortable | Ngram | Type |
|-------|----------|------------|----------|-------|------|
| `id` | Core | ✅ | ✅ | ❌ | integer |
| `uuid` | Metadata | ❌ | ❌ | ❌ | string |
| `first_name` | Core | ✅ | ❌ | ✅ (5-10) | string |
| `last_name` | Core | ✅ | ❌ | ✅ (5-10) | string |
| `company_id` | Core | ✅ | ✅ | ❌ | string |
| `email` | Core | ✅ | ✅ | ❌ | string |
| `title` | Core | ✅ | ❌ | ✅ (5-10) | string |
| `departments` | Core | ✅ | ✅ | ❌ | array[string] |
| `mobile_phone` | Core | ✅ | ✅ | ❌ | string |
| `email_status` | Core | ✅ | ✅ | ❌ | string |
| `seniority` | Core | ✅ | ✅ | ❌ | string |
| `city` | Core | ✅ | ❌ | ❌ | string |
| `state` | Core | ✅ | ❌ | ❌ | string |
| `country` | Core | ✅ | ❌ | ❌ | string |
| `linkedin_url` | Core | ✅ | ❌ | ❌ | string |
| `created_at` | Core | ✅ | ✅ | ❌ | datetime |
| `work_direct_phone` | Metadata | ❌ | ❌ | ❌ | string |
| `home_phone` | Metadata | ❌ | ❌ | ❌ | string |
| `other_phone` | Metadata | ❌ | ❌ | ❌ | string |
| `facebook_url` | Metadata | ❌ | ❌ | ❌ | string |
| `twitter_url` | Metadata | ❌ | ❌ | ❌ | string |
| `website` | Metadata | ❌ | ❌ | ❌ | string |
| `stage` | Metadata | ❌ | ❌ | ❌ | string |
| `updated_at` | Metadata | ❌ | ✅ | ❌ | datetime |
| `deleted_at` | Metadata | ❌ | ❌ | ❌ | datetime |
| `company_name` | Denormalized | ✅ | ❌* | ✅ (5-10) | string |
| `company_employees_count` | Denormalized | ✅ | ❌* | ❌ | integer |
| `company_industries` | Denormalized | ✅ | ❌* | ❌ | array[string] |
| `company_keywords` | Denormalized | ✅ | ❌* | ❌ | array[string] |
| `company_address` | Denormalized | ✅ | ❌* | ❌ | string |
| `company_annual_revenue` | Denormalized | ✅ | ❌* | ❌ | integer |
| `company_total_funding` | Denormalized | ✅ | ❌* | ❌ | integer |
| `company_technologies` | Denormalized | ✅ | ❌* | ❌ | array[string] |
| `company_city` | Denormalized | ✅ | ❌* | ❌ | string |
| `company_state` | Denormalized | ✅ | ❌* | ❌ | string |
| `company_country` | Denormalized | ✅ | ❌* | ❌ | string |
| `company_linkedin_url` | Denormalized | ✅ | ❌* | ❌ | string |
| `company_website` | Denormalized | ✅ | ❌* | ✅ (5-10) | string |
| `company_normalized_domain` | Denormalized | ✅ | ❌* | ✅ (5-10) | string |

*Denormalized fields are **ONLY for filtering** in `where` clauses. Use `company_config.select_columns` with direct field names (e.g., `name`, `employees_count`) to select company data.

---

## Contact Core Fields (16 fields)

These fields are indexed in Elasticsearch and can be used in `where` clauses for filtering. They can also be selected using `select_columns`.

### Field List

1. `id` - Contact ID (integer, sortable)
2. `uuid` - Contact UUID (string, response-only)
3. `first_name` - First name (text, ngram support 5-10)
4. `last_name` - Last name (text, ngram support 5-10)
5. `company_id` - Company UUID (keyword, sortable)
6. `email` - Email address (keyword, sortable)
7. `title` - Job title (text, ngram support 5-10)
8. `departments` - Departments array (keyword array, sortable)
9. `mobile_phone` - Mobile phone (keyword, sortable)
10. `email_status` - Email verification status (keyword, sortable)
11. `seniority` - Seniority level (keyword, sortable)
12. `city` - City name (text)
13. `state` - State/Province (text)
14. `country` - Country name (text)
15. `linkedin_url` - LinkedIn URL (text)
16. `created_at` - Creation date (datetime, sortable)

### Examples

#### Example 1: Select All Core Fields

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
    "uuid",
    "first_name",
    "last_name",
    "company_id",
    "email",
    "title",
    "departments",
    "mobile_phone",
    "email_status",
    "seniority",
    "city",
    "state",
    "country",
    "linkedin_url",
    "created_at"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 2: Select Essential Core Fields Only

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
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "company_id",
    "mobile_phone",
    "email_status"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 3: Filter and Select by ID

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [43171040, 43171041, 43171042]
      }
    }
  },
  "select_columns": ["id", "first_name", "last_name", "email", "title"]
}
```

#### Example 4: Filter by Name and Select Core Fields

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
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "departments",
    "seniority"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 5: Filter by Date Range and Select

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
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "created_at"
  ],
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

---

## Contact Metadata Fields (9 fields)

> **⚠️ Critical**: These fields are **NOT** indexed in Elasticsearch and **CANNOT be used in `where` clauses**. They can only be selected using `select_columns` after filtering by other fields.

### Field List

1. `uuid` - Contact UUID (string)
2. `work_direct_phone` - Work direct phone number (string)
3. `home_phone` - Home phone number (string)
4. `other_phone` - Other phone number (string)
5. `facebook_url` - Facebook profile URL (string)
6. `twitter_url` - Twitter profile URL (string)
7. `website` - Personal website URL (string)
8. `stage` - Contact stage/status (string)
9. `updated_at` - Last update timestamp (datetime, sortable)
10. `deleted_at` - Soft delete timestamp (datetime, null if active)

### Examples

#### Example 1: Select All Metadata Fields

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified",
        "seniority": ["Senior", "Lead", "Principal"]
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "uuid",
    "work_direct_phone",
    "home_phone",
    "other_phone",
    "facebook_url",
    "twitter_url",
    "website",
    "stage",
    "updated_at",
    "deleted_at"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 2: Select Phone Numbers Only

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Sales", "Business Development"],
        "email_status": "verified"
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "mobile_phone",
    "work_direct_phone",
    "home_phone",
    "other_phone"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 3: Select Social Media URLs

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": ["Executive", "Principal"]
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "linkedin_url",
    "facebook_url",
    "twitter_url",
    "website"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 4: Select Contact Stage and Metadata

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified"
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "stage",
    "updated_at",
    "deleted_at"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 5: Select UUID for Reference

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [43171040, 43171041]
      }
    }
  },
  "select_columns": [
    "id",
    "uuid",
    "first_name",
    "last_name",
    "email"
  ]
}
```

---

## Denormalized Company Fields (13 fields) - FOR FILTERING ONLY

These fields are denormalized into the contact index with the `company_` prefix. They allow **filtering contacts by company attributes** in a single query without needing a separate company lookup.

> **⚠️ CRITICAL**: 
> - **Denormalized fields** (`company_*`) are **ONLY for filtering** in `where` clauses
> - **Denormalized fields** are **NOT available** in `select_columns` - they will NOT be returned in responses
> - To get company data in responses, you **MUST** use `company_config.select_columns` with direct field names (e.g., `name`, `employees_count`)
> - **DO NOT** use `company_*` prefix in `company_config.select_columns` - it will cause errors!

### Field List

1. `company_name` - Company name (text, ngram support 5-10)
2. `company_employees_count` - Employee count (integer, range)
3. `company_industries` - Industries (keyword array)
4. `company_keywords` - Keywords (keyword array)
5. `company_address` - Company address (text)
6. `company_annual_revenue` - Annual revenue in cents (integer, range)
7. `company_total_funding` - Total funding in cents (integer, range)
8. `company_technologies` - Technologies (keyword array)
9. `company_city` - Company city (text)
10. `company_state` - Company state (text)
11. `company_country` - Company country (text)
12. `company_linkedin_url` - Company LinkedIn URL (text)
13. `company_website` - Company website (text, ngram support 5-10)
14. `company_normalized_domain` - Normalized domain (text, ngram support 5-10)

### Examples - Filtering Only

> **Note**: These examples show how to **filter** using denormalized fields. To **select** company data in responses, see the [Company Config section](#company-config---populate-company-objects-27-fields).

#### Example 1: Filter by Company Industries

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": ["Senior", "Lead"],
        "company_industries": ["Software", "SaaS"]  // ✅ Filtering with denormalized field
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
      "keywords",
      "address",
      "annual_revenue",
      "total_funding",
      "technologies",
      "city",
      "state",
      "country",
      "linkedin_url",
      "website",
      "normalized_domain"
    ]
  },
  "page": 1,
  "limit": 25
}
```

#### Example 2: Filter by Company Size

```json
{
  "where": {
    "range_query": {
      "must": {
        "company_employees_count": {  // ✅ Filtering with denormalized field
          "gte": 100,
          "lte": 1000
        },
        "company_annual_revenue": {  // ✅ Filtering with denormalized field
          "gte": 5000000
        }
      }
    },
    "keyword_match": {
      "must": {
        "seniority": ["Senior", "Lead", "Principal"]
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
      "annual_revenue",
      "city",
      "state",
      "country"
    ]
  },
  "page": 1,
  "limit": 25
}
```

#### Example 3: Filter by Company Technologies

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_technologies": ["Python", "Go", "JavaScript"],  // ✅ Filtering with denormalized field
        "departments": ["Engineering"]
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
      "technologies",
      "website"
    ]
  },
  "page": 1,
  "limit": 25
}
```

#### Example 4: Filter by Company Name (Substring)

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "micros",
          "filter_key": "company_name",  // ✅ Filtering with denormalized field
          "search_type": "substring",
          "operator": "and"
        }
      ]
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
      "industries"
    ]
  },
  "page": 1,
  "limit": 25
}
```

#### Example 5: Filter by Company Location

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "San Francisco",
          "filter_key": "company_city",  // ✅ Filtering with denormalized field
          "search_type": "exact",
          "slop": 0
        },
        {
          "text_value": "California",
          "filter_key": "company_state",  // ✅ Filtering with denormalized field
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
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
      "name",
      "city",
      "state",
      "country",
      "address"
    ]
  },
  "page": 1,
  "limit": 25
}
```

---

## Company Config - Populate Company Objects (27 fields)

The `company_config` feature allows you to populate **full company objects** from PostgreSQL alongside contact data. This is different from denormalized `company_*` fields:

- **Denormalized fields** (`company_*`): Already in contact index, can filter by them, limited fields
- **Company Config** (`company_config.populate`): Fetches full company objects separately, all 27 fields available

### How company_config Works

1. **Elasticsearch Search**: Filter contacts (can use denormalized `company_*` fields)
2. **Extract Company IDs**: Get `company_id` values from matched contacts
3. **Parallel Fetch**: Fetch full company records from PostgreSQL (in parallel with contacts)
4. **Attach to Response**: Company objects are attached to each contact in the response

### When to Use company_config

✅ **Use `company_config.populate` when:**
- You need full company objects (not just denormalized fields)
- You need company metadata fields (e.g., `phone_number`, `latest_funding`, `facebook_url`)
- You need all company fields for detail views
- You want structured company data separate from contact data

❌ **Use denormalized `company_*` fields when:**
- You only need basic company info (name, employees, revenue)
- You want to filter by company attributes in a single query
- Performance is critical (denormalized is faster)

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

### Complete Field Reference (27 fields)

#### Field Name Mapping: Denormalized vs Company Config

| Denormalized Field (in `select_columns`) | Company Config Field (in `company_config.select_columns`) | Available in Denormalized? |
|------------------------------------------|----------------------------------------------------------|----------------------------|
| `company_name` | `name` | ✅ Yes |
| `company_employees_count` | `employees_count` | ✅ Yes |
| `company_industries` | `industries` | ✅ Yes |
| `company_keywords` | `keywords` | ✅ Yes |
| `company_address` | `address` | ✅ Yes |
| `company_annual_revenue` | `annual_revenue` | ✅ Yes |
| `company_total_funding` | `total_funding` | ✅ Yes |
| `company_technologies` | `technologies` | ✅ Yes |
| `company_city` | `city` | ✅ Yes |
| `company_state` | `state` | ✅ Yes |
| `company_country` | `country` | ✅ Yes |
| `company_linkedin_url` | `linkedin_url` | ✅ Yes |
| `company_website` | `website` | ✅ Yes |
| `company_normalized_domain` | `normalized_domain` | ✅ Yes |
| N/A | `id` | ❌ No |
| N/A | `uuid` | ❌ No |
| N/A | `created_at` | ❌ No |
| N/A | `facebook_url` | ❌ No |
| N/A | `twitter_url` | ❌ No |
| N/A | `company_name_for_emails` | ❌ No |
| N/A | `phone_number` | ❌ No |
| N/A | `latest_funding` | ❌ No |
| N/A | `latest_funding_amount` | ❌ No |
| N/A | `last_raised_at` | ❌ No |
| N/A | `updated_at` | ❌ No |
| N/A | `deleted_at` | ❌ No |
| N/A | `linkedin_sales_url` | ❌ No |

**Key Points:**
- 13 fields available in both (with different names)
- 14 fields ONLY available via `company_config` (metadata and additional fields)
- Denormalized fields use `company_*` prefix
- Company Config fields use direct names (no prefix)

#### Core Company Fields (17 fields) - Filterable in Company API

1. `id` - Company ID (bigint, primary key) - **Only in company_config**
2. `uuid` - Company UUID (text, unique) - **Only in company_config**
3. `name` - Company name (text, ngram support 3-10) - Maps to `company_name` (denormalized)
4. `employees_count` - Employee count (bigint, range filterable) - Maps to `company_employees_count` (denormalized)
5. `industries` - Industries array (text[], keyword filterable) - Maps to `company_industries` (denormalized)
6. `keywords` - Keywords array (text[], keyword filterable) - Maps to `company_keywords` (denormalized)
7. `address` - Company address (text) - Maps to `company_address` (denormalized)
8. `annual_revenue` - Annual revenue in cents (bigint, range filterable) - Maps to `company_annual_revenue` (denormalized)
9. `total_funding` - Total funding in cents (bigint, range filterable) - Maps to `company_total_funding` (denormalized)
10. `technologies` - Technologies array (text[], keyword filterable) - Maps to `company_technologies` (denormalized)
11. `city` - City name (text) - Maps to `company_city` (denormalized)
12. `state` - State/Province (text) - Maps to `company_state` (denormalized)
13. `country` - Country name (text) - Maps to `company_country` (denormalized)
14. `linkedin_url` - LinkedIn URL (text) - Maps to `company_linkedin_url` (denormalized)
15. `website` - Website URL (text) - Maps to `company_website` (denormalized)
16. `normalized_domain` - Normalized domain (text) - Maps to `company_normalized_domain` (denormalized)
17. `created_at` - Creation date (timestamp, range filterable) - **Only in company_config**

#### Company Metadata Fields (10 fields) - Response-Only

18. `facebook_url` - Facebook page URL (text) - **Only in company_config**
19. `twitter_url` - Twitter profile URL (text) - **Only in company_config**
20. `company_name_for_emails` - Company name formatted for emails (text) - **Only in company_config**
21. `phone_number` - Company phone number (text) - **Only in company_config**
22. `latest_funding` - Latest funding round (text, e.g., "Series B") - **Only in company_config**
23. `latest_funding_amount` - Latest funding amount in cents (bigint) - **Only in company_config**
24. `last_raised_at` - Date of last funding round (text) - **Only in company_config**
25. `updated_at` - Last update timestamp (timestamp) - **Only in company_config**
26. `deleted_at` - Soft delete timestamp (timestamp, null if active) - **Only in company_config**
27. `linkedin_sales_url` - LinkedIn Sales Navigator URL (text) - **Only in company_config**

### Examples

#### Example 1: Basic Company Population

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "Altitude",
          "filter_key": "company_name",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    }
  },
  "company_config": {
    "populate": true,
    "select_columns": ["uuid", "name", "website"]
  },
  "page": 1,
  "limit": 100
}
```

**Response Structure:**
```json
{
  "data": [
    {
      "id": 43171040,
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "company_id": "c0a8012e-1111-2222-3333-444455556666",
      "company": {
        "uuid": "c0a8012e-1111-2222-3333-444455556666",
        "name": "Altitude Software",
        "website": "https://altitude.com"
      }
    }
  ]
}
```

#### Example 2: Select All Core Company Fields

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": ["Senior", "Lead"],
        "email_status": "verified"
      }
    }
  },
  "company_config": {
    "populate": true,
    "select_columns": [
      "id",
      "uuid",
      "name",
      "employees_count",
      "industries",
      "keywords",
      "address",
      "annual_revenue",
      "total_funding",
      "technologies",
      "city",
      "state",
      "country",
      "linkedin_url",
      "website",
      "normalized_domain",
      "created_at"
    ]
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "company_id"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 3: Select Company Metadata Fields

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Sales", "Marketing"]
      }
    }
  },
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",
      "phone_number",
      "facebook_url",
      "twitter_url",
      "company_name_for_emails",
      "latest_funding",
      "latest_funding_amount",
      "last_raised_at"
    ]
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "company_id"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 4: Complete Company Object (All 27 Fields)

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified"
      }
    }
  },
  "company_config": {
    "populate": true,
    "select_columns": [
      "id",
      "uuid",
      "name",
      "employees_count",
      "industries",
      "keywords",
      "address",
      "annual_revenue",
      "total_funding",
      "technologies",
      "city",
      "state",
      "country",
      "linkedin_url",
      "website",
      "normalized_domain",
      "created_at",
      "facebook_url",
      "twitter_url",
      "company_name_for_emails",
      "phone_number",
      "latest_funding",
      "latest_funding_amount",
      "last_raised_at",
      "updated_at",
      "deleted_at",
      "linkedin_sales_url"
    ]
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "company_id"
  ],
  "page": 1,
  "limit": 25
}
```

#### Example 5: Filter by Denormalized Field, Populate Company Object

```json
{
  "where": {
    "range_query": {
      "must": {
        "company_employees_count": {
          "gte": 100,
          "lte": 1000
        },
        "company_annual_revenue": {
          "gte": 5000000
        }
      }
    },
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS"],
        "seniority": ["Senior", "Lead"]
      }
    }
  },
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",
      "employees_count",
      "annual_revenue",
      "industries",
      "website",
      "phone_number",
      "latest_funding"
    ]
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "company_id"
  ],
  "page": 1,
  "limit": 25
}
```

**Key Point**: You can filter by denormalized `company_*` fields (fast, single query) and then populate full company objects with additional fields not available in denormalized data.

---

## Denormalized vs Company Config Comparison

### Side-by-Side Comparison

| Feature | Denormalized Fields (`company_*`) | Company Config (`company_config.populate`) |
|---------|-----------------------------------|--------------------------------------------|
| **Data Source** | Elasticsearch (denormalized in contact index) | PostgreSQL (full company table) |
| **Query Performance** | ⚡ Fast (single Elasticsearch query) | 🐢 Slower (parallel PostgreSQL query) |
| **Filtering** | ✅ Can filter by them in `where` clauses | ❌ Cannot filter (used for response only) |
| **Selection** | ❌ **NOT available** in `select_columns` | ✅ **ONLY way** to get company data in responses |
| **Fields Available** | 13 limited fields (for filtering) | 27 complete fields (for selection) |
| **Response Structure** | ❌ NOT returned in response | ✅ Nested `company` object |
| **Use Case** | Filter contacts by company attributes | Get full company details in response |
| **Metadata Fields** | ❌ Not available | ✅ Available (phone, funding, social, etc.) |
| **Field Names** | `company_*` prefix (e.g., `company_name`) | Direct names (e.g., `name`) |

### Response Structure Comparison

> **⚠️ IMPORTANT**: Denormalized `company_*` fields are **NOT available** in `select_columns`. They can **ONLY** be used for filtering in `where` clauses. To get company data in responses, you **MUST** use `company_config.select_columns`.

#### Using Company Config (Recommended)

**Query:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software"]
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",
      "employees_count",
      "annual_revenue",
      "phone_number",
      "latest_funding"
    ]
  }
}
```

**Response:**
```json
{
  "data": [
    {
      "id": 43171040,
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "company_id": "c0a8012e-1111-2222-3333-444455556666",
      "company": {
        "uuid": "c0a8012e-1111-2222-3333-444455556666",
        "name": "Altitude Software",
        "employees_count": 500,
        "annual_revenue": 10000000,
        "phone_number": "+1-555-123-4567",
        "latest_funding": "Series B"
      }
    }
  ]
}
```

### When to Use Each Approach

#### Use Denormalized Fields When:
- ✅ **Filtering contacts by company attributes** (fast, single query)
- ✅ Performance is critical for filtering
- ✅ You want to filter without needing company data in response

**Example Use Case - Filtering Only:**
```json
{
  "where": {
    "range_query": {
      "must": {
        "company_employees_count": {"gte": 100}  // ✅ Filtering with denormalized field
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
      "name",
      "employees_count"
    ]
  }
}
```

#### Use Company Config When:
- ✅ You need full company objects with all fields
- ✅ You need company metadata (phone, funding, social media)
- ✅ Building detail views with complete company information
- ✅ You want structured company data separate from contact
- ✅ You need fields not available in denormalized data

**Example Use Case:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified"
      }
    }
  },
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",
      "phone_number",
      "website",
      "latest_funding",
      "facebook_url",
      "twitter_url"
    ]
  }
}
```

### Recommended Approach: Filter with Denormalized, Select with Company Config

**Best Practice**: Use denormalized fields for fast filtering, then use `company_config.select_columns` to get company data in responses:

```json
{
  "where": {
    "range_query": {
      "must": {
        "company_employees_count": {"gte": 100}  // ✅ Filter by denormalized field (fast)
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
      "name",
      "employees_count",
      "phone_number",
      "latest_funding",
      "facebook_url"
    ]
  }
}
```

**Response Structure:**
```json
{
  "data": [
    {
      "id": 43171040,
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "company_id": "c0a8012e-1111-2222-3333-444455556666",
      "company": {
        "uuid": "c0a8012e-1111-2222-3333-444455556666",
        "name": "Altitude Software",
        "employees_count": 500,
        "phone_number": "+1-555-123-4567",
        "latest_funding": "Series B",
        "facebook_url": "https://facebook.com/altitude"
      }
    }
  ]
}
```

**Key Points:**
- Denormalized fields (`company_*`) are **ONLY for filtering** in `where` clauses
- Company Config fields appear in nested `company` object
- Field names differ: `company_name` (filtering) vs `name` (selection)
- **Always use `company_config.select_columns`** to get company data in responses

**Benefits:**
- ⚡ Fast filtering using denormalized fields
- 📊 Complete company data via populated objects
- 🎯 Get all company fields (27 total) including metadata

---

## Usage Patterns

### Pattern 1: Simple List View (8 fields)

**Use Case**: Basic contact list with essential information

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified"
      }
    }
  },
  "select_columns": [
    "uuid",
    "first_name",
    "last_name",
    "email",
    "title",
    "mobile_phone",
    "email_status",
    "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "name"
    ]
  },
  "page": 1,
  "limit": 25
}
```

### Pattern 2: Full Detail View (30+ fields)

**Use Case**: Comprehensive contact information for detail pages

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [43171040]
      }
    }
  },
  "select_columns": [
    "id",
    "uuid",
    "first_name",
    "last_name",
    "company_id",
    "email",
    "title",
    "departments",
    "mobile_phone",
    "email_status",
    "seniority",
    "city",
    "state",
    "country",
    "linkedin_url",
    "created_at",
    "work_direct_phone",
    "home_phone",
    "other_phone",
    "facebook_url",
    "twitter_url",
    "website",
    "stage",
    "updated_at"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",
      "employees_count",
      "industries",
      "annual_revenue",
      "city",
      "state",
      "country"
    ]
  }
}
```

### Pattern 3: Export Mode (All Fields)

**Use Case**: Data export - get all available fields

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
  "limit": 100
}
```

**Note**: Omit `select_columns` to get all fields. For exports, use `limit: 100` (maximum) and paginate.

### Pattern 4: Performance-Optimized View (15 fields)

**Use Case**: Balanced view with essential fields for better performance

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": ["Senior", "Lead"],
        "email_status": "verified"
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "departments",
    "seniority",
    "mobile_phone",
    "email_status",
    "company_name",
    "company_employees_count",
    "company_industries",
    "company_annual_revenue",
    "city",
    "country"
  ],
  "page": 1,
  "limit": 50
}
```

### Pattern 5: Mobile-Optimized View (10 fields)

**Use Case**: Minimal data for mobile applications

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified"
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "mobile_phone",
    "email_status",
    "city",
    "country",
    "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "name"
    ]
  },
  "page": 1,
  "limit": 25
}
```

---

## Complete Examples

### Example 1: Comprehensive Filter with All Field Categories

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
          "text_value": "engineer",
          "filter_key": "title",
          "search_type": "shuffle",
          "fuzzy": true
        },
        {
          "text_value": "micros",
          "filter_key": "company_name",
          "search_type": "substring",
          "operator": "and"
        }
      ]
    },
    "keyword_match": {
      "must": {
        "departments": ["Engineering"],
        "seniority": ["Senior", "Lead"],
        "email_status": "verified",
        "company_industries": ["Software", "SaaS"],
        "company_technologies": ["Python", "Go"]
      }
    },
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2023-01-01T00:00:00Z"
        },
        "company_employees_count": {
          "gte": 100,
          "lte": 1000
        },
        "company_annual_revenue": {
          "gte": 5000000
        }
      }
    }
  },
  "select_columns": [
    "id",
    "uuid",
    "first_name",
    "last_name",
    "company_id",
    "email",
    "title",
    "departments",
    "mobile_phone",
    "email_status",
    "seniority",
    "city",
    "state",
    "country",
    "linkedin_url",
    "created_at",
    "work_direct_phone",
    "home_phone",
    "other_phone",
    "facebook_url",
    "twitter_url",
    "website",
    "stage",
    "updated_at",
    "deleted_at"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",
      "employees_count",
      "industries",
      "keywords",
      "address",
      "annual_revenue",
      "total_funding",
      "technologies",
      "city",
      "state",
      "country",
      "linkedin_url",
      "website",
      "normalized_domain"
    ]
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

### Example 2: Account-Based Marketing Query

```json
{
  "where": {
    "range_query": {
      "must": {
        "company_employees_count": {
          "gte": 500,
          "lte": 5000
        },
        "company_annual_revenue": {
          "gte": 10000000
        }
      }
    },
    "keyword_match": {
      "must": {
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified",
        "company_industries": ["Software", "SaaS", "Technology"]
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "seniority",
    "departments",
    "mobile_phone",
    "work_direct_phone",
    "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",  // Direct name - NO company_ prefix!
      "employees_count",  // Direct name - NO company_ prefix!
      "annual_revenue",  // Direct name - NO company_ prefix!
      "industries",  // Direct name - NO company_ prefix!
      "city",  // Direct name - NO company_ prefix!
      "state",  // Direct name - NO company_ prefix!
      "country",  // Direct name - NO company_ prefix!
      "website",  // Direct name - NO company_ prefix!
      "phone_number",  // Metadata field - not in denormalized
      "latest_funding"  // Metadata field - not in denormalized
    ]
  },
  "order_by": [
    {
      "order_by": "seniority",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 50
}
```

### Example 3: Lead Qualification Query

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
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "title",
    "seniority",
    "departments",
    "mobile_phone",
    "email_status",
    "city",
    "state",
    "country",
    "linkedin_url",
    "work_direct_phone",
    "stage",
    "company_id",
    "created_at"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",  // Direct name - NO company_ prefix!
      "employees_count",  // Direct name - NO company_ prefix!
      "industries",  // Direct name - NO company_ prefix!
      "phone_number",  // Metadata - not in denormalized
      "latest_funding",  // Metadata - not in denormalized
      "facebook_url"  // Metadata - not in denormalized
    ]
  },
  "order_by": [
    {
      "order_by": "seniority",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 50
}
```

---

## Best Practices

### 1. Performance Optimization

**✅ DO:**
- Use `select_columns` for list views (8-15 fields)
- Select only fields you need
- Use simple mode for basic lists
- Omit `select_columns` only when you need all fields

**❌ DON'T:**
- Select all 38 fields for simple list views
- Include fields you won't use
- Forget to use `select_columns` for mobile apps

### 2. Field Selection Strategy

**Simple List View (8 fields):**
```json
{
  "select_columns": [
    "uuid", "first_name", "last_name", "email", 
    "title", "mobile_phone", "email_status", "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": ["name"]
  }
}
```

**Detail View (20-30 fields):**
```json
"select_columns": [
  // Core fields + metadata + company fields you need
]
```

**Export Mode:**
```json
// Omit select_columns entirely
```

### 3. Combining Filterable and Response-Only Fields

**✅ Correct:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified"  // Filterable field
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "email",
    "work_direct_phone",  // Response-only field - OK to select
    "facebook_url"        // Response-only field - OK to select
  ]
}
```

**❌ Incorrect:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "work_direct_phone": "1234567890"  // ERROR: Cannot filter by response-only field
      }
    }
  }
}
```

### 4. Common Patterns

**Pattern: Get Contact with All Phone Numbers**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email": "john.smith@example.com"
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "mobile_phone",
    "work_direct_phone",
    "home_phone",
    "other_phone"
  ]
}
```

**Pattern: Get Contact with Company Info**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [43171040]
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
      "city",
      "state",
      "country",
      "phone_number",
      "latest_funding"
    ]
  }
}
```

**Pattern: Get Contact with Social Media**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": ["Executive", "Principal"]
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "email",
    "linkedin_url",
    "facebook_url",
    "twitter_url",
    "website"
  ]
}
```

### 5. Error Prevention

**Always include `id` in select_columns:**
```json
"select_columns": ["id", "first_name", "last_name", ...]
```

**Don't include denormalized company fields in select_columns:**
```json
// ❌ These are NOT available in select_columns (denormalized fields are for filtering only):
"company_name"
"company_employees_count"
"company_industries"
// etc.

// ✅ Use company_config.select_columns instead:
{
  "company_config": {
    "populate": true,
    "select_columns": ["name", "employees_count", "industries"]
  }
}
```

---

## Common Errors

### Error 1: Filtering by Response-Only Fields

**❌ Incorrect:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "work_direct_phone": "1234567890"  // ERROR
      }
    }
  }
}
```

**✅ Correct:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email": "john.smith@example.com"  // Filter by filterable field
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "last_name",
    "work_direct_phone"  // Then select response-only field
  ]
}
```

### Error 2: Including Non-Existent Fields

**❌ Incorrect:**
```json
{
  "select_columns": [
    "id",
    "first_name",
    "company_latest_funding",      // ERROR: Causes 500 error
    "company_latest_funding_amount", // ERROR: Causes 500 error
    "company_last_raised_at"        // ERROR: Causes 500 error
  ]
}
```

**✅ Correct (Using Company Config):**
```json
{
  "select_columns": [
    "id",
    "first_name",
    "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",
      "total_funding",
      "latest_funding",
      "latest_funding_amount",
      "last_raised_at"
    ]
  }
}
```

### Error 5: Using company_ Prefix in company_config.select_columns

**❌ Incorrect:**
```json
{
  "company_config": {
    "populate": true,
    "select_columns": [
      "company_name",  // ERROR: Wrong! Should be "name"
      "company_employees_count",  // ERROR: Wrong! Should be "employees_count"
      "company_annual_revenue"  // ERROR: Wrong! Should be "annual_revenue"
    ]
  }
}
```

**✅ Correct:**
```json
{
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",  // Correct - direct field name
      "employees_count",  // Correct - direct field name
      "annual_revenue"  // Correct - direct field name
    ]
  }
}
```

**Key Rule:**
- In `select_columns` (contact level): Use `company_*` prefix for denormalized fields
- In `company_config.select_columns` (company level): Use direct field names WITHOUT prefix

### Error 3: Typo in Field Names

**❌ Incorrect:**
```json
{
  "select_columns": [
    "firstname",    // ERROR: Should be "first_name"
    "lastname",     // ERROR: Should be "last_name"
    "companyName"   // ERROR: Should be "company_name"
  ]
}
```

**✅ Correct:**
```json
{
  "select_columns": [
    "first_name",
    "last_name",
    "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "name"
    ]
  }
}
```

### Error 4: Using select_columns to Filter

**❌ Incorrect Understanding:**
- `select_columns` does NOT filter which documents are returned
- It only filters which fields are returned for matched documents

**✅ Correct Understanding:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified"  // This filters documents
      }
    }
  },
  "select_columns": ["id", "email"]  // This only selects fields
}
```

### Error 6: Using Denormalized Fields in select_columns

**❌ Incorrect - Denormalized fields are NOT in select_columns:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software"]  // ✅ OK - filtering
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "company_name",  // ❌ ERROR: Denormalized fields NOT in select_columns!
    "company_employees_count"  // ❌ ERROR: Denormalized fields NOT in select_columns!
  ]
}
```

**✅ Correct - Use company_config.select_columns:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software"]  // ✅ OK - filtering with denormalized
      }
    }
  },
  "select_columns": [
    "id",
    "first_name",
    "company_id"
  ],
  "company_config": {
    "populate": true,
    "select_columns": [
      "uuid",
      "name",  // ✅ Correct - direct name, NO prefix!
      "employees_count",  // ✅ Correct - direct name, NO prefix!
      "phone_number"  // ✅ Correct - metadata field
    ]
  }
}
```

**Remember:**
- **Denormalized fields** (`company_*`): **ONLY for filtering** in `where` clauses
- **Company Config fields** (direct names): **ONLY for selection** in `company_config.select_columns`
- Denormalized fields are **NOT available** in `select_columns` - they will NOT be returned!

---

## Summary

### Quick Reference

| Category | Fields | Filterable | Selectable | Use Case |
|----------|--------|------------|------------|----------|
| **Contact Core** | 16 | ✅ Yes | ✅ Yes | Primary filtering and selection |
| **Contact Metadata** | 9 | ❌ No | ✅ Yes | Additional contact information |
| **Company Denormalized** | 13 | ✅ Yes | ❌ No* | **ONLY for filtering** in `where` clauses |
| **Company Config (Populate)** | 27 | ❌ No** | ✅ Yes | Full company objects in response |
| **Total Contact Fields** | **25** | **16 filterable** | **25 selectable** | Complete contact data |
| **Total Company Fields** | **27** | **17 filterable** | **27 selectable** | Complete company data |

*Denormalized fields are **ONLY for filtering** in `where` clauses, NOT in `select_columns`
**Company Config fields cannot be used in contact filters, but can be used in company filters when querying companies directly

*Company Config fields cannot be used in contact filters, but can be used in company filters when querying companies directly.

### Key Takeaways

1. **25 total contact fields** available for `select_columns` (16 core + 9 metadata)
2. **27 total company fields** available for `company_config.select_columns`
3. **16 contact fields** can be used in filters (core fields)
4. **13 denormalized company fields** can be used in filters (in `where` clauses only)
5. **9 contact fields** are response-only (cannot filter, can select)
6. **17 company fields** are filterable in Company API
7. **10 company fields** are response-only (cannot filter, can select)
8. **Denormalized fields** (`company_*`): **ONLY for filtering** in `where` clauses, NOT in `select_columns`
9. **Company Config** (`company_config.populate`): **ONLY way** to get company data in responses
10. **Always include `id`** in your `select_columns`
11. **Use simple mode** (8 fields) for list views
12. **Omit `select_columns`** for exports to get all fields
13. **Always use `company_config.select_columns`** for company data in responses

### ⚠️ CRITICAL Field Naming Rules

**For Denormalized Fields (ONLY in `where` clauses for filtering):**
- ✅ Use `company_*` prefix: `company_name`, `company_employees_count`, etc.
- ✅ Used ONLY for filtering contacts by company attributes
- ❌ **NOT available** in `select_columns` - they will NOT be returned
- ✅ Fast - already in Elasticsearch index

**For Company Config Fields (in `company_config.select_columns`):**
- ✅ Use direct names WITHOUT prefix: `name`, `employees_count`, etc.
- ❌ DO NOT use `company_*` prefix - it will cause errors!
- ✅ These appear in nested `company` object
- ✅ **ONLY way** to get company data in responses
- ⚠️ Slower - requires separate PostgreSQL query

**Field Name Mapping:**
- `company_name` (filtering) ↔ `name` (selection)
- `company_employees_count` (filtering) ↔ `employees_count` (selection)
- `company_annual_revenue` (filtering) ↔ `annual_revenue` (selection)
- And so on...

**Remember:** 
- Denormalized fields (`company_*`) = **Filtering only** (in `where` clauses)
- Company Config fields (direct names) = **Selection only** (in `company_config.select_columns`)
- They serve different purposes and cannot be interchanged!

## Authentication

All API requests require authentication using the `X-API-Key` header:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {...},
    "select_columns": ["id", "first_name", "last_name"]
  }'
```

**Security Best Practices**:
- Store API keys securely (environment variables, secret management systems)
- Always use HTTPS in production
- Rotate API keys periodically
- Use different API keys for different environments

**See**: [API Reference - Authentication](./06-api-reference.md#authentication) for complete authentication details.

---

## Using select_columns with Write Operations

When creating or updating records, you can use `select_columns` in the response to limit which fields are returned, but note that write operations don't support `select_columns` in the request body - they always return the full record.

### Create Operations

When creating a company or contact, the response includes all fields by default. You cannot use `select_columns` in create requests, but you can filter the response in your application code.

**Example: Create company (full response)**

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "name": "TechStart Inc",
    "normalized_domain": "techstart.com",
    "employees_count": 50
  }'
```

**Response** (includes all fields):
```json
{
  "data": {
    "uuid": "...",
    "name": "TechStart Inc",
    "normalized_domain": "techstart.com",
    "employees_count": 50,
    "industries": null,
    "annual_revenue": null,
    // ... all other fields
  },
  "success": true
}
```

### Update Operations

Update operations also return the full updated record. You cannot use `select_columns` in update requests.

**Example: Update company (full response)**

```bash
curl -X PUT https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/company-uuid \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "employees_count": 75,
    "annual_revenue": 7500000
  }'
```

### Filtering After Write Operations

If you need to retrieve only specific fields after a write operation, use a filter query with `select_columns`:

**Example: Create company, then retrieve with select_columns**

```bash
# Step 1: Create company
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "name": "TechStart Inc",
    "normalized_domain": "techstart.com"
  }'

# Step 2: Retrieve with select_columns
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "normalized_domain": ["techstart.com"]
        }
      }
    },
    "select_columns": ["uuid", "name", "employees_count", "annual_revenue"]
  }'
```

### Best Practices for Write Operations

1. **Write operations return full records** - Always expect complete data in responses
2. **Use filter queries with select_columns** - If you need specific fields, query after writing
3. **Minimize data in requests** - Only send fields you want to create/update
4. **Use partial updates** - Only include fields you want to change in update requests

**See**: 
- [Company API - Write Operations](../company.md#write-operations)
- [Contact API - Write Operations](../contacts.md#write-operations)
- [Field Reference - Fields in Write Operations](./04-filter-field-reference.md#fields-in-write-operations)

---

### Related Documentation

- [Contact Filters Complete Guide](./02-contact-filters-complete-guide.md)
- [Company Filters Complete Guide](./01-company-filters-complete-guide.md)
- [API Reference](./06-api-reference.md)
- [Examples and Use Cases](./05-examples-use-cases.md)

---

**Last Updated**: 2025-01-XX
**Version**: 1.2
