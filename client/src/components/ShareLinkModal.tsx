import React, { useState } from 'react';
import Modal from './Modal';
import api from '../services/axios';

interface ShareLinkModalProps {
  isOpen: boolean;
  onClose: () => void;
  mealPlanId: string;
  mealPlanName: string;
}

const ShareLinkModal: React.FC<ShareLinkModalProps> = ({
  isOpen,
  onClose,
  mealPlanId,
  mealPlanName,
}) => {
  const [role, setRole] = useState('viewer');
  const [expiresIn, setExpiresIn] = useState(7*24); // 7 days in hours
  const [shareLink, setShareLink] = useState<string | null>(null);
  const [error, setError] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [isCopied, setIsCopied] = useState(false);

  const handleGenerateLink = async () => {
    setIsGenerating(true);
    setError(''); 

    try {
      const response = await api.post('/api/meal-plans/generate-link', {
        mealPlanId,
        role,
        expiresIn,
      });
      
      setShareLink(response.data.code);
    } catch (err: any) {
      setError(err.response?.data || 'Failed to generate share link');
    } finally {
      setIsGenerating(false);
    }
  };

  const handleCopyLink = () => {
    if (shareLink) {
      navigator.clipboard.writeText(shareLink);
      setIsCopied(true);
      setTimeout(() => setIsCopied(false), 2000);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className="share-link-modal">
        <h2>Share Meal Plan</h2>
        <p>Generate a code to share "{mealPlanName}" with others</p>
        
        {error && <div className="error-message">{error}</div>}
        
        {!shareLink ? (
          <>
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
            
            <div className="form-group">
              <label htmlFor="expiresIn">Link Expires In</label>
              <select 
                id="expiresIn"
                value={expiresIn}
                onChange={(e) => setExpiresIn(Number(e.target.value))}
              >
                <option value={24}>1 Day</option>
                <option value={72}>3 Days</option>
                <option value={168}>7 Days</option>
                <option value={720}>30 Days</option>
              </select>
            </div>
            
            <div className="form-buttons">
              <button 
                className="cancel-button" 
                onClick={onClose}
                disabled={isGenerating}
              >
                Cancel
              </button>
              <button 
                className="generate-button" 
                onClick={handleGenerateLink}
                disabled={isGenerating}
              >
                {isGenerating ? 'Generating...' : 'Generate Link'}
              </button>
            </div>
          </>
        ) : (
          <div className="share-link-container">
            <p>Share this code with others to give them access to your meal plan:</p>
            <div className="share-link-box">
              <input 
                type="text" 
                value={shareLink} 
                readOnly 
                className="share-link-input"
                onClick={(e) => (e.target as HTMLInputElement).select()}
              />
              <button 
                className="copy-link-button" 
                onClick={handleCopyLink}
              >
                {isCopied ? 'Copied!' : 'Copy'}
              </button>
            </div>
            <p className="share-link-note">
              Anyone with this code can join your meal plan.
            </p>
            <div className="form-buttons">
              <button 
                className="close-button" 
                onClick={onClose}
              >
                Close
              </button>
              <button 
                className="generate-new-button" 
                onClick={() => setShareLink(null)}
              >
                Generate New Link
              </button>
            </div>
            <div>

            </div>
          </div>
        )}
      </div>
    </Modal>
  );
};

export default ShareLinkModal; 