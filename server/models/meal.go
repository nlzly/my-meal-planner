package models

import "time"

// Meal represents a single meal in the meal planner
type Meal struct {
	ID          string    `json:"id"`
	MealPlanID  string    `json:"mealPlanId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Day         string    `json:"day"`
	MealType    string    `json:"mealType"` // Breakfast, Lunch, Dinner
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// MealRequest is used for creating or updating a meal
type MealRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Day         string `json:"day"`
	MealType    string `json:"mealType"`
}

// MealPlan represents a collection of meals
type MealPlan struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedBy   string    `json:"createdBy"` // User ID who created the plan
}

// MealPlanAccess represents a user's access to a meal plan
type MealPlanAccess struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	MealPlanID string `json:"mealPlanId"`
	Role       string `json:"role"` // "owner", "editor", "viewer"
}

// ShareLink represents a link that can be used to join a meal plan
type ShareLink struct {
	ID         string    `json:"id"` // The unique code for the link
	MealPlanID string    `json:"mealPlanId"`
	CreatedBy  string    `json:"createdBy"` // User ID who created the link
	Role       string    `json:"role"`      // Role to assign when joining: "editor" or "viewer"
	ExpiresAt  time.Time `json:"expiresAt"` // When the link expires
	CreatedAt  time.Time `json:"createdAt"`
}
