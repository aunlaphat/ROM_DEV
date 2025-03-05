import { createContext, useContext, useCallback, useEffect, useState } from 'react';
import { useDispatch } from "react-redux";
import { NavigateTo, windowNavigateReplaceTo } from "../utils/navigation";
import { ROUTES, ROUTE_LOGIN } from "../resources/routes";
import { login, logout } from "../redux/auth/action";
import { getCookies, setCookies, removeCookies } from "../store/useCookies"; // ลบ decodeToken
import { logger } from "../utils/logger";

// Interface สำหรับข้อมูลที่จะเก็บใน Context
interface AuthContextType {
  // User Data
  userID: string;
  roleID: number;
  channelID: number;
  customerID: string;
  
  // Auth Methods
  login: (values: any) => Promise<void>;
  logout: () => Promise<void>;
  toLogin: () => void;
  redirectToMain: () => void;
}

// สร้าง Context
const AuthContext = createContext<AuthContextType | null>(null);

// Provider Component
export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const dispatch = useDispatch();
  const [userID, setUserID] = useState<string>('');
  const [roleID, setRoleID] = useState<number>(0);
  const [channelID, setChannelID] = useState<number>(0);
  const [customerID, setCustomerID] = useState<string>('');

  // อ่านข้อมูล user จาก cookies
  useEffect(() => {
    const token = getCookies('jwt');
    if (token) {
      // ไม่ต้อง decode token ในฝั่ง frontend
      // ข้อมูล user จะถูกดึงจาก backend โดยอัตโนมัติ
    }
  }, []);

  const toLogin = useCallback(() => {
    logger.auth("info", "Navigating to Login page");
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
  }, []);

  const handleLogin = useCallback(async (values: any) => {
    try {
      logger.auth("info", "Dispatching login request", { username: values.username });
      dispatch(login(values));
    } catch (error) {
      logger.auth("error", "Login failed", { error });
      throw error;
    }
  }, [dispatch]);

  const handleLogout = useCallback(async () => {
    try {
      logger.auth("info", "Processing logout request");
      dispatch(logout());
      removeCookies("jwt");
      logger.auth("info", "User logged out successfully");
    } catch (error) {
      logger.auth("error", "Logout failed", { error });
    }
  }, [dispatch]);

  const redirectToMain = useCallback(() => {
    if (!getCookies("jwt")) {
      logger.auth("info", "Redirecting to Home page");
      NavigateTo({ pathname: ROUTES.ROUTE_MAIN.PATH });
    }
  }, []);

  const value: AuthContextType = {
    userID,
    roleID,
    channelID,
    customerID,
    login: handleLogin,
    logout: handleLogout,
    toLogin,
    redirectToMain
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

// Hook สำหรับใช้งาน Context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// Re-export AuthProvider เพื่อให้ใช้งานง่าย
export { AuthContext };
