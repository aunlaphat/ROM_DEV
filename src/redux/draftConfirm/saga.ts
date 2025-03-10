// src/redux/draftConfirm/saga.ts
import { takeLatest, all } from 'redux-saga/effects';
import { DraftConfirmActionTypes } from './types';
import {
  fetchOrdersSaga,
  fetchOrderDetailsSaga,
  fetchCodeRSaga,
  addItemSaga,
  removeItemSaga,
  confirmDraftOrderSaga,
  cancelOrderSaga
} from './api';

export default function* draftConfirmSaga() {
  yield all([
    takeLatest(DraftConfirmActionTypes.FETCH_ORDERS_REQUEST, fetchOrdersSaga),
    takeLatest(DraftConfirmActionTypes.FETCH_ORDER_DETAILS_REQUEST, fetchOrderDetailsSaga),
    takeLatest(DraftConfirmActionTypes.FETCH_CODE_R_REQUEST, fetchCodeRSaga),
    takeLatest(DraftConfirmActionTypes.ADD_ITEM_REQUEST, addItemSaga),
    takeLatest(DraftConfirmActionTypes.REMOVE_ITEM_REQUEST, removeItemSaga),
    takeLatest(DraftConfirmActionTypes.CONFIRM_DRAFT_ORDER_REQUEST, confirmDraftOrderSaga),
    takeLatest(DraftConfirmActionTypes.CANCEL_ORDER_REQUEST, cancelOrderSaga)
  ]);
}