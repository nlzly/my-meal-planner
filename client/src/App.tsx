import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import axios from "axios";
import * as localMealService from "./services/localMealService";
import LoginButton from "./components/LoginButton";
import JoinMealPlan from "./components/JoinMealPlan";
import MealPlannerContainer from "./components/MealPlannerContainer"; // Import the new component
import { Meal, Day, MealType, MealPlan } from "./features/meals/types";
import "./App.css";
import { fetchMealPlans, fetchMeals } from "./features/meals/mealsApi";


function App() {
  const [status, setStatus] = useState<string>("Loading...");
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [authError, setAuthError] = useState<string | null>(null);
  const [weekStartDate] = useState(new Date()); // Keep for now, might be needed elsewhere or passed down
  const [mealPlans, setMealPlans] = useState<MealPlan[]>([]);
  const [selectedMealPlanId, setSelectedMealPlanId] = useState<string>(""); // Keep: Used for selection
  const [showCreatePlanModal, setShowCreatePlanModal] = useState(false); // Keep: Controls modal in App scope
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

     fetchMealPlans().then( response => {
        if(response.status == 200) {
          setMealPlans(response.data)
          setStatus("")
        } else {
          setMealPlans([])
          setStatus(response.error)
        }
     });
    
  }, []);

  const handleCreateMealPlan = async () => { // Keep: Manages meal plan list and creation modal
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
      // setError("Failed to create meal plan. Please try again."); // Error state moved
      setStatus("Failed to create meal plan."); // Use status for general messages
    }
  };

  const handleLoginSuccess = (): void => { // Keep: Manages auth state
    setIsAuthenticated(true);
    setAuthError("");
    // getMeals(); // Removed: MealPlannerContainer fetches its own meals
  };

  const handleLoginFailure = (error: string): void => { // Keep: Manages auth state
    setAuthError(error);
  };

  const handleLogout = () => { // Keep: Manages auth state
    localStorage.removeItem('token');
    delete axios.defaults.headers.common['Authorization'];
    setIsAuthenticated(false);
    // setMeals([]); // meals state moved
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
                  <MealPlannerContainer
                    selectedMealPlanId={selectedMealPlanId}
                    setSelectedMealPlanId={setSelectedMealPlanId} // Still needed to change plan
                    mealPlans={mealPlans} // List of plans
                    setShowCreatePlanModal={setShowCreatePlanModal} // Show create modal
                    setShowShareModal={setShowShareModal} // Show share modal
                    showCreatePlanModal={showCreatePlanModal} // Create modal state
                    newPlanName={newPlanName} // Create modal state
                    setNewPlanName={setNewPlanName} // Create modal state
                    newPlanDescription={newPlanDescription} // Create modal state
                    setNewPlanDescription={setNewPlanDescription} // Create modal state
                    handleCreateMealPlan={handleCreateMealPlan} // Create plan handler
                    showShareModal={showShareModal} // Share modal state
                  />
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
