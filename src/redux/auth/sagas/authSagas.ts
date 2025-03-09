import { put, call, Effect, delay } from 'redux-saga/effects';
import { AxiosResponse } from 'axios';
import { notification } from 'antd';
import { AuthActionTypes, LoginPayload, LoginResponse, AuthCheckResponse } from '../../../types/auth.types';
import { authAPI } from '../../../services/api/auth.api';
import { ROUTES, ROUTES_NO_AUTH } from '../../../resources/routes';
import { windowNavigateReplaceTo } from '../../../utils/navigation';
import { ROUTE_LOGIN } from '../../../resources/routes';
import { closeLoading, openLoading } from '../../../components/alert/useAlert';
import { logger } from '../../../utils/logger';

export function* loginSaga(action: { type: AuthActionTypes; payload: LoginPayload }) {
  try {
    openLoading();
    logger.log('info', '[Auth] Starting login process');

    const response: AxiosResponse<LoginResponse> = yield call(
      authAPI.login, 
      action.payload
    );

    if (!response.data.success) {
      throw new Error(response.data.message || 'Login failed');
    }

    // เก็บ token ใน localStorage เป็น backup (backend จัดการ cookie ให้แล้ว)
    const token = response.data.data;
    localStorage.setItem("access_token", token);

    // ตรวจสอบข้อมูลผู้ใช้หลังจาก login สำเร็จ
    const authResponse: AxiosResponse<AuthCheckResponse> = yield call(authAPI.checkAuth);
    const { user } = authResponse.data.data;

    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGIN_SUCCESS, 
      users: user
    });

    notification.success({ 
      message: "Login Successful! 🎉",
      description: `ยินดีต้อนรับ คุณ ${user.fullName || user.fullNameTH || user.userName} 👋`,
      placement: "topLeft",
    });

    yield delay(1000);
    window.location.href = ROUTES.ROUTE_MAIN.PATH;

  } catch (error: any) {
    logger.log('error', '[Auth] Login failed', { error: error.message });

    notification.error({ 
      message: "Login Failed! ❌",
      description: error.response?.data?.message || "เข้าสู่ระบบไม่สำเร็จ กรุณาตรวจสอบข้อมูลและลองใหม่ 🔄",
      placement: "topLeft",
    });

    yield put({ type: AuthActionTypes.AUTHEN_LOGIN_FAIL, message: error.message });
  } finally {
    yield delay(300);
    closeLoading();
  }
}

export function* loginLarkSaga(action: { type: AuthActionTypes; payload: any }) {
  try {
    openLoading();
    logger.log('info', '[Auth] Starting Lark login process');

    const response: AxiosResponse<LoginResponse> = yield call(
      authAPI.loginLark, 
      action.payload
    );

    if (!response.data.success) {
      throw new Error(response.data.message || 'Lark login failed');
    }

    const token = response.data.data;
    localStorage.setItem("access_token", token);

    const authResponse: AxiosResponse<AuthCheckResponse> = yield call(authAPI.checkAuth);
    const { user } = authResponse.data.data;

    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGIN_LARK_SUCCESS, 
      users: user
    });

    notification.success({ 
      message: "Lark Login Successful! 🎉",
      description: `ยินดีต้อนรับ คุณ ${user.fullName || user.fullNameTH || user.userName} 👋`,
      placement: "topLeft",
    });

    yield delay(1000);
    window.location.href = ROUTES.ROUTE_MAIN.PATH;

  } catch (error: any) {
    logger.log('error', '[Auth] Lark login failed', { error: error.message });

    notification.error({ 
      message: "Lark Login Failed! ❌",
      description: error.response?.data?.message || "เข้าสู่ระบบไม่สำเร็จ กรุณาตรวจสอบข้อมูลและลองใหม่ 🔄",
      placement: "topLeft",
    });

    yield put({ type: AuthActionTypes.AUTHEN_LOGIN_LARK_FAIL, message: error.message });
  } finally {
    yield delay(300);
    closeLoading();
  }
}

export function* logoutSaga(): Generator<Effect, void, AxiosResponse> {
  try {
    openLoading();
    logger.log("info", '[Auth] Processing logout');

    // ส่ง request ไปยัง API เพื่อ invalidate JWT token บน server
    try {
      yield call(authAPI.logout);
    } catch (e) {
      logger.log("warn", '[Auth] Non-critical logout API error', e);
      // ทำการ logout ต่อไปแม้ว่า API จะล้มเหลว เพราะว่า backend จะลบ cookie แล้ว
    }

    // ลบข้อมูลการ login ฝั่ง client
    localStorage.removeItem("access_token");
    // ลบ cookie ชื่อ jwt ถ้ามี (แต่ส่วนใหญ่ backend จะจัดการให้แล้ว)
    document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

    yield put({ type: AuthActionTypes.AUTHEN_LOGOUT_SUCCESS });

    notification.success({
      message: "Logged Out Successfully! 👋",
      description: "ขอบคุณที่ทำงานหนัก แล้วพบกันใหม่! 🚀",
      placement: "topLeft",
    });

    yield delay(1000);
    window.location.href = ROUTES_NO_AUTH.ROUTE_LOGIN.PATH;

  } catch (error: any) {
    logger.log("error", '[Auth] Critical logout error', { error: error.message });

    notification.error({
      message: "Logout Failed! ❌",
      description: "เกิดข้อผิดพลาด ไม่สามารถออกจากระบบได้ กรุณาลองใหม่อีกครั้ง 🔄",
      placement: "topLeft",
    });

    yield delay(1500);
  } finally {
    closeLoading();
  }
}

export function* checkAuthSaga(): Generator<Effect, void, AxiosResponse<AuthCheckResponse>> {
  try {
    logger.log('info', '[Auth] Starting auth check...');
    
    const response = (yield call(authAPI.checkAuth)) as AxiosResponse<AuthCheckResponse>;
    
    if (!response?.data?.success) {
      throw new Error(response?.data?.message || 'Authentication failed');
    }

    const { user } = response.data.data;
    logger.log('info', '[Auth] Auth check success', { user });

    yield put({
      type: AuthActionTypes.AUTHEN_CHECK_SUCCESS,
      users: user
    });

  } catch (error: any) {
    logger.log('error', '[Auth] Auth check error', { error: error.message });

    // เช็คว่ามีการแสดงข้อความแล้วหรือยัง เพื่อไม่ให้แสดงซ้ำ
    const errorStatus = error.response?.status;
    
    // ถ้าเป็น 401 Unauthorized หรือ 403 Forbidden แสดงว่า token หมดอายุหรือไม่ถูกต้อง
    if (errorStatus === 401 || errorStatus === 403) {
      notification.error({
        message: "Session Expired! ⏳",
        description: "เซสชันหมดอายุ กรุณาเข้าสู่ระบบใหม่ 🔐",
        placement: "topLeft",
      });
      
      // ลบข้อมูลการ login ทั้งหมด
      localStorage.removeItem("access_token");
      document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

      yield put({ 
        type: AuthActionTypes.AUTHEN_CHECK_FAIL, 
        message: error.message 
      });

      // นำทางกลับไปหน้า login
      yield delay(1000);
      windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
    } else {
      // กรณีมี error อื่นๆ เช่น server error
      notification.error({
        message: "Authentication Error",
        description: "เกิดข้อผิดพลาดในการตรวจสอบสถานะผู้ใช้ กรุณาลองใหม่อีกครั้ง",
        placement: "topLeft",
      });

      yield put({ 
        type: AuthActionTypes.AUTHEN_CHECK_FAIL, 
        message: error.message 
      });
    }
  } finally {
    closeLoading();
  }
}