// src/services/tagService.js
import apiClient from './apiClient';

export const tagService = {
  list: async () => {
    try {
      const response = await apiClient.get('/tag/list');
      return response;
    } catch (error) {
      throw new Error(`Failed to fetch tags: ${error.message}`);
    }
  },

  create: async (tagData) => {
    try {
      const response = await apiClient.post('/tag', tagData);
      return response;
    } catch (error) {
      throw new Error(`Failed to create tag: ${error.message}`);
    }
  },

  rename: async (id, tagData) => {
    try {
      const response = await apiClient.patch(`/tag/${id}`, tagData);
      return response;
    } catch (error) {
      throw new Error(`Failed to rename tag: ${error.message}`);
    }
  },

  delete: async (id) => {
    try {
      const response = await apiClient.delete(`/tag/${id}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to delete tag: ${error.message}`);
    }
  },

  restore: async (id) => {
    try {
      const response = await apiClient.post(`/tag/${id}/restore`);
      return response;
    } catch (error) {
      throw new Error(`Failed to restore tag: ${error.message}`);
    }
  },

  // Get a single tag
  get: async (id) => {
    try {
      const response = await apiClient.get(`/tag/${id}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to fetch tag: ${error.message}`);
    }
  },
};