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
      console.log('[Auth] Processing login request');
      return { ...state, loading: true };
      
    case AuthActionTypes.AUTHEN_LOGIN_SUCCESS:
      logger.auth('info', 'Updating auth state after login', {
        user: action.users,
        previous: state.user
      });
      return {
        ...state,
        isAuthenticated: true,
        user: action.users,
        loading: false,
        error: null
      };
    case AuthActionTypes.AUTHEN_CHECK_SUCCESS:
      console.log('[Auth] Processing auth success:', {
        incoming: action.users,
        current: state.user
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
      console.error('[Auth] Login failed:', action.message);
      return {
        ...state,
        isAuthenticated: false,
        user: null,
        token: null,
        loading: false,
        error: action.message,
      };
    case AuthActionTypes.AUTHEN_LOGIN_LARK_REQ:
      return { ...state, loading: true };
    case AuthActionTypes.AUTHEN_LOGIN_LARK_SUCCESS:
      return {
        ...state,
        loading: false,
        isAuthenticated: true,
        user: action.users,
      };
    case AuthActionTypes.AUTHEN_LOGIN_LARK_FAIL:
      return {
        ...state,
        loading: false,
        isAuthenticated: false,
        error: action.message,
      };
    case AuthActionTypes.AUTHEN_LOGOUT_SUCCESS:
      return { ...state, loading: false, isAuthenticated: false, user: null };
    case AuthActionTypes.AUTHEN_LOGOUT_FAIL:
      return { ...state, loading: false, error: action.message };
    case AuthActionTypes.AUTHEN_CHECK_REQ:
      return { 
        ...state, 
        loading: true,
        error: null
      };
    case AuthActionTypes.AUTHEN_CHECK_SUCCESS:
      return {
        ...state,
        isAuthenticated: true,
        user: action.users,
        loading: false,
        error: null
      };
    case AuthActionTypes.AUTHEN_CHECK_FAIL:
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
