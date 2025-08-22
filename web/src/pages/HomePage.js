// src/pages/HomePage.js
import React, { useState } from 'react';
import { useBasicItems } from '../hooks/useItems';
import ItemList from '../components/ItemList';
import AddItemDialog from '../components/AddItemDialog';
import { Container, Typography, Box, Button, Snackbar } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';

const HomePage = () => {
  const defaultParams = { page: 1, page_size: 10 };
  const { data, isLoading, error, refetch } = useBasicItems(defaultParams);
  const items = data?.data?.list || [];
  const total = data?.data?.total || 0;

  const [openAddDialog, setOpenAddDialog] = useState(false);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  const handlePageChange = (newPage) => {
    refetch({ ...defaultParams, page: newPage });
  };

  const handleItemAdded = () => {
    // Refresh the item list after adding a new item
    refetch();
    setSnackbar({ open: true, message: 'Item added successfully!', severity: 'success' });
  };

  const handleCloseSnackbar = (event, reason) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackbar({ ...snackbar, open: false });
  };

  return (
    <Container maxWidth="lg">
      <Box my={4}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h4" component="h1" gutterBottom>
            Home - Recent Items
          </Typography>
          <Button
            variant="contained"
            color="primary"
            startIcon={<AddIcon />}
            onClick={() => setOpenAddDialog(true)}
          >
            Add Item
          </Button>
        </Box>
        <ItemList
          items={items}
          total={total}
          loading={isLoading}
          error={error}
          page={defaultParams.page}
          pageSize={defaultParams.page_size}
          onPageChange={handlePageChange}
          showCategory={true} // Show category name on home page
        />
      </Box>

      <AddItemDialog
        open={openAddDialog}
        onClose={() => setOpenAddDialog(false)}
        onItemAdded={handleItemAdded}
      />
      
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

export default HomePage;