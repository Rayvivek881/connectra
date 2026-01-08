# Connectra API Documentation

This document provides an overview of all available APIs with request/response examples for frontend integration.

## Base URL
```
http://localhost:8000
```

## Authentication
All endpoints require Bearer token authentication:
```
Authorization: Bearer <your_token>
```

---

## 🔄 Understanding UPSERT Operations

### What is UPSERT?
**UPSERT** = **UP**date + in**SERT**

UPSERT is a database operation that:
- **INSERTS** a new record if it doesn't exist
- **UPDATES** the existing record if it already exists

The system determines whether to insert or update based on the **UUID** (unique identifier).

### How UUID Works
- **Companies**: UUID is generated from `name + linkedin_url`
- **Contacts**: UUID is generated from `first_name + last_name + linkedin_url`

If you send the same company/contact data twice, the second request will **update** the existing record instead of creating a duplicate.

### Supported Operations

| Operation | Supported | How to Achieve |
|-----------|-----------|----------------|
| **CREATE (Insert)** | ✅ Yes | Send data without existing UUID - new record created |
| **READ (Get)** | ✅ Yes | Use filter APIs (`POST /companies/`, `POST /contacts/`) |
| **UPDATE** | ✅ Yes | Send data with same identifier fields - existing record updated |
| **DELETE** | ❌ No | Not supported via API (soft delete managed internally) |

### UPSERT Examples

**Insert New Record (UUID doesn't exist):**
```json
// POST /companies/batch-upsert
[{
  "name": "New Company",
  "linkedin_url": "https://linkedin.com/company/new-company",
  "employees_count": 100
}]
// Result: New company created with auto-generated UUID
```

**Update Existing Record (UUID exists):**
```json
// POST /companies/batch-upsert
[{
  "name": "New Company",  // Same name
  "linkedin_url": "https://linkedin.com/company/new-company",  // Same LinkedIn
  "employees_count": 200  // Updated value
}]
// Result: Existing company updated (employees_count changed from 100 to 200)
```

**Explicit UUID Update:**
```json
// POST /companies/batch-upsert
[{
  "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",  // Existing UUID
  "name": "Updated Company Name",
  "employees_count": 300
}]
// Result: Record with this UUID is updated
```

---

## 📁 Common APIs (`http://localhost:8000/common`)

### 1. Batch Upsert (CSV Data)
**POST** `http://localhost:8000/common/batch-upsert`

Insert companies and contacts from raw CSV-like data.

**Request:**
```json
{
  "data": [
    {
      "company": "Acme Corporation",
      "company_linkedin_url": "https://linkedin.com/company/acme-corp",
      "email": "john.doe@acme.com",
      "first_name": "John",
      "last_name": "Doe",
      "title": "Senior Software Engineer",
      "departments": "Engineering, Product",
      "employees": "500",
      "industry": "Technology, Software",
      "company_city": "San Francisco",
      "company_state": "California",
      "company_country": "United States",
      "mobile_phone": "+1-555-123-4567",
      "person_linkedin_url": "https://linkedin.com/in/johndoe"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Batch upsert successful",
  "company_uuids": ["a1b2c3d4-e5f6-7890-abcd-ef1234567890"],
  "contact_uuids": ["c3d4e5f6-a7b8-9012-cdef-123456789012"]
}
```

---

### 2. Get Upload URL
**GET** `http://localhost:8000/common/upload-url?filename=contacts.csv`

Generate presigned S3 URL for file upload.

**Response:**
```json
{
  "success": true,
  "upload_url": "https://s3.amazonaws.com/bucket/...",
  "s3_key": "uploads/uuid_contacts.csv",
  "expires_in": "24h0m0s"
}
```

---

### 3. Create Job
**POST** `http://localhost:8000/common/jobs/create`

Create background job for CSV import or data export.

**Request (Import CSV):**
```json
{
  "job_type": "insert_csv_file",
  "job_data": {
    "s3_key": "uploads/uuid_contacts.csv",
    "s3_bucket": "your-bucket"
  },
  "retry_count": 3
}
```

**Request (Export CSV):**
```json
{
  "job_type": "export_csv_file",
  "job_data": {
    "s3_bucket": "your-bucket",
    "service": "contact",
    "vql": {
      "where": {
        "keyword_match": {
          "must": { "country": ["united states"] }
        }
      },
      "select_columns": ["first_name", "last_name", "email"],
      "limit": 1000
    }
  },
  "retry_count": 2
}
```

**Response:**
```json
{
  "success": true,
  "message": "Job created successfully"
}
```

---

### 4. List Jobs
**POST** `http://localhost:8000/common/jobs`

**Request:**
```json
{
  "job_type": "insert_csv_file",
  "status": ["open", "in_queue", "processing", "completed", "failed"],
  "limit": 20
}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "uuid": "job-uuid-1234",
      "job_type": "insert_csv_file",
      "status": "completed",
      "job_response": {
        "messages": "Successfully imported 150 records"
      },
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

---

### 5. Get Filters
**GET** `http://localhost:8000/common/:service/filters`

Get available filters for a service.

**Example:** `GET http://localhost:8000/common/contact/filters`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "key": "country",
      "filter_type": "keyword",
      "display_name": "Country",
      "active": true
    },
    {
      "key": "seniority",
      "filter_type": "keyword",
      "display_name": "Seniority Level",
      "active": true
    }
  ]
}
```

---

### 6. Get Filter Data
**POST** `http://localhost:8000/common/:service/filters/data`

Get filter values with search.

**Request:**
```json
{
  "filter_key": "country",
  "search_text": "united",
  "page": 1,
  "limit": 20
}
```

**Response:**
```json
{
  "success": true,
  "data": [
    { "value": "united states", "display_value": "United States" },
    { "value": "united kingdom", "display_value": "United Kingdom" }
  ]
}
```

---

## 🏢 Companies APIs (`http://localhost:8000/companies`)

### 1. Get Companies By Filter
**POST** `http://localhost:8000/companies/`

**Request:**
```json
{
  "where": {
    "text_matches": {
      "must": [
        {
          "text_value": "artificial intelligence",
          "filter_key": "keywords",
          "search_type": "match_phrase"
        }
      ]
    },
    "keyword_match": {
      "must": {
        "country": ["united states"],
        "industries": ["technology"]
      }
    },
    "range_query": {
      "must": {
        "employees_count": { "gte": 100, "lte": 5000 }
      }
    }
  },
  "order_by": [
    { "order_by": "employees_count", "order_direction": "desc" }
  ],
  "select_columns": ["uuid", "name", "industries", "employees_count", "website"],
  "limit": 25
}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "Acme Corporation",
      "industries": ["Technology", "Software"],
      "employees_count": 500,
      "website": "https://acme.com",
      "cursor": ["500", "Acme Corporation"]
    }
  ]
}
```

---

### 2. Get Companies Count
**POST** `http://localhost:8000/companies/count`

**Request:**
```json
{
  "where": {
    "keyword_match": {
      "must": { "country": ["united states"] }
    }
  }
}
```

**Response:**
```json
{
  "success": true,
  "count": 1542
}
```

---

### 3. Batch Upsert Companies
**POST** `http://localhost:8000/companies/batch-upsert`

**Request:**
```json
[
  {
    "name": "Acme Corporation",
    "employees_count": 500,
    "industries": ["Technology", "Software"],
    "city": "San Francisco",
    "state": "California",
    "country": "United States",
    "linkedin_url": "https://linkedin.com/company/acme-corp",
    "website": "https://acme.com"
  }
]
```

**Response:**
```json
{
  "success": true,
  "company_uuids": ["a1b2c3d4-e5f6-7890-abcd-ef1234567890"]
}
```

---

## 👤 Contacts APIs (`http://localhost:8000/contacts`)

### 1. Get Contacts By Filter
**POST** `http://localhost:8000/contacts/`

**Request (with company population):**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "seniority": ["senior", "manager"],
        "departments": ["engineering"]
      }
    }
  },
  "select_columns": ["uuid", "first_name", "last_name", "email", "title"],
  "company_config": {
    "populate": true,
    "select_columns": ["uuid", "name", "website"]
  },
  "limit": 25
}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "uuid": "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "first_name": "john",
      "last_name": "doe",
      "email": "john.doe@acme.com",
      "title": "Senior Software Engineer",
      "company": {
        "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "name": "Acme Corporation",
        "website": "https://acme.com"
      },
      "cursor": ["doe", "john"]
    }
  ]
}
```

---

### 2. Get Contacts Count
**POST** `http://localhost:8000/contacts/count`

**Request:**
```json
{
  "where": {
    "keyword_match": {
      "must": {
        "country": ["united states"],
        "seniority": ["senior"]
      }
    }
  }
}
```

**Response:**
```json
{
  "success": true,
  "count": 8743
}
```

---

### 3. Batch Upsert Contacts
**POST** `http://localhost:8000/contacts/batch-upsert`

**Request:**
```json
[
  {
    "first_name": "John",
    "last_name": "Doe",
    "company_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "email": "john.doe@acme.com",
    "title": "Senior Software Engineer",
    "departments": ["Engineering"],
    "seniority": "senior",
    "city": "San Francisco",
    "country": "United States",
    "linkedin_url": "https://linkedin.com/in/johndoe"
  }
]
```

**Response:**
```json
{
  "success": true,
  "contact_uuids": ["c3d4e5f6-a7b8-9012-cdef-123456789012"]
}
```

---

## 📖 VQL Query Reference

### Query Structure
```json
{
  "where": {
    "text_matches": { "must": [], "must_not": [] },
    "keyword_match": { "must": {}, "must_not": {} },
    "range_query": { "must": {}, "must_not": {} }
  },
  "order_by": [{ "order_by": "field", "order_direction": "asc|desc" }],
  "cursor": ["value1", "value2"],
  "select_columns": ["field1", "field2"],
  "company_config": { "populate": true, "select_columns": [] },
  "page": 1,
  "limit": 25
}
```

### Text Match Types
| Type | Description |
|------|-------------|
| `match` | Basic word matching |
| `match_phrase` | Phrase matching with word order |
| `match_phrase_prefix` | Prefix matching for autocomplete |

### Range Operators
| Operator | Description |
|----------|-------------|
| `gte` | Greater than or equal |
| `gt` | Greater than |
| `lte` | Less than or equal |
| `lt` | Less than |

---

## ❌ Error Response Format
All error responses follow this format:
```json
{
  "success": false,
  "error": "ERR_CODE: Error description with details"
}
```

### Common Error Codes
| Code | Description |
|------|-------------|
| `ERR_EMPTY_PAYLOAD` | Data array is empty |
| `ERR_BATCH_TOO_LARGE` | Batch size exceeds limit |
| `ERR_PAGE_SIZE_EXCEEDED` | Page size too large |
| `ERR_MISSING_JOB_TYPE` | Job type not provided |
| `ERR_MISSING_FILENAME` | Filename query param missing |

---

## 📝 Notes for Frontend Integration

1. **UUID Generation**: UUIDs are auto-generated if not provided:
   - Companies: `name + linkedin_url`
   - Contacts: `first_name + last_name + linkedin_url`

2. **Cursor Pagination**: Use `cursor` values from response for efficient pagination with `order_by`.

3. **Data Cleaning**: All string fields are automatically cleaned (trimmed, normalized).

4. **Lowercase Fields**: Some fields are stored lowercase: `city`, `state`, `country`, `email`, `linkedin_url`, `seniority`, `email_status`.

5. **Array Fields**: Use comma-separated values in CSV data for arrays (industries, departments, technologies).

6. **Valid Service Values**: Use `"contact"` or `"company"` (singular) for service endpoints.

7. **Valid Job Types**: `"insert_csv_file"`, `"export_csv_file"`

8. **Valid Job Statuses**: `"open"`, `"in_queue"`, `"processing"`, `"completed"`, `"failed"`, `"retry_in_queued"`, `"retrying"`

