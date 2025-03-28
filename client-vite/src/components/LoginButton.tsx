import React from 'react';

interface LoginButtonProps {
  onLoginSuccess: (token: string) => void;
  onLoginFailure: (error: string) => void;
}

const LoginButton: React.FC<LoginButtonProps> = ({ onLoginSuccess, onLoginFailure }) => {
  const handleLogin = async () => {
    try {
      // For now, we'll just simulate a successful login
      const mockToken = 'mock-auth-token';
      localStorage.setItem('auth_token', mockToken);
      onLoginSuccess(mockToken);
    } catch (error) {
      onLoginFailure('Failed to login. Please try again.');
    }
  };

  return (
    <button className="login-button" onClick={handleLogin}>
      Login
    </button>
  );
};

export default LoginButton; 