import React, { useState, useEffect } from 'react';
import { Meal, Day, MealType, MealRequest } from '../types/meal';
import * as localMealService from '../services/localMealService';

interface AddMealFormProps {
  onMealAdded: (meal: Meal) => void;
  initialDay: Day;
  initialMealType: MealType;
  mealToEdit?: Meal;
}

const AddMealForm: React.FC<AddMealFormProps> = ({
  onMealAdded,
  initialDay,
  initialMealType,
  mealToEdit
}) => {
  const [formData, setFormData] = useState<MealRequest>({
    name: '',
    description: '',
    chef: '',
    day: initialDay,
    mealType: initialMealType
  });

  useEffect(() => {
    if (mealToEdit) {
      setFormData({
        name: mealToEdit.name,
        description: mealToEdit.description || '',
        chef: mealToEdit.chef || '',
        day: mealToEdit.day,
        mealType: mealToEdit.mealType
      });
    }
  }, [mealToEdit]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const meal: Meal = {
      ...formData,
      id: mealToEdit?.id || crypto.randomUUID(),
      createdAt: mealToEdit?.createdAt || new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };

    if (mealToEdit) {
      localMealService.updateMeal(meal);
    } else {
      localMealService.addMeal(meal);
    }

    onMealAdded(meal);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  return (
    <form onSubmit={handleSubmit} className="add-meal-form">
      <div className="form-group">
        <label htmlFor="name">Meal Name</label>
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
        <label htmlFor="description">Description</label>
        <textarea
          id="description"
          name="description"
          value={formData.description}
          onChange={handleChange}
        />
      </div>

      <div className="form-group">
        <label htmlFor="chef">Chef</label>
        <input
          type="text"
          id="chef"
          name="chef"
          value={formData.chef}
          onChange={handleChange}
        />
      </div>

      <div className="form-row">
        <div className="form-group">
          <label htmlFor="day">Day</label>
          <select
            id="day"
            name="day"
            value={formData.day}
            onChange={handleChange}
            required
          >
            <option value="Monday">Monday</option>
            <option value="Tuesday">Tuesday</option>
            <option value="Wednesday">Wednesday</option>
            <option value="Thursday">Thursday</option>
            <option value="Friday">Friday</option>
            <option value="Saturday">Saturday</option>
            <option value="Sunday">Sunday</option>
          </select>
        </div>

        <div className="form-group">
          <label htmlFor="mealType">Meal Type</label>
          <select
            id="mealType"
            name="mealType"
            value={formData.mealType}
            onChange={handleChange}
            required
          >
            <option value="Breakfast">Breakfast</option>
            <option value="Lunch">Lunch</option>
            <option value="Dinner">Dinner</option>
          </select>
        </div>
      </div>

      <button type="submit" className="submit-button">
        {mealToEdit ? 'Update Meal' : 'Add Meal'}
      </button>
    </form>
  );
};

export default AddMealForm; 