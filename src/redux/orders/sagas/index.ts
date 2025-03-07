import { all, takeLatest } from 'redux-saga/effects';
import { ReturnOrderActionTypes } from '../../../types/order.types';
import { 
  searchOrderSaga, 
  createReturnSaga, 
  generateSRSaga, 
  confirmReturnSaga 
} from './orderSagas';

/**
 * รวม Saga ทั้งหมดของ Return Order
 */
export default function* returnOrderSaga(): Generator {
  yield all([
    takeLatest(ReturnOrderActionTypes.SEARCH_ORDER_REQUEST, searchOrderSaga),
    takeLatest(ReturnOrderActionTypes.CREATE_RETURN_REQUEST, createReturnSaga),
    takeLatest(ReturnOrderActionTypes.GENERATE_SR_REQUEST, generateSRSaga),
    takeLatest(ReturnOrderActionTypes.CONFIRM_RETURN_REQUEST, confirmReturnSaga),
  ]);
}
