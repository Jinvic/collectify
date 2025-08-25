// src/hooks/useItemAssociations.js
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { itemAssociationService } from '../services/itemAssociationService';

// Custom hook for adding a tag to an item
export const useAddTagToItem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ itemId, tagId }) => itemAssociationService.addTag(itemId, tagId),
    onSuccess: (_, variables) => {
      // Invalidate and refetch the specific item
      queryClient.invalidateQueries({ queryKey: ['item', variables.itemId] });
    },
  });
};

// Custom hook for removing a tag from an item
export const useRemoveTagFromItem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ itemId, tagId }) => itemAssociationService.removeTag(itemId, tagId),
    onSuccess: (_, variables) => {
      // Invalidate and refetch the specific item
      queryClient.invalidateQueries({ queryKey: ['item', variables.itemId] });
    },
  });
};

// Custom hook for adding an item to a collection
export const useAddItemToCollection = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ itemId, collectionId }) => itemAssociationService.addToCollection(itemId, collectionId),
    onSuccess: (_, variables) => {
      // Invalidate and refetch the specific item
      queryClient.invalidateQueries({ queryKey: ['item', variables.itemId] });
    },
  });
};

// Custom hook for removing an item from a collection
export const useRemoveItemFromCollection = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ itemId, collectionId }) => itemAssociationService.removeFromCollection(itemId, collectionId),
    onSuccess: (_, variables) => {
      // Invalidate and refetch the specific item
      queryClient.invalidateQueries({ queryKey: ['item', variables.itemId] });
    },
  });
};