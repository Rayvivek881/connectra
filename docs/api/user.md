# User API Documentation

Complete API documentation for user authentication and profile management endpoints.

**Related Documentation:**
- [Auth API](./auth.md) - For authentication endpoints
- [Billing API](./billing.md) - For subscription and billing management
- [Usage API](./usage.md) - For feature usage tracking
- [Contacts API](./contacts.md) - For contact management endpoints
- [Companies API](./companies.md) - For company management endpoints

## Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Role-Based Access Control](#role-based-access-control)
- [CORS Testing](#cors-testing)
- [Authentication Endpoints](#authentication-endpoints)
  - [POST /api/v2/auth/register/](#post-apiv2authregister---user-registration)
  - [POST /api/v2/auth/login/](#post-apiv2authlogin---user-login)
  - [POST /api/v2/auth/logout/](#post-apiv2authlogout---user-logout)
  - [GET /api/v2/auth/session/](#get-apiv2authsession---get-current-session)
  - [POST /api/v2/auth/refresh/](#post-apiv2authrefresh---refresh-access-token)
- [User Profile Endpoints](#user-profile-endpoints)
  - [GET /api/v2/users/profile/](#get-apiv2usersprofile---get-current-user-profile)
  - [PUT /api/v2/users/profile/](#put-apiv2usersprofile---update-current-user-profile)
  - [POST /api/v2/users/profile/avatar/](#post-apiv2usersprofileavatar---upload-user-avatar)
  - [POST /api/v2/users/promote-to-admin/](#post-apiv2userspromote-to-admin---promote-user-to-admin)
  - [POST /api/v2/users/promote-to-super-admin/](#post-apiv2userspromote-to-super-admin---promote-user-to-super-admin)
- [Super Admin Endpoints](#super-admin-endpoints)
  - [GET /api/v2/users/](#get-apiv2users---list-all-users)
  - [PUT /api/v2/users/{user_id}/role/](#put-apiv2usersuser_idrole---update-user-role)
  - [PUT /api/v2/users/{user_id}/credits/](#put-apiv2usersuser_idcredits---update-user-credits)
  - [DELETE /api/v2/users/{user_id}/](#delete-apiv2usersuser_id---delete-user)
  - [GET /api/v2/users/stats/](#get-apiv2usersstats---get-user-statistics)
  - [GET /api/v2/users/history/](#get-apiv2usershistory---get-user-history)
- [Error Responses](#error-responses)
- [Notes](#notes)

---

## Base URL

```txt
http://localhost:8000
```

**API Version:** `/api/v2/`
- Authentication endpoints: `/api/v2/auth/`
- User profile endpoints: `/api/v2/users/`

## Authentication

Most endpoints require JWT authentication via the `Authorization` header:

```txt
Authorization: Bearer <access_token>
```

Tokens are obtained through the login or register endpoints.

---

## Role-Based Access Control

The API implements a comprehensive role-based access control (RBAC) system with four distinct user roles:

### User Roles

1. **SuperAdmin** (`SuperAdmin`)
   - Full control over all users, UI, and plan details
   - Can manage user roles, credits, and delete users
   - Can view user statistics
   - Has access to all features
   - **Unlimited credits** (no credit deduction for any operations)

2. **Admin** (`Admin`)
   - Full control over all UI pages
   - Cannot manage users (no user management capabilities)
   - Can view user statistics
   - Has access to all features
   - **Unlimited credits** (no credit deduction for any operations)

3. **ProUser** (`ProUser`)
   - Full CRUD access to contacts and companies
   - Can purchase subscription plans
   - Has access to all UI features
   - Can modify (update/delete) resources
   - Credits are deducted for operations (1 credit per search, 1 credit per item exported)

4. **FreeUser** (`FreeUser`)
   - Default role for new registrations
   - Receives 50 initial credits upon registration
   - Can create and read contacts and companies
   - Cannot update or delete resources (read-only for modifications)
   - Has access to contact page, company page, LinkedIn, email, and AI assistants
   - Credits are deducted for operations (1 credit per search, 1 credit per item exported)

### Role Permissions Summary

| Action | FreeUser | ProUser | Admin | SuperAdmin |
|--------|----------|---------|-------|------------|
| Create contacts/companies | ✅ | ✅ | ✅ | ✅ |
| Read contacts/companies | ✅ | ✅ | ✅ | ✅ |
| Update contacts/companies | ❌ | ✅ | ✅ | ✅ |
| Delete contacts/companies | ❌ | ✅ | ✅ | ✅ |
| Purchase plans | ❌ | ✅ | ✅ | ✅ |
| Manage users | ❌ | ❌ | ❌ | ✅ |
| View user statistics | ❌ | ❌ | ✅ | ✅ |
| All UI features | ❌ | ✅ | ✅ | ✅ |

### Role Assignment

- **New Registrations**: Automatically assigned `FreeUser` role with 50 initial credits
- **Role Changes**: Only SuperAdmin can change user roles via `PUT /api/v2/users/{user_id}/role/`
- **Self-Promotion**: Users can self-promote to Admin via `POST /api/v2/users/promote-to-admin/` (not recommended for production)

### Error Responses for Role Restrictions

When a user attempts to access an endpoint without the required role, they will receive:

**Error (403 Forbidden):**

```json
{
  "detail": "You do not have permission to perform this action. [Role] role required."
}
```

---

## CORS Testing

All endpoints support CORS (Cross-Origin Resource Sharing) for browser-based requests. For testing CORS headers, you can include an optional `Origin` header in your requests:

**Optional Header:**

- `Origin: http://localhost:3000` (or your frontend origin)

**Expected CORS Response Headers:**

- `Access-Control-Allow-Origin: http://localhost:3000` (matches the Origin header)
- `Access-Control-Allow-Credentials: true`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH`
- `Access-Control-Allow-Headers: *`
- `Access-Control-Max-Age: 3600`

**Note:** The Origin header is optional and only needed when testing CORS behavior. The API automatically handles CORS preflight (OPTIONS) requests.

---

## Authentication Endpoints

### POST /api/v2/auth/register/ - User Registration

Register a new user account and receive access tokens.

**Headers:**

- `Content-Type: application/json`

**Request Body:**

```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123",
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

**Field Requirements:**

- `name` (string, required): User's full name (max 255 characters)
- `email` (string, required): Valid email address (must be unique)
- `password` (string, required): Password with minimum 8 characters and maximum 72 characters (bcrypt limitation)
- `geolocation` (object, optional): IP geolocation data from frontend. All fields within geolocation are optional.

**Response:**

**Success (201 Created):**

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

**Error (400 Bad Request) - Email Already Exists:**

```json
{
  "email": ["Email already exists"]
}
```

**Error (400 Bad Request) - Invalid Password:**

```json
{
  "password": [
    "This password is too short. It must contain at least 8 characters."
  ]
}
```

**Status Codes:**

- `201 Created`: Registration successful
- `400 Bad Request`: Email already exists or password validation failed
- `422 Unprocessable Entity`: Invalid request data format or missing required fields

**Notes:**

- A user profile is automatically created upon registration with default values:
  - `role`: `FreeUser` (default role for new users)
  - `credits`: `50` (initial credits for free users)
  - `subscription_plan`: `free`
  - `subscription_status`: `active`
  - `notifications`: `{"weeklyReports": true, "newLeadAlerts": true}`
- The email is used as the username
- Tokens are immediately returned for automatic login after registration
- The `geolocation` field is optional. If provided, it will be stored in the user history table for audit purposes.

---

### POST /api/v2/auth/login/ - User Login

Authenticate a user and receive access tokens.

**Headers:**

- `Content-Type: application/json`

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123",
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

**Field Requirements:**

- `email` (string, required): User's email address
- `password` (string, required): User's password
- `geolocation` (object, optional): IP geolocation data from frontend

**Response:**

**Success (200 OK):**

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

**Error (400 Bad Request) - Invalid Credentials:**

```json
{
  "detail": "Invalid email or password"
}
```

**Error (400 Bad Request) - Account Disabled:**

```json
{
  "detail": {
    "non_field_errors": ["User account is disabled"]
  }
}
```

**Status Codes:**

- `200 OK`: Login successful
- `400 Bad Request`: Invalid credentials or account disabled
- `422 Unprocessable Entity`: Missing required fields or invalid data format

**Notes:**

- The user's `last_sign_in_at` timestamp is updated upon successful login
- Tokens are generated using JWT (JSON Web Tokens) with HS256 algorithm
- Access tokens expire after 30 minutes (configurable via `ACCESS_TOKEN_EXPIRE_MINUTES`)
- Refresh tokens expire after 7 days (configurable via `REFRESH_TOKEN_EXPIRE_DAYS`)
- The `geolocation` field is optional. If provided, it will be stored in the user history table for audit purposes.

---

### POST /api/v2/auth/logout/ - User Logout

Logout the current user and invalidate refresh token.

**Headers:**

- `Authorization: Bearer <access_token>` (required)

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Field Requirements:**

- `refresh_token` (string, optional): Refresh token to blacklist. If not provided, logout still succeeds but token may not be invalidated.

**Response:**

**Success (200 OK):**

```json
{
  "message": "Logout successful"
}
```

**Error (401 Unauthorized):**

```json
{
  "detail": "Authentication credentials were not provided."
}
```

**Status Codes:**

- `200 OK`: Logout successful
- `401 Unauthorized`: Authentication required or invalid token

**Notes:**

- The refresh token (if provided) is added to a blacklist to prevent reuse
- Access tokens cannot be blacklisted but will expire naturally
- Logout succeeds even if refresh token is invalid or not provided

---

### GET /api/v2/auth/session/ - Get Current Session

Get the current authenticated user's session information.

**Headers:**

- `Authorization: Bearer <access_token>` (required)

**Response:**

**Success (200 OK):**

```json
{
  "user": {
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "last_sign_in_at": "2024-01-15T12:00:00Z"
  }
}
```

**Error (401 Unauthorized):**

```json
{
  "detail": "Authentication credentials were not provided."
}
```

**Status Codes:**

- `200 OK`: Session valid
- `401 Unauthorized`: Invalid or expired token

**Notes:**

- `last_sign_in_at` will be `null` if the user has never logged in
- This endpoint is useful for checking token validity
- The timestamp is in ISO 8601 format (UTC)

---

### POST /api/v2/auth/refresh/ - Refresh Access Token

Refresh an expired access token using a refresh token.

**Headers:**

- `Content-Type: application/json`

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Field Requirements:**

- `refresh_token` (string, required): Valid refresh token

**Response:**

**Success (200 OK):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error (400 Bad Request) - Invalid Refresh Token:**

```json
{
  "detail": "Invalid refresh token"
}
```

**Error (400 Bad Request) - Token Error:**

```json
{
  "detail": "Token is invalid or expired"
}
```

**Status Codes:**

- `200 OK`: Token refreshed successfully
- `400 Bad Request`: Invalid refresh token or missing field

**Notes:**

- A new refresh token is also returned (token rotation)
- The old refresh token is not automatically invalidated - both tokens remain valid until they expire
- This endpoint does not require authentication (uses refresh token instead)
- Token rotation provides better security by limiting the lifetime of refresh tokens

---

## User Profile Endpoints

### GET /api/v2/users/profile/ - Get Current User Profile

Get the profile information for the currently authenticated user.

**Headers:**

- `Authorization: Bearer <access_token>` (required)

**Response:**

**Success (200 OK):**

```json
{
  "uuid": "123e4567-e89b-12d3-a456-426614174000",
  "name": "John Doe",
  "email": "user@example.com",
  "role": "Member",
  "avatar_url": "https://picsum.photos/seed/123/40/40",
  "is_active": true,
  "job_title": "Software Engineer",
  "bio": "Passionate developer",
  "timezone": "America/New_York",
  "notifications": {
    "weeklyReports": true,
    "newLeadAlerts": true
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Error (401 Unauthorized):**

```json
{
  "detail": "Authentication credentials were not provided."
}
```

**Status Codes:**

- `200 OK`: Profile retrieved successfully
- `401 Unauthorized`: Authentication required

**Notes:**

- If a profile doesn't exist, it will be automatically created with default values
- The profile is linked to the authenticated user via a one-to-one relationship
- All timestamps are in ISO 8601 format (UTC)

---

### PUT /api/v2/users/profile/ - Update Current User Profile

Update the profile information for the currently authenticated user. All fields are optional - only provided fields will be updated.

**Headers:**

- `Authorization: Bearer <access_token>` (required)
- `Content-Type: application/json`

**Request Body:**

```json
{
  "name": "John Doe Updated",
  "job_title": "Senior Software Engineer",
  "bio": "Updated bio",
  "timezone": "America/Los_Angeles",
  "avatar_url": "https://picsum.photos/seed/123/40/40",
  "notifications": {
    "weeklyReports": false,
    "newLeadAlerts": true
  },
  "role": "Admin"
}
```

**Field Requirements:**

All fields are optional:

- `name` (string, optional): User's full name (max 255 characters)
- `job_title` (string, optional): User's job title (max 255 characters)
- `bio` (string, optional): User's biography (text field)
- `timezone` (string, optional): User's timezone (max 100 characters)
- `avatar_url` (string, optional): URL to user's avatar image
- `notifications` (object, optional): User notification preferences (merged with existing preferences)
- `role` (string, optional): User's role (max 50 characters)

**Response:**

**Success (200 OK):**

```json
{
  "uuid": "123e4567-e89b-12d3-a456-426614174000",
  "name": "John Doe Updated",
  "email": "user@example.com",
  "role": "Admin",
  "avatar_url": "https://picsum.photos/seed/123/40/40",
  "is_active": true,
  "job_title": "Senior Software Engineer",
  "bio": "Updated bio",
  "timezone": "America/Los_Angeles",
  "notifications": {
    "weeklyReports": false,
    "newLeadAlerts": true
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

**Status Codes:**

- `200 OK`: Profile updated successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required

**Notes:**

- This is a partial update (PATCH-like behavior) - only provided fields are updated
- The `notifications` field is merged with existing preferences, not replaced
- If a profile doesn't exist, it will be automatically created
- The `email` field cannot be updated through this endpoint (it's read-only)
- The `updated_at` timestamp is automatically updated

---

### POST /api/v2/users/profile/avatar/ - Upload User Avatar

Upload an avatar image file for the currently authenticated user. The image will be stored in the media directory or S3 and the user's `avatar_url` will be updated automatically.

**Headers:**

- `Authorization: Bearer <access_token>` (required)
- `Content-Type: multipart/form-data`

**Request Body:**

Form data with a file field named `avatar`:

```txt
avatar: [image file]
```

**File Requirements:**

- **File Types**: JPEG, PNG, GIF, or WebP
- **Maximum Size**: 5MB (5,242,880 bytes)
- **Validation**: Both file extension and file content (magic bytes) are validated
- **Allowed Extensions**: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`

**Response:**

**Success (200 OK):**

```json
{
  "avatar_url": "http://localhost:8000/media/avatars/123e4567-e89b-12d3-a456-426614174000_20240115T120000123456.jpg",
  "profile": {
    "uuid": "123e4567-e89b-12d3-a456-426614174000",
    "name": "John Doe",
    "email": "user@example.com",
    "role": "Member",
    "avatar_url": "http://localhost:8000/media/avatars/123e4567-e89b-12d3-a456-426614174000_20240115T120000123456.jpg",
    "is_active": true,
    "job_title": "Software Engineer",
    "bio": "Passionate developer",
    "timezone": "America/New_York",
    "notifications": {
      "weeklyReports": true,
      "newLeadAlerts": true
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-15T12:00:00Z"
  },
  "message": "Avatar uploaded successfully"
}
```

**Error (400 Bad Request) - File Too Large:**

```json
{
  "avatar": [
    "Image file too large. Maximum size is 5.0MB"
  ]
}
```

**Error (400 Bad Request) - Invalid File Type:**

```json
{
  "avatar": [
    "Invalid file type. Allowed types: .jpg, .jpeg, .png, .gif, .webp"
  ]
}
```

**Status Codes:**

- `200 OK`: Avatar uploaded successfully
- `400 Bad Request`: Invalid file (wrong type, too large, or not a valid image)
- `401 Unauthorized`: Authentication required
- `500 Internal Server Error`: Server error while saving file

**Notes:**

- The old avatar file (if it exists and is a local file) will be automatically deleted when a new avatar is uploaded
- External avatar URLs (not stored locally, starting with `http://` or `https://`) will not be deleted
- The filename format is: `{user_id}_{timestamp}.{extension}` where timestamp is in format `YYYYMMDDTHHMMSSffffff` (UTC)
- Files are stored in the `uploads/avatars/` directory (configurable via `UPLOAD_DIR` setting) or S3 if configured
- File validation checks both extension and magic bytes (file signature) to ensure it's actually an image
- If S3 is configured, files are uploaded to S3 with fallback to local storage
- If a profile doesn't exist, it will be automatically created with default values

---

### POST /api/v2/users/promote-to-admin/ - Promote User to Admin

Promote the currently authenticated user to admin role. This endpoint allows authenticated users to change their role to "Admin". The operation is logged for audit purposes.

**Headers:**

- `Authorization: Bearer <access_token>` (required)

**Request Body:**

No request body required. The endpoint uses the authenticated user from the Bearer token.

**Response:**

**Success (200 OK):**

```json
{
  "uuid": "123e4567-e89b-12d3-a456-426614174000",
  "name": "John Doe",
  "email": "user@example.com",
  "role": "Admin",
  "avatar_url": null,
  "is_active": true,
  "job_title": "Software Engineer",
  "bio": "Passionate developer",
  "timezone": "America/New_York",
  "notifications": {
    "weeklyReports": true,
    "newLeadAlerts": true
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

**Error (400 Bad Request) - User Account Disabled:**

```json
{
  "non_field_errors": ["User account is disabled"]
}
```

**Status Codes:**

- `200 OK`: User promoted to admin successfully
- `400 Bad Request`: User account is disabled
- `401 Unauthorized`: Authentication required
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error while promoting user

**Notes:**

- This endpoint allows authenticated users to self-promote to admin role
- The operation is logged for audit purposes
- If a profile doesn't exist, it will be automatically created with default values before promotion
- The `role` field in the profile is updated to "Admin"
- The `updated_at` timestamp is automatically updated
- This is a self-service endpoint with no additional security checks (consider adding rate limiting or admin approval workflow in production)

---

### POST /api/v2/users/promote-to-super-admin/ - Promote User to Super Admin

Promote a user to super admin role (Super Admin only). This endpoint allows super admins to promote any user to "SuperAdmin" role. The operation is logged for audit purposes.

**Headers:**

- `Authorization: Bearer <access_token>` (required, SuperAdmin role)

**Query Parameters:**

- `user_id` (string, UUID, required): User ID to promote to super admin

**Request Body:**

No request body required. The target user UUID is specified via the `user_id` query parameter.

**Response:**

**Success (200 OK):**

```json
{
  "uuid": "123e4567-e89b-12d3-a456-426614174000",
  "name": "John Doe",
  "email": "user@example.com",
  "role": "SuperAdmin",
  "avatar_url": null,
  "is_active": true,
  "job_title": "Software Engineer",
  "bio": "Passionate developer",
  "timezone": "America/New_York",
  "notifications": {
    "weeklyReports": true,
    "newLeadAlerts": true
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

**Error (403 Forbidden) - Not Super Admin:**

```json
{
  "detail": "You do not have permission to perform this action. SuperAdmin role required."
}
```

**Status Codes:**

- `200 OK`: User promoted to super admin successfully
- `400 Bad Request`: User account is disabled
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error while promoting user

**Notes:**

- Only SuperAdmin can promote users to SuperAdmin role
- The operation is logged for audit purposes
- If a profile doesn't exist, it will be automatically created with default values before promotion
- The `role` field in the profile is updated to "SuperAdmin"
- The `updated_at` timestamp is automatically updated
- Requires `user_id` query parameter to specify the target user

---

## Super Admin Endpoints

All Super Admin endpoints require the `SuperAdmin` role. These endpoints allow full user management capabilities.

### GET /api/v2/users/ - List All Users

List all users in the system with their profiles. This endpoint is restricted to Super Admin only.

**Headers:**

- `Authorization: Bearer <access_token>` (required, SuperAdmin role)

**Query Parameters:**

- `limit` (integer, optional, default: 100, min: 1, max: 1000): Maximum number of users to return
- `offset` (integer, optional, default: 0, min: 0): Number of users to skip (for pagination)

**Response:**

**Success (200 OK):**

```json
{
  "users": [
    {
      "uuid": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "FreeUser",
      "is_active": true,
      "credits": 50,
      "subscription_plan": "free",
      "subscription_status": "active",
      "created_at": "2024-01-01T00:00:00Z",
      "last_sign_in_at": "2024-01-15T12:00:00Z"
    }
  ],
  "total": 150
}
```

**Error (403 Forbidden) - Not Super Admin:**

```json
{
  "detail": "You do not have permission to perform this action. SuperAdmin role required."
}
```

**Status Codes:**

- `200 OK`: Users retrieved successfully
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: SuperAdmin role required
- `500 Internal Server Error`: Failed to list users

**Notes:**

- Returns paginated list of all users with their profile information
- Includes role, credits, subscription plan, and status for each user
- Only SuperAdmin can access this endpoint

---

### PUT /api/v2/users/{user_id}/role/ - Update User Role

Update a user's role. This endpoint is restricted to Super Admin only.

**Headers:**

- `Authorization: Bearer <access_token>` (required, SuperAdmin role)
- `Content-Type: application/json`

**Path Parameters:**

- `user_id` (string, UUID, required): User ID to update

**Request Body:**

```json
{
  "role": "ProUser"
}
```

**Field Requirements:**

- `role` (string, required): New role for the user. Valid values: `SuperAdmin`, `Admin`, `FreeUser`, `ProUser`

**Response:**

**Success (200 OK):**

Returns a `ProfileResponse` object with the updated role.

**Error (400 Bad Request) - Invalid Role:**

```json
{
  "detail": "Invalid role: invalid_role. Valid roles: SuperAdmin, Admin, FreeUser, ProUser"
}
```

**Error (403 Forbidden) - Not Super Admin:**

```json
{
  "detail": "You do not have permission to perform this action. SuperAdmin role required."
}
```

**Status Codes:**

- `200 OK`: Role updated successfully
- `400 Bad Request`: Invalid role value
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: User not found
- `500 Internal Server Error`: Failed to update user role

**Notes:**

- Only SuperAdmin can change user roles
- Valid roles: `SuperAdmin`, `Admin`, `FreeUser`, `ProUser`
- The operation is logged for audit purposes

---

### PUT /api/v2/users/{user_id}/credits/ - Update User Credits

Update a user's credit balance. This endpoint is restricted to Super Admin only.

**Headers:**

- `Authorization: Bearer <access_token>` (required, SuperAdmin role)
- `Content-Type: application/json`

**Path Parameters:**

- `user_id` (string, UUID, required): User ID to update

**Request Body:**

```json
{
  "credits": 1000
}
```

**Field Requirements:**

- `credits` (integer, required, min: 0): New credit balance for the user

**Response:**

**Success (200 OK):**

Returns a `ProfileResponse` object with the updated credits.

**Error (403 Forbidden) - Not Super Admin:**

```json
{
  "detail": "You do not have permission to perform this action. SuperAdmin role required."
}
```

**Status Codes:**

- `200 OK`: Credits updated successfully
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: User not found
- `500 Internal Server Error`: Failed to update user credits

**Notes:**

- Only SuperAdmin can modify user credits
- Credits must be a non-negative integer
- Useful for manual credit adjustments or promotional credits

---

### DELETE /api/v2/users/{user_id}/ - Delete User

Delete a user and their profile. This endpoint is restricted to Super Admin only.

**Headers:**

- `Authorization: Bearer <access_token>` (required, SuperAdmin role)

**Path Parameters:**

- `user_id` (string, UUID, required): User ID to delete

**Response:**

**Success (204 No Content):**

No response body.

**Error (400 Bad Request) - Cannot Delete Self:**

```json
{
  "detail": "Cannot delete your own account"
}
```

**Error (403 Forbidden) - Not Super Admin:**

```json
{
  "detail": "You do not have permission to perform this action. SuperAdmin role required."
}
```

**Status Codes:**

- `204 No Content`: User deleted successfully
- `400 Bad Request`: Cannot delete own account
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: SuperAdmin role required
- `404 Not Found`: User not found
- `500 Internal Server Error`: Failed to delete user

**Notes:**

- Only SuperAdmin can delete users
- Cannot delete your own account (prevents accidental self-deletion)
- This will cascade delete the user's profile and all related data
- This operation cannot be undone

---

### GET /api/v2/users/stats/ - Get User Statistics

Get aggregated statistics about users in the system. This endpoint is restricted to Admin or Super Admin.

**Headers:**

- `Authorization: Bearer <access_token>` (required, Admin or SuperAdmin role)

**Response:**

**Success (200 OK):**

```json
{
  "total_users": 150,
  "active_users": 120,
  "users_by_role": {
    "FreeUser": 100,
    "ProUser": 30,
    "Admin": 15,
    "SuperAdmin": 5
  },
  "users_by_plan": {
    "free": 100,
    "5k": 20,
    "25k": 10,
    "100k": 5
  }
}
```

**Response Fields:**

- `total_users` (integer): Total number of users in the system
- `active_users` (integer): Number of active users (is_active = true)
- `users_by_role` (object): Count of users grouped by role
- `users_by_plan` (object): Count of users grouped by subscription plan

**Error (403 Forbidden) - Not Admin or Super Admin:**

```json
{
  "detail": "Admin or Super Admin role required"
}
```

**Status Codes:**

- `200 OK`: Statistics retrieved successfully
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Admin or SuperAdmin role required
- `500 Internal Server Error`: Failed to get user statistics

**Notes:**

- Both Admin and SuperAdmin can access this endpoint
- Statistics are calculated in real-time from the database
- Useful for dashboard and analytics purposes

---

### GET /api/v2/users/history/ - Get User History

Get user history records (Super Admin only). Returns paginated list of user registration and login events with IP geolocation data. Supports filtering by user_id and event_type.

**Headers:**

- `Authorization: Bearer <access_token>` (required, SuperAdmin role)

**Query Parameters:**

- `user_uuid` (string, UUID, optional): Filter by user UUID
- `event_type` (string, optional): Filter by event type. Valid values: `registration`, `login`
- `limit` (integer, optional, default: 100, min: 1, max: 1000): Maximum number of records to return
- `offset` (integer, optional, default: 0, min: 0): Number of records to skip (for pagination)

**Response:**

**Success (200 OK):**

```json
{
  "items": [
    {
      "uuid": "123e4567-e89b-12d3-a456-426614174000",
      "user_id": "223e4567-e89b-12d3-a456-426614174001",
      "event_type": "registration",
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "country": "United States",
      "city": "New York",
      "created_at": "2024-01-01T00:00:00Z"
    },
    {
      "uuid": "323e4567-e89b-12d3-a456-426614174002",
      "user_id": "223e4567-e89b-12d3-a456-426614174001",
      "event_type": "login",
      "ip_address": "192.168.1.2",
      "user_agent": "Mozilla/5.0...",
      "country": "United States",
      "city": "San Francisco",
      "created_at": "2024-01-15T12:00:00Z"
    }
  ],
  "total": 250,
  "limit": 100,
  "offset": 0
}
```

**Response Fields:**

- `items` (array): List of user history records
  - `uuid` (string, UUID): History record UUID
  - `user_uuid` (string, UUID): User UUID associated with this event
  - `event_type` (string): Event type - `registration` or `login`
  - `ip_address` (string, optional): IP address from the request
  - `user_agent` (string, optional): User agent string from the request
  - `country` (string, optional): Country from IP geolocation
  - `city` (string, optional): City from IP geolocation
  - `created_at` (datetime, ISO 8601): When the event occurred
- `total` (integer): Total number of history records matching the filters
- `limit` (integer): Maximum number of records returned
- `offset` (integer): Number of records skipped

**Error (403 Forbidden) - Not Super Admin:**

```json
{
  "detail": "You do not have permission to perform this action. SuperAdmin role required."
}
```

**Status Codes:**

- `200 OK`: User history retrieved successfully
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: SuperAdmin role required
- `500 Internal Server Error`: Failed to get user history

**Notes:**

- Only SuperAdmin can access this endpoint
- Returns paginated list of user registration and login events
- Includes IP geolocation data (country, city, IP, ISP, etc.) when available
- Supports filtering by `user_id` and `event_type`
- Events are recorded automatically during registration and login
- IP geolocation data is provided by the frontend in the request body (optional field). If not provided, history records will be created without geolocation data.

---

## Error Responses

All endpoints may return the following common error responses:

### 400 Bad Request

```json
{
  "detail": "Error message describing what went wrong"
}
```

Or field-specific errors:

```json
{
  "field_name": ["Error message for this field"]
}
```

### 401 Unauthorized

```json
{
  "detail": "Authentication credentials were not provided."
}
```

Or:

```json
{
  "detail": "Given token not valid for any token type"
}
```

### 403 Forbidden

Returned when the user does not have the required role:

```json
{
  "detail": "You do not have permission to perform this action. [Role] role required."
}
```

### 500 Internal Server Error

```json
{
  "detail": "An error occurred while processing the request."
}
```

---

## Notes

- All timestamps are in ISO 8601 format (UTC): `YYYY-MM-DDTHH:MM:SSZ`
- JWT tokens use HS256 algorithm with configurable expiration:
  - Access tokens: 30 minutes (default, configurable via `ACCESS_TOKEN_EXPIRE_MINUTES`)
  - Refresh tokens: 7 days (default, configurable via `REFRESH_TOKEN_EXPIRE_DAYS`)
- Profile creation happens automatically upon user registration with default values:
  - `role`: `FreeUser` (default role for new users)
  - `credits`: `50` (initial credits for free users)
  - `subscription_plan`: `free`
  - `subscription_status`: `active`
  - `notifications`: `{"weeklyReports": true, "newLeadAlerts": true}`
- Token refresh implements token rotation (new tokens issued, old tokens remain valid until expiration)
- The `email` field in the profile is read-only and synced from the User model
- Avatar uploads validate both file extension and file content (magic bytes) for security
- Password hashing uses bcrypt with automatic salt generation
- Password length is limited to 72 characters due to bcrypt's internal limitation
- The base URL in examples (`http://localhost:8000`) should be replaced with your actual API base URL

