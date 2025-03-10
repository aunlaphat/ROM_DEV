// src/redux/auth/saga.tsx
import { takeLatest, all } from "redux-saga/effects";
import { login, loginLark, logout, checkAuthen } from "./api";
import { AuthActionTypes } from "./types";
import { logger } from '../../utils/logger';

/**
 * Auth Saga - รวมทุก watcher สำหรับการจัดการ authentication
 */
function* authSaga() {
  // ใช้ logger เพื่อบันทึกการเริ่มต้น saga
  logger.log('info', '[Auth Saga] Initialized auth watchers');
  
  yield all([
    // Login watcher - สำหรับการล็อกอินปกติ
    takeLatest(AuthActionTypes.AUTHEN_LOGIN_REQ, login),
    
    // Lark Login watcher - สำหรับการล็อกอินผ่าน Lark
    takeLatest(AuthActionTypes.AUTHEN_LOGIN_LARK_REQ, loginLark),
    
    // Logout watcher - สำหรับการล็อกเอาท์
    takeLatest(AuthActionTypes.AUTHEN_LOGOUT_REQ, logout),
    
    // Auth Check watcher - สำหรับตรวจสอบการยืนยันตัวตน
    takeLatest(AuthActionTypes.AUTHEN_CHECK_REQ, checkAuthen)
  ]);
}

export default authSaga;