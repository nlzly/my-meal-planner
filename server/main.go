package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"my-meal-planner/api"
	"my-meal-planner/db"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
		// Continue anyway, as env vars might be set directly
	}

	// Get Google OAuth credentials from environment
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")

	// Debug logging for OAuth credentials
	log.Println("OAuth Configuration:")
	log.Println("Client ID:", clientID)
	if clientSecret != "" {
		log.Println("Client Secret:", clientSecret[:5]+"..."+clientSecret[len(clientSecret)-5:])
	} else {
		log.Println("Client Secret: <not set>")
	}
	log.Println("Redirect URL:", redirectURL)

	if redirectURL == "" {
		redirectURL = "http://localhost:8080/auth/google/callback"
	}

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my-meal-planner-secret-key" // Default secret for development
	}

	// Create OAuth config
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	// Create store
	store := db.NewMemoryStore(oauthConfig, []byte(jwtSecret))

	// Create handler
	handler := api.NewHandler(store)

	// Create mux
	mux := http.NewServeMux()

	// Register routes
	handler.RegisterRoutes(mux)

	// Serve static files from the client directory
	fs := http.FileServer(http.Dir("../client/dist"))
	mux.Handle("/", fs)

	// Add CORS middleware
	corsHandler := enableCORS(mux)

	// Start server
	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatal(err)
	}
}

// enableCORS wraps a handler with CORS support
func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := map[string]bool{
			"http://localhost:5173":                true,
			"https://your-frontend.up.railway.app": true,
		}
		origin := r.Header.Get("Origin")
		// Set CORS headers
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin") // Required for varying by Origin
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the original handler
		h.ServeHTTP(w, r)
	})
}
