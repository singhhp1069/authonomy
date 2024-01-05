package services

import (
	"authonomy/models"
	"authonomy/store"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SsiClient is the client for interacting with the SSI service
type SsiClient struct {
	serviceUrl string
	httpClient *http.Client
}

// NewSsiClient creates a new instance of SsiClient
func NewSsiClient(url string) *SsiClient {
	if url == "" {
		panic("ssi service url not set")
	}
	return &SsiClient{
		serviceUrl: url,
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
	policy, err := client.createPolicyFromFile(filepath.Join(path, "ssi", "schemas", "rbac.json"), db)
	if err != nil {
		return fmt.Errorf("error creating policy from file %v", err)
	}
	db.SetPolicy(policy)
	// set ABAC policy
	policy, err = client.createPolicyFromFile(filepath.Join(path, "ssi", "schemas", "abac.json"), db)
	if err != nil {
		return fmt.Errorf("error creating policy from file %v", err)
	}
	db.SetPolicy(policy)

	// set OAuth2 user info schema
	policy, err = client.createPolicyFromFile(filepath.Join(path, "ssi", "schemas", "oauth_info.json"), db)
	if err != nil {
		return fmt.Errorf("error creating policy from file %v", err)
	}
	// only facebook is supported yet for the demo
	return db.SetProviderSchema(models.ProviderSchema{ProviderName: "facebook", SchemaID: policy.ID})
}

// createPolicyFromFile create a policy from a json file path.
func (client *SsiClient) createPolicyFromFile(filePath string, db *store.Store) (schemaRep models.PolicySchemaResponse, err error) {
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

	req, err := http.NewRequest("PUT", client.serviceUrl+"/dids/key", bytes.NewBuffer(requestBody))
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
	req, err := http.NewRequest("PUT", client.serviceUrl+"/schemas", bytes.NewBuffer(schemaBody))
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

	req, err := http.NewRequest("PUT", client.serviceUrl+"/credentials", bytes.NewBuffer(requestBody))
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

// // VerifyCredential verifies a credential JWT and returns the verification result
// func (client *SsiClient) VerifyCredential(credentialJWT string) (models.VerificationResponse, error) {
// 	// Prepare the request body
// 	requestBody, err := json.Marshal(map[string]string{"credentialJwt": credentialJWT})
// 	if err != nil {
// 		return models.VerificationResponse{}, err
// 	}

// 	// Create and send the request
// 	req, err := http.NewRequest("PUT", SSI_BASE_URL+"/credentials/verification", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		return models.VerificationResponse{}, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.httpClient.Do(req)
// 	if err != nil {
// 		return models.VerificationResponse{}, err
// 	}
// 	defer resp.Body.Close()

// 	// Decode the response
// 	var verificationResp models.VerificationResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&verificationResp); err != nil {
// 		return models.VerificationResponse{}, err
// 	}

// 	return verificationResp, nil
// }

// IsSchemaExists checks if a schema exists
func (client *SsiClient) IsSchemaExists(schema string) bool {
	resp, err := client.httpClient.Get(client.serviceUrl + "/schemas/" + schema)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// IsDIDExists checks if a DID exists
func (client *SsiClient) IsDIDExists(did string) bool {
	resp, err := client.httpClient.Get(client.serviceUrl + "/dids/key/" + did)
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
