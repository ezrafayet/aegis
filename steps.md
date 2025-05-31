# Auth Service Development Steps

## Project Overview

**Goal**: Build a simple, dockerized drop-in auth service that handles GitHub OAuth, JWT sessions, API tokens, role-based auth, and user metadata with PostgreSQL backend.

## Development Roadmap

### 1. Project Setup & Infrastructure

- [ ] **Initialize Project Structure**
  - Choose technology stack (Go recommended for simplicity and performance)
  - Set up directory structure (`cmd/`, `internal/`, `pkg/`, `migrations/`, `docker/`)
  - Initialize go modules or package.json
  - Create `.gitignore` and basic project files

- [ ] **Configuration Management**
  - Implement config loader for `config.json`
  - Add environment variable overrides
  - Validate configuration on startup
  - Set up structured logging based on `log_level`

- [ ] **Docker Setup**
  - Create `Dockerfile` for the auth service
  - Create `docker-compose.yml` with PostgreSQL
  - Set up development environment
  - Configure health checks

### 2. Database Schema Design

- [ ] **Create Migration System**
  - Set up database migration tools
  - Create initial migration scripts

- [ ] **Design Core Tables**
  ```sql
  -- Users table
  CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    github_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    avatar_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
  );

  -- User roles (many-to-many with predefined roles)
  CREATE TABLE user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    PRIMARY KEY (user_id, role)
  );

  -- User metadata (flexible key-value store)
  CREATE TABLE user_metadata (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (user_id, key)
  );

  -- API tokens for programmatic access
  CREATE TABLE api_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
  );
  ```

### 3. Core Authentication Implementation

- [ ] **JWT Service**
  - Implement JWT token generation with configurable expiration
  - Create refresh token mechanism
  - Add token validation and parsing
  - Handle token blacklisting/revocation

- [ ] **GitHub OAuth Integration**
  - Implement OAuth flow initiation
  - Handle callback processing
  - Extract user information from GitHub API
  - Create or update user records

- [ ] **Session Management**
  - Cookie-based session handling
  - Secure cookie configuration from config
  - Session cleanup and expiration

### 4. API Endpoints Implementation

#### **Authentication Endpoints**
- [ ] `GET /auth/github` - Initiate GitHub OAuth flow
- [ ] `GET /auth/github/callback` - Handle GitHub OAuth callback
- [ ] `POST /auth/logout` - Logout user (invalidate session)
- [ ] `GET /auth/me` - Get current user info from JWT
- [ ] `POST /auth/refresh` - Refresh JWT token using refresh token
- [ ] `GET /auth/verify` - Verify JWT token (for other services)

#### **API Token Management**
- [ ] `GET /api/tokens` - List user's API tokens
- [ ] `POST /api/tokens` - Create new API token
- [ ] `DELETE /api/tokens/:id` - Revoke API token

#### **User Management**
- [ ] `GET /api/users/me` - Get current user profile
- [ ] `PUT /api/users/me` - Update user profile
- [ ] `GET /api/users/me/roles` - Get user roles
- [ ] `PUT /api/users/me/metadata` - Update user metadata

#### **Admin Endpoints** (platform_admin role only)
- [ ] `GET /admin/users` - List all users
- [ ] `PUT /admin/users/:id/roles` - Update user roles
- [ ] `GET /admin/users/:id` - Get user details
- [ ] `DELETE /admin/users/:id` - Delete user

#### **System Endpoints**
- [ ] `GET /health` - Health check
- [ ] `GET /info` - Service info and version

### 5. Middleware Development

- [ ] **Security Middleware**
  - CORS middleware using `allowed_origins` from config
  - JWT validation middleware
  - Role-based access control middleware
  - Rate limiting middleware
  - Request logging middleware

- [ ] **Validation Middleware**
  - Input validation and sanitization
  - Request size limiting
  - Content-type validation

### 6. Security Implementation

- [ ] **Token Security**
  - Secure random token generation
  - Password hashing for API tokens
  - Token rotation mechanisms

- [ ] **Application Security**
  - CSRF protection
  - Secure headers middleware
  - Input sanitization
  - SQL injection prevention

- [ ] **Configuration Security**
  - Validate JWT secret strength
  - Secure cookie settings enforcement
  - Environment-based security configs

### 7. Error Handling & Logging

- [ ] **Structured Logging**
  - Implement configurable log levels
  - Add request tracing
  - Security event logging
  - Performance metrics

- [ ] **Error Handling**
  - Consistent error response format
  - Proper HTTP status codes
  - Error recovery mechanisms
  - Graceful degradation

### 8. Testing Implementation

- [ ] **Unit Tests**
  - Test all service functions
  - Mock external dependencies
  - Test edge cases and error conditions

- [ ] **Integration Tests**
  - Test full OAuth flow
  - Database integration tests
  - API endpoint tests

- [ ] **Security Tests**
  - Test authentication bypass attempts
  - Test authorization controls
  - Test input validation

### 9. Documentation

- [ ] **API Documentation**
  - OpenAPI/Swagger specification
  - Endpoint documentation with examples
  - Authentication guide

- [ ] **Deployment Documentation**
  - Docker deployment guide
  - Configuration reference
  - Environment setup guide

- [ ] **Integration Documentation**
  - Client-side integration examples
  - SDK/library examples for different languages

### 10. Client Integration Snippets

- [ ] **Next.js Integration**
  ```typescript
  // useAuth hook example
  // API route protection example
  // Login/logout components
  ```

- [ ] **Go Integration**
  ```go
  // JWT validation middleware
  // User context extraction
  // Role checking utilities
  ```

- [ ] **Node.js Integration**
  ```javascript
  // Express middleware
  // Authentication helpers
  // API client example
  ```

### 11. Production Readiness

- [ ] **Performance Optimization**
  - Connection pooling
  - Caching strategies
  - Query optimization

- [ ] **Monitoring & Observability**
  - Health check endpoints
  - Metrics collection
  - Log aggregation setup

- [ ] **Deployment Configuration**
  - Production Dockerfile optimization
  - Kubernetes manifests (optional)
  - CI/CD pipeline setup

## Configuration Requirements

Based on `config.example.json`, ensure support for:

- **Application settings**: name, URL, log level
- **Database**: PostgreSQL connection URL
- **JWT**: secret, access/refresh token expiration
- **OAuth**: GitHub client ID/secret, enabled flag
- **Auth**: redirect URLs, allowed origins
- **Cookies**: domain, security settings
- **Users**: default roles, metadata schema validation

## Success Criteria

- [ ] Single Docker container deployment
- [ ] Configuration-driven setup
- [ ] GitHub OAuth working end-to-end
- [ ] JWT-based session management
- [ ] Role-based access control
- [ ] API token generation and validation
- [ ] User metadata management
- [ ] Ready-to-use client integration examples

## Next Steps

1. Start with project structure and configuration loading
2. Set up database and migrations
3. Implement core JWT functionality
4. Build GitHub OAuth flow
5. Create API endpoints incrementally
6. Add security middleware and testing
7. Create documentation and integration examples 