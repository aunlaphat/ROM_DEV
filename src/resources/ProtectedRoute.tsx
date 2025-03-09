import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../hooks/auth';
import { ROUTES_NO_AUTH } from './routes';
import Loading from '../components/loading';
import { RoleID } from '../constants/roles';
import { logger } from '../utils/logger';

interface ProtectedRouteProps {
  allowedRoles?: RoleID[];
  children?: React.ReactNode;
}

/**
 * คอมโพเนนต์สำหรับป้องกันเส้นทางที่ต้องการการยืนยันตัวตน
 * และสามารถจำกัดเฉพาะบทบาทที่กำหนดได้
 */
const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  allowedRoles = [],
  children,
}) => {
  const {
    isAuthenticated,
    loading,
    hasInitialized,
    roleID,
  } = useAuth();

  // ถ้ากำลังโหลดข้อมูลหรือยังไม่เริ่มต้น ให้แสดง loading
  if (loading || !hasInitialized) {
    return <Loading />;
  }

  // ถ้ายังไม่ได้ยืนยันตัวตน ให้ redirect ไปหน้า login
  if (!isAuthenticated) {
    logger.log('warn', 'Access to protected route denied - not authenticated');
    return <Navigate to={ROUTES_NO_AUTH.ROUTE_LOGIN.PATH} replace />;
  }

  // ถ้ามีการระบุบทบาทที่อนุญาตและผู้ใช้ไม่ได้อยู่ในบทบาทที่อนุญาต
  if (allowedRoles.length > 0 && roleID && !allowedRoles.includes(roleID)) {
    logger.log('warn', 'Access to protected route denied - insufficient permissions', {
      userRoleID: roleID,
      allowedRoles,
    });
    return <Navigate to="/unauthorized" replace />;
  }

  // คืนค่า children หรือ Outlet ในกรณีที่ใช้กับ nested routes
  return <>{children || <Outlet />}</>;
};

export default ProtectedRoute;