// src/hooks/useCategories.js
import { useState, useEffect } from 'react';
import { categoryService } from '../services/categoryService';

export const useCategories = () => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchCategories = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await categoryService.list();
      setCategories(data.data?.list || []);
    } catch (err) {
      console.error("Hook: Failed to fetch categories:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const createCategory = async (categoryData) => {
    setLoading(true);
    setError(null);
    try {
      const data = await categoryService.create(categoryData);
      // Refetch categories after creation
      await fetchCategories();
      return data;
    } catch (err) {
      console.error("Hook: Failed to create category:", err);
      setError(err.message);
      throw err; // Re-throw so caller can handle
    } finally {
      setLoading(false);
    }
  };

  // Fetch categories on initial load
  useEffect(() => {
    fetchCategories();
  }, []);

  return { categories, loading, error, fetchCategories, createCategory };
};