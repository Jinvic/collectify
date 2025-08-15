// src/components/Navbar.js
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';

const Navbar = () => {
  const location = useLocation();

  const navItems = [
    { label: 'Index', path: '/' },
    { label: 'Categories', path: '/categories' },
    // { label: 'Tags', path: '/tags' }, // Placeholder
    // { label: 'Collections', path: '/collections' }, // Placeholder
    // { label: 'Search', path: '/search' }, // Placeholder
  ];

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          Collectify
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          {navItems.map((item) => (
            <Button
              key={item.path}
              component={Link}
              to={item.path}
              color="inherit"
              variant={location.pathname === item.path ? 'outlined' : 'text'} // Highlight active link
            >
              {item.label}
            </Button>
          ))}
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Navbar;