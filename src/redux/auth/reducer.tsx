// src/redux/auth/reducer.tsx
import { AuthActionTypes } from "./types";
import { logger } from '../../utils/logger';
import { AuthState, User } from "./interface";

// สร้าง initial state
const initialState: AuthState = {
  user: null,
  token: localStorage.getItem("access_token") || null,
  isAuthenticated: false,
  loading: false,
  error: null
};

/**
 * ฟังก์ชันสำหรับแปลงข้อมูล user ให้สอดคล้องกับแอพพลิเคชัน
 * รองรับทั้งกรณีที่ backend ส่ง fullName และกรณีที่ frontend ใช้ fullNameTH
 */
const normalizeUser = (userData: any): User => {
  return {
    userID: userData.userID,
    userName: userData.userName,
    fullNameTH: userData.fullNameTH,
    nickName: userData.nickName,
    roleID: userData.roleID || userData.userRoleID,
    roleName: userData.roleName,
    departmentNo: userData.departmentNo,
    platform: userData.platform
  };
};

/**
 * Auth Reducer
 */
export default function authReducer(state: AuthState = initialState, action: any): AuthState {
  switch (action.type) {
    // Request actions
    case AuthActionTypes.AUTHEN_LOGIN_REQ:
    case AuthActionTypes.AUTHEN_LOGIN_LARK_REQ:
    case AuthActionTypes.AUTHEN_CHECK_REQ:
      logger.redux.action(action.type);
      return { 
        ...state, 
        loading: true,
        error: null
      };
      
    // Success login actions (รวมทั้ง login ปกติและ Lark)
    case AuthActionTypes.AUTHEN_LOGIN_SUCCESS:
    case AuthActionTypes.AUTHEN_LOGIN_LARK_SUCCESS:
      logger.redux.action(action.type, { 
        userID: action.users?.userID,
        userName: action.users?.userName
      });
      
      // นำข้อมูล user มาแปลงเป็นรูปแบบที่แอพต้องการ
      const normalizedUser = normalizeUser(action.users);
      
      logger.state.update('Auth State', {
        isAuthenticated: true,
        user: {
          userID: normalizedUser.userID,
          userName: normalizedUser.userName,
          roleID: normalizedUser.roleID
        }
      });
      
      return {
        ...state,
        isAuthenticated: true,
        user: normalizedUser,
        loading: false,
        error: null
      };
      
    // Success check auth
    case AuthActionTypes.AUTHEN_CHECK_SUCCESS:
      logger.redux.action(action.type, {
        userID: action.users?.userID,
        userName: action.users?.userName
      });
      
      // นำข้อมูล user มาแปลงเป็นรูปแบบที่แอพต้องการ
      const checkedUser = normalizeUser(action.users);
      
      logger.state.update('Auth State (Check)', {
        isAuthenticated: true,
        user: {
          userID: checkedUser.userID,
          userName: checkedUser.userName
        }
      });
      
      return {
        ...state,
        isAuthenticated: true,
        user: checkedUser,
        loading: false,
        error: null
      };
      
    // Failure actions
    case AuthActionTypes.AUTHEN_LOGIN_FAIL:
    case AuthActionTypes.AUTHEN_LOGIN_LARK_FAIL:
    case AuthActionTypes.AUTHEN_CHECK_FAIL:
      logger.redux.action(action.type, { error: action.message });
      
      logger.state.error('Auth State Error', { 
        message: action.message,
        previousState: {
          isAuthenticated: state.isAuthenticated,
          hasUser: !!state.user
        }
      });
      
      return {
        ...state,
        isAuthenticated: false,
        user: null,
        token: null,
        loading: false,
        error: action.message
      };
      
    // Logout success
    case AuthActionTypes.AUTHEN_LOGOUT_SUCCESS:
      logger.redux.action(action.type);
      
      logger.state.update('Auth State', {
        isAuthenticated: false,
        message: 'User logged out'
      });
      
      return { 
        ...state, 
        loading: false, 
        isAuthenticated: false, 
        user: null,
        token: null,
        error: null
      };
      
    // Logout failure - แต่ยังคงสถานะเดิม
    case AuthActionTypes.AUTHEN_LOGOUT_FAIL:
      logger.redux.action(action.type, { error: action.message });
      
      logger.state.error('Logout Failed', { message: action.message });
      
      return { 
        ...state, 
        loading: false, 
        error: action.message
      };
      
    // Default case
    default:
      return state;
  }
}