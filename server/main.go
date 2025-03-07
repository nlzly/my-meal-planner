package main

import (
	"encoding/json"
	"log"
	"net/http"

	"my-meal-planner/api"
	"my-meal-planner/db"
)

func main() {
	// Create a new store
	store := db.NewStore()

	// Create a new API handler
	handler := api.NewHandler(store)

	// Create a new HTTP server mux
	mux := http.NewServeMux()

	// Register health check route
	mux.HandleFunc("/api/health", healthCheckHandler)

	// Register API routes
	handler.RegisterRoutes(mux)

	// Enable CORS
	corsHandler := enableCORS(mux)

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "ok",
		"message": "Server is running",
	}

	json.NewEncoder(w).Encode(response)
}

// enableCORS wraps a handler with CORS support
func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the original handler
		h.ServeHTTP(w, r)
	})
}
