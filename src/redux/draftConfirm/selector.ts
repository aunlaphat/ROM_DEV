// src/redux/draftConfirm/selector.ts
import { createSelector } from 'reselect';
import { RootState } from '../types';
import { DraftConfirmState, Order, OrderItem } from './types';

// Base selector - ดึง draftConfirm state จาก root state
const getDraftConfirmState = (state: RootState): DraftConfirmState => state.draftConfirm;

// ดึงรายการคำสั่งทั้งหมด
export const getOrders = createSelector(
  [getDraftConfirmState],
  (draftConfirmState) => draftConfirmState.orders
);

// ดึงคำสั่งที่เลือก
export const getSelectedOrder = createSelector(
  [getDraftConfirmState],
  (draftConfirmState) => draftConfirmState.selectedOrder
);

// ดึงรายการ CodeR ทั้งหมด
export const getCodeRList = createSelector(
  [getDraftConfirmState],
  (draftConfirmState) => draftConfirmState.codeRList
);

// ดึงสถานะการโหลด
export const getLoading = createSelector(
  [getDraftConfirmState],
  (draftConfirmState) => draftConfirmState.loading
);

// ดึงข้อผิดพลาด
export const getError = createSelector(
  [getDraftConfirmState],
  (draftConfirmState) => draftConfirmState.error
);

// ดึงรายการสินค้าของคำสั่งที่เลือก
export const getSelectedOrderItems = createSelector(
  [getSelectedOrder],
  (selectedOrder): OrderItem[] => selectedOrder?.items || []
);

// คำนวณมูลค่ารวมของคำสั่งที่เลือก
export const getSelectedOrderTotalAmount = createSelector(
  [getSelectedOrderItems],
  (items): number => {
    return items.reduce((total, item) => {
      const qty = item.qty || 0;
      const price = item.price || 0;
      return total + (qty * price);
    }, 0);
  }
);

// คำนวณจำนวนรายการของคำสั่งที่เลือก
export const getSelectedOrderItemCount = createSelector(
  [getSelectedOrderItems],
  (items): number => items.length
);

// ดึง orderNo ของคำสั่งที่เลือก
export const getSelectedOrderNo = createSelector(
  [getSelectedOrder],
  (selectedOrder): string | undefined => selectedOrder?.orderNo
);

// ดึง srNo ของคำสั่งที่เลือก
export const getSelectedOrderSrNo = createSelector(
  [getSelectedOrder],
  (selectedOrder): string | null | undefined => selectedOrder?.srNo
);