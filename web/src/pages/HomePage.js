// src/pages/HomePage.js
import React from 'react';
import { useBasicItems } from '../hooks/useItems';
import ItemList from '../components/ItemList';
import { Container, Typography, Box } from '@mui/material';

const HomePage = () => {
  const defaultParams = { page: 1, page_size: 10 };
  const { data, isLoading, error, refetch } = useBasicItems(defaultParams);
  const items = data?.data?.list || [];
  const total = data?.data?.total || 0;

  const handlePageChange = (newPage) => {
    refetch({ ...defaultParams, page: newPage });
  };

  return (
    <Container maxWidth="lg">
      <Box my={4}>
        <Typography variant="h4" component="h1" gutterBottom>
          Home - Recent Items
        </Typography>
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
    </Container>
  );
};

export default HomePage;