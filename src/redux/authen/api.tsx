import { put } from "redux-saga/effects";
import { GET, POST } from "../../services";
import { CHECKAUTH, LARK_LOGIN, LOGIN } from "../../services/path";
import * as type from "./types";
import { windowNavigateReplaceTo } from "../../utils/navigation";
import { ROUTES_PATH, ROUTE_LOGIN } from "../../resources/routes-name";
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

export function* logout() {
  try {
    openAlert({ type: "info", title: "ออกจากระบบสำเร็จ" });
    openLoading();
    removeCookies("jwt");
    yield put({ type: type.AUTHEN_LOGOUT_SUCCESS });
  } catch (e: any) {
    openAlert({ type: "error", message: e.data });
    yield put({ type: type.AUTHEN_LOGOUT_FAIL, message: e.message });
  } finally {
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
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
    yield put({ type: type.AUTHEN_CHECK_FAIL, message: e.message });
  } finally {
    closeLoading();
  }
}
