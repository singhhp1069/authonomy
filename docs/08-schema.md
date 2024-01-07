# Supported Policy Schemas and Authentication Methods

## Policy Schemas

### 1. RBAC Policy (Role-Based Access Control)

#### Description

Defines access control based on roles and their associated permissions.

#### Schema

- **Type**: Object
- **Properties**:
  - `credentialSubject` (object):
    - `roles` (array): An array of objects defining roles and permissions.
      - `roleName` (string): Name of the role.
      - `permissions` (array): List of permissions as strings.

### 2. ABAC Policy (Attribute-Based Access Control)

#### Description

Specifies access control based on attributes associated with subjects, actions, and resources.

#### Schema

- **Type**: Object
- **Properties**:
  - `credentialSubject` (object):
    - `rules` (array): An array of objects defining rules for access control.
      - `ruleId` (string): Identifier for the rule.
      - `subjectAttributes`, `actionAttributes`, `resourceAttributes` (object): Attributes for subject, action, and resource.
      - `effect` (string): Specifies the effect ('permit' or 'deny').

## Authentication Methods

### Supported

1. **OAuth2**
   - **Type**: Social
   - **Provider**: Facebook
   - **Schema**:
     - **Type**: Object
     - **Properties**:
       - `credentialSubject` (object):
         - `user_id` (string): User identifier.
         - `name` (string): User's name.
         - `email` (string, optional): User's email.

### In Development

- OAuth2, OpenID, Password-based, Password-less, FIDO, AuthN.

## Remarks

- The RBAC and ABAC policies are designed to provide flexible and robust access control mechanisms.
- The OAuth2 schema for Facebook is currently supported, focusing on essential user attributes.
- Future development includes expanding authentication methods to cover a broader range of protocols and standards.

---

These policies and authentication methods are integral to the security and access management of the system, ensuring that only authorized entities have access to specific resources based on defined criteria.
