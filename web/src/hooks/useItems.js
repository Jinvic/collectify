// src/hooks/useItems.js
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { itemService } from '../services/itemService';

// Custom hook for fetching basic item list
export const useBasicItems = (params = {}) => {
  return useQuery({
    queryKey: ['basicItems', params], // Key includes params for caching different queries
    queryFn: () => itemService.list(params),
  });
};

// Custom hook for searching items
export const useSearchItems = (searchParams = {}) => {
  return useQuery({
    queryKey: ['searchItems', searchParams],
    queryFn: () => itemService.search(searchParams),
    // enabled: Object.keys(searchParams).length > 0, // Optional: only run if params are provided
  });
};

// Custom hook for fetching a single item
export const useItem = (id) => {
  return useQuery({
    queryKey: ['item', id],
    queryFn: () => itemService.get(id),
    enabled: !!id,
  });
};

// Custom hook for creating an item
export const useCreateItem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: itemService.create,
    onSuccess: () => {
      // Invalidate queries that might be affected
      queryClient.invalidateQueries({ queryKey: ['basicItems'] });
      queryClient.invalidateQueries({ queryKey: ['searchItems'] });
    },
  });
};

// Custom hook for updating an item
export const useUpdateItem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }) => itemService.update(id, data),
    onSuccess: (data, variables) => {
      // Invalidate and refetch the specific item
      queryClient.invalidateQueries({ queryKey: ['item', variables.id] });
      // Invalidate list queries as the item might have changed position or data
      queryClient.invalidateQueries({ queryKey: ['basicItems'] });
      queryClient.invalidateQueries({ queryKey: ['searchItems'] });
    },
  });
};

// Custom hook for deleting an item
export const useDeleteItem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: itemService.delete,
    onSuccess: () => {
      // Invalidate list queries
      queryClient.invalidateQueries({ queryKey: ['basicItems'] });
      queryClient.invalidateQueries({ queryKey: ['searchItems'] });
      // Note: The specific item query will become stale, which is usually fine
    },
  });
};