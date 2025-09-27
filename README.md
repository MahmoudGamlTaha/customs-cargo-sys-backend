# Request Management System

A comprehensive Go web application with PostgreSQL database that provides user authentication, role-based access control, and a complete request management workflow. The system supports multiple user roles (client, staff, admin) with specific permissions and features automated request approval/rejection with unique serial number generation.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [User Roles and Permissions](#user-roles-and-permissions)
- [Database Schema](#database-schema)
- [Development](#development)
- [Deployment](#deployment)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

### Core Functionality

- **User Management**: Complete user registration, authentication, and profile management
- **Role-Based Access Control**: Three distinct user roles with specific permissions
- **Request Management**: Full lifecycle request handling with approval/rejection workflow
- **Serial Number Generation**: Automatic unique serial number generation for approved requests
- **Company and Branch Management**: Organizational structure support
- **Audit Trail**: Complete request history tracking

### Technical Features

- **RESTful API**: Clean, well-documented REST API endpoints
- **JWT Authentication**: Secure token-based authentication
- **Password Security**: Argon2id password hashing
- **Database Migrations**: Automated database setup and updates
- **Docker Support**: Complete containerization with Docker Compose
- **Production Ready**: Nginx reverse proxy, rate limiting, and security headers
- **Modular Architecture**: Clean separation of concerns with repository pattern

## Architecture

The application follows a modular, layered architecture:

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handler/         # HTTP handlers (controllers)
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models and DTOs
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic layer
│   └── utils/           # Utility functions
├── pkg/
│   ├── auth/            # Authentication utilities
│   ├── database/        # Database connection
│   └── logger/          # Logging utilities
├── api/                 # API documentation
├── migrations/          # Database migrations
├── scripts/             # Deployment and utility scripts
└── docs/                # Additional documentation
```

### Design Patterns

- **Repository Pattern**: Abstracts data access logic
- **Service Layer**: Encapsulates business logic
- **Dependency Injection**: Loose coupling between components
- **Middleware Pattern**: Cross-cutting concerns handling

## Prerequisites

- **Go**: Version 1.21 or higher
- **PostgreSQL**: Version 12 or higher
- **Docker**: Version 20.10 or higher (optional)
- **Docker Compose**: Version 2.0 or higher (optional)

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd request-management-system
   ```

2. **Start the application**:
   ```bash
   docker-compose up -d
   ```

3. **Verify the installation**:
   ```bash
   curl http://localhost:8080/health
   ```

4. **Access the application**:
   - API Base URL: `http://localhost:8080/api/v1`
   - Health Check: `http://localhost:8080/health`
   - pgAdmin (optional): `http://localhost:5050`

### Default Users

The system comes with pre-configured users for testing:

| Username | Password | Role | Description |
|----------|----------|------|-------------|
| admin | admin123 | admin | System administrator |
| staff1 | staff123 | staff | Staff member (Main Branch) |
| staff2 | staff123 | staff | Staff member (North Branch) |
| client1 | client123 | client | Client (Tech Solutions Inc) |
| client2 | client123 | client | Client (Global Services Ltd) |

## Installation

### Manual Installation

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL database**:
   ```bash
   # Create database
   createdb request_management_system
   
   # Run migrations
   psql -d request_management_system -f database_schema.sql
   ```

3. **Configure environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Build and run**:
   ```bash
   go build -o app ./cmd/server
   ./app
   ```

### Using Make Commands

The project includes a Makefile for common operations:

```bash
# Install development tools
make install-tools

# Run tests
make test

# Build application
make build

# Run with live reload
make dev

# Set up database
make db-setup

# Deploy with Docker
make docker-run
```

## Configuration

### Environment Variables

The application uses environment variables for configuration. Copy `.env.example` to `.env` and modify as needed:

```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
SERVER_IDLE_TIMEOUT=60s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=request_management_system
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=24h
JWT_ISSUER=request-management-system

# Application Configuration
APP_NAME=Request Management System
APP_VERSION=1.0.0
APP_ENV=development
LOG_LEVEL=info
```

### Production Configuration

For production deployment, ensure you:

1. **Change default passwords and secrets**
2. **Use strong JWT secret key**
3. **Configure SSL/TLS certificates**
4. **Set up proper database backups**
5. **Configure monitoring and logging**

## API Documentation

### Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

### Base URL

```
http://localhost:8080/api/v1
```

### Endpoints Overview

#### Authentication Endpoints

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| POST | `/auth/register` | Register new client user | Public |
| POST | `/auth/login` | User login | Public |

#### User Management

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/users/profile` | Get current user profile | Authenticated |
| PUT | `/users/profile` | Update current user profile | Authenticated |
| POST | `/users/change-password` | Change password | Authenticated |

#### Request Management

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| POST | `/requests` | Create new request | Authenticated |
| GET | `/requests` | List requests | Authenticated |
| GET | `/requests/{id}` | Get request details | Authenticated |
| PUT | `/requests/{id}` | Update request | Authenticated |
| POST | `/requests/{id}/approve` | Approve request | Staff/Admin |
| POST | `/requests/{id}/reject` | Reject request | Staff/Admin |
| GET | `/requests/{id}/history` | Get request history | Authenticated |

#### Company Management

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/companies` | List companies | Authenticated |
| GET | `/companies/all` | Get all companies | Authenticated |
| GET | `/companies/{id}` | Get company details | Authenticated |

#### Branch Management

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/branches` | List branches | Authenticated |
| GET | `/branches/all` | Get all branches | Authenticated |
| GET | `/branches/{id}` | Get branch details | Authenticated |

#### Admin Endpoints

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| POST | `/admin/users` | Create user | Admin |
| GET | `/admin/users` | List all users | Admin |
| GET | `/admin/users/{id}` | Get user details | Admin |
| PUT | `/admin/users/{id}/role` | Update user role | Admin |
| DELETE | `/admin/users/{id}` | Delete user | Admin |
| POST | `/admin/companies` | Create company | Admin |
| PUT | `/admin/companies/{id}` | Update company | Admin |
| DELETE | `/admin/companies/{id}` | Delete company | Admin |
| POST | `/admin/branches` | Create branch | Admin |
| PUT | `/admin/branches/{id}` | Update branch | Admin |
| DELETE | `/admin/branches/{id}` | Delete branch | Admin |
| DELETE | `/admin/requests/{id}` | Delete request | Admin |

### Example API Calls

#### User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newclient",
    "email": "client@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "role": "client",
    "company_id": "company-uuid-here"
  }'
```

#### User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

#### Create Request

```bash
curl -X POST http://localhost:8080/api/v1/requests \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "title": "New Service Request",
    "description": "I need assistance with setting up a new service account for our department."
  }'
```

#### Approve Request

```bash
curl -X POST http://localhost:8080/api/v1/requests/{request_id}/approve \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt_token>" \
  -d '{
    "approved_by_staff_id": "staff-uuid-here"
  }'
```

## User Roles and Permissions

### Client Role

**Capabilities:**
- Register for an account (self-registration)
- Create requests for themselves
- View and edit their own pending requests
- View history of their requests
- Update their profile information

**Restrictions:**
- Cannot see other clients' requests
- Cannot approve or reject requests
- Cannot access admin functions

### Staff Role

**Capabilities:**
- All client capabilities
- Create requests on behalf of clients
- View all requests in the system
- Approve or reject pending requests
- View request history for all requests
- Must be assigned to a branch by admin

**Restrictions:**
- Cannot manage users, companies, or branches
- Cannot delete requests

### Admin Role

**Capabilities:**
- All staff capabilities
- Create, update, and delete users
- Assign roles and manage user permissions
- Create, update, and delete companies
- Create, update, and delete branches
- Delete requests (only pending ones)
- Full system access

**Special Notes:**
- Admin users are created by other admins
- At least one admin must exist in the system

## Database Schema

### Core Tables

#### Users Table
- Stores user information and authentication data
- Links to companies (for clients) and branches (for staff)
- Includes role-based access control fields

#### Companies Table
- Stores company information
- Each client must belong to a company
- Includes unique company codes

#### Branches Table
- Stores branch information
- Each staff member must belong to a branch
- Includes unique branch codes

#### Requests Table
- Stores request information and status
- Links to client and staff users
- Includes approval/rejection workflow fields
- Auto-generates serial numbers for approved requests

#### Request History Table
- Audit trail for request status changes
- Tracks who made changes and when
- Stores change reasons

### Key Features

#### Serial Number Generation
- Automatic generation for approved requests
- Format: YYYY-NNNNNN (e.g., 2024-000001)
- Unique across the entire system
- Year-based sequence numbering

#### Audit Trail
- Complete history of request status changes
- Tracks user actions and timestamps
- Immutable audit records

#### Data Integrity
- Foreign key constraints
- Check constraints for business rules
- Unique constraints for codes and usernames
- Automatic timestamp updates

## Development

### Project Structure

The project follows Go best practices with a clean architecture:

- **cmd/**: Application entry points
- **internal/**: Private application code
- **pkg/**: Public library code
- **api/**: API specifications
- **scripts/**: Utility scripts
- **docs/**: Documentation

### Development Workflow

1. **Set up development environment**:
   ```bash
   make install-tools
   make dev-up
   ```

2. **Run tests**:
   ```bash
   make test
   make test-coverage
   ```

3. **Code formatting and linting**:
   ```bash
   make fmt
   make lint
   ```

4. **Security scanning**:
   ```bash
   make security
   ```

### Adding New Features

1. **Database Changes**: Add migration files to `migrations/`
2. **Models**: Update models in `internal/models/`
3. **Repository**: Add data access methods in `internal/repository/`
4. **Service**: Implement business logic in `internal/service/`
5. **Handler**: Add HTTP handlers in `internal/handler/`
6. **Routes**: Update routes in `cmd/server/main.go`

### Testing

The project includes comprehensive testing:

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test package
go test ./internal/service/...
```

## Deployment

### Development Deployment

```bash
# Using Docker Compose
docker-compose up -d

# Using Make
make docker-run

# Manual deployment
make dev-up
```

### Production Deployment

1. **Prepare environment**:
   ```bash
   cp .env.example .env
   # Edit .env with production values
   ```

2. **Deploy with script**:
   ```bash
   ENVIRONMENT=production ./scripts/deploy.sh
   ```

3. **Or deploy manually**:
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
   ```

### Production Considerations

#### Security
- Change all default passwords
- Use strong JWT secrets
- Configure SSL/TLS certificates
- Set up firewall rules
- Enable audit logging

#### Performance
- Configure connection pooling
- Set up database indexing
- Enable response compression
- Configure rate limiting
- Monitor resource usage

#### Monitoring
- Set up health checks
- Configure log aggregation
- Monitor database performance
- Set up alerting
- Track API metrics

### Scaling

#### Horizontal Scaling
- Multiple application instances behind load balancer
- Database read replicas
- Redis for session storage
- CDN for static assets

#### Vertical Scaling
- Increase server resources
- Optimize database queries
- Tune connection pools
- Cache frequently accessed data

## Testing

### Test Coverage

The project includes comprehensive test coverage:

- **Unit Tests**: Individual component testing
- **Integration Tests**: Database and API testing
- **End-to-End Tests**: Complete workflow testing

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/service

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Test Data

Test data is automatically created during database setup:
- Sample companies and branches
- Test users for each role
- Sample requests in various states

## Contributing

### Development Setup

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

### Code Standards

- Follow Go best practices
- Use meaningful variable names
- Add comments for complex logic
- Include tests for new features
- Update documentation as needed

### Commit Guidelines

- Use clear, descriptive commit messages
- Reference issue numbers when applicable
- Keep commits focused and atomic
- Follow conventional commit format

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

For support and questions:

1. Check the documentation
2. Search existing issues
3. Create a new issue with detailed information
4. Include logs and error messages
5. Provide steps to reproduce problems

## Changelog

### Version 1.0.0
- Initial release
- Complete user management system
- Request workflow implementation
- Role-based access control
- Docker deployment support
- Comprehensive API documentation

---

**Built with ❤️ using Go, PostgreSQL, and modern development practices.**

