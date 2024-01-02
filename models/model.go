package models

import (
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
	ID                                 string `json:"id"`
	FullyQualifiedVerificationMethodID string `json:"fullyQualifiedVerificationMethodId"`
	Credential                         struct {
		Context           []string               `json:"@context"`
		ID                string                 `json:"id"`
		Type              []string               `json:"type"`
		Issuer            string                 `json:"issuer"`
		IssuanceDate      string                 `json:"issuanceDate"`
		CredentialSubject map[string]interface{} `json:"credentialSubject"`
		CredentialSchema  struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"credentialSchema"`
	} `json:"credential"`
	CredentialJwt string `json:"credentialJwt"`
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

type Config struct {
	App        *ApplicationResponse
	Auth       *AuthProvider
	Cred       *ApplicationPolicyResponse
	SessionKey string
}

type ProviderSchema struct {
	ProviderName string `json:"provider_name" `
	SchemaID     string `json:"schema_id" `
}

type IssueOAuthCredentialRequest struct {
	AppDID         string   `json:"app_did" validate:"required"`
	Provider       string   `json:"provider" validate:"required"`
	UserInfo       UserInfo `json:"user_info,omitempty"`
	CredentialType string   `json:"credential_type" validate:"required"`
	UserDID        string   `json:"user_did" validate:"required"`
}

type IssueOAuthCredentialResponse struct {
	ID                                 string                       `json:"id" binding:"required"`
	Credential                         credsdk.VerifiableCredential `json:"credential"`
	DID                                didsdk.Document              `json:"did"`
	CredentailJWT                      string                       `json:"credentialJwt" binding:"required"`
	FullyQualifiedVerificationMethodId string                       `json:"fullyQualifiedVerificationMethodId" binding:"required"`
}

type UserInfo struct {
	UserID string `json:"user_id" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email"`
	// more fields soon
}
