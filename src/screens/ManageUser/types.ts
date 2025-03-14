// src/screens/ManageUser/types.ts
/**
 * Type definitions for the ManageUser module
 * Updated to match backend structures
 */

/**
 * User data interface
 */
export interface User {
  userID: string;
  userName: string;
  nickName: string;
  fullNameTH: string;
  departmentNo: string;
  roleID: number;
  roleName: string;
  warehouseID: string; // Changed to string to match backend
  warehouseName: string;
  description: string;
  isActive: boolean;
  lastLoginAt?: string;
  createdAt: string;
  updatedAt?: string;
}

/**
 * Role data interface
 */
export interface Role {
  roleID: number;
  roleName: string;
  description?: string;
}

/**
 * Warehouse data interface
 */
export interface Warehouse {
  warehouseID: string; // Changed to string to match backend
  warehouseName: string;
  location?: string;
}

/**
 * Form data for adding a new user - matches AddUserRequest
 */
export interface AddUserRequest {
  userID: string;
  roleID: number;
  warehouseID: string; // String type to match backend
}

/**
 * Form data for editing a user - matches EditUserRequest
 */
export interface EditUserRequest {
  userID: string;
  roleID?: number; // Optional to match backend
  warehouseID?: string; // Optional string to match backend
}

/**
 * Response for adding a user - matches AddUserResponse
 */
export interface AddUserResponse {
  userID: string;
  roleID: number;
  warehouseID: string; // String type to match backend
  createdBy: string;
}

/**
 * Response for editing a user - matches EditUserResponse
 */
export interface EditUserResponse {
  userID: string;
  roleID?: number; // Optional to match backend
  roleName: string;
  warehouseID?: string; // Optional string to match backend
  warehouseName: string;
  updatedBy: string;
  updatedAt: string;
}

/**
 * Response for deleting a user - matches DeleteUserResponse
 */
export interface DeleteUserResponse {
  userID: string;
  userName: string;
  roleID: number;
  roleName: string;
  warehouseID: string; // String type to match backend
  warehouseName: string;
  deactivatedBy: string;
  deactivatedAt: string;
  message: string;
}

/**
 * Pagination configuration
 */
export interface PaginationConfig {
  current: number;
  pageSize: number;
  total: number;
  showSizeChanger?: boolean;
  showQuickJumper?: boolean;
  pageSizeOptions?: string[];
}

/**
 * Generic API response
 */
export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

// Component Props Interfaces

/**
 * Props for UserHeader component
 */
export interface UserHeaderProps {
  onSearch: (value: string) => void;
  onAddUser: () => void;
}

/**
 * Props for UserTable component
 */
export interface UserTableProps {
  users: User[];
  loading: boolean;
  pagination: PaginationConfig;
  onEdit: (user: User) => void;
  onDelete: (userID: string) => void;
  onChange: (pagination: any) => void;
}

/**
 * Props for UserActions component
 */
export interface UserActionsProps {
  user: User;
  onEdit: () => void;
  onDelete: () => void;
}

/**
 * Props for UserForm component
 */
export interface UserFormProps {
  visible: boolean;
  user: User | null;
  roles: Role[];
  warehouses: Warehouse[];
  onSave: (values: AddUserRequest | EditUserRequest) => void;
  onCancel: () => void;
}