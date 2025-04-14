import axios from "axios";
import { MealPlan } from "../types/types";

export interface MealPlanResponse {
    data : MealPlan[],
    error: string,
    status: number
}

export const fetchMealPlans = async (): Promise<MealPlanResponse> => {
    try {
      const response = await axios.get<MealPlan[]>("/api/meal-plans");
        return {
            data: response.data,
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