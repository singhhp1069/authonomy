package utils

import (
	"fmt"

	"github.com/TBD54566975/ssi-sdk/credential"
)

// IsRoleExists check the role exists on RBAC or not.
func IsRoleExists(cred credential.CredentialSubject, r string) bool {
	// Extract the roles
	roles, ok := cred["roles"].([]any)
	if !ok {
		fmt.Println("roles not found or not in expected format")
		return false
	}

	for _, roleInterface := range roles {
		role, ok := roleInterface.(map[string]any)
		if !ok {
			fmt.Println("role is not in expected format")
			continue
		}

		roleName, ok := role["roleName"].(string)
		if ok {
			return (roleName == r)
		}
	}
	return false
}
