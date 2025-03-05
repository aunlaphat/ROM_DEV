import { AuthActionTypes } from "./types";
import { logger } from '../../utils/logger';

type User = {
  userID: string;
  userName: string;
  fullNameTH: string;
  nickName: string;
  roleID: number;
  roleName: string;
  departmentNo: string;
  platform: string;
};

type State = {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error?: string;
};

const initialState: State = {
  user: null,
  token: localStorage.getItem("access_token") || null,
  isAuthenticated: false,
  loading: false,
};

export default function authReducer(state = initialState, action: any) {
  switch (action.type) {
    case AuthActionTypes.AUTHEN_LOGIN_REQ:
      logger.log('info', '[Auth] Processing login request');
      return { ...state, loading: true };
      
    case AuthActionTypes.AUTHEN_LOGIN_SUCCESS:
      logger.log('info', '[Auth] Login successful', {
        user: action.users,
        previousUser: state.user
      });
      return {
        ...state,
        isAuthenticated: true,
        user: action.users,
        loading: false,
        error: null
      };
      
    case AuthActionTypes.AUTHEN_CHECK_SUCCESS:
      logger.log('info', '[Auth] Auth check successful', {
        incomingUser: action.users,
        currentUser: state.user
      });
      return {
        ...state,
        isAuthenticated: true,
        user: {
          ...action.users,
          roleID: action.users.roleID || action.users.userRoleID,
          userID: action.users.userID,
          userName: action.users.userName,
          roleName: action.users.roleName
        },
        loading: false,
      };
      
    case AuthActionTypes.AUTHEN_LOGIN_FAIL:
      logger.error('[Auth] Login failed', { message: action.message });
      return {
        ...state,
        isAuthenticated: false,
        user: null,
        token: null,
        loading: false,
        error: action.message,
      };
      
    case AuthActionTypes.AUTHEN_LOGIN_LARK_REQ:
      logger.log('info', '[Auth] Processing Lark login request');
      return { ...state, loading: true };
      
    case AuthActionTypes.AUTHEN_LOGIN_LARK_SUCCESS:
      logger.log('info', '[Auth] Lark login successful', { user: action.users });
      return {
        ...state,
        loading: false,
        isAuthenticated: true,
        user: action.users,
      };
      
    case AuthActionTypes.AUTHEN_LOGIN_LARK_FAIL:
      logger.error('[Auth] Lark login failed', { message: action.message });
      return {
        ...state,
        loading: false,
        isAuthenticated: false,
        error: action.message,
      };
      
    case AuthActionTypes.AUTHEN_LOGOUT_SUCCESS:
      logger.log('info', '[Auth] Logout successful');
      return { ...state, loading: false, isAuthenticated: false, user: null };
      
    case AuthActionTypes.AUTHEN_LOGOUT_FAIL:
      logger.error('[Auth] Logout failed', { message: action.message });
      return { ...state, loading: false, error: action.message };
      
    case AuthActionTypes.AUTHEN_CHECK_REQ:
      logger.log('info', '[Auth] Processing auth check request');
      return { 
        ...state, 
        loading: true,
        error: null
      };
      
    case AuthActionTypes.AUTHEN_CHECK_SUCCESS:
      logger.log('info', '[Auth] Auth check successful', { user: action.users });
      return {
        ...state,
        isAuthenticated: true,
        user: action.users,
        loading: false,
        error: null
      };
      
    case AuthActionTypes.AUTHEN_CHECK_FAIL:
      logger.error('[Auth] Auth check failed', { message: action.message });
      return {
        ...state,
        isAuthenticated: false,
        user: null,
        loading: false,
        error: action.message
      };
      
    default:
      return state;
  }
}
