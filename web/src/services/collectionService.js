// src/services/collectionService.js
import apiClient from './apiClient';

export const collectionService = {
  list: async () => {
    try {
      const response = await apiClient.get('/collection/list');
      return response;
    } catch (error) {
      throw new Error(`Failed to fetch collections: ${error.message}`);
    }
  },

  create: async (collectionData) => {
    try {
      const response = await apiClient.post('/collection', collectionData);
      return response;
    } catch (error) {
      throw new Error(`Failed to create collection: ${error.message}`);
    }
  },

  update: async (id, collectionData) => {
    try {
      const response = await apiClient.patch(`/collection/${id}`, collectionData);
      return response;
    } catch (error) {
      throw new Error(`Failed to update collection: ${error.message}`);
    }
  },

  delete: async (id) => {
    try {
      const response = await apiClient.delete(`/collection/${id}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to delete collection: ${error.message}`);
    }
  },

  restore: async (id) => {
    try {
      const response = await apiClient.post(`/collection/${id}/restore`);
      return response;
    } catch (error) {
      throw new Error(`Failed to restore collection: ${error.message}`);
    }
  },

  // Get a single collection
  get: async (id) => {
    try {
      const response = await apiClient.get(`/collection/${id}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to fetch collection: ${error.message}`);
    }
  },
};