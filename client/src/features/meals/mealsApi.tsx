import { Meal, MealPlan } from "./types";
import api from "../../services/axios";

export interface MealPlanResponse {
    data : MealPlan[],
    error: string,
    status: number
}

export interface MealResponse {
    data : Meal[]
    error : string,
    status : number
}

export const fetchMealPlans = async (): Promise<MealPlanResponse> => {
    try {
      const response = await api.get<MealPlan[]>("/api/meal-plans");
        return {
            data: response.data?.length > 0 ? response.data : [],
            error: "",
            status: response.status
        }
    } catch (error) {
        console.error("Error fetching meal plans:", error);
        return {
            data: [],
            error: "Failed to load meal plans. Please try again.",
            status: 400
        }
    }
  };

  export const fetchMeals = async(selectedMealPlanId : string) : Promise<MealResponse> => {
    try {
        const response = await api.get<Meal[]>(`/api/meals?mealPlanId=${selectedMealPlanId}`)
        return {
            data: response.data?.length > 0 ? response.data : [],
            error: "",
            status: response.status
        }
    } catch(error) {
        console.error("Error fetching meal plans:", error);
        return {
            data: [],
            error: "Failed to load meals. Please try again.",
            status: 400
        }
    }
  }

  export const addMeal = async(meal : Meal, selectedMealPlanId : string) : Promise<void> => {
    try {
        await api.post<Meal>("/api/meals", {
            meal: meal,
            mealPlanId: selectedMealPlanId,
          });
    } catch (error) {
        console.error("Error adding meal:", error);
    }
  }

  export const deleteMeal = async(mealId : string) : Promise<void> => {
    try {
        await api.delete<Meal>(`/api/meals/${mealId}`);
    } catch (error) {
        console.error("Error deleting meal:", error);
    }
  }

  export const updateMeal = async(meal : Meal) : Promise<void> => {
    try {
        
        await api.put<Meal>(`/api/meals/${meal.id}`, {
            name:  meal.name,
            description : meal.description,
            mealType: meal.mealType,
            day: meal.day
          });
    } catch (error) {
        console.error("Error adding meal:", error);
    }
  }