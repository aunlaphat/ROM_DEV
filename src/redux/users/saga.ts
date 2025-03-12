// src/redux/users/saga.ts

import { takeLatest, all } from 'redux-saga/effects';
import { UserActionTypes } from './types';
import {
  fetchUsersSaga,
  fetchUserSaga,
  addUserSaga,
  editUserSaga,
  deleteUserSaga,
  fetchRolesSaga,
  fetchWarehousesSaga
} from './api';
import { logger } from '../../utils/logger';

/**
 * Users Saga - รวมทุก watcher สำหรับการจัดการผู้ใช้
 */
function* usersSaga() {
  // บันทึก log การเริ่มต้น saga
  logger.log('info', '[Users Saga] Initialized user watchers');
  
  yield all([
    // Fetch Users watcher
    takeLatest(UserActionTypes.FETCH_USERS_REQUEST, fetchUsersSaga),
    
    // Fetch Single User watcher
    takeLatest(UserActionTypes.FETCH_USER_REQUEST, fetchUserSaga),
    
    // Add User watcher
    takeLatest(UserActionTypes.ADD_USER_REQUEST, addUserSaga),
    
    // Edit User watcher
    takeLatest(UserActionTypes.EDIT_USER_REQUEST, editUserSaga),
    
    // Delete User watcher
    takeLatest(UserActionTypes.DELETE_USER_REQUEST, deleteUserSaga),
    
    // Fetch Roles watcher
    takeLatest(UserActionTypes.FETCH_ROLES_REQUEST, fetchRolesSaga),
    
    // Fetch Warehouses watcher
    takeLatest(UserActionTypes.FETCH_WAREHOUSES_REQUEST, fetchWarehousesSaga)
  ]);
}

export default usersSaga;