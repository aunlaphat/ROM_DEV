import { AuthActionTypes, LoginPayload } from '../../types/auth.types';

export const login = (payload: LoginPayload) => ({
  type: AuthActionTypes.AUTHEN_LOGIN_REQ,
  payload,
});

export const logout = () => ({
  type: AuthActionTypes.AUTHEN_LOGOUT_REQ,
});

export const checkAuth = () => ({
  type: AuthActionTypes.AUTHEN_CHECK_REQ,
});

export const loginLark = (payload: any) => ({
  type: AuthActionTypes.AUTHEN_LOGIN_LARK_REQ,
  payload,
});