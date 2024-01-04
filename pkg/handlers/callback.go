package handlers

import (
	"authonomy/pkg/providers"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Assuming CallbackHandler is defined elsewhere in your package
type CallbackHandler struct {
}

// NewCallbackHandler creates a new instance of AppHandler
func NewCallbackHandler() *CallbackHandler {
	return &CallbackHandler{}
}

// HandleCallback handles the callback route
func (r *CallbackHandler) HandleCallback(w http.ResponseWriter, req *http.Request) {
	// Split the URL path to get the parameters
	pathSegments := strings.Split(req.URL.Path, "/")
	if len(pathSegments) < 4 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	provider := pathSegments[2]
	did := pathSegments[3]

	fmt.Println("did", did)
	fmt.Println("provider", provider)

	// Redirect to the web page with query parameters
	redirectURL := "/web/index.html?provider=" + provider + "&did=" + did
	http.Redirect(w, req, redirectURL, http.StatusMovedPermanently)
}

func (r *CallbackHandler) HandleMe(w http.ResponseWriter, req *http.Request) {
	// Split the URL path to get the parameters
	pathSegments := strings.Split(req.URL.Path, "/")
	if len(pathSegments) < 4 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	provider := pathSegments[2]
	accessToken := pathSegments[3]

	fmt.Println("provider", provider)
	fmt.Println("accessToken", accessToken)

	userInfo, err := providers.GetUserInfo(provider, accessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
