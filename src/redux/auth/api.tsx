import { put, call, Effect } from "redux-saga/effects";
import { GET, POST } from "../../services";
import { CHECKAUTH, LOGIN } from "../../services/path";
import { AuthActionTypes } from "./types";
import { windowNavigateReplaceTo } from "../../utils/navigation";
import { ROUTES_PATH, ROUTE_LOGIN } from "../../resources/routes";
import { closeLoading, openAlert, openLoading } from "../../components/alert/useAlert";
import { AxiosResponse } from "axios";
import { getCookies, setCookies } from "../../store/useCookies";
import { logger } from '../../utils/logger';
import { delay } from 'redux-saga/effects';

// ✅ กำหนด Type ของ Payload และ Response
interface LoginPayload {
  username: string;
  password: string;
}

interface LoginResponse {
  success: boolean;
  message: string;
  data: string;  // JWT token from backend
  statusCode?: number;
}

interface AuthCheckResponse {
  success: boolean;
  message: string;
  data: {
    source: string;
    user: {
      userID: string;
      userName: string;
      fullName: string;
      nickName: string;
      roleID: number;
      roleName: string;
      departmentNo: string;
      platform: string;
    };
  };
}

export function* login(action: { type: AuthActionTypes; payload: LoginPayload }) {
  try {
    openLoading();
    logger.auth('info', 'Starting login process');

    // 1. Login และรับ token
    const response: AxiosResponse<LoginResponse> = yield call(POST, LOGIN, {
      userName: action.payload.username,
      password: action.payload.password
    });

    if (!response.data.success) {
      throw new Error(response.data.message || 'Login failed');
    }

    // 2. บันทึก token
    const token = response.data.data;
    localStorage.setItem("access_token", token);

    // 3. ดึงข้อมูล user
    const authResponse: AxiosResponse<AuthCheckResponse> = yield call(GET, CHECKAUTH);
    const { user } = authResponse.data.data;

    // 4. อัพเดท Redux state ก่อน redirect
    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGIN_SUCCESS, 
      users: {
        userID: user.userID,
        userName: user.userName,
        roleID: user.roleID,
        fullName: user.fullName,
        nickName: user.nickName,
        roleName: user.roleName,
        departmentNo: user.departmentNo,
        platform: user.platform
      }
    });

    // 5. แสดง success message
    openAlert({ type: "success", title: "เข้าสู่ระบบสำเร็จ" });

    // 6. รอให้ state อัพเดทเสร็จก่อน redirect
    yield delay(100);
    yield call(windowNavigateReplaceTo, { pathname: "/home" });

  } catch (error: any) {
    logger.auth('error', 'Login failed', {
      error: error.message,
      details: error.response?.data
    });
    console.error('[Auth] Login error details:', {
      request: {
        username: action.payload.username,
        url: LOGIN
      },
      error: {
        message: error.message,
        status: error.response?.status,
        data: error.response?.data
      }
    });

    const errorMessage = error.response?.data?.message || error.message;
    openAlert({ 
      type: "error", 
      title: "เข้าสู่ระบบไม่สำเร็จ", 
      message: errorMessage 
    });
    
    yield put({ type: AuthActionTypes.AUTHEN_LOGIN_FAIL, message: errorMessage });
  } finally {
    closeLoading();
  }
}

// ✅ ปรับปรุงฟังก์ชัน Logout
export function* logout(): Generator<Effect, void, AxiosResponse> {
  try {
    openLoading();
    logger.auth('info', 'Processing logout');

    // 1. Clear all auth data first
    localStorage.removeItem("access_token");
    document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    
    // 2. Update redux state
    yield put({ type: AuthActionTypes.AUTHEN_LOGOUT_SUCCESS });
    
    // 3. Call logout API (ถ้ามี error ก็ไม่เป็นไร เพราะเคลียร์ข้อมูลไปแล้ว)
    yield call(POST, "/auth/logout", {});

    logger.auth('info', 'Logout successful');
    openAlert({ type: "success", title: "ออกจากระบบสำเร็จ" });
    
    // 4. Redirect to login
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });

  } catch (error: any) {
    logger.auth('error', 'Logout error', { error: error.message });
    openAlert({ type: "error", message: "เกิดข้อผิดพลาดระหว่างออกจากระบบ" });
  } finally {
    closeLoading();
  }
}

// ✅ ปรับปรุงฟังก์ชัน Check Authentication
export function* checkAuthen(): Generator<Effect, void, AxiosResponse<AuthCheckResponse>> {
  try {
    console.log('🔍 Starting auth check...');
    
    const response = (yield call(GET, CHECKAUTH)) as AxiosResponse<AuthCheckResponse>;
    
    if (!response?.data?.success) {
      console.error('❌ Auth check failed:', response?.data?.message);
      throw new Error(response?.data?.message || 'Authentication failed');
    }

    const { user } = response.data.data;
    console.log('✅ Auth check success:', user);

    yield put({
      type: AuthActionTypes.AUTHEN_CHECK_SUCCESS,
      users: {
        userID: user.userID,
        userName: user.userName,
        roleID: user.roleID,
        fullName: user.fullName,
        nickName: user.nickName,
        roleName: user.roleName,
        departmentNo: user.departmentNo,
        platform: user.platform
      }
    });

  } catch (error: any) {
    console.error('❌ Auth check error:', error);
    yield put({ 
      type: AuthActionTypes.AUTHEN_CHECK_FAIL, 
      message: error.message 
    });
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
  } finally {
    closeLoading();
  }
}
