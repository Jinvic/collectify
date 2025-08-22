// src/components/Navbar.js
import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { 
  AppBar, Toolbar, Typography, Button, Box, IconButton,
  Menu, MenuItem, Divider
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';
import LoginDialog from './LoginDialog';
import LogoutDialog from './LogoutDialog';
import UpdateUserDialog from './UpdateUserDialog';

const Navbar = () => {
  const location = useLocation();
  const { user, logout, authEnabled } = useAuth();
  const [loginOpen, setLoginOpen] = useState(false);
  const [logoutOpen, setLogoutOpen] = useState(false);
  const [updateUserOpen, setUpdateUserOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);

  const navItems = [
    { label: 'Index', path: '/' },
    { label: 'Categories', path: '/categories' },
    // { label: 'Tags', path: '/tags' }, // Placeholder
    // { label: 'Collections', path: '/collections' }, // Placeholder
    // { label: 'Search', path: '/search' }, // Placeholder
  ];

  const handleLogout = () => {
    logout();
    setLogoutOpen(false);
  };

  const handleMenuOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  // 如果未启用认证，不显示登录相关组件
  if (!authEnabled) {
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
  }

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
          {user ? (
            <>
              <Button
                color="inherit"
                variant="outlined"
                onClick={handleMenuOpen}
                aria-controls="user-menu"
                aria-haspopup="true"
              >
                {user.username}
              </Button>
              <Menu
                id="user-menu"
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleMenuClose}
                MenuListProps={{
                  'aria-labelledby': 'basic-button',
                }}
              >
                <MenuItem onClick={() => {
                  setUpdateUserOpen(true);
                  handleMenuClose();
                }}>
                  Update
                </MenuItem>
                <Divider />
                <MenuItem onClick={() => {
                  setLogoutOpen(true);
                  handleMenuClose();
                }}>
                  Logout
                </MenuItem>
              </Menu>
            </>
          ) : (
            <Button
              color="inherit"
              variant="outlined"
              onClick={() => setLoginOpen(true)}
            >
              Login
            </Button>
          )}
        </Box>
      </Toolbar>
      
      <LoginDialog open={loginOpen} onClose={() => setLoginOpen(false)} />
      <LogoutDialog open={logoutOpen} onClose={() => setLogoutOpen(false)} onConfirm={handleLogout} />
      <UpdateUserDialog open={updateUserOpen} onClose={() => setUpdateUserOpen(false)} />
    </AppBar>
  );
};

export default Navbar;