// src/pages/ItemDetailPage.js
import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useItem, useUpdateItem, useDeleteItem } from '../hooks/useItems';
import { formatFieldValue, getFieldTypeName } from '../utils/itemUtils';
import { formatDate } from '../utils/formatDate';
import {
  Container, Typography, Box, CircularProgress, Alert, Button, Dialog, DialogTitle,
  DialogContent, DialogActions, TextField, FormControl, InputLabel, Select, MenuItem,
  FormControlLabel, Checkbox, Paper, Divider, IconButton, Snackbar, Alert as MuiAlert
} from '@mui/material';
import { Edit, Delete, Save, Cancel } from '@mui/icons-material';

const ItemDetailPage = () => {
  const { id } = useParams();
  const itemId = parseInt(id, 10);

  const { data, isLoading, error, refetch } = useItem(itemId);
  const item = data?.data;
  
  const { mutate: updateItem, isLoading: isUpdating, error: updateError } = useUpdateItem();
  const { mutate: deleteItem, error: deleteError } = useDeleteItem();

  const [editMode, setEditMode] = useState(false);
  const [editedItem, setEditedItem] = useState({});
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  // Sync edited item when data loads or edit mode changes
  React.useEffect(() => {
    if (item) {
      // Initialize editedItem with current item data
      // This is a simplified approach. A more robust way would be to have a form library.
      const initialValues = {
        name: item.name || '',
        status: item.status || 1,
        rating: item.rating !== null ? item.rating : '',
        description: item.description || '',
        notes: item.notes || '',
        cover_url: item.cover_url || '',
        source_url: item.source_url || '',
        priority: item.priority || 0,
        // Values for custom fields will be handled separately
        values: item.values ? [...item.values] : []
      };
      setEditedItem(initialValues);
    }
  }, [item, editMode]);

  const handleInputChange = (field, value) => {
    setEditedItem(prev => ({ ...prev, [field]: value }));
  };

  const handleValueChange = (fieldId, newValue) => {
    setEditedItem(prev => {
      const newValues = [...(prev.values || [])];
      const index = newValues.findIndex(v => v.field_id === fieldId);
      if (index >= 0) {
        newValues[index] = { ...newValues[index], value: newValue };
      } else {
        newValues.push({ field_id: fieldId, value: newValue });
      }
      return { ...prev, values: newValues };
    });
  };

  const handleSave = () => {
    // Prepare data for API
    const updateData = {
      // id: itemId, // ID is in the URL
      item: {
        name: editedItem.name,
        status: editedItem.status,
        rating: editedItem.rating === '' ? null : parseFloat(editedItem.rating),
        description: editedItem.description,
        notes: editedItem.notes,
        cover_url: editedItem.cover_url,
        source_url: editedItem.source_url,
        priority: parseInt(editedItem.priority, 10),
        values: editedItem.values // This includes the custom field values
      },
      category_id: item.category_id // Assuming category_id is not editable here
    };

    updateItem({ id: itemId, ...updateData }, {
      onSuccess: () => {
        setEditMode(false);
        setSnackbar({ open: true, message: 'Item updated successfully!', severity: 'success' });
        refetch(); // Refetch to get updated data
      },
      onError: (error) => {
        console.error("Failed to update item:", error);
        setSnackbar({ open: true, message: `Error: ${error.message}`, severity: 'error' });
      }
    });
  };

  const handleDelete = () => {
    deleteItem(itemId, {
      onSuccess: () => {
        setOpenDeleteDialog(false);
        setSnackbar({ open: true, message: 'Item deleted successfully!', severity: 'success' });
        // TODO: Navigate back or to a list page
        // This requires React Router's useNavigate hook
      },
      onError: (error) => {
        console.error("Failed to delete item:", error);
        setSnackbar({ open: true, message: `Error: ${error.message}`, severity: 'error' });
      }
    });
  };

  const handleCloseSnackbar = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackbar({ ...snackbar, open: false });
  };

  if (isLoading) return <CircularProgress />;
  
  if (error) {
    return (
      <Container maxWidth="md">
        <Box my={4}>
          <Alert severity="error">{error.message}</Alert>
        </Box>
      </Container>
    );
  }
  
  if (!item) {
    return (
      <Container maxWidth="md">
        <Box my={4}>
          <Alert severity="warning">Item not found.</Alert>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="md">
      <Box my={4}>
        {/* Error alerts */}
        {(updateError || deleteError) && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {updateError?.message || deleteError?.message}
          </Alert>
        )}
        
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h4" component="h1">
            {editMode ? 'Edit Item' : item.name}
          </Typography>
          {!editMode && (
            <Box>
              <Button
                variant="contained"
                color="primary"
                startIcon={<Edit />}
                onClick={() => setEditMode(true)}
                sx={{ mr: 1 }}
              >
                Edit
              </Button>
              <Button
                variant="outlined"
                color="error"
                startIcon={<Delete />}
                onClick={() => setOpenDeleteDialog(true)}
              >
                Delete
              </Button>
            </Box>
          )}
        </Box>

        {editMode ? (
          <Paper elevation={3} sx={{ p: 3 }}>
            <TextField
              label="Name"
              fullWidth
              margin="normal"
              value={editedItem.name || ''}
              onChange={(e) => handleInputChange('name', e.target.value)}
            />
            <FormControl fullWidth margin="normal">
              <InputLabel>Status</InputLabel>
              <Select
                value={editedItem.status || 1}
                onChange={(e) => handleInputChange('status', e.target.value)}
                label="Status"
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
              value={editedItem.rating}
              onChange={(e) => handleInputChange('rating', e.target.value)}
              inputProps={{ min: 0, max: 10, step: 0.1 }}
            />
            <TextField
              label="Description"
              fullWidth
              margin="normal"
              multiline
              rows={3}
              value={editedItem.description || ''}
              onChange={(e) => handleInputChange('description', e.target.value)}
            />
            <TextField
              label="Notes"
              fullWidth
              margin="normal"
              multiline
              rows={3}
              value={editedItem.notes || ''}
              onChange={(e) => handleInputChange('notes', e.target.value)}
            />
            <TextField
              label="Cover URL"
              fullWidth
              margin="normal"
              value={editedItem.cover_url || ''}
              onChange={(e) => handleInputChange('cover_url', e.target.value)}
            />
            <TextField
              label="Source URL"
              fullWidth
              margin="normal"
              value={editedItem.source_url || ''}
              onChange={(e) => handleInputChange('source_url', e.target.value)}
            />
            <TextField
              label="Priority"
              type="number"
              fullWidth
              margin="normal"
              value={editedItem.priority || 0}
              onChange={(e) => handleInputChange('priority', e.target.value)}
              inputProps={{ min: 0 }}
            />

            <Divider sx={{ my: 2 }} />
            <Typography variant="h6" gutterBottom>
              Custom Fields
            </Typography>
            {item.values && item.values.length > 0 ? (
              item.values.map((fieldValue) => {
                const field = fieldValue; // fieldValue object contains field info
                return (
                  <Box key={field.field_id} mb={2}>
                    <Typography variant="subtitle1">{field.field_name} ({getFieldTypeName(field.field_type)})</Typography>
                    {field.field_type === 1 && ( // String
                      <TextField
                        fullWidth
                        variant="outlined"
                        value={fieldValue.value || ''}
                        onChange={(e) => handleValueChange(field.field_id, e.target.value)}
                      />
                    )}
                    {field.field_type === 2 && ( // Integer
                      <TextField
                        fullWidth
                        type="number"
                        variant="outlined"
                        value={fieldValue.value || ''}
                        onChange={(e) => handleValueChange(field.field_id, parseInt(e.target.value, 10))}
                      />
                    )}
                    {field.field_type === 3 && ( // Boolean
                      <FormControlLabel
                        control={
                          <Checkbox
                            checked={!!fieldValue.value}
                            onChange={(e) => handleValueChange(field.field_id, e.target.checked)}
                          />
                        }
                        label="True/False"
                      />
                    )}
                    {field.field_type === 4 && ( // Datetime
                      <TextField
                        fullWidth
                        type="datetime-local"
                        variant="outlined"
                        value={fieldValue.value ? new Date(fieldValue.value).toISOString().slice(0, 16) : ''}
                        onChange={(e) => handleValueChange(field.field_id, e.target.value)}
                        InputLabelProps={{
                          shrink: true,
                        }}
                      />
                    )}
                    {/* Add handling for array types if needed */}
                  </Box>
                );
              })
            ) : (
              <Typography variant="body2" color="textSecondary">
                No custom fields defined for this item's category.
              </Typography>
            )}

            <Box mt={3} display="flex" justifyContent="flex-end">
              <Button
                variant="outlined"
                onClick={() => setEditMode(false)}
                sx={{ mr: 1 }}
                disabled={isUpdating}
                startIcon={<Cancel />}
              >
                Cancel
              </Button>
              <Button
                variant="contained"
                color="primary"
                onClick={handleSave}
                disabled={isUpdating}
                startIcon={<Save />}
              >
                {isUpdating ? 'Saving...' : 'Save'}
              </Button>
            </Box>
          </Paper>
        ) : (
          <Paper elevation={3} sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>Details</Typography>
            <Typography><strong>Status:</strong> {getStatusText(item.status)}</Typography>
            <Typography><strong>Rating:</strong> {item.rating !== null ? item.rating : 'N/A'}</Typography>
            <Typography><strong>Description:</strong> {item.description || 'N/A'}</Typography>
            <Typography><strong>Notes:</strong> {item.notes || 'N/A'}</Typography>
            <Typography><strong>Priority:</strong> {item.priority}</Typography>
            <Typography><strong>Created:</strong> {formatDate(item.created_at)}</Typography>
            <Typography><strong>Updated:</strong> {formatDate(item.updated_at)}</Typography>
            {item.completed_at && <Typography><strong>Completed:</strong> {formatDate(item.completed_at)}</Typography>}
            
            {item.cover_url && (
              <Box mt={2}>
                <Typography><strong>Cover:</strong></Typography>
                <img src={item.cover_url} alt="Cover" style={{ maxWidth: '100%', height: 'auto' }} />
              </Box>
            )}
            {item.source_url && (
              <Box mt={1}>
                <Typography><strong>Source:</strong> <a href={item.source_url} target="_blank" rel="noopener noreferrer">{item.source_url}</a></Typography>
              </Box>
            )}

            <Divider sx={{ my: 2 }} />
            <Typography variant="h6" gutterBottom>Custom Fields</Typography>
            {item.values && item.values.length > 0 ? (
              <Box component="dl"> {/* Using description list for semantics */}
                {item.values.map((fieldValue) => (
                  <Box key={fieldValue.field_id} display="flex" alignItems="flex-start" mb={1}>
                    <Typography component="dt" variant="subtitle1" sx={{ fontWeight: 'bold', minWidth: 150 }}>
                      {fieldValue.field_name}:
                    </Typography>
                    <Typography component="dd" variant="body1" sx={{ ml: 1 }}>
                      {formatFieldValue(fieldValue.value, fieldValue.field_type)}
                    </Typography>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography variant="body2" color="textSecondary">
                No custom fields for this item.
              </Typography>
            )}
          </Paper>
        )}
      </Box>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={openDeleteDialog}
        onClose={() => setOpenDeleteDialog(false)}
        aria-labelledby="delete-dialog-title"
        aria-describedby="delete-dialog-description"
      >
        <DialogTitle id="delete-dialog-title">Confirm Delete</DialogTitle>
        <DialogContent>
          <Typography id="delete-dialog-description">
            Are you sure you want to delete "{item.name}"? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDeleteDialog(false)}>Cancel</Button>
          <Button onClick={handleDelete} color="error" variant="contained">
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      {/* Snackbar for notifications */}
      <Snackbar open={snackbar.open} autoHideDuration={6000} onClose={handleCloseSnackbar}>
        <MuiAlert onClose={handleCloseSnackbar} severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </MuiAlert>
      </Snackbar>
    </Container>
  );
};

// Simple helper to map status number to text
const getStatusText = (status) => {
  switch (status) {
    case 1: return 'To Do';
    case 2: return 'In Progress';
    case 3: return 'Paused';
    case 4: return 'Abandoned';
    case 5: return 'Completed';
    default: return 'Unknown';
  }
};

export default ItemDetailPage;