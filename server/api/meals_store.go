package api

import (
	"encoding/json"
	"my-meal-planner/db"
	"my-meal-planner/models"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// listMeals returns all meals for a specific meal plan
func (h *Handler) listMeals(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := h.store.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	mealPlanID := r.URL.Query().Get("mealPlanId")
	if mealPlanID == "" {
		http.Error(w, "Meal plan ID is required", http.StatusBadRequest)
		return
	}

	// Check if user has access to this meal plan
	hasAccess, err := h.store.CheckMealPlanAccess(claims.UserID, mealPlanID)
	if err != nil {
		http.Error(w, "Failed to check access", http.StatusInternalServerError)
		return
	}

	if !hasAccess {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	meals := h.store.ListMealsByPlan(mealPlanID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

// createMeal creates a new meal in a meal plan
func (h *Handler) createMeal(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := h.store.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var req struct {
		models.MealRequest `json:"meal"`
		MealPlanID         string `json:"mealPlanId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.MealPlanID == "" {
		http.Error(w, "Meal plan ID is required", http.StatusBadRequest)
		return
	}

	// Check if user has edit access to this meal plan
	hasAccess, err := h.store.CheckMealPlanAccess(claims.UserID, req.MealPlanID)
	if err != nil {
		http.Error(w, "Failed to check access", http.StatusInternalServerError)
		return
	}

	if !hasAccess {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	meal := &models.Meal{
		ID:          uuid.New().String(),
		MealPlanID:  req.MealPlanID,
		Name:        req.Name,
		Description: req.Description,
		Day:         req.Day,
		MealType:    req.MealType,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.store.CreateMeal(meal); err != nil {
		http.Error(w, "Failed to create meal", http.StatusInternalServerError)
		return
	}

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

	// Get existing meal
	existingMeal, err := h.store.GetMeal(id)
	if err != nil {
		if err == db.ErrMealNotFound {
			http.Error(w, "Meal not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Update meal fields
	existingMeal.Name = req.Name
	existingMeal.Description = req.Description
	existingMeal.Day = req.Day
	existingMeal.MealType = req.MealType
	existingMeal.UpdatedAt = time.Now()

	// Update the meal
	if err := h.store.UpdateMeal(existingMeal); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingMeal)
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
