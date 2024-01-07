# Architecture

## Overview

This document provides an architectural overview of Authonomy, an innovative middleware designed to revolutionize authorization management for application developers. Authonomy separates authentication from authorization, leveraging decentralized authorization through DIDs, VCs, and VPs.

## System Architecture

### Overall System Architecture

The following diagram illustrates the high-level architecture of Authonomy, showing how it integrates with client applications, authentication providers, and handles authorization management.

```mermaid
graph LR
A[Client Application] --> B(Authonomy Middleware)
B --> C{Authentication Providers}
B --> D{Authorization Management}
D --> E[DIDs]
D --> F[VCs/VPs]
C --> G{type: social}
C --> H{type: passwordless}
G --> J[facebook:OAuth2, google: OpenID, ...]
H --> K[email:SMTP, phoneOTP: verification, ...]
F --> I[Access Control Policies e.g.; RBAC, ABAC, KYC etc]
```

## Authentication and Authorization Flow

This sequence diagram details the flow of authentication and authorization within the Authonomy system, from the user's request to the final access decision.

```mermaid
sequenceDiagram
participant U as User
participant A as Client Application
participant M as Authonomy Middleware
participant P as Authentication Provider
participant D as Decentralized Authorization

U->>A: Request Access
A->>M: Initiate Authentication
M->>P: Use Authentication Provider
P-->>M: Authentication Response
M->>D: Request Authorization (DID, VCs/VPs)
D-->>M: Authorization Response
M-->>A: Issue Policy Credential
```

## Integration in Application Development

The diagram below shows how Authonomy is integrated into the application development process, outlining its interaction with third-party authentication providers and decentralized authorization mechanisms.

```mermaid
graph LR
A[Client Application] --> B{Authonomy Integration}
B -->|Authentication| C[Third-party Auth Providers]
B -->|Authorization| D[Decentralized Authorization via Authonomy]
C --> E[OAuth2, OpenID, etc.]
D --> F[DID/VC Based Access Control]
```

## Conclusion

Authonomy introduces a novel approach to digital identity management by separating authentication and authorization. Its architecture leverages modern decentralized technologies to provide a flexible, secure, and user-centric solution for application developers.
