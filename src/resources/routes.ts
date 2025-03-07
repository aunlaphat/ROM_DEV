import { lazy } from 'react';
import { RoleID } from '../constants/roles';
import { IconKey } from './icon';

// Lazy loading components
const Home = lazy(() => import("../screens/Home"));
const Login = lazy(() => import("../screens/auth").then(module => ({ default: module.Login })));
const NotFound = lazy(() => import("../screens/NotFound").then(module => ({ default: module.NotFound })));
const CreateReturnOrderMKP = lazy(() => import("../screens/Orders/Marketplace"));
const DraftAndConfirm = lazy(() => import("../screens/Draft&Confirm").then(module => ({ default: module.DraftAndConfirm })));
const ConfirmReturnTrade = lazy(() => import("../screens/ConfirmReturnTrade").then(module => ({ default: module.ConfirmReturnTrade })));
const Report = lazy(() => import("../screens/Report").then(module => ({ default: module.Report })));
const ManageUser = lazy(() => import("../screens/ManageUser"));

export const ROUTE_LOGIN = process.env.REACT_APP_FRONTEND_URL + "/";

// โครงสร้าง Route ที่มีการเพิ่ม allowedRoles
interface RouteType {
  KEY: string;
  PATH: string;
  LABEL: string;
  COMPONENT: React.ComponentType<any>;
  allowedRoles: RoleID[];  // เพิ่ม allowedRoles สำหรับการตรวจสอบสิทธิ์
  icon?: IconKey;           // เพิ่ม icon key สำหรับแสดงใน sidebar (ใช้ชื่อฟังก์ชันจาก Icon object)
}

// แยก routes เป็น 2 ประเภท: ไม่ต้องการยืนยันตัวตน และ ต้องการยืนยันตัวตน
export interface Routes {
  [key: string]: RouteType;
}

// Routes ที่ไม่ต้องการยืนยันตัวตน
export const PUBLIC_ROUTES: Routes = {
  ROUTE_LOGIN: {
    KEY: "login",
    PATH: "/",
    LABEL: "Login",
    COMPONENT: Login,
    allowedRoles: [],  // ไม่มีการตรวจสอบสิทธิ์
  },
};

// Routes ที่ต้องการยืนยันตัวตน
export const PROTECTED_ROUTES: Routes = {
  ROUTE_MAIN: {
    KEY: "home",
    PATH: "/home",
    LABEL: "Home",
    COMPONENT: Home,
    allowedRoles: [RoleID.ADMIN, RoleID.TRADE_CONSIGN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.VIEWER],
    icon: 'Home',
  },
  ROUTE_CREATERETURNORDERMKP: {
    KEY: "createReturnOrderMKP",
    PATH: "/create-return-order-mkp",
    LABEL: "Create Return Order MKP",
    COMPONENT: CreateReturnOrderMKP,
    allowedRoles: [RoleID.ADMIN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.TRADE_CONSIGN, RoleID.VIEWER],
    icon: 'Edit1',
  },
  ROUTE_CONFIRMRETURNTRADE: {
    KEY: "confirmReturnTrade",
    PATH: "/confirm-return-trade",
    LABEL: "Confirm Return Trade",
    COMPONENT: ConfirmReturnTrade,
    allowedRoles: [RoleID.ADMIN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.TRADE_CONSIGN, RoleID.VIEWER],
    icon: 'Confirm',
  },
  ROUTE_REPORT: {
    KEY: "report",
    PATH: "/report",
    LABEL: "Report",
    COMPONENT: Report,
    allowedRoles: [RoleID.ADMIN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.TRADE_CONSIGN, RoleID.VIEWER],
    icon: 'Report',
  },
  ROUTE_DRAFTANDCONFIRM: {
    KEY: "draftAndConfirm",
    PATH: "/draft-and-confirm",
    LABEL: "Draft and Confirm",
    COMPONENT: DraftAndConfirm,
    allowedRoles: [RoleID.ADMIN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.TRADE_CONSIGN, RoleID.VIEWER],
    icon: 'Draft',
  },
  ROUTE_MANAGEUSER: {
    KEY: "manageUser",
    PATH: "/manage-user",
    LABEL: "Manage User",
    COMPONENT: ManageUser,
    allowedRoles: [RoleID.ADMIN], // เฉพาะ Admin เท่านั้น
    icon: 'manageUser',
  },
  ROUTE_NOTFOUND: {
    KEY: "notFound",
    PATH: "*",
    LABEL: "Page Not Found",
    COMPONENT: NotFound,
    allowedRoles: [RoleID.ADMIN, RoleID.TRADE_CONSIGN, RoleID.ACCOUNTING, RoleID.WAREHOUSE, RoleID.VIEWER],
  },
};

// นำออกเพื่อการใช้งานร่วมกันทั้งสองประเภท
export const ROUTES = {
  ...PUBLIC_ROUTES,
  ...PROTECTED_ROUTES,
};

// Export เพื่อความเข้ากันได้กับโค้ดเดิม
export const ROUTES_NO_AUTH = PUBLIC_ROUTES;

// Utility function สำหรับตรวจสอบว่า route นั้นเข้าถึงได้สำหรับบทบาทนั้นๆ หรือไม่
export const isRouteAccessible = (routeKey: string, userRoleID?: RoleID): boolean => {
  if (!userRoleID) return false;
  
  const route = PROTECTED_ROUTES[routeKey];
  if (!route) return false;
  
  return route.allowedRoles.includes(userRoleID);
};

// Helper function สำหรับ filter routes ตามบทบาท
export const getAccessibleRoutes = (userRoleID?: RoleID): RouteType[] => {
  if (!userRoleID) return [];
  
  return Object.values(PROTECTED_ROUTES).filter(route => 
    route.allowedRoles.includes(userRoleID) && route.KEY !== 'notFound'
  );
};