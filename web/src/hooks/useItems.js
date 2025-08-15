// src/hooks/useItems.js
import { useState, useEffect } from 'react';
import { itemService } from '../services/itemService';

export const useItems = (initialParams = {}) => {
  const [items, setItems] = useState([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [params, setParams] = useState({ page: 1, page_size: 10, ...initialParams });

  const fetchItems = async (fetchParams = params) => {
    setLoading(true);
    setError(null);
    try {
      const data = await itemService.list(fetchParams);
      setItems(data.data?.list || []);
      setTotal(data.data?.total || 0);
    } catch (err) {
      console.error("Hook: Failed to fetch items:", err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Fetch items when params change
  useEffect(() => {
    fetchItems(params);
  }, [params]);

  const searchItems = async (searchData) => {
    setLoading(true);
    setError(null);
    try {
      const data = await itemService.search(searchData);
      setItems(data.data?.list || []);
      setTotal(data.data?.total || 0);
      return data;
    } catch (err) {
      console.error("Hook: Failed to search items:", err);
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { items, total, loading, error, fetchItems, searchItems, setParams, params };
};