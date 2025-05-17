package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"my-meal-planner/db"
	"my-meal-planner/models"
)

// Handler contains all the dependencies for the API handlers
type Handler struct {
	store db.Store
}

// NewHandler creates a new API handler
func NewHandler(store db.Store) *Handler {
	return &Handler{
		store: store,
	}
}

// handleMealPlans handles GET and POST requests for /api/meal-plans
func (h *Handler) handleMealPlans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listMealPlans(w, r)
	case http.MethodPost:
		h.createMealPlan(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
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

// handleMealPlanByID handles GET, PUT, and DELETE requests for /api/meal-plans/{id}
func (h *Handler) handleMealPlanByID(w http.ResponseWriter, r *http.Request) {
	// Extract meal plan ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/api/meal-plans/")
	if id == "" {
		http.Error(w, "Invalid meal plan ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getMealPlan(w, r, id)
	case http.MethodPut:
		h.updateMealPlan(w, r, id)
	case http.MethodDelete:
		h.deleteMealPlan(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGenerateShareLink handles generating a sharing link for a meal plan
func (h *Handler) handleGenerateShareLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate token
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := h.store.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Parse request
	var req struct {
		MealPlanID string `json:"mealPlanId"`
		Role       string `json:"role"`      // "editor" or "viewer"
		ExpiresIn  int    `json:"expiresIn"` // expiration in hours (optional)
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.MealPlanID == "" {
		http.Error(w, "Meal plan ID is required", http.StatusBadRequest)
		return
	}

	// Default role to "viewer" if not specified or invalid
	if req.Role != "editor" && req.Role != "viewer" {
		req.Role = "viewer"
	}

	// Default expiration to 7 days if not specified or invalid
	if req.ExpiresIn <= 0 {
		req.ExpiresIn = 7 * 24 // 7 days in hours
	}

	// Check if user has owner access to this meal plan
	isOwner, err := h.store.CheckMealPlanOwnership(claims.UserID, req.MealPlanID)
	if err != nil || !isOwner {
		http.Error(w, "Only the owner can share a meal plan", http.StatusForbidden)
		return
	}

	// Generate a unique share code
	shareCode := uuid.New().String()

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)

	// Create share link information
	shareLink := &models.ShareLink{
		ID:         shareCode,
		MealPlanID: req.MealPlanID,
		CreatedBy:  claims.UserID,
		Role:       req.Role,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
	}

	// Store the share link
	if err := h.store.CreateShareLink(shareLink); err != nil {
		http.Error(w, "Failed to create share link", http.StatusInternalServerError)
		return
	}

	// Return the share link
	shareURL := fmt.Sprintf("%s/#/join?code=%s", "http://localhost:5173", shareCode)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"shareLink": shareURL,
		"code":      shareCode,
		"expiresAt": expiresAt.Format(time.RFC3339),
	})
}
