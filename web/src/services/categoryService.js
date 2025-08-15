// src/services/categoryService.js
import apiClient from '../api';

export const categoryService = {
  list: async () => {
    try {
      const response = await apiClient.get('/category/list');
      return response.data;
    } catch (error) {
      throw new Error(`Failed to fetch categories: ${error.message}`);
    }
  },

  create: async (categoryData) => {
    try {
      const response = await apiClient.post('/category', categoryData);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to create category: ${error.message}`);
    }
  },

  // Add more methods for update, delete, etc. as needed
  // get: async (id) => { ... }
  // update: async (id, categoryData) => { ... }
  // delete: async (id) => { ... }
};