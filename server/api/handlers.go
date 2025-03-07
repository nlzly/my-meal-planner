package api

import (
	"encoding/json"
	"net/http"
	"strings"

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
	mux.HandleFunc("/api/meals", h.handleMeals)
	mux.HandleFunc("/api/meals/", h.handleMealByID)
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
