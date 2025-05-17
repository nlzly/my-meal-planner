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

// listMealPlans returns all meal plans the user has access to
func (h *Handler) listMealPlans(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := h.store.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	mealPlans := h.store.ListMealPlansByUser(claims.UserID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mealPlans)
}

// createMealPlan creates a new meal plan
func (h *Handler) createMealPlan(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := h.store.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mealPlan := &models.MealPlan{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   claims.UserID,
	}

	if err := h.store.CreateMealPlan(mealPlan); err != nil {
		http.Error(w, "Failed to create meal plan", http.StatusInternalServerError)
		return
	}

	// Create owner access for the creator
	access := &models.MealPlanAccess{
		UserID:     claims.UserID,
		MealPlanID: mealPlan.ID,
		Role:       "owner",
	}

	if err := h.store.CreateMealPlanAccess(access); err != nil {
		http.Error(w, "Failed to create meal plan access", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mealPlan)
}

// getMealPlan returns a meal plan by ID
func (h *Handler) getMealPlan(w http.ResponseWriter, r *http.Request, id string) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := h.store.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if user has access to this meal plan
	hasAccess, err := h.store.CheckMealPlanAccess(claims.UserID, id)
	if err != nil {
		http.Error(w, "Failed to check access", http.StatusInternalServerError)
		return
	}

	if !hasAccess {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	mealPlan, err := h.store.GetMealPlan(id)
	if err != nil {
		if err == db.ErrMealPlanNotFound {
			http.Error(w, "Meal plan not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mealPlan)
}

// updateMealPlan updates a meal plan by ID
func (h *Handler) updateMealPlan(w http.ResponseWriter, r *http.Request, id string) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := h.store.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if user has edit access to this meal plan
	hasAccess, err := h.store.CheckMealPlanAccess(claims.UserID, id)
	if err != nil {
		http.Error(w, "Failed to check access", http.StatusInternalServerError)
		return
	}

	if !hasAccess {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get existing meal plan
	existingPlan, err := h.store.GetMealPlan(id)
	if err != nil {
		if err == db.ErrMealPlanNotFound {
			http.Error(w, "Meal plan not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Update meal plan fields
	existingPlan.Name = req.Name
	existingPlan.Description = req.Description
	existingPlan.UpdatedAt = time.Now()

	// Update the meal plan
	if err := h.store.UpdateMealPlan(existingPlan); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingPlan)
}

// deleteMealPlan deletes a meal plan by ID
func (h *Handler) deleteMealPlan(w http.ResponseWriter, r *http.Request, id string) {
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims, err := h.store.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if user has owner access to this meal plan
	hasAccess, err := h.store.CheckMealPlanAccess(claims.UserID, id)
	if err != nil {
		http.Error(w, "Failed to check access", http.StatusInternalServerError)
		return
	}

	if !hasAccess {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	err = h.store.DeleteMealPlan(id)
	if err != nil {
		if err == db.ErrMealPlanNotFound {
			http.Error(w, "Meal plan not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleShareMealPlan handles sharing a meal plan with another user
func (h *Handler) handleShareMealPlan(w http.ResponseWriter, r *http.Request) {
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
		Email      string `json:"email"`
		Role       string `json:"role"` // "editor" or "viewer"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.MealPlanID == "" || req.Email == "" {
		http.Error(w, "Meal plan ID and email are required", http.StatusBadRequest)
		return
	}

	// Validate role
	if req.Role != "editor" && req.Role != "viewer" {
		http.Error(w, "Invalid role. Must be 'editor' or 'viewer'", http.StatusBadRequest)
		return
	}

	// Check if user has owner access to this meal plan
	isOwner, err := h.store.CheckMealPlanOwnership(claims.UserID, req.MealPlanID)
	if err != nil || !isOwner {
		http.Error(w, "Only the owner can share a meal plan", http.StatusForbidden)
		return
	}

	// Find the user by email
	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Don't allow sharing with yourself
	if user.ID == claims.UserID {
		http.Error(w, "Cannot share with yourself", http.StatusBadRequest)
		return
	}

	// Create access record
	access := &models.MealPlanAccess{
		ID:         uuid.New().String(),
		UserID:     user.ID,
		MealPlanID: req.MealPlanID,
		Role:       req.Role,
	}

	if err := h.store.CreateMealPlanAccess(access); err != nil {
		http.Error(w, "Failed to share meal plan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Meal plan shared successfully",
	})
}

// handleJoinMealPlan handles joining a meal plan via a share link
func (h *Handler) handleJoinMealPlan(w http.ResponseWriter, r *http.Request) {
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
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Code == "" {
		http.Error(w, "Share code is required", http.StatusBadRequest)
		return
	}

	// Get the share link information
	shareLink, err := h.store.GetShareCode(req.Code)
	if err != nil {
		http.Error(w, "Invalid share code", http.StatusNotFound)
		return
	}

	// Check if the share link has expired
	if shareLink.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Share link has expired", http.StatusForbidden)
		return
	}

	// Check if the user is the owner of the meal plan (can't join their own plan)
	isOwner, _ := h.store.CheckMealPlanOwnership(claims.UserID, shareLink.MealPlanID)
	if isOwner {
		http.Error(w, "You already own this meal plan", http.StatusBadRequest)
		return
	}

	// Check if the user already has access to the meal plan
	hasAccess, _ := h.store.CheckMealPlanAccess(claims.UserID, shareLink.MealPlanID)
	if hasAccess {
		http.Error(w, "You already have access to this meal plan", http.StatusBadRequest)
		return
	}

	// Create access record
	access := &models.MealPlanAccess{
		ID:         uuid.New().String(),
		UserID:     claims.UserID,
		MealPlanID: shareLink.MealPlanID,
		Role:       shareLink.Role,
	}

	if err := h.store.CreateMealPlanAccess(access); err != nil {
		http.Error(w, "Failed to join meal plan", http.StatusInternalServerError)
		return
	}

	// Get the meal plan information to return to the client
	mealPlan, err := h.store.GetMealPlan(shareLink.MealPlanID)
	if err != nil {
		http.Error(w, "Failed to get meal plan information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Successfully joined meal plan",
		"mealPlan": mealPlan,
		"role":     shareLink.Role,
	})
}
