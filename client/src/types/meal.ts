export interface Meal {
  id: string;
  name: string;
  description: string;
  day: string;
  mealType: string; // 'Breakfast', 'Lunch', 'Dinner'
  createdAt: string;
  updatedAt: string;
  chef: string;
}

export interface MealRequest {
  name: string;
  description: string;
  day: string;
  mealType: string;
  chef: string;
}

export type Day = 'Monday' | 'Tuesday' | 'Wednesday' | 'Thursday' | 'Friday' | 'Saturday' | 'Sunday';
export type MealType = 'Breakfast' | 'Lunch' | 'Dinner';
