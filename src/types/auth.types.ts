// 1. src/types/auth.types.ts - กำหนด Types ทั้งหมดที่เกี่ยวข้องกับ Auth
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

export interface User {
    userID: string;
    userName: string;
    fullName?: string;
    fullNameTH?: string;
    nickName?: string;
    roleID: number;
    roleName: string;
    departmentNo?: string;
    platform?: string;
    userRoleID?: number;
}

export interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
    loading: boolean;
    error?: string;
}

export interface LoginPayload {
    username: string;
    password: string;
}

export interface LoginResponse {
    success: boolean;
    message: string;
    data: string; // token
}

export interface AuthCheckResponse {
    success: boolean;
    message: string;
    data: {
      user: User;
    };
}

export interface AuthActions {
    login: (payload: LoginPayload) => Promise<void>;
    logout: () => Promise<void>;
    checkAuth: () => Promise<void>;
    loginLark: (payload: any) => Promise<void>;
}