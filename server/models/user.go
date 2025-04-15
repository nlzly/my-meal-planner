package models

// User represents a user in the system
type User struct {
	ID       string `json:"id"` // Google's 'sub' claim
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	CreateAt string `json:"created_at,omitempty"`
	UpdateAt string `json:"updated_at,omitempty"`
}
