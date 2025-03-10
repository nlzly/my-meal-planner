package models

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	GoogleID string `json:"google_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}
