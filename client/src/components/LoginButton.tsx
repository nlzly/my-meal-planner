import React from 'react';

const LoginButton: React.FC = () => {
  const handleLogin = () => {
    // Redirect to the server's OAuth login endpoint
    window.location.href = '/auth/google/login';
  };

  return (
    <button 
      className="login-button"
      onClick={handleLogin}
    >
      Sign in with Google
    </button>
  );
};

export default LoginButton; 