// src/redux/draftConfirm/hook.ts
import { useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  fetchOrdersRequest,
  fetchOrderDetailsRequest,
  fetchCodeRRequest,
  addItemRequest,
  removeItemRequest,
  confirmDraftOrderRequest,
  cancelOrderRequest,
  clearSelectedOrder
} from './action';
import {
  getOrders,
  getSelectedOrder,
  getCodeRList,
  getLoading,
  getError,
  getSelectedOrderItems,
  getSelectedOrderTotalAmount,
  getSelectedOrderItemCount,
  getSelectedOrderNo,
  getSelectedOrderSrNo
} from './selector';
import {
  FetchOrdersRequest,
  FetchOrderDetailsRequest,
  AddItemRequest,
  RemoveItemRequest,
  ConfirmDraftOrderRequest,
  CancelOrderRequest
} from './types';

/**
 * Custom hook สำหรับใช้ฟังก์ชันและข้อมูลต่างๆ ของ Draft & Confirm
 * @returns ฟังก์ชันและข้อมูลที่เกี่ยวข้องกับ Draft & Confirm
 */
export const useDraftConfirm = () => {
  const dispatch = useDispatch();
  
  // Selectors
  const orders = useSelector(getOrders);
  const selectedOrder = useSelector(getSelectedOrder);
  const codeRList = useSelector(getCodeRList);
  const loading = useSelector(getLoading);
  const error = useSelector(getError);
  const selectedOrderItems = useSelector(getSelectedOrderItems);
  const totalAmount = useSelector(getSelectedOrderTotalAmount);
  const itemCount = useSelector(getSelectedOrderItemCount);
  const orderNo = useSelector(getSelectedOrderNo);
  const srNo = useSelector(getSelectedOrderSrNo);
  
  // Actions
  const fetchOrders = useCallback((params: FetchOrdersRequest) => {
    dispatch(fetchOrdersRequest(params));
  }, [dispatch]);
  
  const fetchOrderDetails = useCallback((params: FetchOrderDetailsRequest) => {
    dispatch(fetchOrderDetailsRequest(params));
  }, [dispatch]);
  
  const fetchCodeR = useCallback(() => {
    dispatch(fetchCodeRRequest());
  }, [dispatch]);
  
  const addItem = useCallback((params: AddItemRequest) => {
    dispatch(addItemRequest(params));
  }, [dispatch]);
  
  const removeItem = useCallback((params: RemoveItemRequest) => {
    dispatch(removeItemRequest(params));
  }, [dispatch]);
  
  const confirmDraftOrder = useCallback((params: ConfirmDraftOrderRequest) => {
    dispatch(confirmDraftOrderRequest(params));
  }, [dispatch]);
  
  const cancelOrder = useCallback((params: CancelOrderRequest) => {
    dispatch(cancelOrderRequest(params));
  }, [dispatch]);
  
  const clearOrder = useCallback(() => {
    dispatch(clearSelectedOrder());
  }, [dispatch]);
  
  return {
    // State data
    orders,
    selectedOrder,
    codeRList,
    loading,
    error,
    selectedOrderItems,
    totalAmount,
    itemCount,
    orderNo,
    srNo,
    
    // Actions
    fetchOrders,
    fetchOrderDetails,
    fetchCodeR,
    addItem,
    removeItem,
    confirmDraftOrder,
    cancelOrder,
    clearOrder
  };
};