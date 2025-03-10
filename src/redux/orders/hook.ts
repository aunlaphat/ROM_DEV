// src/redux/orders/hook.ts
import { useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  searchOrder,
  createReturnOrder,
  generateSr,
  updateSr,
  updateStatus,
  cancelOrder,
  markOrderEdited,
  resetOrder,
  setCurrentStep
} from './action';
import {
  getSearchResult,
  getReturnOrder,
  getCurrentStep,
  getLoading,
  getError,
  getSrCreated,
  getIsEdited,
  getSrNumber,
  getOrderNumber,
  getSoNumber,
  hasReturnItems,
  getTotalReturnAmount,
  getTotalReturnQuantity,
  getOrderData
} from './selector';
import {
  SearchOrderRequest,
  CreateReturnOrderRequest,
  UpdateSrRequest,
  UpdateStatusRequest,
  CancelOrderRequest,
  OrderStep
} from './types';

/**
 * Custom hook สำหรับใช้ฟังก์ชันและข้อมูลต่างๆ ของ Order
 * @returns ฟังก์ชันและข้อมูลที่เกี่ยวข้องกับ Order
 */
export const useOrder = () => {
  const dispatch = useDispatch();
  
  // Selectors
  const orderData = useSelector(getOrderData);
  const searchResult = useSelector(getSearchResult);
  const returnOrder = useSelector(getReturnOrder);
  const currentStep = useSelector(getCurrentStep);
  const loading = useSelector(getLoading);
  const error = useSelector(getError);
  const srCreated = useSelector(getSrCreated);
  const isEdited = useSelector(getIsEdited);
  const srNumber = useSelector(getSrNumber);
  const orderNumber = useSelector(getOrderNumber);
  const soNumber = useSelector(getSoNumber);
  const hasItems = useSelector(hasReturnItems);
  const totalAmount = useSelector(getTotalReturnAmount);
  const totalQuantity = useSelector(getTotalReturnQuantity);
  
  // Actions
  const handleSearchOrder = useCallback((params: SearchOrderRequest) => {
    dispatch(searchOrder(params));
  }, [dispatch]);
  
  const handleCreateReturnOrder = useCallback((data: CreateReturnOrderRequest) => {
    dispatch(createReturnOrder(data));
  }, [dispatch]);
  
  const handleGenerateSr = useCallback((orderNo: string) => {
    dispatch(generateSr(orderNo));
  }, [dispatch]);
  
  const handleUpdateSr = useCallback((data: UpdateSrRequest) => {
    dispatch(updateSr(data));
  }, [dispatch]);
  
  const handleUpdateStatus = useCallback((data: UpdateStatusRequest) => {
    dispatch(updateStatus(data));
  }, [dispatch]);
  
  const handleCancelOrder = useCallback((data: CancelOrderRequest) => {
    dispatch(cancelOrder(data));
  }, [dispatch]);
  
  const handleMarkOrderEdited = useCallback((orderNo: string) => {
    dispatch(markOrderEdited(orderNo));
  }, [dispatch]);
  
  const handleResetOrder = useCallback(() => {
    dispatch(resetOrder());
  }, [dispatch]);
  
  const handleSetStep = useCallback((step: OrderStep) => {
    dispatch(setCurrentStep(step));
  }, [dispatch]);
  
  return {
    // State data
    orderData,
    searchResult,
    returnOrder,
    currentStep,
    loading,
    error,
    srCreated,
    isEdited,
    srNumber,
    orderNumber,
    soNumber,
    hasItems,
    totalAmount,
    totalQuantity,
    
    // Actions
    searchOrder: handleSearchOrder,
    createReturnOrder: handleCreateReturnOrder,
    generateSr: handleGenerateSr,
    updateSr: handleUpdateSr,
    updateStatus: handleUpdateStatus,
    cancelOrder: handleCancelOrder,
    markOrderEdited: handleMarkOrderEdited,
    resetOrder: handleResetOrder,
    setStep: handleSetStep
  };
};