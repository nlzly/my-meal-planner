package db

import (
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"my-meal-planner/models"
)

var (
	ErrMealNotFound = errors.New("meal not found")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidToken = errors.New("invalid token")
)

// Store is an in-memory database for meals and users
type Store struct {
	meals           map[string]models.Meal
	users           map[string]models.User
	usersByGoogleID map[string]string // Maps Google IDs to user IDs
	mutex           sync.RWMutex
	nextMealID      int
	nextUserID      int
	jwtSecret       []byte
	oauthConfig     *oauth2.Config
}

// NewStore creates a new in-memory store
func NewStore() *Store {
	// Get Google OAuth credentials from environment
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/auth/google/callback"
	}

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my-meal-planner-secret-key" // Default secret for development
	}

	// Create OAuth config
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &Store{
		meals:           make(map[string]models.Meal),
		users:           make(map[string]models.User),
		usersByGoogleID: make(map[string]string),
		nextMealID:      1,
		nextUserID:      1,
		jwtSecret:       []byte(jwtSecret),
		oauthConfig:     oauthConfig,
	}
}

// GetOAuthConfig returns the OAuth2 config
func (s *Store) GetOAuthConfig() *oauth2.Config {
	return s.oauthConfig
}

// CreateMeal adds a new meal to the store
func (s *Store) CreateMeal(req models.MealRequest) models.Meal {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := s.generateID()
	now := time.Now()

	meal := models.Meal{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Day:         req.Day,
		MealType:    req.MealType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.meals[id] = meal
	return meal
}

// GetMeal retrieves a meal by ID
func (s *Store) GetMeal(id string) (models.Meal, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	meal, exists := s.meals[id]
	if !exists {
		return models.Meal{}, ErrMealNotFound
	}

	return meal, nil
}

// UpdateMeal updates an existing meal
func (s *Store) UpdateMeal(id string, req models.MealRequest) (models.Meal, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	meal, exists := s.meals[id]
	if !exists {
		return models.Meal{}, ErrMealNotFound
	}

	meal.Name = req.Name
	meal.Description = req.Description
	meal.Day = req.Day
	meal.MealType = req.MealType
	meal.UpdatedAt = time.Now()

	s.meals[id] = meal
	return meal, nil
}

// DeleteMeal removes a meal from the store
func (s *Store) DeleteMeal(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.meals[id]; !exists {
		return ErrMealNotFound
	}

	delete(s.meals, id)
	return nil
}

// ListMeals returns all meals in the store
func (s *Store) ListMeals() []models.Meal {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	meals := make([]models.Meal, 0, len(s.meals))
	for _, meal := range s.meals {
		meals = append(meals, meal)
	}

	return meals
}

// generateID creates a new unique ID for a meal
func (s *Store) generateID() string {
	id := s.nextMealID
	s.nextMealID++
	return strconv.Itoa(id)
}

// CreateOrGetUser creates a new user or returns an existing one
func (s *Store) CreateOrGetUser(user models.User) (models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if user already exists by Google ID
	if existingUserID, exists := s.usersByGoogleID[user.GoogleID]; exists {
		return s.users[existingUserID], nil
	}

	// Create new user
	userID := strconv.Itoa(s.nextUserID)
	s.nextUserID++

	user.ID = userID
	s.users[userID] = user
	s.usersByGoogleID[user.GoogleID] = userID

	return user, nil
}

// GenerateToken creates a new JWT token for a user
func (s *Store) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *Store) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
	}

	return "", ErrInvalidToken
}
