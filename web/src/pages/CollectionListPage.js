// src/pages/CollectionListPage.js
import React, { useState } from 'react';
import { useCollections, useCreateCollection } from '../hooks/useCollections';
import { Link } from 'react-router-dom';
import {
  Container, Typography, Box, List, ListItem, ListItemText, ListItemSecondaryAction,
  IconButton, Button, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, CircularProgress, Alert, Snackbar
} from '@mui/material';
import { Edit, Add } from '@mui/icons-material';

const CollectionListPage = () => {
  const { data, isLoading, error } = useCollections();
  const collections = data?.data?.list || [];
  const { mutate: createCollection, isLoading: isCreating, error: createError } = useCreateCollection();

  const [open, setOpen] = useState(false);
  const [newCollectionName, setNewCollectionName] = useState('');
  const [newCollectionDescription, setNewCollectionDescription] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setNewCollectionName('');
    setNewCollectionDescription('');
  };

  const handleCreate = () => {
    if (newCollectionName.trim()) {
      createCollection({ 
        name: newCollectionName,
        description: newCollectionDescription
      }, {
        onSuccess: () => {
          handleClose();
          setSnackbar({ open: true, message: 'Collection created successfully!', severity: 'success' });
        },
        onError: (error) => {
          console.error("Failed to create collection:", error);
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
          Collections
        </Typography>

        {error && <Alert severity="error">{error.message}</Alert>}
        {createError && <Alert severity="error">{createError.message}</Alert>}
        {isLoading && <CircularProgress />}

        <List>
          {collections.map((collection) => (
            <ListItem key={collection.id} divider>
              <ListItemText
                primary={collection.name}
                secondary={collection.description || 'No description'}
              />
              <ListItemSecondaryAction>
                <IconButton edge="end" aria-label="edit" component={Link} to={`/collections/${collection.id}`}>
                  <Edit />
                </IconButton>
                {/* 
                // Placeholder for delete action
                <IconButton edge="end" aria-label="delete" onClick={() => handleDelete(collection.id)}>
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
            Create New Collection
          </Button>
        </Box>
      </Box>

      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Create New Collection</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Collection Name"
            type="text"
            fullWidth
            variant="outlined"
            value={newCollectionName}
            onChange={(e) => setNewCollectionName(e.target.value)}
            disabled={isCreating}
          />
          <TextField
            margin="dense"
            label="Description (optional)"
            type="text"
            fullWidth
            variant="outlined"
            multiline
            rows={3}
            value={newCollectionDescription}
            onChange={(e) => setNewCollectionDescription(e.target.value)}
            disabled={isCreating}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} disabled={isCreating}>
            Cancel
          </Button>
          <Button onClick={handleCreate} disabled={isCreating || !newCollectionName.trim()}>
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

export default CollectionListPage;