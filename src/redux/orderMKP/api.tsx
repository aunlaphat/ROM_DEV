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
  UPDATESTATUS
} from "../../services/path";
import { calculateReturnStatus } from '../../utils/calculateStatus';
import { STATUS } from "../../constants/returnOrder";

// Request Types
export interface SearchOrderRequest {
  soNo?: string;
  orderNo?: string;
}

export interface CreateBeforeReturnOrderRequest {
  orderNo: string;
  soNo: string;
  channelID: number;
  customerID: string;
  reason: string;
  soStatus?: string;
  mkpStatus?: string;
  warehouseID: number;
  returnDate: string; // ISO date string
  trackingNo: string;
  logistic: string;
  items: CreateBeforeReturnOrderItemRequest[];
}

interface CreateBeforeReturnOrderItemRequest {
  orderNo: string;
  sku: string;
  itemName: string;
  qty: number;
  returnQty: number;
  price: number;
  createBy: string;
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
  currentStep: 'search' | 'create' | 'confirm';
  isEdited: boolean;
  orderLines: any[];
  srCreated: boolean;
}

// เพิ่ม type สำหรับ API calls
type ApiFunction = (...args: any[]) => Promise<AxiosResponse>;

export function* searchOrder(action: { 
  type: ReturnOrderActionTypes; 
  payload: SearchOrderRequest 
}): SagaIterator {
  try {
    openLoading();
    logger.group('Search Order');
    
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

    logger.timeEnd('Search Duration');

    if (!response.data.success) {
      throw new Error(response.data.message || 'Search failed');
    }

    logger.api('info', '✅ Search completed:', response.data);
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_SEARCH_SUCCESS,
      payload: response.data.data // เปลี่ยนเป็น response.data.data
    });

    notification.success({
      message: 'ค้นหาสำเร็จ',
      description: 'พบข้อมูลที่ค้นหา'
    });

  } catch (error: any) {
    logger.critical('Search Order Failed', {
      error: error.message,
      payload: action.payload
    });
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_SEARCH_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'ค้นหาไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    logger.groupEnd();
    yield delay(300);
    closeLoading();
  }
}

export function* createBeforeReturnOrder(action: { 
  type: ReturnOrderActionTypes; 
  payload: CreateBeforeReturnOrderRequest 
}): SagaIterator {
  try {
    openLoading();
    logger.group('Create Return Order');
    logger.api('info', '📝 Creating new return order:', {
      orderNo: action.payload.orderNo,
      soNo: action.payload.soNo
    });
    logger.time('Create Duration');

    const formattedData = {
      ...action.payload,
      returnDate: new Date(action.payload.returnDate).toISOString()
    };

    const response: AxiosResponse<APIResponse<BeforeReturnOrderResponse>> = yield call(
      POST as unknown as ApiFunction, 
      CREATEBEFORERETURNORDER,  // ใช้ path ที่มีอยู่แล้ว
      formattedData
    );
    logger.timeEnd('Create Duration');

    if (!response.data.success) {
      throw new Error(response.data.message);
    }

    logger.api('info', '✅ Order created successfully:', {
      orderNo: response.data.orderNo,
      srNo: response.data.srNo
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
    logger.critical('Create Order Failed', {
      error: error.message,
      orderNo: action.payload.orderNo
    });
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_CREATE_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'สร้างรายการไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    logger.groupEnd();
    yield delay(300);
    closeLoading();
  }
}

export interface CreateSRRequest {
  orderNo: string;
  warehouseFrom: string;
  returnDate: string;
  trackingNo: string;
  transportType: string;
}

export function* updateSR(action: { 
  type: ReturnOrderActionTypes; 
  payload: CreateSRRequest
}): SagaIterator {
  try {
    openLoading();
    logger.group('Update SR');
    logger.api('info', '🔄 Updating SR for order:', action.payload.orderNo);
    logger.time('Update SR Duration');

    const response = yield call(
      POST as unknown as ApiFunction, 
      `${UPDATESR}/${action.payload.orderNo}`,
      action.payload
    );
    logger.timeEnd('Update SR Duration');

    if (!response.data.success) {
      throw new Error(response.data.message);
    }

    logger.api('info', '✅ SR updated successfully:', {
      orderNo: response.data.orderNo,
      newSrNo: response.data.srNo
    });
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_SUCCESS,
      payload: response.data
    });

    notification.success({
      message: 'อัพเดท SR สำเร็จ'
    });

  } catch (error: any) {
    logger.critical('Update SR Failed', {
      error: error.message,
      orderNo: action.payload.orderNo
    });
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_SR_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'อัพเดท SR ไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    logger.groupEnd();
    yield delay(300);
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
    logger.group('Update Return Order Status');
    logger.api('info', '🔄 Calculating status for:', action.payload);

    const { statusReturnID, statusConfID } = calculateReturnStatus({
      roleID: action.payload.roleID,
      isCNCreated: action.payload.isCNCreated,
      isEdited: action.payload.isEdited
    });

    logger.api('info', 'Calculated status:', { statusReturnID, statusConfID });

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

    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_SUCCESS,
      payload: response.data
    });

    notification.success({
      message: 'ยืนยันคำสั่งคืนสินค้าสำเร็จ',
      description: `สถานะถูกอัพเดทเป็น ${statusReturnID === STATUS.RETURN.BOOKING ? 'Booking' : 'Pending'}`
    });

  } catch (error: any) {
    logger.critical('Update status failed:', {
      error: error.message,
      payload: action.payload
    });
    yield put({
      type: ReturnOrderActionTypes.RETURN_ORDER_UPDATE_STATUS_FAIL,
      payload: error.message
    });
    notification.error({
      message: 'อัพเดทสถานะไม่สำเร็จ',
      description: error.response?.data?.message || 'กรุณาลองใหม่อีกครั้ง'
    });
  } finally {
    logger.groupEnd();
    yield delay(300);
    closeLoading();
  }
}

export function* cancelOrder(action: { 
  type: ReturnOrderActionTypes; 
  payload: CancelOrderRequest 
}): SagaIterator {
  try {
    openLoading();
    logger.api('info', 'Cancelling order:', action.payload);

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
    logger.critical('Cancel order failed:', error);
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
    logger.api('info', 'Marking order as edited:', action.payload);

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
    logger.critical('Mark as edited failed', { 
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
