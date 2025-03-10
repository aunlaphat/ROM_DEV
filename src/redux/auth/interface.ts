// src/redux/auth/interface.ts

// User Interface - ตรงกับข้อมูลที่ backend ส่งมา (Login struct)
export interface User {
  userID: string;        // จาก UserID ใน Login struct
  userName: string;      // จาก UserName ใน Login struct
  fullNameTH: string;    // จาก FullNameTH ใน Login struct
  nickName: string;      // จาก NickName ใน Login struct
  roleID: number;        // จาก RoleID ใน Login struct
  roleName: string;      // จาก RoleName ใน Login struct
  departmentNo: string;  // จาก DepartmentNo ใน Login struct
  platform: string;      // จาก Platform ใน Login struct
}

// State Type สำหรับ Redux Store
export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
}

// Request Types - ตรงกับที่ backend ต้องการ (request.LoginWeb)
export interface LoginRequest {
  userName: string;  // ตรงกับ UserName ใน request.LoginWeb
  password: string;  // ตรงกับ Password ใน request.LoginWeb
}

// request.LoginLark
export interface LarkLoginRequest {
  userID: string;     // ตรงกับ UserID ใน request.LoginLark
  userName: string;   // ตรงกับ UserName ใน request.LoginLark
}

// Response Types - ตรงกับที่ backend ส่งกลับมา
export interface ApiResponse<T> {
  success: boolean;  // จาก handleResponse ใน backend
  message: string;   // จาก handleResponse ใน backend
  data: T;           // จาก handleResponse ใน backend
}

// Response จาก Login และ LoginLark (เป็น token string)
export interface LoginResponse {
  token: string;  // เป็น string ของ JWT token
}

// Response จาก CheckAuthen (มี user object)
export interface AuthCheckResponse {
  source: string;  // จาก jwt_source ใน middleware
  user: {
    userID: string;        // จาก claims["userID"]
    userName: string;      // จาก claims["userName"]
    fullNameTH: string;      // จาก claims["fullNameTH"]
    nickName: string;      // จาก claims["nickName"]
    roleID: number;        // จาก claims["roleID"]
    roleName: string;      // จาก claims["roleName"]
    departmentNo: string;  // จาก claims["departmentNo"]
    platform: string;      // จาก claims["platform"]
  };
}