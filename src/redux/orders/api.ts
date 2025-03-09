// src/redux/orders/api.ts
import { call, put, delay } from 'redux-saga/effects';
import { SagaIterator } from 'redux-saga';
import { 
  searchOrderSuccess,
  searchOrderFailure,
  createReturnOrderSuccess,
  createReturnOrderFailure,
  generateSrSuccess,
  generateSrFailure,
  updateSrSuccess,
  updateSrFailure,
  updateStatusSuccess,
  updateStatusFailure,
  cancelOrderSuccess,
  cancelOrderFailure,
  markOrderEditedSuccess,
  markOrderEditedFailure,
  setCurrentStep
} from './action';
import { GET, POST, PATCH } from '../../services';
import { logger } from '../../utils/logger';
import { notification } from 'antd';
import { openLoading, closeLoading } from '../../components/alert/useAlert';
import { 
  SearchOrderRequest, 
  CreateReturnOrderRequest, 
  UpdateSrRequest,
  UpdateStatusRequest,
  CancelOrderRequest,
  OrderActionTypes
} from './types';

// Search Order Saga
export function* searchOrderSaga(action: { type: typeof OrderActionTypes.SEARCH_ORDER_REQUEST; payload: SearchOrderRequest }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Search Order');
    
    // ตรวจสอบว่ามีข้อมูลสำหรับค้นหา
    if (!action.payload.soNo && !action.payload.orderNo) {
      throw new Error('กรุณากรอกเลข SO หรือ Order');
    }
    
    // สร้าง query params
    const queryParams = new URLSearchParams();
    if (action.payload.soNo) {
      queryParams.append('soNo', action.payload.soNo);
    }
    if (action.payload.orderNo) {
      queryParams.append('orderNo', action.payload.orderNo);
    }
    
    // เรียก API endpoint ให้ตรงกับ backend (GET /order/search)
    logger.api.request(`/order/search?${queryParams.toString()}`);
    const response = yield call(() => GET(`/order/search?${queryParams.toString()}`));
    
    if (!response.data.success) {
      throw new Error(response.data.message || 'Search failed');
    }
    
    logger.api.success('/order/search', response.data.data);
    
    // ส่ง action เมื่อค้นหาสำเร็จ
    yield put(searchOrderSuccess(response.data.data));
    
    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'ค้นหาสำเร็จ',
      description: 'พบข้อมูลที่ค้นหา'
    });
    
  } catch (error: any) {
    logger.error('Search Order Error', error);
    
    // ส่ง action เมื่อค้นหาล้มเหลว
    yield put(searchOrderFailure(error.message));
    
    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: 'ค้นหาไม่สำเร็จ',
      description: error.response?.data?.message || error.message
    });
  } finally {
    closeLoading();
    logger.perf.end('Search Order');
  }
}

// Create Return Order Saga
export function* createReturnOrderSaga(action: { type: typeof OrderActionTypes.CREATE_RETURN_ORDER_REQUEST; payload: CreateReturnOrderRequest }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Create Return Order');
    
    // แปลง date ให้อยู่ในรูปแบบที่ถูกต้อง
    const formattedData = {
      ...action.payload,
      returnDate: new Date(action.payload.returnDate).toISOString()
    };
    
    // เรียก API endpoint ให้ตรงกับ backend (POST /order/create)
    logger.api.request('/order/create', {
      orderNo: formattedData.orderNo,
      itemCount: formattedData.items.length
    });
    
    const response = yield call(() => POST('/order/create', formattedData));
    
    if (!response.data.success) {
      throw new Error(response.data.message || 'Create return order failed');
    }
    
    logger.api.success('/order/create', response.data.data);
    
    // ส่ง action เมื่อสร้างสำเร็จ
    yield put(createReturnOrderSuccess(response.data.data));
    
    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'สร้างรายการสำเร็จ',
      description: 'สร้างรายการคืนสินค้าสำเร็จ'
    });
    
  } catch (error: any) {
    logger.error('Create Return Order Error', error);
    
    // ส่ง action เมื่อสร้างล้มเหลว
    yield put(createReturnOrderFailure(error.message));
    
    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: 'สร้างรายการไม่สำเร็จ',
      description: error.response?.data?.message || error.message
    });
  } finally {
    closeLoading();
    logger.perf.end('Create Return Order');
  }
}

// Generate SR Number Saga
export function* generateSrSaga(action: { type: typeof OrderActionTypes.GENERATE_SR_REQUEST; payload: { orderNo: string } }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Generate SR');
    
    const { orderNo } = action.payload;
    
    // เรียก API endpoint ให้ตรงกับ backend (POST /order/generate-sr/:orderNo)
    logger.api.request(`/order/generate-sr/${orderNo}`);
    
    // Fixed: Added empty object as second parameter to POST
    const response = yield call(() => POST(`/order/generate-sr/${orderNo}`, {}));
    
    if (!response.data.success) {
      throw new Error(response.data.message || 'Generate SR failed');
    }
    
    logger.api.success(`/order/generate-sr/${orderNo}`, response.data.data);
    
    // ส่ง action เมื่อสร้าง SR สำเร็จ
    yield put(generateSrSuccess(response.data.data));
    
    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'สร้าง SR สำเร็จ',
      description: `หมายเลข SR: ${response.data.data}`
    });
    
  } catch (error: any) {
    logger.error('Generate SR Error', error);
    
    // ส่ง action เมื่อสร้าง SR ล้มเหลว
    yield put(generateSrFailure(error.message));
    
    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: 'สร้าง SR ไม่สำเร็จ',
      description: error.response?.data?.message || error.message
    });
  } finally {
    closeLoading();
    logger.perf.end('Generate SR');
  }
}

// Update SR Number Saga
export function* updateSrSaga(action: { type: typeof OrderActionTypes.UPDATE_SR_REQUEST; payload: UpdateSrRequest }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Update SR');
    
    const { orderNo, srNo } = action.payload;
    
    // เรียก API endpoint ให้ตรงกับ backend (POST /order/update-sr/:orderNo)
    logger.api.request(`/order/update-sr/${orderNo}`, { srNo });
    
    const response = yield call(() => POST(`/order/update-sr/${orderNo}`, { srNo }));
    
    if (!response.data.success) {
      throw new Error(response.data.message || 'Update SR failed');
    }
    
    logger.api.success(`/order/update-sr/${orderNo}`, response.data.data);
    
    // ส่ง action เมื่ออัพเดท SR สำเร็จ
    yield put(updateSrSuccess(response.data.data));
    
    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'อัพเดท SR สำเร็จ',
      description: `หมายเลข SR: ${srNo} ถูกบันทึกเรียบร้อยแล้ว`
    });
    
    // เปลี่ยนขั้นตอนเป็น preview
    yield put(setCurrentStep('preview'));
    
  } catch (error: any) {
    logger.error('Update SR Error', error);
    
    // ส่ง action เมื่ออัพเดท SR ล้มเหลว
    yield put(updateSrFailure(error.message));
    
    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: 'อัพเดท SR ไม่สำเร็จ',
      description: error.response?.data?.message || error.message
    });
  } finally {
    closeLoading();
    logger.perf.end('Update SR');
  }
}

// Update Order Status Saga
export function* updateStatusSaga(action: { type: typeof OrderActionTypes.UPDATE_STATUS_REQUEST; payload: UpdateStatusRequest }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Update Status');
    
    const { orderNo, roleID, userID } = action.payload;
    
    // เรียก API endpoint ให้ตรงกับ backend (POST /order/update-status/:orderNo)
    logger.api.request(`/order/update-status/${orderNo}`, {
      orderNo,
      roleID,
      userID
    });
    
    const response = yield call(() => POST(`/order/update-status/${orderNo}`, {
      orderNo,
      roleID,
      userID
    }));
    
    if (!response.data.success) {
      throw new Error(response.data.message || 'Update status failed');
    }
    
    logger.api.success(`/order/update-status/${orderNo}`, response.data.data);
    
    // ส่ง action เมื่ออัพเดทสถานะสำเร็จ
    yield put(updateStatusSuccess(response.data.data));
    
    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'ยืนยันคำสั่งคืนสินค้าสำเร็จ',
      description: 'อัพเดทสถานะสำเร็จ'
    });
    
    // เปลี่ยนขั้นตอนเป็น confirm
    yield put(setCurrentStep('confirm'));
    
  } catch (error: any) {
    logger.error('Update Status Error', error);
    
    // ส่ง action เมื่ออัพเดทสถานะล้มเหลว
    yield put(updateStatusFailure(error.message));
    
    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: 'ยืนยันคำสั่งคืนสินค้าไม่สำเร็จ',
      description: error.response?.data?.message || error.message
    });
  } finally {
    closeLoading();
    logger.perf.end('Update Status');
  }
}

// Cancel Order Saga
export function* cancelOrderSaga(action: { type: typeof OrderActionTypes.CANCEL_ORDER_REQUEST; payload: CancelOrderRequest }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Cancel Order');
    
    // เรียก API endpoint ให้ตรงกับ backend (POST /order/cancel)
    logger.api.request('/order/cancel', action.payload);
    
    const response = yield call(() => POST('/order/cancel', action.payload));
    
    if (!response.data.success) {
      throw new Error(response.data.message || 'Cancel order failed');
    }
    
    logger.api.success('/order/cancel', response.data.data);
    
    // ส่ง action เมื่อยกเลิกสำเร็จ
    yield put(cancelOrderSuccess(response.data.data));
    
    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'ยกเลิกรายการสำเร็จ',
      description: 'คำสั่งคืนสินค้าถูกยกเลิกเรียบร้อยแล้ว'
    });
    
  } catch (error: any) {
    logger.error('Cancel Order Error', error);
    
    // ส่ง action เมื่อยกเลิกล้มเหลว
    yield put(cancelOrderFailure(error.message));
    
    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: 'ยกเลิกรายการไม่สำเร็จ',
      description: error.response?.data?.message || error.message
    });
  } finally {
    closeLoading();
    logger.perf.end('Cancel Order');
  }
}

// Mark Order as Edited Saga
export function* markOrderEditedSaga(action: { type: typeof OrderActionTypes.MARK_EDITED_REQUEST; payload: string }): SagaIterator {
  try {
    openLoading();
    logger.perf.start('Mark Order as Edited');
    
    const orderNo = action.payload;
    
    // เรียก API endpoint ให้ตรงกับ backend (PATCH /order/mark-edited/:orderNo)
    logger.api.request(`/order/mark-edited/${orderNo}`);
    
    // Fixed: Added empty object as second parameter to PATCH
    const response = yield call(() => PATCH(`/order/mark-edited/${orderNo}`, {}));
    
    if (!response.data.success) {
      throw new Error(response.data.message || 'Mark as edited failed');
    }
    
    logger.api.success(`/order/mark-edited/${orderNo}`, response.data.data);
    
    // ส่ง action เมื่อทำเครื่องหมายสำเร็จ
    yield put(markOrderEditedSuccess());
    
    // แสดงข้อความแจ้งเตือน
    notification.success({
      message: 'ทำเครื่องหมายการแก้ไขสำเร็จ',
      description: 'รายการถูกทำเครื่องหมายว่ามีการแก้ไขเรียบร้อยแล้ว'
    });
    
  } catch (error: any) {
    logger.error('Mark as Edited Error', error);
    
    // ส่ง action เมื่อทำเครื่องหมายล้มเหลว
    yield put(markOrderEditedFailure(error.message));
    
    // แสดงข้อความแจ้งเตือน
    notification.error({
      message: 'ทำเครื่องหมายการแก้ไขไม่สำเร็จ',
      description: error.response?.data?.message || error.message
    });
  } finally {
    closeLoading();
    logger.perf.end('Mark Order as Edited');
  }
}