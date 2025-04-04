package db

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"my-meal-planner/models"
)

var (
	ErrMealNotFound     = errors.New("meal not found")
	ErrMealPlanNotFound = errors.New("meal plan not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidToken     = errors.New("invalid token")
	ErrAccessDenied     = errors.New("access denied")
)

// Store defines the interface for data storage operations
type Store interface {
	// User operations
	CreateOrUpdateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)

	// Token operations
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)

	// OAuth operations
	GetOAuthConfig() *oauth2.Config

	// Meal plan operations
	CreateMealPlan(plan *models.MealPlan) error
	GetMealPlan(id string) (*models.MealPlan, error)
	UpdateMealPlan(plan *models.MealPlan) error
	DeleteMealPlan(id string) error
	ListMealPlansByUser(userID string) []*models.MealPlan
	CreateMealPlanAccess(access *models.MealPlanAccess) error
	CheckMealPlanAccess(userID, mealPlanID string) (bool, error)
	CheckMealPlanOwnership(userID, mealPlanID string) (bool, error)

	// Share link operations
	CreateShareLink(link *models.ShareLink) error
	GetShareLink(id string) (*models.ShareLink, error)
	DeleteShareLink(id string) error

	// Meal operations
	CreateMeal(meal *models.Meal) error
	GetMeal(id string) (*models.Meal, error)
	UpdateMeal(meal *models.Meal) error
	DeleteMeal(id string) error
	ListMealsByPlan(mealPlanID string) []*models.Meal
}

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// MemoryStore implements the Store interface using in-memory storage
type MemoryStore struct {
	users          map[string]*models.User
	meals          map[string]*models.Meal
	mealPlans      map[string]*models.MealPlan
	mealPlanAccess map[string]*models.MealPlanAccess
	shareLinks     map[string]*models.ShareLink
	mutex          sync.RWMutex
	oauthConfig    *oauth2.Config
	jwtSecret      []byte
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore(oauthConfig *oauth2.Config, jwtSecret []byte) *MemoryStore {
	return &MemoryStore{
		users:          make(map[string]*models.User),
		meals:          make(map[string]*models.Meal),
		mealPlans:      make(map[string]*models.MealPlan),
		mealPlanAccess: make(map[string]*models.MealPlanAccess),
		shareLinks:     make(map[string]*models.ShareLink),
		oauthConfig:    oauthConfig,
		jwtSecret:      jwtSecret,
	}
}

// generateID generates a unique ID
func (s *MemoryStore) generateID() string {
	return time.Now().Format("20060102150405") + "-" + uuid.New().String()
}

// GetOAuthConfig returns the OAuth2 config
func (s *MemoryStore) GetOAuthConfig() *oauth2.Config {
	return s.oauthConfig
}

// CreateMeal adds a new meal to the store
func (s *MemoryStore) CreateMeal(meal *models.Meal) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if meal.ID == "" {
		meal.ID = s.generateID()
	}
	s.meals[meal.ID] = meal
	return nil
}

// GetMeal retrieves a meal by ID
func (s *MemoryStore) GetMeal(id string) (*models.Meal, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	meal, exists := s.meals[id]
	if !exists {
		return nil, ErrMealNotFound
	}
	return meal, nil
}

// UpdateMeal updates an existing meal
func (s *MemoryStore) UpdateMeal(meal *models.Meal) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	existingMeal, exists := s.meals[meal.ID]
	if !exists {
		return ErrMealNotFound
	}

	existingMeal.Name = meal.Name
	existingMeal.Description = meal.Description
	existingMeal.Day = meal.Day
	existingMeal.MealType = meal.MealType
	existingMeal.UpdatedAt = time.Now()

	s.meals[meal.ID] = existingMeal
	return nil
}

// DeleteMeal removes a meal from the store
func (s *MemoryStore) DeleteMeal(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.meals[id]; !exists {
		return ErrMealNotFound
	}

	delete(s.meals, id)
	return nil
}

// ListMealsByPlan returns all meals in a specific meal plan
func (s *MemoryStore) ListMealsByPlan(mealPlanID string) []*models.Meal {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var meals []*models.Meal
	for _, meal := range s.meals {
		if meal.MealPlanID == mealPlanID {
			meals = append(meals, meal)
		}
	}
	return meals
}

// CreateOrGetUser creates a new user or returns an existing one
func (s *MemoryStore) CreateOrGetUser(user models.User) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if user already exists by Google ID
	if existingUser, exists := s.users[user.ID]; exists {
		return existingUser, nil
	}

	// Create new user
	userID := s.generateID()
	user.ID = userID
	s.users[userID] = &user

	return &user, nil
}

// GenerateToken creates a new JWT token for a user
func (s *MemoryStore) GenerateToken(userID string) (string, error) {
	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateToken validates a JWT token and returns the claims
func (s *MemoryStore) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// CreateMealPlan creates a new meal plan
func (s *MemoryStore) CreateMealPlan(plan *models.MealPlan) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if plan.ID == "" {
		plan.ID = s.generateID()
	}
	s.mealPlans[plan.ID] = plan
	return nil
}

// GetMealPlan retrieves a meal plan by ID
func (s *MemoryStore) GetMealPlan(id string) (*models.MealPlan, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	plan, exists := s.mealPlans[id]
	if !exists {
		return nil, ErrMealPlanNotFound
	}
	return plan, nil
}

// UpdateMealPlan updates an existing meal plan
func (s *MemoryStore) UpdateMealPlan(plan *models.MealPlan) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	existingPlan, exists := s.mealPlans[plan.ID]
	if !exists {
		return ErrMealPlanNotFound
	}

	existingPlan.Name = plan.Name
	existingPlan.Description = plan.Description
	existingPlan.UpdatedAt = time.Now()

	s.mealPlans[plan.ID] = existingPlan
	return nil
}

// DeleteMealPlan removes a meal plan from the store
func (s *MemoryStore) DeleteMealPlan(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.mealPlans[id]; !exists {
		return ErrMealPlanNotFound
	}

	delete(s.mealPlans, id)
	return nil
}

// ListMealPlansByUser returns all meal plans a user has access to
func (s *MemoryStore) ListMealPlansByUser(userID string) []*models.MealPlan {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var plans []*models.MealPlan
	for _, access := range s.mealPlanAccess {
		if access.UserID == userID {
			if plan, exists := s.mealPlans[access.MealPlanID]; exists {
				plans = append(plans, plan)
			}
		}
	}
	return plans
}

// CreateMealPlanAccess creates a new meal plan access record
func (s *MemoryStore) CreateMealPlanAccess(access *models.MealPlanAccess) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if access.ID == "" {
		access.ID = s.generateID()
	}
	s.mealPlanAccess[access.ID] = access
	return nil
}

// CheckMealPlanAccess checks if a user has access to a meal plan
// If checkOwner is true, it only returns true if the user is the owner
func (s *MemoryStore) CheckMealPlanAccess(userID, mealPlanID string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// First check if the user created the meal plan (which makes them an owner)
	plan, exists := s.mealPlans[mealPlanID]
	if exists && plan.CreatedBy == userID {
		return true, nil
	}

	// Then check for explicit access grants
	for _, access := range s.mealPlanAccess {
		if access.UserID == userID && access.MealPlanID == mealPlanID {
			return true, nil
		}
	}

	return false, ErrAccessDenied
}

// CheckMealPlanOwnership checks if a user is the owner of a meal plan
func (s *MemoryStore) CheckMealPlanOwnership(userID, mealPlanID string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Check if the user created the meal plan
	plan, exists := s.mealPlans[mealPlanID]
	if exists && plan.CreatedBy == userID {
		return true, nil
	}

	// Check for explicit owner access
	for _, access := range s.mealPlanAccess {
		if access.UserID == userID && access.MealPlanID == mealPlanID && access.Role == "owner" {
			return true, nil
		}
	}

	return false, ErrAccessDenied
}

// GetUserByGoogleID retrieves a user by their Google ID
func (s *MemoryStore) GetUserByGoogleID(googleID string) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, user := range s.users {
		if user.ID == googleID {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// GetUserByID retrieves a user by their ID
func (s *MemoryStore) GetUserByID(id string) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// CreateOrUpdateUser creates a new user or updates an existing one
func (s *MemoryStore) CreateOrUpdateUser(user *models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if user exists by ID
	if existingUser, exists := s.users[user.ID]; exists {
		// Update existing user
		existingUser.Email = user.Email
		existingUser.Name = user.Name
		s.users[user.ID] = existingUser
		return nil
	}

	// Create new user
	s.users[user.ID] = user
	return nil
}

// GetUserByEmail returns a user by email
func (s *MemoryStore) GetUserByEmail(email string) (*models.User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Find user by email
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// CreateShareLink creates a new share link
func (s *MemoryStore) CreateShareLink(link *models.ShareLink) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if link.ID == "" {
		link.ID = s.generateID()
	}
	s.shareLinks[link.ID] = link
	return nil
}

// GetShareLink retrieves a share link by ID
func (s *MemoryStore) GetShareLink(id string) (*models.ShareLink, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	link, exists := s.shareLinks[id]
	if !exists {
		return nil, errors.New("share link not found")
	}
	return link, nil
}

// DeleteShareLink removes a share link from the store
func (s *MemoryStore) DeleteShareLink(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.shareLinks[id]; !exists {
		return errors.New("share link not found")
	}

	delete(s.shareLinks, id)
	return nil
}
