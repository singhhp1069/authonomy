package handlers

import (
	"authonomy/models"
	"authonomy/pkg/providers"
	"authonomy/services"
	"authonomy/store"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

// CredentialHandler handles credential-related requests
type CredentialHandler struct {
	ssiService *services.SsiClient
	db         *store.Store
}

// NewCredentialHandler creates a new instance of CredentialHandler
func NewCredentialHandler(ssiService *services.SsiClient, db *store.Store) *CredentialHandler {
	return &CredentialHandler{ssiService: ssiService, db: db}
}

func (h *CredentialHandler) IssueOAuthCredential(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var validate = validator.New()
	var credReq models.IssueOAuthCredentialRequest

	err := json.NewDecoder(r.Body).Decode(&credReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(credReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	schema, err := h.db.GetProviderSchema("facebook")
	if err != nil {
		http.Error(w, "Failed to get schema ID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	isSchemaExist := h.ssiService.IsSchemaExists(schema.SchemaID)
	if !isSchemaExist {
		http.Error(w, "Failed to get schema ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	app, err := h.db.GetApp(credReq.AppDID)
	if err != nil {
		http.Error(w, "Failed to get application DID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	policy, err := h.db.GetIssuedPolicy(app.AppDID)
	if err != nil {
		http.Error(w, "Failed to get application policy: "+err.Error(), http.StatusInternalServerError)
		return
	}

	isDIDExits := h.ssiService.IsDIDExists(credReq.AppDID)
	if !isDIDExits {
		http.Error(w, "Failed to get application DID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	userInfo, err := providers.GetUserInfo("facebook", credReq.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userCredMap, err := models.StructToMap(models.UserInfo{UserID: userInfo.ID, Name: userInfo.Name})
	if err != nil {
		http.Error(w, "Failed to convert to map: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// TODO:: revokable
	userCredential, err := h.ssiService.IssueCredentialBySchemaID(app.AppDID, credReq.UserDID, schema.SchemaID, userCredMap)
	if err != nil {
		http.Error(w, "Failed to issue userCredential: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// default role as user (hardcoded for the hackathon demo)
	userRole := models.Role{
		RoleName:    "user",
		Permissions: []string{"view_content", "comment"},
	}
	policyCredMap, err := models.StructToMap(models.RolesWrapper{Roles: []models.Role{userRole}})
	if err != nil {
		http.Error(w, "Failed to convert to map: "+err.Error(), http.StatusInternalServerError)
		return
	}
	policyCredential, err := h.ssiService.IssueCredentialBySchemaID(app.AppDID, credReq.UserDID, policy.SchemaID, policyCredMap)
	if err != nil {
		http.Error(w, "Failed to issue policyCredential: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.IssueOAuthCredential{OAuthCredential: userCredential, PolicyCredential: policyCredential})
}

// RevokeOAuthCredential godoc
// @Summary Revoke OAuth Credential
// @Description Revoke an existing OAuth credential.
// @Tags Authentication Management
// @Accept  json
// @Produce  json
// @Param revokeOAuthCredentialRequest body interface{} true "Revoke Credential Request"
// @Success 200 {string} string "Credential successfully revoked"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /revoke-credential [post]
func (h *CredentialHandler) RevokeOAuthCredential(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("implementation pending")
}
