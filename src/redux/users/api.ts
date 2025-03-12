// src/redux/users/api.ts

import { call, put, delay } from "redux-saga/effects";
import { SagaIterator } from "redux-saga";
import { notification } from "antd";
import { GET, POST, PATCH, DELETE } from "../../services";
import { USER, CONSTANT } from "../../services/path";
import { logger } from "../../utils/logger";
import { openLoading, closeLoading } from "../../components/alert/useAlert";
import { UserActionTypes } from "./types";
import {
  ApiResponse,
  UserResponse,
  AddUserResponse,
  EditUserResponse,
  DeleteUserResponse,
  RoleResponse,
  WarehouseResponse,
} from "./types";
import {
  fetchUsersSuccess,
  fetchUsersFailure,
  fetchUserSuccess,
  fetchUserFailure,
  addUserSuccess,
  addUserFailure,
  editUserSuccess,
  editUserFailure,
  deleteUserSuccess,
  deleteUserFailure,
  fetchRolesSuccess,
  fetchRolesFailure,
  fetchWarehousesSuccess,
  fetchWarehousesFailure,
} from "./action";

// Fetch Users Saga
export function* fetchUsersSaga(action: {
  type: typeof UserActionTypes.FETCH_USERS_REQUEST;
  payload: any;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Fetch Users");

    const { isActive = true, limit = 100, offset = 0 } = action.payload;

    // สร้าง query params
    const queryParams = new URLSearchParams();
    queryParams.append("isActive", String(isActive));
    queryParams.append("limit", String(limit));
    queryParams.append("offset", String(offset));

    // เรียก API endpoint
    logger.api.request(`${USER.LIST}?${queryParams.toString()}`);
    const response = yield call(() =>
      GET(`${USER.LIST}?${queryParams.toString()}`)
    );

    const apiResponse = response.data as ApiResponse<UserResponse[]>;

    if (!apiResponse.success) {
      throw new Error(apiResponse.message || "Failed to fetch users");
    }

    logger.api.success(USER.LIST, {
      count: apiResponse.data.length,
    });

    // ส่ง action เมื่อดึงข้อมูลสำเร็จ
    // ในกรณีนี้ backend อาจไม่ส่ง total กลับมา เลยใช้ค่าประมาณ
    const total =
      apiResponse.data.length < limit ? offset + apiResponse.data.length : 1000;
    yield put(fetchUsersSuccess(apiResponse.data, total));
  } catch (error: any) {
    logger.error("Fetch Users Error", error);

    // ส่ง action เมื่อดึงข้อมูลล้มเหลว
    yield put(fetchUsersFailure(error.message));

    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: "ดึงข้อมูลผู้ใช้ไม่สำเร็จ",
      description: error.response?.data?.message || error.message,
    });
  } finally {
    closeLoading();
    logger.perf.end("Fetch Users");
  }
}

// Fetch Single User Saga
export function* fetchUserSaga(action: {
  type: typeof UserActionTypes.FETCH_USER_REQUEST;
  payload: string;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Fetch User Details");

    const userID = action.payload;

    // เรียก API endpoint
    logger.api.request(USER.GET(userID));
    const response = yield call(() => GET(USER.GET(userID)));

    const apiResponse = response.data as ApiResponse<UserResponse>;

    if (!apiResponse.success) {
      throw new Error(apiResponse.message || "Failed to fetch user details");
    }

    logger.api.success(USER.GET(userID), {
      userID: apiResponse.data.userID,
      userName: apiResponse.data.userName,
    });

    // ส่ง action เมื่อดึงข้อมูลสำเร็จ
    yield put(fetchUserSuccess(apiResponse.data));
  } catch (error: any) {
    logger.error("Fetch User Details Error", error);

    // ส่ง action เมื่อดึงข้อมูลล้มเหลว
    yield put(fetchUserFailure(error.message));

    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: "ดึงข้อมูลผู้ใช้ไม่สำเร็จ",
      description: error.response?.data?.message || error.message,
    });
  } finally {
    closeLoading();
    logger.perf.end("Fetch User Details");
  }
}

// Add User Saga
export function* addUserSaga(action: {
  type: typeof UserActionTypes.ADD_USER_REQUEST;
  payload: any;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Add User');

    const userData = {
      ...action.payload,
      warehouseID: Number(action.payload.warehouseID) // แปลงเป็น number
    };
    
    // แสดงการใช้ name ใน log
    logger.log('info', 'Add User with names', {
      roleID: userData.roleID,
      roleName: userData.roleName,
      warehouseID: userData.warehouseID,
      warehouseName: userData.warehouseName
    });
    
    // เรียก API endpoint
    logger.api.request(USER.ADD, {
      userID: userData.userID,
      roleID: userData.roleID,
      roleName: userData.roleName,
      warehouseID: userData.warehouseID,
      warehouseName: userData.warehouseName
    });
    
    const response = yield call(() => POST(USER.ADD, userData));
    
    const apiResponse = response.data as ApiResponse<AddUserResponse>;
    
    if (!apiResponse.success) {
      throw new Error(apiResponse.message || 'Failed to add user');
    }

    logger.api.success(USER.ADD, {
      userID: apiResponse.data.userID,
      roleID: apiResponse.data.roleID,
      roleName: userData.roleName || "Unknown"
    });

    // ส่ง action เมื่อเพิ่มผู้ใช้สำเร็จ
    yield put(addUserSuccess({
      ...apiResponse.data,
      // เพิ่ม name ที่ส่งไปเมื่อ response ไม่มี name กลับมา
      roleName: userData.roleName,
      warehouseName: userData.warehouseName
    }));

    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'เพิ่มผู้ใช้สำเร็จ',
      description: 'ผู้ใช้ใหม่ถูกเพิ่มเข้าสู่ระบบแล้ว',
    });

    // โหลดรายการผู้ใช้ใหม่
    yield put({ 
      type: UserActionTypes.FETCH_USERS_REQUEST, 
      payload: { isActive: true } 
    });

  } catch (error: any) {
    logger.error('Add User Error', error);
    yield put(addUserFailure(error.message));
    notification.error({
      message: 'เพิ่มผู้ใช้ไม่สำเร็จ',
      description: error.response?.data?.message || error.message,
    });
  } finally {
    closeLoading();
    logger.perf.end('Add User');
  }
}

// Edit User Saga
export function* editUserSaga(action: {
  type: typeof UserActionTypes.EDIT_USER_REQUEST;
  payload: { userID: string, userData: any };
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Edit User');

    const { userID, userData } = action.payload;

    // แสดงการใช้ name ใน log
    logger.log('info', 'Edit User with names', {
      roleID: userData.roleID,
      roleName: userData.roleName,
      warehouseID: userData.warehouseID,
      warehouseName: userData.warehouseName
    });
    
    // เรียก API endpoint ส่งเฉพาะ ID ที่แน่ใจว่า backend ต้องการ
    logger.api.request(USER.EDIT(userID), {
      userID: userData.userID,
      roleID: userData.roleID,
      warehouseID: userData.warehouseID
    });
    
    const response = yield call(() => PATCH(USER.EDIT(userID), userData));
    
    const apiResponse = response.data as ApiResponse<EditUserResponse>;
    
    if (!apiResponse.success) {
      throw new Error(apiResponse.message || 'Failed to edit user');
    }

    logger.api.success(USER.EDIT(userID), {
      userID: apiResponse.data.userID,
      updatedBy: apiResponse.data.updatedBy
    });

    // สร้างข้อมูลที่ถูกต้องตาม type definition
    const combinedData: EditUserResponse = {
      ...apiResponse.data,
      // ใช้ค่าจาก apiResponse ถ้ามี หรือใช้ค่าจาก userData ถ้าไม่มี
      roleName: apiResponse.data.roleName || userData.roleName || '',
      warehouseName: apiResponse.data.warehouseName || userData.warehouseName || ''
    };

    // ส่ง action เมื่อแก้ไขผู้ใช้สำเร็จ
    yield put(editUserSuccess(combinedData));

    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'แก้ไขผู้ใช้สำเร็จ',
      description: 'ข้อมูลผู้ใช้ถูกอัปเดตแล้ว',
    });

    // โหลดรายการผู้ใช้ใหม่
    yield put({ 
      type: UserActionTypes.FETCH_USERS_REQUEST, 
      payload: { isActive: true } 
    });

  } catch (error: any) {
    logger.error('Edit User Error', error);
    yield put(editUserFailure(error.message));
    notification.error({
      message: 'แก้ไขผู้ใช้ไม่สำเร็จ',
      description: error.response?.data?.message || error.message,
    });
  } finally {
    closeLoading();
    logger.perf.end('Edit User');
  }
}

// Delete User Saga
export function* deleteUserSaga(action: {
  type: typeof UserActionTypes.DELETE_USER_REQUEST;
  payload: string;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Delete User");

    const userID = action.payload;

    // เรียก API endpoint - ส่ง empty object เป็น body
    logger.api.request(USER.DELETE(userID));

    const response = yield call(() => DELETE(USER.DELETE(userID), {}));

    const apiResponse = response.data as ApiResponse<DeleteUserResponse>;

    if (!apiResponse.success) {
      throw new Error(apiResponse.message || "Failed to delete user");
    }

    logger.api.success(USER.DELETE(userID), {
      userID: apiResponse.data.userID,
      userName: apiResponse.data.userName,
    });

    // ส่ง action เมื่อลบผู้ใช้สำเร็จ
    yield put(deleteUserSuccess(apiResponse.data));

    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: "ลบผู้ใช้สำเร็จ",
      description: apiResponse.data.message || "ผู้ใช้ถูกลบออกจากระบบแล้ว",
    });

    // โหลดรายการผู้ใช้ใหม่
    yield put({
      type: UserActionTypes.FETCH_USERS_REQUEST,
      payload: { isActive: true },
    });
  } catch (error: any) {
    logger.error("Delete User Error", error);

    // ส่ง action เมื่อลบผู้ใช้ล้มเหลว
    yield put(deleteUserFailure(error.message));

    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: "ลบผู้ใช้ไม่สำเร็จ",
      description: error.response?.data?.message || error.message,
    });
  } finally {
    closeLoading();
    logger.perf.end("Delete User");
  }
}

// Fetch Roles Saga
export function* fetchRolesSaga(): SagaIterator {
  try {
    logger.perf.start("Fetch Roles");

    // เรียก API endpoint ที่ถูกต้อง
    logger.api.request(CONSTANT.ROLES);

    const response = yield call(() => GET(CONSTANT.ROLES));

    const apiResponse = response.data as ApiResponse<RoleResponse[]>;

    if (!apiResponse.success) {
      throw new Error(apiResponse.message || "Failed to fetch roles");
    }

    logger.api.success(CONSTANT.ROLES, {
      count: apiResponse.data.length,
    });

    // ส่ง action เมื่อดึงข้อมูลสำเร็จ
    yield put(fetchRolesSuccess(apiResponse.data));
  } catch (error: any) {
    logger.error("Fetch Roles Error", error);

    // ส่ง action เมื่อดึงข้อมูลล้มเหลว
    yield put(
      fetchRolesFailure(
        error.response?.data?.message ||
          error.message ||
          "ไม่สามารถดึงข้อมูลบทบาทได้"
      )
    );

    // แสดงข้อความแจ้งเตือน (เพิ่มเติม)
    notification.error({
      message: "ดึงข้อมูลบทบาทไม่สำเร็จ",
      description:
        error.response?.data?.message ||
        error.message ||
        "กรุณาลองใหม่อีกครั้ง",
    });
  } finally {
    logger.perf.end("Fetch Roles");
  }
}

// Fetch Warehouses Saga
export function* fetchWarehousesSaga(): SagaIterator {
  try {
    logger.perf.start("Fetch Warehouses");

    // เรียก API endpoint
    logger.api.request(CONSTANT.WAREHOUSES);

    const response = yield call(() => GET(CONSTANT.WAREHOUSES));

    const apiResponse = response.data as ApiResponse<WarehouseResponse[]>;

    if (!apiResponse.success) {
      throw new Error(apiResponse.message || "Failed to fetch warehouses");
    }

    logger.api.success(CONSTANT.WAREHOUSES, {
      count: apiResponse.data.length,
    });

    // ส่ง action เมื่อดึงข้อมูลสำเร็จ
    yield put(fetchWarehousesSuccess(apiResponse.data));
  } catch (error: any) {
    logger.error("Fetch Warehouses Error", error);

    // ส่ง action เมื่อดึงข้อมูลล้มเหลว
    yield put(fetchWarehousesFailure(error.message));
  } finally {
    logger.perf.end("Fetch Warehouses");
  }
}
