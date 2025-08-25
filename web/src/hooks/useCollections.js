// src/hooks/useCollections.js
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { collectionService } from '../services/collectionService';

// Custom hook for fetching collections
export const useCollections = () => {
  return useQuery({
    queryKey: ['collections'], // Unique key for the query
    queryFn: collectionService.list,
  });
};

// Custom hook for creating a collection
export const useCreateCollection = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: collectionService.create,
    onSuccess: () => {
      // Invalidate and refetch the collections list query
      queryClient.invalidateQueries({ queryKey: ['collections'] });
    },
  });
};

// Custom hook for updating a collection
export const useUpdateCollection = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }) => collectionService.update(id, data),
    onSuccess: () => {
      // Invalidate and refetch the collections list query
      queryClient.invalidateQueries({ queryKey: ['collections'] });
    },
  });
};

// Custom hook for deleting a collection
export const useDeleteCollection = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: collectionService.delete,
    onSuccess: () => {
      // Invalidate and refetch the collections list query
      queryClient.invalidateQueries({ queryKey: ['collections'] });
    },
  });
};

// Custom hook for fetching a single collection
export const useCollection = (id) => {
  return useQuery({
    queryKey: ['collection', id], // Unique key includes the ID
    queryFn: () => collectionService.get(id),
    enabled: !!id, // Only run the query if id is truthy
  });
};