// src/redux/draftConfirm/action.ts
import { logger } from '../../utils/logger';
import { AddItemResponse, CodeR, Order, UpdateOrderStatusResponse } from './types/response';
import { DraftConfirmActionTypes } from './types/action';
import { AddItemRequest, CancelOrderRequest, ConfirmDraftOrderRequest, FetchOrderDetailsRequest, FetchOrdersRequest, RemoveItemRequest } from './types/request';

// Fetch Orders
export const fetchOrdersRequest = (params: FetchOrdersRequest) => {
    logger.log('info', '[Action] Fetch Orders Request', { params });
    return {
        type: DraftConfirmActionTypes.FETCH_ORDERS_REQUEST,
        payload: params
    };
};

export const fetchOrdersSuccess = (orders: Order[]) => ({
    type: DraftConfirmActionTypes.FETCH_ORDERS_SUCCESS,
    payload: orders
});

export const fetchOrdersFailure = (error: string) => ({
    type: DraftConfirmActionTypes.FETCH_ORDERS_FAILURE,
    payload: error
});

// Fetch Order Details
export const fetchOrderDetailsRequest = (params: FetchOrderDetailsRequest) => {
    logger.log('info', '[Action] Fetch Order Details Request', {
        orderNo: params.orderNo,
        statusConfID: params.statusConfID
    });
    return {
        type: DraftConfirmActionTypes.FETCH_ORDER_DETAILS_REQUEST,
        payload: params
    };
};

export const fetchOrderDetailsSuccess = (order: Order) => ({
    type: DraftConfirmActionTypes.FETCH_ORDER_DETAILS_SUCCESS,
    payload: order
});

export const fetchOrderDetailsFailure = (error: string) => ({
    type: DraftConfirmActionTypes.FETCH_ORDER_DETAILS_FAILURE,
    payload: error
});

// Fetch CodeR List
export const fetchCodeRRequest = () => {
    logger.log('info', '[Action] Fetch CodeR List Request');
    return {
        type: DraftConfirmActionTypes.FETCH_CODE_R_REQUEST
    };
};

export const fetchCodeRSuccess = (codeRList: CodeR[]) => ({
    type: DraftConfirmActionTypes.FETCH_CODE_R_SUCCESS,
    payload: codeRList
});

export const fetchCodeRFailure = (error: string) => ({
    type: DraftConfirmActionTypes.FETCH_CODE_R_FAILURE,
    payload: error
});

// Add Item to Order
export const addItemRequest = (params: AddItemRequest) => {
    logger.log('info', '[Action] Add Item Request', {
        orderNo: params.orderNo,
        sku: params.sku
    });
    return {
        type: DraftConfirmActionTypes.ADD_ITEM_REQUEST,
        payload: params
    };
};

export const addItemSuccess = (items: AddItemResponse[]) => ({
    type: DraftConfirmActionTypes.ADD_ITEM_SUCCESS,
    payload: items
});

export const addItemFailure = (error: string) => ({
    type: DraftConfirmActionTypes.ADD_ITEM_FAILURE,
    payload: error
});

// Remove Item from Order
export const removeItemRequest = (params: RemoveItemRequest) => {
    logger.log('info', '[Action] Remove Item Request', {
        orderNo: params.orderNo,
        sku: params.sku
    });
    return {
        type: DraftConfirmActionTypes.REMOVE_ITEM_REQUEST,
        payload: params
    };
};

export const removeItemSuccess = () => ({
    type: DraftConfirmActionTypes.REMOVE_ITEM_SUCCESS
});

export const removeItemFailure = (error: string) => ({
    type: DraftConfirmActionTypes.REMOVE_ITEM_FAILURE,
    payload: error
});

// Confirm Draft Order
export const confirmDraftOrderRequest = (params: ConfirmDraftOrderRequest) => {
    logger.log('info', '[Action] Confirm Draft Order Request', {
        orderNo: params.orderNo
    });
    return {
        type: DraftConfirmActionTypes.CONFIRM_DRAFT_ORDER_REQUEST,
        payload: params
    };
};

export const confirmDraftOrderSuccess = (response: UpdateOrderStatusResponse) => ({
    type: DraftConfirmActionTypes.CONFIRM_DRAFT_ORDER_SUCCESS,
    payload: response
});

export const confirmDraftOrderFailure = (error: string) => ({
    type: DraftConfirmActionTypes.CONFIRM_DRAFT_ORDER_FAILURE,
    payload: error
});

// Cancel Order
export const cancelOrderRequest = (params: CancelOrderRequest) => {
    logger.log('info', '[Action] Cancel Order Request', {
        orderNo: params.orderNo,
        cancelReason: params.cancelReason
    });
    return {
        type: DraftConfirmActionTypes.CANCEL_ORDER_REQUEST,
        payload: params
    };
};

export const cancelOrderSuccess = () => ({
    type: DraftConfirmActionTypes.CANCEL_ORDER_SUCCESS
});

export const cancelOrderFailure = (error: string) => ({
    type: DraftConfirmActionTypes.CANCEL_ORDER_FAILURE,
    payload: error
});

// Clear Selected Order
export const clearSelectedOrder = () => {
    logger.log('info', '[Action] Clear Selected Order');
    return {
        type: DraftConfirmActionTypes.CLEAR_SELECTED_ORDER
    };
};