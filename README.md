# Mock Google OAuth2 Server

![CI Status](https://github.com/chrisw-dev/golang-mock-oauth2-server/actions/workflows/ci.yml/badge.svg)
A lightweight server written in Go that simulates Google's OAuth2 authentication endpoints for development and testing purposes.

## Overview
This mock server recreates Google's OAuth2 flow without requiring an internet connection or real Google credentials. It implements the following endpoints:

- `/authorize` - Authorization endpoint where users are redirected to authenticate
- `/token` - Token exchange endpoint to obtain access tokens
- `/userinfo` - User profile information endpoint

## Use Cases
- Develop OAuth2 clients offline
- Write integration tests without depending on external services
- Simulate authentication flows without creating real Google credentials
- Reproduce specific OAuth2 scenarios for debugging

## Installation
Prerequisites
- Go 1.16 or higher

## Getting Started

```bash
# Clone the repository
git clone https://github.com/yourusername/golang-mock-oauth2-server.git
cd golang-mock-oauth2-server

# Install dependencies
go mod download

# Build the server
go build -o mock-oauth2-server
```

## Running the Server

```bash
# Start with default settings (port 8080)
./mock-oauth2-server

# Specify a custom port
./mock-oauth2-server -port 9000
```

## Architecture
The mock OAuth2 server is designed with a modular architecture that separates concerns and makes the codebase maintainable and extensible.

### Project Structure

```bash
golang-mock-oauth2-server/
├── cmd/
│   └── server/
│       └── main.go
│       └── main_test.go                  # Test server initialization
├── internal/
│   ├── config/
│   │   └── config.go
│   │   └── config_test.go                # Test configuration loading
│   ├── handlers/
│   │   ├── authorize.go
│   │   ├── authorize_test.go             # Test authorize handler
│   │   ├── token.go
│   │   ├── token_test.go                 # Test token handler
│   │   ├── userinfo.go
│   │   ├── userinfo_test.go              # Test userinfo handler
│   │   ├── config.go
│   │   └── config_test.go                # Test config handler
│   ├── middleware/
│   │   └── auth.go
│   │   └── auth_test.go                  # Test auth middleware
│   ├── models/
│   │   ├── token.go
│   │   ├── token_test.go                 # Test token models
│   │   ├── user.go
│   │   └── user_test.go                  # Test user models
│   ├── store/
│   │   ├── memory.go
│   │   └── memory_test.go                # Test in-memory store
│   └── server/
│       ├── server.go
│       └── server_test.go                # Test server setup
├── pkg/
│   └── oauth/
│       ├── provider.go
│       ├── provider_test.go              # Test provider interface
│       ├── google.go
│       └── google_test.go                # Test Google implementation
├── test/
│   ├── integration/                      # Integration tests
│   │   └── oauth_flow_test.go            # Test complete OAuth flow
│   └── fixtures/                         # Test data
│       ├── users.json                    # Sample user profiles
│       └── tokens.json                   # Sample tokens
├── go.mod
└── go.sum
```

### Core Components
#### Server (`internal/server/`)
The server component initializes and configures the HTTP server, registering routes and middleware.

#### Handlers (`internal/handlers/`)
Handler functions process incoming HTTP requests for each OAuth2 endpoint:

- `authorize.go` - Handles user authentication and generates authorization codes
- `token.go` - Exchanges authorization codes for access tokens
- `userinfo.go` - Returns user profile information
- `config.go` - Manages dynamic configuration for testing

#### Configuration (`internal/config/`)
Manages server settings from environment variables, command-line flags, and dynamic configuration changes.

#### Models (`internal/models/`)
Data structures representing tokens, user profiles, and configuration settings.

#### Store (`internal/store/`)
An in-memory data store maintains the server's state, including issued tokens and user profiles.

#### State Management
The server maintains state using an in-memory store that:

- Maps authorization codes to client information
- Tracks issued tokens and their expiration
- Stores configurable user profiles and token responses

#### Dynamic Configuration
The `/config` endpoint enables runtime modification of:

- User profile information returned by `/userinfo`
- Token responses from the `/token` endpoint
- Error scenarios for testing error handling

### Endpoint Documentation
#### Authorization Endpoint (`/authorize`)
Simulates Google's authorization page.

##### Parameters:

- `client_id` - OAuth2 client ID
- `redirect_uri` - URL to redirect after authorization
- `scope` - Requested permission scopes
- `response_type` - Must be "code"
- `state` - Optional state parameter

**Response**: Redirects to the provided `redirect_uri` with an authorization code.

#### Token Endpoint (`/token`)

Exchange authorization codes for access tokens.

##### Parameters:

- `grant_type` - Must be "authorization_code"
- `code` - The authorization code from the /authorize endpoint
- `client_id` - OAuth2 client ID
- `client_secret` - OAuth2 client secret
- `redirect_uri` - Must match the URI used in the authorization request

**Response**:

```json
{
  "access_token": "mock-access-token",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "mock-refresh-token",
  "id_token": "mock-id-token"
}
```

#### User Info Endpoint (`/userinfo`)

Retrieves mock user profile information.

##### Headers:

- `Authorization: Bearer {access_token}`

**Response**:

```json
{
  "sub": "12345678",
  "name": "Test User",
  "given_name": "Test",
  "family_name": "User",
  "email": "testuser@example.com",
  "picture": "https://example.com/photo.jpg",
  "email_verified": true
}
```

#### Configuration Endpoint

##### Dynamic Configuration Endpoint (`/config`)
Allows test code to dynamically configure the responses returned by the OAuth2 endpoints.

**Method**: POST

**Request Body**:

```json
{
  "user_info": {
    "sub": "custom-id-123",
    "name": "Custom Test User",
    "email": "custom@example.com",
    "email_verified": false
  },
  "tokens": {
    "access_token": "custom-access-token",
    "id_token": "custom-id-token",
    "expires_in": 1800
  },
  "error_scenario": {
    "endpoint": "token", 
    "error": "invalid_grant",
    "error_description": "Custom error for testing"
  }
}
```

**Response**:

```json
{
  "status": "success",
  "message": "Configuration updated"
}
```

This enables testing scenarios like:

- Testing how your application handles different user profiles
- Simulating error responses from OAuth endpoints
- Testing token expiration scenarios
- Creating custom authentication states for specific test cases



### Configuration

The server can be configured by setting environment variables:

MOCK_OAUTH_PORT - Server port (default: 8080)
MOCK_USER_EMAIL - Email for the mock user (default: testuser@example.com)
MOCK_USER_NAME - Name for the mock user (default: Test User)
MOCK_TOKEN_EXPIRY - Token expiry in seconds (default: 3600)

### Example Usage

In your OAuth2 client application:

```bash
const (
    clientID     = "test-client-id"
    clientSecret = "test-client-secret"
    redirectURL  = "http://localhost:8081/callback"
    authURL      = "http://localhost:8080/authorize"
    tokenURL     = "http://localhost:8080/token"
    userInfoURL  = "http://localhost:8080/userinfo"
)

// Configure OAuth2 client to use the mock server
oauth2Config := &oauth2.Config{
    ClientID:     clientID,
    ClientSecret: clientSecret,
    RedirectURL:  redirectURL,
    Scopes:       []string{"openid", "email", "profile"},
    Endpoint: oauth2.Endpoint{
        AuthURL:  authURL,
        TokenURL: tokenURL,
    },
}
```

## Testing Strategies

1. Mock dependencies: Create mock versions of interfaces to isolate components during testing
1. Table-driven tests: Use Go's table-driven test pattern for testing multiple scenarios
1. HTTP tests: Use httptest package to test handlers
1. Config tests: Test different server configurations
1. Error handling: Test error scenarios using the /config endpoint
