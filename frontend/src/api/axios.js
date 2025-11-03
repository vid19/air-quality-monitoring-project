import axios from "axios";

// Create a new axios instance
const api = axios.create();

// This is the interceptor. It's a function that runs before any request is sent.
api.interceptors.request.use(
  (config) => {
    // Get the token from localStorage (where our AuthContext saved it)
    const token = localStorage.getItem("token");

    if (token) {
      // If the token exists, add it to the request headers
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default api;
