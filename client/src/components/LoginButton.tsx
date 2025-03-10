import React from 'react';
import { GoogleLogin } from '@react-oauth/google';
import axios from 'axios';

interface LoginButtonProps {
  onLoginSuccess: (token: string) => void;
  onLoginFailure: (error: string) => void;
}

const LoginButton: React.FC<LoginButtonProps> = ({ onLoginSuccess, onLoginFailure }) => {
  const handleSuccess = async (credentialResponse: any) => {
    try {
      // Exchange Google token for our JWT
      const response = await axios.post('/auth/google/callback', {
        credential: credentialResponse.credential
      });
      
      // Save token to localStorage
      localStorage.setItem('auth_token', response.data.token);
      
      // Set Authorization header for future requests
      axios.defaults.headers.common['Authorization'] = `Bearer ${response.data.token}`;
      
      // Notify parent component
      onLoginSuccess(response.data.token);
    } catch (error: any) {
      console.error('Login error:', error);
      onLoginFailure(error.response?.data?.message || 'Login failed');
    }
  };

  return (
    <div className="login-button-container">
      <GoogleLogin
        onSuccess={handleSuccess}
        onError={() => onLoginFailure('Google login failed')}
        useOneTap
      />
    </div>
  );
};

export default LoginButton;
