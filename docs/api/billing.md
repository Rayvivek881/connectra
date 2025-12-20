# Billing API

## Overview

The Billing API provides endpoints for managing subscriptions, addon packages, invoices, and billing information. It supports both user-facing operations and administrative management.

**Related Documentation:**
- [Auth API](./auth.md) - For authentication
- [User API](./user.md) - For user profile management
- [Usage API](./usage.md) - For feature usage tracking

**Base URL**: `/api/v2/billing`

**Authentication**: Most endpoints require JWT Bearer token authentication. Public endpoints for viewing plans and addons do not require authentication.

---

## Endpoints

### Public Endpoints (No Authentication Required)

#### 1. Get Subscription Plans

Retrieve all available subscription plans (public endpoint).

**Endpoint**: `GET /api/v2/billing/plans/`

**Authentication**: Not required

**Request Body**: None

**Response** (200 OK):

```json
{
  "plans": [
    {
      "tier": "5k",
      "name": "5K Credits Plan",
      "category": "starter",
      "periods": {
        "monthly": {
          "period": "monthly",
          "credits": 5000,
          "rate_per_credit": 0.01,
          "price": 50.0,
          "savings": {
            "amount": 0.0,
            "percentage": 0
          }
        },
        "yearly": {
          "period": "yearly",
          "credits": 60000,
          "rate_per_credit": 0.008,
          "price": 480.0,
          "savings": {
            "amount": 120.0,
            "percentage": 20
          }
        }
      }
    }
  ]
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `plans` | array | List of subscription plans |
| `plans[].tier` | string | Plan tier identifier (e.g., "5k", "25k", "100k") |
| `plans[].name` | string | Plan display name |
| `plans[].category` | string | Plan category |
| `plans[].periods` | object | Available subscription periods for this plan |
| `plans[].periods[period].period` | string | Period identifier (e.g., "monthly", "yearly") |
| `plans[].periods[period].credits` | integer | Number of credits included |
| `plans[].periods[period].rate_per_credit` | float | Rate per credit |
| `plans[].periods[period].price` | float | Total price |
| `plans[].periods[period].savings` | object | Savings information (optional) |
| `plans[].periods[period].savings.amount` | float | Savings amount |
| `plans[].periods[period].savings.percentage` | integer | Savings percentage |

**Error Responses**:

- `500 Internal Server Error`: Failed to retrieve subscription plans

---

#### 2. Get Addon Packages

Retrieve all available addon packages (public endpoint).

**Endpoint**: `GET /api/v2/billing/addons/`

**Authentication**: Not required

**Request Body**: None

**Response** (200 OK):

```json
{
  "packages": [
    {
      "id": "addon_1k",
      "name": "1K Credits Addon",
      "credits": 1000,
      "rate_per_credit": 0.012,
      "price": 12.0
    },
    {
      "id": "addon_5k",
      "name": "5K Credits Addon",
      "credits": 5000,
      "rate_per_credit": 0.01,
      "price": 50.0
    }
  ]
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `packages` | array | List of addon packages |
| `packages[].id` | string | Package identifier |
| `packages[].name` | string | Package display name |
| `packages[].credits` | integer | Number of credits in package |
| `packages[].rate_per_credit` | float | Rate per credit |
| `packages[].price` | float | Total price |

**Error Responses**:

- `500 Internal Server Error`: Failed to retrieve addon packages

---

### Authenticated Endpoints (JWT Required)

#### 3. Get Billing Info

Get billing information for the current authenticated user.

**Endpoint**: `GET /api/v2/billing/`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Request Body**: None

**Response** (200 OK):

```json
{
  "credits": 5000,
  "credits_used": 250,
  "credits_limit": 5000,
  "subscription_plan": "5k",
  "subscription_period": "monthly",
  "subscription_status": "active",
  "subscription_started_at": "2024-01-01T00:00:00Z",
  "subscription_ends_at": "2024-02-01T00:00:00Z",
  "usage_percentage": 5.0
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `credits` | integer | Current credit balance |
| `credits_used` | integer | Credits used in current period |
| `credits_limit` | integer | Credit limit for current subscription |
| `subscription_plan` | string | Current subscription plan tier |
| `subscription_period` | string | Current subscription period (e.g., "monthly", "yearly") |
| `subscription_status` | string | Subscription status (e.g., "active", "cancelled") |
| `subscription_started_at` | string | Subscription start date (ISO 8601, UTC) |
| `subscription_ends_at` | string | Subscription end date (ISO 8601, UTC) |
| `usage_percentage` | float | Percentage of credits used |

**Error Responses**:

- `401 Unauthorized`: Not authenticated
- `404 Not Found`: User profile not found
- `500 Internal Server Error`: Failed to retrieve billing information

---

#### 4. Subscribe to Plan

Subscribe the current user to a subscription plan.

**Endpoint**: `POST /api/v2/billing/subscribe/`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "tier": "5k",
  "period": "monthly"
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `tier` | string | Yes | Subscription plan tier (e.g., "5k", "25k", "100k") |
| `period` | string | Yes | Subscription period (e.g., "monthly", "yearly") |

**Response** (200 OK):

```json
{
  "message": "Successfully subscribed to plan",
  "subscription_plan": "5k",
  "subscription_period": "monthly",
  "credits": 5000,
  "subscription_ends_at": "2024-02-01T00:00:00Z"
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Success message |
| `subscription_plan` | string | Subscribed plan tier |
| `subscription_period` | string | Subscription period |
| `credits` | integer | Credits added to account |
| `subscription_ends_at` | string | Subscription end date (ISO 8601, UTC) |

**Error Responses**:

- `400 Bad Request`: Invalid tier or period
  ```json
  {
    "detail": "invalid tier: invalid_tier"
  }
  ```
  ```json
  {
    "detail": "invalid period: invalid_period"
  }
  ```
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: User profile not found
- `500 Internal Server Error`: Failed to subscribe to plan

---

#### 5. Purchase Addon

Purchase addon credits for the current user.

**Endpoint**: `POST /api/v2/billing/addon/`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "package_id": "addon_1k"
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `package_id` | string | Yes | Addon package identifier |

**Response** (200 OK):

```json
{
  "message": "Addon credits purchased successfully",
  "package": "addon_1k",
  "credits_added": 1000,
  "total_credits": 6000
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Success message |
| `package` | string | Package identifier |
| `credits_added` | integer | Credits added to account |
| `total_credits` | integer | Total credits after purchase |

**Error Responses**:

- `400 Bad Request`: Invalid package ID
  ```json
  {
    "detail": "invalid package ID: invalid_package"
  }
  ```
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: User profile not found
- `500 Internal Server Error`: Failed to purchase addon credits

---

#### 6. Cancel Subscription

Cancel the current user's subscription.

**Endpoint**: `POST /api/v2/billing/cancel/`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Request Body**: None

**Response** (200 OK):

```json
{
  "message": "Subscription cancelled successfully",
  "subscription_status": "cancelled"
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Success message |
| `subscription_status` | string | Updated subscription status |

**Error Responses**:

- `400 Bad Request`: Subscription already cancelled
  ```json
  {
    "detail": "subscription is already cancelled"
  }
  ```
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: User profile not found
- `500 Internal Server Error`: Failed to cancel subscription

---

#### 7. Get Invoices

Get invoice history for the current user.

**Endpoint**: `GET /api/v2/billing/invoices/`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Query Parameters**:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `limit` | integer | No | 10 | Maximum number of invoices to return (min: 1, max: 100) |
| `offset` | integer | No | 0 | Number of invoices to skip (min: 0) |

**Request Body**: None

**Response** (200 OK):

```json
{
  "invoices": [
    {
      "id": "inv_123456",
      "amount": 50.0,
      "status": "paid",
      "created_at": "2024-01-15T10:30:00Z",
      "description": "5K Credits Plan - Monthly"
    },
    {
      "id": "inv_123457",
      "amount": 12.0,
      "status": "paid",
      "created_at": "2024-01-10T08:15:00Z",
      "description": "1K Credits Addon"
    }
  ],
  "total": 2
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `invoices` | array | List of invoice items |
| `invoices[].id` | string | Invoice identifier |
| `invoices[].amount` | float | Invoice amount |
| `invoices[].status` | string | Invoice status (e.g., "paid", "pending", "failed") |
| `invoices[].created_at` | string | Invoice creation date (ISO 8601, UTC) |
| `invoices[].description` | string | Invoice description (optional) |
| `total` | integer | Total number of invoices |

**Error Responses**:

- `401 Unauthorized`: Not authenticated
- `404 Not Found`: User profile not found
- `500 Internal Server Error`: Failed to retrieve invoices

---

### Admin Endpoints (SuperAdmin Required)

All admin endpoints require SuperAdmin role and are prefixed with `/admin`.

#### 8. Admin Get Subscription Plans

Get all subscription plans for admin management (including inactive).

**Endpoint**: `GET /api/v2/billing/admin/plans/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `include_inactive` | boolean | No | Include inactive plans (default: false) |

**Request Body**: None

**Response**: Same as Get Subscription Plans

**Error Responses**:

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `500 Internal Server Error`: Failed to retrieve subscription plans

---

#### 9. Admin Create Subscription Plan

Create a new subscription plan.

**Endpoint**: `POST /api/v2/billing/admin/plans/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "tier": "10k",
  "name": "10K Credits Plan",
  "category": "professional",
  "is_active": true,
  "periods": [
    {
      "period": "monthly",
      "credits": 10000,
      "rate_per_credit": 0.009,
      "price": 90.0,
      "savings_amount": 0.0,
      "savings_percentage": 0
    },
    {
      "period": "yearly",
      "credits": 120000,
      "rate_per_credit": 0.007,
      "price": 840.0,
      "savings_amount": 240.0,
      "savings_percentage": 22
    }
  ]
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `tier` | string | Yes | Plan tier identifier (must be unique) |
| `name` | string | Yes | Plan display name |
| `category` | string | Yes | Plan category |
| `is_active` | boolean | No | Whether plan is active (default: true) |
| `periods` | array | Yes | List of subscription periods |
| `periods[].period` | string | Yes | Period identifier (e.g., "monthly", "yearly") |
| `periods[].credits` | integer | Yes | Number of credits (min: 1) |
| `periods[].rate_per_credit` | float | Yes | Rate per credit (min: 0) |
| `periods[].price` | float | Yes | Total price (min: 0) |
| `periods[].savings_amount` | float | No | Savings amount |
| `periods[].savings_percentage` | integer | No | Savings percentage |

**Response** (201 Created):

```json
{
  "message": "Subscription plan created successfully",
  "tier": "10k"
}
```

**Error Responses**:

- `400 Bad Request`: Plan with tier already exists
  ```json
  {
    "detail": "plan with tier 10k already exists"
  }
  ```
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `500 Internal Server Error`: Failed to create subscription plan

---

#### 10. Admin Update Subscription Plan

Update an existing subscription plan.

**Endpoint**: `PUT /api/v2/billing/admin/plans/:tier/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Path Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tier` | string | Yes | Plan tier identifier |

**Request Body**:

```json
{
  "name": "Updated Plan Name",
  "category": "premium",
  "is_active": false
}
```

**Request Parameters** (all optional):

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Plan display name |
| `category` | string | Plan category |
| `is_active` | boolean | Whether plan is active |

**Response** (200 OK):

```json
{
  "message": "Subscription plan updated successfully",
  "tier": "10k"
}
```

**Error Responses**:

- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: Plan not found
  ```json
  {
    "detail": "plan with tier 10k not found"
  }
  ```
- `500 Internal Server Error`: Failed to update subscription plan

---

#### 11. Admin Delete Subscription Plan

Delete a subscription plan.

**Endpoint**: `DELETE /api/v2/billing/admin/plans/:tier/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Path Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tier` | string | Yes | Plan tier identifier |

**Request Body**: None

**Response** (200 OK):

```json
{
  "message": "Subscription plan deleted successfully",
  "tier": "10k"
}
```

**Error Responses**:

- `400 Bad Request`: Invalid tier parameter
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: Plan not found
  ```json
  {
    "detail": "plan with tier 10k not found"
  }
  ```
- `500 Internal Server Error`: Failed to delete subscription plan

---

#### 12. Admin Create Subscription Plan Period

Create or update a subscription plan period.

**Endpoint**: `POST /api/v2/billing/admin/plans/:tier/periods/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Path Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tier` | string | Yes | Plan tier identifier |

**Request Body**:

```json
{
  "period": "quarterly",
  "credits": 15000,
  "rate_per_credit": 0.0085,
  "price": 127.5,
  "savings_amount": 7.5,
  "savings_percentage": 6
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `period` | string | Yes | Period identifier |
| `credits` | integer | Yes | Number of credits (min: 1) |
| `rate_per_credit` | float | Yes | Rate per credit (min: 0) |
| `price` | float | Yes | Total price (min: 0) |
| `savings_amount` | float | No | Savings amount |
| `savings_percentage` | integer | No | Savings percentage |

**Response** (200 OK):

```json
{
  "message": "Period created/updated successfully",
  "tier": "10k",
  "period": "quarterly"
}
```

**Error Responses**:

- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: Plan not found
  ```json
  {
    "detail": "plan with tier 10k not found"
  }
  ```
- `500 Internal Server Error`: Failed to create/update period

---

#### 13. Admin Delete Subscription Plan Period

Delete a subscription plan period.

**Endpoint**: `DELETE /api/v2/billing/admin/plans/:tier/periods/:period/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Path Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tier` | string | Yes | Plan tier identifier |
| `period` | string | Yes | Period identifier |

**Request Body**: None

**Response** (200 OK):

```json
{
  "message": "Period deleted successfully",
  "tier": "10k",
  "period": "quarterly"
}
```

**Error Responses**:

- `400 Bad Request`: Missing tier or period parameter
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: Period not found
  ```json
  {
    "detail": "period quarterly not found for plan 10k"
  }
  ```
- `500 Internal Server Error`: Failed to delete period

---

#### 14. Admin Get Addon Packages

Get all addon packages for admin management (including inactive).

**Endpoint**: `GET /api/v2/billing/admin/addons/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `include_inactive` | boolean | No | Include inactive packages (default: false) |

**Request Body**: None

**Response**: Same as Get Addon Packages

**Error Responses**:

- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `500 Internal Server Error`: Failed to retrieve addon packages

---

#### 15. Admin Create Addon Package

Create a new addon package.

**Endpoint**: `POST /api/v2/billing/admin/addons/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body**:

```json
{
  "id": "addon_10k",
  "name": "10K Credits Addon",
  "credits": 10000,
  "rate_per_credit": 0.009,
  "price": 90.0,
  "is_active": true
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Package identifier (must be unique) |
| `name` | string | Yes | Package display name |
| `credits` | integer | Yes | Number of credits (min: 1) |
| `rate_per_credit` | float | Yes | Rate per credit (min: 0) |
| `price` | float | Yes | Total price (min: 0) |
| `is_active` | boolean | No | Whether package is active (default: true) |

**Response** (201 Created):

```json
{
  "message": "Addon package created successfully",
  "id": "addon_10k"
}
```

**Error Responses**:

- `400 Bad Request`: Package with ID already exists
  ```json
  {
    "detail": "package with id addon_10k already exists"
  }
  ```
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `500 Internal Server Error`: Failed to create addon package

---

#### 16. Admin Update Addon Package

Update an existing addon package.

**Endpoint**: `PUT /api/v2/billing/admin/addons/:package_id/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Path Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `package_id` | string | Yes | Package identifier |

**Request Body** (all fields optional):

```json
{
  "name": "Updated Package Name",
  "credits": 12000,
  "rate_per_credit": 0.008,
  "price": 96.0,
  "is_active": false
}
```

**Request Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Package display name |
| `credits` | integer | Number of credits (min: 1) |
| `rate_per_credit` | float | Rate per credit (min: 0) |
| `price` | float | Total price (min: 0) |
| `is_active` | boolean | Whether package is active |

**Response** (200 OK):

```json
{
  "message": "Addon package updated successfully",
  "id": "addon_10k"
}
```

**Error Responses**:

- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: Package not found
  ```json
  {
    "detail": "package with id addon_10k not found"
  }
  ```
- `500 Internal Server Error`: Failed to update addon package

---

#### 17. Admin Delete Addon Package

Delete an addon package.

**Endpoint**: `DELETE /api/v2/billing/admin/addons/:package_id/`

**Authentication**: Required (JWT Bearer Token, SuperAdmin role)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Path Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `package_id` | string | Yes | Package identifier |

**Request Body**: None

**Response** (200 OK):

```json
{
  "message": "Addon package deleted successfully",
  "id": "addon_10k"
}
```

**Error Responses**:

- `400 Bad Request`: Invalid package_id parameter
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: Package not found
  ```json
  {
    "detail": "package with id addon_10k not found"
  }
  ```
- `500 Internal Server Error`: Failed to delete addon package

---

## Example cURL Requests

### Get Subscription Plans (Public)

```bash
curl -X GET http://localhost:8000/api/v2/billing/plans/
```

### Get Billing Info

```bash
curl -X GET http://localhost:8000/api/v2/billing/ \
  -H "Authorization: Bearer <access_token>"
```

### Subscribe to Plan

```bash
curl -X POST http://localhost:8000/api/v2/billing/subscribe/ \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "tier": "5k",
    "period": "monthly"
  }'
```

### Purchase Addon

```bash
curl -X POST http://localhost:8000/api/v2/billing/addon/ \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "package_id": "addon_1k"
  }'
```

### Cancel Subscription

```bash
curl -X POST http://localhost:8000/api/v2/billing/cancel/ \
  -H "Authorization: Bearer <access_token>"
```

### Get Invoices

```bash
curl -X GET "http://localhost:8000/api/v2/billing/invoices/?limit=20&offset=0" \
  -H "Authorization: Bearer <access_token>"
```

### Admin Create Subscription Plan

```bash
curl -X POST http://localhost:8000/api/v2/billing/admin/plans/ \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "tier": "10k",
    "name": "10K Credits Plan",
    "category": "professional",
    "is_active": true,
    "periods": [
      {
        "period": "monthly",
        "credits": 10000,
        "rate_per_credit": 0.009,
        "price": 90.0
      }
    ]
  }'
```

---

## Notes

- All prices are in the base currency (e.g., USD)
- Credits are added immediately upon subscription or addon purchase
- Subscription periods are calculated from the subscription start date
- Cancelled subscriptions remain active until the end of the current billing period
- Admin endpoints require SuperAdmin role
- Public endpoints (plans and addons) do not require authentication
- All timestamps are in ISO 8601 format (UTC)
- Invoice pagination uses limit/offset with default limit of 10 and maximum of 100

