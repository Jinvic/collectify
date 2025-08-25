// src/pages/CollectionDetailPage.js
import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useCollection, useUpdateCollection, useDeleteCollection } from '../hooks/useCollections';
import { useSearchItems } from '../hooks/useItems';
import ItemList from '../components/ItemList';
import {
  Container, Typography, Box, CircularProgress, Alert, Button, Dialog, DialogTitle,
  DialogContent, DialogActions, TextField, IconButton, Snackbar
} from '@mui/material';
import { Edit, Delete } from '@mui/icons-material';

const CollectionDetailPage = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const collectionId = parseInt(id, 10);

  const { data: collectionData, isLoading: isCollectionLoading, error: collectionError } = useCollection(collectionId);
  const collection = collectionData?.data;
  
  const { mutate: updateCollection, error: updateError } = useUpdateCollection();
  const { mutate: deleteCollection, error: deleteCollectionError } = useDeleteCollection();

  const defaultSearchParams = { collection_ids: [collectionId], page: 1, page_size: 10 };
  const { data: itemsData, isLoading: isItemsLoading, error: itemsError, refetch: refetchItems } = useSearchItems(defaultSearchParams);
  const items = itemsData?.data?.list || [];
  const totalItems = itemsData?.data?.total || 0;

  const [openEditDialog, setOpenEditDialog] = useState(false);
  const [editName, setEditName] = useState(collection?.name || '');
  const [editDescription, setEditDescription] = useState(collection?.description || '');
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  // Sync collection data when it loads
  React.useEffect(() => {
    if (collection?.name) {
      setEditName(collection.name);
      setEditDescription(collection.description || '');
    }
  }, [collection?.name, collection?.description]);

  const handleUpdate = () => {
    if (editName.trim() && (editName !== collection.name || editDescription !== collection.description)) {
      updateCollection({ 
        id: collectionId, 
        name: editName,
        description: editDescription
      }, {
        onSuccess: () => {
          setOpenEditDialog(false);
          setSnackbar({ open: true, message: 'Collection updated successfully!', severity: 'success' });
        },
        onError: (error) => {
          console.error("Failed to update collection:", error);
          setSnackbar({ open: true, message: `Error: ${error.message}`, severity: 'error' });
        }
      });
    } else {
      setOpenEditDialog(false);
    }
  };

  const handleDelete = () => {
    deleteCollection(collectionId, {
      onSuccess: () => {
        setOpenDeleteDialog(false);
        setSnackbar({ open: true, message: 'Collection deleted successfully!', severity: 'success' });
        // Navigate back to collections list
        navigate('/collections');
      },
      onError: (error) => {
        console.error("Failed to delete collection:", error);
        setSnackbar({ open: true, message: `Error: ${error.message}`, severity: 'error' });
        setOpenDeleteDialog(false);
      }
    });
  };

  const handleItemsPageChange = (newPage) => {
    refetchItems({ ...defaultSearchParams, page: newPage });
  };

  const handleCloseSnackbar = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackbar({ ...snackbar, open: false });
  };

  if (isCollectionLoading) return <CircularProgress />;
  
  if (collectionError) {
    return (
      <Container maxWidth="lg">
        <Box my={4}>
          <Alert severity="error">{collectionError.message}</Alert>
        </Box>
      </Container>
    );
  }
  
  if (!collection) {
    return (
      <Container maxWidth="lg">
        <Box my={4}>
          <Alert severity="warning">Collection not found.</Alert>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box my={4}>
        {/* Error alerts */}
        {(updateError || deleteCollectionError) && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {updateError?.message || deleteCollectionError?.message}
          </Alert>
        )}
        
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h4" component="h1">
            {collection.name}
          </Typography>
          <Box>
            <IconButton onClick={() => setOpenEditDialog(true)} size="small" sx={{ mr: 1 }}>
              <Edit />
            </IconButton>
            <IconButton 
              onClick={() => setOpenDeleteDialog(true)} 
              size="small"
              color="error"
            >
              <Delete />
            </IconButton>
          </Box>
        </Box>

        <Typography variant="body1" gutterBottom>
          {collection.description || 'No description'}
        </Typography>

        <Typography variant="h6" gutterBottom mt={4}>
          Items in this Collection
        </Typography>
        <ItemList
          items={items}
          total={totalItems}
          loading={isItemsLoading}
          error={itemsError}
          page={defaultSearchParams.page}
          pageSize={defaultSearchParams.page_size}
          onPageChange={handleItemsPageChange}
          showCategory={true}
        />
      </Box>

      {/* Edit Collection Dialog */}
      <Dialog open={openEditDialog} onClose={() => setOpenEditDialog(false)}>
        <DialogTitle>Edit Collection</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Name"
            type="text"
            fullWidth
            variant="outlined"
            value={editName}
            onChange={(e) => setEditName(e.target.value)}
          />
          <TextField
            margin="dense"
            label="Description"
            type="text"
            fullWidth
            variant="outlined"
            multiline
            rows={3}
            value={editDescription}
            onChange={(e) => setEditDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenEditDialog(false)}>Cancel</Button>
          <Button onClick={handleUpdate} disabled={!editName.trim()}>
            Save
          </Button>
        </DialogActions>
      </Dialog>

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
            Are you sure you want to delete "{collection.name}"? This action cannot be undone.
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
      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        message={snackbar.message}
      />
    </Container>
  );
};

export default CollectionDetailPage;