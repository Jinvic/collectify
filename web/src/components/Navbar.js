// src/components/Navbar.js
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { AppBar, Toolbar, Typography, Button, Box, IconButton } from '@mui/material';

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
        <IconButton
          component={Link}
          to="/"
          color="inherit"
          sx={{ mr: 1 }}
        >
          <img src="/logo192.png" alt="Collectify Logo" style={{ height: 40 }} />
        </IconButton>
        <Typography 
          variant="h6" 
          component={Link} 
          to="/" 
          sx={{ 
            flexGrow: 1, 
            textDecoration: 'none', 
            color: 'inherit',
            display: 'flex',
            alignItems: 'center'
          }}
        >
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