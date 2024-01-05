# Access Token Generation and Validation Process

## Overview

The utility functions for creating and validating access tokens employ JWT (JSON Web Tokens). These tokens are used for secure authentication and authorization within the application.

## Configuration

- **Secret Key**: A secret key used for JWT encryption is obtained from the application's configuration (`service.jwt_encryption_key`).

## CustomClaims Structure

- **Description**: A custom JWT claims structure that extends `jwt.StandardClaims`.
- **Fields**:
  - `AppDID` (string): Application DID.
  - `CredentialJWTs` (`models.IssueOAuthCredential`): OAuth Credentials included in the token.

## Function: CreateAccessToken

### Purpose

Generates a JWT token with custom claims.

### Input

- `appDID` (string): Application DID.
- `credentialJWTs` (`models.IssueOAuthCredential`): Struct of issued OAuth credentials.

### Process

1. Sets the expiration time for the token (24 hours by default).
2. Creates a JWT token with custom claims and the HS256 signing method.
3. Signs the token using the encryption key.

### Output

- **Success**: A string representing the signed JWT token.
- **Failure**: An error if the token signing process fails.

### Example

```go
tokenString, err := CreateAccessToken(appDID, credentialJWTs)
```

## Function: ValidateAccessToken

### Purpose

Validates a given JWT token string and extracts the custom claims.

### Input

- `tokenString` (string): The JWT token string to validate.

### Process

1. Parses the token string with the custom claims structure.
2. Verifies the token using the provided encryption key.
3. Validates the token's authenticity and its expiration.

### Output

- **Success**: The CustomClaims if the token is valid.
- **Failure**: An error if the token is invalid or parsing fails.

```go
claims, err := ValidateAccessToken(tokenString)
```

## Security Considerations

- The JWT secret key should be securely stored and managed.
- Tokens have a set expiration time and should be refreshed as needed.
- Proper error handling is crucial for ensuring security and correct access control.
