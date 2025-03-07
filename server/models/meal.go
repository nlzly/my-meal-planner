package models

import "time"

// Meal represents a single meal in the meal planner
type Meal struct {
	ID          string    `json:"id"`
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
