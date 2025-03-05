import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
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
      data: error.response?.data, // เพิ่มการแสดงผล data
    });
    return Promise.reject(error);
  }
);

export default axiosInstance;

export const POST = (url: string, data: any, config?: AxiosRequestConfig) => 
  axiosInstance.post(url, data, { ...config });

export const GET = (url: string, config?: AxiosRequestConfig) => 
  axiosInstance.get(url, { ...config });

export const PATCH = (url: string, data: any, config?: AxiosRequestConfig) => 
  axiosInstance.patch(url, data, { ...config });

export const PUT = (url: string, data: any, config?: AxiosRequestConfig) => 
  axiosInstance.put(url, data, { ...config });

export const DELETE = (url: string, config?: AxiosRequestConfig) => 
  axiosInstance.delete(url, { ...config });

export const UPLOAD = (url: string, formData: FormData, config?: AxiosRequestConfig) => 
  axiosInstance.post(url, formData, { ...config });