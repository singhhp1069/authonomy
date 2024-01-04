package handlers

import (
	"encoding/json"
	"identitysphere-api/models"
	"identitysphere-api/services"
	"identitysphere-api/store"
	"net/http"

	"github.com/go-playground/validator"
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

// @Summary Get all applications
// @Description Retrieves a list of all applications
// @Tags Application Management
// @Accept json
// @Produce json
// @Param x-api-key header string true "API Key"
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
// @Param x-api-key header string true "API Key"
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
