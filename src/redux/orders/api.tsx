import { put, call, Effect } from "redux-saga/effects";
import { GET, POST } from "../../services";
import { ReturnOrderActionTypes } from "./types";
import { logger } from "../../utils/logger";
import { notification } from "antd";
import { delay } from 'redux-saga/effects';
import { openLoading, closeLoading } from "../../components/alert/useAlert";
import { AxiosResponse } from "axios";
import { SagaIterator } from 'redux-saga';
import { 
  SEARCHORDER,
  CREATEBEFORERETURNORDER,
  UPDATESR,
  MARKEDITED,
  CANCELORDER,
  UPDATESTATUS,
  GENERATESR
} from "../../services/path";
import { calculateReturnStatus } from '../../utils/calculateStatus';
import { STATUS } from "../../constants/returnOrder";
import { ReactNode } from "react";

// Request Types
export interface SearchOrderRequest {
  soNo?: string;
  orderNo?: string;
}

export interface CreateBeforeReturnOrderRequest {
  success: any;
  message: ReactNode;
  orderNo: string;
  soNo: string;
  channelID: number;
  customerID: string;
  reason: string;
  soStatus?: string;
  mkpStatus?: string;
  warehouseID: number;
  returnDate: string; // จะถูกแปลงเป็น time.Time ที่ backend
  trackingNo: string;
  logistic: string;
  items: CreateBeforeReturnOrderItemRequest[];
  userID?: string; // Optional เพราะจะถูกเติมจาก JWT ที่ backend
}

export interface CreateBeforeReturnOrderItemRequest {
  orderNo: string;
  sku: string;
  itemName: string;
  qty: number;
  returnQty: number;
  price: number;
  createBy?: string; // เปลี่ยนเป็น optional
  trackingNo?: string;
  alterSKU?: string;
}

export interface CancelOrderRequest {
  refID: string;
  sourceTable: string;
  cancelReason: string;
}

export interface SearchOrderResponse {
  success: boolean;
  message?: string;
  soNo: string;
  orderNo: string;
  statusMKP: string;
  salesStatus: string;
  createDate: string;
  items: SearchOrderItem[];
}

export interface SearchOrderItem {
  sku: string;
  itemName: string;
  qty: number;
  price: number;
}

export interface BeforeReturnOrderResponse {
  orderNo: string;
  soNo: string;
  srNo: string | null;
  channelId: number;
  reason: string;
  customerId: string;
  trackingNo: string;
  logistic: string;
  warehouseId: number;
  soStatus: string | null;
  mkpStatus: string | null;
  returnDate: string | null;
  statusReturnId: number | null;
  statusConfId: number | null;
  confirmBy: string | null;
  confirmDate: string | null;
  createBy: string;
  createDate: string;
  updateBy: string | null;
  updateDate: string | null;
  cancelId: number | null;
  isCNCreated: boolean;
  isEdited: boolean;
  items: BeforeReturnOrderItem[];
}

interface BeforeReturnOrderItem {
  orderNo: string;
  sku: string;
  itemName: string;
  qty: number;
  returnQty: number;
  price: number;
  createBy: string;
  createDate: string;
  trackingNo?: string;
  alterSKU?: string;
}

export interface UpdateSrNoResponse {
  orderNo: string;
  srNo: string;
  statusReturnID?: number;
  statusConfID?: number;
  updateBy: string;
  updateDate: string;
}

export interface UpdateOrderStatusResponse {
  statusReturnID: number;
  statusConfID: number;
  confirmBy: string;
  confirmDate: string;
}

export interface CancelOrderResponse {
  refID: string;
  sourceTable: string;
  cancelReason: string;
  cancelBy: string;
  cancelDate: string;
}

interface APIResponse<T> {
  orderNo: any;
  srNo: any;
  success: boolean;
  message?: string;
  data: T;
}

// เพิ่ม ReturnOrderState interface
export interface ReturnOrderState {
  orderData: {
    head: {
      orderNo: string;
      soNo: string;
      srNo: string | null;
      salesStatus: string;
      mkpStatus: string;
      locationTo: string;
      statusReturnID?: any;
      statusConfID?: any;
      confirmBy?: any;
      confirmDate?: any;
    };
    lines: {
      sku: string;
      itemName: string;
      qty: number;
      price: number;
    }[];
  } | null;
  searchResult: SearchOrderResponse | null;
  returnOrder: BeforeReturnOrderResponse | null;
  loading: boolean;
  error: string | null;
  currentStep: 'search' | 'create' | 'sr' | 'preview' | 'confirm';
  isEdited: boolean;
  orderLines: any[];
  srCreated: boolean;
}

// เพิ่ม type สำหรับ API calls
type ApiFunction = (...args: any[]) => Promise<AxiosResponse>;

// ตัวอย่างการใช้งาน logger ที่ปรับปรุงแล้ว
export function* searchOrder(action: { type: ReturnOrderActionTypes; payload: SearchOrderRequest }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Search Order');
    logger.api.request(SEARCHORDER, action.payload);

    // ตรวจสอบว่ามีการส่งค่าใดค่าหนึ่งมา
    if (!action.payload.soNo && !action.payload.orderNo) {
      throw new Error('กรุณากรอกเลข SO หรือ Order');
    }

    // สร้าง query parameters แยกตาม field
    const queryParams = new URLSearchParams();
    if (action.payload.soNo) {
      queryParams.append('soNo', action.payload.soNo);
    }
    if (action.payload.orderNo) {
      queryParams.append('orderNo', action.payload.orderNo);
    }

    const response: AxiosResponse<APIResponse<SearchOrderResponse>> = yield call(
      GET as unknown as ApiFunction, 
      `${SEARCHORDER}?${queryParams.toString()}`
    );

    if (!response.data.success) {
      throw new Error(response.data.message || 'Search failed');
    }

    logger.api.success(SEARCHORDER, response.data.data);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_SEARCH_SUCCESS,
      payload: response.data.data // เปลี่ยนเป็น response.data.data
    });

    notification.success({
      message: 'ค้นหาสำเร็จ',
      description: 'พบข้อมูลที่ค้นหา'
    });

  } catch (error: any) {
    logger.error('Search Order', error);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_SEARCH_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'ค้นหาไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    logger.perf.end('Search Order');
    closeLoading();
  }
}

export function* createBeforeReturnOrder(action: { type: ReturnOrderActionTypes; payload: CreateBeforeReturnOrderRequest }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Create Return Order');
    logger.api.request(CREATEBEFORERETURNORDER, {
      orderNo: action.payload.orderNo,
      items: action.payload.items.length
    });

    const formattedData = {
      ...action.payload,
      returnDate: new Date(action.payload.returnDate).toISOString()
    };

    const response: AxiosResponse<APIResponse<BeforeReturnOrderResponse>> = yield call(
      POST as unknown as ApiFunction, 
      CREATEBEFORERETURNORDER,  // ใช้ path ที่มีอยู่แล้ว
      formattedData
    );

    if (!response.data.success) {
      throw new Error(response.data.message);
    }

    logger.api.success(CREATEBEFORERETURNORDER, {
      orderNo: response.data.data.orderNo,
      srNo: response.data.data.srNo
    });

    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_CREATE_SUCCESS,
      payload: response.data.data
    });

    notification.success({
      message: 'สร้างรายการสำเร็จ',
      description: 'สร้างรายการคืนสินค้าสำเร็จ'
    });

  } catch (error: any) {
    logger.error('Create Return Order', error);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_CREATE_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'สร้างรายการไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    logger.perf.end('Create Return Order');
    closeLoading();
  }
}

// แก้ไข interface เป็น
export interface CreateSRRequest {
  srNo: any;
  orderNo: string;
  warehouseFrom: string;
  returnDate: string;
  trackingNo: string;
  transportType: string;
  // ลบ srNo ออก เพราะควรถูกสร้างที่ backend
}

export function* updateSrNo(action: { 
  type: ReturnOrderActionTypes; 
  payload: CreateSRRequest;
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Generate SR');
    
    // Generate SR Number
    logger.api.request(GENERATESR, { orderNo: action.payload.orderNo });
    const srResponse = yield call(POST as unknown as ApiFunction, `${GENERATESR}/${action.payload.orderNo}`);
    
    if (!srResponse.data.success) {
      throw new Error('SR Generation failed');
    }
    
    const srNo = srResponse.data.data;
    logger.api.success(GENERATESR, { srNo });

    // Update SR
    logger.api.request(UPDATESR, { ...action.payload, srNo });
    const response = yield call(POST as unknown as ApiFunction, `${UPDATESR}/${action.payload.orderNo}`, { 
      ...action.payload, 
      srNo 
    });

    if (!response.data.success) {
      throw new Error(response.data.message);
    }

    logger.api.success(UPDATESR, response.data.data);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_SUCCESS,
      payload: response.data.data
    });

    notification.success({
      message: 'อัพเดท SR สำเร็จ',
      description: response.data.message
    });

  } catch (error: any) {
    logger.error('SR Update', error);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'อัพเดท SR ไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    logger.perf.end('Generate SR');
    closeLoading();
  }
}

// Step 1: updateStatus - อัพเดทสถานะตาม Role
export interface UpdateStatusRequest {
  orderNo: string;
  statusReturnID: number;
  statusConfID: number;
  userID: string;
}

export function* updateStatus(action: { 
  type: ReturnOrderActionTypes; 
  payload: {
    orderNo: string;
    roleID: number;
    isCNCreated: boolean;
    isEdited: boolean;
    userID: string;
  }
}): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Update Status');
    logger.api.request(UPDATESTATUS, action.payload);

    const { statusReturnID, statusConfID } = calculateReturnStatus({
      roleID: action.payload.roleID,
      isCNCreated: action.payload.isCNCreated,
      isEdited: action.payload.isEdited
    });

    logger.state.update('Calculated Status', { statusReturnID, statusConfID });

    const response = yield call(
      POST as unknown as ApiFunction,
      `${UPDATESTATUS}/${action.payload.orderNo}`,
      {
        orderNo: action.payload.orderNo,
        statusReturnID,
        statusConfID,
        userID: action.payload.userID
      }
    );

    if (!response.data.success) {
      throw new Error(response.data.message);
    }

    logger.api.success(UPDATESTATUS, response.data.data);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_SUCCESS,
      payload: response.data.data
    });

    notification.success({
      message: 'ยืนยันคำสั่งคืนสินค้าสำเร็จ',
      description: `สถานะถูกอัพเดทเป็น ${statusReturnID === STATUS.RETURN.BOOKING ? 'Booking' : 'Pending'}`
    });

  } catch (error: any) {
    logger.error('Status Update', error);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'อัพเดทสถานะไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    logger.perf.end('Update Status');
    closeLoading();
  }
}

export function* cancelOrder(action: { 
  type: ReturnOrderActionTypes; 
  payload: CancelOrderRequest 
}): SagaIterator {
  try {
    openLoading();
    logger.api.request('Cancelling order', action.payload);

    const response = yield call(
      POST as unknown as ApiFunction, 
      CANCELORDER, 
      action.payload
    );

    if (!response.data.success) {
      throw new Error(response.data.message);
    }

    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_CANCEL_SUCCESS,
      payload: response.data
    });

    notification.success({
      message: 'ยกเลิกรายการสำเร็จ'
    });

  } catch (error: any) {
    logger.error('Cancel order failed:', error);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_CANCEL_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'ยกเลิกรายการไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    yield delay(300);
    closeLoading();
  }
}

export function* markOrderAsEdited(action: { 
  type: ReturnOrderActionTypes; 
  payload: string 
}): SagaIterator {
  try {
    openLoading();
    logger.log('info', 'Marking order as edited:', action.payload);

    const response = yield call(
      POST as unknown as ApiFunction, 
      `${MARKEDITED}/${action.payload}`
    );

    if (!response.data.success) {
      throw new Error(response.data.message);
    }

    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_MARK_EDITED_SUCCESS
    });

  } catch (error: any) {
    logger.error('Mark as edited failed', { 
      error: error.message,
      orderNo: action.payload 
    });
    
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_MARK_EDITED_FAIL,
      payload: error.message
    });
  } finally {
    yield delay(300);
    closeLoading();
  }
}

// ฟังก์ชันสำหรับเรียก API เพื่อสร้าง SrNo
export const generateSrNo = async (orderNo: string): Promise<string> => {
  try {
    const response: AxiosResponse<APIResponse<string>> = await POST(`${GENERATESR}/${orderNo}`, {});
    if (response.data.success) {
      return response.data.data;
    } else {
      throw new Error(response.data.message);
    }
  } catch (error) {
    throw new Error('Failed to generate SrNo');
  }
};
