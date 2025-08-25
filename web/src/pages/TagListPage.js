// src/pages/TagListPage.js
import React, { useState } from 'react';
import { useTags, useCreateTag } from '../hooks/useTags';
import { Link } from 'react-router-dom';
import {
  Container, Typography, Box, List, ListItem, ListItemText, ListItemSecondaryAction,
  IconButton, Button, Dialog, DialogTitle, DialogContent, DialogActions,
  TextField, CircularProgress, Alert, Snackbar
} from '@mui/material';
import { Edit, Add } from '@mui/icons-material';

const TagListPage = () => {
  const { data, isLoading, error } = useTags();
  const tags = data?.data?.list || [];
  const { mutate: createTag, isLoading: isCreating, error: createError } = useCreateTag();

  const [open, setOpen] = useState(false);
  const [newTagName, setNewTagName] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setNewTagName('');
  };

  const handleCreate = () => {
    if (newTagName.trim()) {
      createTag({ name: newTagName }, {
        onSuccess: () => {
          handleClose();
          setSnackbar({ open: true, message: 'Tag created successfully!', severity: 'success' });
        },
        onError: (error) => {
          console.error("Failed to create tag:", error);
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
          Tags
        </Typography>

        {error && <Alert severity="error">{error.message}</Alert>}
        {createError && <Alert severity="error">{createError.message}</Alert>}
        {isLoading && <CircularProgress />}

        <List>
          {tags.map((tag) => (
            <ListItem key={tag.id} divider>
              <ListItemText
                primary={tag.name}
                // secondary={`ID: ${tag.id}`} // Optional: show ID
              />
              <ListItemSecondaryAction>
                <IconButton edge="end" aria-label="edit" component={Link} to={`/tags/${tag.id}`}>
                  <Edit />
                </IconButton>
                {/* 
                // Placeholder for delete action
                <IconButton edge="end" aria-label="delete" onClick={() => handleDelete(tag.id)}>
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
            Create New Tag
          </Button>
        </Box>
      </Box>

      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Create New Tag</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Tag Name"
            type="text"
            fullWidth
            variant="outlined"
            value={newTagName}
            onChange={(e) => setNewTagName(e.target.value)}
            disabled={isCreating}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} disabled={isCreating}>
            Cancel
          </Button>
          <Button onClick={handleCreate} disabled={isCreating || !newTagName.trim()}>
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

export default TagListPage;