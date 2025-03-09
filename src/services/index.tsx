import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { getCookies } from "../store/useCookies";
import { logger } from '../utils/logger';

/**
 * Axios Instance พร้อม Configuration
 * - baseURL: จาก Environment Variables
 * - timeout: 10 วินาที
 * - withCredentials: ส่ง cookies ในการร้องขอข้ามโดเมน
 * - headers: ตั้งค่า Content-Type เป็น JSON
 */
const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_BACKEND_URL,
  timeout: 10000, // ⏳ Timeout 10 วินาที
  withCredentials: true, // 🍪 ใช้ Cookies สำหรับ Authentication
  headers: { "Content-Type": "application/json" } // 📌 Default Header เป็น JSON
});

/**
 * Interceptor สำหรับการส่งคำขอ
 * - เพิ่ม Authorization header ถ้ามี Token
 * - บันทึก Log ทุกคำขอ
 */
axiosInstance.interceptors.request.use((config) => {
  // ดึง Token จาก Cookie หรือ localStorage
  const token = getCookies("jwt") || localStorage.getItem("access_token");
  
  // เพิ่ม Authorization header ถ้ามี token
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  
  // ใช้ logger แทน console.log
  logger.api.request(config.url || 'unknown', {
    method: config.method,
    data: config.data
  });
  
  return config;
});

/**
 * Interceptor สำหรับการรับการตอบกลับ
 * - บันทึก Log การตอบกลับที่สำเร็จ
 * - จัดการข้อผิดพลาดและบันทึก Log
 */
axiosInstance.interceptors.response.use(
  (response: AxiosResponse) => {
    // ใช้ logger แทน console.log
    logger.api.success(response.config.url || 'unknown', {
      status: response.status,
      data: response.data
    });
    return response;
  },
  (error) => {
    // ใช้ logger.api.error แทน console.error
    logger.api.error(error.config?.url || 'unknown', {
      message: error.message,
      response: error.response?.data,
      status: error.response?.status
    });
    
    // ตรวจสอบสถานะ 401 Unauthorized (Token หมดอายุหรือไม่ถูกต้อง)
    if (error.response?.status === 401) {
      // อาจจะเรียก action logout หรือเปลี่ยนเส้นทางไปยังหน้า login
      // แต่ต้องระวังการเรียกวนซ้ำ (infinite loop)
      logger.error('Unauthorized API request - Token may be invalid or expired');
    }
    
    return Promise.reject(error);
  }
);

export default axiosInstance;

/**
 * ส่งคำขอ POST ไปยัง API
 * @param url เส้นทาง API endpoint
 * @param data ข้อมูลที่จะส่ง
 * @param config การตั้งค่าเพิ่มเติมสำหรับ Axios
 */
export const POST = (url: string, data: any, config?: AxiosRequestConfig) => 
  axiosInstance.post(url, data, { ...config });

/**
 * ส่งคำขอ GET ไปยัง API
 * @param url เส้นทาง API endpoint
 * @param config การตั้งค่าเพิ่มเติมสำหรับ Axios
 */
export const GET = (url: string, config?: AxiosRequestConfig) => 
  axiosInstance.get(url, { ...config });

/**
 * ส่งคำขอ PATCH ไปยัง API (ใช้สำหรับอัพเดทบางส่วน)
 * @param url เส้นทาง API endpoint
 * @param data ข้อมูลที่จะส่ง
 * @param config การตั้งค่าเพิ่มเติมสำหรับ Axios
 */
export const PATCH = (url: string, data: any, config?: AxiosRequestConfig) => 
  axiosInstance.patch(url, data, { ...config });

/**
 * ส่งคำขอ PUT ไปยัง API (ใช้สำหรับอัพเดททั้งหมด)
 * @param url เส้นทาง API endpoint
 * @param data ข้อมูลที่จะส่ง
 * @param config การตั้งค่าเพิ่มเติมสำหรับ Axios
 */
export const PUT = (url: string, data: any, config?: AxiosRequestConfig) => 
  axiosInstance.put(url, data, { ...config });

/**
 * ส่งคำขอ DELETE ไปยัง API
 * @param url เส้นทาง API endpoint
 * @param config การตั้งค่าเพิ่มเติมสำหรับ Axios
 */
export const DELETE = (url: string, config?: AxiosRequestConfig) => 
  axiosInstance.delete(url, { ...config });

/**
 * ส่งคำขอ POST พร้อม FormData ไปยัง API (ใช้สำหรับอัพโหลดไฟล์)
 * @param url เส้นทาง API endpoint
 * @param formData ข้อมูล FormData ที่จะส่ง
 * @param config การตั้งค่าเพิ่มเติมสำหรับ Axios
 */
export const UPLOAD = (url: string, formData: FormData, config?: AxiosRequestConfig) => 
  axiosInstance.post(url, formData, { 
    ...config,
    headers: {
      ...config?.headers,
      'Content-Type': 'multipart/form-data'
    }
  });