import React, { useState } from 'react';
import {
  Dialog, DialogTitle, DialogContent, DialogActions,
  Button, TextField, FormControl, InputLabel, Select, MenuItem,
  Typography, Box, CircularProgress, Alert, FormControlLabel, Checkbox
} from '@mui/material';
import { useCreateItem } from '../hooks/useItems';
import { useCategories } from '../hooks/useCategories';

const AddItemDialog = ({ open, onClose, onItemAdded }) => {
  const { data: categoriesData, isLoading: categoriesLoading, error: categoriesError } = useCategories();
  const categories = categoriesData?.data?.list || [];
  
  const { mutate: createItem, isLoading: isCreating } = useCreateItem();
  
  const [selectedCategory, setSelectedCategory] = useState('');
  const [itemData, setItemData] = useState({
    name: '',
    status: 1,
    rating: '',
    description: '',
    notes: '',
    cover_url: '',
    source_url: '',
    priority: 0
  });

  const handleCategoryChange = (categoryId) => {
    setSelectedCategory(categoryId);
  };

  const handleInputChange = (field, value) => {
    setItemData(prev => ({ ...prev, [field]: value }));
  };

  const handleCreate = () => {
    if (selectedCategory && itemData.name.trim()) {
      const dataToSubmit = {
        category_id: selectedCategory,
        item: {
          ...itemData,
          rating: itemData.rating === '' ? null : parseFloat(itemData.rating),
          priority: parseInt(itemData.priority, 10)
        }
      };

      createItem(dataToSubmit, {
        onSuccess: () => {
          // Reset form
          setSelectedCategory('');
          setItemData({
            name: '',
            status: 1,
            rating: '',
            description: '',
            notes: '',
            cover_url: '',
            source_url: '',
            priority: 0
          });
          onItemAdded();
          onClose();
        },
        onError: (error) => {
          console.error("Failed to create item:", error);
          // TODO: Show error to user
        }
      });
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Add New Item</DialogTitle>
      <DialogContent>
        {categoriesError && <Alert severity="error">{categoriesError.message}</Alert>}
        
        <Box mt={1}>
          <FormControl fullWidth margin="normal" required>
            <InputLabel>Select Category</InputLabel>
            <Select
              value={selectedCategory}
              onChange={(e) => handleCategoryChange(e.target.value)}
              label="Select Category"
              disabled={categoriesLoading || isCreating}
            >
              {categoriesLoading ? (
                <MenuItem value="">
                  <CircularProgress size={24} />
                </MenuItem>
              ) : (
                categories.map((category) => (
                  <MenuItem key={category.id} value={category.id}>
                    {category.name}
                  </MenuItem>
                ))
              )}
            </Select>
          </FormControl>
          
          <TextField
            label="Name"
            fullWidth
            margin="normal"
            required
            value={itemData.name}
            onChange={(e) => handleInputChange('name', e.target.value)}
            disabled={isCreating}
          />
          
          <FormControl fullWidth margin="normal">
            <InputLabel>Status</InputLabel>
            <Select
              value={itemData.status}
              onChange={(e) => handleInputChange('status', e.target.value)}
              label="Status"
              disabled={isCreating}
            >
              <MenuItem value={1}>To Do</MenuItem>
              <MenuItem value={2}>In Progress</MenuItem>
              <MenuItem value={3}>Paused</MenuItem>
              <MenuItem value={4}>Abandoned</MenuItem>
              <MenuItem value={5}>Completed</MenuItem>
            </Select>
          </FormControl>
          
          <TextField
            label="Rating"
            type="number"
            fullWidth
            margin="normal"
            value={itemData.rating}
            onChange={(e) => handleInputChange('rating', e.target.value)}
            inputProps={{ min: 0, max: 10, step: 0.1 }}
            disabled={isCreating}
          />
          
          <TextField
            label="Description"
            fullWidth
            margin="normal"
            multiline
            rows={3}
            value={itemData.description}
            onChange={(e) => handleInputChange('description', e.target.value)}
            disabled={isCreating}
          />
          
          <TextField
            label="Notes"
            fullWidth
            margin="normal"
            multiline
            rows={3}
            value={itemData.notes}
            onChange={(e) => handleInputChange('notes', e.target.value)}
            disabled={isCreating}
          />
          
          <TextField
            label="Cover URL"
            fullWidth
            margin="normal"
            value={itemData.cover_url}
            onChange={(e) => handleInputChange('cover_url', e.target.value)}
            disabled={isCreating}
          />
          
          <TextField
            label="Source URL"
            fullWidth
            margin="normal"
            value={itemData.source_url}
            onChange={(e) => handleInputChange('source_url', e.target.value)}
            disabled={isCreating}
          />
          
          <TextField
            label="Priority"
            type="number"
            fullWidth
            margin="normal"
            value={itemData.priority}
            onChange={(e) => handleInputChange('priority', e.target.value)}
            inputProps={{ min: 0 }}
            disabled={isCreating}
          />
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={isCreating}>
          Cancel
        </Button>
        <Button 
          onClick={handleCreate} 
          disabled={isCreating || !selectedCategory || !itemData.name.trim()}
          variant="contained"
          color="primary"
        >
          {isCreating ? <CircularProgress size={24} /> : 'Create'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddItemDialog;