// src/redux/orders/actions.ts
import { 
  OrderActionTypes, 
  SearchOrderRequest, 
  CreateReturnOrderRequest,
  UpdateSrRequest,
  UpdateStatusRequest,
  CancelOrderRequest,
  OrderStep,
  SearchOrderResponse,
  ReturnOrderResponse,
  UpdateSrResponse,
  UpdateStatusResponse,
  CancelOrderResponse
} from './types';
import { logger } from '../../utils/logger';

// ขั้นตอนที่ 6: กำหนด Actions ให้ตรงกับการทำงานของแต่ละ API

// Search Order Actions
export const searchOrder = (searchParams: SearchOrderRequest) => {
  logger.log('info', '[Action] Search Order Request', { params: searchParams });
  return {
    type: OrderActionTypes.SEARCH_ORDER_REQUEST,
    payload: searchParams
  };
};

export const searchOrderSuccess = (data: SearchOrderResponse) => ({
  type: OrderActionTypes.SEARCH_ORDER_SUCCESS,
  payload: data
});

export const searchOrderFailure = (error: string) => ({
  type: OrderActionTypes.SEARCH_ORDER_FAILURE,
  payload: error
});

// Create Return Order Actions
export const createReturnOrder = (data: CreateReturnOrderRequest) => {
  logger.log('info', '[Action] Create Return Order Request', { 
    orderNo: data.orderNo,
    itemCount: data.items.length 
  });
  return {
    type: OrderActionTypes.CREATE_RETURN_ORDER_REQUEST,
    payload: data
  };
};

export const createReturnOrderSuccess = (data: ReturnOrderResponse) => ({
  type: OrderActionTypes.CREATE_RETURN_ORDER_SUCCESS,
  payload: data
});

export const createReturnOrderFailure = (error: string) => ({
  type: OrderActionTypes.CREATE_RETURN_ORDER_FAILURE,
  payload: error
});

// Generate SR Number Actions
export const generateSr = (orderNo: string) => {
  logger.log('info', '[Action] Generate SR Request', { orderNo });
  return {
    type: OrderActionTypes.GENERATE_SR_REQUEST,
    payload: { orderNo }
  };
};

export const generateSrSuccess = (srNo: string) => ({
  type: OrderActionTypes.GENERATE_SR_SUCCESS,
  payload: srNo
});

export const generateSrFailure = (error: string) => ({
  type: OrderActionTypes.GENERATE_SR_FAILURE,
  payload: error
});

// Update SR Number Actions
export const updateSr = (data: UpdateSrRequest) => {
  logger.log('info', '[Action] Update SR Request', { 
    orderNo: data.orderNo,
    srNo: data.srNo 
  });
  return {
    type: OrderActionTypes.UPDATE_SR_REQUEST,
    payload: data
  };
};

export const updateSrSuccess = (data: UpdateSrResponse) => ({
  type: OrderActionTypes.UPDATE_SR_SUCCESS,
  payload: data
});

export const updateSrFailure = (error: string) => ({
  type: OrderActionTypes.UPDATE_SR_FAILURE,
  payload: error
});

// Update Order Status Actions
export const updateStatus = (data: UpdateStatusRequest) => {
  logger.log('info', '[Action] Update Status Request', { 
    orderNo: data.orderNo,
    roleID: data.roleID,
    userID: data.userID 
  });
  return {
    type: OrderActionTypes.UPDATE_STATUS_REQUEST,
    payload: data
  };
};

export const updateStatusSuccess = (data: UpdateStatusResponse) => ({
  type: OrderActionTypes.UPDATE_STATUS_SUCCESS,
  payload: data
});

export const updateStatusFailure = (error: string) => ({
  type: OrderActionTypes.UPDATE_STATUS_FAILURE,
  payload: error
});

// Cancel Order Actions
export const cancelOrder = (data: CancelOrderRequest) => {
  logger.log('info', '[Action] Cancel Order Request', { refID: data.refID });
  return {
    type: OrderActionTypes.CANCEL_ORDER_REQUEST,
    payload: data
  };
};

export const cancelOrderSuccess = (data: CancelOrderResponse) => ({
  type: OrderActionTypes.CANCEL_ORDER_SUCCESS,
  payload: data
});

export const cancelOrderFailure = (error: string) => ({
  type: OrderActionTypes.CANCEL_ORDER_FAILURE,
  payload: error
});

// Mark Order as Edited Actions
export const markOrderEdited = (orderNo: string) => {
  logger.log('info', '[Action] Mark Order as Edited Request', { orderNo });
  return {
    type: OrderActionTypes.MARK_EDITED_REQUEST,
    payload: orderNo
  };
};

export const markOrderEditedSuccess = () => ({
  type: OrderActionTypes.MARK_EDITED_SUCCESS
});

export const markOrderEditedFailure = (error: string) => ({
  type: OrderActionTypes.MARK_EDITED_FAILURE,
  payload: error
});

// Set Current Step Action
export const setCurrentStep = (step: OrderStep) => {
  logger.log('info', '[Action] Set Current Step', { step });
  return {
    type: OrderActionTypes.SET_CURRENT_STEP,
    payload: step
  };
};

// Reset Order Action
export const resetOrder = () => {
  logger.log('info', '[Action] Reset Order State');
  return {
    type: OrderActionTypes.RESET_ORDER
  };
};

// Clear Error Action
export const clearError = () => ({
  type: OrderActionTypes.CLEAR_ERROR
});