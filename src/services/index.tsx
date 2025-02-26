import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { env } from "../utils/env/config";
import { getCookies } from "../store/useCookies";

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_BACKEND_URL,
  timeout: 10000, // ⏳ Timeout 10 วินาที
  withCredentials: true, // 🍪 ใช้ Cookies สำหรับ Authentication
  headers: { "Content-Type": "application/json" }, // 📌 Default Header เป็น JSON
});

// ✅ Interceptor: Request
axiosInstance.interceptors.request.use((config) => {
  console.log("📤 [API] Request:", {
    method: config.method,
    url: config.url,
    data: config.data,
  });
  return config;
});

// ✅ Interceptor: Response
axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log("📥 [API] Response:", {
      status: response.status,
      data: response.data,
      cookies: document.cookie,
    });
    return response;
  },
  (error) => {
    console.error("❌ [API] Error:", {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status,
    });
    return Promise.reject(error);
  }
);

export default axiosInstance;


export const POST = (url: string, data: any, config?: AxiosRequestConfig) => 
  axiosInstance.post(url, data, { ...config });

export const GET = (url: string, config?: AxiosRequestConfig) => 
  axiosInstance.get(url, { ...config });

export const PATCH = async (path: string, data: any) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await axiosInstance.patch(`${path}`, data, config);
    data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};

export const PUT = async (path: string, data: any, header?: any) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await axiosInstance.put(`${path}`, data, header ? header : config);
    data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};

export const DELETE = async (path: string, data: any, header?: any) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await axiosInstance.delete(`${path}`, {
      headers: header ? header : config,
      data,
    });
    data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};

export const UPLOAD = async (
  path: string,
  formdata: FormData,
  header?: any
) => {
  const config = {
    headers: { Authorization: `Bearer ${getCookies("jwt")}` },
  };
  try {
    const response = await axiosInstance.post(
      `${path}`,
      formdata,
      header ? header : config
    );
    let data = await response.data;
    return data;
  } catch (error: any) {
    console.log(error.response);
    throw error.response;
  }
};