// src/redux/orders/saga.ts
import { takeLatest } from 'redux-saga/effects';
import { OrderActionTypes } from './types';
import {
  searchOrderSaga,
  createReturnOrderSaga,
  generateSrSaga,
  updateSrSaga,
  updateStatusSaga,
  cancelOrderSaga,
  markOrderEditedSaga
} from './api';

export default function* orderSaga() {
  yield takeLatest(OrderActionTypes.SEARCH_ORDER_REQUEST, searchOrderSaga);
  yield takeLatest(OrderActionTypes.CREATE_RETURN_ORDER_REQUEST, createReturnOrderSaga);
  yield takeLatest(OrderActionTypes.GENERATE_SR_REQUEST, generateSrSaga);
  yield takeLatest(OrderActionTypes.UPDATE_SR_REQUEST, updateSrSaga);
  yield takeLatest(OrderActionTypes.UPDATE_STATUS_REQUEST, updateStatusSaga);
  yield takeLatest(OrderActionTypes.CANCEL_ORDER_REQUEST, cancelOrderSaga);
  yield takeLatest(OrderActionTypes.MARK_EDITED_REQUEST, markOrderEditedSaga);
}