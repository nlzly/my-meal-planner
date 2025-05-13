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
	shareLink, err := h.store.GetShareLink(req.Code)
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
