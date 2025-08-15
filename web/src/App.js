// src/App.js
import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import HomePage from './pages/HomePage';
import CategoryListPage from './pages/CategoryListPage';
import CategoryDetailPage from './pages/CategoryDetailPage';
import ItemDetailPage from './pages/ItemDetailPage';
// Import other pages as they are created
// import TagListPage from './pages/TagListPage';
// import CollectionListPage from './pages/CollectionListPage';
// import SearchPage from './pages/SearchPage';

function App() {
  return (
    <div>
      <Navbar />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/categories" element={<CategoryListPage />} />
        <Route path="/categories/:id" element={<CategoryDetailPage />} />
        <Route path="/items/:id" element={<ItemDetailPage />} />
        {/* Add routes for other pages */}
        {/* <Route path="/tags" element={<TagListPage />} /> */}
        {/* <Route path="/collections" element={<CollectionListPage />} /> */}
        {/* <Route path="/search" element={<SearchPage />} /> */}
      </Routes>
    </div>
  );
}

export default App;