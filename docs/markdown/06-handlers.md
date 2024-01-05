# Service Handlers

## Overview

The `handlers` package in the Authonomy service provides HTTP handlers for various functionalities like application management, user authentication, and access control. This package is a crucial part of the service's REST API.

## Handlers

### AppHandler

Handles application-related requests.

#### NewAppHandler

- **Purpose**: Creates a new instance of `AppHandler`.
- **Parameters**: `ssiService` (*services.SsiClient), `db` (*store.Store).

#### HandleApplications

- **Purpose**: Routes application requests based on the HTTP method.
- **Methods**: `GET` (getApplications), `POST` (createApplication).

#### getApplications

- **Endpoint**: `/applications` (GET)
- **Description**: Retrieves a list of all applications.
- **Responses**: 200 (Array of `models.ApplicationResponse`), 500 (Internal Server Error).

#### createApplication

- **Endpoint**: `/applications` (POST)
- **Description**: Creates a new application with provided details.
- **Responses**: 200 (`models.ApplicationResponse`), 400 (Bad Request), 500 (Internal Server Error).

### AuthHandler

Handles authentication-related requests.

#### NewAuthHandler

- **Purpose**: Creates a new instance of `AuthHandler`.
- **Parameters**: `ssiService` (*services.SsiClient), `db` (*store.Store).

#### SignUpHandler

- **Endpoint**: `/signup` (GET)
- **Description**: Handles the sign-up process by providing a redirect URL for authentication.
- **Responses**: 200 (Redirect URL for sign-up), 400 (Bad Request), 500 (Internal Server Error).

#### GetAccessToken

- **Endpoint**: `/get-access-token` (POST)
- **Description**: Handles the sign-in process using application DID and credential JWT.
- **Responses**: 200 (`models.GetAccessTokenResponse`), 400 (Bad Request), 500 (Internal Server Error).

#### RequestAccess

- **Endpoint**: `/request-access` (POST)
- **Description**: Initiates a request for user access (implementation pending).
- **Responses**: 200 (Success), 405 (Method Not Allowed).

#### GrandAccess

- **Endpoint**: `/grant-access` (PUT)
- **Description**: Grants access based on a valid request (implementation pending).
- **Responses**: 200 (Success), 405 (Method Not Allowed).

#### RevokeAccess

- **Endpoint**: `/revoke-access` (PUT)
- **Description**: Revokes the access of a user (implementation pending).
- **Responses**: 200 (Success), 405 (Method Not Allowed).

#### VerifyAccess

- **Endpoint**: `/verify-access` (GET)
- **Description**: Verifies if a user has access to a specific resource based on their role.
- **Responses**: 200 (Success), 400 (Bad Request), 401 (Unauthorized), 500 (Internal Server Error).

#### GetAccessList

- **Endpoint**: `/get-access-list` (GET)
- **Description**: Lists the access for the user on the resource.
- **Responses**: 200 (Success), 400 (Bad Request), 401 (Unauthorized), 500 (Internal Server Error).

### CallbackHandler

#### NewCallbackHandler

- **Purpose**: Creates a new instance of `CallbackHandler`.

#### HandleCallback

- **Endpoint**: Dynamic, based on provider and DID.
- **Method**: GET
- **Description**: Handles the OAuth callback, redirecting to a web page with query parameters including the provider and DID.

#### HandleMe

- **Endpoint**: Dynamic, based on provider and access token.
- **Method**: GET
- **Description**: Retrieves user information based on the provider and access token.
- **Responses**: 200 (User Information), 400 (Bad Request), 500 (Internal Server Error).

### CredentialHandler

#### NewCredentialHandler

- **Purpose**: Creates a new instance of `CredentialHandler`.
- **Parameters**: `ssiService` (*services.SsiClient), `db` (*store.Store).

#### IssueOAuthCredential

- **Endpoint**: Unspecified (handled dynamically)
- **Method**: POST
- **Description**: Issues OAuth credentials based on provided request parameters.
- **Responses**: 200 (Issued Credentials), 400 (Bad Request), 500 (Internal Server Error).

#### RevokeOAuthCredential

- **Endpoint**: `/revoke-credential` (POST)
- **Description**: Revokes an existing OAuth credential.
- **Responses**: 200 (Success Message), 400 (Bad Request), 500 (Internal Server Error).

### MiddlewareService

#### NewMiddlewareService

- **Purpose**: Creates a new instance of `MiddlewareService`.
- **Parameters**: `apikey` (string).

#### EnableCORS

- **Purpose**: Middleware to enable CORS (Cross-Origin Resource Sharing).
- **Description**: Sets CORS headers and handles preflight requests.

#### XApiKeyMiddleware

- **Purpose**: Middleware to validate the `x-api-key` in request headers.
- **Description**: Checks for a valid API key in the request headers.

#### LoggingMiddleware

- **Purpose**: Middleware for logging each request.
- **Description**: Logs the HTTP method and URL path of each request.

#### ChainMiddleware

- **Purpose**: Chains multiple middleware functions.
- **Description**: Allows for easy combination of multiple middleware functions.

### PolicyHandler

#### NewPolicyHandler

- **Purpose**: Creates a new instance of `PolicyHandler`.
- **Parameters**: `ssiService` (*services.SsiClient), `db` (*store.Store).

#### GetPolicyHandler

- **Endpoint**: `/policies` (GET)
- **Description**: Retrieves a list of all policies.
- **Responses**: 200 (Array of `models.PolicySchemaResponse`), 500 (Internal Server Error).

#### CreatePolicyHandler

- **Endpoint**: `/create-policy` (POST)
- **Description**: Creates a new policy based on the provided schema.
- **Responses**: 200 (`models.PolicySchemaResponse`), 400 (Bad Request), 500 (Internal Server Error).

#### AttachPolicyHandler

- **Endpoint**: `/attach-policy` (POST)
- **Description**: Attaches a policy to an application using the provided application and issuer DID, and schema ID.
- **Responses**: 200 (`models.ApplicationPolicyResponse`), 400 (Bad Request), 500 (Internal Server Error).

### AuthProviderHandler

#### NewAuthProviderHandler

- **Purpose**: Creates a new instance of `AuthProviderHandler`.
- **Parameters**: `ssiService` (*services.SsiClient), `db` (*store.Store).

#### GetAuthConnectorHandler

- **Endpoint**: `/auth-provider` (GET)
- **Description**: Retrieves a list of all auth providers.
- **Responses**: 200 (Array of `models.AvailableProvider`), 500 (Internal Server Error).

#### LinkAuthProviderHandler

- **Endpoint**: `/auth-provider/link` (POST)
- **Description**: Links an OAuth provider to an application by its DID.
- **Responses**: 200 (`models.AuthProvider`), 400 (Bad Request), 500 (Internal Server Error).

#### UnLinkAuthProviderHandler

- **Endpoint**: `/auth-provider/unlink` (POST)
- **Description**: Unlinks an authentication provider from an application.
- **Responses**: 200 (Success Message), 400 (Bad Request), 500 (Internal Server Error).

---
