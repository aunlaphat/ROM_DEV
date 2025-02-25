export enum RoleID {
  ADMIN = 1,
  ACCOUNTING = 2,
  WAREHOUSE = 3,
  TRADE_CONSIGN = 4,
  VIEWER = 5
}

export const RoleLabel = {
  [RoleID.ADMIN]: 'แอดมิน',
  [RoleID.ACCOUNTING]: 'เจ้าหน้าที่บัญชี',
  [RoleID.WAREHOUSE]: 'เจ้าหน้าที่คลัง',
  [RoleID.TRADE_CONSIGN]: 'เจ้าหน้าที่ฝ่ายค้าขาย',
  [RoleID.VIEWER]: 'ผู้ดูข้อมูล'
};
