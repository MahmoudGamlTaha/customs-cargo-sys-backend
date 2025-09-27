# API Documentation

## Overview

The Request Management System API provides a comprehensive REST interface for managing users, requests, companies, and branches. The API follows RESTful principles and uses JSON for data exchange.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

The API uses JWT (JSON Web Token) for authentication. After successful login, include the token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

## Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error": {
    "code": "ERROR_CODE",
    "details": "Detailed error information"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request data |
| 401 | Unauthorized - Authentication required |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 422 | Unprocessable Entity - Validation error |
| 500 | Internal Server Error - Server error |

## Pagination

List endpoints support pagination with the following query parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number |
| page_size | integer | 20 | Items per page (max 100) |

### Pagination Response
```json
{
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total_items": 150,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

## Filtering and Sorting

List endpoints support filtering and sorting:

| Parameter | Description |
|-----------|-------------|
| search | Search term for text fields |
| sort_by | Field to sort by |
| sort_order | Sort order (asc/desc) |

## Authentication Endpoints

### Register User

Register a new client user.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "company_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "role": "client",
    "company": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Tech Solutions Inc",
      "code": "TECH001"
    },
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Login User

Authenticate user and receive JWT token.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-16T10:30:00Z",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "username": "johndoe",
      "email": "john@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "client"
    }
  }
}
```

## User Management Endpoints

### Get User Profile

Get current user's profile information.

**Endpoint:** `GET /users/profile`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "message": "User profile retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "role": "client",
    "company": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Tech Solutions Inc",
      "code": "TECH001"
    },
    "is_active": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Update User Profile

Update current user's profile information.

**Endpoint:** `PUT /users/profile`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "first_name": "John",
  "last_name": "Smith",
  "phone": "+1234567891"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    // Updated user object
  }
}
```

### Change Password

Change current user's password.

**Endpoint:** `POST /users/change-password`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

## Request Management Endpoints

### Create Request

Create a new request.

**Endpoint:** `POST /requests`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "New Service Request",
  "description": "I need assistance with setting up a new service account for our department. This includes configuring access permissions and setting up the necessary integrations.",
  "client_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Note:** `client_id` is only required when staff/admin creates request for a client. Clients creating their own requests should omit this field.

**Response:**
```json
{
  "success": true,
  "message": "Request created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "title": "New Service Request",
    "description": "I need assistance with setting up a new service account...",
    "status": "pending",
    "serial_number": null,
    "client_id": "550e8400-e29b-41d4-a716-446655440001",
    "created_by_staff_id": null,
    "approved_by_staff_id": null,
    "rejection_reason": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### List Requests

List requests with pagination and filtering.

**Endpoint:** `GET /requests`

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (integer): Page number (default: 1)
- `page_size` (integer): Items per page (default: 20, max: 100)
- `search` (string): Search in title and description
- `status` (string): Filter by status (pending, approved, rejected)
- `start_date` (string): Filter by creation date (YYYY-MM-DD)
- `end_date` (string): Filter by creation date (YYYY-MM-DD)
- `sort_by` (string): Sort field (created_at, updated_at, title)
- `sort_order` (string): Sort order (asc, desc)

**Response:**
```json
{
  "success": true,
  "message": "Requests retrieved",
  "data": {
    "requests": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440002",
        "title": "New Service Request",
        "description": "I need assistance with...",
        "status": "pending",
        "serial_number": null,
        "client": {
          "id": "550e8400-e29b-41d4-a716-446655440001",
          "username": "johndoe",
          "first_name": "John",
          "last_name": "Doe",
          "company": {
            "name": "Tech Solutions Inc",
            "code": "TECH001"
          }
        },
        "created_by_staff": null,
        "approved_by_staff": null,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_items": 1,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

### Get Request Details

Get detailed information about a specific request.

**Endpoint:** `GET /requests/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "message": "Request retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "title": "New Service Request",
    "description": "I need assistance with setting up a new service account...",
    "status": "approved",
    "serial_number": "2024-000001",
    "client": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "username": "johndoe",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@example.com",
      "company": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Tech Solutions Inc",
        "code": "TECH001"
      }
    },
    "created_by_staff": null,
    "approved_by_staff": {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "username": "staff1",
      "first_name": "Alice",
      "last_name": "Johnson",
      "branch": {
        "name": "Main Branch",
        "code": "MAIN001"
      }
    },
    "rejection_reason": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### Update Request

Update a pending request.

**Endpoint:** `PUT /requests/{id}`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "Updated Service Request",
  "description": "Updated description with more details about the service account requirements."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Request updated successfully",
  "data": {
    // Updated request object
  }
}
```

### Approve Request

Approve a pending request (Staff/Admin only).

**Endpoint:** `POST /requests/{id}/approve`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "approved_by_staff_id": "550e8400-e29b-41d4-a716-446655440003"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Request approved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "title": "New Service Request",
    "description": "I need assistance with...",
    "status": "approved",
    "serial_number": "2024-000001",
    "client_id": "550e8400-e29b-41d4-a716-446655440001",
    "approved_by_staff_id": "550e8400-e29b-41d4-a716-446655440003",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### Reject Request

Reject a pending request (Staff/Admin only).

**Endpoint:** `POST /requests/{id}/reject`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "approved_by_staff_id": "550e8400-e29b-41d4-a716-446655440003",
  "rejection_reason": "Insufficient information provided. Please provide more details about the specific service account requirements and business justification."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Request rejected successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "title": "New Service Request",
    "description": "I need assistance with...",
    "status": "rejected",
    "serial_number": null,
    "client_id": "550e8400-e29b-41d4-a716-446655440001",
    "approved_by_staff_id": "550e8400-e29b-41d4-a716-446655440003",
    "rejection_reason": "Insufficient information provided...",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### Get Request History

Get the history of status changes for a request.

**Endpoint:** `GET /requests/{id}/history`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "message": "Request history retrieved",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "request_id": "550e8400-e29b-41d4-a716-446655440002",
      "old_status": null,
      "new_status": "pending",
      "changed_by": {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "username": "johndoe",
        "first_name": "John",
        "last_name": "Doe"
      },
      "change_reason": "Request created",
      "created_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440011",
      "request_id": "550e8400-e29b-41d4-a716-446655440002",
      "old_status": "pending",
      "new_status": "approved",
      "changed_by": {
        "id": "550e8400-e29b-41d4-a716-446655440003",
        "username": "staff1",
        "first_name": "Alice",
        "last_name": "Johnson"
      },
      "change_reason": "Request meets all requirements",
      "created_at": "2024-01-15T11:00:00Z"
    }
  ]
}
```

## Company Management Endpoints

### List Companies

List all companies with pagination.

**Endpoint:** `GET /companies`

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (integer): Page number
- `page_size` (integer): Items per page
- `search` (string): Search in name and code

**Response:**
```json
{
  "success": true,
  "message": "Companies retrieved",
  "data": {
    "companies": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Tech Solutions Inc",
        "code": "TECH001",
        "address": "123 Tech Street, Silicon Valley, CA 94000",
        "phone": "+1-555-0123",
        "email": "info@techsolutions.com",
        "created_at": "2024-01-15T09:00:00Z",
        "updated_at": "2024-01-15T09:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_items": 1,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

### Get All Companies

Get all companies without pagination (for dropdown lists).

**Endpoint:** `GET /companies/all`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "message": "Companies retrieved",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Tech Solutions Inc",
      "code": "TECH001",
      "address": "123 Tech Street, Silicon Valley, CA 94000",
      "phone": "+1-555-0123",
      "email": "info@techsolutions.com"
    }
  ]
}
```

### Get Company Details

Get detailed information about a specific company.

**Endpoint:** `GET /companies/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "message": "Company retrieved",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Tech Solutions Inc",
    "code": "TECH001",
    "address": "123 Tech Street, Silicon Valley, CA 94000",
    "phone": "+1-555-0123",
    "email": "info@techsolutions.com",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

## Branch Management Endpoints

### List Branches

List all branches with pagination.

**Endpoint:** `GET /branches`

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (integer): Page number
- `page_size` (integer): Items per page
- `search` (string): Search in name and code

**Response:**
```json
{
  "success": true,
  "message": "Branches retrieved",
  "data": {
    "branches": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440020",
        "name": "Main Branch",
        "code": "MAIN001",
        "address": "456 Main Street, Downtown, NY 10001",
        "phone": "+1-555-0456",
        "email": "main@company.com",
        "created_at": "2024-01-15T09:00:00Z",
        "updated_at": "2024-01-15T09:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_items": 1,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

### Get All Branches

Get all branches without pagination (for dropdown lists).

**Endpoint:** `GET /branches/all`

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "success": true,
  "message": "Branches retrieved",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "name": "Main Branch",
      "code": "MAIN001",
      "address": "456 Main Street, Downtown, NY 10001",
      "phone": "+1-555-0456",
      "email": "main@company.com"
    }
  ]
}
```

## Admin Endpoints

### Create User

Create a new user (Admin only).

**Endpoint:** `POST /admin/users`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "username": "newstaff",
  "email": "staff@example.com",
  "password": "securepassword123",
  "first_name": "Jane",
  "last_name": "Smith",
  "phone": "+1234567890",
  "role": "staff",
  "company_id": "550e8400-e29b-41d4-a716-446655440000",
  "branch_id": "550e8400-e29b-41d4-a716-446655440020"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440004",
    "username": "newstaff",
    "email": "staff@example.com",
    "first_name": "Jane",
    "last_name": "Smith",
    "phone": "+1234567890",
    "role": "staff",
    "company": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Tech Solutions Inc",
      "code": "TECH001"
    },
    "branch": {
      "id": "550e8400-e29b-41d4-a716-446655440020",
      "name": "Main Branch",
      "code": "MAIN001"
    },
    "is_active": true,
    "created_at": "2024-01-15T12:00:00Z"
  }
}
```

### List All Users

List all users with pagination and filtering (Admin only).

**Endpoint:** `GET /admin/users`

**Headers:** `Authorization: Bearer <token>`

**Query Parameters:**
- `page` (integer): Page number
- `page_size` (integer): Items per page
- `search` (string): Search in username, email, first_name, last_name
- `role` (string): Filter by role (client, staff, admin)
- `company_id` (string): Filter by company ID
- `branch_id` (string): Filter by branch ID

**Response:**
```json
{
  "success": true,
  "message": "Users retrieved",
  "data": {
    "users": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "username": "johndoe",
        "email": "john@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "phone": "+1234567890",
        "role": "client",
        "company": {
          "id": "550e8400-e29b-41d4-a716-446655440000",
          "name": "Tech Solutions Inc",
          "code": "TECH001"
        },
        "branch": null,
        "is_active": true,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_items": 1,
      "total_pages": 1,
      "has_next": false,
      "has_prev": false
    }
  }
}
```

### Update User Role

Update a user's role and assignments (Admin only).

**Endpoint:** `PUT /admin/users/{id}/role`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "role": "staff",
  "company_id": "550e8400-e29b-41d4-a716-446655440000",
  "branch_id": "550e8400-e29b-41d4-a716-446655440020"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User role updated successfully",
  "data": {
    // Updated user object
  }
}
```

### Create Company

Create a new company (Admin only).

**Endpoint:** `POST /admin/companies`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "name": "New Company Ltd",
  "code": "NEW001",
  "address": "789 Business Ave, Corporate City, TX 75001",
  "phone": "+1-555-0789",
  "email": "info@newcompany.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Company created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440030",
    "name": "New Company Ltd",
    "code": "NEW001",
    "address": "789 Business Ave, Corporate City, TX 75001",
    "phone": "+1-555-0789",
    "email": "info@newcompany.com",
    "created_at": "2024-01-15T13:00:00Z",
    "updated_at": "2024-01-15T13:00:00Z"
  }
}
```

### Create Branch

Create a new branch (Admin only).

**Endpoint:** `POST /admin/branches`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "name": "South Branch",
  "code": "SOUTH001",
  "address": "321 South Street, Southern District, FL 33101",
  "phone": "+1-555-0321",
  "email": "south@company.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Branch created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440040",
    "name": "South Branch",
    "code": "SOUTH001",
    "address": "321 South Street, Southern District, FL 33101",
    "phone": "+1-555-0321",
    "email": "south@company.com",
    "created_at": "2024-01-15T13:30:00Z",
    "updated_at": "2024-01-15T13:30:00Z"
  }
}
```

## Error Handling

### Common Error Responses

#### Validation Error (422)
```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": {
      "validation_errors": "Field 'email' is required"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Unauthorized (401)
```json
{
  "success": false,
  "message": "Authentication required",
  "error": {
    "code": "UNAUTHORIZED",
    "details": "Invalid or missing authentication token"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Forbidden (403)
```json
{
  "success": false,
  "message": "Access denied",
  "error": {
    "code": "FORBIDDEN",
    "details": "Insufficient permissions for this operation"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Not Found (404)
```json
{
  "success": false,
  "message": "Resource not found",
  "error": {
    "code": "NOT_FOUND",
    "details": "The requested resource could not be found"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Conflict (409)
```json
{
  "success": false,
  "message": "Resource already exists",
  "error": {
    "code": "CONFLICT",
    "details": "Username already exists"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Authentication endpoints**: 5 requests per second
- **General API endpoints**: 10 requests per second
- **Burst allowance**: 20 requests

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 9
X-RateLimit-Reset: 1642248600
```

## Data Models

### User Model
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "phone": "string",
  "role": "client|staff|admin",
  "company_id": "uuid",
  "branch_id": "uuid|null",
  "is_active": "boolean",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Request Model
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "status": "pending|approved|rejected",
  "serial_number": "string|null",
  "client_id": "uuid",
  "created_by_staff_id": "uuid|null",
  "approved_by_staff_id": "uuid|null",
  "rejection_reason": "string|null",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Company Model
```json
{
  "id": "uuid",
  "name": "string",
  "code": "string",
  "address": "string",
  "phone": "string",
  "email": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Branch Model
```json
{
  "id": "uuid",
  "name": "string",
  "code": "string",
  "address": "string",
  "phone": "string",
  "email": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## SDK Examples

### JavaScript/Node.js Example

```javascript
const axios = require('axios');

class RequestManagementAPI {
  constructor(baseURL, token = null) {
    this.baseURL = baseURL;
    this.token = token;
    this.client = axios.create({
      baseURL: this.baseURL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add auth interceptor
    this.client.interceptors.request.use((config) => {
      if (this.token) {
        config.headers.Authorization = `Bearer ${this.token}`;
      }
      return config;
    });
  }

  async login(username, password) {
    const response = await this.client.post('/auth/login', {
      username,
      password,
    });
    this.token = response.data.data.token;
    return response.data;
  }

  async createRequest(title, description) {
    const response = await this.client.post('/requests', {
      title,
      description,
    });
    return response.data;
  }

  async listRequests(page = 1, pageSize = 20) {
    const response = await this.client.get('/requests', {
      params: { page, page_size: pageSize },
    });
    return response.data;
  }

  async approveRequest(requestId, staffId) {
    const response = await this.client.post(`/requests/${requestId}/approve`, {
      approved_by_staff_id: staffId,
    });
    return response.data;
  }
}

// Usage
const api = new RequestManagementAPI('http://localhost:8080/api/v1');

async function example() {
  // Login
  await api.login('admin', 'admin123');
  
  // Create request
  const request = await api.createRequest(
    'New Service Request',
    'Need help with account setup'
  );
  
  // List requests
  const requests = await api.listRequests();
  console.log(requests);
}
```

### Python Example

```python
import requests
import json

class RequestManagementAPI:
    def __init__(self, base_url, token=None):
        self.base_url = base_url
        self.token = token
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json'})
    
    def _get_headers(self):
        headers = {}
        if self.token:
            headers['Authorization'] = f'Bearer {self.token}'
        return headers
    
    def login(self, username, password):
        response = self.session.post(
            f'{self.base_url}/auth/login',
            json={'username': username, 'password': password}
        )
        response.raise_for_status()
        data = response.json()
        self.token = data['data']['token']
        return data
    
    def create_request(self, title, description):
        response = self.session.post(
            f'{self.base_url}/requests',
            json={'title': title, 'description': description},
            headers=self._get_headers()
        )
        response.raise_for_status()
        return response.json()
    
    def list_requests(self, page=1, page_size=20):
        response = self.session.get(
            f'{self.base_url}/requests',
            params={'page': page, 'page_size': page_size},
            headers=self._get_headers()
        )
        response.raise_for_status()
        return response.json()
    
    def approve_request(self, request_id, staff_id):
        response = self.session.post(
            f'{self.base_url}/requests/{request_id}/approve',
            json={'approved_by_staff_id': staff_id},
            headers=self._get_headers()
        )
        response.raise_for_status()
        return response.json()

# Usage
api = RequestManagementAPI('http://localhost:8080/api/v1')

# Login
api.login('admin', 'admin123')

# Create request
request = api.create_request(
    'New Service Request',
    'Need help with account setup'
)

# List requests
requests = api.list_requests()
print(json.dumps(requests, indent=2))
```

## Postman Collection

A Postman collection is available for testing the API. Import the collection and set the following environment variables:

- `base_url`: `http://localhost:8080/api/v1`
- `token`: JWT token (automatically set after login)

The collection includes pre-request scripts for automatic token management and comprehensive test cases for all endpoints.

## Changelog

### Version 1.0.0
- Initial API release
- Complete user management endpoints
- Request workflow implementation
- Company and branch management
- Admin functionality
- Comprehensive error handling
- Rate limiting implementation

---

For more information, see the [README](README.md) file or contact the development team.

