import React from 'react';
import * as localMealService from '../services/localMealService';
import { Meal } from '../features/meals/types';

interface MealItemProps {
  meal: Meal;
  onDelete: (id: string) => void;
  onUpdate?: (meal: Meal) => void;
  onDragStart?: (e: React.DragEvent<HTMLDivElement>, meal: Meal) => void;
  onDragEnd?: (e: React.DragEvent<HTMLDivElement>) => void;
  onMouseEnter?: () => void;
  onMouseLeave?: () => void;
}

function MealItem({ meal, onDelete, onUpdate, onDragStart, onDragEnd, onMouseEnter, onMouseLeave }: MealItemProps) {
  const handleDelete = (): void => {
      try {
        const success = localMealService.deleteMeal(meal.id);
        if (success && onDelete) {
          onDelete(meal.id);
        }
      } catch (err) {
        console.error('Error deleting meal:', err);
        alert('Failed to delete meal. Please try again.');
      }
  };

  const handleEdit = (): void => {
    if (onUpdate) {
      onUpdate(meal);
    }
  };

  return (
    <div 
      className="meal-item"
      draggable
      onDragStart={(e) => onDragStart?.(e, meal)}
      onDragEnd={onDragEnd}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <div className="meal-header">
        <h4>{meal.name}</h4>
        <div className="meal-actions">
          <button 
            className="edit-button" 
            onClick={handleEdit}
            aria-label="Edit meal"
          >
            ✎
          </button>
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
      {meal.chef && (
        <p className="meal-chef">Chef: {meal.chef}</p>
      )}
    </div>
  );
}

export default MealItem; 