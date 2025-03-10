import React from 'react';
import axios from 'axios';
import { Meal } from '../types/meal';

interface MealItemProps {
  meal: Meal;
  onDelete: (id: string) => void;
  onUpdate?: (meal: Meal) => void;
}

function MealItem({ meal, onDelete, onUpdate }: MealItemProps) {
  const handleDelete = async (): Promise<void> => {
    if (window.confirm(`Are you sure you want to delete "${meal.name}"?`)) {
      try {
        await axios.delete(`/api/meals/${meal.id}`);
        if (onDelete) {
          onDelete(meal.id);
        }
      } catch (err) {
        console.error('Error deleting meal:', err);
        alert('Failed to delete meal. Please try again.');
      }
    }
  };

  return (
    <div className="meal-item">
      <div className="meal-header">
        <h4>{meal.name}</h4>
        <div className="meal-actions">
          <button 
            className="delete-button" 
            onClick={handleDelete}
            aria-label="Delete meal"
          >
            ×
          </button>
        </div>
      </div>
      {meal.description && (
        <p className="meal-description">{meal.description}</p>
      )}
    </div>
  );
}

export default MealItem;
