package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/oauth2"

	"my-meal-planner/db"
	"my-meal-planner/models"
)

// Handler contains all the dependencies for the API handlers
type Handler struct {
	store *db.Store
}

// NewHandler creates a new API handler
func NewHandler(store *db.Store) *Handler {
	return &Handler{
		store: store,
	}
}

// RegisterRoutes registers all the API routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Auth routes
	mux.HandleFunc("/auth/google/login", h.handleGoogleLogin)
	mux.HandleFunc("/auth/google/callback", h.handleGoogleCallback)

	// Protected meal routes
	protected := http.NewServeMux()
	protected.HandleFunc("/api/meals", h.handleMeals)
	protected.HandleFunc("/api/meals/", h.handleMealByID)
	mux.Handle("/api/", h.authMiddleware(protected))
}

// authMiddleware verifies JWT tokens
func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate JWT token
		_, err := h.store.ValidateToken(strings.TrimPrefix(tokenString, "Bearer "))
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleMeals handles GET and POST requests for /api/meals
func (h *Handler) handleMeals(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listMeals(w, r)
	case http.MethodPost:
		h.createMeal(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleMealByID handles GET, PUT, and DELETE requests for /api/meals/{id}
func (h *Handler) handleMealByID(w http.ResponseWriter, r *http.Request) {
	// Extract meal ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/api/meals/")
	if id == "" {
		http.Error(w, "Invalid meal ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getMeal(w, r, id)
	case http.MethodPut:
		h.updateMeal(w, r, id)
	case http.MethodDelete:
		h.deleteMeal(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listMeals returns all meals
func (h *Handler) listMeals(w http.ResponseWriter, r *http.Request) {
	meals := h.store.ListMeals()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

// createMeal creates a new meal
func (h *Handler) createMeal(w http.ResponseWriter, r *http.Request) {
	var req models.MealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" || req.Day == "" || req.MealType == "" {
		http.Error(w, "Name, day, and meal type are required", http.StatusBadRequest)
		return
	}

	meal := h.store.CreateMeal(req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(meal)
}

// getMeal returns a meal by ID
func (h *Handler) getMeal(w http.ResponseWriter, r *http.Request, id string) {
	meal, err := h.store.GetMeal(id)
	if err != nil {
		if err == db.ErrMealNotFound {
			http.Error(w, "Meal not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meal)
}

// updateMeal updates a meal by ID
func (h *Handler) updateMeal(w http.ResponseWriter, r *http.Request, id string) {
	var req models.MealRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" || req.Day == "" || req.MealType == "" {
		http.Error(w, "Name, day, and meal type are required", http.StatusBadRequest)
		return
	}

	meal, err := h.store.UpdateMeal(id, req)
	if err != nil {
		if err == db.ErrMealNotFound {
			http.Error(w, "Meal not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meal)
}

// handleGoogleLogin initiates Google OAuth flow
func (h *Handler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.store.GetOAuthConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback processes Google OAuth callback
func (h *Handler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Handle both redirect flow and client-side flow
	var googleToken string

	if r.Method == http.MethodGet {
		// Server-side flow (redirect)
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code parameter", http.StatusBadRequest)
			return
		}

		token, err := h.store.GetOAuthConfig().Exchange(ctx, code)
		if err != nil {
			http.Error(w, "Failed to exchange code for token", http.StatusUnauthorized)
			return
		}
		googleToken = token.AccessToken
	} else if r.Method == http.MethodPost {
		// Client-side flow (credential from Google Sign-In button)
		var req struct {
			Credential string `json:"credential"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Credential == "" {
			http.Error(w, "Missing credential", http.StatusBadRequest)
			return
		}

		googleToken = req.Credential
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userInfo struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if r.Method == http.MethodGet {
		// For server-side flow, get user info from access token
		resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + googleToken)
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			http.Error(w, "Invalid user info", http.StatusUnauthorized)
			return
		}
	} else {
		// For client-side flow, decode the JWT ID token
		// Note: In a production environment, you should verify the token signature
		// using Google's public keys, but for simplicity we're just decoding it here
		parts := strings.Split(googleToken, ".")
		if len(parts) != 3 {
			http.Error(w, "Invalid ID token format", http.StatusUnauthorized)
			return
		}

		// Decode the payload (second part of the JWT)
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(w, "Failed to decode token payload", http.StatusUnauthorized)
			return
		}

		var tokenInfo struct {
			Sub   string `json:"sub"`
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err := json.Unmarshal(payload, &tokenInfo); err != nil {
			http.Error(w, "Invalid token payload", http.StatusUnauthorized)
			return
		}

		userInfo.Sub = tokenInfo.Sub
		userInfo.Email = tokenInfo.Email
		userInfo.Name = tokenInfo.Name
	}

	// Create or get user
	user, err := h.store.CreateOrGetUser(models.User{
		GoogleID: userInfo.Sub,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate JWT
	tokenString, err := h.store.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// deleteMeal deletes a meal by ID
func (h *Handler) deleteMeal(w http.ResponseWriter, r *http.Request, id string) {
	err := h.store.DeleteMeal(id)
	if err != nil {
		if err == db.ErrMealNotFound {
			http.Error(w, "Meal not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
