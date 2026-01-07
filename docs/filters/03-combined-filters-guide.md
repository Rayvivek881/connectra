# Combined Filters Guide - Company and Contact

## Table of Contents

1. [Overview](#overview)
2. [Using Company and Contact Filters Together](#using-company-and-contact-filters-together)
3. [Account-Based Filtering Patterns](#account-based-filtering-patterns)
4. [Cross-Entity Filtering Strategies](#cross-entity-filtering-strategies)
5. [Best Practices for Combined Filters](#best-practices-for-combined-filters)
6. [Authentication and Security](#authentication-and-security)
7. [VQL Syntax Reference](#vql-syntax-reference)
8. [Pagination Strategies](#pagination-strategies)
9. [Field Selection Optimization](#field-selection-optimization)
10. [Error Handling](#error-handling)

## Overview

This guide explains how to effectively combine company and contact filters to create powerful, account-based search queries. The contact index includes denormalized company fields (with `company_` prefix), enabling you to filter contacts directly by company attributes in a single query, or use separate company and contact queries for more complex scenarios.

> **Authentication Required**: All examples in this guide show only the JSON request body. When making actual API calls, you must include the `X-API-Key` header for authentication. See [06-api-reference.md](./06-api-reference.md#authentication) for complete HTTP request examples with headers.

## Using Company and Contact Filters Together

### Strategy 1: Company-First Approach

**Step 1: Find target companies**
```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "SaaS"],
        "country": ["USA"]
      }
    },
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
  "limit": 50
}
```

**Step 2: Find contacts at those companies**
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "company-uuid-1",
          "company-uuid-2",
          "company-uuid-3"
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

### Strategy 2: Contact-First Approach

**Step 1: Find target contacts**
```json
POST /contacts
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
        "departments": ["Sales", "Marketing"],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified",
        "country": ["USA"]
      }
    }
  },
  "page": 1,
  "limit": 100
}
```

**Step 2: Get their companies**
```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [1, 2, 3, 4, 5]
      }
    }
  },
  "page": 1,
  "limit": 50
}
```

### Strategy 3: Direct Company Field Filtering (Single Query)

**Filter contacts directly by denormalized company fields - no separate company query needed**

The contact index includes denormalized company data, allowing you to filter contacts by company attributes in a single query. This is the most efficient approach when you only need contacts and want to filter by company criteria.

**Example: Find senior contacts at high-revenue software companies**
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS", "Technology"],
        "company_country": ["USA", "Canada", "UK"],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified"
      }
    },
    "range_query": {
      "must": {
        "company_employees_count": {
          "gte": 50,
          "lte": 1000
        },
        "company_annual_revenue": {
          "gte": 1000000
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

**Example: Find engineering contacts at companies using specific technologies**
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_technologies": ["Python", "Go", "React", "AWS"],
        "departments": ["Engineering"],
        "seniority": ["Senior", "Lead", "Principal"]
      }
    },
    "range_query": {
      "must": {
        "company_employees_count": {
          "gte": 100
        }
      }
    }
  },
  "page": 1,
  "limit": 50
}
```

**Benefits of Direct Company Field Filtering**:

- Single API call instead of two (company query + contact query)
- Faster execution (no need to fetch company IDs first)
- Simpler code (no need to extract and pass company IDs)
- More efficient for account-based filtering scenarios

**When to Use**: Use this strategy when you only need contacts and want to filter by company attributes. Use Strategy 1 (Company-First) when you also need the company data itself.

### Strategy 4: Parallel Filtering

**Filter companies and contacts in parallel with matching criteria**

**Companies Query**:
```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology"],
        "country": ["USA", "Canada"]
      }
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50,
          "lte": 1000
        }
      }
    }
  },
  "page": 1,
  "limit": 50
}
```

**Contacts Query** (matching criteria):
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "country": ["USA", "Canada"],
        "seniority": ["Senior", "Lead"],
        "email_status": "verified"
      }
    }
  },
  "page": 1,
  "limit": 100
}
```

**Then join on application side** using `company_id` field.

---

## Account-Based Filtering Patterns

### Pattern 1: High-Value Account Targeting

**Goal**: Find decision-makers at high-value companies

**Approach A: Using Denormalized Company Fields (Single Query)**
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS", "Technology"],
        "company_country": ["USA", "Canada", "UK"],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "departments": ["Sales", "Marketing", "Engineering", "Product"],
        "email_status": "verified"
      }
    },
    "range_query": {
      "must": {
        "company_employees_count": {
          "gte": 100,
          "lte": 1000
        },
        "company_annual_revenue": {
          "gte": 5000000
        },
        "company_total_funding": {
          "gte": 10000000
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
  "limit": 200
}
```

**Approach B: Two-Step Process (Company-First)**

**Step 1: Identify high-value companies**
```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "SaaS", "Technology"],
        "country": ["USA", "Canada", "UK"]
      }
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 100,
          "lte": 1000
        },
        "annual_revenue": {
          "gte": 5000000
        },
        "total_funding": {
          "gte": 10000000
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
  "limit": 100
}
```

**Step 2: Find decision-makers at those companies**
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "company-uuid-1",
          "company-uuid-2",
          "..."
        ],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "departments": ["Sales", "Marketing", "Engineering", "Product"],
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
  "limit": 200
}
```

### Pattern 2: Industry-Specific Targeting

**Goal**: Find contacts in specific industries with matching company criteria

**Companies Query**:
```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Healthcare", "Medical Devices"],
        "technologies": ["AI", "Machine Learning", "Cloud"]
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
  "page": 1,
  "limit": 50
}
```

**Contacts Query**:
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "company-uuid-1",
          "..."
        ],
        "departments": ["Engineering", "Research", "Product"],
        "seniority": ["Senior", "Lead", "Principal"],
        "email_status": "verified"
      }
    },
    "text_matches": {
      "must": [
        {
          "text_value": "healthcare medical",
          "filter_key": "title",
          "search_type": "shuffle",
          "operator": "or",
          "fuzzy": true
        }
      ]
    }
  },
  "page": 1,
  "limit": 100
}
```

### Pattern 3: Geographic Targeting

**Goal**: Find companies and contacts in specific regions

**Companies Query**:
```json
POST /companies
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "Silicon Valley",
          "filter_key": "address",
          "search_type": "shuffle",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "state": ["CA"],
        "country": ["USA"]
      }
    }
  },
  "page": 1,
  "limit": 50
}
```

**Contacts Query**:
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "company-uuid-1",
          "..."
        ],
        "state": ["CA"],
        "country": ["USA"]
      }
    }
  },
  "page": 1,
  "limit": 100
}
```

### Pattern 4: Technology Stack Targeting

**Goal**: Find companies using specific technologies and their technical contacts

**Companies Query**:
```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "technologies": ["Python", "Go", "React", "AWS"],
        "industries": ["Software", "SaaS"]
      }
    }
  },
  "page": 1,
  "limit": 50
}
```

**Contacts Query**:
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "company-uuid-1",
          "..."
        ],
        "departments": ["Engineering", "Product"],
        "seniority": ["Senior", "Lead", "Principal"]
      }
    },
    "text_matches": {
      "must": [
        {
          "text_value": "engineer developer architect",
          "filter_key": "title",
          "search_type": "shuffle",
          "operator": "or",
          "fuzzy": true
        }
      ]
    }
  },
  "page": 1,
  "limit": 100
}
```

### Pattern 5: Growth Stage Targeting

**Goal**: Find well-funded companies and their key contacts

**Companies Query**:
```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology", "SaaS"]
      }
    },
    "range_query": {
      "must": {
        "total_funding": {
          "gte": 10000000
        },
        "annual_revenue": {
          "gte": 5000000
        },
        "employees_count": {
          "gte": 50,
          "lte": 500
        }
      }
    }
  },
  "order_by": [
    {
      "order_by": "total_funding",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 50
}
```

**Contacts Query**:
```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "company-uuid-1",
          "..."
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

---

## Cross-Entity Filtering Strategies

### Strategy 1: Matching Geographic Criteria

**Use Case**: Find companies and contacts in the same geographic region

**Companies**:

- Filter by: `country`, `state`, `city`, `address` (text search)

**Contacts**:

- Filter by: `country`, `state`, `city` (text or keyword)

**Example**:
```json
// Companies in California
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "state": ["CA"],
        "country": ["USA"]
      }
    }
  }
}

// Contacts in California
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "state": ["CA"],
        "country": ["USA"]
      }
    }
  }
}
```

### Strategy 2: Industry + Department Alignment

**Use Case**: Find companies in specific industries and contacts in relevant departments

**Companies**:

- Filter by: `industries`, `keywords`

**Contacts**:

- Filter by: `departments`, `title` (text search)

**Example**:
```json
// Software companies
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "SaaS"]
      }
    }
  }
}

// Engineering contacts
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering"],
        "company_id": ["..."]
      }
    }
  }
}
```

### Strategy 3: Technology Stack + Technical Roles

**Use Case**: Find companies using specific technologies and their technical staff

**Companies**:

- Filter by: `technologies`, `keywords`

**Contacts**:

- Filter by: `title` (text search for technical terms), `departments`

**Example**:
```json
// Companies using Python and Go
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "technologies": ["Python", "Go"]
      }
    }
  }
}

// Technical contacts
POST /contacts
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "engineer developer",
          "filter_key": "title",
          "search_type": "shuffle",
          "operator": "or"
        }
      ]
    },
    "keyword_match": {
      "must": {
        "company_id": ["..."],
        "departments": ["Engineering"]
      }
    }
  }
}
```

### Strategy 4: Company Size + Contact Seniority

**Use Case**: Find mid-size companies and their senior decision-makers

**Companies**:

- Filter by: `employees_count` (range query)

**Contacts**:

- Filter by: `seniority` (keyword)

**Example**:
```json
// Mid-size companies (50-500 employees)
POST /companies
{
  "where": {
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50,
          "lte": 500
        }
      }
    }
  }
}

// Senior contacts
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": ["..."],
        "seniority": ["Senior", "Lead", "Principal", "Executive"]
      }
    }
  }
}
```

---

## Best Practices for Combined Filters

### 1. Choose the Right Strategy

**For Single-Query Efficiency**: Use denormalized company fields (`company_*`) in contact queries when you only need contacts filtered by company attributes. This requires only one API call.

**For Company Data Access**: Use the company-first approach when you also need the company records themselves, or when you need to perform complex company-level analysis before filtering contacts.

**Pattern A (Single Query - Recommended)**:

1. Filter contacts directly using `company_*` fields (e.g., `company_industries`, `company_employees_count`)
2. Combine with contact filters (e.g., `seniority`, `departments`, `email_status`)

**Pattern B (Two-Step Process)**:

1. Filter companies by industry, size, revenue, location
2. Extract company IDs/UUIDs
3. Filter contacts by those company IDs plus role/seniority criteria

### 2. Use Matching Criteria Across Entities

**Why**: Ensures consistency in your targeting.

**Example**:

- Companies: `country: ["USA"]`
- Contacts: `country: ["USA"]` (same countries)

### 3. Leverage Company Data for Contact Filtering

**Why**: Company attributes (industry, size, revenue) help identify relevant contacts.

**Pattern**:

- High-revenue companies → Target senior/executive contacts
- Technology companies → Target engineering/product contacts
- Growth-stage companies → Target sales/marketing contacts

### 4. Combine Multiple Filter Types

**Why**: More precise targeting.

**Example**:
```json
// Companies: Industry + Size + Location
POST /companies
{
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
  }
}

// Contacts: Company + Role + Seniority + Email Status
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": ["..."],
        "departments": ["Engineering"],
        "seniority": ["Senior", "Lead"],
        "email_status": "verified"
      }
    }
  }
}
```

### 5. Use Pagination Efficiently

**Why**: Large result sets require efficient pagination.

**Best Practice**:

- Use `search_after` for large datasets
- Keep `limit` reasonable (25-50)
- Use count endpoints when you only need totals
- Use `select_columns` to limit returned fields and improve performance

### 6. Cache Company Results

**Why**: Company filters are often reused for multiple contact queries.

**Pattern**:

1. Query companies once
2. Cache company IDs/UUIDs
3. Use cached IDs for multiple contact queries

### 7. Filter by Email Status for Contacts

**Why**: Ensures contactability.

**Always include**:
```json
{
  "keyword_match": {
    "must": {
      "email_status": "verified"
    }
  }
}
```

### 8. Use Seniority for Decision-Makers

**Why**: Targets people with decision-making authority.

**Recommended values**:

- `["Senior", "Lead", "Principal", "Executive"]`

### 9. Combine Text and Keyword Filters

**Why**: Text filters provide flexibility, keyword filters provide precision.

**Example**:
```json
{
  "text_matches": {
    "must": [
      {
        "text_value": "director manager",
        "filter_key": "title",
        "search_type": "shuffle",
        "operator": "or"
      }
    ]
  },
  "keyword_match": {
    "must": {
      "seniority": ["Senior", "Lead"],
      "email_status": "verified"
    }
  }
}
```

### 10. Validate Company-Contact Relationships

**Why**: Ensure contacts belong to target companies.

**Option A: Filter by company_id (when you have specific company UUIDs)**
```json
{
  "keyword_match": {
    "must": {
      "company_id": ["company-uuid-1", "company-uuid-2", "..."]
    }
  }
}
```

**Option B: Filter by denormalized company fields (when filtering by company attributes)**
```json
{
  "keyword_match": {
    "must": {
      "company_industries": ["Software", "SaaS"],
      "company_country": ["USA"]
    }
  },
  "range_query": {
    "must": {
      "company_employees_count": {"gte": 100}
    }
  }
}
```

**Note**: Denormalized company fields (`company_*`) are automatically kept in sync with company data, so filtering by these fields ensures contacts match the company criteria without needing a separate company query.

---

## Common Combined Filter Patterns

### Pattern A: Lead Generation

**Companies**: High-value, growth-stage companies
**Contacts**: Verified, senior decision-makers

### Pattern B: Competitive Intelligence

**Companies**: Competitors in same space
**Contacts**: Key personnel at competitor companies

### Pattern C: Partnership Development

**Companies**: Complementary businesses
**Contacts**: Partnership/business development contacts

### Pattern D: Market Research

**Companies**: Companies in target market
**Contacts**: Industry experts and thought leaders

### Pattern E: Investment Targeting

**Companies**: Well-funded, high-growth companies
**Contacts**: C-level executives and founders

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

## CRUD for Account-Based Workflows

When working with account-based marketing, you'll often need to create or update companies and contacts together. The following CRUD operations support these workflows.

### Creating Companies and Contacts

**1. Create Company First**:
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "name": "Acme Software Corp",
    "normalized_domain": "acme.com",
    "industries": ["Software"],
    "country": "USA"
  }'
```

**2. Create Contacts for the Company**:
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "email": "john.doe@acme.com",
    "first_name": "John",
    "last_name": "Doe",
    "company_id": "company-uuid-from-step-1",
    "title": "VP Engineering"
  }'
```

### Bulk Account Setup

Use `batch-upsert` to efficiently set up multiple accounts:

```bash
# Bulk upsert companies
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/batch-upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "pg_companies": [...],
    "es_companies": [...]
  }'

# Bulk upsert contacts
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/batch-upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "pg_contacts": [...],
    "es_contacts": [...]
  }'
```

**See**: 
- [Company API - Write Operations](../company.md#write-operations)
- [Contact API - Write Operations](../contacts.md#write-operations)
- [CRUD Implementation Plan](../CRUD_IMPLEMENTATION_PLAN.md)

---

## Related Documentation

- [Company Filters Guide](./01-company-filters-complete-guide.md)
- [Contact Filters Guide](./02-contact-filters-complete-guide.md)
- [Filter Field Reference](./04-filter-field-reference.md)
- [Examples and Use Cases](./05-examples-use-cases.md)
- [API Reference](./06-api-reference.md)

---

**Last Updated**: 2025-01-XX  
**Version**: 1.2

