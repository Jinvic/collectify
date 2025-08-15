// src/services/itemService.js
import apiClient from './apiClient';

export const itemService = {
  list: async (params = {}) => {
    try {
      // Default pagination if not provided
      const defaultParams = { page: 1, page_size: 10, ...params };
      const response = await apiClient.get('/item/list', { params: defaultParams });
      return response.data;
    } catch (error) {
      throw new Error(`Failed to fetch items: ${error.response?.data?.msg || error.message}`);
    }
  },

  search: async (searchData) => {
    try {
      const response = await apiClient.post('/item/search', searchData);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to search items: ${error.response?.data?.msg || error.message}`);
    }
  },

  get: async (id) => {
    try {
      const response = await apiClient.get(`/item/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to fetch item ${id}: ${error.response?.data?.msg || error.message}`);
    }
  },

  create: async (itemData) => {
    try {
      const response = await apiClient.post('/item', itemData);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to create item: ${error.response?.data?.msg || error.message}`);
    }
  },

  update: async (id, itemData) => {
    try {
      const response = await apiClient.put(`/item/${id}`, itemData);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to update item ${id}: ${error.response?.data?.msg || error.message}`);
    }
  },

  delete: async (id) => {
    try {
      const response = await apiClient.delete(`/item/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to delete item ${id}: ${error.response?.data?.msg || error.message}`);
    }
  },

  restore: async (id) => {
    try {
      const response = await apiClient.post(`/item/${id}/restore`);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to restore item ${id}: ${error.response?.data?.msg || error.message}`);
    }
  },
  // Add more methods for other item operations as needed
};