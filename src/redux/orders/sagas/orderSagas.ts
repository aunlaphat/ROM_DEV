import { takeLatest, put, call } from 'redux-saga/effects';
import { 
  searchOrderSuccess, searchOrderFailure, 
  createReturnSuccess, createReturnFailure, 
  generateSRSuccess, generateSRFailure, 
  confirmReturnSuccess, confirmReturnFailure 
} from '../action';
import { returnOrderAPI } from '../../../services/api/order.api';
import { logger } from '../../../utils/logger';

/**
 * Saga สำหรับค้นหา Order
 */
export function* searchOrderSaga(action: any): Generator {
  try {
    logger.log('info', '[ReturnOrder] Searching Order:', action.payload);
    const response = yield call(returnOrderAPI.searchOrder, action.payload.soNo, action.payload.orderNo);
    yield put(searchOrderSuccess(response));
  } catch (error: any) {
    logger.log('error', '[ReturnOrder] Search Failed:', error);
    yield put(searchOrderFailure(error.message || 'Unknown error occurred.'));
  }
}

/**
 * Saga สำหรับสร้าง Return Order
 */
export function* createReturnSaga(action: any): Generator {
  try {
    logger.log('info', '[ReturnOrder] Creating Return Order:', action.payload);
    const response = yield call(returnOrderAPI.createReturnOrder, action.payload);
    yield put(createReturnSuccess(response));
  } catch (error: any) {
    logger.log('error', '[ReturnOrder] Create Failed:', error);
    yield put(createReturnFailure(error.message || 'Unknown error occurred.'));
  }
}

/**
 * Saga สำหรับ Generate SR Number
 */
export function* generateSRSaga(action: any): Generator {
  try {
    logger.log('info', '[ReturnOrder] Generating SR for Order No:', action.payload.orderNo);
    const response = yield call(returnOrderAPI.generateSR, action.payload.orderNo);
    yield put(generateSRSuccess(response));
  } catch (error: any) {
    logger.log('error', '[ReturnOrder] Generate SR Failed:', error);
    yield put(generateSRFailure(error.message || 'Unknown error occurred.'));
  }
}

/**
 * Saga สำหรับ Confirm Return Order
 */
export function* confirmReturnSaga(action: any): Generator {
  try {
    logger.log('info', '[ReturnOrder] Confirming Return Order:', action.payload);
    yield call(returnOrderAPI.confirmReturnOrder, action.payload);
    yield put(confirmReturnSuccess());
  } catch (error: any) {
    logger.log('error', '[ReturnOrder] Confirm Return Order Failed:', error);
    yield put(confirmReturnFailure(error.message || 'Unknown error occurred.'));
  }
}
