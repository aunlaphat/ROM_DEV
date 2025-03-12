// src/redux/users/hooks.ts

import { useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  fetchUsers,
  fetchUser,
  addUser,
  editUser,
  deleteUser,
  fetchRoles,
  fetchWarehouses,
  resetUserError
} from './action';
import {
  getAllUsers,
  getCurrentUser,
  getUserLoading,
  getUserError,
  getUserPagination,
  getRoles,
  getWarehouses
} from './selector';
import { AddUserRequest, EditUserRequest, GetUsersRequest } from './types';

/**
 * Custom hook สำหรับใช้งานในหน้า User Management
 */
export const useUsers = () => {
  const dispatch = useDispatch();
  
  // Selectors
  const users = useSelector(getAllUsers);
  const currentUser = useSelector(getCurrentUser);
  const loading = useSelector(getUserLoading);
  const error = useSelector(getUserError);
  const pagination = useSelector(getUserPagination);
  const roles = useSelector(getRoles);
  const warehouses = useSelector(getWarehouses);
  
  // Actions
  const handleFetchUsers = useCallback((params: GetUsersRequest = {}) => {
    dispatch(fetchUsers(params));
  }, [dispatch]);
  
  const handleFetchUser = useCallback((userID: string) => {
    dispatch(fetchUser(userID));
  }, [dispatch]);
  
  const handleAddUser = useCallback((userData: AddUserRequest) => {
    dispatch(addUser(userData));
  }, [dispatch]);
  
  const handleEditUser = useCallback((userID: string, userData: EditUserRequest) => {
    dispatch(editUser(userID, userData));
  }, [dispatch]);
  
  const handleDeleteUser = useCallback((userID: string) => {
    dispatch(deleteUser(userID));
  }, [dispatch]);
  
  const handleFetchRoles = useCallback(() => {
    dispatch(fetchRoles());
  }, [dispatch]);
  
  const handleFetchWarehouses = useCallback(() => {
    dispatch(fetchWarehouses());
  }, [dispatch]);
  
  const handleResetError = useCallback(() => {
    dispatch(resetUserError());
  }, [dispatch]);
  
  return {
    // State
    users,
    currentUser,
    loading,
    error,
    pagination,
    roles,
    warehouses,
    
    // Actions
    fetchUsers: handleFetchUsers,
    fetchUser: handleFetchUser,
    addUser: handleAddUser,
    editUser: handleEditUser,
    deleteUser: handleDeleteUser,
    fetchRoles: handleFetchRoles,
    fetchWarehouses: handleFetchWarehouses,
    resetError: handleResetError
  };
};