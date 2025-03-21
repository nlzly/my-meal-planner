import { Meal, MealRequest } from '../types/meal';
import { v4 as uuidv4 } from 'uuid';

// Local storage key
const MEALS_STORAGE_KEY = 'local_meals';

// Helper to get current date in ISO format
const getCurrentDate = (): string => {
  return new Date().toISOString();
};

// Load meals from local storage
export const loadMeals = (): Meal[] => {
  const storedMeals = localStorage.getItem(MEALS_STORAGE_KEY);
  return storedMeals ? JSON.parse(storedMeals) : [];
};

// Save meals to local storage
export const saveMeals = (meals: Meal[]): void => {
  localStorage.setItem(MEALS_STORAGE_KEY, JSON.stringify(meals));
};

// Add a new meal
export const addMeal = (mealRequest: MealRequest): Meal => {
  const meals = loadMeals();
  
  const newMeal: Meal = {
    id: uuidv4(),
    ...mealRequest,
    createdAt: getCurrentDate(),
    updatedAt: getCurrentDate()
  };
  
  meals.push(newMeal);
  saveMeals(meals);
  
  return newMeal;
};

// Delete a meal
export const deleteMeal = (mealId: string): boolean => {
  const meals = loadMeals();
  const updatedMeals = meals.filter(meal => meal.id !== mealId);
  
  if (updatedMeals.length < meals.length) {
    saveMeals(updatedMeals);
    return true;
  }
  
  return false;
};

// Get all meals
export const getAllMeals = (): Meal[] => {
  return loadMeals();
};

// Update a meal
export const updateMeal = (updatedMeal: Meal): boolean => {
  const meals = loadMeals();
  const mealIndex = meals.findIndex(meal => meal.id === updatedMeal.id);
  
  if (mealIndex !== -1) {
    meals[mealIndex] = {
      ...updatedMeal,
      updatedAt: getCurrentDate()
    };
    saveMeals(meals);
    return true;
  }
  
  return false;
};
