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

export interface Role {
    id: string;
    name: string;
}

export interface Warehouse {
    id: string;
    name: string;
}

export interface ApiResponse<T> {
    success: boolean;
    message: string;
    data: T;
}
