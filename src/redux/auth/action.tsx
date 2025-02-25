import { AuthActionTypes } from "./types";

export const login = (payload: { username: string; password: string }) => ({
  type: AuthActionTypes.AUTHEN_LOGIN_REQ,
  payload,
});

export const logout = () => ({
  type: AuthActionTypes.AUTHEN_LOGOUT_REQ,
});

export const checkAuthen = () => ({
  type: AuthActionTypes.AUTHEN_CHECK_REQ,
});

export const loginLark = (payload: any) => ({
  type: AuthActionTypes.AUTHEN_LOGIN_LARK_REQ,
  payload,
});
