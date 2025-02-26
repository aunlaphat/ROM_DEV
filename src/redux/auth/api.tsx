import { put, call, Effect } from "redux-saga/effects";
import { GET, POST } from "../../services";
import { CHECKAUTH, LOGIN } from "../../services/path";
import { AuthActionTypes } from "./types";
import { windowNavigateReplaceTo } from "../../utils/navigation";
import { ROUTES, ROUTES_NO_AUTH, ROUTE_LOGIN } from "../../resources/routes";
import { closeLoading, openAlert, openLoading } from "../../components/alert/useAlert";
import { AxiosResponse } from "axios";
import { getCookies, setCookies } from "../../store/useCookies";
import { logger } from '../../utils/logger';
import { delay } from 'redux-saga/effects';
import { notification } from "antd";

interface LoginPayload {
  username: string;
  password: string;
}

interface LoginResponse {
  success: boolean;
  message: string;
  data: string;
}

interface AuthCheckResponse {
  success: boolean;
  message: string;
  data: {
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

    const response: AxiosResponse<LoginResponse> = yield call(POST, LOGIN, {
      userName: action.payload.username,
      password: action.payload.password
    });

    if (!response.data.success) {
      throw new Error(response.data.message || 'Login failed');
    }

    const token = response.data.data;
    localStorage.setItem("access_token", token);

    const authResponse: AxiosResponse<AuthCheckResponse> = yield call(GET, CHECKAUTH);
    const { user } = authResponse.data.data;

    yield put({ 
      type: AuthActionTypes.AUTHEN_LOGIN_SUCCESS, 
      users: user
    });

    notification.success({ 
      message: "Login Successful! 🎉",
      description: `ยินดีต้อนรับ คุณ ${user.fullName} 👋`,
      placement: "topLeft",
    });

    yield delay(1000);
    window.location.href = ROUTES.ROUTE_MAIN.PATH;

  } catch (error: any) {
    logger.auth('error', 'Login failed', { error: error.message });

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

export function* logout(): Generator<Effect, void, AxiosResponse> {
  try {
    openLoading();
    logger.auth("info", "Processing logout");

    notification.success({
      message: "Logged Out Successfully! 👋",
      description: "ขอบคุณที่ทำงานหนัก แล้วพบกันใหม่! 🚀",
      placement: "topLeft",
    });

    yield delay(1000);

    localStorage.removeItem("access_token");
    document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";

    yield put({ type: AuthActionTypes.AUTHEN_LOGOUT_SUCCESS });

    try {
      yield call(POST, "/auth/logout", {});
    } catch (e) {
      logger.auth("warn", "Non-critical logout API error", e);
    }

    window.location.href = ROUTES_NO_AUTH.ROUTE_LOGIN.PATH;

  } catch (error: any) {
    logger.auth("error", "Critical logout error", { error: error.message });

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

export function* checkAuthen(): Generator<Effect, void, AxiosResponse<AuthCheckResponse>> {
  try {
    console.log('🔍 Starting auth check...');
    
    const response = (yield call(GET, CHECKAUTH)) as AxiosResponse<AuthCheckResponse>;
    
    if (!response?.data?.success) {
      throw new Error(response?.data?.message || 'Authentication failed');
    }

    const { user } = response.data.data;
    console.log('✅ Auth check success:', user);

    yield put({
      type: AuthActionTypes.AUTHEN_CHECK_SUCCESS,
      users: user
    });

  } catch (error: any) {
    console.error('❌ Auth check error:', error);

    notification.error({
      message: "Session Expired! ⏳",
      description: "เซสชันหมดอายุ กรุณาเข้าสู่ระบบใหม่ 🔐",
      placement: "topLeft",
        });

    yield put({ 
      type: AuthActionTypes.AUTHEN_CHECK_FAIL, 
      message: error.message 
    });

    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
  } finally {
    closeLoading();
  }
}
