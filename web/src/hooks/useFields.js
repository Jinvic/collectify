// src/hooks/useFields.js
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { fieldService } from '../services/fieldService';

// Custom hook for creating a field
export const useCreateField = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: fieldService.create,
    onSuccess: () => {
      // Invalidate queries that might be affected (e.g., category details which include fields)
      // Assuming category details are cached under 'category' key
      queryClient.invalidateQueries({ queryKey: ['category'] }); 
    },
  });
};

// Custom hook for deleting a field
export const useDeleteField = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: fieldService.delete,
    onSuccess: () => {
      // Invalidate queries that might be affected
      queryClient.invalidateQueries({ queryKey: ['category'] });
    },
  });
};