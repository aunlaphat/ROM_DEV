// src/redux/draftConfirm/api.ts
import { call, put } from "redux-saga/effects";
import { SagaIterator } from "redux-saga";
import {
  fetchOrdersSuccess,
  fetchOrdersFailure,
  fetchOrderDetailsSuccess,
  fetchOrderDetailsFailure,
  fetchCodeRSuccess,
  fetchCodeRFailure,
  addItemSuccess,
  addItemFailure,
  removeItemSuccess,
  removeItemFailure,
  confirmDraftOrderSuccess,
  confirmDraftOrderFailure,
  cancelOrderSuccess,
  cancelOrderFailure,
} from "./action";
import { GET, POST, PATCH, DELETE } from "../../services";
import { logger } from "../../utils/logger";
import {
  DraftConfirmActionTypes,
  FetchOrdersRequest,
  FetchOrderDetailsRequest,
  AddItemRequest,
  RemoveItemRequest,
  ConfirmDraftOrderRequest,
  CancelOrderRequest,
} from "./types";
import { openLoading, closeLoading } from "../../components/alert/useAlert";

// Fetch Orders Saga
export function* fetchOrdersSaga(action: {
  type: typeof DraftConfirmActionTypes.FETCH_ORDERS_REQUEST;
  payload: FetchOrdersRequest;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Fetch Orders");

    const { statusConfID, startDate, endDate } = action.payload;
    const query = new URLSearchParams();
    query.append("statusConfID", statusConfID.toString());

    if (startDate) query.append("startDate", startDate);
    if (endDate) query.append("endDate", endDate);

    logger.api.request(`/draft-confirm/orders?${query.toString()}`);
    const response = yield call(() =>
      GET(`/draft-confirm/orders?${query.toString()}`)
    );

    if (!response.data.success) {
      throw new Error(response.data.message || "Fetch orders failed");
    }

    logger.api.success(`/draft-confirm/orders`, response.data.data);
    yield put(fetchOrdersSuccess(response.data.data));
  } catch (error: any) {
    logger.error("Fetch Orders Error", error);
    yield put(fetchOrdersFailure(error.message || "Failed to fetch orders"));
  } finally {
    closeLoading();
    logger.perf.end("Fetch Orders");
  }
}

// Fetch Order Details Saga
export function* fetchOrderDetailsSaga(action: {
  type: typeof DraftConfirmActionTypes.FETCH_ORDER_DETAILS_REQUEST;
  payload: FetchOrderDetailsRequest;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Fetch Order Details");

    const { orderNo, statusConfID } = action.payload;
    const query = new URLSearchParams();
    query.append("statusConfID", statusConfID.toString());
    query.append("orderNo", orderNo);

    logger.api.request(`/draft-confirm/order/details?${query.toString()}`);
    const response = yield call(() =>
      GET(`/draft-confirm/order/details?${query.toString()}`)
    );

    if (!response.data.success) {
      throw new Error(response.data.message || "Fetch order details failed");
    }

    logger.api.success(`/draft-confirm/order/details`, response.data.data);
    yield put(fetchOrderDetailsSuccess(response.data.data));
  } catch (error: any) {
    logger.error("Fetch Order Details Error", error);
    yield put(
      fetchOrderDetailsFailure(error.message || "Failed to fetch order details")
    );
  } finally {
    closeLoading();
    logger.perf.end("Fetch Order Details");
  }
}

// Fetch CodeR List Saga
export function* fetchCodeRSaga(): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Fetch CodeR List");

    logger.api.request("/draft-confirm/list-codeR");
    const response = yield call(() => GET("/draft-confirm/list-codeR"));

    if (!response.data.success) {
      throw new Error(response.data.message || "Fetch CodeR list failed");
    }

    logger.api.success("/draft-confirm/list-codeR", response.data.data);
    yield put(fetchCodeRSuccess(response.data.data));
  } catch (error: any) {
    logger.error("Fetch CodeR List Error", error);
    yield put(fetchCodeRFailure(error.message || "Failed to fetch CodeR list"));
  } finally {
    closeLoading();
    logger.perf.end("Fetch CodeR List");
  }
}

// Add Item Saga
export function* addItemSaga(action: {
  type: typeof DraftConfirmActionTypes.ADD_ITEM_REQUEST;
  payload: AddItemRequest;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Add Item");

    const { orderNo, ...itemData } = action.payload;

    logger.api.request(`/draft-confirm/add-item/${orderNo}`, itemData);
    const response = yield call(() =>
      POST(`/draft-confirm/add-item/${orderNo}`, itemData)
    );

    if (!response.data.success) {
      throw new Error(response.data.message || "Add item failed");
    }

    logger.api.success(
      `/draft-confirm/add-item/${orderNo}`,
      response.data.data
    );
    yield put(addItemSuccess(response.data.data));

    // Re-fetch order details to get updated items
    yield put({
      type: DraftConfirmActionTypes.FETCH_ORDER_DETAILS_REQUEST,
      payload: { orderNo, statusConfID: 1 }, // 1 = Draft status
    });
  } catch (error: any) {
    logger.error("Add Item Error", error);
    yield put(addItemFailure(error.message || "Failed to add item"));
  } finally {
    closeLoading();
    logger.perf.end("Add Item");
  }
}

// Remove Item Saga
export function* removeItemSaga(action: {
  type: typeof DraftConfirmActionTypes.REMOVE_ITEM_REQUEST;
  payload: RemoveItemRequest;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Remove Item");

    const { orderNo, sku } = action.payload;

    logger.api.request(`/draft-confirm/remove-item/${orderNo}/${sku}`);
    const response = yield call(() =>
      DELETE(`/draft-confirm/remove-item/${orderNo}/${sku}`)
    );

    if (!response.data.success) {
      throw new Error(response.data.message || "Remove item failed");
    }

    logger.api.success(`/draft-confirm/remove-item/${orderNo}/${sku}`);
    yield put(removeItemSuccess());

    // Re-fetch order details to get updated items
    yield put({
      type: DraftConfirmActionTypes.FETCH_ORDER_DETAILS_REQUEST,
      payload: { orderNo, statusConfID: 1 }, // 1 = Draft status
    });
  } catch (error: any) {
    logger.error("Remove Item Error", error);
    yield put(removeItemFailure(error.message || "Failed to remove item"));
  } finally {
    closeLoading();
    logger.perf.end("Remove Item");
  }
}

// Confirm Draft Order Saga
export function* confirmDraftOrderSaga(action: {
  type: typeof DraftConfirmActionTypes.CONFIRM_DRAFT_ORDER_REQUEST;
  payload: ConfirmDraftOrderRequest;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Confirm Draft Order');

    const { orderNo } = action.payload;
    
    logger.api.request(`/draft-confirm/update-status/${orderNo}`);
    const response = yield call(() => PATCH(`/draft-confirm/update-status/${orderNo}`, {}));

    if (!response.data.success) {
      throw new Error(response.data.message || 'Confirm draft order failed');
    }

    logger.api.success(`/draft-confirm/update-status/${orderNo}`, response.data.data);
    yield put(confirmDraftOrderSuccess(response.data.data));
    
    // ใช้วันที่ปัจจุบันเป็นค่าเริ่มต้นสำหรับการค้นหา
    const today = new Date();
    const formattedDate = today.toISOString().split('T')[0]; // ได้รูปแบบ 'YYYY-MM-DD'
    
    // Re-fetch orders with updated status และส่งวันที่ไปด้วย
    yield put({
      type: DraftConfirmActionTypes.FETCH_ORDERS_REQUEST,
      payload: { 
        statusConfID: 2, // 2 = Confirm status
        startDate: formattedDate,
        endDate: formattedDate
      }
    });
  } catch (error: any) {
    logger.error('Confirm Draft Order Error', error);
    yield put(confirmDraftOrderFailure(error.message || 'Failed to confirm draft order'));
  } finally {
    closeLoading();
    logger.perf.end('Confirm Draft Order');
  }
}

// Cancel Order Saga
export function* cancelOrderSaga(action: {
  type: typeof DraftConfirmActionTypes.CANCEL_ORDER_REQUEST;
  payload: CancelOrderRequest;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start("Cancel Order");

    const { orderNo, cancelReason } = action.payload;

    logger.api.request("/order/cancel", {
      refID: orderNo,
      sourceTable: "BeforeReturnOrder",
      cancelReason: cancelReason,
    });

    const response = yield call(() =>
      POST("/order/cancel", {
        refID: orderNo,
        sourceTable: "BeforeReturnOrder",
        cancelReason: cancelReason,
      })
    );

    if (!response.data.success) {
      throw new Error(response.data.message || "Cancel order failed");
    }

    logger.api.success("/order/cancel", response.data.data);
    yield put(cancelOrderSuccess());

    const today = new Date();
    const formattedDate = today.toISOString().split("T")[0];

    yield put({
      type: DraftConfirmActionTypes.FETCH_ORDERS_REQUEST,
      payload: {
        statusConfID: 1, // 1 = Draft status
        startDate: formattedDate,
        endDate: formattedDate,
      },
    });
  } catch (error: any) {
    logger.error("Cancel Order Error", error);
    yield put(cancelOrderFailure(error.message || "Failed to cancel order"));
  } finally {
    closeLoading();
    logger.perf.end("Cancel Order");
  }
}
