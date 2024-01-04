package cmd

import (
	"authonomy/pkg/handlers"
	"authonomy/services"
	"authonomy/store"
	"fmt"
	"log"
	"net/http"

	_ "authonomy/docs" // Swaggo generates docs in this package

	httpSwagger "github.com/swaggo/http-swagger"
)

func Start(dbPath, secret, apiKey, port, ssiUrl string, reset bool) {
	// Initialize the data store (e.g., database connection)
	store, err := store.NewStore(dbPath, secret)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer store.Close()
	// Initialize services with dependencies
	ssiService := services.NewSsiClient(ssiUrl)
	// clear db before start (for the demo) or use the reset flag
	if reset {
		err = store.ClearDB()
		if err != nil {
			log.Fatalf("Failed to clean the database: %v", err)
		}

		// Dummy policies
		err = services.CreateDemoPolicies(ssiService, store)
		if err != nil {
			log.Fatalf("Failed to create policies: %v", err)
		}
	}

	fmt.Println("=======================")
	fmt.Println("\033[32m", "------x-api-key------", "\033[0m")
	fmt.Println("\033[32m", apiKey, "\033[0m")
	fmt.Println("=======================")
	fmt.Println("=====Swagger=======")
	fmt.Println("\033[32m", fmt.Sprintf("http://localhost%s/swagger", port), "\033[0m")
	fmt.Println("=======================")
	// Initialize handlers with services
	m := handlers.NewMiddlewareService(apiKey)
	appHandler := handlers.NewAppHandler(ssiService, store)
	authProviderHandler := handlers.NewAuthProviderHandler(ssiService, store)
	policyHandler := handlers.NewPolicyHandler(ssiService, store)
	callbackHandler := handlers.NewCallbackHandler()
	credentialHandler := handlers.NewCredentialHandler(ssiService, store)
	authHandler := handlers.NewAuthHandler(ssiService, store)
	// Swagger endpoint
	url := httpSwagger.URL("http://localhost:8081/swagger/doc.json") // The url pointing to API definition
	http.Handle("/swagger/", httpSwagger.Handler(
		url, //The url pointing to API definition
	))
	// Set up routes
	// application owner access
	http.HandleFunc("/applications", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(appHandler.HandleApplications))

	http.HandleFunc("/auth-provider", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(authProviderHandler.GetAuthConnectorHandler))
	http.HandleFunc("/auth-provider/link", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(authProviderHandler.LinkAuthProviderHandler))
	http.HandleFunc("/auth-provider/unlink", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(authProviderHandler.UnLinkAuthProviderHandler))

	http.HandleFunc("/policies", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(policyHandler.GetPolicyHandler))
	http.HandleFunc("/create-policy", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(policyHandler.CreatePolicyHandler))
	http.HandleFunc("/attach-policy", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(policyHandler.AttachPolicyHandler))

	http.HandleFunc("/grant-access", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(authHandler.GrandAccess))
	http.HandleFunc("/revoke-access", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(authHandler.RevokeAccess))
	http.HandleFunc("/revoke-credential", m.ChainMiddleware(m.XApiKeyMiddleware, m.LoggingMiddleware)(credentialHandler.RevokeOAuthCredential))

	// application itself access
	http.HandleFunc("/validate-access", m.ChainMiddleware(m.LoggingMiddleware)(authHandler.VerifyAccess))
	http.HandleFunc("/issue-credential", m.ChainMiddleware(m.EnableCORS, m.LoggingMiddleware)(credentialHandler.IssueOAuthCredential))
	// application user access
	http.HandleFunc("/callback/", m.ChainMiddleware(m.EnableCORS, m.LoggingMiddleware)(callbackHandler.HandleCallback))
	http.HandleFunc("/me/", m.ChainMiddleware(m.EnableCORS, m.LoggingMiddleware)(callbackHandler.HandleMe))
	http.HandleFunc("/signup", m.ChainMiddleware(m.EnableCORS, m.LoggingMiddleware)(authHandler.SignUpHandler))

	http.HandleFunc("/get-access-token", m.ChainMiddleware(m.LoggingMiddleware)(authHandler.GetAccessToken))
	http.HandleFunc("/request-access", m.ChainMiddleware(m.LoggingMiddleware)(authHandler.RequestAccess))

	// static web page for access_token
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))
	// Start the server
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
