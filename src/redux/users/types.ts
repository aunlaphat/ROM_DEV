// src/redux/users/types.ts

export enum UserActionTypes {
  // Get Users List
  FETCH_USERS_REQUEST = "@USER/FETCH_USERS_REQUEST",
  FETCH_USERS_SUCCESS = "@USER/FETCH_USERS_SUCCESS",
  FETCH_USERS_FAILURE = "@USER/FETCH_USERS_FAILURE",

  // Get Single User
  FETCH_USER_REQUEST = "@USER/FETCH_USER_REQUEST",
  FETCH_USER_SUCCESS = "@USER/FETCH_USER_SUCCESS",
  FETCH_USER_FAILURE = "@USER/FETCH_USER_FAILURE",

  // Add User
  ADD_USER_REQUEST = "@USER/ADD_USER_REQUEST",
  ADD_USER_SUCCESS = "@USER/ADD_USER_SUCCESS",
  ADD_USER_FAILURE = "@USER/ADD_USER_FAILURE",

  // Edit User
  EDIT_USER_REQUEST = "@USER/EDIT_USER_REQUEST",
  EDIT_USER_SUCCESS = "@USER/EDIT_USER_SUCCESS",
  EDIT_USER_FAILURE = "@USER/EDIT_USER_FAILURE",

  // Delete User
  DELETE_USER_REQUEST = "@USER/DELETE_USER_REQUEST",
  DELETE_USER_SUCCESS = "@USER/DELETE_USER_SUCCESS",
  DELETE_USER_FAILURE = "@USER/DELETE_USER_FAILURE",

  // Fetch Roles
  FETCH_ROLES_REQUEST = "@USER/FETCH_ROLES_REQUEST",
  FETCH_ROLES_SUCCESS = "@USER/FETCH_ROLES_SUCCESS",
  FETCH_ROLES_FAILURE = "@USER/FETCH_ROLES_FAILURE",

  // Fetch Warehouses
  FETCH_WAREHOUSES_REQUEST = "@USER/FETCH_WAREHOUSES_REQUEST",
  FETCH_WAREHOUSES_SUCCESS = "@USER/FETCH_WAREHOUSES_SUCCESS",
  FETCH_WAREHOUSES_FAILURE = "@USER/FETCH_WAREHOUSES_FAILURE",

  // Reset Error
  RESET_USER_ERROR = "@USER/RESET_USER_ERROR",
}

// Response Types - สอดคล้องกับ Backend
export interface UserResponse {
  userID: string;
  userName: string;
  nickName: string;
  fullNameTH: string;
  departmentNo: string;
  roleID: number;          // int ตาม backend
  roleName: string;
  warehouseID: number;     // int ตาม backend
  warehouseName: string;
  description: string;
  isActive: boolean;
  lastLoginAt?: string;    // optional
  createdAt: string;       // ไม่ใช่ optional
  updatedAt?: string;      // optional
}

export interface AddUserResponse {
  userID: string;
  roleID: number;          // int ตาม backend
  warehouseID: number;     // int ตาม backend
  createdBy: string;
}

export interface EditUserResponse {
  userID: string;
  roleID?: number;         // optional int ตาม backend
  roleName: string;
  warehouseID?: number;    // optional int ตาม backend
  warehouseName: string;
  updatedBy: string;
  updatedAt: string;
}

export interface DeleteUserResponse {
  userID: string;
  userName: string;
  roleID: number;          // int ตาม backend
  roleName: string;
  warehouseID: number;     // int ตาม backend
  warehouseName: string;
  deactivatedBy: string;
  deactivatedAt: string;
  message: string;
}

// Role และ Warehouse Response
export interface RoleResponse {
  roleID: number;
  roleName: string;
  description: string;
}

export interface WarehouseResponse {
  warehouseID: number;
  warehouseName: string;
  location: string;
}

// Request Types - สอดคล้องกับ Backend
export interface GetUsersRequest {
  isActive?: boolean;
  limit?: number;
  offset?: number;
}

export interface AddUserRequest {
  userID: string;
  roleID: number;          // int ตาม backend
  warehouseID: number;     // int ตาม backend
}

export interface EditUserRequest {
  userID: string;
  roleID?: number;         // optional int ตาม backend
  warehouseID?: number;    // optional int ตาม backend
}

// API Response Format
export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

// User State
export interface UserState {
  users: UserResponse[];
  currentUser: UserResponse | null;
  loading: boolean;
  error: string | null;
  pagination: {
    current: number;
    pageSize: number;
    total: number;
  };
  roles: RoleResponse[];
  warehouses: WarehouseResponse[];
}

// Initial state
export const initialUserState: UserState = {
  users: [],
  currentUser: null,
  loading: false,
  error: null,
  pagination: {
    current: 1,
    pageSize: 100,
    total: 0,
  },
  roles: [],
  warehouses: [],
};