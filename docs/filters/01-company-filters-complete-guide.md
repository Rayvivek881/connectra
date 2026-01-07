# Company Filters - Complete Guide

**Version**: 1.2  
**Last Updated**: 2025-01-XX

## Table of Contents

1. [Overview](#overview)
   - [VQL (Vivek Query Language) Overview](#vql-vivek-query-language-overview)
2. [Prerequisites](#prerequisites)
3. [Filter Structure](#filter-structure)
   - [Base VQL Query Structure](#base-vql-query-structure)
   - [Company Filter Field Categories](#company-filter-field-categories)
4. [Text Match Filters](#text-match-filters)
   - [Comprehensive VQL Syntax Reference](#comprehensive-vql-syntax-reference)
   - [Search Type: exact](#search-type-exact)
   - [Search Type: shuffle](#search-type-shuffle)
   - [Search Type: substring](#search-type-substring)
   - [Fuzzy Matching Details](#fuzzy-matching-details)
   - [Operator Behavior](#operator-behavior)
   - [Slop Parameter Details](#slop-parameter-details)
5. [Keyword Match Filters](#keyword-match-filters)
6. [Range Query Filters](#range-query-filters)
7. [Combined Filter Patterns](#combined-filter-patterns)
8. [Sorting and Pagination](#sorting-and-pagination)
   - [Pagination Strategies](#pagination-strategies)
   - [Page-Based Pagination](#page-based-pagination)
   - [Cursor-Based Pagination (search_after)](#cursor-based-pagination-search_after)
   - [Comparison: Page-Based vs Cursor-Based](#comparison-page-based-vs-cursor-based)
   - [Pagination Best Practices](#pagination-best-practices)
9. [Field Selection Optimization](#field-selection-optimization-with-select_columns)
10. [Real-World Use Cases](#real-world-use-cases)
11. [Field Reference](#field-reference)
12. [Best Practices](#best-practices)
13. [Error Handling](#error-handling)
14. [Related Documentation](#related-documentation)

## Overview

The Company API supports comprehensive filtering using VQL (Vivek Query Language). Filters are organized into three main categories:

- **Text Matches**: Full-text search on text fields (supports exact, shuffle, and substring search types)
- **Keyword Matches**: Exact matching on keyword/array fields
- **Range Queries**: Numeric and date range filtering

All filters can be combined using `must` (AND logic) and `must_not` (NOT logic) conditions.

> **Authentication Required**: All examples in this guide show only the JSON request body. When making actual API calls, you must include the `X-API-Key` header for authentication. See [06-api-reference.md](./06-api-reference.md#authentication) for complete HTTP request examples with headers.

### VQL (Vivek Query Language) Overview

VQL is a powerful query language that converts user-friendly JSON queries into Elasticsearch queries. It provides:

- **Flexible Text Search**: Multiple search types (exact, shuffle, substring) with fuzzy matching
- **Precise Keyword Matching**: Exact value matching for categorical data
- **Range Filtering**: Numeric and date range queries with multiple operators
- **Boolean Logic**: Complex combinations using `must` (AND) and `must_not` (NOT)
- **Performance Optimization**: Efficient query execution with proper indexing

**Key Benefits**:

- Simple JSON syntax - no need to learn Elasticsearch query DSL
- Type-safe field validation
- Automatic query optimization
- Support for complex nested conditions

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

### Company Filter Field Categories

**Text Search Fields** (use in `text_matches`):

- `name` - Company name
- `address` - Company address
- `city` - City name
- `state` - State/Province
- `country` - Country name
- `linkedin_url` - LinkedIn URL
- `website` - Website URL
- `normalized_domain` - Domain name

**Keyword Fields** (use in `keyword_match`):

- `id` - Company ID
- `industries` - Industries array
- `keywords` - Keywords array
- `technologies` - Technologies array
- `country` - Country (can be text or keyword)
- `city` - City (can be text or keyword)
- `state` - State (can be text or keyword)

**Range Query Fields**:

- `employees_count` - Employee count (integer)
- `annual_revenue` - Annual revenue (integer)
- `total_funding` - Total funding (integer)
- `created_at` - Creation date (ISO 8601 string)

---

## Text Match Filters

### Single Text Match - Name

**Example: Search by company name**

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

### Single Text Match - Name (Exact)

**Example: Exact phrase search**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "artificial intelligence",
          "filter_key": "name",
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

**Parameters**:

- `search_type: "exact"` - Phrase matching with word order
- `slop: 2` - Allows 2 words between terms

### Single Text Match - Address

**Example: Search by address**

```json
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
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Text Matches - Same Field

**Example: Multiple name searches**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "software engineer",
          "filter_key": "name",
          "search_type": "exact",
          "slop": 2
        },
        {
          "text_value": "Python developer",
          "filter_key": "name",
          "search_type": "exact",
          "slop": 1
        }
      ]
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Text Matches - Different Fields

**Example: Name and address search**

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
        },
        {
          "text_value": "New York",
          "filter_key": "address",
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

### Text Match with must_not

**Example: Include "tech", exclude "consulting"**

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
      ],
      "must_not": [
        {
          "text_value": "consulting",
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

### URL Text Search - LinkedIn

**Example: Search by LinkedIn URL**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "linkedin.com/company/acme",
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

### URL Text Search - Website

**Example: Search by website**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "acme.com",
          "filter_key": "website",
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

### Domain Text Search

**Example: Search by normalized domain**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "acme.com",
          "filter_key": "normalized_domain",
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

### Substring Text Search

**Example: Partial name matching using substring search**

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
  "page": 1,
  "limit": 25
}
```

**Use Cases for Substring Search**:

- Finding companies with partial name matches (e.g., "soft" matches "Microsoft", "Software Corp")
- Searching for partial words within company names
- Autocomplete-style searches where users type partial text
- Finding variations of company names

**Important**:

- Only the `name` field supports substring search (via ngram analysis)
- Minimum search text length is 3 characters (min_gram: 3)
- Maximum substring match is 10 characters (max_gram: 10)
- Other text fields (`address`, `city`, `state`, `country`, `linkedin_url`, `website`, `normalized_domain`) do not support substring search

**Note**: Substring search uses ngram analysis and automatically queries the `.ngram` field. For company filters, only the `name` field supports ngram matching (min_gram: 3, max_gram: 10). This enables efficient partial text matching for company names.

### Text Match Parameters Reference

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `text_value` | string | Yes | Text to search for |
| `filter_key` | string | Yes | Field name (name, address, city, etc.) |
| `search_type` | string | Yes | `"exact"`, `"shuffle"`, or `"substring"` |
| `slop` | integer | No | Word distance for exact search (default: 0) |
| `operator` | string | No | `"and"` or `"or"` for shuffle/substring search |
| `fuzzy` | boolean | No | Enable fuzzy matching for exact/shuffle (default: false) |

### Comprehensive VQL Syntax Reference

#### Search Type: `"exact"`

**Behavior**: Phrase matching with word order preservation. Matches documents where all words appear in the specified order.

**Use Cases**:

- Multi-word phrases where order matters: "artificial intelligence", "New York"
- Company names with specific word order: "Acme Software Corp"
- Addresses: "Silicon Valley"

**Parameters**:

- `slop`: Number of words that can appear between terms (default: 0)
  - `slop: 0` - Exact phrase match
  - `slop: 1` - Allows 1 word between terms (e.g., "artificial intelligence" matches "artificial general intelligence")
  - `slop: 2` - Allows 2 words between terms
- `fuzzy`: Enable typo tolerance (default: false)
  - `fuzzy: true` - Allows 1-2 character edits (e.g., "inteligence" matches "intelligence")

**Example - Exact Phrase**:

```json
{
  "text_value": "artificial intelligence",
  "filter_key": "name",
  "search_type": "exact",
  "slop": 0  // Exact phrase, no words between
}
```

**Example - Exact with Slop**:

```json
{
  "text_value": "software engineer",
  "filter_key": "name",
  "search_type": "exact",
  "slop": 2  // Matches "senior software engineer", "software development engineer"
}
```

**Example - Exact with Fuzzy**:

```json
{
  "text_value": "artifical inteligence",  // Typos
  "filter_key": "name",
  "search_type": "exact",
  "slop": 1,
  "fuzzy": true  // Tolerates typos
}
```

**Performance**: Fast for short phrases, slower for long phrases with high slop values.

#### Search Type: `"shuffle"`

**Behavior**: Word matching where order doesn't matter. Matches documents containing all specified words in any order.

**Use Cases**:

- Flexible text search: "software engineer" matches "engineer software"
- Multi-word searches where order is flexible
- General text search with typo tolerance

**Parameters**:

- `operator`: `"and"` (all words must match) or `"or"` (any word matches) (default: `"and"`)
  - `operator: "and"` - All words must be present
  - `operator: "or"` - Any word can match
- `fuzzy`: Enable typo tolerance (default: false)

**Example - Shuffle with AND**:

```json
{
  "text_value": "software technology",
  "filter_key": "name",
  "search_type": "shuffle",
  "operator": "and",  // Both "software" AND "technology" must be present
  "fuzzy": true
}
```

**Example - Shuffle with OR**:

```json
{
  "text_value": "software technology saas",
  "filter_key": "name",
  "search_type": "shuffle",
  "operator": "or",  // Any of "software" OR "technology" OR "saas" can match
  "fuzzy": true
}
```

**Example - Shuffle with Fuzzy**:

```json
{
  "text_value": "softwar technolgy",  // Typos
  "filter_key": "name",
  "search_type": "shuffle",
  "operator": "and",
  "fuzzy": true  // Tolerates typos in words
}
```

**Performance**: Generally faster than exact for multi-word searches, especially with `operator: "or"`.

#### Search Type: `"substring"`

**Behavior**: Partial text matching using ngram analysis. Matches documents containing the search text as a substring.

**Use Cases**:

- Autocomplete-style searches
- Partial word matching: "soft" matches "software", "Microsoft"
- Finding variations: "tech" matches "technology", "technologies"

**Parameters**:

- `operator`: `"and"` (all characters must match) or `"or"` (any character matches) (default: `"and"`)
  - For substring, `"and"` is typically used (all characters in sequence)
- **Minimum Length**: 3 characters for company `name` field (min_gram: 3)
- **Maximum Length**: 10 characters for efficient matching (max_gram: 10)

**Example - Substring Search**:

```json
{
  "text_value": "soft",
  "filter_key": "name",
  "search_type": "substring",
  "operator": "and"  // All characters in sequence
}
```

**Matches**: "Software", "Microsoft", "SoftTech", "SoftCorp"

**Example - Longer Substring**:

```json
{
  "text_value": "micros",
  "filter_key": "name",
  "search_type": "substring",
  "operator": "and"
}
```

**Matches**: "Microsoft", "Microsystems", "MicroCorp"

**Ngram Configuration**:

- **Company `name` field**: min_gram: 3, max_gram: 10
- **Minimum search text**: 3 characters
- **Efficient for**: 3-10 character substrings

**Performance**: Very fast for short substrings, uses pre-indexed ngram tokens.

**Limitations**:

- Only `name` field supports substring search for companies
- Minimum 3 characters required
- Maximum efficient length is 10 characters (longer searches may be slower)

### Fuzzy Matching Details

**Fuzzy Matching** is available for `"exact"` and `"shuffle"` search types only. It provides typo tolerance using Levenshtein distance.

**How It Works**:

- Allows 1-2 character edits (insertions, deletions, substitutions, transpositions)
- Automatically calculates edit distance
- Matches words with similar spelling

**Example - Fuzzy Matching**:

```json
{
  "text_value": "inteligence",  // Missing 'l'
  "filter_key": "name",
  "search_type": "shuffle",
  "fuzzy": true  // Matches "intelligence"
}
```

**Fuzzy Matching Examples**:

- "softwar" → matches "software" (1 deletion)
- "tecnology" → matches "technology" (1 substitution)
- "artifical" → matches "artificial" (1 insertion)
- "inteligence" → matches "intelligence" (1 insertion)

**Performance Impact**: Fuzzy matching is slower than exact matching but provides better user experience for typos.

### Operator Behavior

**`operator: "and"`** (Default):

- All words/characters must be present
- For shuffle: All words must appear (order doesn't matter)
- For substring: All characters must appear in sequence

**`operator: "or"`**:

- Any word/character can match
- For shuffle: Any word can appear
- For substring: Less commonly used (typically use "and")

**Example - AND vs OR**:

```json
// AND: Both words required
{
  "text_value": "software technology",
  "search_type": "shuffle",
  "operator": "and"  // Matches: "software technology", "technology software"
  // Does NOT match: "software" alone or "technology" alone
}

// OR: Either word can match
{
  "text_value": "software technology",
  "search_type": "shuffle",
  "operator": "or"  // Matches: "software", "technology", or both
}
```

### Slop Parameter Details

**Slop** is only used with `search_type: "exact"`. It controls how many words can appear between search terms.

**Slop Values**:

- `slop: 0` - Exact phrase, no words between (strictest)
- `slop: 1` - Allows 1 word between terms
- `slop: 2` - Allows 2 words between terms
- `slop: 3+` - Allows more words (use sparingly, impacts performance)

**Example - Slop Impact**:

```json
// Search: "software engineer"
// slop: 0 - Matches: "software engineer" only
// slop: 1 - Matches: "software engineer", "senior software engineer"
// slop: 2 - Matches: "software engineer", "senior software engineer", "software development engineer"
```

**Performance**: Higher slop values increase query complexity and may slow down searches.

**Search Type Guide**:

- **`"exact"`**: Phrase matching with word order. Use for multi-word phrases where order matters (e.g., "artificial intelligence")
- **`"shuffle"`**: Word matching where order doesn't matter. Use for flexible text search (e.g., "software engineer" matches "engineer software")
- **`"substring"`**: Partial text matching using ngram analysis. Use for finding partial matches within words (e.g., "soft" matches "software", "Microsoft"). For company filters, only the `name` field supports substring search via ngram (min_gram: 3, max_gram: 10). The search text must be at least 3 characters long.

---

## Prerequisites

Before using company filters, ensure you have:

- ✅ Valid API key configured
- ✅ Understanding of VQL query structure
- ✅ Knowledge of available filterable fields
- ✅ Access to the Company API endpoint

**See**: [API Reference](./06-api-reference.md) for authentication and endpoint details.

---

## Keyword Match Filters

### Single Keyword - Industries

**Example: Single industry**

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
  "limit": 25
}
```

### Multiple Industries

**Example: Multiple industries**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology", "SaaS"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Single Keyword - Technologies

**Example: Single technology**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "technologies": ["Python"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Technologies

**Example: Multiple technologies**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "technologies": ["Python", "Go", "React", "JavaScript"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Single Keyword - Keywords

**Example: Single keyword**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "keywords": ["AI"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Keywords

**Example: Multiple keywords**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "keywords": ["AI", "Machine Learning", "Cloud Computing"]
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

### Multiple Keyword Filters Combined

**Example: Industries, technologies, and country**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology"],
        "technologies": ["Python", "Go"],
        "country": ["USA", "Canada"],
        "keywords": ["AI", "Machine Learning"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Keyword Match with must_not

**Example: Include industries, exclude keywords**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology"],
        "country": ["USA"]
      },
      "must_not": {
        "keywords": ["Legacy", "Outdated"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Company ID Filter

**Example: Specific company ID**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": 1
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Company IDs

**Example: Multiple company IDs**

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [1, 2, 3, 4, 5]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

---

## Range Query Filters

### Single Range - Employees Count (gte)

**Example: Minimum employee count**

```json
{
  "where": {
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Single Range - Employees Count (Range)

**Example: Employee count range**

```json
{
  "where": {
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
  "limit": 25
}
```

### Single Range - Annual Revenue (gte)

**Example: Minimum annual revenue**

```json
{
  "where": {
    "range_query": {
      "must": {
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

### Single Range - Annual Revenue (Range)

**Example: Annual revenue range**

```json
{
  "where": {
    "range_query": {
      "must": {
        "annual_revenue": {
          "gte": 1000000,
          "lte": 50000000
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Single Range - Total Funding (gte)

**Example: Minimum total funding**

```json
{
  "where": {
    "range_query": {
      "must": {
        "total_funding": {
          "gte": 5000000
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Single Range - Total Funding (Range)

**Example: Total funding range**

```json
{
  "where": {
    "range_query": {
      "must": {
        "total_funding": {
          "gte": 1000000,
          "lte": 100000000
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Multiple Range Queries

**Example: Employee count and revenue**

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

### All Numeric Ranges Combined

**Example: Employees, revenue, and funding**

```json
{
  "where": {
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50,
          "lte": 1000
        },
        "annual_revenue": {
          "gte": 1000000,
          "lte": 50000000
        },
        "total_funding": {
          "gte": 5000000
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Date Range - Created At

**Example: Companies created in date range**

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

**Example: Companies created after date**

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

### Combined Numeric and Date Ranges

**Example: Employees, revenue, and creation date**

```json
{
  "where": {
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 100
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
  "page": 1,
  "limit": 25
}
```

### Range Query Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `gte` | Greater than or equal | `{"gte": 50}` |
| `lte` | Less than or equal | `{"lte": 1000}` |
| `gt` | Greater than | `{"gt": 50}` |
| `lt` | Less than | `{"lt": 1000}` |

**Date Format**: ISO 8601 (RFC3339) - `"2024-01-01T00:00:00Z"`

---

## Combined Filter Patterns

### Text Match + Keyword Match

**Example: Name search + industry filter**

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
    },
    "keyword_match": {
      "must": {
        "industries": ["Software", "Technology"],
        "country": ["USA"]
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Text Match + Range Query

**Example: Name search + employee count**

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
    },
    "range_query": {
      "must": {
        "employees_count": {
          "gte": 50,
          "lte": 500
        }
      }
    }
  },
  "page": 1,
  "limit": 25
}
```

### Keyword Match + Range Query

**Example: Industry + revenue filter**

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
        "annual_revenue": {
          "gte": 1000000
        },
        "employees_count": {
          "gte": 50
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
        "technologies": ["Python", "Go"],
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

---

## Sorting and Pagination

### Single Sort

**Example: Sort by revenue**

```json
{
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
    }
  ],
  "page": 1,
  "limit": 25
}
```

### Multiple Sorts

**Example: Sort by revenue then employees**

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
      "order_by": "annual_revenue",
      "order_direction": "desc"
    },
    {
      "order_by": "employees_count",
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

### Sortable Fields

**Can sort by**:

- `id`
- `employees_count`
- `annual_revenue`
- `total_funding`
- `created_at`
- `industries` (keyword field)
- `keywords` (keyword field)
- `technologies` (keyword field)

**Cannot sort by** (text fields):

- `name`
- `address`
- `city`
- `state`
- `country`
- `linkedin_url`
- `website`
- `normalized_domain`

### Pagination Strategies

Connectra supports two pagination methods: **page-based** and **cursor-based** (using `search_after`). Each has different use cases and performance characteristics.

#### Page-Based Pagination

**How It Works**:

- Uses `page` and `limit` parameters
- Elasticsearch calculates offset: `offset = (page - 1) * limit`
- Returns results starting from that offset

**Example: First page**

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
  "limit": 25
}
```

**Pagination Limits**:

- `page`: Maximum 10 (pages 1-10)
- `limit`: Maximum 100, default 25

**When to Use**:

- ✅ Small result sets (< 250 records total)
- ✅ User-facing pagination (page 1, 2, 3, etc.)
- ✅ Random access to specific pages
- ✅ Simple implementation

**Limitations**:

- ❌ Maximum 10 pages (250 records with limit 25, 1000 records with limit 100)
- ❌ Performance degrades with higher page numbers
- ❌ Not suitable for large datasets
- ❌ Deep pagination (page 8-10) is slower

**Performance Characteristics**:

- Pages 1-3: Fast
- Pages 4-7: Moderate
- Pages 8-10: Slower (Elasticsearch must scan more documents)

**Example - Multi-Page Navigation**:

```json
// Page 1
{
  "where": { "keyword_match": { "must": { "industries": ["Software"] } } },
  "page": 1,
  "limit": 25
}

// Page 2
{
  "where": { "keyword_match": { "must": { "industries": ["Software"] } } },
  "page": 2,
  "limit": 25
}

// Page 10 (maximum)
{
  "where": { "keyword_match": { "must": { "industries": ["Software"] } } },
  "page": 10,
  "limit": 25
}
```

#### Cursor-Based Pagination (search_after)

**How It Works**:

- Uses `search_after` parameter with sort field values
- Elasticsearch uses these values as a cursor to find the next page
- More efficient for large datasets
- No page limit (can paginate through millions of records)

**Example: Using search_after**

```json
{
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
      "order_by": "created_at",
      "order_direction": "asc"
    },
    {
      "order_by": "id",
      "order_direction": "asc"
    }
  ],
  "search_after": [5000000, "2024-01-15T08:00:00Z", 123],
  "limit": 25
}
```

**How to Get search_after Values**:

1. Make initial request with `order_by` (no `search_after`)
2. Get the last document from the response
3. Extract sort field values in the same order as `order_by`
4. Use those values as `search_after` for the next page

**Example - Multi-Step Cursor Pagination**:

**Step 1: Initial Request**

```json
{
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
  "limit": 25
}
```

**Response** (last document):

```json
{
  "data": [
    // ... 24 documents ...
    {
      "id": 123,
      "annual_revenue": 5000000,
      "name": "Example Corp"
    }
  ]
}
```

**Step 2: Next Page Using search_after**

```json
{
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
  "search_after": [5000000, 123],  // Values from last document
  "limit": 25
}
```

**When to Use**:

- ✅ Large result sets (> 250 records)
- ✅ Data export scenarios
- ✅ Infinite scroll implementations
- ✅ When you need to paginate beyond page 10
- ✅ Performance-critical pagination

**Best Practices**:

- ✅ Always include `id` as the last sort field for stable pagination
- ✅ Use consistent `order_by` across all pages
- ✅ Store `search_after` values from the last document
- ✅ Use `limit: 100` for exports (maximum allowed)

**Important Notes**:

- `search_after` requires `order_by` to be specified
- Values must match the sort fields in the same order
- Use `id` as tiebreaker for stable pagination
- Works with any number of sort fields

**Example - Complete Cursor Pagination Flow**:

```javascript
// Initial request
const response1 = await fetch('/companies', {
  method: 'POST',
  body: JSON.stringify({
    where: { keyword_match: { must: { industries: ["Software"] } } },
    order_by: [
      { order_by: "annual_revenue", order_direction: "desc" },
      { order_by: "id", order_direction: "asc" }
    ],
    limit: 25
  })
});

const data1 = await response1.json();
const lastDoc = data1.data[data1.data.length - 1];

// Next page
const response2 = await fetch('/companies', {
  method: 'POST',
  body: JSON.stringify({
    where: { keyword_match: { must: { industries: ["Software"] } } },
    order_by: [
      { order_by: "annual_revenue", order_direction: "desc" },
      { order_by: "id", order_direction: "asc" }
    ],
    search_after: [lastDoc.annual_revenue, lastDoc.id],  // Values from last document
    limit: 25
  })
});
```

### Comparison: Page-Based vs Cursor-Based

| Feature | Page-Based | Cursor-Based (search_after) |
|---------|------------|----------------------------|
| **Maximum Pages** | 10 pages | Unlimited |
| **Maximum Records** | 250 (limit 25) or 1000 (limit 100) | Unlimited |
| **Performance** | Degrades with page number | Consistent performance |
| **Use Case** | Small datasets, user pagination | Large datasets, exports |
| **Implementation** | Simple (`page` parameter) | Requires storing cursor values |
| **Random Access** | ✅ Yes (jump to page 5) | ❌ No (sequential only) |
| **Sorting Required** | No | ✅ Yes (required) |
| **Stability** | May change if data updates | Stable (uses sort values) |

### Pagination Best Practices

1. **For Small Datasets (< 250 records)**:
   - Use page-based pagination
   - Simple and user-friendly
   - Example: `page: 1, limit: 25`

2. **For Large Datasets (> 250 records)**:
   - Use cursor-based pagination (`search_after`)
   - Consistent performance
   - Example: `search_after: [...], limit: 100`

3. **For Data Exports**:
   - Use cursor-based pagination
   - Use maximum `limit: 100`
   - Iterate until no more results
   - Example: `search_after: [...], limit: 100`

4. **For User-Facing Pagination**:
   - Use page-based for pages 1-10
   - Switch to cursor-based if more pages needed
   - Cache `search_after` values for each page

5. **Always Include ID in Sort**:
   - Add `id` as last sort field for stable pagination
   - Prevents duplicate/missing records
   - Example: `order_by: [{order_by: "annual_revenue"}, {order_by: "id"}]`

### Common Pagination Errors

**Error: Page Number Exceeded**

```json
{
  "error": "ERR_PAGE_OUT_OF_RANGE: the requested page number is beyond the available range; verify total pages before requesting",
  "success": false
}
```

**Solution**: Use `search_after` for pagination beyond page 10

**Error: Page Size Exceeded**

```json
{
  "error": "ERR_PAGE_SIZE_EXCEEDED: the requested page size surpasses the maximum allowed limit; consider using pagination with smaller batches",
  "success": false
}
```

**Solution**: Reduce `limit` to maximum 100

**Error: Missing search_after Values**

- Ensure `search_after` values match `order_by` fields
- Extract values from the last document in previous response
- Include all sort fields in the same order

---

## Real-World Use Cases

### Lead Generation

**Find high-value software companies**

```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "software technology",
          "filter_key": "name",
          "search_type": "shuffle",
          "fuzzy": true,
          "operator": "and"
        }
      ]
    },
    "keyword_match": {
      "must": {
        "industries": ["Software", "SaaS", "Technology"],
        "country": ["USA", "Canada", "UK"],
        "technologies": ["Cloud", "AI", "Machine Learning"]
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
    }
  ],
  "page": 1,
  "limit": 50
}
```

### Competitive Analysis

**Find competitors in same space**

```json
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

### Partnership Opportunities

**Find potential partners**

```json
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

### Market Research

**Find companies by location and size**

```json
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

### Investment Targeting

**Find well-funded companies**

```json
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

---

## Field Reference

### Filterable Fields (Elasticsearch Index)

These fields are indexed in Elasticsearch and can be used in filter queries:

| Field | Type | Search Type | Sortable | Ngram Support | Description |
|-------|------|-------------|----------|---------------|-------------|
| `id` | integer | keyword | Yes | No | Company ID |
| `name` | string | text | No | Yes (3-10) | Company name |
| `address` | string | text | No | No | Company address |
| `city` | string | text | No | No | City name |
| `state` | string | text | No | No | State/Province |
| `country` | string | text | No | No | Country name |
| `industries` | array[string] | keyword | Yes | No | Industries array |
| `keywords` | array[string] | keyword | Yes | No | Keywords array |
| `technologies` | array[string] | keyword | Yes | No | Technologies array |
| `employees_count` | integer | range | Yes | No | Employee count |
| `annual_revenue` | integer | range | Yes | No | Annual revenue (in cents) |
| `total_funding` | integer | range | Yes | No | Total funding (in cents) |
| `linkedin_url` | string | text | No | No | LinkedIn URL |
| `website` | string | text | No | No | Website URL |
| `normalized_domain` | string | text | No | No | Domain name |
| `created_at` | datetime | range | Yes | No | Creation date (ISO 8601) |

**Ngram Support**: Only the `name` field has ngram analysis enabled (min_gram: 3, max_gram: 10), allowing efficient substring/partial matching. Other text fields do not support substring search.

**Note**: `city`, `state`, and `country` are text fields in Elasticsearch. They can be used in `text_matches` for fuzzy/flexible search, but for exact matching, use `keyword_match` with the exact values.

### Response-Only Fields (PostgreSQL Only)

> **⚠️ Important**: These fields are stored in PostgreSQL but are **NOT** indexed in Elasticsearch and **cannot be used in filters**. They are only returned in API responses:

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

**Important**: These fields cannot be used in `where` clauses. They are only available in response data and can be selected using `select_columns` if needed.

### Filter Availability

**Direct-Derived Filters** (`direct_derived: true`):

- `address`
- `annual_revenue`
- `employees_count`
- `linkedin_url`
- `normalized_domain`
- `total_funding`
- `website`

**Stored Filters** (`direct_derived: false`):

- `city`
- `country`
- `industries`
- `keywords`
- `state`
- `technologies`
- `uuid` (displayed as "Name" in filter UI, but not filterable in Elasticsearch)

---

## Best Practices

1. **Use Appropriate Search Types**:
   - Use `exact` for phrases where word order matters
   - Use `shuffle` for general text search
   - Use `substring` for partial text matching and autocomplete-style queries
   - Enable `fuzzy` for user-generated search terms (exact and shuffle only)

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
   - Use `select_columns` to limit returned fields and reduce payload size

### Field Selection Optimization with select_columns

The `select_columns` parameter allows you to specify which fields are returned in the API response, significantly reducing payload size and improving performance.

**How It Works**:

1. Elasticsearch search executes first (finds matching documents)
2. PostgreSQL retrieves full records for matched document IDs
3. If `select_columns` is specified, only those fields are returned
4. Response payload is reduced, improving transfer speed

**Performance Impact**:

- **Without `select_columns`**: All fields returned (~30+ fields per company)
- **With `select_columns`**: Only specified fields returned
- **Payload Reduction**: 50-80% reduction for typical list views
- **Transfer Speed**: Faster response times, especially for large result sets

**Example - Optimized List View**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "industries": ["Software"]
      }
    }
  },
  "select_columns": [
    "id",
    "name",
    "employees_count",
    "annual_revenue",
    "industries",
    "country",
    "website"
  ],
  "page": 1,
  "limit": 50
}
```

**Example - Full Detail View**:

```json
{
  "where": {
    "keyword_match": {
      "must": {
        "id": [1, 2, 3]
      }
    }
  },
  "select_columns": [
    "id",
    "uuid",
    "name",
    "employees_count",
    "industries",
    "keywords",
    "technologies",
    "annual_revenue",
    "total_funding",
    "address",
    "city",
    "state",
    "country",
    "linkedin_url",
    "website",
    "normalized_domain",
    "created_at",
    "facebook_url",
    "twitter_url",
    "phone_number",
    "latest_funding",
    "latest_funding_amount",
    "last_raised_at"
  ],
  "page": 1,
  "limit": 25
}
```

**Best Practices**:

- ✅ Use `select_columns` for list views (8-15 fields)
- ✅ Select only fields you need
- ✅ Omit `select_columns` for exports (get all fields)
- ✅ Always include `id` in `select_columns`
- ❌ Don't select fields you won't use

**See**: [Select Columns Guide](./select_columns_filter.md) for complete documentation on field selection optimization.

5. **Combining Filters**:
   - All filters in `must` use AND logic
   - Filters in `must_not` exclude matching records
   - Combine text, keyword, and range filters for precise results

---

## Error Handling

### Common Errors and Solutions

#### 400 Bad Request - Invalid Request Body

**Error Response**:

```json
{
  "error": "invalid request body",
  "success": false
}
```

**Common Causes**:

- Malformed JSON in request body
- Missing required fields
- Invalid field types

**Solution**:

- Validate JSON syntax
- Ensure all required fields are present
- Check field types match expected values

#### 400 Bad Request - Page Size Exceeded

**Error Response**:

```json
{
  "error": "ERR_PAGE_SIZE_EXCEEDED: the requested page size surpasses the maximum allowed limit; consider using pagination with smaller batches",
  "success": false
}
```

**Cause**: `limit` parameter exceeds maximum value (100)

**Solution**:

```json
{
  "limit": 100  // Maximum allowed
}
```

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
- Elasticsearch cluster unavailable
- Query syntax errors

**Solution**:

- Verify field names using `/filters` endpoint
- Check Elasticsearch cluster health
- Review query structure

**See**: [Error Handling Guide](./06-api-reference.md#error-handling) for complete error reference and troubleshooting.

---

## Write Operations (CRUD)

> **Status**: The following CRUD operations are documented but currently only `batch-upsert` is implemented. See [CRUD Implementation Plan](../CRUD_IMPLEMENTATION_PLAN.md) for implementation details.

### Create Company

Create a new company record with automatic Elasticsearch indexing.

**Endpoint**: `POST /companies/create`

**Request**:
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/create \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "name": "Acme Software Corp",
    "normalized_domain": "acme.com",
    "employees_count": 120,
    "industries": ["Software", "Technology"],
    "country": "USA"
  }'
```

**Response** (201 Created):
```json
{
  "data": {
    "uuid": "c0a8012e-1111-2222-3333-444455556666",
    "name": "Acme Software Corp",
    "normalized_domain": "acme.com",
    "employees_count": 120,
    "industries": ["Software", "Technology"],
    "country": "USA",
    "created_at": "2025-12-24T10:30:00Z"
  },
  "success": true
}
```

### Update Company

Update an existing company by UUID.

**Endpoint**: `PUT /companies/:uuid`

**Request**:
```bash
curl -X PUT https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/c0a8012e-1111-2222-3333-444455556666 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "employees_count": 150,
    "annual_revenue": 7500000
  }'
```

### Delete Company

Soft delete a company by UUID.

**Endpoint**: `DELETE /companies/:uuid`

**Request**:
```bash
curl -X DELETE https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/c0a8012e-1111-2222-3333-444455556666 \
  -H "X-API-Key: your-secret-api-key"
```

### Upsert Company

Create or update a company (identified by UUID or normalized_domain).

**Endpoint**: `POST /companies/upsert`

**Request**:
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "normalized_domain": "acme.com",
    "name": "Acme Software Corp",
    "employees_count": 120
  }'
```

### Bulk Upsert Companies

Efficiently create or update multiple companies.

**Endpoint**: `POST /companies/batch-upsert` (Currently Implemented)

**Request**:
```bash
curl -X POST https://iarj32v8e1.execute-api.us-east-1.amazonaws.com/companies/batch-upsert \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "pg_companies": [
      {
        "uuid": "c0a8012e-1111-2222-3333-444455556666",
        "name": "Acme Software Corp",
        "normalized_domain": "acme.com"
      }
    ],
    "es_companies": [
      {
        "uuid": "c0a8012e-1111-2222-3333-444455556666",
        "name": "Acme Software Corp",
        "normalized_domain": "acme.com"
      }
    ]
  }'
```

**See**: [Company API - Write Operations](../company.md#write-operations) for complete CRUD documentation.

---

## Related Documentation

### Filter Documentation

- [Contact Filters Guide](./02-contact-filters-complete-guide.md) - Complete contact filtering guide
- [Combined Filters Guide](./03-combined-filters-guide.md) - Account-based filtering strategies
- [Filter Field Reference](./04-filter-field-reference.md) - Complete field reference
- [Examples and Use Cases](./05-examples-use-cases.md) - Real-world examples
- [API Reference](./06-api-reference.md) - Complete API endpoint reference
- [Select Columns Guide](./select_columns_filter.md) - Field selection optimization

### Main Documentation

- [System Documentation](../system.md) - System architecture and setup
- [Company API](../company.md) - Company API documentation
- [Contact API](../contacts.md) - Contact API documentation
- [Main README](../README.md) - Documentation index

---

**Last Updated**: 2025-01-XX  
**Version**: 1.2
