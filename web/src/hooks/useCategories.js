// src/hooks/useCategories.js
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { categoryService } from '../services/categoryService';

// Custom hook for fetching categories
export const useCategories = () => {
  return useQuery({
    queryKey: ['categories'], // Unique key for the query
    queryFn: categoryService.list,
    // Initial data can be provided if needed, e.g., from a placeholder or cache
    // initialData: { data: { list: [] } },
  });
};

// Custom hook for creating a category
export const useCreateCategory = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: categoryService.create,
    onSuccess: () => {
      // Invalidate and refetch the categories list query
      queryClient.invalidateQueries({ queryKey: ['categories'] });
    },
  });
};

// Custom hook for renaming a category
export const useRenameCategory = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }) => categoryService.rename(id, data),
    onSuccess: () => {
      // Invalidate and refetch the categories list query
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      // Also invalidate any single category queries if they exist
      // queryClient.invalidateQueries({ queryKey: ['category'] }); 
    },
  });
};

// Custom hook for deleting a category
export const useDeleteCategory = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: categoryService.delete,
    onSuccess: () => {
      // Invalidate and refetch the categories list query
      queryClient.invalidateQueries({ queryKey: ['categories'] });
    },
  });
};

// Custom hook for fetching a single category (with fields)
export const useCategory = (id) => {
  return useQuery({
    queryKey: ['category', id], // Unique key includes the ID
    queryFn: () => categoryService.get(id),
    enabled: !!id, // Only run the query if id is truthy
  });
};