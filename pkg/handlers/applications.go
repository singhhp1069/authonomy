package handlers

import (
	"encoding/json"
	"identitysphere-api/models"
	"identitysphere-api/services"
	"identitysphere-api/store"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// AppHandler handles application-related requests
type AppHandler struct {
	ssiService *services.SsiClient
	db         *store.Store
}

// NewAppHandler creates a new instance of AppHandler
func NewAppHandler(ssiService *services.SsiClient, db *store.Store) *AppHandler {
	return &AppHandler{ssiService: ssiService, db: db}
}

// HandleApplications routes the request to the appropriate function based on the HTTP method
func (h *AppHandler) HandleApplications(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getApplications(w, r)
	case "POST":
		h.createApplication(w, r)
	default:
		http.Error(w, "Unsupported HTTP Method", http.StatusMethodNotAllowed)
	}
}

// GetConfig godoc
// @Summary Get application configuration
// @Description Retrieve the configuration for a specific application.
// @Tags Application Management
// @Accept  json
// @Produce  json
// @Param app_did path string true "Application DID"
// @Success 200 {object} models.Config "Configuration retrieved successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Application not found"
// @Failure 500 {string} string "Internal server error"
// @Router /application/{app_did}/config [get]
func (h *AppHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// Split the URL path into segments
	pathSegments := strings.Split(r.URL.Path, "/")
	// The URL format is /application/{app_did}/config
	// Checking for minimum segments and specific segments
	if len(pathSegments) < 4 || pathSegments[1] != "application" || pathSegments[3] != "config" {
		http.NotFound(w, r)
		return
	}

	appID := pathSegments[2]
	app, err := h.db.GetApp(appID)
	if err != nil {
		http.Error(w, "app did is not configured yet", http.StatusInternalServerError)
		return
	}

	auth, err := h.db.GetAuthProvider(appID)
	if err != nil {
		http.Error(w, "app authentication is not configured yet", http.StatusInternalServerError)
		return
	}

	cred, err := h.db.GetIssuedPolicy(appID)
	if err != nil {
		http.Error(w, "policy is not configured yet", http.StatusInternalServerError)
		return
	}
	config := models.Config{
		App:        app,
		Auth:       auth,
		Cred:       cred,
		SessionKey: uuid.New().String(),
	}
	err = h.db.SetConfig(config)
	if err != nil {
		http.Error(w, "config not set", http.StatusInternalServerError)
		return
	}
	// Send the app config as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config) // Assuming 'app' holds the configuration
}

// @Summary Get all applications
// @Description Retrieves a list of all applications
// @Tags Application Management
// @Accept json
// @Produce json
// @Success 200 {array} models.ApplicationResponse
// @Failure 500 {object} string "Internal Server Error"
// @Router /applications [get]
func (h *AppHandler) getApplications(w http.ResponseWriter, r *http.Request) {
	apps, err := h.db.GetAllApps()
	if err != nil {
		http.Error(w, "Failed to get applications: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

// createApplication handles POST requests to create a new application
// @Summary Create a new application
// @Description Create a new application with the given details
// @Tags Application Management
// @Accept json
// @Produce json
// @Param application body models.ApplicationRequest true "Application to create"
// @Success 200 {object} models.ApplicationResponse
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /applications [post]
func (h *AppHandler) createApplication(w http.ResponseWriter, r *http.Request) {
	var validate = validator.New()
	var appReq models.ApplicationRequest

	err := json.NewDecoder(r.Body).Decode(&appReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(appReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	did, err := h.ssiService.CreateDid()
	if err != nil {
		http.Error(w, "Failed to create DID: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response := models.ApplicationResponse{
		AppDID:     did,
		AppName:    appReq.AppName,
		AppDetails: appReq.AppDetails,
	}
	err = h.db.SetApp(response)
	if err != nil {
		http.Error(w, "Failed to save application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
