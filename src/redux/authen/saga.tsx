import { takeEvery } from "redux-saga/effects";
import { login, logout, checkAuthen } from "./api";
import * as type from "./types";

function* authenSaga() {
  yield takeEvery(type.AUTHEN_LOGIN_REQ, login);
  yield takeEvery(type.AUTHEN_LOGOUT_REQ, logout);
  yield takeEvery(type.AUTHEN_CHECK_REQ, checkAuthen);
}

export default authenSaga;
