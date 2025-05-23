# Mock Google OAuth2 Server

[![Security Policy](https://img.shields.io/badge/Security-Policy-green.svg)](./SECURITY.md)

![CI Status](https://github.com/chrisw-dev/golang-mock-oauth2-server/actions/workflows/ci.yml/badge.svg)
A lightweight server written in Go that simulates Google's OAuth2 authentication endpoints for development and testing purposes.

> **WARNING**: This server is intended for development and testing purposes only. It should never be used in production environments as it does not implement proper security measures required for real OAuth2 authentication.

## Overview
This mock server recreates Google's OAuth2 flow without requiring an internet connection or real Google credentials. It implements the following endpoints:

- `/authorize` - Authorization endpoint where users are redirected to authenticate
- `/token` - Token exchange endpoint to obtain access tokens
- `/userinfo` - User profile information endpoint
- `/.well-known/openid-configuration` - OpenID Connect discovery endpoint

## Use Cases
- Develop OAuth2 clients offline
- Write integration tests without depending on external services
- Simulate authentication flows without creating real Google credentials
- Reproduce specific OAuth2 scenarios for debugging

## Installation
Prerequisites
- Go 1.23 or higher

## Getting Started

```bash
# Clone the repository
git clone https://github.com/chrisw-dev/golang-mock-oauth2-server.git
cd golang-mock-oauth2-server

# Install dependencies
go mod download

# Build the server
go build -o mock-oauth2-server ./cmd/server
```

## Running the Server

```bash
# Start with default settings (port 8080)
./mock-oauth2-server

# Specify a custom port using the command-line flag (highest priority)
./mock-oauth2-server --port 9088

# Specify a custom hostname for the server URLs
./mock-oauth2-server --host http://mock-oauth2-server:9088

# Specify a custom port using environment variable (used if no command-line flag is provided)
MOCK_OAUTH_PORT=9088 ./mock-oauth2-server

# Specify a custom issuer URL using environment variable (useful in containerized environments)
MOCK_ISSUER_URL=http://mock-oauth2:8080 ./mock-oauth2-server
```

## Running with Docker

This project includes a Dockerfile that can be used to build and run the server in a container.

### Using GitHub Container Registry

The mock OAuth2 server is available as a pre-built Docker image from GitHub Container Registry (GHCR). This allows you to use the server without building it yourself.

```bash
# Pull the latest image
docker pull ghcr.io/chrisw-dev/golang-mock-oauth2-server:latest

# Run the container
docker run -p 8080:8080 ghcr.io/chrisw-dev/golang-mock-oauth2-server:latest
```

The following image tags are available:
- `latest` - The most recent build from the main branch
- `x.y.z` - Specific version (e.g., `1.0.0`)
- `x.y` - Latest patch version of a specific minor version (e.g., `1.0`)
- `sha-abcdef` - Specific commit SHA

You can also use the image in your Docker Compose file:

```yaml
version: '3'

services:
  mock-oauth2:
    image: ghcr.io/chrisw-dev/golang-mock-oauth2-server:latest
    ports:
      - "8080:8080"
    environment:
      - MOCK_USER_EMAIL=custom@example.com
      - MOCK_USER_NAME=Custom User
```

If you're using the image in a private environment, you may need to authenticate with GitHub Container Registry:

```bash
# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

### Using Docker

```bash
# Build the Docker image
docker build -t mock-oauth2-server .

# Run the container
docker run -p 8080:8080 mock-oauth2-server
```

### Using Docker Compose

For a more convenient setup, you can use Docker Compose. Create a `docker-compose.yml` file with the following content:

```yaml
version: '3'

services:
  mock-oauth2:
    build: .
    container_name: mock-oauth2-server
    ports:
      - "8080:8080"
    environment:
      - MOCK_OAUTH_PORT=8080
      - MOCK_USER_EMAIL=testuser@example.com
      - MOCK_USER_NAME=Test User
      - MOCK_TOKEN_EXPIRY=3600
      - MOCK_ISSUER_URL=http://mock-oauth2:8080
    # Mount a volume for custom fixtures if needed
    # volumes:
    #   - ./test/fixtures:/app/test/fixtures
    restart: unless-stopped

  # Example of how to integrate with your application
  # your-app:
  #   image: your-application-image
  #   container_name: your-app
  #   depends_on:
  #     - mock-oauth2
  #   environment:
  #     - OAUTH_AUTH_URL=http://mock-oauth2:8080/authorize
  #     - OAUTH_TOKEN_URL=http://mock-oauth2:8080/token
  #     - OAUTH_USERINFO_URL=http://mock-oauth2:8080/userinfo
  #   ports:
  #     - "8081:8081"
```

Then run:

```bash
# Start the services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the services
docker-compose down
```

This setup allows you to:
- Configure the mock server using environment variables
- Integrate it with your application in the same Docker network
- Access the mock server endpoints from your host at http://localhost:8080
- Access the mock server from other containers using the service name (e.g., http://mock-oauth2:8080)

> **REMINDER**: This server is for testing purposes only and should not be used in production environments.

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
│   │   ├── openid_config.go              # OpenID Connect discovery endpoint
│   │   ├── openid_config_test.go         # Test OpenID Connect discovery
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
- `openid_config.go` - Provides OpenID Connect discovery metadata
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

#### OpenID Connect Discovery Endpoint (`/.well-known/openid-configuration`)

Provides OpenID Connect (OIDC) configuration metadata for client auto-configuration.

**Method**: GET

**Response**: A JSON document with standard OIDC configuration

```json
{
  "issuer": "http://localhost:8080",
  "authorization_endpoint": "http://localhost:8080/authorize",
  "token_endpoint": "http://localhost:8080/token",
  "userinfo_endpoint": "http://localhost:8080/userinfo",
  "jwks_uri": "http://localhost:8080/jwks",
  "response_types_supported": ["code"],
  "subject_types_supported": ["public"],
  "id_token_signing_alg_values_supported": ["RS256"],
  "scopes_supported": ["openid", "email", "profile"],
  "token_endpoint_auth_methods_supported": ["client_secret_post", "client_secret_basic"],
  "claims_supported": [
    "sub", "iss", "name", "given_name", 
    "family_name", "email", "email_verified", "picture"
  ]
}
```

#### Configuration Endpoint

##### Dynamic Configuration Endpoint (`/config`)
Allows test code to dynamically configure the responses returned by the OAuth2 endpoints.

**Method**: POST

**Request Body**:

```json
{
  "tokens": {
    "access_token": "custom-access-token",
    "id_token": "custom-id-token",
    "expires_in": 1800,
    "user_info": {
      "sub": "custom-id-456",
      "name": "Updated Test User",
      "email": "updated@example.com",
      "email_verified": true,
      "picture": "https://example.com/updated.jpg"
    }
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

**IMPORTANT**: To update user information that will be returned by the `/userinfo` endpoint, you must include the user profile data inside the `tokens.user_info` object, not in the top-level `user_info` field. The top-level `user_info` field updates a different user object that is not used by the `/userinfo` endpoint.

This enables testing scenarios like:

- Testing how your application handles different user profiles
- Simulating error responses from OAuth endpoints
- Testing token expiration scenarios
- Creating custom authentication states for specific test cases

### Configuration

The server can be configured through multiple methods, with the following precedence order (highest to lowest):

1. Command-line flags (highest priority)
2. Environment variables
3. Default values (lowest priority)

Available configuration options:

- Port:
  - Command-line: `--port 9088`
  - Environment: `MOCK_OAUTH_PORT=9088`
  - Default: `8080`

- Issuer URL:
  - Command-line: `--host http://custom-hostname:9088`
  - Environment: `MOCK_ISSUER_URL=http://mock-oauth2:9088`
  - Default: `http://localhost:[port]`

- Other settings (environment variables only):
  - `MOCK_USER_EMAIL` - Email for the mock user (default: testuser@example.com)
  - `MOCK_USER_NAME` - Name for the mock user (default: Test User)
  - `MOCK_TOKEN_EXPIRY` - Token expiry in seconds (default: 3600)

The issuer URL is particularly important in containerized environments where the service name differs from "localhost". It affects the URLs returned in the OpenID Connect discovery document and needs to match what your OAuth client is configured to use.

### Example Usage

In your OAuth2 client application:

```go
// For a Go application using the golang.org/x/oauth2 package
import (
    "context"
    "golang.org/x/oauth2"
)

const (
    clientID     = "test-client-id"
    clientSecret = "test-client-secret"
    redirectURL  = "http://localhost:8081/callback"
)

// You can either specify endpoints manually:
oauth2Config := &oauth2.Config{
    ClientID:     clientID,
    ClientSecret: clientSecret,
    RedirectURL:  redirectURL,
    Scopes:       []string{"openid", "email", "profile"},
    Endpoint: oauth2.Endpoint{
        AuthURL:  "http://localhost:8080/authorize",
        TokenURL: "http://localhost:8080/token",
    },
}

// Or use the OpenID Connect discovery document:
provider, err := oidc.NewProvider(context.Background(), "http://localhost:8080")
if err != nil {
    // handle error
}

oauth2Config := &oauth2.Config{
    ClientID:     clientID,
    ClientSecret: clientSecret,
    RedirectURL:  redirectURL,
    Scopes:       []string{"openid", "email", "profile"},
    Endpoint:     provider.Endpoint(),
}
```

## Testing Strategies

1. Mock dependencies: Create mock versions of interfaces to isolate components during testing
1. Table-driven tests: Use Go's table-driven test pattern for testing multiple scenarios
1. HTTP tests: Use httptest package to test handlers
1. Config tests: Test different server configurations
1. Error handling: Test error scenarios using the /config endpoint


