import { Meal } from '../features/meals/types';

const MEALS_KEY = 'meals';

export const getAllMeals = (): Meal[] => {
  const mealsJson = localStorage.getItem(MEALS_KEY);
  return mealsJson ? JSON.parse(mealsJson) : [];
};

export const addMeal = (meal: Meal): boolean => {
  try {
    const meals = getAllMeals();
    meals.push(meal);
    localStorage.setItem(MEALS_KEY, JSON.stringify(meals));
    return true;
  } catch (err) {
    console.error('Error adding meal:', err);
    return false;
  }
};

export const deleteMeal = (id: string): boolean => {
  try {
    const meals = getAllMeals();
    const filteredMeals = meals.filter(meal => meal.id !== id);
    localStorage.setItem(MEALS_KEY, JSON.stringify(filteredMeals));
    return true;
  } catch (err) {
    console.error('Error deleting meal:', err);
    return false;
  }
};

export const updateMeal = (updatedMeal: Meal): boolean => {
  try {
    const meals = getAllMeals();
    const index = meals.findIndex(meal => meal.id === updatedMeal.id);
    if (index !== -1) {
      meals[index] = updatedMeal;
      localStorage.setItem(MEALS_KEY, JSON.stringify(meals));
      return true;
    }
    return false;
  } catch (err) {
    console.error('Error updating meal:', err);
    return false;
  }
};

export const clearMealsForWeek = (weekStartDate: Date): boolean => {
  try {
    const meals = getAllMeals();
    const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
    const currentDayIndex = weekStartDate.getDay();
    const weekDays = [...days.slice(currentDayIndex), ...days.slice(0, currentDayIndex)];

    const filteredMeals = meals.filter(meal => !weekDays.includes(meal.day));
    localStorage.setItem(MEALS_KEY, JSON.stringify(filteredMeals));
    return true;
  } catch (err) {
    console.error('Error clearing meals for week:', err);
    return false;
  }
}; 