package providers

import (
	"authonomy/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type FacebookUserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetUserInfo(provider, accessToken string) (*FacebookUserInfo, error) {
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
	responseMap, ok := httpResponse.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}
	// Marshal map back to JSON bytes
	jsonBytes, err := json.Marshal(responseMap)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response: %s", err)
	}
	// Unmarshal JSON bytes into FacebookUserInfo struct
	var userInfo FacebookUserInfo
	err = json.Unmarshal(jsonBytes, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response to struct: %s", err)
	}

	return &userInfo, nil
}
func GetLoginUrl(clientID, redirectUrl string) string {
	return fmt.Sprintf("https://www.facebook.com/v14.0/dialog/oauth?client_id=%s&redirect_uri=%s&display=popup&response_type=token&auth_type=reauthenticate", clientID, redirectUrl)
}
