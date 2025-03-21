import React, { useState } from "react";
import MealItem from "./MealItem";
import { Meal, Day, MealType } from "../types/meal";

type MealGridProps = {
  days: Day[];
  mealTypes: MealType[];
  getMealsForSlot: (day: Day, mealType: MealType) => Meal[];
  onDeleteMeal: (mealId: string) => void;
  onUpdateMeal: (meal: Meal) => void;
  openModal: (day: Day, mealType: MealType) => void;
  onMoveMeal: (mealId: string, newDay: Day, newMealType: MealType) => void;
};

const MealGrid: React.FC<MealGridProps> = ({
  days,
  mealTypes,
  getMealsForSlot,
  onDeleteMeal,
  onUpdateMeal,
  openModal,
  onMoveMeal,
}) => {
  const [draggedMeal, setDraggedMeal] = useState<Meal | null>(null);

  const handleDragStart = (e: React.DragEvent<HTMLDivElement>, meal: Meal) => {
    setDraggedMeal(meal);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragEnd = () => {
    setDraggedMeal(null);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>, day: Day, mealType: MealType) => {
    e.preventDefault();
    if (draggedMeal && (draggedMeal.day !== day || draggedMeal.mealType !== mealType)) {
      onMoveMeal(draggedMeal.id, day, mealType);
    }
  };

  return (
    <div className="meal-days">
      {days.map((day) => (
        <div key={day} className="meal-day">
          <h3>{day}</h3>
          <div className="meal-slots">
            {mealTypes.map((mealType) => {
              const mealsForSlot = getMealsForSlot(day, mealType);
              return (
                <div 
                  key={mealType} 
                  className="meal-slot" 
                  onClick={mealsForSlot.length < 1 ? () => openModal(day, mealType) : undefined}
                  onDragOver={handleDragOver}
                  onDrop={(e) => handleDrop(e, day, mealType)}
                >
                  <h4>{mealType}</h4>
                  {mealsForSlot.length > 0 ? (
                    mealsForSlot.map((meal) => (
                      <MealItem 
                        key={meal.id} 
                        meal={meal} 
                        onDelete={onDeleteMeal}
                        onUpdate={onUpdateMeal}
                        onDragStart={handleDragStart}
                        onDragEnd={handleDragEnd}
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
  );
};

export default MealGrid;
