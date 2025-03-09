// src/redux/orders/selector.ts
import { createSelector } from 'reselect';
import { RootState } from '../types';
import { OrderState, OrderStep } from './types';

// ขั้นตอนที่ 9: สร้าง Selectors สำหรับดึงข้อมูลจาก Redux Store

// Base selector - ดึง order state
const getOrderState = (state: RootState): OrderState => state.order;

export const getOrderData = createSelector(
  [getOrderState],
  (orderState) => orderState.orderData
);

// ดึงผลการค้นหา
export const getSearchResult = createSelector(
  [getOrderState],
  (orderState) => orderState.searchResult
);

// ดึงข้อมูล return order
export const getReturnOrder = createSelector(
  [getOrderState],
  (orderState) => orderState.returnOrder
);

// ดึงขั้นตอนปัจจุบัน
export const getCurrentStep = createSelector(
  [getOrderState],
  (orderState) => orderState.currentStep
);

// ดึงสถานะการโหลด
export const getLoading = createSelector(
  [getOrderState],
  (orderState) => orderState.loading
);

// ดึงข้อผิดพลาด
export const getError = createSelector(
  [getOrderState],
  (orderState) => orderState.error
);

// ดึงสถานะการสร้าง SR
export const getSrCreated = createSelector(
  [getOrderState],
  (orderState) => orderState.srCreated
);

// ดึงสถานะการแก้ไข
export const getIsEdited = createSelector(
  [getOrderState],
  (orderState) => orderState.isEdited
);

// ดึงหมายเลข SR
export const getSrNumber = createSelector(
  [getReturnOrder],
  (returnOrder) => returnOrder?.srNo || null
);

// ดึงหมายเลข Order
export const getOrderNumber = createSelector(
  [getReturnOrder],
  (returnOrder) => returnOrder?.orderNo || null
);

// ดึงหมายเลข SO
export const getSoNumber = createSelector(
  [getReturnOrder],
  (returnOrder) => returnOrder?.soNo || null
);

// ดึงสถานะของแต่ละ step (process, finish, wait)
export const getStepStatus = (step: OrderStep) => createSelector(
  [getCurrentStep],
  (currentStep): 'process' | 'finish' | 'wait' => {
    const steps: OrderStep[] = ['search', 'create', 'sr', 'preview', 'confirm'];
    const currentIndex = steps.indexOf(currentStep);
    const stepIndex = steps.indexOf(step);

    if (stepIndex < currentIndex) return 'finish';
    if (stepIndex === currentIndex) return 'process';
    return 'wait';
  }
);

// ตรวจสอบว่ามีรายการคืนสินค้าอย่างน้อย 1 รายการ
export const hasReturnItems = createSelector(
  [getReturnOrder],
  (returnOrder): boolean => {
    if (!returnOrder || !returnOrder.items) return false;
    return returnOrder.items.some(item => item.returnQty > 0);
  }
);

// คำนวณมูลค่าการคืนสินค้าทั้งหมด
export const getTotalReturnAmount = createSelector(
  [getReturnOrder],
  (returnOrder): number => {
    if (!returnOrder || !returnOrder.items) return 0;
    return returnOrder.items.reduce((total, item) => {
      return total + (item.returnQty * item.price);
    }, 0);
  }
);

// คำนวณจำนวนสินค้าที่คืนทั้งหมด
export const getTotalReturnQuantity = createSelector(
  [getReturnOrder],
  (returnOrder): number => {
    if (!returnOrder || !returnOrder.items) return 0;
    return returnOrder.items.reduce((total, item) => {
      return total + item.returnQty;
    }, 0);
  }
);