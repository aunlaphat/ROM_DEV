import * as type from "./types";

export function login(payload: any) {
  return {
    type: type.AUTHEN_LOGIN_REQ,
    payload: payload,
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
