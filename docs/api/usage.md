# Usage API

## Overview

The Usage API provides endpoints for tracking and retrieving feature usage information. It allows users to monitor their usage of various features and track consumption of credits.

**Related Documentation:**
- [Auth API](./auth.md) - For authentication
- [User API](./user.md) - For user profile management
- [Billing API](./billing.md) - For subscription and billing

**Base URL**: `/api/v2/usage`

**Authentication**: All endpoints require JWT Bearer token authentication.

---

## Endpoints

### 1. Get Current Usage

Get current feature usage for the authenticated user.

**Endpoint**: `GET /api/v2/usage/current/`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Request Body**: None

**Response** (200 OK):

```json
{
  "AI_CHAT": {
    "used": 50,
    "limit": 100
  },
  "BULK_EXPORT": {
    "used": 25,
    "limit": 50
  },
  "EMAIL_FINDER": {
    "used": 200,
    "limit": 500
  },
  "DATA_SEARCH": {
    "used": 150,
    "limit": 1000
  },
  "LINKEDIN": {
    "used": 75,
    "limit": 200
  }
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `[feature]` | object | Feature usage data (key is feature name) |
| `[feature].used` | integer | Number of times feature has been used |
| `[feature].limit` | integer | Usage limit for this feature |

**Supported Features**:

The following features are tracked:

- `AI_CHAT` - AI chat interactions
- `BULK_EXPORT` - Bulk data exports
- `API_KEYS` - API key usage
- `TEAM_MANAGEMENT` - Team management operations
- `EMAIL_FINDER` - Email finder searches
- `VERIFIER` - Email verification
- `LINKEDIN` - LinkedIn operations
- `DATA_SEARCH` - Data search operations
- `ADVANCED_FILTERS` - Advanced filter usage
- `AI_SUMMARIES` - AI summary generation
- `SAVE_SEARCHES` - Saved searches
- `BULK_VERIFICATION` - Bulk email verification

**Error Responses**:

- `401 Unauthorized`: Not authenticated
  ```json
  {
    "detail": "Not authenticated"
  }
  ```
- `404 Not Found`: User profile not found
  ```json
  {
    "detail": "User profile not found"
  }
  ```
- `500 Internal Server Error`: Failed to retrieve feature usage
  ```json
  {
    "detail": "Failed to retrieve feature usage"
  }
  ```

---

### 2. Track Usage

Track feature usage for the authenticated user. This endpoint is used to record when a user uses a specific feature.

**Endpoint**: `POST /api/v2/usage/track/`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "feature": "DATA_SEARCH",
  "amount": 1
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `feature` | string | Yes | Feature name (must be one of the supported features) |
| `amount` | integer | No | Amount of usage to track (default: 1, min: 1) |

**Response** (200 OK):

```json
{
  "feature": "DATA_SEARCH",
  "used": 151,
  "limit": 1000,
  "success": true
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `feature` | string | Feature name that was tracked |
| `used` | integer | Updated usage count after tracking |
| `limit` | integer | Usage limit for this feature |
| `success` | boolean | Whether tracking was successful |

**Error Responses**:

- `400 Bad Request`: Invalid feature or amount
  ```json
  {
    "detail": "Invalid feature: INVALID_FEATURE"
  }
  ```
  ```json
  {
    "detail": [
      {
        "type": "greater_than_equal",
        "loc": ["body", "amount"],
        "msg": "Input should be greater than or equal to 1",
        "input": 0
      }
    ]
  }
  ```
- `401 Unauthorized`: Not authenticated
  ```json
  {
    "detail": "Not authenticated"
  }
  ```
- `404 Not Found`: User profile not found
  ```json
  {
    "detail": "User profile not found"
  }
  ```
- `500 Internal Server Error`: Failed to track feature usage
  ```json
  {
    "detail": "Failed to track feature usage"
  }
  ```

---

## Feature Types Reference

The following feature types are supported for usage tracking:

| Feature | Description |
|---------|-------------|
| `AI_CHAT` | AI chat interactions |
| `BULK_EXPORT` | Bulk data exports |
| `API_KEYS` | API key usage |
| `TEAM_MANAGEMENT` | Team management operations |
| `EMAIL_FINDER` | Email finder searches |
| `VERIFIER` | Email verification |
| `LINKEDIN` | LinkedIn operations |
| `DATA_SEARCH` | Data search operations |
| `ADVANCED_FILTERS` | Advanced filter usage |
| `AI_SUMMARIES` | AI summary generation |
| `SAVE_SEARCHES` | Saved searches |
| `BULK_VERIFICATION` | Bulk email verification |

---

## Example cURL Requests

### Get Current Usage

```bash
curl -X GET http://localhost:8000/api/v2/usage/current/ \
  -H "Authorization: Bearer <access_token>"
```

### Track Usage

```bash
curl -X POST http://localhost:8000/api/v2/usage/track/ \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "feature": "DATA_SEARCH",
    "amount": 1
  }'
```

### Track Usage (Default Amount)

```bash
curl -X POST http://localhost:8000/api/v2/usage/track/ \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "feature": "EMAIL_FINDER"
  }'
```

---

## Notes

- All endpoints require JWT authentication
- Usage limits are determined by the user's subscription plan
- The `amount` parameter in track usage defaults to 1 if not provided
- Usage tracking increments the used count for the specified feature
- Feature names are case-sensitive and must match exactly
- Usage data is stored per user and per feature
- The current usage endpoint returns all features with their usage statistics
- Usage limits may vary based on subscription tier and plan

