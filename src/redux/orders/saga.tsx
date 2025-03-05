import { takeLatest } from 'redux-saga/effects';
import { ReturnOrderActionTypes } from './types';
import {
  searchOrder,
  createBeforeReturnOrder, 
  updateSrNo,
  updateStatus,
  cancelOrder,
  markOrderAsEdited
} from './api';

export default function* returnOrderSaga() {
  yield takeLatest(ReturnOrderActionTypes.RETURN_ORDER_SEARCH_REQ, searchOrder);
  yield takeLatest(ReturnOrderActionTypes.RETURN_ORDER_CREATE_REQ, createBeforeReturnOrder);
  yield takeLatest(ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_REQ, updateSrNo);
  yield takeLatest(ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_REQ, updateStatus);
  yield takeLatest(ReturnOrderActionTypes.RETURN_ORDER_CANCEL_REQ, cancelOrder);
  yield takeLatest(ReturnOrderActionTypes.RETURN_ORDER_MARK_EDITED_REQ, markOrderAsEdited);
}
