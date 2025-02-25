import { takeLatest, Effect } from "redux-saga/effects";
import { login, logout, checkAuthen } from "./api";
import { AuthActionTypes } from "./types"; // ✅ ใช้ enum แทน string

function* authenSaga(): Generator<Effect, void, unknown> {
  yield takeLatest(AuthActionTypes.AUTHEN_LOGIN_REQ, login);
  // yield takeLatest(AuthActionTypes.AUTHEN_LOGIN_LARK_REQ, login_lark);
  yield takeLatest(AuthActionTypes.AUTHEN_LOGOUT_REQ, logout);
  yield takeLatest(AuthActionTypes.AUTHEN_CHECK_REQ, checkAuthen);
}

export default authenSaga;
