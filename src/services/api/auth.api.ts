import axios, { AxiosResponse } from 'axios';
import { LoginPayload, LoginResponse, AuthCheckResponse } from '../../types/auth.types';
import { GET, POST } from '../../services';
import { CHECKAUTH, LOGIN, LOGIN_LARK } from '../../services/path';

export const authAPI = {
  /**
   * ทำการ login ด้วย username และ password
   * Backend endpoint: POST /auth/login
   */
  login: (credentials: LoginPayload): Promise<AxiosResponse<LoginResponse>> => {
    return POST(LOGIN, {
      userName: credentials.username,
      password: credentials.password
    });
  },

  /**
   * ตรวจสอบสถานะการ authentication ของผู้ใช้ปัจจุบัน
   * Backend endpoint: GET /auth/
   */
  checkAuth: (): Promise<AxiosResponse<AuthCheckResponse>> => {
    return GET(CHECKAUTH);
  },

  /**
   * ทำการ logout จากระบบ
   * Backend endpoint: POST /auth/logout
   */
  logout: (): Promise<AxiosResponse<any>> => {
    return POST("/auth/logout", {});
  },

  /**
   * ทำการ login ผ่าน Lark
   * Backend endpoint: POST /auth/login-lark
   */
  loginLark: (payload: any): Promise<AxiosResponse<LoginResponse>> => {
    return POST(LOGIN_LARK, payload);
  }
};