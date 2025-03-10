// src/screens/Orders/Marketplace/index.tsx
import React, { useEffect, useState } from "react";
import { Form, message, notification } from "antd";
import ReturnOrderForm from "./components/ReturnOrderForm";
import { useAuth } from "../../../hooks/useAuth";
import { useOrder } from "../../../hooks/useOrder";

// Custom Hooks
import { useReturnItemsManager } from "./hooks/useReturnItemsManager";
import { useReturnOrderCreate } from "./hooks/useReturnOrderCreate";
import { useReturnOrderSR } from "./hooks/useReturnOrderSR";
import { useReturnOrderConfirm } from "./hooks/useReturnOrderConfirm";
import { useReturnOrderNavigation } from "./hooks/useReturnOrderNavigation";

// Utils
import { isCreateReturnOrderDisabled, isCreateSRDisabled, validateStepTransition as validateStepTransitionUtil } from "./utils/validation";
import { useNavigate } from "react-router-dom";
import { closeLoading, openLoading } from "../../../components/alert/useAlert";

const CreateReturnOrderMKP: React.FC = () => {
  const [form] = Form.useForm();
  const auth = useAuth();
  const navigate = useNavigate();
  
  // สร้าง state เพื่อติดตามสถานะการส่ง action
  const [hasSearched, setHasSearched] = useState(false);
  const [hasCreatedOrder, setHasCreatedOrder] = useState(false);
  const [hasGeneratedSr, setHasGeneratedSr] = useState(false);
  const [hasConfirmedOrder, setHasConfirmedOrder] = useState(false);
  const [stepLoading, setStepLoading] = useState(false);
  const [selectedSalesOrder, setSelectedSalesOrder] = useState('');

  // ใช้ custom hook ที่แก้ไขแล้ว
  const {
    orderData, 
    returnOrder, 
    searchResult,
    currentStep, 
    loading, 
    error,
    searchOrder,
    createReturnOrder,
    updateStatus,
    setStep,
    generateSr
  } = useOrder();

  // Custom hooks สำหรับการจัดการข้อมูล Return Items
  const {
    returnItems,
    setReturnItems,
    getReturnQty,
    updateReturnQty,
    calculateTotalAmount
  } = useReturnItemsManager(orderData);

  // จัดการการเปลี่ยนแปลงข้อมูลช่องค้นหา
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedSalesOrder(e.target.value.trim());
  };

  // จัดการการค้นหา
  const handleSearch = () => {
    if (!selectedSalesOrder) {
      notification.error({ message: "กรุณากรอกเลข SO/Order" });
      return;
    }

    setStepLoading(true);
    try {
      const isSoNo = selectedSalesOrder.startsWith("SO");
      const searchPayload = {
        [isSoNo ? "soNo" : "orderNo"]: selectedSalesOrder.trim(),
      };

      setHasSearched(true);
      searchOrder(searchPayload);
    } catch (error) {
      console.error('[ReturnOrder] Search error:', error);
      setHasSearched(false);
    } finally {
      setStepLoading(false);
    }
  };

  // Custom hook สำหรับการสร้าง Return Order
  const { handleCreateReturnOrder } = useReturnOrderCreate(
    orderData,
    form,
    returnItems,
    createReturnOrder,
    setHasCreatedOrder,
    setStepLoading,
    calculateTotalAmount
  );

  // Custom hook สำหรับการสร้าง SR
  const { handleCreateSr } = useReturnOrderSR(
    orderData,
    form,
    generateSr,
    setHasGeneratedSr,
    setStepLoading
  );

  // Custom hook สำหรับการยืนยัน
  const { handleConfirm } = useReturnOrderConfirm(
    orderData,
    auth,
    updateStatus,
    setHasConfirmedOrder
  );

  // Custom hook สำหรับการนำทาง
  const {
    handleCancel,
    handleNext,
    renderBackButton,
    getStepStatus
  } = useReturnOrderNavigation(
    orderData,
    returnOrder,
    currentStep,
    form,
    loading,
    setStep,
    setSelectedSalesOrder,
    setReturnItems,
    setStepLoading
  );

  // ติดตามการเปลี่ยนแปลงของสถานะจากการค้นหา
  useEffect(() => {
    if (hasSearched && !loading && searchResult) {
      setStep('create');
      setHasSearched(false);
      
      // แสดงแจ้งเตือนเมื่อค้นหาสำเร็จ
      notification.success({
        message: "ค้นหาสำเร็จ",
        description: "พบข้อมูลคำสั่งซื้อ กรุณากรอกข้อมูลสำหรับการคืนสินค้า"
      });
    }
  }, [hasSearched, loading, searchResult, setStep]);

  // ติดตามการเปลี่ยนแปลงของสถานะจากการสร้าง order
  useEffect(() => {
    if (hasCreatedOrder && !loading && returnOrder) {
      setStep('sr');
      setHasCreatedOrder(false);
      
      // แสดงแจ้งเตือนเมื่อสร้าง order สำเร็จ
      notification.success({
        message: "สร้างคำสั่งคืนสินค้าสำเร็จ",
        description: "กรุณาดำเนินการขั้นตอนถัดไป"
      });
    }
  }, [hasCreatedOrder, loading, returnOrder, setStep]);

  // ติดตามการเปลี่ยนแปลงของสถานะจากการสร้าง SR
  useEffect(() => {
    if (hasGeneratedSr && !loading && orderData?.head.srNo) {
      setStep('preview');
      setHasGeneratedSr(false);
      setStepLoading(false); // เพิ่มบรรทัดนี้เพื่อรีเซ็ต loading state
      
      // แสดงแจ้งเตือนเมื่อสร้าง SR สำเร็จ
      notification.success({
        message: "สร้างเลข SR สำเร็จ",
        description: `SR Number: ${orderData.head.srNo}`,
        duration: 5,
      });
    }
  }, [hasGeneratedSr, loading, orderData, setStep, setStepLoading]);

  // ติดตามการเปลี่ยนแปลงของสถานะจากการยืนยัน order
  useEffect(() => {
    if (hasConfirmedOrder && !loading) {
      setHasConfirmedOrder(false);
      setStepLoading(false); // รีเซ็ต loading state
      
      // แสดง global loading (ใช้ openLoading จาก alert)
      openLoading();
      
      // ปิด message loading และแสดงข้อความสำเร็จ
      message.success({
        content: 'อัพเดตสถานะสำเร็จ',
        key: 'confirmStatus', // ใช้ key เดียวกับ message.loading
        duration: 2, // ลดเวลาลง เพื่อให้ redirect เร็วขึ้น
      });
      
      notification.success({
        message: 'อัพเดตสถานะสำเร็จ',
        description: 'การคืนสินค้าถูกดำเนินการเรียบร้อยแล้ว กำลังนำทางไปยังหน้า Draft & Confirm',
        duration: 2,
      });

      // เพิ่ม delay เล็กน้อยก่อน redirect เพื่อให้ผู้ใช้เห็น notification
      setTimeout(() => {
        // ปิด loading ก่อน navigate เพื่อป้องกันปัญหา
        closeLoading();
        // Redirect ไปยัง /draft-and-confirm
        navigate('/draft-and-confirm');
      }, 2000); // delay 2 วินาที
    }
  }, [hasConfirmedOrder, loading, setStepLoading, navigate]);

  // ตรวจสอบความถูกต้องในการเปลี่ยนขั้นตอน (wrapper function)
  const validateStepTransition = (fromStep: string, toStep: string): boolean => {
    return validateStepTransitionUtil(fromStep, toStep, orderData, returnOrder);
  };

  // ตรวจสอบการ disable ปุ่ม Create Return Order (wrapper function)
  const isCreateReturnOrderDisabledWrapper = (): boolean => {
    return isCreateReturnOrderDisabled(orderData, returnItems, form, loading, stepLoading);
  };

  // ตรวจสอบการ disable ปุ่ม Create SR (wrapper function)
  const isCreateSRDisabledWrapper = (): boolean => {
    return isCreateSRDisabled(returnOrder, orderData, loading, stepLoading);
  };

  return (
    <ReturnOrderForm
      currentStep={currentStep}
      orderData={orderData}
      loading={loading}
      error={error}
      form={form}
      selectedSalesOrder={selectedSalesOrder}
      handleInputChange={handleInputChange}
      handleSearch={handleSearch}
      handleCreateReturnOrder={handleCreateReturnOrder}
      handleCreateSr={handleCreateSr}
      handleCancel={handleCancel}
      getReturnQty={getReturnQty}
      updateReturnQty={updateReturnQty}
      isCreateReturnOrderDisabled={isCreateReturnOrderDisabledWrapper}
      getStepStatus={getStepStatus}
      renderBackButton={renderBackButton}
      handleNext={handleNext}
      returnItems={returnItems}
      handleConfirm={handleConfirm}
      validateStepTransition={validateStepTransition}
      stepLoading={stepLoading}
      isCreateSRDisabled={isCreateSRDisabledWrapper}
    />
  );
};

export default CreateReturnOrderMKP;