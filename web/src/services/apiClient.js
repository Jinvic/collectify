// src/services/apiClient.js
import axios from 'axios';

// Determine the base URL based on the environment
// In development, the 'proxy' in package.json will handle '/api' requests
// In production, the Go backend will serve the API from the same origin
const API_BASE_URL = process.env.NODE_ENV === 'production' ? '/api' : '/api';

// Create an axios instance with the base URL
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000, // 10 seconds timeout
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor (optional)
apiClient.interceptors.request.use(
  (config) => {
    // You can add auth tokens here if needed in the future
    // const token = localStorage.getItem('token');
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor (optional, for global error handling)
apiClient.interceptors.response.use(
  (response) => {
    // Any status code that lie within the range of 2xx cause this function to trigger
    return response;
  },
  (error) => {
    // Any status codes that falls outside the range of 2xx cause this function to trigger
    console.error('API Error:', error.response?.data?.msg || error.message);
    // You could dispatch a global error action here if using a state management library
    return Promise.reject(error);
  }
);

export default apiClient;