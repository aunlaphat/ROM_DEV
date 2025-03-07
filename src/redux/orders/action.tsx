import { ReturnOrderActionTypes, SearchOrderRequest, SearchOrderResponse } from '../../types/order.types';

/**
 * ค้นหา Order จาก SO No หรือ Order No
 */
export const searchOrderRequest = (payload: SearchOrderRequest) => ({
  type: ReturnOrderActionTypes.SEARCH_ORDER_REQUEST,
  payload,
});

export const searchOrderSuccess = (response: SearchOrderResponse) => ({
  type: ReturnOrderActionTypes.SEARCH_ORDER_SUCCESS,
  payload: response,
});

export const searchOrderFailure = (error: string) => ({
  type: ReturnOrderActionTypes.SEARCH_ORDER_FAILURE,
  payload: error,
});

/**
 * สร้าง Return Order
 */
export const createReturnRequest = (data: any) => ({
  type: ReturnOrderActionTypes.CREATE_RETURN_REQUEST,
  payload: data,
});

export const createReturnSuccess = (response: any) => ({
  type: ReturnOrderActionTypes.CREATE_RETURN_SUCCESS,
  payload: response,
});

export const createReturnFailure = (error: string) => ({
  type: ReturnOrderActionTypes.CREATE_RETURN_FAILURE,
  payload: error,
});

/**
 * Generate SR Number
 */
export const generateSRRequest = (orderNo: string) => ({
  type: ReturnOrderActionTypes.GENERATE_SR_REQUEST,
  payload: { orderNo },
});

export const generateSRSuccess = (response: any) => ({
  type: ReturnOrderActionTypes.GENERATE_SR_SUCCESS,
  payload: response,
});

export const generateSRFailure = (error: string) => ({
  type: ReturnOrderActionTypes.GENERATE_SR_FAILURE,
  payload: error,
});

/**
 * Confirm Return Order
 */
export const confirmReturnRequest = (data: any) => ({
  type: ReturnOrderActionTypes.CONFIRM_RETURN_REQUEST,
  payload: data,
});

export const confirmReturnSuccess = () => ({
  type: ReturnOrderActionTypes.CONFIRM_RETURN_SUCCESS,
});

export const confirmReturnFailure = (error: string) => ({
  type: ReturnOrderActionTypes.CONFIRM_RETURN_FAILURE,
  payload: error,
});

/**
 * จัดการขั้นตอนของ Module
 */
export const setStep = (step: string) => ({
  type: ReturnOrderActionTypes.SET_STEP,
  payload: step,
});

export const resetReturnOrder = () => ({
  type: ReturnOrderActionTypes.RESET,
});
