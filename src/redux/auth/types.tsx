import { User, LoginRequest, LarkLoginRequest } from './interface';

// ✅ Action Types Enum
export enum AuthActionTypes {
  AUTHEN_LOGIN_REQ = "@AUTHEN/LOGIN_REQ",
  AUTHEN_LOGIN_SUCCESS = "@AUTHEN/LOGIN_SUCCESS",
  AUTHEN_LOGIN_FAIL = "@AUTHEN/LOGIN_FAIL",

  AUTHEN_LOGOUT_REQ = "@AUTHEN/LOGOUT_REQ",
  AUTHEN_LOGOUT_SUCCESS = "@AUTHEN/LOGOUT_SUCCESS",
  AUTHEN_LOGOUT_FAIL = "@AUTHEN/LOGOUT_FAIL",

  AUTHEN_CHECK_REQ = "@AUTHEN/CHECK_REQ",
  AUTHEN_CHECK_SUCCESS = "@AUTHEN/CHECK_SUCCESS",
  AUTHEN_CHECK_FAIL = "@AUTHEN/CHECK_FAIL",

  AUTHEN_LOGIN_LARK_REQ = "@AUTHEN/LOGIN_LARK_REQ",
  AUTHEN_LOGIN_LARK_SUCCESS = "@AUTHEN/LOGIN_LARK_SUCCESS",
  AUTHEN_LOGIN_LARK_FAIL = "@AUTHEN/LOGIN_LARK_FAIL",
}

// ✅ Action Interfaces - กำหนด Type ให้ Redux Actions
export interface LoginAction {
  type: AuthActionTypes.AUTHEN_LOGIN_REQ;
  payload: LoginRequest;
}

export interface LarkLoginAction {
  type: AuthActionTypes.AUTHEN_LOGIN_LARK_REQ;
  payload: LarkLoginRequest;
}

export interface LoginSuccessAction {
  type: AuthActionTypes.AUTHEN_LOGIN_SUCCESS;
  users: User;
}

export interface AuthCheckAction {
  type: AuthActionTypes.AUTHEN_CHECK_REQ;
}

// ✅ รวม Type ของ Actions ทั้งหมด
export type AuthAction = 
  | LoginAction
  | LarkLoginAction
  | LoginSuccessAction
  | AuthCheckAction;