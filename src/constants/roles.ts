export enum RoleID {
  ADMIN = 1,
  ACCOUNTING = 2,
  WAREHOUSE = 3,
  TRADE_CONSIGN = 4,
  VIEWER = 5,
}

// ✅ Label ของแต่ละ Role
export const RoleLabel = {
  [RoleID.ADMIN]: "แอดมิน",
  [RoleID.ACCOUNTING]: "เจ้าหน้าที่บัญชี",
  [RoleID.WAREHOUSE]: "เจ้าหน้าที่คลัง",
  [RoleID.TRADE_CONSIGN]: "เจ้าหน้าที่ฝ่ายค้าขาย",
  [RoleID.VIEWER]: "ผู้ดูข้อมูล",
};

// ✅ Role-based Route Access (กำหนดหน้าแต่ละ Role สามารถเข้าถึงได้)
export const ROLE_ACCESS = {
  [RoleID.ADMIN]: ["home", "manage_user", "return_order"],
  [RoleID.ACCOUNTING]: ["home", "return_order"],
  [RoleID.WAREHOUSE]: ["home", "return_order"],
  [RoleID.TRADE_CONSIGN]: ["home"],
  [RoleID.VIEWER]: ["home"],
};

// ✅ ฟังก์ชันตรวจสอบว่า Role สามารถเข้าถึง Route นี้ได้หรือไม่
export const canAccessRoute = (roleID: RoleID, routeKey: string): boolean => {
  return ROLE_ACCESS[roleID]?.includes(routeKey) || false;
};
