import React, { useState, ChangeEvent, FormEvent } from 'react';
import axios from 'axios';
import { Meal, MealRequest, Day, MealType } from '../types/meal';

const DAYS: Day[] = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
const MEAL_TYPES: MealType[] = ['Breakfast', 'Lunch', 'Dinner'];

interface AddMealFormProps {
  onMealAdded: (meal: Meal) => void;
}

function AddMealForm({ onMealAdded }: AddMealFormProps) {
  const [formData, setFormData] = useState<MealRequest>({
    name: '',
    description: '',
    day: DAYS[0],
    mealType: MEAL_TYPES[0]
  });
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError('');

    try {
      const response = await axios.post<Meal>('/api/meals', formData);
      setFormData({
        name: '',
        description: '',
        day: formData.day,
        mealType: formData.mealType
      });
      if (onMealAdded) {
        onMealAdded(response.data);
      }
    } catch (err: any) {
      setError(err.response?.data || 'Failed to add meal. Please try again.');
      console.error('Error adding meal:', err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="add-meal-form">
      <h3>Add New Meal</h3>
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
          {isSubmitting ? 'Adding...' : 'Add Meal'}
        </button>
      </form>
    </div>
  );
}

export default AddMealForm;
