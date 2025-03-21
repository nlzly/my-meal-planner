import React, { useState, ChangeEvent, FormEvent, useEffect } from 'react';
import * as localMealService from '../services/localMealService';
import { Meal, MealRequest, Day, MealType } from '../types/meal';

const DAYS: Day[] = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
const MEAL_TYPES: MealType[] = ['Breakfast', 'Lunch', 'Dinner'];

interface AddMealFormProps {
  onMealAdded: (meal: Meal) => void;
  initialDay?: Day;
  initialMealType?: MealType;
  mealToEdit?: Meal;
}

function AddMealForm({ onMealAdded, initialDay = DAYS[0], initialMealType = MEAL_TYPES[0], mealToEdit }: AddMealFormProps) {
  const [formData, setFormData] = useState<MealRequest>({
    name: '',
    description: '',
    day: initialDay,
    mealType: initialMealType
  });
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    if (mealToEdit) {
      setFormData({
        name: mealToEdit.name,
        description: mealToEdit.description || '',
        day: mealToEdit.day,
        mealType: mealToEdit.mealType
      });
    }
  }, [mealToEdit]);

  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError('');

    try {
      let meal: Meal;
      if (mealToEdit) {
        const success = localMealService.updateMeal({
          ...mealToEdit,
          ...formData
        });
        if (!success) {
          throw new Error('Failed to update meal');
        }
        meal = {
          ...mealToEdit,
          ...formData
        };
      } else {
        meal = localMealService.addMeal(formData);
      }
      
      if (onMealAdded) {
        onMealAdded(meal);
      }
    } catch (err: any) {
      setError(mealToEdit ? 'Failed to update meal. Please try again.' : 'Failed to add meal. Please try again.');
      console.error('Error saving meal:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="add-meal-form">
      <h3>{mealToEdit ? 'Edit Meal' : 'Add New Meal'}</h3>
      {error && <div className="error-message">{error}</div>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="name">Meal Name:</label>
          <input
            type="text"
            id="name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="description">Description:</label>
          <textarea
            id="description"
            name="description"
            value={formData.description}
            onChange={handleChange}
            rows={3}
          />
        </div>

        <div className="form-row">
          <div className="form-group">
            <label htmlFor="day">Day:</label>
            <select
              id="day"
              name="day"
              value={formData.day}
              onChange={handleChange}
            >
              {DAYS.map(day => (
                <option key={day} value={day}>{day}</option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="mealType">Meal Type:</label>
            <select
              id="mealType"
              name="mealType"
              value={formData.mealType}
              onChange={handleChange}
            >
              {MEAL_TYPES.map(type => (
                <option key={type} value={type}>{type}</option>
              ))}
            </select>
          </div>
        </div>

        <button 
          type="submit" 
          className="submit-button"
          disabled={isSubmitting}
        >
          {isSubmitting ? (mealToEdit ? 'Updating...' : 'Adding...') : (mealToEdit ? 'Update Meal' : 'Add Meal')}
        </button>
      </form>
    </div>
  );
}

export default AddMealForm;
