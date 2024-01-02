package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"identitysphere-api/models"
	"identitysphere-api/store"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	SSI_BASE_URL = "http://localhost:8080/v1" // Replace with actual SSI service URL
)

// SsiClient is the client for interacting with the SSI service
type SsiClient struct {
	httpClient *http.Client
}

// NewSsiClient creates a new instance of SsiClient
func NewSsiClient() *SsiClient {
	return &SsiClient{
		httpClient: &http.Client{},
	}
}

// CreateDemoPolicies to create a demo policy for the demo.
func CreateDemoPolicies(client *SsiClient, db *store.Store) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	// set RBAC policy
	policy, err := client.createPolicyFromFile(filepath.Join(path, "services", "sample_schema", "rbac.json"), db)
	if err != nil {
		return fmt.Errorf("error creating policy from file %v", err)
	}
	db.SetPolicy(policy)
	// set ABAC policy
	policy, err = client.createPolicyFromFile(filepath.Join(path, "services", "sample_schema", "abac.json"), db)
	if err != nil {
		return fmt.Errorf("error creating policy from file %v", err)
	}
	db.SetPolicy(policy)

	// set OAuth2 user info schema
	policy, err = client.createPolicyFromFile(filepath.Join(path, "services", "sample_schema", "oauth_info.json.json"), db)
	if err != nil {
		return fmt.Errorf("error creating policy from file %v", err)
	}
	// only facebook is supported yet for the demo
	db.SetProviderSchema(models.ProviderSchema{ProviderName: "facebook", SchemaID: policy.ID})
	return nil
}

func (client *SsiClient) createPolicyFromFile(filePath string, db *store.Store) (schemaRep models.PolicySchemaResponse, err error) {
	fmt.Println("filePath", filePath)
	// Read the JSON file
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return schemaRep, err
	}
	// Unmarshal JSON into PolicySchemaRequest struct
	var schema models.PolicySchemaRequest
	err = json.Unmarshal(fileBytes, &schema)
	if err != nil {
		return schemaRep, err
	}
	// Create the policy
	policy, err := client.CreatePolicy(schema)
	if err != nil {
		return schemaRep, err
	}
	// Store policies
	return policy, nil
}

// CreateDid creates a new DID and returns its identifier
func (client *SsiClient) CreateDid() (string, error) {
	requestBody, err := json.Marshal(map[string]string{"keyType": "Ed25519"})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("PUT", SSI_BASE_URL+"/dids/key", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var didResp models.DidCreationResponse
	if err := json.NewDecoder(resp.Body).Decode(&didResp); err != nil {
		return "", err
	}

	return didResp.Did.ID, nil
}

// CreatePolicy creates a new policy and returns its response
func (client *SsiClient) CreatePolicy(schema models.PolicySchemaRequest) (policy models.PolicySchemaResponse, err error) {
	schemaBody, err := json.Marshal(schema)
	if err != nil {
		return policy, err
	}

	req, err := http.NewRequest("PUT", SSI_BASE_URL+"/schemas", bytes.NewBuffer(schemaBody))
	if err != nil {
		return policy, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return policy, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return policy, fmt.Errorf("failed to create schema, status code: %d", resp.StatusCode)
	}

	var policyResp models.PolicySchemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&policyResp); err != nil {
		return policy, err
	}

	return policyResp, nil
}

// IssueCredentialBySchemaID issues a credential based on a schema ID
func (client *SsiClient) IssueCredentialBySchemaID(issuer, subject, schemaID string, data map[string]interface{}) (cred models.CredentialResponse, err error) {
	credRequest := models.CredentialRequest{
		Issuer:               issuer,
		VerificationMethodID: prepareVerificationMethod(issuer),
		Subject:              subject,
		SchemaID:             schemaID,
		Data:                 data,
	}

	requestBody, err := json.Marshal(credRequest)
	if err != nil {
		return cred, err
	}

	req, err := http.NewRequest("PUT", SSI_BASE_URL+"/credentials", bytes.NewBuffer(requestBody))
	if err != nil {
		return cred, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return cred, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return cred, fmt.Errorf("failed to issue credential, status code: %d", resp.StatusCode)
	}

	var credentialResp models.CredentialResponse
	if err := json.NewDecoder(resp.Body).Decode(&credentialResp); err != nil {
		return cred, err
	}

	return credentialResp, nil
}

// IsSchemaExists checks if a schema exists
func (client *SsiClient) IsSchemaExists(schema string) bool {
	resp, err := client.httpClient.Get(SSI_BASE_URL + "/schemas/" + schema)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// IsDIDExists checks if a DID exists
func (client *SsiClient) IsDIDExists(did string) bool {
	resp, err := client.httpClient.Get(SSI_BASE_URL + "/dids/key/" + did)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// prepareVerificationMethod formats the verification method for the credential request
func prepareVerificationMethod(appDID string) string {
	segments := strings.Split(appDID, ":")
	if len(segments) < 3 {
		return "" // Error handling can be improved here
	}

	lastSegment := segments[len(segments)-1]
	return fmt.Sprintf("%s#%s", appDID, lastSegment)
}
