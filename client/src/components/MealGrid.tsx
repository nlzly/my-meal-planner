import React from "react";
import MealItem from "./MealItem";
import { Meal, Day, MealType } from "../types/meal";

type MealGridProps = {
  days: Day[];
  mealTypes: MealType[];
  getMealsForSlot: (day: Day, mealType: MealType) => Meal[];
  onDeleteMeal: (mealId: string) => void;
  openModal: () => void;
};

const MealGrid: React.FC<MealGridProps> = ({
  days,
  mealTypes,
  getMealsForSlot,
  onDeleteMeal,
  openModal,
}) => {
  return (
    <div className="meal-days">
      {days.map((day) => (
        <div key={day} className="meal-day">
          <h3>{day}</h3>
          <div className="meal-slots">
            {mealTypes.map((mealType) => {
              const mealsForSlot = getMealsForSlot(day, mealType);
              return (
                <div key={mealType} className="meal-slot" onClick={openModal}>
                  <h4>{mealType}</h4>
                  {mealsForSlot.length > 0 ? (
                    mealsForSlot.map((meal) => (
                      <MealItem key={meal.id} meal={meal} onDelete={onDeleteMeal} />
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
