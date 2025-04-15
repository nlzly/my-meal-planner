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

export interface MealPlan {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
  createdBy: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
}