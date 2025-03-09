// src/redux/auth/action.tsx
import { logger } from '../../utils/logger';
import { LarkLoginRequest, LoginRequest } from "./interface";
import { AuthActionTypes } from './types';

/**
 * Action creator สำหรับการล็อกอินปกติ
 * @param payload ข้อมูลการล็อกอิน (userName, password)
 */
export const login = (payload: LoginRequest) => {
  logger.auth.login(payload.userName, false, { 
    status: 'attempting',
    timestamp: new Date().toISOString()
  });
  
  return {
    type: AuthActionTypes.AUTHEN_LOGIN_REQ,
    payload,
  };
};

/**
 * Action creator สำหรับการล็อกอินผ่าน Lark
 * @param payload ข้อมูลจาก Lark (userID, userName, fullName, email)
 */
export const loginLark = (payload: LarkLoginRequest) => {
  logger.auth.login(payload.userName || payload.userID, false, { 
    status: 'attempting',
    method: 'Lark',
    provider: 'Lark',
    timestamp: new Date().toISOString(),
    userID: payload.userID
  });
  
  return {
    type: AuthActionTypes.AUTHEN_LOGIN_LARK_REQ,
    payload,
  };
};

/**
 * Action creator สำหรับการล็อกเอาท์
 */
export const logout = () => {
  logger.auth.logout('User initiated logout');
  
  return {
    type: AuthActionTypes.AUTHEN_LOGOUT_REQ,
  };
};

/**
 * Action creator สำหรับตรวจสอบการยืนยันตัวตน
 */
export const checkAuthen = () => {
  logger.auth.check(false, { 
    status: 'checking',
    timestamp: new Date().toISOString()
  });
  
  return {
    type: AuthActionTypes.AUTHEN_CHECK_REQ,
  };
};