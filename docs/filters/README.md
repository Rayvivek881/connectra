# Filter Documentation Index

## Overview

This directory contains comprehensive documentation for all filter capabilities in the Connectra API. The documentation covers company filters, contact filters, combined filtering strategies, field references, examples, and API reference.

> **Note**: All filter endpoints require authentication using the `X-API-Key` header. See [API Reference](./06-api-reference.md#authentication) for details.

## Documentation Files

### 1. [Company Filters Complete Guide](./01-company-filters-complete-guide.md)

**Comprehensive guide to all company filter capabilities**

- Complete filter structure overview
- Text match filters (name, address, location, URLs) - supports exact, shuffle, and substring search
- Keyword match filters (industries, technologies, keywords)
- Range query filters (employees, revenue, funding, dates)
- Field selection with `select_columns`
- Combined filter patterns
- Sorting and pagination
- Real-world use cases
- Field reference

**Best for**: Understanding all possible company filter combinations

---

### 2. [Contact Filters Complete Guide](./02-contact-filters-complete-guide.md)

**Comprehensive guide to all contact filter capabilities**

- Complete filter structure overview
- Text match filters (names, titles, location, LinkedIn) - supports exact, shuffle, and substring search
- Keyword match filters (departments, seniority, email status)
- Range query filters (dates)
- Denormalized company fields (`company_*`) for filtering contacts by company attributes
- Field selection with `select_columns`
- Combined filter patterns
- Sorting and pagination
- Real-world use cases
- Field reference

**Best for**: Understanding all possible contact filter combinations, including filtering by company attributes

---

### 3. [Combined Filters Guide](./03-combined-filters-guide.md)

**Guide to using company and contact filters together**

- Direct company field filtering (single query using denormalized `company_*` fields)
- Company-first filtering approach (two-step process)
- Contact-first filtering approach
- Account-based filtering patterns
- Cross-entity filtering strategies
- Best practices for combined filters
- Common combined filter patterns

**Best for**: Account-based marketing, lead generation, competitive intelligence

---

### 4. [Filter Field Reference](./04-filter-field-reference.md)

**Complete reference for all fields**

- Company field reference (text, keyword, range)
- Contact field reference (text, keyword, range)
- Denormalized company fields in contact index (`company_*` prefix)
- Response-only fields (PostgreSQL only, not filterable)
- Field type explanations
- Ngram fields documentation with accurate min/max values
- Filter availability (direct-derived vs stored)
- Sorting capabilities
- Field usage examples

**Best for**: Quick reference when building queries, understanding which fields are filterable vs response-only

---

### 5. [Examples and Use Cases](./05-examples-use-cases.md)

**Real-world examples and scenarios**

- Company filter examples
- Contact filter examples
- Combined use cases
- Real-world scenarios (B2B sales, recruiting, email marketing, etc.)
- Best practices from examples

**Best for**: Learning from practical examples

---

### 6. [API Reference](./06-api-reference.md)

**Complete API endpoint reference**

- Company filter endpoints
- Company CRUD operations (create, update, delete, upsert, bulk)
- Contact filter endpoints
- Contact CRUD operations (create, update, delete, upsert, bulk)
- Jobs API endpoints
- Authentication and error handling
- Deployment information (Lambda and server modes)
- Request/response formats
- Error handling
- Rate limiting
- Authentication
- API endpoint summary

**Best for**: Implementation and integration

---

### 7. [Jobs API Guide](./jobs.md)

**Complete guide to background job processing**

- Job types (CSV import/export)
- Job states and lifecycle
- API endpoints for job creation and listing
- Job runners (first_time vs retry)
- Configuration and best practices
- Error handling and retry logic
- Examples for common scenarios

**Best for**: Understanding asynchronous job processing for large data operations

---

## Quick Start Guide

### For Developers

1. **Start with**: [API Reference](./06-api-reference.md) - Understand the endpoints
2. **Then read**: [Filter Field Reference](./04-filter-field-reference.md) - Know what fields are available
3. **Review**: [Examples and Use Cases](./05-examples-use-cases.md) - See practical examples
4. **Deep dive**: [Company Filters Guide](./01-company-filters-complete-guide.md) or [Contact Filters Guide](./02-contact-filters-complete-guide.md) - Master the filters

### For Business Users

1. **Start with**: [Examples and Use Cases](./05-examples-use-cases.md) - See what's possible
2. **Then read**: [Combined Filters Guide](./03-combined-filters-guide.md) - Learn filtering strategies
3. **Reference**: [Filter Field Reference](./04-filter-field-reference.md) - Understand available filters

### For Account-Based Marketing

1. **Start with**: [Combined Filters Guide](./03-combined-filters-guide.md) - Account-based patterns
2. **Review**: [Examples and Use Cases](./05-examples-use-cases.md) - Real-world scenarios
3. **Reference**: [Company Filters Guide](./01-company-filters-complete-guide.md) and [Contact Filters Guide](./02-contact-filters-complete-guide.md) - Detailed field information

---

## Key Concepts

### Filter Types

1. **Text Matches**: Full-text search on text fields using VQL (Vivek Query Language)
   - `search_type: "exact"` - Phrase matching with word order
   - `search_type: "shuffle"` - Word matching (order doesn't matter)
   - `search_type: "substring"` - Partial text matching using ngram analysis
   - Supports fuzzy matching for typos (exact and shuffle only)
   - **Ngram Configuration**:
     - Company: Only `name` field has ngram (min_gram: 3, max_gram: 10)
     - Contact: `first_name`, `last_name`, `title` have ngram (min_gram: 5, max_gram: 10)
     - Denormalized: `company_name`, `company_website`, `company_normalized_domain` have ngram (min_gram: 5, max_gram: 10)

2. **Keyword Matches**: Exact matching on keyword/array fields
   - Single values or arrays
   - Faster than text searches
   - Use for exact filtering

3. **Range Queries**: Numeric and date range filtering
   - Operators: `gte`, `lte`, `gt`, `lt`
   - Dates in ISO 8601 format

### Filter Logic

- **`must`**: AND logic - all conditions must match
- **`must_not`**: NOT logic - exclude matching records
- All filter types can be combined

### Filter Availability

- **Filterable Fields**: Fields indexed in Elasticsearch that can be used in filter queries
- **Response-Only Fields**: Fields stored in PostgreSQL but not indexed (e.g., `facebook_url`, `twitter_url`, `phone_number`, `stage`). These cannot be used in filters but are available in API responses.
- **Direct-Derived**: Filter values extracted from database in real-time
- **Stored**: Filter values pre-computed and cached for faster access

### Denormalized Company Fields in Contact Index

The contact index includes denormalized company data with the `company_*` prefix, allowing you to filter contacts directly by company attributes (e.g., `company_industries`, `company_employees_count`, `company_annual_revenue`) without needing a separate company query. This enables efficient single-query account-based filtering.

### Sorting

- Only keyword fields and date fields can be sorted
- Text fields cannot be sorted (analyzed in Elasticsearch)
- Denormalized company fields in contact index are not sortable

---

## Common Use Cases

### Lead Generation

- Find high-value companies
- Identify decision-makers
- Filter by industry, size, revenue
- Target verified contacts

### Competitive Intelligence

- Find competitors
- Identify key personnel
- Analyze technology stacks
- Monitor market segments

### Partnership Development

- Find complementary businesses
- Identify partnership contacts
- Filter by technology alignment
- Target growth-stage companies

### Market Research

- Analyze market segments
- Identify industry experts
- Geographic analysis
- Technology trend analysis

### Investment Targeting

- Find well-funded companies
- Identify executives
- Filter by growth metrics
- Target high-value accounts

---

## Best Practices

1. **Always filter by `email_status: "verified"`** for contacts
2. **Use `seniority` filters** to target decision-makers
3. **Combine multiple filter types** for precise targeting
4. **Use `search_after`** for large result sets
5. **Use denormalized company fields** (`company_*`) for single-query account-based filtering when you only need contacts
6. **Filter companies first**, then find contacts when you also need company data
7. **Use stored filters** when available for better performance
8. **Sort by relevant fields** to prioritize results
9. **Use `substring` search** for autocomplete-style queries and partial text matching
   - Company `name`: minimum 3 characters (min_gram: 3)
   - Contact `first_name`, `last_name`, `title`: minimum 5 characters (min_gram: 5)
   - Denormalized `company_name`, `company_website`, `company_normalized_domain`: minimum 5 characters (min_gram: 5)
10. **Use `select_columns`** to limit returned fields and improve performance
11. **Choose the right search type**: Use `exact` for phrases, `shuffle` for flexible word matching, `substring` for partial matches
12. **Remember response-only fields** (like `facebook_url`, `twitter_url`, `phone_number`, `stage`) cannot be used in filters but are available in responses

---

## Write Operations

> **New Feature** (December 2024): Connectra now supports full CRUD operations for data management.

**Available Write Operations**:

- **Create**: `POST /companies/create`, `POST /contacts/create`
- **Update**: `PUT /companies/:uuid`, `PUT /contacts/:uuid`
- **Delete**: `DELETE /companies/:uuid`, `DELETE /contacts/:uuid`
- **Upsert**: `POST /companies/upsert`, `POST /contacts/upsert`
- **Bulk**: `POST /companies/bulk`, `POST /contacts/bulk`

**Key Features**:

- Automatic Elasticsearch indexing on all writes
- Comprehensive validation layer
- Bulk operations for efficient data imports
- Soft delete (preserves data with `deleted_at` timestamp)

**See**: 

- [Company API - Write Operations](../company.md#write-operations)
- [Contact API - Write Operations](../contacts.md#write-operations)
- [System Documentation - Write Operations](../system.md#6-write-operations-and-data-management)

---

## Related Documentation

- [Company API Documentation](../company.md)
- [Contact API Documentation](../contacts.md)
- [System Documentation](../system.md)

---

## Support

For questions or issues:

1. Review the relevant documentation file
2. Check the [Examples and Use Cases](./05-examples-use-cases.md)
3. Refer to the [API Reference](./06-api-reference.md) for endpoint details

---

## Documentation Structure

```
filters/
├── README.md (this file)
├── 01-company-filters-complete-guide.md
├── 02-contact-filters-complete-guide.md
├── 03-combined-filters-guide.md
├── 04-filter-field-reference.md
├── 05-examples-use-cases.md
├── 06-api-reference.md
├── jobs.md
└── select_columns_filter.md
```

---

**Last Updated**: 2025-01-XX

**Version**: 1.2

**Recent Updates**:

- **2025-12-24**: ✅ **Write Operations Added**: Full CRUD operations (create, update, delete, upsert, bulk) for contacts and companies - see [Company API](../company.md#write-operations) and [Contact API](../contacts.md#write-operations)
- **2025-01-XX**: Enhanced authentication and security documentation
- **2025-01-XX**: Added comprehensive VQL syntax reference
- **2025-01-XX**: Enhanced pagination documentation with cursor-based strategies
- **2025-01-XX**: Added field selection optimization documentation
- **2025-01-XX**: Standardized documentation structure across all files
- **2025-12-17**: Fixed critical spelling error: `direct_drived` → `direct_derived` throughout codebase and documentation
- **2025-12-17**: Enhanced `select_columns` documentation with clearer explanation of PostgreSQL-only behavior
- **2025-12-17**: Added warning boxes for response-only fields to prevent filter usage errors
- Updated ngram configuration documentation with accurate min/max values per index
- Added comprehensive documentation for denormalized company fields in contact index
- Documented all response-only fields (PostgreSQL only, not filterable)
- Added distinction between filterable and response-only fields throughout documentation
- Updated examples to show denormalized company field filtering
- Added `select_columns` usage examples with response-only fields
- Clarified field types match Elasticsearch mappings exactly
- Added examples showing single-query account-based filtering using denormalized fields

