// src/redux/users/actions.ts

import { UserActionTypes } from './types';
import { 
  GetUsersRequest, 
  AddUserRequest, 
  EditUserRequest,
  UserResponse,
  AddUserResponse,
  EditUserResponse,
  DeleteUserResponse,
  RoleResponse,
  WarehouseResponse
} from './types';
import { logger } from '../../utils/logger';

// Get Users List
export const fetchUsers = (params: GetUsersRequest = {}) => {
  logger.log('info', '[Action] Fetching users list', { params });
  return {
    type: UserActionTypes.FETCH_USERS_REQUEST,
    payload: params
  };
};

export const fetchUsersSuccess = (users: UserResponse[], total: number) => ({
  type: UserActionTypes.FETCH_USERS_SUCCESS,
  payload: { users, total }
});

export const fetchUsersFailure = (error: string) => ({
  type: UserActionTypes.FETCH_USERS_FAILURE,
  payload: error
});

// Get Single User
export const fetchUser = (userID: string) => {
  logger.log('info', '[Action] Fetching user details', { userID });
  return {
    type: UserActionTypes.FETCH_USER_REQUEST,
    payload: userID
  };
};

export const fetchUserSuccess = (user: UserResponse) => ({
  type: UserActionTypes.FETCH_USER_SUCCESS,
  payload: user
});

export const fetchUserFailure = (error: string) => ({
  type: UserActionTypes.FETCH_USER_FAILURE,
  payload: error
});

// Add User
export const addUser = (user: AddUserRequest) => {
  logger.log('info', '[Action] Adding new user', { userID: user.userID });
  return {
    type: UserActionTypes.ADD_USER_REQUEST,
    payload: user
  };
};

export const addUserSuccess = (response: AddUserResponse) => ({
  type: UserActionTypes.ADD_USER_SUCCESS,
  payload: response
});

export const addUserFailure = (error: string) => ({
  type: UserActionTypes.ADD_USER_FAILURE,
  payload: error
});

// Edit User
export const editUser = (userID: string, userData: EditUserRequest) => {
  logger.log('info', '[Action] Editing user', { userID });
  return {
    type: UserActionTypes.EDIT_USER_REQUEST,
    payload: { userID, userData }
  };
};

export const editUserSuccess = (response: EditUserResponse) => ({
  type: UserActionTypes.EDIT_USER_SUCCESS,
  payload: response
});

export const editUserFailure = (error: string) => ({
  type: UserActionTypes.EDIT_USER_FAILURE,
  payload: error
});

// Delete User
export const deleteUser = (userID: string) => {
  logger.log('info', '[Action] Deleting user', { userID });
  return {
    type: UserActionTypes.DELETE_USER_REQUEST,
    payload: userID
  };
};

export const deleteUserSuccess = (response: DeleteUserResponse) => ({
  type: UserActionTypes.DELETE_USER_SUCCESS,
  payload: response
});

export const deleteUserFailure = (error: string) => ({
  type: UserActionTypes.DELETE_USER_FAILURE,
  payload: error
});

// Fetch Roles
export const fetchRoles = () => {
  logger.log('info', '[Action] Fetching roles');
  return {
    type: UserActionTypes.FETCH_ROLES_REQUEST
  };
};

export const fetchRolesSuccess = (roles: RoleResponse[]) => ({
  type: UserActionTypes.FETCH_ROLES_SUCCESS,
  payload: roles
});

export const fetchRolesFailure = (error: string) => ({
  type: UserActionTypes.FETCH_ROLES_FAILURE,
  payload: error
});

// Fetch Warehouses
export const fetchWarehouses = () => {
  logger.log('info', '[Action] Fetching warehouses');
  return {
    type: UserActionTypes.FETCH_WAREHOUSES_REQUEST
  };
};

export const fetchWarehousesSuccess = (warehouses: WarehouseResponse[]) => ({
  type: UserActionTypes.FETCH_WAREHOUSES_SUCCESS,
  payload: warehouses
});

export const fetchWarehousesFailure = (error: string) => ({
  type: UserActionTypes.FETCH_WAREHOUSES_FAILURE,
  payload: error
});

// Reset Error
export const resetUserError = () => ({
  type: UserActionTypes.RESET_USER_ERROR
});