// src/redux/auth/api.tsx
import { put, call, Effect } from "redux-saga/effects";
import { GET, POST } from "../../services";
import { AUTH } from "../../services/path";
import { windowNavigateReplaceTo } from "../../utils/navigation";
import { ROUTES, ROUTES_NO_AUTH, ROUTE_LOGIN } from "../../resources/routes";
import { closeLoading, openLoading } from "../../components/alert/useAlert";
import { AxiosResponse } from "axios";
import { setCookies, removeCookies } from "../../store/useCookies";
import { logger } from '../../utils/logger';
import { delay } from 'redux-saga/effects';
import { notification } from "antd";
import { AuthActionTypes } from "./types";
import { ApiResponse, AuthCheckResponse, LarkLoginRequest, LoginRequest } from "./interface";

/**
 * จัดการ Login ปกติ
 */
export function* login(action: { type: AuthActionTypes; payload: LoginRequest }): Generator<Effect, void, any> {
  try {
    openLoading();
    logger.perf.start('Auth: Login Process');
    
    logger.log('info', '[Auth] Starting login process', {
      username: action.payload.userName,
      timestamp: new Date().toISOString()
    });

    // เรียก API Login
    logger.api.request(AUTH.LOGIN, {
      userName: action.payload.userName,
      // ไม่เก็บ password ใน log
      password: '********'
    });

    const response: AxiosResponse<ApiResponse<string>> = yield call(POST, AUTH.LOGIN, action.payload);

    if (!response.data.success) {
      throw new Error(response.data.message || 'Login failed');
    }

    logger.api.success(AUTH.LOGIN, {
      success: response.data.success,
      hasToken: !!response.data.data
    });

    // เก็บ token ที่ได้จาก backend
    const token = response.data.data;
    localStorage.setItem("access_token", token);
    setCookies("jwt", token);

    logger.log('debug', '[Auth] Token stored', {
      storage: 'localStorage + cookies',
      tokenLength: token.length
    });

    // ตรวจสอบข้อมูลผู้ใช้ด้วย Token ที่ได้มา
    logger.api.request(AUTH.CHECK);
    const authResponse: AxiosResponse<ApiResponse<AuthCheckResponse>> = yield call(GET, AUTH.CHECK);
    
    if (!authResponse.data.success) {
      throw new Error(authResponse.data.message || 'Failed to get user data');
    }
    
    logger.api.success(AUTH.CHECK, {
      success: authResponse.data.success,
      user: authResponse.data.data.user
    });

    // แปลงข้อมูลผู้ใช้ให้ตรงกับที่ frontend ต้องการ
    const { user } = authResponse.data.data;
    
    // ตรวจสอบและบันทึก response จาก API เพื่อดูข้อมูลที่ส่งกลับมา
    logger.log('debug', '[Auth] User data received', {
      rawUserData: user,
      hasFullNameTH: !!user.fullNameTH
    });
    
    // แก้ไข: ตรวจสอบข้อมูลชื่อผู้ใช้และจัดการกรณีที่อาจเป็น undefined
    const normalizedUser = {
      userID: user.userID || '',
      userName: user.userName || '',
      fullNameTH: user.fullNameTH || '',
      nickName: user.nickName || '',
      roleID: user.roleID || 0,
      roleName: user.roleName || '',
      departmentNo: user.departmentNo || '',
      platform: user.platform || ''
    };
    
    // บันทึก log ข้อมูลที่ normalize แล้ว
    logger.log('debug', '[Auth] Normalized user data', {
      normalizedUser: {
        userID: normalizedUser.userID,
        userName: normalizedUser.userName,
        fullNameTH: normalizedUser.fullNameTH,
        roleID: normalizedUser.roleID,
        roleName: normalizedUser.roleName
      }
    });

    // อัปเดต Redux state
    logger.redux.action(AuthActionTypes.AUTHEN_LOGIN_SUCCESS, { userID: normalizedUser.userID });
    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGIN_SUCCESS, 
      users: normalizedUser
    });
    
    // สร้างข้อความต้อนรับที่ปลอดภัย
    const welcomeName = normalizedUser.fullNameTH || normalizedUser.userName || 'ผู้ใช้งาน';

    // แสดงข้อความแจ้งเตือน
    notification.success({ 
      message: "Login Successful! 🎉",
      description: `ยินดีต้อนรับ คุณ ${welcomeName} 👋`,
      placement: "topLeft",
    });

    yield delay(1000);
    
    // นำทางไปหน้าหลัก
    logger.navigation.to(ROUTES.ROUTE_MAIN.PATH, { 
      method: 'window.location', 
      from: 'login saga'
    });
    window.location.href = ROUTES.ROUTE_MAIN.PATH;

  } catch (error: any) {
    logger.auth.login(action.payload.userName, false, { 
      error: error.message,
      response: error.response?.data,
      timestamp: new Date().toISOString()
    });

    logger.redux.action(AuthActionTypes.AUTHEN_LOGIN_FAIL, { 
      error: error.message 
    });

    // แสดงข้อความแจ้งเตือนข้อผิดพลาด
    notification.error({ 
      message: "Login Failed! ❌",
      description: error.response?.data?.message || "เข้าสู่ระบบไม่สำเร็จ กรุณาตรวจสอบข้อมูลและลองใหม่ 🔄",
      placement: "topLeft",
    });

    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGIN_FAIL, 
      message: error.message 
    });
  } finally {
    yield delay(300);
    closeLoading();
    logger.perf.end('Auth: Login Process');
  }
}

/**
 * จัดการ Login ผ่าน Lark
 */
export function* loginLark(action: { type: AuthActionTypes; payload: LarkLoginRequest }): Generator<Effect, void, any> {
  try {
    openLoading();
    logger.perf.start('Auth: Lark Login Process');
    
    logger.log('info', '[Auth] Starting Lark login process', {
      userID: action.payload.userID,
      userName: action.payload.userName,
      timestamp: new Date().toISOString()
    });

    // เรียก API Login Lark
    logger.api.request(AUTH.LOGIN_LARK, {
      userID: action.payload.userID,
      userName: action.payload.userName,
    });

    const response: AxiosResponse<ApiResponse<string>> = yield call(POST, AUTH.LOGIN_LARK, action.payload);

    if (!response.data.success) {
      throw new Error(response.data.message || 'Lark login failed');
    }

    logger.api.success(AUTH.LOGIN_LARK, {
      success: response.data.success,
      hasToken: !!response.data.data
    });

    // เก็บ token ที่ได้จาก backend
    const token = response.data.data;
    localStorage.setItem("access_token", token);
    setCookies("jwt", token);

    logger.log('debug', '[Auth] Token stored from Lark login', {
      storage: 'localStorage + cookies',
      tokenLength: token.length
    });

    // ตรวจสอบข้อมูลผู้ใช้ด้วย Token ที่ได้มา
    logger.api.request(AUTH.CHECK);
    const authResponse: AxiosResponse<ApiResponse<AuthCheckResponse>> = yield call(GET, AUTH.CHECK);
    
    if (!authResponse.data.success) {
      throw new Error(authResponse.data.message || 'Failed to get user data after Lark login');
    }
    
    // แปลงข้อมูลผู้ใช้ให้ตรงกับที่ frontend ต้องการ
    const { user } = authResponse.data.data;
    
    // ตรวจสอบและบันทึก response จาก API เพื่อดูข้อมูลที่ส่งกลับมา
    logger.log('debug', '[Auth] Lark user data received', {
      rawUserData: user,
      hasFullNameTH: !!user.fullNameTH
    });
    
    // แก้ไข: ตรวจสอบข้อมูลชื่อผู้ใช้และจัดการกรณีที่อาจเป็น undefined
    const normalizedUser = {
      userID: user.userID || '',
      userName: user.userName || '',
      fullNameTH: user.fullNameTH || '',
      nickName: user.nickName || '',
      roleID: user.roleID || 0,
      roleName: user.roleName || '',
      departmentNo: user.departmentNo || '',
      platform: user.platform || ''
    };
    
    // บันทึก log ข้อมูลที่ normalize แล้ว
    logger.log('debug', '[Auth] Normalized Lark user data', {
      normalizedUser: {
        userID: normalizedUser.userID,
        userName: normalizedUser.userName,
        fullNameTH: normalizedUser.fullNameTH,
        roleID: normalizedUser.roleID
      }
    });

    // อัปเดต Redux state
    logger.redux.action(AuthActionTypes.AUTHEN_LOGIN_LARK_SUCCESS, { userID: normalizedUser.userID });
    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGIN_LARK_SUCCESS, 
      users: normalizedUser
    });
    
    // สร้างข้อความต้อนรับที่ปลอดภัย
    const welcomeName = normalizedUser.fullNameTH || normalizedUser.userName || 'ผู้ใช้งาน';

    // แสดงข้อความแจ้งเตือน
    notification.success({ 
      message: "Lark Login Successful! 🎉",
      description: `ยินดีต้อนรับ คุณ ${welcomeName} 👋`,
      placement: "topLeft",
    });

    yield delay(1000);
    
    // นำทางไปหน้าหลัก
    window.location.href = ROUTES.ROUTE_MAIN.PATH;

  } catch (error: any) {
    logger.error('[Auth] Lark login failed', {
      error: error.message,
      userID: action.payload.userID,
      timestamp: new Date().toISOString(),
      response: error.response?.data
    });

    // แสดงข้อความแจ้งเตือนข้อผิดพลาด
    notification.error({ 
      message: "Lark Login Failed! ❌",
      description: error.response?.data?.message || "เข้าสู่ระบบผ่าน Lark ไม่สำเร็จ กรุณาลองใหม่อีกครั้ง 🔄",
      placement: "topLeft",
    });

    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGIN_LARK_FAIL, 
      message: error.message 
    });
  } finally {
    yield delay(300);
    closeLoading();
    logger.perf.end('Auth: Lark Login Process');
  }
}

/**
 * จัดการ Logout
 */
export function* logout(): Generator<Effect, void, any> {
  try {
    openLoading();
    logger.perf.start('Auth: Logout Process');
    logger.log("info", '[Auth] Processing logout', {
      timestamp: new Date().toISOString()
    });

    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: "Logged Out Successfully! 👋",
      description: "ขอบคุณที่ทำงานหนัก แล้วพบกันใหม่! 🚀",
      placement: "topLeft",
    });

    yield delay(1000);

    // ลบ token
    localStorage.removeItem("access_token");
    removeCookies("jwt");
    
    logger.log('debug', '[Auth] Tokens cleared', {
      storage: 'localStorage + cookies'
    });

    // เรียก logout API 
    try {
      logger.api.request(AUTH.LOGOUT);
      yield call(POST, AUTH.LOGOUT, {});
      logger.api.success(AUTH.LOGOUT);
    } catch (e) {
      logger.log("warn", '[Auth] Non-critical logout API error', e);
    }

    // Update Redux state
    logger.redux.action(AuthActionTypes.AUTHEN_LOGOUT_SUCCESS);
    yield put({ type: AuthActionTypes.AUTHEN_LOGOUT_SUCCESS });

    // นำทางไปหน้า login
    logger.navigation.to(ROUTES_NO_AUTH.ROUTE_LOGIN.PATH, {
      method: 'window.location',
      from: 'logout saga'
    });
    window.location.href = ROUTES_NO_AUTH.ROUTE_LOGIN.PATH;

  } catch (error: any) {
    logger.error('[Auth] Critical logout error', error);

    // แสดงข้อความแจ้งเตือนข้อผิดพลาด
    notification.error({
      message: "Logout Failed! ❌",
      description: "เกิดข้อผิดพลาด ไม่สามารถออกจากระบบได้ กรุณาลองใหม่อีกครั้ง 🔄",
      placement: "topLeft",
    });

    // Update Redux state with error
    logger.redux.action(AuthActionTypes.AUTHEN_LOGOUT_FAIL, { error: error.message });
    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGOUT_FAIL, 
      message: error.message 
    });

    yield delay(1500);
  } finally {
    closeLoading();
    logger.perf.end('Auth: Logout Process');
  }
}

/**
 * ตรวจสอบการยืนยันตัวตน
 */
export function* checkAuthen(): Generator<Effect, void, any> {
  try {
    logger.perf.start('Auth: Check Authentication');
    logger.log('info', '[Auth] Starting auth check...', {
      timestamp: new Date().toISOString()
    });
    
    // เรียก API
    logger.api.request(AUTH.CHECK);
    const response: AxiosResponse<ApiResponse<AuthCheckResponse>> = yield call(GET, AUTH.CHECK);
    
    if (!response?.data?.success) {
      throw new Error(response?.data?.message || 'Authentication failed');
    }

    // API success
    logger.api.success(AUTH.CHECK, {
      user: response.data.data.user
    });

    // แก้ไข: ตรวจสอบว่า user มีจริงหรือไม่
    const { user } = response.data.data;
    
    if (!user) {
      throw new Error('No user data returned from API');
    }

    // Log authenticated user
    logger.auth.check(true, {
      user: {
        userID: user.userID || '',
        userName: user.userName || '',
        roleID: user.roleID || 0,
        roleName: user.roleName || ''
      },
      timestamp: new Date().toISOString()
    });
    
    // แก้ไข: เพิ่มการ normalize ข้อมูลผู้ใช้
    const normalizedUser = {
      userID: user.userID || '',
      userName: user.userName || '',
      fullNameTH: user.fullNameTH || '',
      nickName: user.nickName || '',
      roleID: user.roleID || 0,
      roleName: user.roleName || '',
      departmentNo: user.departmentNo || '',
      platform: user.platform || ''
    };

    // Update Redux state
    logger.redux.action(AuthActionTypes.AUTHEN_CHECK_SUCCESS, { userID: normalizedUser.userID });
    yield put({
      type: AuthActionTypes.AUTHEN_CHECK_SUCCESS,
      users: normalizedUser
    });

  } catch (error: any) {
    logger.auth.check(false, { 
      error: error.message,
      timestamp: new Date().toISOString()
    });
    
    logger.error('[Auth] Auth check error', error);

    // แสดงข้อความแจ้งเตือนข้อผิดพลาด
    notification.error({
      message: "Session Expired! ⏳",
      description: "เซสชันหมดอายุ กรุณาเข้าสู่ระบบใหม่ 🔐",
      placement: "topLeft",
    });

    // Cleanup tokens ในกรณีตรวจสอบล้มเหลว
    localStorage.removeItem("access_token");
    removeCookies("jwt");

    // Update Redux state with error
    logger.redux.action(AuthActionTypes.AUTHEN_CHECK_FAIL, { error: error.message });
    yield put({ 
      type: AuthActionTypes.AUTHEN_CHECK_FAIL, 
      message: error.message 
    });

    // นำทางไปหน้า login
    logger.navigation.to(ROUTE_LOGIN, {
      method: 'windowNavigateReplaceTo',
      from: 'auth check failure'
    });
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
  } finally {
    closeLoading();
    logger.perf.end('Auth: Check Authentication');
  }
}