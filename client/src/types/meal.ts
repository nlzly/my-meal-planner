export interface Meal {
  id: string;
  name: string;
  description: string;
  day: string;
  mealType: string; // 'Breakfast', 'Lunch', 'Dinner'
  createdAt: string;
  updatedAt: string;
}

export interface MealRequest {
  name: string;
  description: string;
  day: string;
  mealType: string;
}

export type Day = 'Monday' | 'Tuesday' | 'Wednesday' | 'Thursday' | 'Friday' | 'Saturday' | 'Sunday';
export type MealType = 'Breakfast' | 'Lunch' | 'Dinner';
