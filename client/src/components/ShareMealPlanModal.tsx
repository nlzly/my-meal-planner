import React, { useState } from 'react';
import Modal from './Modal';
import api from '../services/axios';

interface ShareMealPlanModalProps {
  isOpen: boolean;
  onClose: () => void;
  mealPlanId: string;
  mealPlanName: string;
}

const ShareMealPlanModal: React.FC<ShareMealPlanModalProps> = ({
  isOpen,
  onClose,
  mealPlanId,
  mealPlanName,
}) => {
  const [email, setEmail] = useState('');
  const [role, setRole] = useState('viewer');
  const [error, setError] = useState('');
  const [isSharing, setIsSharing] = useState(false);
  const [shareSuccess, setShareSuccess] = useState(false);

  const handleShare = async () => {
    if (!email) {
      setError('Email is required');
      return;
    }

    setIsSharing(true);
    setError('');

    try {
      await api.post('/api/meal-plans/share', {
        mealPlanId,
        email,
        role,
      });
      
      setShareSuccess(true);
      setEmail('');
      setTimeout(() => {
        setShareSuccess(false);
      }, 3000);
    } catch (err: any) {
      setError(err.response?.data || 'Failed to share meal plan');
    } finally {
      setIsSharing(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className="share-meal-plan-modal">
        <h2>Share Meal Plan</h2>
        <p>Share "{mealPlanName}" with another user</p>
        
        {shareSuccess && (
          <div className="success-message">Meal plan shared successfully!</div>
        )}
        
        {error && <div className="error-message">{error}</div>}
        
        <div className="form-group">
          <label htmlFor="email">User Email</label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Enter email address"
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="role">Access Level</label>
          <select
            id="role"
            value={role}
            onChange={(e) => setRole(e.target.value)}
          >
            <option value="viewer">Viewer (can only view)</option>
            <option value="editor">Editor (can edit meals)</option>
          </select>
        </div>
        
        <div className="form-buttons">
          <button 
            className="cancel-button" 
            onClick={onClose}
            disabled={isSharing}
          >
            Cancel
          </button>
          <button 
            className="share-button" 
            onClick={handleShare}
            disabled={isSharing || !email}
          >
            {isSharing ? 'Sharing...' : 'Share'}
          </button>
        </div>
      </div>
    </Modal>
  );
};

export default ShareMealPlanModal; 