import { createContext, useContext, useCallback, useEffect, useState } from 'react';
import { NavigateTo, windowNavigateReplaceTo } from '../../utils/navigation';
import { ROUTES, ROUTE_LOGIN } from '../../resources/routes';
import { getCookies, removeCookies } from '../../store/useCookies';
import { logger } from '../../utils/logger';
import { useAuthState } from './useAuthState';
import { useAuthActions } from './useAuthActions';
import { LoginPayload } from '../../types/auth.types';

// Interface สำหรับข้อมูลที่จะเก็บใน Context
interface AuthContextType {
  // User Data
  userID: string;
  roleID: number;
  userName: string;
  roleName: string;
  fullName: string;
  nickName: string;
  departmentNo: string;
  platform: string;
  isAuthenticated: boolean;
  loading: boolean;
  
  // Auth Methods
  login: (values: LoginPayload) => void;
  logout: () => void;
  checkAuth: () => void;
  loginLark: (payload: any) => void;
  toLogin: () => void;
  redirectToMain: () => void;
  
  // Auth Status
  hasInitialized: boolean;
}

// สร้าง Context
export const AuthContext = createContext<AuthContextType | null>(null);

// Provider Component
export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isAuthenticated, loading } = useAuthState();
  const { login, logout, checkAuth, loginLark } = useAuthActions();
  const [hasInitialized, setHasInitialized] = useState<boolean>(false);

  // ตรวจสอบ token เมื่อ component mount
  useEffect(() => {
    const initializeAuth = async () => {
      try {
        const token = getCookies('jwt') || localStorage.getItem('access_token');
        logger.log('info', 'Authentication Status', {
          tokenExists: !!token,
          timestamp: new Date().toISOString(),
        });
        
        if (token && !isAuthenticated) {
          checkAuth();
        }
      } finally {
        setHasInitialized(true);
      }
    };
    
    initializeAuth();
  }, [checkAuth, isAuthenticated]);

  const toLogin = useCallback(() => {
    logger.perf.start('Navigation: Login Page');
    logger.log('info', 'Redirecting', {
      destination: 'Login Page',
      route: ROUTE_LOGIN,
    });
    
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
    logger.perf.end('Navigation: Login Page');
  }, []);

  const redirectToMain = useCallback(() => {
    const hasToken = getCookies('jwt') || localStorage.getItem('access_token');
    logger.log(
      hasToken ? 'warn' : 'info',
      'Main Page Redirect',
      {
        status: hasToken ? 'blocked' : 'proceeding',
        destination: ROUTES.ROUTE_MAIN.PATH,
        hasToken,
        timestamp: new Date().toISOString(),
      }
    );

    if (!hasToken) {
      NavigateTo({ pathname: ROUTES.ROUTE_MAIN.PATH });
    }
  }, []);

  const handleLogout = useCallback(() => {
    logout();
    removeCookies('jwt');
  }, [logout]);

  const value: AuthContextType = {
    userID: user?.userID || '',
    roleID: user?.roleID || 0,
    userName: user?.userName || '',
    roleName: user?.roleName || '',
    fullName: user?.fullName || user?.fullNameTH || '',
    nickName: user?.nickName || '',
    departmentNo: user?.departmentNo || '',
    platform: user?.platform || '',
    isAuthenticated,
    loading,
    hasInitialized,
    login,
    logout: handleLogout,
    checkAuth,
    loginLark,
    toLogin,
    redirectToMain,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Hook สำหรับใช้งาน Context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};