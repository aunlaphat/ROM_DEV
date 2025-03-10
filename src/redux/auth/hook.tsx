// src/redux/auth/hooks.tsx
import { useSelector, useDispatch } from 'react-redux';
import { useCallback, useEffect } from 'react';
import { RootState } from '../types';
import { login, loginLark, logout, checkAuthen } from './action';
import { logger } from '../../utils/logger';
import { LarkLoginRequest, LoginRequest } from './interface';

/**
 * Custom hook สำหรับใช้งาน authentication
 * รวมฟังก์ชันและข้อมูลที่เกี่ยวข้องกับการยืนยันตัวตน
 * 
 * @returns ข้อมูลและฟังก์ชันที่เกี่ยวข้องกับ auth
 */
export const useAuth = () => {
  const dispatch = useDispatch();
  
  // ดึงข้อมูล auth จาก redux store
  const auth = useSelector((state: RootState) => state.auth);
  
  // ทุกครั้งที่มีการเรียกใช้ hook จะตรวจสอบ token
  useEffect(() => {
    const token = localStorage.getItem('access_token');
    
    logger.log('debug', '[useAuth] Hook initialized', {
      hasToken: !!token,
      isAuthenticated: auth.isAuthenticated,
      timestamp: new Date().toISOString()
    });
    
    // ถ้ามี token แต่ยังไม่ได้ยืนยันตัวตน จะทำการตรวจสอบอัตโนมัติ
    if (token && !auth.isAuthenticated) {
      logger.log('info', '[useAuth] Auto checking authentication');
      dispatch(checkAuthen());
    }
  }, [dispatch, auth.isAuthenticated]);

    // ฟังก์ชันสำหรับล็อกอิน
  const handleLogin = useCallback((credentials: LoginRequest) => {
    logger.log('info', '[useAuth] Login requested', {
      userName: credentials.userName,
      timestamp: new Date().toISOString()
    });
    dispatch(login(credentials));
  }, [dispatch]);
  
  // ฟังก์ชันสำหรับล็อกอินผ่าน Lark
  const handleLarkLogin = useCallback((larkData: LarkLoginRequest) => {
    logger.log('info', '[useAuth] Lark login requested', {
      userID: larkData.userID,
      timestamp: new Date().toISOString()
    });
    dispatch(loginLark(larkData));
  }, [dispatch]);
  
  // ฟังก์ชันสำหรับล็อกเอาท์
  const handleLogout = useCallback(() => {
    logger.log('info', '[useAuth] Logout requested', {
      timestamp: new Date().toISOString()
    });
    dispatch(logout());
  }, [dispatch]);
  
  // ฟังก์ชันสำหรับตรวจสอบการยืนยันตัวตน
  const handleCheckAuth = useCallback(() => {
    logger.log('info', '[useAuth] Manual auth check requested', {
      timestamp: new Date().toISOString()
    });
    dispatch(checkAuthen());
  }, [dispatch]);
  
  return {
    // ข้อมูล auth
    isAuthenticated: auth.isAuthenticated,
    user: auth.user,
    loading: auth.loading,
    error: auth.error,
    
    // ฟังก์ชันการทำงาน
    login: handleLogin,
    loginLark: handleLarkLogin,
    logout: handleLogout,
    checkAuth: handleCheckAuth
  };
};

/**
 * Custom hook สำหรับตรวจสอบบทบาทของผู้ใช้
 * 
 * @param requiredRoleIDs บทบาทที่ต้องการตรวจสอบ (array)
 * @returns true ถ้าผู้ใช้มีบทบาทที่ต้องการ
 */
export const useHasRole = (requiredRoleIDs: number[]) => {
  const auth = useSelector((state: RootState) => state.auth);
  
  // ถ้าไม่ได้ยืนยันตัวตน หรือไม่มีข้อมูลผู้ใช้ จะไม่มีบทบาทใดๆ
  if (!auth.isAuthenticated || !auth.user) {
    return false;
  }
  
  // ตรวจสอบว่าผู้ใช้มีบทบาทที่ต้องการหรือไม่
  return requiredRoleIDs.includes(auth.user.roleID);
};