# Service Models

## Overview

The `models` package defines the data structures used in the Authonomy service. These models are used for managing and interacting with applications, policies, credentials, authentication providers, and user information.

## Models

### ApplicationRequest

Represents a request to create a new application.

- `AppName`: Name of the application.
- `AppDetails`: Details of the application including description and contact email.

### AppDetails

Contains details about an application.

- `Description`: Description of the application.
- `ContactEmail`: Contact email address.

### ApplicationResponse

Response structure for application creation or query.

- `AppDID`: Application DID.
- `AppSceret`: Application secret.
- `AppName`: Application name.
- `AppDetails`: Application details.

### DidCreationResponse

Response for DID creation.

- `Did`: DID document.

### PolicySchemaResponse

Response structure for policy schema queries.

- `ID`: Schema ID.
- `Name`: Schema name.
- `Schema`: JSON schema.

### PolicySchemaRequest

Request structure for creating a policy schema.

- `Name`: Schema name.
- `Schema`: JSON schema.

### JsonSchema

Defines a JSON schema.

- `DollarSchema`: Schema definition.
- `Type`: Type of the schema.
- `Properties`: Schema properties.
- `Required`: Required fields.

### ApplicationPolicyRequest

Request for creating an application policy.

- `ApplicationDID`: Application DID.
- `SchemaID`: Schema ID.
- `IssuerDID`: Issuer DID.
- `Credential`: Credential data.

### ApplicationPolicyResponse

Response for application policy queries.

- `ApplicationDID`: Application DID.
- `SchemaID`: Schema ID.
- `IssuerDID`: Issuer DID.
- `CredentialID`: Credential ID.
- `CredentialSubject`: Credential subject.

### CredentialRequest

Request for creating a credential.

- `Issuer`: Issuer DID.
- `VerificationMethodID`: Verification method ID.
- `Subject`: Subject DID.
- `SchemaID`: Schema ID.
- `Data`: Credential data.

### CredentialResponse

Response for credential queries.

- `ID`: Credential ID.
- `FullyQualifiedVerificationMethodID`: Verification method ID.
- `Credential`: Verifiable credential.
- `CredentialJwt`: Credential JWT.

### AuthProvider

Represents an authentication provider.

- `AppDID`: Application DID.
- `Provider`: Available provider.
- `Config`: OAuth configuration.

### AvailableProvider

Details of an available authentication provider.

- `ProviderName`: Provider name.
- `ProviderType`: Provider type.
- `ProviderProtocol`: Provider protocol.
- `ProviderSchemaID`: Schema ID.

### OAuthConfig

Configuration for OAuth authentication.

- `ClientID`: Client ID.
- `RedirectURL`: Redirect URL.

### ProviderSchema

Schema for an authentication provider.

- `ProviderName`: Provider name.
- `SchemaID`: Schema ID.

### IssueOAuthCredentialRequest

Request for issuing OAuth credentials.

- `AppDID`: Application DID.
- `Provider`: Provider name.
- `AccessToken`: Access token.
- `UserDID`: User DID.

### IssueOAuthCredential

Represents OAuth and policy credentials.

- `OAuthCredential`: OAuth credential.
- `PolicyCredential`: Policy credential.

### UserInfo

User information.

- `UserID`: User ID.
- `Name`: User name.
- `Email`: User email.

### GetAccessTokenResponse

Response for access token queries.

- `AccessToken`: Access token.

### VerifyAccessRequest

Request for verifying access.

- `CredentialJWT`: Credential JWT.

### VerifyAccessResponse

Response for access verification.

- `Status`: Verification status.

### RBAC

Role-based access control structure.

- `ID`: ID.
- `Roles`: Roles.

### Role

Role structure.

- `RoleName`: Role name.
- `Permissions`: Permissions.

### RolesWrapper

Wrapper for roles.

- `Roles`: Roles.

### AccessList

Structure for access lists.

- `ApplicationPolicy`: Application policy.
- `UserAccessList`: User access list.

### Helper Functions

#### StructToMap

Converts a struct to a map.

- `data`: Struct to convert.
- `return`: Map and error.

## Usage

These models are used throughout the Authonomy service for handling requests and responses in various operations like managing applications, policies, and credentials.
