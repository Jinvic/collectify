// src/components/ItemList.js
import React from 'react';
import { Link } from 'react-router-dom';
import {
  Table, TableBody, TableCell, TableContainer, TableHead, TableRow,
  Paper, CircularProgress, Typography, Button, Box, TablePagination, Alert
} from '@mui/material';

const ItemList = ({ items, total, loading, error, page, pageSize, onPageChange, showCategory = false }) => {

  const handleChangePage = (event, newPage) => {
    onPageChange(newPage + 1); // Backend uses 1-based indexing
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" my={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box my={4}>
        <Alert severity="error">Error loading items: {error.message}</Alert>
      </Box>
    );
  }

  if (!items || items.length === 0) {
    return (
      <Box my={4}>
        <Typography>No items found.</Typography>
      </Box>
    );
  }

  return (
    <Box my={4}>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="items table">
          <TableHead>
            <TableRow>
              <TableCell>Title</TableCell>
              {showCategory && <TableCell>Category</TableCell>}
              <TableCell>Status</TableCell>
              <TableCell>Rating</TableCell>
              <TableCell>Updated At</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {items.map((item) => (
              <TableRow key={item.id} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                <TableCell component="th" scope="row">
                  {item.title}
                </TableCell>
                {showCategory && (
                  <TableCell>
                    {item.category?.name || 'N/A'}
                  </TableCell>
                )}
                <TableCell>{getStatusText(item.status)}</TableCell>
                <TableCell>{item.rating !== null ? item.rating : 'N/A'}</TableCell>
                <TableCell>{new Date(item.updated_at).toLocaleString()}</TableCell>
                <TableCell>
                  <Button
                    component={Link}
                    to={`/items/${item.id}`}
                    variant="outlined"
                    size="small"
                  >
                    View
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      {total > 0 && (
        <TablePagination
          component="div"
          count={total}
          page={page - 1} // MUI uses 0-based indexing
          onPageChange={handleChangePage}
          rowsPerPage={pageSize}
          rowsPerPageOptions={[5, 10, 20]} // Allow user to change page size?
          onRowsPerPageChange={(e) => {
            // Handle rows per page change if needed
            // This would require updating the parent component's state
          }}
        />
      )}
    </Box>
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

export default ItemList;