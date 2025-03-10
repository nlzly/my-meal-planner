import React, { useState, useEffect } from 'react';
import './App.css';
import axios from 'axios';
import AddMealForm from './components/AddMealForm';
import MealItem from './components/MealItem';
import { Meal, Day, MealType } from './types/meal';

const DAYS: Day[] = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
const MEAL_TYPES: MealType[] = ['Breakfast', 'Lunch', 'Dinner'];

function App() {
  const [status, setStatus] = useState<string>('Loading...');
  const [meals, setMeals] = useState<Meal[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [showAddForm, setShowAddForm] = useState<boolean>(false);

  // Check server status and fetch meals
  useEffect(() => {
    const checkServerStatus = async (): Promise<void> => {
      try {
        const response = await axios.get<{ message: string }>('/api/health');
        setStatus('Server is running: ' + response.data.message);
        fetchMeals();
      } catch (error: any) {
        setStatus('Error connecting to server: ' + error.message);
        setLoading(false);
        setError('Could not connect to server. Please try again later.');
      }
    };

    checkServerStatus();
  }, []);

  // Fetch all meals from the API
  const fetchMeals = async (): Promise<void> => {
    try {
      const response = await axios.get<Meal[]>('/api/meals');
      setMeals(response.data);
      setLoading(false);
    } catch (error: any) {
      console.error('Error fetching meals:', error);
      setError('Failed to load meals. Please try again.');
      setLoading(false);
    }
  };

  // Handle adding a new meal
  const handleMealAdded = (newMeal: Meal): void => {
    setMeals(prevMeals => [...prevMeals, newMeal]);
    setShowAddForm(false);
  };

  // Handle deleting a meal
  const handleDeleteMeal = (mealId: string): void => {
    setMeals(prevMeals => prevMeals.filter(meal => meal.id !== mealId));
  };

  // Get meals for a specific day and meal type
  const getMealsForSlot = (day: Day, mealType: MealType): Meal[] => {
    return meals.filter(meal => meal.day === day && meal.mealType === mealType);
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>My Meal Planner</h1>
        <p className="server-status">{status}</p>
      </header>
      <main>
        <section className="meal-planner-container">
          <div className="meal-planner-header">
            <h2>Weekly Meal Plan</h2>
            <button 
              className="add-meal-button"
              onClick={() => setShowAddForm(!showAddForm)}
            >
              {showAddForm ? 'Cancel' : 'Add Meal'}
            </button>
          </div>

          {showAddForm && (
            <AddMealForm onMealAdded={handleMealAdded} />
          )}

          {error && <div className="error-message">{error}</div>}

          {loading ? (
            <div className="loading">Loading meals...</div>
          ) : (
            <div className="meal-days">
              {DAYS.map(day => (
                <div key={day} className="meal-day">
                  <h3>{day}</h3>
                  <div className="meal-slots">
                    {MEAL_TYPES.map(mealType => {
                      const mealsForSlot = getMealsForSlot(day, mealType);
                      return (
                        <div key={mealType} className="meal-slot">
                          <h4>{mealType}</h4>
                          {mealsForSlot.length > 0 ? (
                            mealsForSlot.map(meal => (
                              <MealItem 
                                key={meal.id} 
                                meal={meal} 
                                onDelete={handleDeleteMeal} 
                              />
                            ))
                          ) : (
                            <p className="no-meal">No meal planned</p>
                          )}
                        </div>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}

export default App;
