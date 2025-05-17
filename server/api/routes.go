package api

import "net/http"

// RegisterRoutes registers all the API routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Auth routes
	mux.HandleFunc("/auth/google/login", h.handleGoogleLogin)
	mux.HandleFunc("/auth/google/callback", h.handleGoogleCallback)

	// Protected routes
	protected := http.NewServeMux()

	// Meal plan routes
	protected.HandleFunc("/api/meal-plans", h.handleMealPlans)
	protected.HandleFunc("/api/meal-plans/", h.handleMealPlanByID)
	protected.HandleFunc("/api/meal-plans/share", h.handleShareMealPlan)
	protected.HandleFunc("/api/meal-plans/generate-link", h.handleGenerateShareLink)
	protected.HandleFunc("/api/meal-plans/join", h.handleJoinMealPlan)

	// Meal routes
	protected.HandleFunc("/api/meals", h.handleMeals)
	protected.HandleFunc("/api/meals/", h.handleMealByID)

	mux.Handle("/api/", h.authMiddleware(protected))
}
