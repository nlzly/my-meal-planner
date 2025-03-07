package db

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"my-meal-planner/models"
)

var (
	ErrMealNotFound = errors.New("meal not found")
)

// Store is an in-memory database for meals
type Store struct {
	meals  map[string]models.Meal
	mutex  sync.RWMutex
	nextID int
}

// NewStore creates a new in-memory store
func NewStore() *Store {
	return &Store{
		meals:  make(map[string]models.Meal),
		nextID: 1,
	}
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
	id := s.nextID
	s.nextID++
	return strconv.Itoa(id)
}
