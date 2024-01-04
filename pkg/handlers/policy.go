package handlers

import (
	"authonomy/models"
	"authonomy/services"
	"authonomy/store"
	"encoding/json"
	"net/http"
)

// PolicyHandler handles policy-related requests
type PolicyHandler struct {
	ssiService *services.SsiClient
	db         *store.Store
}

// NewPolicyHandler creates a new instance of PolicyHandler
func NewPolicyHandler(ssiService *services.SsiClient, db *store.Store) *PolicyHandler {
	return &PolicyHandler{ssiService: ssiService, db: db}
}

// GetPolicyHandler retrieves all policies
// @Summary Get all policies
// @Description Retrieves a list of all policies
// @Tags Authorization Management
// @Accept json
// @Produce json
// @Param x-api-key header string true "API Key"
// @Success 200 {array} models.PolicySchemaResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /policies [get]
func (h *PolicyHandler) GetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	policies, err := h.db.GetAllPolicies()
	if err != nil {
		http.Error(w, "Failed to get policy: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policies)
}

// CreatePolicyHandler creates a new policy
// @Summary Create a new policy
// @Description Creates a new policy based on the provided schema
// @Tags Authorization Management
// @Accept json
// @Produce json
// @Param x-api-key header string true "API Key"
// @Param schema body models.PolicySchemaRequest true "Policy Schema"
// @Success 200 {object} models.PolicySchemaResponse "Successfully created policy"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /create-policy [post]
func (h *PolicyHandler) CreatePolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var schema models.PolicySchemaRequest
	err := json.NewDecoder(r.Body).Decode(&schema)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respSchema, err := h.ssiService.CreatePolicy(schema)
	if err != nil {
		http.Error(w, "Failed to create policy: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.db.SetPolicy(respSchema)
	if err != nil {
		http.Error(w, "Failed to store policy: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respSchema)
}

// AttachPolicyHandler attaches a policy to an application
// @Summary Attach policy to application
// @Description Attaches a policy to an application using the provided application and issuer DID, and schema ID
// @Tags Authorization Management
// @Accept json
// @Produce json
// @Param x-api-key header string true "API Key"
// @Param appPolicy body models.ApplicationPolicyRequest true "Application Policy Request"
// @Success 200 {object} models.ApplicationPolicyResponse "Successfully attached policy"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /attach-policy [post]
func (h *PolicyHandler) AttachPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var appPolicy models.ApplicationPolicyRequest
	err := json.NewDecoder(r.Body).Decode(&appPolicy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isExist := h.ssiService.IsSchemaExists(appPolicy.SchemaID)
	if !isExist {
		http.Error(w, "Schema does not exists id: "+appPolicy.SchemaID, http.StatusInternalServerError)
		return
	}
	isExist = h.ssiService.IsDIDExists(appPolicy.ApplicationDID)
	if !isExist {
		http.Error(w, "app DID does not exists id: "+appPolicy.ApplicationDID, http.StatusInternalServerError)
		return
	}
	isExist = h.ssiService.IsDIDExists(appPolicy.IssuerDID)
	if !isExist {
		http.Error(w, "issuer DID does not exists id: "+appPolicy.IssuerDID, http.StatusInternalServerError)
		return
	}

	respSchema, err := h.ssiService.IssueCredentialBySchemaID(appPolicy.IssuerDID, appPolicy.ApplicationDID, appPolicy.SchemaID, appPolicy.Credential)
	if err != nil {
		http.Error(w, "Failed to issue credential: "+err.Error(), http.StatusInternalServerError)
		return
	}
	policyResponse := models.ApplicationPolicyResponse{
		ApplicationDID: appPolicy.ApplicationDID,
		SchemaID:       appPolicy.SchemaID,
		IssuerDID:      appPolicy.IssuerDID,
		CredentialID:   respSchema.ID,
	}
	// Save the new application details to BadgerDB
	err = h.db.SetIssuedPolicy(policyResponse)
	if err != nil {
		http.Error(w, "Failed to save application policy: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policyResponse)
}
