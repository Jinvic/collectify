// src/services/itemAssociationService.js
import apiClient from './apiClient';

export const itemAssociationService = {
  // Add tag to item
  addTag: async (itemId, tagId) => {
    try {
      const response = await apiClient.post(`/item/${itemId}/tag/${tagId}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to add tag to item: ${error.message}`);
    }
  },

  // Remove tag from item
  removeTag: async (itemId, tagId) => {
    try {
      const response = await apiClient.delete(`/item/${itemId}/tag/${tagId}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to remove tag from item: ${error.message}`);
    }
  },

  // Add item to collection
  addToCollection: async (itemId, collectionId) => {
    try {
      const response = await apiClient.post(`/item/${itemId}/collection/${collectionId}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to add item to collection: ${error.message}`);
    }
  },

  // Remove item from collection
  removeFromCollection: async (itemId, collectionId) => {
    try {
      const response = await apiClient.delete(`/item/${itemId}/collection/${collectionId}`);
      return response;
    } catch (error) {
      throw new Error(`Failed to remove item from collection: ${error.message}`);
    }
  }
};