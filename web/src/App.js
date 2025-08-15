import React, { useState, useEffect } from 'react';

// Base URL for the API
// In development, this relies on the 'proxy' setting in package.json
// In production, it should be relative to the same origin
const API_BASE_URL = '/api'; // Proxy handles '/api' in dev, and it's relative in prod

function App() {
  const [categories, setCategories] = useState([]);
  const [newCategoryName, setNewCategoryName] = useState('');
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Fetch categories and items on component mount
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        const [categoriesRes, itemsRes] = await Promise.all([
          fetch(`${API_BASE_URL}/category/list`).then(res => res.json()),
          fetch(`${API_BASE_URL}/item/list?page=1&page_size=10`).then(res => res.json()), // Initial item list
        ]);
        setCategories(categoriesRes.data?.list || []);
        setItems(itemsRes.data?.list || []);
      } catch (err) {
        console.error("Failed to fetch data:", err);
        setError("Failed to load data. Please check your connection and backend.");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleCreateCategory = async () => {
    if (!newCategoryName.trim()) {
      alert("Category name cannot be empty.");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/category`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: newCategoryName }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // Re-fetch categories to update the list
      const res = await fetch(`${API_BASE_URL}/category/list`).then(res => res.json());
      setCategories(res.data?.list || []);
      setNewCategoryName(''); // Clear input
      alert("Category created successfully!");
    } catch (err) {
      console.error("Failed to create category:", err);
      setError("Failed to create category.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App" style={{ padding: '20px' }}>
      <h1>Collectify</h1>
      
      {error && <p style={{ color: 'red' }}>{error}</p>}

      <div>
        <h2>Categories</h2>
        {loading && <p>Loading categories...</p>}
        <ul>
          {categories.map(category => (
            <li key={category.id}>{category.name}</li>
          ))}
        </ul>
        <div>
          <input
            type="text"
            value={newCategoryName}
            onChange={(e) => setNewCategoryName(e.target.value)}
            placeholder="New Category Name"
            disabled={loading}
          />
          <button onClick={handleCreateCategory} disabled={loading || !newCategoryName.trim()}>
            Create Category
          </button>
        </div>
      </div>

      <div>
        <h2>Items</h2>
        {loading && <p>Loading items...</p>}
        <ul>
          {items.map(item => (
            <li key={item.id}>{item.title}</li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;