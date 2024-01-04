package handlers

import (
	"log"
	"net/http"
)

type MiddlewareService struct {
	apikey string
}

// NewAppHandler creates a new instance of AppHandler
func NewMiddlewareService(apikey string) *MiddlewareService {
	return &MiddlewareService{apikey: apikey}
}

// Middleware type defines a function that wraps an http.HandlerFunc
type Middleware func(http.HandlerFunc) http.HandlerFunc

// EnableCORS is a middleware that adds CORS headers to the response
func (m MiddlewareService) EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// XApiKeyMiddleware checks for a valid x-api-key in the request headers
func (m MiddlewareService) XApiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providedApiKey := r.Header.Get("x-api-key")
		if providedApiKey != m.apikey {
			http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// LoggingMiddleware logs each request
func (m MiddlewareService) LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

// ChainMiddleware provides a convenient way to chain multiple middleware functions
func (m MiddlewareService) ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
