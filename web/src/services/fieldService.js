// src/services/fieldService.js
import apiClient from './apiClient';

export const fieldService = {
  create: async (fieldData) => {
    try {
      const response = await apiClient.post('/field', fieldData);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to create field: ${error.response?.data?.msg || error.message}`);
    }
  },

  delete: async (id) => {
    try {
      const response = await apiClient.delete(`/field/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to delete field: ${error.response?.data?.msg || error.message}`);
    }
  },

  restore: async (id) => {
    try {
      const response = await apiClient.post(`/field/${id}/restore`);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to restore field: ${error.response?.data?.msg || error.message}`);
    }
  },
  // Add more methods for other field operations as needed (e.g., update if backend supports it)
};