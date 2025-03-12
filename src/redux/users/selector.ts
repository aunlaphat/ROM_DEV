// src/redux/users/selector.ts

import { createSelector } from 'reselect';
import { RootState } from '../types';
import { UserState } from './types';

// Base selector - ดึง user state
const getUserState = (state: RootState): UserState => state.user;

// ดึงรายการผู้ใช้ทั้งหมด
export const getAllUsers = createSelector(
  [getUserState],
  (userState) => userState.users
);

// ดึงข้อมูลผู้ใช้ปัจจุบัน
export const getCurrentUser = createSelector(
  [getUserState],
  (userState) => userState.currentUser
);

// ดึงสถานะการโหลด
export const getUserLoading = createSelector(
  [getUserState],
  (userState) => userState.loading
);

// ดึงข้อผิดพลาด
export const getUserError = createSelector(
  [getUserState],
  (userState) => userState.error
);

// ดึงข้อมูล pagination
export const getUserPagination = createSelector(
  [getUserState],
  (userState) => userState.pagination
);

// ดึงรายการบทบาททั้งหมด
export const getRoles = createSelector(
  [getUserState],
  (userState) => userState.roles
);

// ดึงรายการคลังสินค้าทั้งหมด
export const getWarehouses = createSelector(
  [getUserState],
  (userState) => userState.warehouses
);

// ค้นหาผู้ใช้ตาม ID
export const getUserById = (userID: string) => createSelector(
  [getAllUsers],
  (users) => users.find(user => user.userID === userID) || null
);

// ตรวจสอบมีผู้ใช้หรือไม่
export const hasUsers = createSelector(
  [getAllUsers],
  (users) => users.length > 0
);