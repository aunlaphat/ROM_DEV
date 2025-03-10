// src/hooks/useOrder.tsx
import { useCallback } from "react";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "../redux/types";
import {
  searchOrder,
  createReturnOrder,
  generateSr,
  updateSr,
  updateStatus,
  setCurrentStep,
  resetOrder,
} from "../redux/orders/action";
import { logger } from "../utils/logger";

/**
 * Custom hook สำหรับจัดการเกี่ยวกับการทำงานของ Order
 */
export const useOrder = () => {
  const dispatch = useDispatch();

  // เลือกข้อมูลจาก Redux store ให้ตรงกับชื่อใน rootReducer (state.order)
  const orderData = useSelector((state: RootState) => state.order?.orderData || null);
  const searchResult = useSelector((state: RootState) => state.order.searchResult);
  const returnOrder = useSelector((state: RootState) => state.order.returnOrder);
  const currentStep = useSelector((state: RootState) => state.order.currentStep);
  const loading = useSelector((state: RootState) => state.order.loading);
  const error = useSelector((state: RootState) => state.order.error);
  const srCreated = useSelector((state: RootState) => state.order.srCreated);
  const isEdited = useSelector((state: RootState) => state.order.isEdited);

  /**
   * ค้นหา Order
   */
  const searchOrderAction = useCallback(
    (params: any) => {
      try {
        logger.log("info", `[useOrder] Searching order: ${JSON.stringify(params)}`);
        dispatch(searchOrder(params));
      } catch (error) {
        logger.error("[useOrder] Search order error:", error);
        throw error;
      }
    },
    [dispatch]
  );

  /**
   * สร้าง Return Order
   */
  const createReturnOrderAction = useCallback(
    (data: any) => {
      try {
        logger.log("info", `[useOrder] Creating return order for: ${data.orderNo}`);
        dispatch(createReturnOrder(data));
      } catch (error) {
        logger.error("[useOrder] Create return order error:", error);
        throw error;
      }
    },
    [dispatch]
  );

  /**
   * สร้าง SR Number และอัพเดตลงในฐานข้อมูล
   */
  const updateSrAction = useCallback(
    (data: any) => {
      try {
        logger.log("info", `[useOrder] Updating SR for order: ${data.orderNo}, SR: ${data.srNo}`);
        dispatch(updateSr(data));
      } catch (error) {
        logger.error("[useOrder] Update SR error:", error);
        throw error;
      }
    },
    [dispatch]
  );

  /**
   * ยืนยันการคืนสินค้า (อัพเดตสถานะ)
   */
  const updateStatusAction = useCallback(
    (data: any) => {
      try {
        logger.log("info", `[useOrder] Confirming return order: ${data.orderNo}`);
        dispatch(updateStatus(data));
      } catch (error) {
        logger.error("[useOrder] Confirm return error:", error);
        throw error;
      }
    },
    [dispatch]
  );

  /**
   * กำหนดขั้นตอนปัจจุบัน
   */
  const setStepAction = useCallback(
    (step: "search" | "create" | "sr" | "preview" | "confirm") => {
      logger.log("info", `[useOrder] Setting step to: ${step}`);
      dispatch(setCurrentStep(step));
    },
    [dispatch]
  );

  /**
   * รีเซ็ตข้อมูล order
   */
  const resetOrderAction = useCallback(() => {
    logger.log("info", `[useOrder] Resetting order data`);
    dispatch(resetOrder());
  }, [dispatch]);

  /**
   * สร้าง SR Number จากระบบ AX
   */
  const generateSrAction = useCallback((orderNo: string) => {
    try {
      logger.log("info", `[useOrder] Generating SR Number for: ${orderNo}`);
      dispatch(generateSr(orderNo));
    } catch (error) {
      logger.error("[useOrder] Generate SR Number error:", error);
      throw error;
    }
  }, [dispatch]);

  // ส่งค่ากลับ - ชื่อฟังก์ชันตรงกับที่ใช้ในหน้า index.tsx
  return {
    // ข้อมูล
    orderData,
    searchResult,
    returnOrder,
    currentStep,
    loading,
    error,
    srCreated,
    isEdited,

    // ฟังก์ชัน
    searchOrder: searchOrderAction,
    createReturnOrder: createReturnOrderAction,
    updateSr: updateSrAction,
    updateStatus: updateStatusAction,
    setStep: setStepAction,
    resetOrder: resetOrderAction,
    generateSr: generateSrAction,
  };
};