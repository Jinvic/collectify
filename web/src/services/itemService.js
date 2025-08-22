// src/services/itemService.js
import apiClient from './apiClient';

export const itemService = {
  list: async (params = {}) => {
    try {
      // Default pagination if not provided
      const defaultParams = { page: 1, page_size: 10, ...params };
      const response = await apiClient.get('/item/list', { params: defaultParams });
      return response;
    } catch (error) {
      throw new Error(`Failed to fetch items: ${error.message}`);
    }
  },

  search: async (searchData) => {
    try {
      const response = await apiClient.post('/item/search', searchData);
      return response;
    } catch (error) {
      throw new Error(`Failed to search items: ${error.message}`);
    }
  },

  get: async (id) => {
    try {
      const response = await apiClient.get(`/item/${id}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to fetch item ${id}: ${error.message}`);
    }
  },

  create: async (itemData) => {
    try {
      const response = await apiClient.post('/item', itemData);
      return response;
    } catch (error) {
      throw new Error(`Failed to create item: ${error.message}`);
    }
  },

  update: async (id, itemData) => {
    try {
      const response = await apiClient.put(`/item/${id}`, itemData);
      return response;
    } catch (error) {
      throw new Error(`Failed to update item ${id}: ${error.message}`);
    }
  },

  delete: async (id) => {
    try {
      const response = await apiClient.delete(`/item/${id}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to delete item ${id}: ${error.message}`);
    }
  },

  restore: async (id) => {
    try {
      const response = await apiClient.post(`/item/${id}/restore`);
      return response;
    } catch (error) {
      throw new Error(`Failed to restore item ${id}: ${error.message}`);
    }
  },
  // Add more methods for other item operations as needed
};