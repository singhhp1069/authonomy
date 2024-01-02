package handlers

import (
	"encoding/json"
	"identitysphere-api/models"
	"identitysphere-api/services"
	"identitysphere-api/store"
	"net/http"

	"github.com/fatih/structs"
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
	}
	app, err := h.db.GetApp(credReq.AppDID)
	if err != nil {
		http.Error(w, "Failed to get application DID: "+err.Error(), http.StatusInternalServerError)
	}
	isDIDExits := h.ssiService.IsDIDExists(credReq.AppDID)
	if !isDIDExits {
		http.Error(w, "Failed to get application DID: "+err.Error(), http.StatusInternalServerError)
	}

	credential, err := h.ssiService.IssueCredentialBySchemaID(app.AppDID, credReq.UserDID, schema.SchemaID, structs.Map(credReq.UserInfo))
	if err != nil {
		http.Error(w, "Failed to issue credential: "+err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(credential)
}

func (h *CredentialHandler) RevokeOAuthCredential(w http.ResponseWriter, r *http.Request) {
}
