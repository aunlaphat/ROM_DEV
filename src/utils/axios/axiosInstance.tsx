import axios from "axios";
import { env } from "../../utils/env/config"; // นำเข้า config.tsx

const api = axios.create({
  baseURL: env.api_base_url, // ใช้ API URL จาก config.tsx
  headers: {
    "Content-Type": "application/json",
  },
});


// เพิ่ม Interceptors เพื่อจัดการ Token หรือ Error Handling ได้
// api.interceptors.request.use(
//   (config) => {
//     const token = localStorage.getItem("token"); // ตัวอย่างการดึง token
//     if (token) {
//       config.headers.Authorization = `Bearer ${token}`;
//     }
//     return config;
//   },
//   (error) => Promise.reject(error)
// );

export default api;
