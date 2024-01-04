package models

import (
	"encoding/json"

	credsdk "github.com/TBD54566975/ssi-sdk/credential"
	didsdk "github.com/TBD54566975/ssi-sdk/did"
)

type ApplicationRequest struct {
	AppName    string     `json:"app_name" validate:"required,min=3,max=100"`
	AppDetails AppDetails `json:"app_details" validate:"required,dive"`
}

type AppDetails struct {
	Description  string `json:"description" validate:"required,min=10,max=500"`
	ContactEmail string `json:"email" validate:"required,email"`
}

type ApplicationResponse struct {
	AppDID     string     `json:"app_did"`
	AppName    string     `json:"app_name"`
	AppDetails AppDetails `json:"app_details"`
}

type DidCreationResponse struct {
	Did didsdk.Document `json:"did,omitempty"`
}

type PolicySchemaResponse struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Schema JsonSchema `json:"schema"`
}

type PolicySchemaRequest struct {
	Name   string     `json:"name"`
	Schema JsonSchema `json:"schema"`
}

type JsonSchema struct {
	DollarSchema string                 `json:"$schema"`
	Type         string                 `json:"type"`
	Properties   map[string]interface{} `json:"properties"`
	Required     []string               `json:"required,omitempty"`
}

type ApplicationPolicyRequest struct {
	ApplicationDID string                 `json:"application_did" validate:"required"`
	SchemaID       string                 `json:"schema_id" validate:"required"`
	IssuerDID      string                 `json:"issuer_did" validate:"required"`
	Credential     map[string]interface{} `json:"credential" validate:"required"`
}
type ApplicationPolicyResponse struct {
	ApplicationDID string `json:"application_did"`
	SchemaID       string `json:"schema_id"`
	IssuerDID      string `json:"issuer_did"`
	CredentialID   string `json:"credential_id"`
}

type CredentialRequest struct {
	Issuer               string                 `json:"issuer"`
	VerificationMethodID string                 `json:"verificationMethodId"`
	Subject              string                 `json:"subject"`
	SchemaID             string                 `json:"schemaId"`
	Data                 map[string]interface{} `json:"data"`
	// TODO:: add revoked
}

type CredentialResponse struct {
	ID                                 string                        `json:"id"`
	FullyQualifiedVerificationMethodID string                        `json:"fullyQualifiedVerificationMethodId"`
	Credential                         *credsdk.VerifiableCredential `json:"credential,omitempty"`
	CredentialJwt                      string                        `json:"credentialJwt"`
}

type AuthProvider struct {
	AppDID   string            `json:"app_did" validate:"required"`
	Provider AvailableProvider `json:"app_details" validate:"required,dive"`
	Config   OAuthConfig       `json:"config" validate:"required,dive"`
}

type AvailableProvider struct {
	ProviderName     string `json:"provider_name" validate:"required"`
	ProviderType     string `json:"provider_type" validate:"required"` // social, email, phone, etc
	ProviderProtocol string `json:"provider_protocol" validate:"required"`
}

// OAuthConfig holds the configuration for OAuth authentication
type OAuthConfig struct {
	ClientID    string `json:"client_id"`
	RedirectURL string `json:"redirect_url"`
	// ClientSecret string   `json:"clientSecret"`
	// Scopes       []string `json:"scopes"`
	// AuthURL      string   `json:"authUrl"`
	// TokenURL     string   `json:"tokenUrl"`
}

type ProviderSchema struct {
	ProviderName string `json:"provider_name" `
	SchemaID     string `json:"schema_id" `
}

type IssueOAuthCredentialRequest struct {
	AppDID      string `json:"app_did" validate:"required"`
	Provider    string `json:"provider" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
	UserDID     string `json:"user_did" validate:"required"`
}

type IssueOAuthCredentialResponse struct {
	OAuthCredential  interface{} `json:"oauth_credential" validate:"required"`
	PolicyCredential interface{} `json:"policy_credential" validate:"required"`
}

type UserInfo struct {
	UserID string `json:"user_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email"`
	// more fields soon
}

type GetAccessTokenRequest struct {
	AppDID           string      `json:"app_did"`
	OAuthCredential  interface{} `json:"oauth_credential" validate:"required"`
	PolicyCredential interface{} `json:"policy_credential" validate:"required"`
}

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type VerifyAccessRequest struct {
	CredentialJWT string `json:"credential_jwt"`
}

type VerifyAccessResponse struct {
	Status bool `json:"status"`
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dataBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Define structs to match the JSON structure
type RBAC struct {
	ID    string `json:"id"`
	Roles []Role `json:"roles"`
}

type Role struct {
	RoleName    string   `json:"roleName"`
	Permissions []string `json:"permissions"`
}

type RolesWrapper struct {
	Roles []Role `json:"roles"`
}
