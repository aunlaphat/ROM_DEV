import { put } from "redux-saga/effects";
import { GET, POST } from "../../services";
import { CHECKAUTH, LARK_LOGIN, LOGIN } from "../../services/path";
import * as type from "./types";
import { windowNavigateReplaceTo } from "../../utils/navigation";
import { ROUTES_PATH, ROUTE_LOGIN } from "../../resources/routes";
import { removeCookies } from "../../store/useCookies";
import {
  closeLoading,
  openAlert,
  openLoading,
} from "../../components/alert/useAlert";

export function* login(payload: any): Generator<any, void, any> {
  try {
    openLoading();
    const users = yield POST(LOGIN, payload.payload);
    openAlert({ type: "success", message: "", title: "เข้าสู่ระบบสำเร็จ" });
    yield put({ type: type.AUTHEN_LOGIN_SUCCESS, users });
    windowNavigateReplaceTo({ pathname: ROUTES_PATH.ROUTE_MAIN.PATH });
  } catch (e: any) {
    openAlert({
      type: "error",
      message: "กรุณาลองใหม่อีกครั้ง",
      title: "เข้าสู่ระบบไม่สำเร็จ",
    });
    yield put({ type: type.AUTHEN_LOGIN_FAIL, message: e.message });
  } finally {
    closeLoading();
  }
}

export function* logout(): Generator<any, void, any> {
  try {
    openAlert({ type: "info", title: "ออกจากระบบสำเร็จ" });
    openLoading();

    // ✅ ใช้ fetch() แทน axios เพื่อให้มั่นใจว่า Cookie ถูกส่งไป
    const response = yield fetch(`${process.env.REACT_APP_BACKEND_URL}/auth/logout`, {
      method: "POST",
      credentials: "include", // ✅ ให้ JWT Cookie ถูกส่งไปกับ API
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      console.warn("🚨 Logout failed with status:", response.status);
    }

    yield put({ type: type.AUTHEN_LOGOUT_SUCCESS });

    // ✅ Redirect ไปหน้า Login หลัง Logout สำเร็จ
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });

  } catch (e: any) {
    openAlert({ type: "error", message: "เกิดข้อผิดพลาดระหว่าง Logout" });
    yield put({ type: type.AUTHEN_LOGOUT_FAIL, message: e.message });
  } finally {
    closeLoading();
  }
}

export function* login_lark(payload: any): Generator<any, void, any> {
  try {
    openLoading();
    const users = yield POST(LARK_LOGIN, payload.payload);
    openAlert({ type: "success", message: "", title: "เข้าสู่ระบบสำเร็จ" });
    yield put({ type: type.AUTHEN_LOGIN_LARK_SUCCESS, users });
    windowNavigateReplaceTo({ pathname: ROUTES_PATH.ROUTE_MAIN.PATH });
  } catch (e: any) {
    openAlert({
      type: "error",
      message: "กรุณาลองใหม่อีกครั้ง",
      title: "เข้าสู่ระบบไม่สำเร็จ",
    });
    yield put({ type: type.AUTHEN_LOGIN_LARK_FAIL, message: e.message });
  } finally {
    closeLoading();
  }
}

export function* checkAuthen(): Generator<any, void, any> {
  try {
    openLoading();
    const result = yield GET(CHECKAUTH);
    yield put({
      type: type.AUTHEN_CHECK_SUCCESS,
      users: {
        userID: result.UserID,
        userFullName: result.Username,
        userRoleID: result.RoleID,
      },
    });
    if (window.location.pathname === "/") {
      window.location.replace(ROUTES_PATH.ROUTE_MAIN.PATH);
    }
  } catch (e: any) {
    console.warn("JWT Invalid → Redirecting to Login");
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
    yield put({ type: type.AUTHEN_CHECK_FAIL, message: e.message });
  } finally {
    closeLoading();
  }
}
