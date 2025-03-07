import { takeLatest, Effect } from 'redux-saga/effects';
import { AuthActionTypes } from '../../../types/auth.types';
import { loginSaga, logoutSaga, checkAuthSaga, loginLarkSaga } from './authSagas';

export default function* authSaga(): Generator<Effect, void, unknown> {
  yield takeLatest(AuthActionTypes.AUTHEN_LOGIN_REQ, loginSaga);
  yield takeLatest(AuthActionTypes.AUTHEN_LOGIN_LARK_REQ, loginLarkSaga);
  yield takeLatest(AuthActionTypes.AUTHEN_LOGOUT_REQ, logoutSaga);
  yield takeLatest(AuthActionTypes.AUTHEN_CHECK_REQ, checkAuthSaga);
}