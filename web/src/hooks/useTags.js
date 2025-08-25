// src/hooks/useTags.js
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { tagService } from '../services/tagService';

// Custom hook for fetching tags
export const useTags = () => {
  return useQuery({
    queryKey: ['tags'], // Unique key for the query
    queryFn: tagService.list,
  });
};

// Custom hook for creating a tag
export const useCreateTag = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: tagService.create,
    onSuccess: () => {
      // Invalidate and refetch the tags list query
      queryClient.invalidateQueries({ queryKey: ['tags'] });
    },
  });
};

// Custom hook for renaming a tag
export const useRenameTag = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, ...data }) => tagService.rename(id, data),
    onSuccess: () => {
      // Invalidate and refetch the tags list query
      queryClient.invalidateQueries({ queryKey: ['tags'] });
    },
  });
};

// Custom hook for deleting a tag
export const useDeleteTag = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: tagService.delete,
    onSuccess: () => {
      // Invalidate and refetch the tags list query
      queryClient.invalidateQueries({ queryKey: ['tags'] });
    },
  });
};

// Custom hook for fetching a single tag
export const useTag = (id) => {
  return useQuery({
    queryKey: ['tag', id], // Unique key includes the ID
    queryFn: () => tagService.get(id),
    enabled: !!id, // Only run the query if id is truthy
  });
};