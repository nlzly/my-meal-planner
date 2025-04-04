import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import axios from "axios";
import * as localMealService from "./services/localMealService";
import AddMealForm from "./components/AddMealForm";
import MealGrid from "./components/MealGrid";
import LoginButton from "./components/LoginButton";
import Modal from "./components/Modal";
import ConfirmModal from "./components/ConfirmModal";
import JoinMealPlan from "./components/JoinMealPlan";
import { Meal, Day, MealType } from "./types/meal";
import "./App.css";
import ShareLinkModal from './components/ShareLinkModal';

const DAYS: Day[] = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];
const MEAL_TYPES: MealType[] = ["Breakfast", "Lunch", "Dinner"];

interface User {
  id: string;
  name: string;
  email: string;
}

interface MealPlan {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
  createdBy: string;
}

function App() {
  const [status, setStatus] = useState<string>("Loading...");
  const [meals, setMeals] = useState<Meal[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");
  const [showAddForm, setShowAddForm] = useState<boolean>(false);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [authError, setAuthError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedDay, setSelectedDay] = useState<Day>(DAYS[0]);
  const [selectedMealType, setSelectedMealType] = useState<MealType>(MEAL_TYPES[0]);
  const [mealToEdit, setMealToEdit] = useState<Meal | undefined>();
  const [weekStartDate] = useState(new Date());
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [mealPlans, setMealPlans] = useState<MealPlan[]>([]);
  const [selectedMealPlanId, setSelectedMealPlanId] = useState<string>("");
  const [showCreatePlanModal, setShowCreatePlanModal] = useState(false);
  const [newPlanName, setNewPlanName] = useState("");
  const [newPlanDescription, setNewPlanDescription] = useState("");
  const [showShareModal, setShowShareModal] = useState(false);

  useEffect(() => {
    // Check URL for token parameter and return URL
    const queryParams = new URLSearchParams(window.location.search);
    const token = queryParams.get('token');
    const returnTo = queryParams.get('returnTo');

    if (token) {
      // Store token in localStorage
      localStorage.setItem('token', token);
      
      // Set auth header for axios
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      
      // Update auth state
      setIsAuthenticated(true);
      
      // Clear query parameters
      const newUrl = returnTo || window.location.pathname;
      window.history.replaceState({}, document.title, newUrl);
    } else {
      // Check if we have a token in localStorage
      const storedToken = localStorage.getItem('token');
      if (storedToken) {
        axios.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`;
        setIsAuthenticated(true);
      }
    }

    fetchMealPlans();
  }, []);

  useEffect(() => {
    if (selectedMealPlanId) {
      fetchMeals();
    }
  }, [selectedMealPlanId]);

  const fetchMealPlans = async () => {
    try {
      const response = await axios.get<MealPlan[]>("/api/meal-plans");
      if (response.data && response.data.length > 0) {
        setMealPlans(response.data);
        setSelectedMealPlanId(response.data[0].id);
      } else {
        // Handle empty response
        setMealPlans([]);
        setStatus("No meal plans found. Create your first meal plan.");
      }
    } catch (error) {
      console.error("Error fetching meal plans:", error);
      setError("Failed to load meal plans. Please try again.");
      setMealPlans([]); // Ensure mealPlans is always an array
    }
  };

  const fetchMeals = async (): Promise<void> => {
    if (!selectedMealPlanId) {
      setMeals([]);
      setLoading(false);
      return;
    }
    
    try {
      const response = await axios.get<Meal[]>(`/api/meals?mealPlanId=${selectedMealPlanId}`);
      setMeals(response.data || []);
      setLoading(false);
    } catch (error: any) {
      console.error("Error fetching meals:", error);
      setError("Failed to load meals. Please try again.");
      setLoading(false);
      setMeals([]); // Ensure meals is always an array
    }
  };

  const handleCreateMealPlan = async () => {
    try {
      const response = await axios.post<MealPlan>("/api/meal-plans", {
        name: newPlanName,
        description: newPlanDescription,
      });
      setMealPlans(prev => [...prev, response.data]);
      setSelectedMealPlanId(response.data.id);
      setShowCreatePlanModal(false);
      setNewPlanName("");
      setNewPlanDescription("");
    } catch (error) {
      console.error("Error creating meal plan:", error);
      setError("Failed to create meal plan. Please try again.");
    }
  };

  const handleMealAdded = (newMeal: Meal): void => {
    setMeals((prevMeals) => {
      if (mealToEdit) {
        return prevMeals.map((meal) => meal.id === mealToEdit.id ? newMeal : meal);
      }
      return [...prevMeals, newMeal];
    });
    setShowAddForm(false);
    setIsModalOpen(false);
    setMealToEdit(undefined);
  };

  const handleDeleteMeal = (mealId: string): void => {
    localMealService.deleteMeal(mealId);
    setMeals((prevMeals) => prevMeals.filter((meal) => meal.id !== mealId));
  };

  const handleUpdateMeal = (updatedMeal: Meal): void => {
    setMealToEdit(updatedMeal);
    setIsModalOpen(true);
    setSelectedDay(updatedMeal.day as Day);
    setSelectedMealType(updatedMeal.mealType as MealType);
  };

  const handleMoveMeal = (mealId: string, newDay: Day, newMealType: MealType): void => {
    const meal = meals.find(m => m.id === mealId);
    if (meal) {
      const updatedMeal = {
        ...meal,
        day: newDay,
        mealType: newMealType
      };
      localMealService.updateMeal(updatedMeal);
      setMeals(prevMeals => 
        prevMeals.map(m => m.id === mealId ? updatedMeal : m)
      );
    }
  };

  const handleCopyMeal = (meal: Meal, newDay: Day, newMealType: MealType): void => {
    const newMeal = {
      ...meal,
      id: crypto.randomUUID(), // Generate a new ID for the copy
      day: newDay,
      mealType: newMealType,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    localMealService.addMeal(newMeal);
    setMeals(prevMeals => [...prevMeals, newMeal]);
  };

  const getMealsForSlot = (day: Day, mealType: MealType): Meal[] => {
    if (!meals || !Array.isArray(meals)) {
      return [];
    }
    return meals.filter((meal) => meal.day === day && meal.mealType === mealType);
  };

  const handleLoginSuccess = (): void => {
    setIsAuthenticated(true);
    setAuthError("");
    fetchMeals();
  };

  const handleLoginFailure = (error: string): void => {
    setAuthError(error);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    delete axios.defaults.headers.common['Authorization'];
    setIsAuthenticated(false);
    setMeals([]);
  };

  const handleAddMeal = () => {
    setShowAddForm(true);
  };

  const handleClearMealPlan = () => {
    setShowConfirmModal(true);
  };

  const handleConfirmClear = () => {
    try {
      const success = localMealService.clearMealsForWeek(weekStartDate);
      if (success) {
        setMeals([]);
      }
    } catch (err) {
      console.error('Error clearing meals:', err);
      alert('Failed to clear meals. Please try again.');
    } finally {
      setShowConfirmModal(false);
    }
  };

  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <h1>My Meal Planner</h1>
          <div className="header-right">
            <p className="server-status">{status}</p>
            {isAuthenticated && (
              <button className="logout-button" onClick={handleLogout}>
                Logout
              </button>
            )}
          </div>
        </header>
        <main>
          {!isAuthenticated ? (
            <section className="login-container">
              <h2>Please login to access your meal planner</h2>
              {authError && <div className="error-message">{authError}</div>}
              <LoginButton />
            </section>
          ) : (
            <Routes>
              <Route 
                path="/join" 
                element={
                  <JoinMealPlan 
                    setSelectedMealPlanId={setSelectedMealPlanId} 
                    refreshMealPlans={fetchMealPlans} 
                  />
                } 
              />
              <Route 
                path="/" 
                element={
                  <section className="meal-planner-container">
                    <div className="meal-planner-header">
                      <h2>Weekly Meal Plan</h2>
                      <div className="header-buttons">
                        <select 
                          className="meal-plan-select"
                          value={selectedMealPlanId}
                          onChange={(e) => setSelectedMealPlanId(e.target.value)}
                        >
                          {mealPlans && mealPlans.length > 0 ? (
                            mealPlans.map(plan => (
                              <option key={plan.id} value={plan.id}>
                                {plan.name}
                              </option>
                            ))
                          ) : (
                            <option value="">No meal plans available</option>
                          )}
                        </select>
                        <button className="create-plan-button" onClick={() => setShowCreatePlanModal(true)}>
                          Create New Plan
                        </button>
                        {selectedMealPlanId && (
                          <button 
                            className="share-button" 
                            onClick={() => setShowShareModal(true)}
                          >
                            Share Plan
                          </button>
                        )}
                        <button 
                          className="clear-button" 
                          onClick={handleClearMealPlan}
                          disabled={!selectedMealPlanId}
                        >
                          Clear Meal Plan
                        </button>
                        <button 
                          className="add-button" 
                          onClick={handleAddMeal}
                          disabled={!selectedMealPlanId}
                        >
                          Add Meal
                        </button>
                      </div>
                    </div>

                    {showAddForm && <AddMealForm onMealAdded={handleMealAdded} initialDay={DAYS[0]} initialMealType={MEAL_TYPES[0]} />}
                    {error && <div className="error-message">{error}</div>}

                    <div>
                      <Modal isOpen={isModalOpen} onClose={() => {
                        setIsModalOpen(false);
                        setMealToEdit(undefined);
                      }}>
                        <AddMealForm 
                          onMealAdded={handleMealAdded} 
                          initialDay={selectedDay} 
                          initialMealType={selectedMealType}
                          mealToEdit={mealToEdit}
                        />
                      </Modal>

                      <ConfirmModal
                        isOpen={showConfirmModal}
                        onClose={() => setShowConfirmModal(false)}
                        onConfirm={handleConfirmClear}
                        title="Clear Meal Plan"
                        message="Are you sure you want to clear all meals for this week? This action cannot be undone."
                      />
                    </div>

                    <Modal isOpen={showCreatePlanModal} onClose={() => setShowCreatePlanModal(false)}>
                      <div className="create-plan-form">
                        <h2>Create New Meal Plan</h2>
                        <div className="form-group">
                          <label htmlFor="planName">Plan Name</label>
                          <input
                            type="text"
                            id="planName"
                            value={newPlanName}
                            onChange={(e) => setNewPlanName(e.target.value)}
                            placeholder="Enter plan name"
                          />
                        </div>
                        <div className="form-group">
                          <label htmlFor="planDescription">Description</label>
                          <textarea
                            id="planDescription"
                            value={newPlanDescription}
                            onChange={(e) => setNewPlanDescription(e.target.value)}
                            placeholder="Enter plan description"
                          />
                        </div>
                        <div className="form-buttons">
                          <button className="cancel-button" onClick={() => setShowCreatePlanModal(false)}>
                            Cancel
                          </button>
                          <button 
                            className="submit-button"
                            onClick={handleCreateMealPlan}
                            disabled={!newPlanName.trim()}
                          >
                            Create Plan
                          </button>
                        </div>
                      </div>
                    </Modal>

                    {loading ? (
                      <div className="loading">Loading meals...</div>
                    ) : (
                      <MealGrid
                        days={DAYS}
                        mealTypes={MEAL_TYPES}
                        getMealsForSlot={getMealsForSlot}
                        onDeleteMeal={handleDeleteMeal}
                        onUpdateMeal={handleUpdateMeal}
                        onMoveMeal={handleMoveMeal}
                        onCopyMeal={handleCopyMeal}
                        openModal={(day, mealType) => {
                          setSelectedDay(day);
                          setSelectedMealType(mealType);
                          setIsModalOpen(true);
                        }}
                      />
                    )}

                    {/* Share Meal Plan Modal */}
                    {selectedMealPlanId && (
                      <ShareLinkModal
                        isOpen={showShareModal}
                        onClose={() => setShowShareModal(false)}
                        mealPlanId={selectedMealPlanId}
                        mealPlanName={mealPlans.find(p => p.id === selectedMealPlanId)?.name || ''}
                      />
                    )}
                  </section>
                } 
              />
              <Route path="*" element={<Navigate to="/" />} />
            </Routes>
          )}
        </main>
      </div>
    </Router>
  );
}

export default App;
