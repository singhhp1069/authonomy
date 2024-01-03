package handlers

import (
	"encoding/json"
	"fmt"
	"identitysphere-api/models"
	"identitysphere-api/pkg/providers"
	"identitysphere-api/pkg/utils"
	"identitysphere-api/services"
	"identitysphere-api/store"
	"net/http"
)

// AuthHandler handles auth-related requests
type AuthHandler struct {
	ssiService *services.SsiClient
	db         *store.Store
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(ssiService *services.SsiClient, db *store.Store) *AuthHandler {
	return &AuthHandler{ssiService: ssiService, db: db}
}

// SignUpHandler godoc
// @Summary Sign up for an application
// @Description Handles the sign-up process by providing a redirect URL for authentication.
// @Tags User Access Management
// @Accept  json
// @Produce  json
// @Param app_did query string true "Application DID"
// @Param session_token query string true "Session Token"
// @Success 200 {object} map[string]string "Redirect URL for sign-up"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /signup [get]
func (h *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract query parameters
	queryParams := r.URL.Query()
	appDid := queryParams.Get("app_did")
	sessionToken := queryParams.Get("session_token")
	fmt.Println("appDid", appDid)
	fmt.Println("queryParams", queryParams)
	auth, err := h.db.GetAuthProvider(appDid)
	if err != nil {
		http.Error(w, "app authentication is not configured yet", http.StatusInternalServerError)
		return
	}

	fmt.Println("authDetails", auth)
	//  session todo validate

	// For demonstration, let's just send back these parameters
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"app_did":       appDid,
		"session_token": sessionToken,
		"redirect_url":  providers.GetLoginUrl(auth.Config.ClientID, auth.Config.RedirectURL),
	}
	json.NewEncoder(w).Encode(response)
}

// SignInHandler godoc
// @Summary Sign in to an application
// @Description Handles the sign-in process using application DID, credential JWT, and session key.
// @Tags User Access Management
// @Accept  json
// @Produce  json
// @Param app_did query string true "Application DID"
// @Param cred_jwt query string true "Credential JWT"
// @Param session_key query string true "Session Key"
// @Success 200 {object} models.SignInResponse "Access Token"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /sign-in [get]
func (h *AuthHandler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// Extract query parameters
	queryParams := r.URL.Query()
	appDid := queryParams.Get("app_did")
	credJwt := queryParams.Get("cred_jwt")
	sessionKey := queryParams.Get("session_key")

	fmt.Println("App DID:", appDid)
	fmt.Println("Credential JWT:", credJwt)
	fmt.Println("Session Key:", sessionKey)

	appDetails, err := h.db.GetApp(appDid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _, cred, err := utils.ParseVerifiableCredentialFromJWT(credJwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cred.Issuer != appDetails.AppDID {
		http.Error(w, "invalid credential", http.StatusInternalServerError)
	}
	// For demonstration, let's just return a success message
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(models.SignInResponse{AccessToken: "xxx"})
}

func (h *AuthHandler) GrandAccess(w http.ResponseWriter, r *http.Request) {
}

func (h *AuthHandler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
}

func (h *AuthHandler) VerifyAccess(w http.ResponseWriter, r *http.Request) {
}
