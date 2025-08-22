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

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    // Get token from localStorage
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = token;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
apiClient.interceptors.response.use(
  (response) => {
    // If response has success code, return data directly
    if (response.data && response.data.code === 0) {
      return response.data;
    }
    
    // If response has fail code, throw error
    if (response.data && response.data.code === 1) {
      throw new Error(response.data.msg || 'Request failed');
    }
    
    return response;
  },
  (error) => {
    // Handle network errors
    if (!error.response) {
      throw new Error('Network error. Please check your connection.');
    }
    
    // Handle HTTP errors
    const { status, data } = error.response;
    
    // Handle unauthorized access
    if (status === 401 || (data && data.code === 1 && data.msg === '未登录')) {
      // Clear local auth data
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      // Redirect to login or show login prompt
      window.dispatchEvent(new Event('unauthorized'));
      
      throw new Error('请先登录');
    }
    
    // Handle other errors
    const message = data?.msg || data?.message || `HTTP Error: ${status}`;
    throw new Error(message);
  }
);

export default apiClient;