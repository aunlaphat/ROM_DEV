// src/redux/users/api.ts

import { call, put } from "redux-saga/effects";
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
  AddUserRequest,
  EditUserRequest,
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
    const total =
      apiResponse.data.length < limit ? offset + apiResponse.data.length : 1000;
    yield put(fetchUsersSuccess(apiResponse.data, total));
  } catch (error: any) {
    logger.error("Fetch Users Error", error);
    yield put(fetchUsersFailure(error.message));
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
    yield put(fetchUserFailure(error.message));
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
  payload: AddUserRequest;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Add User");

    // สร้างข้อมูลที่จะส่งไป API - ตรวจสอบให้แน่ใจว่า warehouseID เป็น number
    const userData = {
      userID: action.payload.userID,
      roleID: action.payload.roleID,
      warehouseID: Number(action.payload.warehouseID), // แปลงเป็น number ให้แน่ใจ
    };

    logger.log("info", "Add User request:", {
      userID: userData.userID,
      roleID: userData.roleID,
      warehouseID: userData.warehouseID,
    });

    // เรียก API endpoint
    logger.api.request(USER.ADD, userData);

    const response = yield call(() => POST(USER.ADD, userData));

    const apiResponse = response.data as ApiResponse<AddUserResponse>;

    if (!apiResponse.success) {
      throw new Error(apiResponse.message || "Failed to add user");
    }

    logger.api.success(USER.ADD, {
      userID: apiResponse.data.userID,
      roleID: apiResponse.data.roleID,
      warehouseID: apiResponse.data.warehouseID, // เป็น number เหมือนส่วนอื่นๆ
    });

    // ส่ง action เมื่อเพิ่มผู้ใช้สำเร็จ
    yield put(addUserSuccess(apiResponse.data));

    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: "เพิ่มผู้ใช้สำเร็จ",
      description: "ผู้ใช้ใหม่ถูกเพิ่มเข้าสู่ระบบแล้ว",
    });

    // โหลดรายการผู้ใช้ใหม่
    yield put({
      type: UserActionTypes.FETCH_USERS_REQUEST,
      payload: { isActive: true },
    });
  } catch (error: any) {
    logger.error("Add User Error", error);
    yield put(addUserFailure(error.message));
    notification.error({
      message: "เพิ่มผู้ใช้ไม่สำเร็จ",
      description: error.response?.data?.message || error.message,
    });
  } finally {
    closeLoading();
    logger.perf.end("Add User");
  }
}

// Edit User Saga
export function* editUserSaga(action: {
  type: typeof UserActionTypes.EDIT_USER_REQUEST;
  payload: { userID: string; userData: EditUserRequest };
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Edit User");

    const { userID, userData } = action.payload;

    // ตรวจสอบให้แน่ใจว่า userData มีข้อมูลครบถ้วน
    if (!userData || !userData.userID) {
      throw new Error("ข้อมูลผู้ใช้ไม่ถูกต้องหรือไม่ครบถ้วน");
    }

    // Log ข้อมูลที่ได้รับเพื่อการตรวจสอบ
    logger.log("info", "Original edit user data:", {
      userID,
      userData,
    });

    // สร้างข้อมูลที่จะส่งไป API - ส่ง userData โดยตรง (ไม่ผ่าน apiRequestData)
    // ข้อมูลควรมี roleID และ warehouseID ที่ถูกต้องแล้วจาก UserForm
    const apiRequestData = {
      userID: userData.userID,
      roleID: Number(userData.roleID),
      warehouseID: Number(userData.warehouseID),
    };

    // เพิ่ม roleID ถ้ามีการกำหนดค่า
    if (userData.roleID !== undefined) {
      apiRequestData.roleID = Number(userData.roleID);
    }

    // เพิ่ม warehouseID ถ้ามีการกำหนดค่า
    if (userData.warehouseID !== undefined) {
      apiRequestData.warehouseID = Number(userData.warehouseID);
    }

    // แสดงข้อมูลที่จะส่งไป API
    logger.log("info", "Edit User API request data:", apiRequestData);

    // เรียก API endpoint
    logger.api.request(USER.EDIT(userID), apiRequestData);

    // ส่งข้อมูลที่ได้ปรับแล้วไปยัง API
    const response = yield call(() => PATCH(USER.EDIT(userID), apiRequestData));

    const apiResponse = response.data as ApiResponse<EditUserResponse>;

    if (!apiResponse.success) {
      throw new Error(apiResponse.message || "Failed to edit user");
    }

    logger.api.success(USER.EDIT(userID), {
      userID: apiResponse.data.userID,
      updatedBy: apiResponse.data.updatedBy,
      responseData: apiResponse.data, // แสดงข้อมูล response ทั้งหมด
    });

    // ส่ง action เมื่อแก้ไขผู้ใช้สำเร็จ
    yield put(editUserSuccess(apiResponse.data));

    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: "แก้ไขผู้ใช้สำเร็จ",
      description: "ข้อมูลผู้ใช้ถูกอัปเดตแล้ว",
    });

    // โหลดรายการผู้ใช้ใหม่
    yield put({
      type: UserActionTypes.FETCH_USERS_REQUEST,
      payload: { isActive: true },
    });
  } catch (error: any) {
    logger.error("Edit User Error", error);
    yield put(editUserFailure(error.message));
    notification.error({
      message: "แก้ไขผู้ใช้ไม่สำเร็จ",
      description: error.response?.data?.message || error.message,
    });
  } finally {
    closeLoading();
    logger.perf.end("Edit User");
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

    // เรียก API endpoint
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
    yield put(deleteUserFailure(error.message));
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

    // เรียก API endpoint
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
    yield put(fetchRolesFailure(error.message));
    notification.error({
      message: "ดึงข้อมูลบทบาทไม่สำเร็จ",
      description: error.response?.data?.message || error.message,
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
    yield put(fetchWarehousesFailure(error.message));
  } finally {
    logger.perf.end("Fetch Warehouses");
  }
}
