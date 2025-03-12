// แก้ไขไฟล์ types.ts
export interface User {
    userID: string;
    userName: string;
    nickName: string;
    fullNameTH: string;
    departmentNo: string;
    roleID: number;
    roleName: string;
    warehouseID: number;
    warehouseName: string;
    description: string;
    isActive: boolean;
    createdAt: string;
}

// แก้ไขจาก id/name เป็น roleID/roleName
export interface Role {
    roleID: number;  // เปลี่ยนจาก id: string
    roleName: string;  // เปลี่ยนจาก name: string
}

// แก้ไขจาก id/name เป็น warehouseID/warehouseName
export interface Warehouse {
    warehouseID: number;  // เปลี่ยนจาก id: string
    warehouseName: string;  // เปลี่ยนจาก name: string
}

export interface ApiResponse<T> {
    success: boolean;
    message: string;
    data: T;
}