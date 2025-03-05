import { createContext, useContext, useCallback, useEffect, useState } from 'react';
import { useDispatch } from "react-redux";
import { NavigateTo, windowNavigateReplaceTo } from "../utils/navigation";
import { ROUTES, ROUTE_LOGIN } from "../resources/routes";
import { login, logout } from "../redux/auth/action";
import { getCookies, setCookies, removeCookies } from "../store/useCookies";
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
    logger.log('info', `Authentication Status`, {
      tokenExists: !!token,
      timestamp: new Date().toISOString()
    });
  }, []);

  const toLogin = useCallback(() => {
    logger.perf.start('Navigation: Login Page');
    logger.log('info', `Redirecting`, {
      destination: 'Login Page',
      route: ROUTE_LOGIN
    });
    
    windowNavigateReplaceTo({ pathname: ROUTE_LOGIN });
    logger.perf.end('Navigation: Login Page');
  }, []);

  const handleLogin = useCallback(async (values: any) => {
    try {
      logger.perf.start('Auth: Login Process');
      logger.log('info', `Login Initiated`, {
        username: values.username,
        timestamp: new Date().toISOString()
      });
      
      dispatch(login(values));
      
      logger.log('success', `Login Successful`, {
        username: values.username,
        timestamp: new Date().toISOString()
      });
      
      logger.state.update('Authentication', {
        status: 'authenticated',
        timestamp: new Date().toISOString()
      });
    } catch (error) {
      logger.error('Login Failed', {
        error,
        username: values.username,
        timestamp: new Date().toISOString()
      });
      throw error;
    } finally {
      logger.perf.end('Auth: Login Process');
    }
  }, [dispatch]);

  const handleLogout = useCallback(async () => {
    try {
      logger.perf.start('Auth: Logout Process');
      logger.log('info', `Logout Initiated`, {
        timestamp: new Date().toISOString()
      });

      dispatch(logout());
      removeCookies("jwt");
      
      logger.log('success', `Logout Completed`, {
        timestamp: new Date().toISOString()
      });
      
      logger.state.update('Authentication', {
        status: 'logged_out',
        timestamp: new Date().toISOString()
      });
    } catch (error) {
      logger.error('Logout Failed', {
        error,
        timestamp: new Date().toISOString()
      });
    } finally {
      logger.perf.end('Auth: Logout Process');
    }
  }, [dispatch]);

  const redirectToMain = useCallback(() => {
    const hasToken = getCookies("jwt");
    logger.log(hasToken ? 'warn' : 'info',
      `Main Page Redirect`, {
        status: hasToken ? 'blocked' : 'proceeding',
        destination: ROUTES.ROUTE_MAIN.PATH,
        hasToken,
        timestamp: new Date().toISOString()
      }
    );

    if (!hasToken) {
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
