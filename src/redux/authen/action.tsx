import * as type from "./types";

export function login(payload: any) {
  return {
    type: type.AUTHEN_LOGIN_REQ,
    payload: payload, // { username, password }
  };
}

export function logout(payload?: any) {
  return {
    type: type.AUTHEN_LOGOUT_REQ,
    payload: payload,
  };
}

export function checkAuthen() {
  return {
    type: type.AUTHEN_CHECK_REQ,
  };
}
export function login_lark(payload: any) {
  return {
    type: type.AUTHEN_LOGIN_LARK_REQ,
    payload: payload,
  };
}