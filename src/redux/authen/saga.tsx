import { takeEvery } from "redux-saga/effects";
import { login, logout, checkAuthen, login_lark } from "./api";
import * as type from "./types";

function* authenSaga() {
  yield takeEvery(type.AUTHEN_LOGIN_REQ, login);
  yield takeEvery(type.AUTHEN_LOGIN_LARK_REQ, login_lark);
  yield takeEvery(type.AUTHEN_LOGOUT_REQ, logout);
  yield takeEvery(type.AUTHEN_CHECK_REQ, checkAuthen);
}

export default authenSaga;
