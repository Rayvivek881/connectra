# Authentication API

## Overview

The Authentication API provides endpoints for user registration, login, logout, session management, and token refresh. All authentication endpoints use JWT (JSON Web Tokens) for secure access.

**Related Documentation:**
- [User API](./user.md) - For user profile management
- [Billing API](./billing.md) - For subscription and billing
- [Usage API](./usage.md) - For feature usage tracking

**Base URL**: `/api/v2/auth`

---

## Endpoints

### 1. Register User

Register a new user account and receive access tokens.

**Endpoint**: `POST /api/v2/auth/register`

**Authentication**: Not required

**Request Body**:

```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "securePassword123",
  "geolocation": {
    "ip": "205.254.184.116",
    "continent": "Asia",
    "continent_code": "AS",
    "country": "India",
    "country_code": "IN",
    "region": "KA",
    "region_name": "Karnataka",
    "city": "Bengaluru",
    "district": "",
    "zip": "",
    "lat": 12.9715,
    "lon": 77.5945,
    "timezone": "Asia/Kolkata",
    "offset": 19800,
    "currency": "INR",
    "isp": "Excitel Broadband Pvt Ltd",
    "org": "Excitel Broadband Pvt Ltd",
    "asname": "",
    "reverse": "",
    "device": "Mozilla/5.0...",
    "proxy": false,
    "hosting": false
  }
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | User's full name (max 255 characters) |
| `email` | string | Yes | User's email address (must be unique, valid email format) |
| `password` | string | Yes | User's password (min 8, max 72 characters, will be hashed) |
| `geolocation` | object | No | IP geolocation data from frontend (all fields optional) |

**Response** (201 Created):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com"
  },
  "message": "Registration successful! Please check your email to verify your account."
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `access_token` | string | JWT access token (valid for 30 minutes) |
| `refresh_token` | string | JWT refresh token (valid for 7 days) |
| `user` | object | User information |
| `user.uuid` | string | User's unique identifier |
| `user.email` | string | User's email address |
| `message` | string | Success message |

**Error Responses**:

- `400 Bad Request`: Email already exists or password validation failed
  ```json
  {
    "email": ["Email already exists"]
  }
  ```
  ```json
  {
    "password": ["This password is too short. It must contain at least 8 characters."]
  }
  ```
- `422 Unprocessable Entity`: Invalid request data format or missing required fields

---

### 2. Login

Authenticate a user and receive access tokens.

**Endpoint**: `POST /api/v2/auth/login`

**Authentication**: Not required

**Request Body**:

```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "geolocation": {
    "ip": "205.254.184.116",
    "continent": "Asia",
    "continent_code": "AS",
    "country": "India",
    "country_code": "IN",
    "region": "KA",
    "region_name": "Karnataka",
    "city": "Bengaluru",
    "district": "",
    "zip": "",
    "lat": 12.9715,
    "lon": 77.5945,
    "timezone": "Asia/Kolkata",
    "offset": 19800,
    "currency": "INR",
    "isp": "Excitel Broadband Pvt Ltd",
    "org": "Excitel Broadband Pvt Ltd",
    "asname": "",
    "reverse": "",
    "device": "Mozilla/5.0...",
    "proxy": false,
    "hosting": false
  }
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User's email address |
| `password` | string | Yes | User's password |
| `geolocation` | object | No | IP geolocation data from frontend (all fields optional) |

**Response** (200 OK):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com"
  }
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `access_token` | string | JWT access token (valid for 30 minutes) |
| `refresh_token` | string | JWT refresh token (valid for 7 days) |
| `user` | object | User information |
| `user.uuid` | string | User's unique identifier |
| `user.email` | string | User's email address |

**Token Claims**:

The JWT tokens contain the following claims:

- `sub`: User UUID
- `email`: User's email address
- `role`: User's role
- `exp`: Token expiration timestamp

**Error Responses**:

- `400 Bad Request`: Invalid credentials or account disabled
  ```json
  {
    "detail": "Invalid email or password"
  }
  ```
  ```json
  {
    "detail": {
      "non_field_errors": ["User account is disabled"]
    }
  }
  ```
- `422 Unprocessable Entity`: Missing required fields or invalid data format

**Usage**:

Include the access token in subsequent API requests using the `Authorization` header:

```http
Authorization: Bearer <access_token>
```

---

### 3. Logout

Logout the current user and invalidate refresh token.

**Endpoint**: `POST /api/v2/auth/logout`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Request Body**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `refresh_token` | string | No | Refresh token to blacklist. If not provided, logout still succeeds but token may not be invalidated. |

**Response** (200 OK):

```json
{
  "message": "Logout successful"
}
```

**Error Responses**:

- `401 Unauthorized`: Authentication credentials were not provided
  ```json
  {
    "detail": "Authentication credentials were not provided."
  }
  ```

**Note**: 
- The refresh token (if provided) is added to a blacklist to prevent reuse
- Access tokens cannot be blacklisted but will expire naturally
- Logout succeeds even if refresh token is invalid or not provided

---

### 4. Get Current Session

Get the current authenticated user's session information.

**Endpoint**: `GET /api/v2/auth/session`

**Authentication**: Required (JWT Bearer Token)

**Headers**:

```http
Authorization: Bearer <access_token>
```

**Request Body**: None

**Response** (200 OK):

```json
{
  "user": {
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "last_sign_in_at": "2024-01-15T12:00:00Z"
  }
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `user` | object | User session information |
| `user.uuid` | string | User's unique identifier |
| `user.email` | string | User's email address |
| `user.last_sign_in_at` | string | Last sign-in timestamp (ISO 8601, UTC). Will be `null` if user has never logged in. |

**Error Responses**:

- `401 Unauthorized`: Invalid or expired token
  ```json
  {
    "detail": "Authentication credentials were not provided."
  }
  ```

**Note**: This endpoint is useful for checking token validity.

---

### 5. Refresh Access Token

Refresh an expired access token using a refresh token.

**Endpoint**: `POST /api/v2/auth/refresh`

**Authentication**: Not required (uses refresh token instead)

**Request Body**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Request Parameters**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `refresh_token` | string | Yes | Valid refresh token |

**Response** (200 OK):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response Parameters**:

| Field | Type | Description |
|-------|------|-------------|
| `access_token` | string | New JWT access token (valid for 30 minutes) |
| `refresh_token` | string | New JWT refresh token (valid for 7 days) |

**Error Responses**:

- `400 Bad Request`: Invalid refresh token
  ```json
  {
    "detail": "Invalid refresh token"
  }
  ```
  ```json
  {
    "detail": "Token is invalid or expired"
  }
  ```

**Note**: 
- A new refresh token is also returned (token rotation)
- The old refresh token is not automatically invalidated - both tokens remain valid until they expire
- This endpoint does not require authentication (uses refresh token instead)
- Token rotation provides better security by limiting the lifetime of refresh tokens

---

## Authentication Flow

1. **Register**: Create a new user account and receive tokens
2. **Login**: Authenticate and receive access/refresh tokens
3. **Use Token**: Include the access token in the `Authorization` header for protected endpoints
4. **Refresh Token**: Use refresh token to get new access token when it expires
5. **Logout**: Invalidate the refresh token when done

## User Roles

The system supports the following user roles:

- `SuperAdmin`: Full system access
- `Admin`: Administrative access
- `ProUser`: Premium user access
- `FreeUser`: Basic user access (default for new registrations)

**Note**: New registrations automatically receive:
- `role`: `FreeUser` (default role)
- `credits`: `50` (initial credits for free users)
- `subscription_plan`: `free`
- `subscription_status`: `active`

## Security Notes

- Passwords are hashed using bcrypt before storage
- Password length: minimum 8 characters, maximum 72 characters (bcrypt limitation)
- Access tokens expire after 30 minutes (configurable via `ACCESS_TOKEN_EXPIRE_MINUTES`)
- Refresh tokens expire after 7 days (configurable via `REFRESH_TOKEN_EXPIRE_DAYS`)
- Blacklisted refresh tokens cannot be reused
- All protected endpoints require valid JWT authentication
- Tokens are validated on every request
- Token refresh implements token rotation (new tokens issued, old tokens remain valid until expiration)

---

## Example cURL Requests

### Register

```bash
curl -X POST http://localhost:8000/api/v2/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "user@example.com",
    "password": "securePassword123",
    "geolocation": {
      "ip": "205.254.184.116",
      "country": "India",
      "city": "Bengaluru",
      "lat": 12.9715,
      "lon": 77.5945
    }
  }'
```

### Login

```bash
curl -X POST http://localhost:8000/api/v2/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securePassword123"
  }'
```

### Logout

```bash
curl -X POST http://localhost:8000/api/v2/auth/logout \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```

### Get Session

```bash
curl -X GET http://localhost:8000/api/v2/auth/session \
  -H "Authorization: Bearer <access_token>"
```

### Refresh Token

```bash
curl -X POST http://localhost:8000/api/v2/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```
