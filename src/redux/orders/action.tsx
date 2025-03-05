import { ReturnOrderActionTypes } from './types';
import { CreateBeforeReturnOrderRequest, CreateSRRequest } from './api';

// 1. ค้นหา Order
export const searchOrder = (searchParams: { soNo?: string; orderNo?: string }) => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_SEARCH_REQ,
  payload: searchParams
});

// 2. สร้าง Return Order ด้วย type ที่ถูกต้อง Order Details)
export const createReturnOrder = (data: CreateBeforeReturnOrderRequest) => {
  return {
    type: ReturnOrderActionTypes.RETURN_ORDER_CREATE_REQ,
    payload: data
  };
};

// 3. สร้าง SR Number (SR)
export const createSrNo = (payload: CreateSRRequest) => {
  return {
    type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_REQ,
    payload
  };
};

// ปรับปรุง interface สำหรับ confirmReturn
export interface ConfirmReturnRequest {
  orderNo: string;
  roleId: number;
  userID: string;
}

// ปรับปรุง confirmReturn action
export const confirmReturn = (data: ConfirmReturnRequest) => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_REQ,
  payload: data
});

// เพิ่ม reset action
export const resetReturnOrder = () => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_RESET
});

export const setCurrentStep = (step: 'search' | 'create' | 'sr' | 'preview' | 'confirm') => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_SET_STEP,
  payload: step
});
