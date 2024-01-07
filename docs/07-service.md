# Function

## Overview

The `Start` function initializes and runs the Authonomy service. It sets up the data store, initializes services with dependencies, and configures HTTP handlers and routes for the service's API. The function also integrates Swagger for API documentation and handles both the application's internal operations and user interactions.

## Function Signature

```go
func Start(dbPath, secret, port, ssiUrl string, reset bool)
```

### Parameters

- `dbPath` (string): Path to the database.
- `secret` (string): Database encryption key.
- `port` (string): Port number for the service to listen on.
- `ssiUrl` (string): URL of the Self-Sovereign Identity (SSI) service.
- `reset` (bool): Flag to reset the database on start.

### Functionality

- `Database Initialization`: Connects to the database using dbPath and secret. If reset is true, the database is cleared.
- `Services Initialization`: Sets up the SSI service client.
- `API Key Generation`: Generates a new API key for the service.
- `Swagger Integration`: Provides a Swagger UI endpoint for API documentation.
- `HTTP Handlers and Routes Setup`: Configures various endpoints for different functionalities like managing applications, authentication providers, policies, credentials, and user access.
- `Static Web Page Hosting`: Hosts a static web page for access token management.
- `Server Startup`: Starts the HTTP server on the specified port.

### Endpoints

```sh
/applications: Manage applications.
/auth-provider: Get, link, and unlink authentication providers.
/policies: Get, create, and attach policies.
/grant-access, /revoke-access: Manage access grants.
/verify-access, /issue-credential: Verify access and issue credentials.
/callback/: Handle callback operations.
/me/: User-related operations.
/signup: Sign up handler.
/get-access-token: Retrieve access tokens.
/request-access: Request access to resources.
/get-access-list: Get a list of access grants.
```

### Usage Example

To start the Authonomy service:

```sh
Start("/path/to/db", "secretKey", ":8080", "http://ssi-service-url", false)
```

This will start the service on port 8080, using the specified database path and SSI service URL, without resetting the database.

For detailed API usage and examples, refer to the Swagger documentation at `http://localhost:[port]/swagger`.
