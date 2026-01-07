# Examples and Use Cases

## Table of Contents

1. [Overview](#overview)
2. [Company Filter Examples](#company-filter-examples)
3. [Contact Filter Examples](#contact-filter-examples)
4. [Combined Use Cases](#combined-use-cases)
5. [Real-World Scenarios](#real-world-scenarios)

## Overview

This document provides comprehensive examples and real-world use cases for filtering companies and contacts. Each example includes the complete JSON request body, curl commands, and explains the use case.

> **Authentication Required**: All examples in this guide include complete HTTP request examples with headers. You must include the `X-API-Key` header for authentication.

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

### Authentication Header

All requests require the `X-API-Key` header:

```
X-API-Key: your-secret-api-key
```

**Note**: Replace `your-secret-api-key` with your actual API key.

---

## Company Filter Examples

### Example 1: Simple Name Search

**Use Case**: Find companies by name

**Complete HTTP Request**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
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
  }'
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

**Expected Response** (200 OK):

```json
{
  "data": [
    {
      "id": 1,
      "uuid": "c0a8012e-1111-2222-3333-444455556666",
      "name": "Acme Software Corp",
      "employees_count": 120,
      "industries": ["Software", "Technology"],
      "annual_revenue": 5000000,
      "country": "USA",
      "created_at": "2024-01-15T08:00:00Z"
    }
  ],
  "success": true
}
```

### Example 1a: Substring Name Search

**Use Case**: Find companies by partial name match (autocomplete-style)

**Complete HTTP Request**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
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
    },
    "select_columns": ["id", "name", "employees_count"],
    "page": 1,
    "limit": 25
  }'
```

**Request Body**:

```json
{
  "where": {
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
  },
  "select_columns": ["id", "name", "employees_count"],
  "page": 1,
  "limit": 25
}
```

**Note**: This will match companies with "soft" anywhere in the name (e.g., "Microsoft", "Software Corp", "SoftTech"). Uses ngram matching for efficient partial text search. The `name` field supports substring search with minimum 3 characters (min_gram: 3). The `select_columns` parameter limits which fields are returned from PostgreSQL after the Elasticsearch search.

**Expected Response** (200 OK):

```json
{
  "data": [
    {
      "id": 1,
      "name": "Acme Software Corp",
      "employees_count": 120
    },
    {
      "id": 2,
      "name": "Microsoft Corporation",
      "employees_count": 221000
    },
    {
      "id": 3,
      "name": "SoftTech Solutions",
      "employees_count": 50
    }
  ],
  "success": true
}
```

**Key Points**:
- Substring search requires minimum 3 characters
- Only `name` field supports substring search for companies
- More efficient than text search for autocomplete scenarios
- `select_columns` reduces response payload size

### Example 2: Industry and Technology Filtering

**Use Case**: Find software companies using specific technologies

**Request**:

```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "SaaS"],
        "technologies": ["Python", "Go", "React"]
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

### Example 3: Size and Revenue Filtering

**Use Case**: Find mid-size companies with good revenue

**Request**:

```json
POST /companies
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

### Example 4: Geographic Filtering

**Use Case**: Find companies in specific regions

**Request**:

```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "country": ["USA", "Canada", "UK"],
        "state": ["CA", "NY", "TX"]
      }
    }
  },
  "page": 1,
  "limit": 50
}
```

### Example 5: Complex Multi-Criteria Search

**Use Case**: Find high-value AI companies in specific regions

**Request**:

```json
POST /companies
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
        "industries": ["Software", "Technology"],
        "technologies": ["Python", "Machine Learning"],
        "country": ["USA"]
      },
      "must_not": {
        "keywords": ["Legacy", "Outdated"]
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
        "created_at": {
          "gte": "2020-01-01T00:00:00Z"
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
      "order_by": "employees_count",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 50
}
```

### Example 6: Funding-Based Filtering

**Use Case**: Find well-funded companies

**Request**:

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
          "gte": 100
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

---

## Contact Filter Examples

### Example 1: Simple Name Search

**Use Case**: Find contacts by name

**Request**:

```json
POST /contacts
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

**Response**:

```json
{
  "data": [
    {
      "id": 43171040,
      "uuid": "021c8c87-1a5b-55a7-86c8-8f6f4710924e",
      "first_name": "John",
      "last_name": "Smith",
      "email": "john.smith@example.com",
      "title": "Senior Software Engineer",
      "departments": ["Engineering"],
      "seniority": "Senior",
      "email_status": "verified",
      "company_id": "c0a8012e-1111-2222-3333-444455556666",
      "created_at": "2024-01-15T08:00:00Z"
    }
  ],
  "success": true
}
```

### Example 2: Job Title Search

**Use Case**: Find engineers and developers

**Request**:

```json
POST /contacts
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
  "limit": 50
}
```

### Example 2a: Substring Title Search

**Use Case**: Find contacts with partial title match (e.g., autocomplete)

**Request**:

```json
POST /contacts
{
  "where": {
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
  },
  "select_columns": ["id", "first_name", "last_name", "title", "email"],
  "page": 1,
  "limit": 50
}
```

**Note**: This will match titles containing "engin" (e.g., "Engineer", "Engineering Manager", "Senior Engineer"). Uses ngram matching for efficient partial text search. The `title` field supports substring search with minimum 5 characters (min_gram: 5). The `select_columns` parameter limits which fields are returned from PostgreSQL after the Elasticsearch search.

### Example 3: Department and Seniority Filtering

**Use Case**: Find senior contacts in specific departments

**Request**:

```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "departments": ["Engineering", "Product"],
        "seniority": ["Senior", "Lead", "Principal"],
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
  "limit": 50
}
```

### Example 4: Verified Contacts in Specific Countries

**Use Case**: Find verified contacts for email campaigns

**Request**:

```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "country": ["USA", "UK", "Canada"],
        "email_status": "verified"
      }
    }
  },
  "page": 1,
  "limit": 100
}
```

### Example 5: Company-Based Contact Filtering (by company_id)

**Use Case**: Find all contacts at specific companies

**Request**:

```json
POST /contacts
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
      "order_by": "seniority",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 100
}
```

### Example 5a: Company-Based Contact Filtering (using denormalized fields)

**Use Case**: Find contacts at companies matching specific criteria - single query approach

**Request**:

```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS"],
        "company_country": ["USA", "Canada"],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
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

**Note**: This example uses denormalized company fields (`company_*`) to filter contacts directly by company attributes without needing a separate company query. This is more efficient than the two-step approach (query companies first, then filter contacts by company_id).

### Example 5b: Filtering Contacts by Company Name (substring)

**Use Case**: Find contacts at companies with partial name match

**Request**:

```json
POST /contacts
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
    },
    "keyword_match": {
      "must": {
        "email_status": "verified"
      }
    }
  },
  "select_columns": ["id", "first_name", "last_name", "title", "email", "company_id"],
  "page": 1,
  "limit": 50
}
```

**Note**: The `company_name` field supports substring search with minimum 5 characters (min_gram: 5). This allows finding contacts at companies like "Microsoft", "Microsystems", etc. The `select_columns` parameter limits the response to only specified fields.

### Example 5c: Filtering Contacts and Populating Company Data (company_config)

**Use Case**: Find contacts at high-revenue software companies and include company details in the response

**Request**:

```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_industries": ["Software", "SaaS"],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified"
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
  },
  "select_columns": ["id", "first_name", "last_name", "title", "email", "company_id"],
  "company_config": {
    "populate": true,
    "select_columns": ["uuid", "name", "employees_count", "industries", "annual_revenue", "website"]
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

**Complete HTTP Request**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "company_industries": ["Software", "SaaS"],
          "seniority": ["Senior", "Lead", "Principal", "Executive"],
          "email_status": "verified"
        }
      },
      "range_query": {
        "must": {
          "company_annual_revenue": {"gte": 5000000},
          "company_employees_count": {"gte": 100}
        }
      }
    },
    "select_columns": ["id", "first_name", "last_name", "title", "email", "company_id"],
    "company_config": {
      "populate": true,
      "select_columns": ["uuid", "name", "employees_count", "industries", "annual_revenue", "website"]
    },
    "page": 1,
    "limit": 25
  }'
```

**Expected Response** (200 OK):

```json
{
  "data": [
    {
      "id": 123,
      "first_name": "John",
      "last_name": "Doe",
      "title": "Senior Software Engineer",
      "email": "john.doe@acme.com",
      "company_id": "company-uuid-here",
      "company": {
        "uuid": "company-uuid-here",
        "name": "Acme Software Corp",
        "employees_count": 150,
        "industries": ["Software", "SaaS"],
        "annual_revenue": 7500000,
        "website": "https://acme.com"
      }
    }
  ],
  "success": true
}
```

**Key Points**:
- ✅ Uses denormalized `company_*` fields in `where` clause for filtering
- ✅ Uses `company_config.populate: true` to fetch company data
- ✅ Uses direct field names (no `company_*` prefix) in `company_config.select_columns`
- ✅ Company data is returned in a nested `company` object

### Example 6: Complex Multi-Criteria Contact Search

**Use Case**: Find senior software engineers in tech companies

**Request**:

```json
POST /contacts
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
      }
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

## Combined Use Cases

### Use Case 1: Lead Generation Campaign

**Goal**: Find high-value companies and their decision-makers

**Step 1: Find target companies**

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
      "order_by": "annual_revenue",
      "order_direction": "desc"
    }
  ],
  "page": 1,
  "limit": 100
}
```

**Step 2: Find decision-makers at those companies**

**Option A: Using company_id (requires company UUIDs from Step 1)**

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

**Option B: Using denormalized company fields (single query, no Step 1 needed)**

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
  "limit": 200
}
```

**Note**: Option B is more efficient as it requires only one API call and filters contacts directly by company attributes using denormalized fields.

### Use Case 2: Competitive Intelligence

**Goal**: Find competitors and their key personnel

**Step 1: Find competitors**

```json
POST /companies
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "cloud computing saas",
          "filter_key": "name",
          "search_type": "shuffle",
          "operator": "or",
          "fuzzy": true
        }
      ]
    },
    "keyword_match": {
      "must": {
        "industries": ["SaaS", "Cloud Computing"],
        "technologies": ["AWS", "Azure", "GCP"]
      },
      "must_not": {
        "keywords": ["Our Company", "Partner"]
      }
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 100,
          "lte": 5000
        },
        "annual_revenue": {
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
  "limit": 50
}
```

**Step 2: Find key personnel**

```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "competitor-uuid-1",
          "..."
        ],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "departments": ["Engineering", "Product", "Sales"]
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

### Use Case 3: Partnership Development

**Goal**: Find potential partners and their business development contacts

**Step 1: Find complementary companies**

```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology"],
        "technologies": ["Python", "Go", "React"],
        "country": ["USA", "Canada"]
      }
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50,
          "lte": 500
        },
        "annual_revenue": {
          "gte": 5000000,
          "lte": 50000000
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

**Step 2: Find partnership contacts**

```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "partner-uuid-1",
          "..."
        ],
        "departments": ["Sales", "Business Development", "Partnerships"],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified"
      }
    },
    "text_matches": {
      "must": [
        {
          "text_value": "partnership business development",
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

### Use Case 4: Market Research

**Goal**: Find companies and contacts in target market segments

**Step 1: Find companies in target market**

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
  "limit": 50
}
```

**Step 2: Find industry experts**

```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "market-uuid-1",
          "..."
        ],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "email_status": "verified"
      }
    },
    "text_matches": {
      "must": [
        {
          "text_value": "director manager lead",
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

### Use Case 5: Investment Targeting

**Goal**: Find well-funded companies and their executives

**Step 1: Find investment targets**

```json
POST /companies
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology", "SaaS"],
        "country": ["USA"]
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
          "gte": 100
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

**Step 2: Find executives**

```json
POST /contacts
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": [
          "investment-uuid-1",
          "..."
        ],
        "seniority": ["Executive", "Principal"],
        "email_status": "verified"
      }
    },
    "text_matches": {
      "must": [
        {
          "text_value": "CEO CTO CFO founder",
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

---

## Real-World Scenarios

### Scenario 1: B2B Sales Outreach

**Objective**: Find decision-makers at high-value accounts

**Company Filter**:

```json
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
          "gte": 100,
          "lte": 1000
        },
        "annual_revenue": {
          "gte": 5000000
        }
      }
    }
  }
}
```

**Contact Filter**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "company_id": ["..."],
        "seniority": ["Senior", "Lead", "Principal", "Executive"],
        "departments": ["Sales", "Marketing", "Product"],
        "email_status": "verified"
      }
    }
  }
}
```

### Scenario 2: Recruiting Campaign

**Objective**: Find senior engineers for recruitment

**Contact Filter**:

```json
{
  "where": {
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
    },
    "keyword_match": {
      "must": {
        "departments": ["Engineering"],
        "seniority": ["Senior", "Lead", "Principal"],
        "email_status": "verified",
        "country": ["USA", "Canada"]
      }
    },
    "range_query": {
      "must": {
        "created_at": {
          "gte": "2023-01-01T00:00:00Z"
        }
      }
    }
  }
}
```

### Scenario 3: Email Marketing Campaign

**Objective**: Find verified contacts for email campaigns

**Contact Filter**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "email_status": "verified",
        "country": ["USA", "UK", "Canada"]
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
  "limit": 100
}
```

### Scenario 4: Account-Based Marketing

**Objective**: Target specific accounts with multiple contacts

**Company Filter**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [1, 2, 3, 4, 5]
      }
    }
  }
}
```

**Contact Filter**:

```json
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
        "email_status": "verified"
      }
    }
  }
}
```

### Scenario 5: Technology Stack Analysis

**Objective**: Find companies using specific technologies

**Company Filter**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "technologies": ["Python", "Go", "React", "AWS"]
      }
    }
  },
  "order_by": [
    {
      "order_by": "employees_count",
      "order_direction": "desc"
    }
  ]
}
```

---

## Error Handling Examples

This section demonstrates common errors and how to fix them.

### Error Example 1: Invalid Field Name

**❌ Incorrect Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "invalid_field": ["Software"]  // ERROR: Field doesn't exist
        }
      }
    }
  }'
```

**Error Response** (400 Bad Request):

```json
{
  "error": "ERR_ELASTICSEARCH_FAILURE: search engine returned status 400; details: ...",
  "success": false
}
```

**✅ Corrected Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "industries": ["Software"]  // ✅ Correct field name
        }
      }
    }
  }'
```

**Solution**: Use `/companies/filters` endpoint to get valid field names before building queries.

---

### Error Example 2: Page Number Exceeded

**❌ Incorrect Query**:

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
    "page": 15,  // ERROR: Maximum is 10
    "limit": 25
  }'
```

**Error Response** (400 Bad Request):

```json
{
  "error": "ERR_PAGE_OUT_OF_RANGE: the requested page number is beyond the available range; verify total pages before requesting",
  "success": false
}
```

**✅ Corrected Query - Option 1 (Page-Based)**:

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
    "page": 10,  // ✅ Maximum allowed
    "limit": 25
  }'
```

**✅ Corrected Query - Option 2 (Cursor-Based)**:

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
    "search_after": [5000000, 123],  // ✅ Use cursor-based pagination
    "limit": 25
  }'
```

**Solution**: Use `search_after` for pagination beyond page 10.

---

### Error Example 3: Page Size Exceeded

**❌ Incorrect Query**:

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
    "limit": 150  // ERROR: Maximum is 100
  }'
```

**Error Response** (400 Bad Request):

```json
{
  "error": "ERR_PAGE_SIZE_EXCEEDED: the requested page size surpasses the maximum allowed limit; consider using pagination with smaller batches",
  "success": false
}
```

**✅ Corrected Query**:

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
    "limit": 100  // ✅ Maximum allowed
  }'
```

**Solution**: Reduce `limit` to maximum 100.

---

### Error Example 4: Missing Required Field

**❌ Incorrect Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "text_matches": {
        "must": [
          {
            "text_value": "software",
            // ERROR: Missing required 'filter_key'
            "search_type": "shuffle"
          }
        ]
      }
    }
  }'
```

**Error Response** (400 Bad Request):

```json
{
  "error": "ERR_INVALID_REQUEST_BODY: the request body is invalid; check JSON syntax and required fields",
  "success": false
}
```

**✅ Corrected Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "text_matches": {
        "must": [
          {
            "text_value": "software",
            "filter_key": "name",  // ✅ Required field
            "search_type": "shuffle",
            "fuzzy": true
          }
        ]
      }
    },
    "page": 1,
    "limit": 25
  }'
```

**Solution**: Ensure all required fields are present in text match objects.

---

### Error Example 5: Invalid Search Type

**❌ Incorrect Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "text_matches": {
        "must": [
          {
            "text_value": "software",
            "filter_key": "name",
            "search_type": "invalid_type"  // ERROR: Invalid search_type
          }
        ]
      }
    }
  }'
```

**Error Response** (500 Internal Server Error):

```json
{
  "error": "ERR_ELASTICSEARCH_FAILURE: search engine returned status 400; details: ...",
  "success": false
}
```

**✅ Corrected Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "text_matches": {
        "must": [
          {
            "text_value": "software",
            "filter_key": "name",
            "search_type": "shuffle",  // ✅ Valid: "exact", "shuffle", or "substring"
            "fuzzy": true
          }
        ]
      }
    },
    "page": 1,
    "limit": 25
  }'
```

**Solution**: Use valid `search_type` values: `"exact"`, `"shuffle"`, or `"substring"`.

---

### Error Example 6: Using Response-Only Field in Filter

**❌ Incorrect Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "facebook_url": "https://facebook.com/company"  // ERROR: Response-only field
        }
      }
    }
  }'
```

**Error Response** (500 Internal Server Error):

```json
{
  "error": "ERR_ELASTICSEARCH_FAILURE: search engine returned status 400; details: ...",
  "success": false
}
```

**✅ Corrected Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "industries": ["Software"]  // ✅ Filter by filterable field
        }
      }
    },
    "select_columns": ["id", "name", "facebook_url"],  // ✅ Select response-only field
    "page": 1,
    "limit": 25
  }'
```

**Solution**: Filter by filterable fields, then select response-only fields using `select_columns`.

---

### Error Example 7: Invalid Date Format

**❌ Incorrect Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "range_query": {
        "must": {
          "created_at": {
            "gte": "2024-01-01"  // ERROR: Missing time component
          }
        }
      }
    }
  }'
```

**Error Response** (500 Internal Server Error):

```json
{
  "error": "ERR_ELASTICSEARCH_FAILURE: search engine returned status 400; details: ...",
  "success": false
}
```

**✅ Corrected Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "range_query": {
        "must": {
          "created_at": {
            "gte": "2024-01-01T00:00:00Z"  // ✅ Correct: ISO 8601 format
          }
        }
      }
    },
    "page": 1,
    "limit": 25
  }'
```

**Solution**: Use ISO 8601 format (RFC3339) for dates: `"YYYY-MM-DDTHH:MM:SSZ"`.

---

### Error Example 8: Invalid Range Query Type

**❌ Incorrect Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "range_query": {
        "must": {
          "employees_count": {
            "gte": "fifty"  // ERROR: Must be integer, not string
          }
        }
      }
    }
  }'
```

**Error Response** (400 Bad Request):

```json
{
  "error": "ERR_INVALID_REQUEST_BODY: the request body is invalid; check JSON syntax and required fields",
  "success": false
}
```

**✅ Corrected Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "where": {
      "range_query": {
        "must": {
          "employees_count": {
            "gte": 50  // ✅ Correct: Integer value
          }
        }
      }
    },
    "page": 1,
    "limit": 25
  }'
```

**Solution**: Use integer values for numeric range queries, not strings.

---

### Error Example 9: Missing Authentication Header

**❌ Incorrect Request**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies \
  -H "Content-Type: application/json" \
  -d '{
    "where": {
      "keyword_match": {
        "must": {
          "industries": ["Software"]
        }
      }
    }
  }'
```

**Error Response** (401 Unauthorized):

```json
{
  "error": "unauthorized",
  "message": "invalid API key"
}
```

**✅ Corrected Request**:

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

**Solution**: Always include `X-API-Key` header with valid API key.

---

### Error Example 10: Rate Limit Exceeded

**Scenario**: Making too many requests in a short time

**Error Response** (429 Too Many Requests):

```json
{
  "error": "rate limit exceeded",
  "message": "too many requests, please try again later"
}
```

**✅ Solution - Implement Exponential Backoff**:

```javascript
async function makeRequestWithRetry(url, options, maxRetries = 5) {
  let delay = 1000; // Start with 1 second
  
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    const response = await fetch(url, options);
    
    if (response.status === 429) {
      // Rate limit exceeded - wait and retry
      console.log(`Rate limited, waiting ${delay}ms before retry...`);
      await sleep(delay);
      delay = Math.min(delay * 2, 30000); // Exponential backoff, cap at 30s
      continue;
    }
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
  }
  
  throw new Error('Max retries exceeded');
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// Usage
const data = await makeRequestWithRetry('/companies', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': 'your-secret-api-key'
  },
  body: JSON.stringify({
    where: { keyword_match: { must: { industries: ["Software"] } } },
    page: 1,
    limit: 25
  })
});
```

**Solution**: Implement exponential backoff when receiving 429 errors. Wait before retrying with increasing delays.

---

### Error Example 11: Using Denormalized Fields in select_columns (Contacts)

**❌ Incorrect Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
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
      "company_name",  // ERROR: Denormalized fields NOT in select_columns!
      "company_employees_count"  // ERROR: Denormalized fields NOT in select_columns!
    ]
  }'
```

**Error Response** (500 Internal Server Error):

```json
{
  "error": "database error: column company_name does not exist",
  "success": false
}
```

**✅ Corrected Query**:

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
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
      "last_name",
      "email",
      "company_id"
    ],
    "company_config": {
      "populate": true,
      "select_columns": [
        "uuid",
        "name",  // ✅ Correct - direct name, NO company_ prefix!
        "employees_count",  // ✅ Correct - direct name, NO company_ prefix!
        "industries"
      ]
    },
    "page": 1,
    "limit": 25
  }'
```

**Solution**: Denormalized `company_*` fields are ONLY for filtering in `where` clauses. Use `company_config.select_columns` with direct field names (no prefix) to get company data in responses.

---

### Error Prevention Checklist

Before sending queries, verify:

- [ ] All required fields are present (`text_value`, `filter_key`, `search_type` for text matches)
- [ ] Field names are valid (use `/filters` endpoint to verify)
- [ ] `page` is between 1 and 10 (or use `search_after`)
- [ ] `limit` is between 1 and 100
- [ ] `search_type` is one of: `"exact"`, `"shuffle"`, `"substring"`
- [ ] Date values use ISO 8601 format: `"YYYY-MM-DDTHH:MM:SSZ"`
- [ ] Numeric range queries use integers, not strings
- [ ] `X-API-Key` header is included
- [ ] Response-only fields are not used in `where` clauses
- [ ] Denormalized `company_*` fields are not used in `select_columns` (contacts only)

**See**: [Error Handling Guide](./06-api-reference.md#error-handling) for complete error reference.

---

## Best Practices from Examples

1. **Always filter by `email_status: "verified"`** for contacts to ensure contactability
2. **Use `seniority` filters** to target decision-makers
3. **Combine multiple filter types** for precise targeting
4. **Use `search_after` for large result sets** instead of page-based pagination
5. **Filter companies first**, then find contacts at those companies for account-based approaches
6. **Use stored filters** (`direct_derived: false`) when available for better performance
7. **Sort by relevant fields** to prioritize results (e.g., revenue, seniority, creation date)
8. **Use `substring` search** for autocomplete-style queries and partial text matching
9. **Use `select_columns`** to limit returned fields and improve performance when you only need specific data
10. **Choose the right search type**: `exact` for phrases, `shuffle` for flexible word matching, `substring` for partial matches

---

## CRUD Operation Examples

> **Status**: The following CRUD operations are documented but currently only `batch-upsert` is implemented. See [CRUD Implementation Plan](../CRUD_IMPLEMENTATION_PLAN.md) for implementation details.

### Creating Companies

**Example: Create a new company**

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "name": "TechStart Inc",
    "normalized_domain": "techstart.com",
    "employees_count": 50,
    "industries": ["Software", "SaaS"],
    "country": "USA",
    "city": "San Francisco",
    "state": "CA",
    "annual_revenue": 5000000,
    "website": "https://techstart.com"
  }'
```

### Creating Contacts

**Example: Create a new contact**

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "email": "jane.smith@techstart.com",
    "first_name": "Jane",
    "last_name": "Smith",
    "title": "VP of Engineering",
    "company_id": "company-uuid-here",
    "seniority": "Executive",
    "departments": ["Engineering"],
    "email_status": "verified",
    "country": "USA"
  }'
```

### Updating Companies

**Example: Update company metrics**

```bash
curl -X PUT https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/company-uuid-here \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "employees_count": 75,
    "annual_revenue": 7500000,
    "industries": ["Software", "SaaS", "AI"]
  }'
```

### Updating Contacts

**Example: Update contact title and seniority**

```bash
curl -X PUT https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/contact-uuid-here \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "title": "Chief Technology Officer",
    "seniority": "Executive"
  }'
```

### Upsert Operations

**Example: Upsert company (create or update by domain)**

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "normalized_domain": "techstart.com",
    "name": "TechStart Inc",
    "employees_count": 75,
    "annual_revenue": 7500000
  }'
```

**Example: Upsert contact (create or update by email)**

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/contacts/upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "email": "jane.smith@techstart.com",
    "first_name": "Jane",
    "last_name": "Smith",
    "title": "CTO"
  }'
```

### Bulk Operations

**Example: Bulk upsert companies (Currently Implemented)**

```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/batch-upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "pg_companies": [
      {
        "uuid": "uuid-1",
        "name": "Company 1",
        "normalized_domain": "company1.com",
        "employees_count": 100
      },
      {
        "uuid": "uuid-2",
        "name": "Company 2",
        "normalized_domain": "company2.com",
        "employees_count": 200
      }
    ],
    "es_companies": [
      {
        "uuid": "uuid-1",
        "name": "Company 1",
        "normalized_domain": "company1.com",
        "employees_count": 100
      },
      {
        "uuid": "uuid-2",
        "name": "Company 2",
        "normalized_domain": "company2.com",
        "employees_count": 200
      }
    ]
  }'
```

**See**: 
- [Company API - Write Operations](../company.md#write-operations)
- [Contact API - Write Operations](../contacts.md#write-operations)
- [Field Reference - Fields in Write Operations](./04-filter-field-reference.md#fields-in-write-operations)

---

## Related Documentation

- [Company Filters Guide](./01-company-filters-complete-guide.md)
- [Contact Filters Guide](./02-contact-filters-complete-guide.md)
- [Combined Filters Guide](./03-combined-filters-guide.md)
- [Filter Field Reference](./04-filter-field-reference.md)
- [API Reference](./06-api-reference.md)

