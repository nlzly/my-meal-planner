export type Day = "Monday" | "Tuesday" | "Wednesday" | "Thursday" | "Friday" | "Saturday" | "Sunday";
export type MealType = "Breakfast" | "Lunch" | "Dinner";

export interface Meal {
  id: string;
  name: string;
  description?: string;
  chef?: string;
  day: Day;
  mealType: MealType;
  createdAt: string;
  updatedAt: string;
}

export interface MealRequest {
  name: string;
  description?: string;
  chef?: string;
  day: Day;
  mealType: MealType;
} 