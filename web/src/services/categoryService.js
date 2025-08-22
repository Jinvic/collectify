// src/services/categoryService.js
import apiClient from './apiClient';

export const categoryService = {
  list: async () => {
    try {
      const response = await apiClient.get('/category/list');
      return response;
    } catch (error) {
      throw new Error(`Failed to fetch categories: ${error.message}`);
    }
  },

  create: async (categoryData) => {
    try {
      const response = await apiClient.post('/category', categoryData);
      return response;
    } catch (error) {
      throw new Error(`Failed to create category: ${error.message}`);
    }
  },

  rename: async (id, categoryData) => {
    try {
      const response = await apiClient.patch(`/category/${id}`, categoryData);
      return response;
    } catch (error) {
      throw new Error(`Failed to rename category: ${error.message}`);
    }
  },

  delete: async (id) => {
    try {
      const response = await apiClient.delete(`/category/${id}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to delete category: ${error.message}`);
    }
  },

  // Get a single category with its fields
  get: async (id) => {
    try {
      const response = await apiClient.get(`/category/${id}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to fetch category: ${error.message}`);
    }
  },
};