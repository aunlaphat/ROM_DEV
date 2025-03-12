// src/redux/users/reducer.ts

import { UserActionTypes } from './types';
import { UserState, UserResponse } from './types';
import { logger } from '../../utils/logger';

// Initial state
const initialState: UserState = {
  users: [],
  currentUser: null,
  loading: false,
  error: null,
  pagination: {
    current: 1,
    pageSize: 100,
    total: 0
  },
  roles: [],
  warehouses: []
};

// User Reducer
const userReducer = (state: UserState = initialState, action: any): UserState => {
  switch (action.type) {
    // Request actions - เริ่มการโหลด
    case UserActionTypes.FETCH_USERS_REQUEST:
    case UserActionTypes.FETCH_USER_REQUEST:
    case UserActionTypes.ADD_USER_REQUEST:
    case UserActionTypes.EDIT_USER_REQUEST:
    case UserActionTypes.DELETE_USER_REQUEST:
      logger.redux.action(action.type);
      return {
        ...state,
        loading: true,
        error: null
      };

    // Fetch Users Success
    case UserActionTypes.FETCH_USERS_SUCCESS:
      logger.redux.action(action.type, {
        count: action.payload.users.length,
        total: action.payload.total
      });
      
      return {
        ...state,
        users: action.payload.users,
        loading: false,
        error: null,
        pagination: {
          ...state.pagination,
          total: action.payload.total
        }
      };

    // Fetch Single User Success
    case UserActionTypes.FETCH_USER_SUCCESS:
      logger.redux.action(action.type, {
        userID: action.payload.userID,
        userName: action.payload.userName
      });
      
      return {
        ...state,
        currentUser: action.payload,
        loading: false,
        error: null
      };

    // Add User Success
    case UserActionTypes.ADD_USER_SUCCESS:
      logger.redux.action(action.type, {
        userID: action.payload.userID,
        roleID: action.payload.roleID,
        warehouseID: action.payload.warehouseID
      });
      
      return {
        ...state,
        loading: false,
        error: null
      };

    // Edit User Success
    case UserActionTypes.EDIT_USER_SUCCESS:
      logger.redux.action(action.type, {
        userID: action.payload.userID,
        updatedBy: action.payload.updatedBy
      });
      
      return {
        ...state,
        loading: false,
        error: null,
        // อัพเดต currentUser ถ้ากำลังแก้ไขผู้ใช้ปัจจุบัน
        currentUser: state.currentUser?.userID === action.payload.userID
          ? { ...state.currentUser, ...action.payload }
          : state.currentUser
      };

    // Delete User Success
    case UserActionTypes.DELETE_USER_SUCCESS:
      logger.redux.action(action.type, {
        userID: action.payload.userID,
        userName: action.payload.userName
      });
      
      return {
        ...state,
        loading: false,
        error: null,
        // ลบผู้ใช้ออกจากรายการ
        users: state.users.filter(user => user.userID !== action.payload.userID),
        // ล้าง currentUser ถ้าเป็นผู้ใช้ที่ถูกลบ
        currentUser: state.currentUser?.userID === action.payload.userID
          ? null
          : state.currentUser
      };

    // Fetch Roles Success
    case UserActionTypes.FETCH_ROLES_SUCCESS:
      logger.redux.action(action.type, {
        count: action.payload.length
      });
      
      return {
        ...state,
        roles: action.payload,
        error: null
      };

    // Fetch Warehouses Success
    case UserActionTypes.FETCH_WAREHOUSES_SUCCESS:
      logger.redux.action(action.type, {
        count: action.payload.length
      });
      
      return {
        ...state,
        warehouses: action.payload,
        error: null
      };

    // Error actions
    case UserActionTypes.FETCH_USERS_FAILURE:
    case UserActionTypes.FETCH_USER_FAILURE:
    case UserActionTypes.ADD_USER_FAILURE:
    case UserActionTypes.EDIT_USER_FAILURE:
    case UserActionTypes.DELETE_USER_FAILURE:
    case UserActionTypes.FETCH_ROLES_FAILURE:
    case UserActionTypes.FETCH_WAREHOUSES_FAILURE:
      logger.redux.action(action.type, {
        error: action.payload
      });
      
      return {
        ...state,
        loading: false,
        error: action.payload
      };

    // Reset Error
    case UserActionTypes.RESET_USER_ERROR:
      return {
        ...state,
        error: null
      };

    default:
      return state;
  }
};

export default userReducer;