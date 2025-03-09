// src/hooks/useAuth.tsx
import { createContext, useContext, useCallback, useEffect } from 'react';
import { useDispatch, useSelector } from "react-redux";
import { NavigateTo, windowNavigateReplaceTo } from "../utils/navigation";
import { ROUTES, ROUTE_LOGIN } from "../resources/routes";
import { login as loginAction, loginLark as loginLarkAction, logout as logoutAction, checkAuthen } from "../redux/auth/action";
import { getCookies } from "../store/useCookies";
import { logger } from "../utils/logger";
import { RootState } from '../redux/types';
import { LoginRequest, LarkLoginRequest } from '../redux/auth/interface';

// Interface สำหรับข้อมูลที่จะเก็บใน Context
interface AuthContextType {
  // User Data
  userID: string;
  userName: string;
  fullNameTH: string;
  roleID: number;
  roleName: string;
  
  // Auth State
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  
  // Auth Methods
  login: (values: LoginRequest) => void;
  loginLark: (values: LarkLoginRequest) => void;
  logout: () => void;
  toLogin: () => void;
  redirectToMain: () => void;
  checkAuth: () => void;
}

// สร้าง Context
const AuthContext = createContext<AuthContextType | null>(null);

// Provider Component
export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const dispatch = useDispatch();
  
  // ใช้ useSelector เพื่อเข้าถึง Redux state
  const auth = useSelector((state: RootState) => state.auth);
  
  // อ่านข้อมูล user จาก cookies และ localStorage เมื่อโหลดคอมโพเนนต์
  useEffect(() => {
    logger.perf.start('Auth Provider: Initial Check');
    
    const token = getCookies('jwt') || localStorage.getItem("access_token");
    
    logger.log('debug', `[Auth Provider] Initial auth check`, {
      hasToken: !!token,
      isAuthenticated: auth.isAuthenticated,
      hasUser: !!auth.user,
      timestamp: new Date().toISOString()
    });
    
    // Auto-check authentication if token exists and no user data
    if (token && !auth.isAuthenticated) {
      logger.log('info', `[Auth Provider] Token found but not authenticated, checking auth...`);
      dispatch(checkAuthen());
    }
    
    logger.perf.end('Auth Provider: Initial Check');
  }, [dispatch, auth.isAuthenticated]);

  // Redirect to login page
  const toLogin = useCallback(() => {
    logger.perf.start('Navigation: To Login Page');
    logger.navigation.to(ROUTE_LOGIN, {
      method: 'windowNavigateReplaceTo',
      from: 'useAuth.toLogin'
    });
    
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
    logger.perf.end('Navigation: To Login Page');
  }, []);

  // Handle login
  const handleLogin = useCallback((values: LoginRequest) => {
    logger.perf.start('Auth: Login Flow');
    logger.log('info', `[useAuth] Login initiated`, {
      username: values.userName,
      timestamp: new Date().toISOString()
    });
    
    dispatch(loginAction(values));
    logger.perf.end('Auth: Login Flow');
  }, [dispatch]);

  // Handle Lark login
  const handleLarkLogin = useCallback((values: LarkLoginRequest) => {
    logger.perf.start('Auth: Lark Login Flow');
    logger.log('info', `[useAuth] Lark login initiated`, {
      userID: values.userID,
      timestamp: new Date().toISOString()
    });
    
    dispatch(loginLarkAction(values));
    logger.perf.end('Auth: Lark Login Flow');
  }, [dispatch]);

  // Handle logout
  const handleLogout = useCallback(() => {
    logger.perf.start('Auth: Logout Flow');
    logger.log('info', `[useAuth] Logout initiated`, {
      timestamp: new Date().toISOString()
    });

    dispatch(logoutAction());
    logger.perf.end('Auth: Logout Flow');
  }, [dispatch]);

  // Force check auth
  const handleCheckAuth = useCallback(() => {
    logger.log('info', '[useAuth] Force check auth');
    dispatch(checkAuthen());
  }, [dispatch]);

  // Redirect to main page if authenticated
  const redirectToMain = useCallback(() => {
    const hasToken = getCookies("jwt") || localStorage.getItem("access_token");
    
    logger.log(hasToken ? 'warn' : 'info',
      `[useAuth] Main page redirect attempt`, {
        status: hasToken ? 'blocked' : 'proceeding',
        destination: ROUTES.ROUTE_MAIN.PATH,
        hasToken,
        timestamp: new Date().toISOString()
      }
    );

    if (!hasToken) {
      logger.navigation.to(ROUTES.ROUTE_MAIN.PATH, {
        method: 'NavigateTo',
        from: 'useAuth.redirectToMain'
      });
      NavigateTo({ pathname: ROUTES.ROUTE_MAIN.PATH });
    }
  }, []);

  // Use values from Redux auth state
  const value: AuthContextType = {
    userID: auth.user?.userID || '',
    roleID: auth.user?.roleID || 0,
    userName: auth.user?.userName || '',
    fullNameTH: auth.user?.fullNameTH || '',  // ใช้ fullNameTH จาก user
    roleName: auth.user?.roleName || '',
    
    isAuthenticated: auth.isAuthenticated,
    isLoading: auth.loading,
    error: auth.error,
    
    login: handleLogin,
    loginLark: handleLarkLogin,
    logout: handleLogout,
    toLogin,
    redirectToMain,
    checkAuth: handleCheckAuth
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
    logger.error('[useAuth] Hook used outside of AuthProvider', {
      timestamp: new Date().toISOString()
    });
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// Re-export AuthProvider เพื่อให้ใช้งานง่าย
export { AuthContext };