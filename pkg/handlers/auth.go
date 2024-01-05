package handlers

import (
	"authonomy/models"
	"authonomy/pkg/providers"
	"authonomy/pkg/utils"
	"authonomy/services"
	"authonomy/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
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
// @Param app_secret query string true "Application secret"
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
	appSecret := queryParams.Get("app_secret")
	fmt.Println("appDid", appDid)
	fmt.Println("queryParams", queryParams)
	appDetails, err := h.db.GetApp(appDid)
	if err != nil {
		http.Error(w, "app is invalid", http.StatusInternalServerError)
		return
	}
	if appDetails.AppSceret != appSecret {
		http.Error(w, "app secret is invalid", http.StatusInternalServerError)
		return
	}
	auth, err := h.db.GetAuthProvider(appDid)
	if err != nil {
		http.Error(w, "app authentication is not configured yet", http.StatusInternalServerError)
		return
	}
	// For demonstration, let's just send back these parameters
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"app_did":      appDid,
		"redirect_url": providers.GetLoginUrl(auth.Config.ClientID, auth.Config.RedirectURL),
	}
	json.NewEncoder(w).Encode(response)
}

// GetAccessToken godoc
// @Summary Sign in or get access token to an application
// @Description Handles the sign-in process using application DID, credential JWT.
// @Tags User Access Management
// @Accept json
// @Produce json
// @Param app_did query string true "Application DID"
// @Param app_secret query string true "Application secret"
// @Param application body models.IssueOAuthCredential true "Application to create"
// @Success 200 {object} models.GetAccessTokenResponse "Access Token"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /get-access-token [post]
func (h *AuthHandler) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract query parameters
	queryParams := r.URL.Query()
	appDid := queryParams.Get("app_did")
	appSecret := queryParams.Get("app_secret")

	appDetails, err := h.db.GetApp(appDid)
	if err != nil {
		http.Error(w, "app is invalid", http.StatusInternalServerError)
		return
	}
	if appDetails.AppSceret != appSecret {
		http.Error(w, "app secret is invalid", http.StatusInternalServerError)
		return
	}

	var validate = validator.New()
	var appReq models.IssueOAuthCredential

	err = json.NewDecoder(r.Body).Decode(&appReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(appReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _, oauthCred, err := utils.ParseVerifiableCredentialFromJWT(appReq.OAuthCredential.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if oauthCred.Issuer != appDid {
		http.Error(w, "incorrect oauth cred", http.StatusBadRequest)
		return
	}

	_, _, policyCred, err := utils.ParseVerifiableCredentialFromJWT(appReq.PolicyCredential.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if policyCred.Issuer != appDid {
		http.Error(w, "incorrect policy cred", http.StatusBadRequest)
		return
	}

	accessToken, err := utils.CreateAccessToken(appDid, models.IssueOAuthCredential{
		OAuthCredential:  appReq.OAuthCredential,
		PolicyCredential: appReq.PolicyCredential,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// For demonstration, let's just return a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.GetAccessTokenResponse{AccessToken: accessToken})
}

// RequestAccess godoc
// @Summary Request access for a user
// @Description Initiates a request for user access.
// @Tags User Access Management
// @Accept  json
// @Produce  json
// @Param app_did query string true "Application DID"
// @Param app_secret query string true "Application secret"
// @Success 200 {string} string "implementation pending"
// @Failure 405 {string} string "Only POST method is allowed"
// @Router /request-access [post]
func (h *AuthHandler) RequestAccess(w http.ResponseWriter, r *http.Request) {
	// request id based tracking and issurance
	// need to better p2p storage?
	// DIDComm exchange (data persistance?)
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("implementation pending")
}

// GrandAccess godoc
// @Summary Grant access to a user
// @Description Grants access based on a valid request.
// @Tags Permission Management
// @Accept  json
// @Produce  json
// @Param app_did query string true "Application DID"
// @Param app_secret query string true "Application secret"
// @Success 200 {string} string "implementation pending"
// @Failure 405 {string} string "Only PUT method is allowed"
// @Router /grant-access [put]
func (h *AuthHandler) GrandAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("implementation pending")
}

// RevokeAccess godoc
// @Summary Revoke access of a user
// @Description Revokes the access of a user.
// @Tags Permission Management
// @Accept  json
// @Produce  json
// @Param app_did query string true "Application DID"
// @Param app_secret query string true "Application secret"
// @Success 200 {string} string "implementation pending"
// @Failure 405 {string} string "Only PUT method is allowed"
// @Router /revoke-access [put]
func (h *AuthHandler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("implementation pending")
}

// VerifyAccess godoc
// @Summary Verify access to a resource
// @Description Verifies if a user has access to a specific resource based on their role.
// @Tags User Access Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer YOUR_ACCESS_TOKEN)
// @Param attribute query string false "e.g.; Role to check access for"
// @Success 200 {string} string "success"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /verify-access [get]
func (h *AuthHandler) VerifyAccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract query parameters
	queryParams := r.URL.Query()
	appDid := queryParams.Get("app_did")
	appSecret := queryParams.Get("app_secret")

	appDetails, err := h.db.GetApp(appDid)
	if err != nil {
		http.Error(w, "app is invalid", http.StatusInternalServerError)
		return
	}
	if appDetails.AppSceret != appSecret {
		http.Error(w, "app secret is invalid", http.StatusInternalServerError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: No Authorization header provided", http.StatusUnauthorized)
		return
	}
	// TODO:: hardcoded, should be based on policy schema ID, and application credential existence and user ownership
	// the can be VP too
	role := queryParams.Get("attribute") // warning: will work with role only
	fmt.Println("role", role)
	// Split the header to get the token part
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		http.Error(w, "Unauthorized: Invalid Authorization header format", http.StatusUnauthorized)
		return
	}
	// headerParts[1] contains the actual token
	token := headerParts[1]

	claims, err := utils.ValidateAccessToken(token)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid access token", http.StatusUnauthorized)
		return
	}
	// TODO:: revokable check
	_, _, oauthCred, err := utils.ParseVerifiableCredentialFromJWT(claims.CredentialJWTs.OAuthCredential.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if oauthCred.Issuer != appDetails.AppDID {
		http.Error(w, "incorrect oauth cred", http.StatusBadRequest)
		return
	}
	_, _, policyCred, err := utils.ParseVerifiableCredentialFromJWT(claims.CredentialJWTs.PolicyCredential.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if policyCred.Issuer != appDetails.AppDID {
		http.Error(w, "incorrect policy cred", http.StatusBadRequest)
		return
	}
	if !utils.IsRoleExists(policyCred.CredentialSubject, role) {
		http.Error(w, "false", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("true")
}

// GetAccessList godoc
// @Summary Get Access list on the resource
// @Description List the access for the user on the resource.
// @Tags User Access Management
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token" default(Bearer YOUR_ACCESS_TOKEN)
// @Success 200 {string} string "success"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /get-access-list [get]
func (h *AuthHandler) GetAccessList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract query parameters
	queryParams := r.URL.Query()
	appDid := queryParams.Get("app_did")
	appSecret := queryParams.Get("app_secret")

	appDetails, err := h.db.GetApp(appDid)
	if err != nil {
		http.Error(w, "app is invalid", http.StatusInternalServerError)
		return
	}
	if appDetails.AppSceret != appSecret {
		http.Error(w, "app secret is invalid", http.StatusInternalServerError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized: No Authorization header provided", http.StatusUnauthorized)
		return
	}
	// Split the header to get the token part
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		http.Error(w, "Unauthorized: Invalid Authorization header format", http.StatusUnauthorized)
		return
	}
	// headerParts[1] contains the actual token
	token := headerParts[1]

	claims, err := utils.ValidateAccessToken(token)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid access token", http.StatusUnauthorized)
		return
	}
	// TODO:: revokable check
	_, _, oauthCred, err := utils.ParseVerifiableCredentialFromJWT(claims.CredentialJWTs.OAuthCredential.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if oauthCred.Issuer != appDetails.AppDID {
		http.Error(w, "incorrect oauth cred", http.StatusBadRequest)
		return
	}
	_, _, policyCred, err := utils.ParseVerifiableCredentialFromJWT(claims.CredentialJWTs.PolicyCredential.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if policyCred.Issuer != appDetails.AppDID {
		http.Error(w, "incorrect policy cred", http.StatusBadRequest)
		return
	}
	appPolicy, _ := h.db.GetIssuedPolicy(appDid)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AccessList{
		ApplicationPolicy: appPolicy.CredentialSubject,
		UserAccessList:    policyCred.CredentialSubject,
	})
}
