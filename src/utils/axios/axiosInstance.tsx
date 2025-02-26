import axios from "axios";
import { env } from "../../utils/env/config"; // นำเข้า config.tsx

const apiClient = axios.create({
  baseURL: env.api_base_url, // ใช้ API URL จาก config.tsx
  headers: {
    "Content-Type": "application/json",
  },
});

export default apiClient;
