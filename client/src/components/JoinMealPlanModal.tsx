import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate, useLocation } from 'react-router-dom';
import api from '../services/axios';
import Modal from './Modal';

interface JoinMealPlanModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const JoinMealPlanModal: React.FC<JoinMealPlanModalProps> = ({ 
  isOpen, 
  onClose 
}) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);
  const [mealPlanName, setMealPlanName] = useState<string>('');
  const [mealPlanId, setMealPlanId] = useState<string>('');
  const [code, setCode] = useState<string>('');
  const navigate = useNavigate();
  const location = useLocation();

  const handleJoin = async () => {
    const response = await api.post('/api/meal-plans/join', { code });
    setMealPlanName(response.data.mealPlan.name);
    setMealPlanId(response.data.mealPlan.id);
  }
  // useEffect(() => {
  //   const joinMealPlan = async () => {
  //     setLoading(true);
  //     setError(null);
  //     console.log(mealPlanId)
  //     // Get the code from the URL query parameter
  //     const queryParams = new URLSearchParams(location.search);
  //     const code = queryParams.get('code');
      
  //     if (!code) {
  //       setError('Invalid link. No share code provided.');
  //       setLoading(false);
  //       return;
  //     }
      
  //     try {
  //       // Check if user is authenticated
  //       const token = localStorage.getItem('token');
  //       if (!token) {
  //         // Redirect to login page with a return URL
  //         navigate(`/?returnTo=${encodeURIComponent(location.pathname + location.search)}`);
  //         return;
  //       }
        
  //       // Set the authorization header for axios
  //       axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        
  //       // Join the meal plan
  //       const response = await api.post('/api/meal-plans/join', { code });
        
  //       // Handle successful join
  //       setSuccess(true);
  //       setMealPlanName(response.data.mealPlan.name);
  //       setMealPlanId(response.data.mealPlan.id);
  //     } catch (err: any) {
  //       setError(err.response?.data || 'Failed to join meal plan');
  //     } finally {
  //       setLoading(false);
  //     }
  //   };
    
  //   joinMealPlan();
  // }, [location, navigate]);
  
  
  // if (loading) {
  //   return (
  //     <div className="join-meal-plan">
  //       <h2>Joining Meal Plan...</h2>
  //       <div className="loading-spinner"></div>
  //     </div>
  //   );
  // }
  
  // if (error) {
  //   return (
  //     <div className="join-meal-plan">
  //       <h2>Error Joining Meal Plan</h2>
  //       <div className="error-message">{error}</div>
  //       <button 
  //         className="return-button"
  //         onClick={() => navigate('/')}
  //       >
  //         Return to Home
  //       </button>
  //     </div>
  //   );
  // }
  
  // if (success) {
  //   return (
  //     <div className="join-meal-plan">
  //       <h2>Successfully Joined Meal Plan</h2>
  //       <div className="success-message">
  //         You have successfully joined "{mealPlanName}".
  //       </div>
  //       <p>You will be redirected to the meal plan shortly...</p>
  //       <button 
  //         className="continue-button"
  //         onClick={() => navigate('/')}
  //       >
  //         Continue to Meal Plan
  //       </button>
  //     </div>
  //   );
  // }
  
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className="join-meal-plan">
        <h2>Join Meal Plan</h2>
        <input type="text"
          id="name"
          name="name"
          value={code}
          onChange={(e) => setCode(e.target.value)}
          required/>
          <button onClick={handleJoin}>
            Join
          </button>
    </div>
    </Modal>
  );
};

export default JoinMealPlanModal; 