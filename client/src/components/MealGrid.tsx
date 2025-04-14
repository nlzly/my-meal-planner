import React, { useState, useEffect, useRef } from "react";
import MealItem from "./MealItem";
import { Meal, Day, MealType } from "../types/types";

type MealGridProps = {
  days: Day[];
  mealTypes: MealType[];
  getMealsForSlot: (day: Day, mealType: MealType) => Meal[];
  onDeleteMeal: (mealId: string) => void;
  onUpdateMeal: (meal: Meal) => void;
  openModal: (day: Day, mealType: MealType) => void;
  onMoveMeal: (mealId: string, newDay: Day, newMealType: MealType) => void;
  onCopyMeal: (meal: Meal, newDay: Day, newMealType: MealType) => void;
};

const MealGrid: React.FC<MealGridProps> = ({
  days,
  mealTypes,
  getMealsForSlot,
  onDeleteMeal,
  onUpdateMeal,
  openModal,
  onMoveMeal,
  onCopyMeal,
}) => {
  const [draggedMeal, setDraggedMeal] = useState<Meal | null>(null);
  const [copiedMeal, setCopiedMeal] = useState<Meal | null>(null);
  const [hoveredMeal, setHoveredMeal] = useState<Meal | null>(null);
  const [hoveredSlot, setHoveredSlot] = useState<{ day: Day; mealType: MealType } | null>(null);
  const gridRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Check if Ctrl/Cmd + C is pressed
      if ((e.ctrlKey || e.metaKey) && e.key === 'c' && hoveredMeal) {
        e.preventDefault();
        setCopiedMeal(hoveredMeal);
      }
      // Check if Ctrl/Cmd + V is pressed
      else if ((e.ctrlKey || e.metaKey) && e.key === 'v' && copiedMeal && hoveredSlot) {
        e.preventDefault();
        const { day, mealType } = hoveredSlot;
        onCopyMeal(copiedMeal, day, mealType);
      }
    };

    const handleClickOutside = (e: MouseEvent) => {
      if (gridRef.current && !gridRef.current.contains(e.target as Node)) {
        setCopiedMeal(null);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('click', handleClickOutside);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
      window.removeEventListener('click', handleClickOutside);
    };
  }, [hoveredMeal, copiedMeal, hoveredSlot, onCopyMeal]);

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

  const handleSlotClick = (day: Day, mealType: MealType, mealsForSlot: Meal[]) => {
    if (mealsForSlot.length < 1) {
      openModal(day, mealType);
    }
  };

  return (
    <div className="meal-days" ref={gridRef}>
      {days.map((day) => (
        <div key={day} className="meal-day">
          <h3>{day}</h3>
          <div className="meal-slots">
            {mealTypes.map((mealType) => {
              const mealsForSlot = getMealsForSlot(day, mealType);
              const isHovered = hoveredSlot?.day === day && hoveredSlot?.mealType === mealType;
              const canPaste = copiedMeal;
              return (
                <div 
                  key={mealType} 
                  className={`meal-slot ${isHovered && canPaste ? 'paste-target' : ''}`}
                  onClick={() => handleSlotClick(day, mealType, mealsForSlot)}
                  onDragOver={handleDragOver}
                  onDrop={(e) => handleDrop(e, day, mealType)}
                  onMouseEnter={() => setHoveredSlot({ day, mealType })}
                  onMouseLeave={() => setHoveredSlot(null)}
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
                        onMouseEnter={() => setHoveredMeal(meal)}
                        onMouseLeave={() => setHoveredMeal(null)}
                      />
                    ))
                  ) : (
                    <p className="no-meal">
                      {isHovered && canPaste ? 'Press Ctrl+V to paste' : 'No meal planned'}
                    </p>
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