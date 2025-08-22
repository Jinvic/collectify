import apiClient from './apiClient';

export const authService = {
  login: async (username, password) => {
    try {
      const response = await apiClient.post('/user/login', { username, password });
      return response;
    } catch (error) {
      throw new Error(`Login failed: ${error.message}`);
    }
  },
  
  // 检查认证是否启用
  isAuthEnabled: async () => {
    try {
      const response = await apiClient.get('/user/enabled');
      return response;
    } catch (error) {
      throw new Error(`Failed to check auth status: ${error.message}`);
    }
  },
  
  // 更新用户信息
  updateUser: async (userData) => {
    try {
      const response = await apiClient.post('/user/update', userData);
      return response;
    } catch (error) {
      throw new Error(`Failed to update user: ${error.message}`);
    }
  }
};