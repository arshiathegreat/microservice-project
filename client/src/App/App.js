// App.js
import React, { useState, useEffect } from 'react';
import Login from '../components/Login/Login';
import Dashboard from '../components/Dashboard/Dashboard';

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));

  const handleLogin = (token, name) => {
    console.log(name);
    setToken(token);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setToken(null);
    
  };

  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token); 
    } 
  }, [token]);

 

  return (
    <div>
      {!token ? (
        <Login setToken={handleLogin} setName={handleLogin} />
      ) : (
        <Dashboard token={token} handleLogout={handleLogout} />
      )}
    </div>
  );
}

export default App;