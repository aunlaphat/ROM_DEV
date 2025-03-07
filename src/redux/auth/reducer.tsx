import { AuthActionTypes, AuthState } from '../../types/auth.types';
import { logger } from '../../utils/logger';

const initialState: AuthState = {
  user: null,
  token: localStorage.getItem("access_token") || null,
  isAuthenticated: false,
  loading: false,
};

export default function authReducer(state = initialState, action: any): AuthState {
  switch (action.type) {
    case AuthActionTypes.AUTHEN_LOGIN_REQ:
    case AuthActionTypes.AUTHEN_LOGIN_LARK_REQ:
    case AuthActionTypes.AUTHEN_CHECK_REQ:
      logger.log('info', `[Auth] Processing ${action.type}`);
      return { ...state, loading: true, error: null };
      
    case AuthActionTypes.AUTHEN_LOGIN_SUCCESS:
    case AuthActionTypes.AUTHEN_LOGIN_LARK_SUCCESS:
      logger.log('info', `[Auth] Success: ${action.type}`, {
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
        error: null
      };
      
    case AuthActionTypes.AUTHEN_LOGIN_FAIL:
    case AuthActionTypes.AUTHEN_LOGIN_LARK_FAIL:
    case AuthActionTypes.AUTHEN_CHECK_FAIL:
      logger.error(`[Auth] Failure: ${action.type}`, { message: action.message });
      return {
        ...state,
        isAuthenticated: false,
        user: null,
        loading: false,
        error: action.message,
      };
      
    case AuthActionTypes.AUTHEN_LOGOUT_SUCCESS:
      logger.log('info', '[Auth] Logout successful');
      return { 
        ...state, 
        loading: false, 
        isAuthenticated: false, 
        user: null,
        token: null,
        error: null
      };
      
    case AuthActionTypes.AUTHEN_LOGOUT_FAIL:
      logger.error('[Auth] Logout failed', { message: action.message });
      return { ...state, loading: false, error: action.message };
      
    default:
      return state;
  }
}