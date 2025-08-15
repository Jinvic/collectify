// src/api.js
import axios from 'axios';

// Determine the base URL based on the environment
// In development, the 'proxy' in package.json will handle '/api' requests
// In production, the Go backend will serve the API from the same origin
const API_BASE_URL = process.env.NODE_ENV === 'production' ? '/api' : '/api';

// Create an axios instance with the base URL
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000, // 10 seconds timeout
});

export default apiClient;