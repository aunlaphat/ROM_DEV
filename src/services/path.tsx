// Authentication endpoints
export const AUTH = {
    LOGIN: `/auth/login`,
    LOGIN_LARK: `/auth/login-lark`, // แก้ไขจาก LARK_LOGIN เป็น LOGIN_LARK และปรับ path ให้ถูกต้อง
    LOGOUT: `/auth/logout`,
    CHECK: `/auth/` // ใช้ endpoint ที่ถูกต้องตาม backend
  };
  
  // Export individual endpoints for backward compatibility
  export const LOGIN = AUTH.LOGIN;
  export const LOGIN_LARK = AUTH.LOGIN_LARK;
  export const LOGOUT = AUTH.LOGOUT;
  export const CHECKAUTH = AUTH.CHECK;
  
  // Return Order endpoints - MKP (Marketplace)
  export const ORDER = {
    SEARCH: `/order/search`,
    CREATE: `/order/create`,
    GENERATE_SR: `/order/generate-sr`,
    UPDATE_SR: `/order/update-sr`,
    UPDATE_STATUS: `/order/update-status`,
    CANCEL: `/order/cancel`,
    MARK_EDITED: `/order/mark-edited`,
    // เพิ่ม function สำหรับ dynamic path
    DETAIL: (id: string) => `/order/${id}`
  };
  
  // Export individual endpoints for backward compatibility
  export const SEARCHORDER = ORDER.SEARCH;
  export const CREATEBEFORERETURNORDER = ORDER.CREATE;
  export const GENERATESR = ORDER.GENERATE_SR;
  export const UPDATESR = ORDER.UPDATE_SR;
  export const UPDATESTATUS = ORDER.UPDATE_STATUS;
  export const CANCELORDER = ORDER.CANCEL;
  export const MARKEDITED = ORDER.MARK_EDITED;
  
  // User Management endpoints
  export const USER = {
    LIST: `/manage-users/`,
    GET: (id: string) => `/manage-users/${id}`,
    ADD: `/manage-users/add`,
    EDIT: (id: string) => `/manage-users/edit/${id}`,
    DELETE: (id: string) => `/manage-users/delete/${id}`
  };
  
  // Export individual endpoints for backward compatibility
  export const FETCHUSERS = USER.LIST;
  export const FETCHUSERBYID = USER.LIST; // Note: should be used with specific ID
  export const ADDUSER = USER.ADD;
  export const EDITUSER = USER.EDIT(""); // Note: requires ID
  export const DELETEUSER = USER.DELETE(""); // Note: requires ID