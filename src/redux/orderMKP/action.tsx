import { ReturnOrderActionTypes } from './types';
import { CreateBeforeReturnOrderRequest } from './api';

// 1. ค้นหา Order
export const searchOrder = (searchTerm: string) => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_SEARCH_REQ,
  payload: searchTerm
});

// 2. สร้าง Return Order ด้วย type ที่ถูกต้อง
export const createReturnOrder = (data: CreateBeforeReturnOrderRequest) => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_CREATE_REQ,
  payload: data
});

// 3. สร้าง SR Number
export const createSrNo = (orderNo: string) => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_REQ,
  payload: orderNo
});

// 4. ยืนยันการคืนสินค้า (update status)
export const confirmReturn = (data: {
  orderNo: string;
  roleId: number;
  isCNCreated?: boolean;
  isEdited?: boolean;
}) => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_REQ,
  payload: data
});

// เพิ่ม reset action
export const resetReturnOrder = () => ({
  type: ReturnOrderActionTypes.RETURN_ORDER_RESET
});
