// src/pages/CategoryListPage.js
import React, { useState } from 'react';
import { useCategories, useCreateCategory } from '../hooks/useCategories';
import { Link } from 'react-router-dom';
import {
  Container, Typography, Box, List, ListItem, ListItemText, ListItemSecondaryAction,
  IconButton, Button, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, CircularProgress, Alert, Snackbar
} from '@mui/material';
import { Edit, Add } from '@mui/icons-material';

const CategoryListPage = () => {
  const { data, isLoading, error } = useCategories();
  const categories = data?.data?.list || [];
  const { mutate: createCategory, isLoading: isCreating, error: createError } = useCreateCategory();

  const [open, setOpen] = useState(false);
  const [newCategoryName, setNewCategoryName] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setNewCategoryName('');
  };

  const handleCreate = () => {
    if (newCategoryName.trim()) {
      createCategory({ name: newCategoryName }, {
        onSuccess: () => {
          handleClose();
          setSnackbar({ open: true, message: 'Category created successfully!', severity: 'success' });
        },
        onError: (error) => {
          console.error("Failed to create category:", error);
          setSnackbar({ open: true, message: `Error: ${error.message}`, severity: 'error' });
        }
      });
    }
  };

  const handleCloseSnackbar = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackbar({ ...snackbar, open: false });
  };

  return (
    <Container maxWidth="md">
      <Box my={4}>
        <Typography variant="h4" component="h1" gutterBottom>
          Categories
        </Typography>

        {error && <Alert severity="error">{error.message}</Alert>}
        {createError && <Alert severity="error">{createError.message}</Alert>}
        {isLoading && <CircularProgress />}

        <List>
          {categories.map((category) => (
            <ListItem key={category.id} divider>
              <ListItemText
                primary={category.name}
                // secondary={`ID: ${category.id}`} // Optional: show ID
              />
              <ListItemSecondaryAction>
                <IconButton edge="end" aria-label="edit" component={Link} to={`/categories/${category.id}`}>
                  <Edit />
                </IconButton>
                {/* 
                // Placeholder for delete action
                <IconButton edge="end" aria-label="delete" onClick={() => handleDelete(category.id)}>
                  <Delete />
                </IconButton> 
                */}
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>

        <Box mt={2}>
          <Button
            variant="contained"
            color="primary"
            startIcon={<Add />}
            onClick={handleClickOpen}
          >
            Create New Category
          </Button>
        </Box>
      </Box>

      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Create New Category</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Category Name"
            type="text"
            fullWidth
            variant="outlined"
            value={newCategoryName}
            onChange={(e) => setNewCategoryName(e.target.value)}
            disabled={isCreating}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} disabled={isCreating}>
            Cancel
          </Button>
          <Button onClick={handleCreate} disabled={isCreating || !newCategoryName.trim()}>
            Create
          </Button>
        </DialogActions>
      </Dialog>
      
      {/* Snackbar for notifications */}
      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        message={snackbar.message}
      />
    </Container>
  );
};

export default CategoryListPage;