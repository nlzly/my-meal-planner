import React, { useState, useEffect } from 'react'; // Added useState, useEffect
import { Meal, Day, MealType, MealPlan } from "../features/meals/types";
import * as localMealService from "../services/localMealService"; // Added localMealService
import { fetchMeals } from "../features/meals/mealsApi"; // Added fetchMeals
import AddMealForm from "./AddMealForm";
import MealGrid from "./MealGrid";
import Modal from "./Modal";
import ConfirmModal from "./ConfirmModal";
import ShareLinkModal from './ShareLinkModal';

const DAYS: Day[] = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];
const MEAL_TYPES: MealType[] = ["Breakfast", "Lunch", "Dinner"];


// Define props interface
interface MealPlannerContainerProps {
  selectedMealPlanId: string;
  setSelectedMealPlanId: (id: string) => void;
  mealPlans: MealPlan[];
  setShowCreatePlanModal: (show: boolean) => void;
  setShowShareModal: (show: boolean) => void;
  showCreatePlanModal: boolean; // Keep modal state for creation here
  newPlanName: string;
  setNewPlanName: (name: string) => void;
  newPlanDescription: string;
  setNewPlanDescription: (desc: string) => void;
  handleCreateMealPlan: () => void; // Keep creation logic in App
  showShareModal: boolean;
}

const MealPlannerContainer: React.FC<MealPlannerContainerProps> = ({
  selectedMealPlanId,
  setSelectedMealPlanId,
  mealPlans,
  setShowCreatePlanModal,
  setShowShareModal,
  showCreatePlanModal,
  newPlanName,
  setNewPlanName,
  newPlanDescription,
  setNewPlanDescription,
  handleCreateMealPlan,
  showShareModal,
}) => {
  // State moved from App.tsx
  const [meals, setMeals] = useState<Meal[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>("");
  const [isModalOpen, setIsModalOpen] = useState(false); // For Add/Edit Meal Modal
  const [selectedDay, setSelectedDay] = useState<Day>(DAYS[0]);
  const [selectedMealType, setSelectedMealType] = useState<MealType>(MEAL_TYPES[0]);
  const [mealToEdit, setMealToEdit] = useState<Meal | undefined>();
  const [showConfirmModal, setShowConfirmModal] = useState(false); // For Clear Plan Modal

  // Fetch meals when selectedMealPlanId changes
  useEffect(() => {
    if (selectedMealPlanId) {
      getMeals();
    } else {
      // Clear meals if no plan is selected
      setMeals([]);
      setLoading(false);
      setError("");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedMealPlanId]);

  // Functions moved from App.tsx
  const getMeals = async (): Promise<void> => {
    if (!selectedMealPlanId) {
      setMeals([]);
      setLoading(false);
      return;
    }
    setLoading(true); // Set loading true before fetch
    setError(""); // Clear previous errors
    fetchMeals(selectedMealPlanId).then( response => {
      if(response.status == 200) {
        setMeals(response.data);
      } else {
        setMeals([]);
        setError(response.error || "Failed to fetch meals.");
      }
      setLoading(false);
    }).catch(err => {
      console.error("Error fetching meals:", err);
      setMeals([]);
      setError("An unexpected error occurred while fetching meals.");
      setLoading(false);
    });
  };

  const handleMealAdded = (newMeal: Meal): void => {
    setMeals((prevMeals) => {
      if (mealToEdit) {
        // This assumes the backend/service handles the actual update
        // and the returned newMeal is the updated one.
        return prevMeals.map((meal) => meal.id === mealToEdit.id ? newMeal : meal);
      }
      // This assumes the backend/service handles the actual add
      // and the returned newMeal has the final ID.
      return [...prevMeals, newMeal];
    });
    // setShowAddForm(false); // Close inline form if open
    setIsModalOpen(false); // Close modal form
    setMealToEdit(undefined); // Clear editing state
  };

  const handleDeleteMeal = (mealId: string): void => {
    // TODO: Call backend API to delete meal
    // For now, using localMealService as placeholder if needed, or just update state
    // localMealService.deleteMeal(mealId); // Example if using local service
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
      // TODO: Call backend API to update meal day/type
      // localMealService.updateMeal(updatedMeal); // Example if using local service
      setMeals(prevMeals =>
        prevMeals.map(m => m.id === mealId ? updatedMeal : m)
      );
    }
  };

  const handleCopyMeal = (meal: Meal, newDay: Day, newMealType: MealType): void => {
    // Prepare data for backend API call, excluding fields generated by the server (id, createdAt, updatedAt)
    const newMealData = {
      name: meal.name,
      description: meal.description,
      day: newDay,
      mealType: newMealType,
      // Include other relevant fields from 'meal' that the backend expects for creation
    };

    // TODO: Call backend API to add the copied meal using newMealData
    // Example: const response = await axios.post(`/api/meal-plans/${selectedMealPlanId}/meals`, newMealData);
    // const actualNewMeal = response.data;

    // For now, simulate adding to state with a temp ID or handle after backend response
    const tempNewMeal = {
      ...newMealData,
      id: crypto.randomUUID(), // Temporary ID for UI update
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    // localMealService.addMeal(tempNewMeal); // Example if using local service
    setMeals(prevMeals => [...prevMeals, tempNewMeal]);
    // Ideally, replace tempNewMeal with the actual meal from backend response later
  };

  const getMealsForSlot = (day: Day, mealType: MealType): Meal[] => {
    if (!meals || !Array.isArray(meals)) {
      return [];
    }
    return meals.filter((meal) => meal.day === day && meal.mealType === mealType);
  };

  const handleAddMealClick = () => { // Renamed from handleAddMeal to avoid conflict
    setMealToEdit(undefined); // Ensure not editing
    setSelectedDay(DAYS[0]); // Reset selections for add modal/form
    setSelectedMealType(MEAL_TYPES[0]);
    // Decide whether to show inline form or modal
    // For simplicity, let's use the modal for adding too
    setIsModalOpen(true);
    // setShowAddForm(true); // Or use inline form
  };

  const handleClearMealPlanClick = () => { // Renamed from handleClearMealPlan
    setShowConfirmModal(true);
  };

  const handleConfirmClear = () => {
    // TODO: Call backend API to clear meals for the selected plan
    try {
      // const success = localMealService.clearMealsForWeek(weekStartDate); // Example if using local service
      // if (success) {
      //   setMeals([]);
      // }
      setMeals([]); // Optimistic update
    } catch (err) {
      console.error('Error clearing meals:', err);
      setError('Failed to clear meals. Please try again.');
    } finally {
      setShowConfirmModal(false);
    }
  };


  return (
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
            onClick={handleClearMealPlanClick} // Use renamed handler
            disabled={!selectedMealPlanId || meals.length === 0} // Also disable if no meals
          >
            Clear Meal Plan
          </button>
          <button
            className="add-button"
            onClick={handleAddMealClick} // Use renamed handler
            disabled={!selectedMealPlanId}
          >
            Add Meal
          </button>
        </div>
      </div>

      {/* Inline AddForm removed for simplicity, using modal for add/edit */}
      {/* {showAddForm && <AddMealForm onMealAdded={handleMealAdded} initialDay={DAYS[0]} initialMealType={MEAL_TYPES[0]} />} */}
      {error && <div className="error-message">{error}</div>}

      <div>
        {/* Modal for Adding/Editing Meals */}
        <Modal isOpen={isModalOpen} onClose={() => {
          setIsModalOpen(false);
          setMealToEdit(undefined); // Clear edit state on close
        }}>
          <AddMealForm
            onMealAdded={handleMealAdded} // Pass the local handler
            initialDay={selectedDay}
            initialMealType={selectedMealType}
            mealToEdit={mealToEdit} // Pass local edit state
            selectedMealPlanId={selectedMealPlanId} // Removed prop - AddMealForm doesn't expect it
          />
        </Modal>

        {/* Confirm Modal for Clearing Plan */}
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
  );
};

export default MealPlannerContainer;
