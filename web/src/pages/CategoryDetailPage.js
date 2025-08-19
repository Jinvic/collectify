// src/pages/CategoryDetailPage.js
import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useCategory, useRenameCategory } from '../hooks/useCategories';
import { useCreateField, useDeleteField } from '../hooks/useFields';
import { useSearchItems } from '../hooks/useItems';
import ItemList from '../components/ItemList';
import { getFieldTypeName } from '../utils/itemUtils';
import {
  Container, Typography, Box, CircularProgress, Alert, Button, Dialog, DialogTitle,
  DialogContent, DialogActions, TextField, FormControl, InputLabel, Select, MenuItem,
  Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, IconButton
} from '@mui/material';
import { Edit, Delete, Add } from '@mui/icons-material';

const CategoryDetailPage = () => {
  const { id } = useParams();
  const categoryId = parseInt(id, 10);

  const { data: categoryData, isLoading: isCategoryLoading, error: categoryError } = useCategory(categoryId);
  const category = categoryData?.data;
  
  const { mutate: renameCategory } = useRenameCategory();
  const { mutate: createField } = useCreateField();
  const { mutate: deleteField } = useDeleteField();

  const defaultSearchParams = { category_id: categoryId, page: 1, page_size: 10 };
  const { data: itemsData, isLoading: isItemsLoading, error: itemsError, refetch: refetchItems } = useSearchItems(defaultSearchParams);
  const items = itemsData?.data?.list || [];
  const totalItems = itemsData?.data?.total || 0;

  const [editName, setEditName] = useState(false);
  const [newName, setNewName] = useState(category?.name || '');
  const [openFieldDialog, setOpenFieldDialog] = useState(false);
  const [newField, setNewField] = useState({ name: '', type: 1, is_array: false, required: false });

  // Sync category name when data loads
  React.useEffect(() => {
    if (category?.name) {
      setNewName(category.name);
    }
  }, [category?.name]);

  const handleRename = () => {
    if (newName.trim() && newName !== category.name) {
      renameCategory({ id: categoryId, name: newName }, {
        onSuccess: () => {
          setEditName(false);
        },
        onError: (error) => {
          console.error("Failed to rename category:", error);
          // TODO: Show error to user
        }
      });
    } else {
      setEditName(false);
    }
  };

  const handleCreateField = () => {
    // Prevent creating boolean array fields
    if (newField.type === 3 && newField.is_array) {
      alert("Boolean type fields cannot be arrays.");
      return;
    }
    
    if (newField.name.trim()) {
      createField({ ...newField, category_id: categoryId }, {
        onSuccess: () => {
          setOpenFieldDialog(false);
          setNewField({ name: '', type: 1, is_array: false, required: false });
        },
        onError: (error) => {
          console.error("Failed to create field:", error);
          alert(`Failed to create field: ${error.message}`);
        }
      });
    }
  };

  const handleDeleteField = (fieldId) => {
    // Simple confirmation, could be improved with a proper dialog
    if (window.confirm('Are you sure you want to delete this field?')) {
      deleteField(fieldId, {
        onError: (error) => {
          console.error("Failed to delete field:", error);
          // TODO: Show error to user
        }
      });
    }
  };

  const handleItemsPageChange = (newPage) => {
    refetchItems({ ...defaultSearchParams, page: newPage });
  };

  if (isCategoryLoading) return <CircularProgress />;
  if (categoryError) return <Alert severity="error">{categoryError.message}</Alert>;
  if (!category) return <Alert severity="warning">Category not found.</Alert>;

  return (
    <Container maxWidth="lg">
      <Box my={4}>
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
          {editName ? (
            <TextField
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              onBlur={handleRename}
              onKeyDown={(e) => e.key === 'Enter' && handleRename()}
              autoFocus
              variant="outlined"
              size="small"
            />
          ) : (
            <Typography variant="h4" component="h1">
              {category.name}
            </Typography>
          )}
          <IconButton onClick={() => setEditName(true)} size="small">
            <Edit />
          </IconButton>
        </Box>

        <Typography variant="h6" gutterBottom>
          Fields
        </Typography>
        {category.fields && category.fields.length > 0 ? (
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Name</TableCell>
                  <TableCell>Type</TableCell>
                  <TableCell>Array</TableCell>
                  <TableCell>Required</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {category.fields.map((field) => (
                  <TableRow key={field.id}>
                    <TableCell>{field.name}</TableCell>
                    <TableCell>{getFieldTypeName(field.type)}</TableCell>
                    <TableCell>{field.is_array ? 'Yes' : 'No'}</TableCell>
                    <TableCell>{field.required ? 'Yes' : 'No'}</TableCell>
                    <TableCell>
                      <IconButton
                        aria-label="delete"
                        size="small"
                        onClick={() => handleDeleteField(field.id)}
                      >
                        <Delete />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        ) : (
          <Typography variant="body2" color="textSecondary">
            No fields defined for this category.
          </Typography>
        )}
        <Box mt={2}>
          <Button
            variant="outlined"
            startIcon={<Add />}
            onClick={() => setOpenFieldDialog(true)}
          >
            Add Field
          </Button>
        </Box>

        <Typography variant="h6" gutterBottom mt={4}>
          Items in this Category
        </Typography>
        <ItemList
          items={items}
          total={totalItems}
          loading={isItemsLoading}
          error={itemsError}
          page={defaultSearchParams.page}
          pageSize={defaultSearchParams.page_size}
          onPageChange={handleItemsPageChange}
          showCategory={false} // Don't show category column as we are already in category context
        />
      </Box>

      {/* Add Field Dialog */}
      <Dialog open={openFieldDialog} onClose={() => setOpenFieldDialog(false)}>
        <DialogTitle>Add New Field</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Field Name"
            type="text"
            fullWidth
            variant="outlined"
            value={newField.name}
            onChange={(e) => setNewField({ ...newField, name: e.target.value })}
          />
          <FormControl fullWidth margin="dense" variant="outlined">
            <InputLabel>Type</InputLabel>
            <Select
              value={newField.type}
              onChange={(e) => setNewField({ ...newField, type: e.target.value })}
              label="Type"
            >
              <MenuItem value={1}>String</MenuItem>
              <MenuItem value={2}>Integer</MenuItem>
              <MenuItem value={3}>Boolean</MenuItem>
              <MenuItem value={4}>Datetime</MenuItem>
            </Select>
          </FormControl>
          <Box display="flex" justifyContent="space-between" mt={2}>
            <Button
              variant={newField.is_array ? "contained" : "outlined"}
              onClick={() => setNewField({ ...newField, is_array: !newField.is_array })}
              disabled={newField.type === 3} // Disable array for Boolean type
            >
              Array
            </Button>
            <Button
              variant={newField.required ? "contained" : "outlined"}
              onClick={() => setNewField({ ...newField, required: !newField.required })}
            >
              Required
            </Button>
          </Box>
          {newField.type === 3 && newField.is_array && (
            <Alert severity="warning" sx={{ mt: 2 }}>
              Boolean type cannot be an array. Array option has been disabled.
            </Alert>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenFieldDialog(false)}>Cancel</Button>
          <Button onClick={handleCreateField} disabled={!newField.name.trim()}>Add</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default CategoryDetailPage;