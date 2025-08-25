// src/pages/TagDetailPage.js
import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useTag, useRenameTag, useDeleteTag } from '../hooks/useTags';
import { useSearchItems } from '../hooks/useItems';
import ItemList from '../components/ItemList';
import {
  Container, Typography, Box, CircularProgress, Alert, Button, Dialog, DialogTitle,
  DialogContent, DialogActions, TextField, IconButton, Snackbar
} from '@mui/material';
import { Edit, Delete } from '@mui/icons-material';

const TagDetailPage = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const tagId = parseInt(id, 10);

  const { data: tagData, isLoading: isTagLoading, error: tagError } = useTag(tagId);
  const tag = tagData?.data;
  
  const { mutate: renameTag, error: renameError } = useRenameTag();
  const { mutate: deleteTag, error: deleteTagError } = useDeleteTag();

  const defaultSearchParams = { tag_ids: [tagId], page: 1, page_size: 10 };
  const { data: itemsData, isLoading: isItemsLoading, error: itemsError, refetch: refetchItems } = useSearchItems(defaultSearchParams);
  const items = itemsData?.data?.list || [];
  const totalItems = itemsData?.data?.total || 0;

  const [openEditDialog, setOpenEditDialog] = useState(false);
  const [editName, setEditName] = useState(tag?.name || '');
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  // Sync tag name when data loads
  React.useEffect(() => {
    if (tag?.name) {
      setEditName(tag.name);
    }
  }, [tag?.name]);

  const handleRename = () => {
    if (editName.trim() && editName !== tag.name) {
      renameTag({ id: tagId, name: editName }, {
        onSuccess: () => {
          setOpenEditDialog(false);
          setSnackbar({ open: true, message: 'Tag renamed successfully!', severity: 'success' });
        },
        onError: (error) => {
          console.error("Failed to rename tag:", error);
          setSnackbar({ open: true, message: `Error: ${error.message}`, severity: 'error' });
        }
      });
    } else {
      setOpenEditDialog(false);
    }
  };

  const handleDelete = () => {
    deleteTag(tagId, {
      onSuccess: () => {
        setOpenDeleteDialog(false);
        setSnackbar({ open: true, message: 'Tag deleted successfully!', severity: 'success' });
        // Navigate back to tags list
        navigate('/tags');
      },
      onError: (error) => {
        console.error("Failed to delete tag:", error);
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

  if (isTagLoading) return <CircularProgress />;
  
  if (tagError) {
    return (
      <Container maxWidth="lg">
        <Box my={4}>
          <Alert severity="error">{tagError.message}</Alert>
        </Box>
      </Container>
    );
  }
  
  if (!tag) {
    return (
      <Container maxWidth="lg">
        <Box my={4}>
          <Alert severity="warning">Tag not found.</Alert>
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box my={4}>
        {/* Error alerts */}
        {(renameError || deleteTagError) && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {renameError?.message || deleteTagError?.message}
          </Alert>
        )}
        
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h4" component="h1">
            {tag.name}
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

        <Typography variant="h6" gutterBottom mt={4}>
          Items with this Tag
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

      {/* Edit Tag Dialog */}
      <Dialog open={openEditDialog} onClose={() => setOpenEditDialog(false)}>
        <DialogTitle>Edit Tag</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Tag Name"
            type="text"
            fullWidth
            variant="outlined"
            value={editName}
            onChange={(e) => setEditName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenEditDialog(false)}>Cancel</Button>
          <Button onClick={handleRename} disabled={!editName.trim() || editName === tag.name}>
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
            Are you sure you want to delete "{tag.name}"? This action cannot be undone.
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

export default TagDetailPage;