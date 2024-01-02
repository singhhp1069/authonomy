// main.go
package main

import (
	"identitysphere-api/cmd"
)

func main() {
	cmd.Start()
}

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// // CORS middleware to set necessary headers
// func enableCORS(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Set headers
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		// If this is a preflight request, we only need to return headers
// 		if r.Method == "OPTIONS" {
// 			return
// 		}

// 		// Call the next handler, which can be another middleware or the final handler
// 		next.ServeHTTP(w, r)
// 	}
// }

// func main() {
// 	InitDB()
// 	defer CloseDB()
// 	http.HandleFunc("/applications", enableCORS(ApplicationHandler)) // get application and there policies
// 	http.HandleFunc("/create-policy", enableCORS(createPolicyHandler))
// 	// http.HandleFunc("/get-policy", getPolicyHandler)  store policy ID and get all policy
// 	http.HandleFunc("/attach-policy", enableCORS(attachPolicyHandler))
// 	http.HandleFunc("/auth/connector", enableCORS(availableConnectorsHandler))
// 	http.HandleFunc("/auth/link", enableCORS(linkAuthProviderHandler))
// 	http.HandleFunc("/auth/unlink", enableCORS(unlinkAuthProviderHandler))
// 	log.Fatal(http.ListenAndServe(":8081", nil))
// }

// func availableConnectorsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	connectors := []AvailableProvider{
// 		{
// 			ProviderName:     "facebook",
// 			ProviderType:     "social",
// 			ProviderProtocol: "oauth2",
// 		},
// 	}
// 	json.NewEncoder(w).Encode(connectors)
// }

// // TODO:: return redirect url
// func linkAuthProviderHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var provider AuthProvider
// 	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	isExist := IsDIDExists(provider.AppID)
// 	if !isExist {
// 		http.Error(w, "app DID does not exists id: "+provider.AppID, http.StatusInternalServerError)
// 		return
// 	}

// 	// TODO: Implement logic to save provider information in the database
// 	err := SetAuthProviderInDB(provider)
// 	if err != nil {
// 		http.Error(w, "Failed to save provide details "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Provider linked successfully"})
// }

// func unlinkAuthProviderHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var provider AuthProvider
// 	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// TODO: Implement logic to remove provider information from the database
// 	// removeAuthProvider(provider)

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Provider unlinked successfully"})
// }

// func createPolicyHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var schema PolicySchemaRequest
// 	err := json.NewDecoder(r.Body).Decode(&schema)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	ssiClient := NewSsiClient()
// 	respSchema, err := ssiClient.CreateSchema(schema)
// 	print("res", respSchema)
// 	if err != nil {
// 		http.Error(w, "Failed to create schema: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(respSchema); err != nil {
// 		http.Error(w, "Failed to encode schema response", http.StatusInternalServerError)
// 	}
// }

// // TODO:: attach batch policies, remove policy
// func attachPolicyHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	var appPolicy ApplicationPolicy // store credential ID too
// 	err := json.NewDecoder(r.Body).Decode(&appPolicy)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	fmt.Println(appPolicy.ApplicationDID, appPolicy.SchemaID)
// 	isExist := IsSchemaExists(appPolicy.SchemaID)
// 	if !isExist {
// 		http.Error(w, "Schema does not exists id: "+appPolicy.SchemaID, http.StatusInternalServerError)
// 		return
// 	}
// 	isExist = IsDIDExists(appPolicy.ApplicationDID)
// 	if !isExist {
// 		http.Error(w, "app DID does not exists id: "+appPolicy.ApplicationDID, http.StatusInternalServerError)
// 		return
// 	}
// 	isExist = IsDIDExists(appPolicy.IssuerDID)
// 	if !isExist {
// 		http.Error(w, "issuer DID does not exists id: "+appPolicy.IssuerDID, http.StatusInternalServerError)
// 		return
// 	}

// 	ssiClient := NewSsiClient()
// 	respSchema, err := ssiClient.IssueCredentialBySchemaID(appPolicy.IssuerDID, appPolicy.ApplicationDID, appPolicy.SchemaID, appPolicy.Credential)
// 	print("res", respSchema)
// 	if err != nil {
// 		http.Error(w, "Failed to create schema: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	// Save the new application details to BadgerDB
// 	err = SetSchemaInDB(appPolicy)
// 	if err != nil {
// 		http.Error(w, "Failed to save application policy: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(respSchema)
// }

// func ApplicationHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		getApplications(w, r)
// 	case "POST":
// 		createApplication(w, r)
// 	default:
// 		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// }
