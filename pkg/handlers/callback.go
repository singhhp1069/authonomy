package handlers

import (
	"encoding/json"
	"fmt"
	"identitysphere-api/pkg/utils"
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

	userInfo, err := getUserInfo(provider, accessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func getUserInfo(provider, accessToken string) (interface{}, error) {
	if provider != "facebook" {
		// Handle other providers or return an error
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	url := fmt.Sprintf("https://graph.facebook.com/v14.0/me?access_token=%s", accessToken)
	fmt.Println("url is", url)
	fmt.Println("provider is", provider)
	fmt.Println("accessToken is", accessToken)
	// Use SendHTTPRequest to make the API call
	httpResponse, err := utils.SendHTTPRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", httpResponse.StatusCode)
	}
	return httpResponse.Body, nil
}
