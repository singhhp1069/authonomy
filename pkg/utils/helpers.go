package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HTTPResponse struct {
	Body       interface{}
	StatusCode int
}

func SendHTTPRequest(method, url string, body *bytes.Buffer) (*HTTPResponse, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseJSON interface{}
	err = json.Unmarshal(responseBody, &responseJSON)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		Body:       responseJSON,
		StatusCode: resp.StatusCode,
	}, nil
}
