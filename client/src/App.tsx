import React, { useState, useEffect } from 'react';
import './App.css';
import axios from 'axios';
import AddMealForm from './components/AddMealForm';
import MealItem from './components/MealItem';
import LoginButton from './components/LoginButton';
import { Meal, Day, MealType } from './types/meal';

const DAYS: Day[] = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
const MEAL_TYPES: MealType[] = ['Breakfast', 'Lunch', 'Dinner'];

function App() {
  const [status, setStatus] = useState<string>('Loading...');
  const [meals, setMeals] = useState<Meal[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [showAddForm, setShowAddForm] = useState<boolean>(false);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [authError, setAuthError] = useState<string>('');

  // Check authentication status on load
  useEffect(() => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      setIsAuthenticated(true);
    }
  }, []);

  // Check server status and fetch meals
  useEffect(() => {
    const checkServerStatus = async (): Promise<void> => {
      try {
        const response = await axios.get<{ message: string }>('/api/health');
        setStatus('Server is running: ' + response.data.message);
        if (isAuthenticated) {
          fetchMeals();
        } else {
          setLoading(false);
        }
      } catch (error: any) {
        setStatus('Error connecting to server: ' + error.message);
        setLoading(false);
        setError('Could not connect to server. Please try again later.');
      }
    };

    checkServerStatus();
  }, [isAuthenticated]);

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

  // Handle successful login
  const handleLoginSuccess = (token: string): void => {
    setIsAuthenticated(true);
    setAuthError('');
    fetchMeals();
  };

  // Handle login failure
  const handleLoginFailure = (error: string): void => {
    setAuthError(error);
  };

  // Handle logout
  const handleLogout = (): void => {
    localStorage.removeItem('auth_token');
    delete axios.defaults.headers.common['Authorization'];
    setIsAuthenticated(false);
    setMeals([]);
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>My Meal Planner</h1>
        <div className="header-right">
          <p className="server-status">{status}</p>
          {isAuthenticated ? (
            <button className="logout-button" onClick={handleLogout}>
              Logout
            </button>
          ) : null}
        </div>
      </header>
      <main>
        {!isAuthenticated ? (
          <section className="login-container">
            <h2>Please login to access your meal planner</h2>
            {authError && <div className="error-message">{authError}</div>}
            <LoginButton 
              onLoginSuccess={handleLoginSuccess} 
              onLoginFailure={handleLoginFailure} 
            />
          </section>
        ) : (
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
        )}
      </main>
    </div>
  );
}

export default App;
