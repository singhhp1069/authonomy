package handlers

import (
	"encoding/json"
	"fmt"
	"identitysphere-api/models"
	"identitysphere-api/services"
	"identitysphere-api/store"
	"net/http"
)

// AuthProviderHandler handles auth-related requests
type AuthProviderHandler struct {
	ssiService *services.SsiClient
	db         *store.Store
}

// NewAuthProviderHandler creates a new instance of AuthProviderHandler
func NewAuthProviderHandler(ssiService *services.SsiClient, db *store.Store) *AuthProviderHandler {
	return &AuthProviderHandler{ssiService: ssiService, db: db}
}

// @Summary Get available auth providers
// @Description Retrieves a list of all auth providers
// @Tags Authentication Management
// @Accept json
// @Produce json
// @Success 200 {array} models.AvailableProvider
// @Failure 500 {object} string "Internal Server Error"
// @Router /auth-provider [get]
func (h *AuthProviderHandler) GetAuthConnectorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// hardcoded connector can be build to supported various mediums.
	connectors := []models.AvailableProvider{
		{
			ProviderName:     "facebook",
			ProviderType:     "social",
			ProviderProtocol: "oauth2",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connectors)
}

// LinkAuthProviderHandler links an authentication provider to an application
// @Summary Link Authentication Provider
// @Description Links an OAuth provider to an application by its DID
// @Tags Authentication Management
// @Accept json
// @Produce json
// @Param provider body models.AuthProvider true "Authentication Provider Details"
// @Success 200 {object} models.AuthProvider "Successfully linked provider"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth-provider/link [post]
func (h *AuthProviderHandler) LinkAuthProviderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var provider models.AuthProvider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isExist := h.ssiService.IsDIDExists(provider.AppDID)
	if !isExist {
		http.Error(w, "app DID does not exists id: "+provider.AppDID, http.StatusInternalServerError)
		return
	}
	// hardcoded for the hackathon
	if provider.Provider.ProviderName != "facebook" {
		http.Error(w, "currently only facebook as a provider supported ", http.StatusInternalServerError)
		return
	}
	provider.Config.RedirectURL = getCallbackUrl(r, provider.AppDID, provider.Provider.ProviderName)
	// more secure way of storing can be applied (encryption)
	err := h.db.SetAuthProvider(provider)
	if err != nil {
		http.Error(w, "Failed to save provide details "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)
}

// LinkAuthProviderHandler links an authentication provider to an application
// @Summary UnLink Authentication Provider
// @Description Links an OAuth provider to an application by its DID
// @Tags Authentication Management
// @Accept json
// @Produce json
// @Param provider body interface{} true "Authentication Provider Details"
// @Success 200 {object} interface{} "Successfully linked provider"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth-provider/unlink [post]
func (h *AuthProviderHandler) UnLinkAuthProviderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("implementation pending")
}

func getCallbackUrl(r *http.Request, did, provider string) string {
	scheme := "http" // Default to HTTP, check if request is HTTPS if needed
	host := r.Host   // This gets the host from the request header
	url := fmt.Sprintf("%s://%s", scheme, host)
	return fmt.Sprintf("%s/%s/%s/%s", url, "callback", provider, did)
}
