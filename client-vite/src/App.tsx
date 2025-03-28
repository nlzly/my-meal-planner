import { useState, useEffect } from "react";
import axios from "axios";
import * as localMealService from "./services/localMealService";
import AddMealForm from "./components/AddMealForm";
import MealGrid from "./components/MealGrid";
import LoginButton from "./components/LoginButton";
import Modal from "./components/Modal";
import ConfirmModal from "./components/ConfirmModal";
import { Meal, Day, MealType } from "./types/meal";
import "./App.css";

const DAYS: Day[] = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];
const MEAL_TYPES: MealType[] = ["Breakfast", "Lunch", "Dinner"];

function App() {
  const [status, setStatus] = useState<string>("Loading...");
  const [meals, setMeals] = useState<Meal[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");
  const [showAddForm, setShowAddForm] = useState<boolean>(false);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(true);
  const [authError, setAuthError] = useState<string>("");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedDay, setSelectedDay] = useState<Day>(DAYS[0]);
  const [selectedMealType, setSelectedMealType] = useState<MealType>(MEAL_TYPES[0]);
  const [mealToEdit, setMealToEdit] = useState<Meal | undefined>();
  const [showAddMealModal, setShowAddMealModal] = useState(false);
  const [weekStartDate, setWeekStartDate] = useState(new Date());
  const [showConfirmModal, setShowConfirmModal] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("auth_token");
    if (token) {
      axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      setIsAuthenticated(true);
    }
  }, []);

  useEffect(() => {
    const checkServerStatus = async (): Promise<void> => {
      try {
        const response = await axios.get<{ message: string }>("/api/health");
        setStatus("Server is running: " + response.data.message);
      } catch (error: any) {
        setStatus("Using local storage mode (no server connection)");
      } finally {
        if (isAuthenticated) {
          fetchMeals();
        } else {
          setLoading(false);
        }
      }
    };

    checkServerStatus();
  }, [isAuthenticated]);

  const fetchMeals = (): void => {
    try {
      const localMeals = localMealService.getAllMeals();
      setMeals(localMeals);
      setLoading(false);
    } catch (error: any) {
      console.error("Error fetching meals:", error);
      setError("Failed to load meals. Please try again.");
      setLoading(false);
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
    return meals.filter((meal) => meal.day === day && meal.mealType === mealType);
  };

  const handleLoginSuccess = (token: string): void => {
    setIsAuthenticated(true);
    setAuthError("");
    fetchMeals();
  };

  const handleLoginFailure = (error: string): void => {
    setAuthError(error);
  };

  const handleLogout = (): void => {
    localStorage.removeItem("auth_token");
    delete axios.defaults.headers.common["Authorization"];
    setIsAuthenticated(false);
    setMeals([]);
  };

  const handleAddMeal = () => {
    setShowAddMealModal(true);
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
    <div className="App">
      <header className="App-header">
        <h1>My Meal Planner</h1>
        <div className="header-right">
          <p className="server-status">{status}</p>
        </div>
      </header>
      <main>
        {!isAuthenticated ? (
          <section className="login-container">
            <h2>Please login to access your meal planner</h2>
            {authError && <div className="error-message">{authError}</div>}
            <LoginButton onLoginSuccess={handleLoginSuccess} onLoginFailure={handleLoginFailure} />
          </section>
        ) : (
          <section className="meal-planner-container">
            <div className="meal-planner-header">
              <h2>Weekly Meal Plan</h2>
              <div className="header-buttons">
                <button className="clear-button" onClick={handleClearMealPlan}>
                  Clear Meal Plan
                </button>
                <button className="add-button" onClick={handleAddMeal}>
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
          </section>
        )}
      </main>
    </div>
  );
}

export default App;
