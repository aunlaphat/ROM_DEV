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

// สร้าง action creator เดียว
export const createSrNo = (payload: CreateSRRequest) => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_REQ,
  payload
});

// ปรับปรุง interface สำหรับ confirmReturn
export interface ConfirmReturnRequest {
  orderNo: string;
  roleId: number;
  userID: string;
}

// แยก interfaces สำหรับ SR
export interface GenerateSRRequest {
  orderNo: string;
}

export interface GenerateSRResponse {
  orderNo: string;
  srNo: string;
}

export interface UpdateSRRequest {
  orderNo: string;
  srNo: string;
  warehouseFrom: string;
  returnDate: string;
  trackingNo: string;
  transportType: string;
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

// แยก actions
export const generateSR = (payload: GenerateSRRequest): {
  type: ReturnOrderActionTypes.RETURN_ORDER_GENERATE_SR_REQ;
  payload: GenerateSRRequest;
} => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_GENERATE_SR_REQ,
  payload
});

export const updateSR = (payload: UpdateSRRequest) => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_REQ,
  payload
});
