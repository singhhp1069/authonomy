package cmd

import (
	"identitysphere-api/pkg/handlers"
	"identitysphere-api/services"
	"identitysphere-api/store"
	"log"
	"net/http"

	_ "identitysphere-api/docs" // Swaggo generates docs in this package

	httpSwagger "github.com/swaggo/http-swagger"
)

func Start() {
	// Initialize the data store (e.g., database connection)
	store, err := store.NewStore()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer store.Close()
	// // clear db before start (for the demo)
	// err = store.ClearDB()
	// if err != nil {
	// 	log.Fatalf("Failed to clean the database: %v", err)
	// }
	// Initialize services with dependencies
	ssiService := services.NewSsiClient()
	// // Dummy policies
	// err = services.CreateDemoPolicies(ssiService, store)
	// if err != nil {
	// 	log.Fatalf("Failed to create policies: %v", err)
	// }
	// Initialize handlers with services
	appHandler := handlers.NewAppHandler(ssiService, store)
	authHandler := handlers.NewAuthHandler(ssiService, store)
	policyHandler := handlers.NewPolicyHandler(ssiService, store)
	callbackHandler := handlers.NewCallbackHandler()
	credentialHandler := handlers.NewCredentialHandler(ssiService, store)
	// Swagger endpoint
	url := httpSwagger.URL("http://localhost:8081/swagger/doc.json") // The url pointing to API definition
	http.Handle("/swagger/", httpSwagger.Handler(
		url, //The url pointing to API definition
	))
	// Set up routes
	http.HandleFunc("/applications", handlers.ChainMiddleware(appHandler.HandleApplications, handlers.EnableCORS, handlers.LoggingMiddleware))
	http.HandleFunc("/application/", appHandler.GetConfig)
	http.HandleFunc("/auth-provider", handlers.EnableCORS(authHandler.GetAuthConnectorHandler))
	http.HandleFunc("/auth-provider/link", handlers.EnableCORS(authHandler.LinkAuthProviderHandler))
	http.HandleFunc("/auth-provider/unlink", handlers.EnableCORS(authHandler.UnLinkAuthProviderHandler))
	http.HandleFunc("/policies", handlers.EnableCORS(policyHandler.GetPolicyHandler))
	http.HandleFunc("/create-policy", handlers.EnableCORS(policyHandler.CreatePolicyHandler))
	http.HandleFunc("/attach-policy", handlers.EnableCORS(policyHandler.AttachPolicyHandler))
	http.HandleFunc("/callback/", callbackHandler.HandleCallback)
	http.HandleFunc("/me/", callbackHandler.HandleMe)
	http.HandleFunc("/issue-credential", credentialHandler.IssueOAuthCredential)
	http.HandleFunc("/revoke-credential", credentialHandler.RevokeOAuthCredential)
	// static web page for access_token
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))
	// Start the server
	log.Println("Starting server on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
