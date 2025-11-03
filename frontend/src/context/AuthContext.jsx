import { createContext, useContext, useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(localStorage.getItem("token"));
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const login = async (username, password) => {
    try {
      setError(null);
      const response = await axios.post("/api/auth/login", {
        username,
        password,
      });

      if (response.data.token) {
        setToken(response.data.token);
        localStorage.setItem("token", response.data.token);
        navigate("/");
      }
    } catch (err) {
      console.error("Login failed:", err.response.data);
      setError(err.response.data.error || "Login failed");
    }
  };

  const register = async (username, password) => {
    try {
      setError(null);
      await axios.post("/api/auth/register", { username, password });
      // After register, send them to the login page
      navigate("/login");
    } catch (err) {
      console.error("Register failed:", err.response.data);
      setError(err.response.data.error || "Registration failed");
    }
  };

  const logout = () => {
    setToken(null);
    localStorage.removeItem("token");
    navigate("/login");
  };

  const value = {
    token,
    error,
    login,
    register,
    logout,
    isAuthenticated: !!token,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  return useContext(AuthContext);
};
