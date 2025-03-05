import React, { useEffect, useState } from "react";
import { Layout, Button, Form, Row, Col, Input, Alert, Modal, message, Spin, notification } from "antd";
import { LeftOutlined } from "@ant-design/icons";
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from "react-router-dom";
import { searchOrder, createSrNo, confirmReturn, createReturnOrder, setCurrentStep } from '../../../redux/orders/action';
import { RootState } from "../../../redux/store";
import { CreateBeforeReturnOrderRequest, generateSrNo, ReturnOrderState } from '../../../redux/orders/api';
import ReturnOrderForm from "./components/ReturnOrderForm";
import { useAuth } from "../../../hooks/useAuth";

const { Content } = Layout;

const CreateReturnOrderMKP = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const { loading, error, orderData, currentStep, returnOrder } = useSelector((state: RootState) => state.returnOrder as ReturnOrderState);
  const [selectedSalesOrder, setSelectedSalesOrder] = useState('');
  const [isChecked, setIsChecked] = useState(false);
  const [returnItems, setReturnItems] = useState<{[key: string]: number}>({});
  const auth = useAuth();
  const [stepLoading, setStepLoading] = useState(false);

  // Handler functions
  const validateStepTransition = (fromStep: string, toStep: string): boolean => {
    switch (toStep) {
      case 'create':
        return !!orderData;
      
      case 'sr':
        // แก้ไขเงื่อนไขการตรวจสอบสำหรับ sr step
        return !!returnOrder;
      
      case 'preview':
        // ตรวจสอบว่ามี SR Number และมาจาก step sr
        return fromStep === 'sr' && !!orderData?.head.srNo;
      
      case 'confirm':
        return fromStep === 'preview' && !!orderData?.head.srNo;
      
      default:
        return true;
    }
  };

  const handleBack = () => {
    const steps = ['search', 'create', 'sr', 'preview', 'confirm'];
    const currentIndex = steps.indexOf(currentStep);
    const prevStep = steps[currentIndex - 1];

    if (!prevStep) {
      navigate("/home");
      return;
    }

    if (validateStepTransition(currentStep, prevStep)) {
      dispatch(setCurrentStep(prevStep as any));
      if (prevStep === 'search') {
        form.resetFields();
        setSelectedSalesOrder("");
        setIsChecked(false);
      }
    } else {
      notification.warning({
        message: 'ไม่สามารถย้อนกลับได้',
        description: 'กรุณาตรวจสอบข้อมูลให้ครบถ้วน'
      });
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedSalesOrder(e.target.value.trim());
  };
  const handleSearch = async () => {
    if (!selectedSalesOrder) {
      message.error("กรุณากรอกเลข SO/Order");
      return;
    }

    setStepLoading(true);
    try {
      const isSoNo = selectedSalesOrder.startsWith("SO");
      const searchPayload = {
        [isSoNo ? "soNo" : "orderNo"]: selectedSalesOrder.trim(),
      };

      await dispatch(searchOrder(searchPayload));
      
      if (validateStepTransition('search', 'create')) {
        dispatch(setCurrentStep("create"));
      }
    } finally {
      setStepLoading(false);
    }
  };

  useEffect(() => {
    if (orderData?.lines) {
      initializeReturnItems(orderData.lines);
    }
  }, [orderData]);

  useEffect(() => {
    if (returnOrder) {
      dispatch(setCurrentStep("sr"));
    }
  }, [returnOrder, dispatch]);

  const handleCancel = () => {
    form.resetFields();
    setSelectedSalesOrder("");
    setIsChecked(false);
  };
  const handleCreateReturnOrder = async () => {
    try {
      if (!orderData?.head.orderNo) {
        message.error("ไม่พบเลขที่ Order");
        return;
      }

      const formValues = form.getFieldsValue();
      const returnItemsList = orderData.lines
        .filter((item) => getReturnQty(item.sku) > 0)
        .map((item) => ({
          orderNo: orderData.head.orderNo,
          sku: item.sku,
          itemName: item.itemName,
          qty: Math.abs(item.qty),
          returnQty: getReturnQty(item.sku),
          price: Math.abs(item.price),
          trackingNo: formValues.trackingNo,
        }));

      if (returnItemsList.length === 0) {
        message.error("กรุณาระบุจำนวนสินค้าที่ต้องการคืน");
        return;
      }

      const warehouseID = Number(formValues.warehouseFrom);
      if (isNaN(warehouseID)) {
        message.error("กรุณาเลือกคลังสินค้าที่ถูกต้อง");
        return;
      }

      const createReturnPayload: CreateBeforeReturnOrderRequest & { success: boolean; message: string } = {
        success: true,
        message: "Return order created successfully",
        orderNo: orderData.head.orderNo,
        soNo: orderData.head.soNo,
        channelID: auth.channelID || 1,
        customerID: auth.customerID || "Customer-002",
        reason: formValues.reason || "Return",
        warehouseID: warehouseID,
        returnDate: formValues.returnDate.toISOString(),
        trackingNo: formValues.trackingNo,
        logistic: formValues.transportType,
        soStatus: orderData.head.salesStatus,
        mkpStatus: orderData.head.mkpStatus,
        items: returnItemsList,
      };

      setStepLoading(true);
      Modal.confirm({
        title: "ยืนยันการสร้างคำสั่งคืนสินค้า",
        content: (
          <div>
            <p>Oreder No: {orderData.head.orderNo}</p>
            <p>SO No: {orderData.head.soNo}</p>
            <p>จำนวนรายการที่จะคืน: {returnItemsList.length} รายการ</p>
            <p>Tracking No: {formValues.trackingNo}</p>
            <p>ขนส่ง: {formValues.transportType}</p>
            <p>วันที่คืน: {formValues.returnDate.format("DD/MM/YYYY HH:mm")}</p>
          </div>
        ),
        okText: "สร้างคำสั่งคืนสินค้า",
        cancelText: "ยกเลิก",
        onOk: async () => {
          const response = await dispatch(createReturnOrder(createReturnPayload));
          if (response?.payload?.success) {
            notification.success({
              message: "สร้างคำสั่งคืนสินค้าสำเร็จ",
              description: response.payload.message,
            });
            dispatch(setCurrentStep("sr"));
          }
          setStepLoading(false);
        },
        onCancel: () => {
          setStepLoading(false);
        }
      });
    } catch (error: any) {
      setStepLoading(false);
      notification.error({
        message: "เกิดข้อผิดพลาด",
        description: error.message,
      });
    }
  };

  const handleCreateSr = async () => {
    try {
      if (!orderData?.head.orderNo) {
        message.error("ไม่พบเลขที่ Order");
        return;
      }

      setStepLoading(true);
      const formValues = form.getFieldsValue();

      // 1. สร้าง SR Number
      const srNo = await generateSrNo(orderData.head.orderNo);

      // 2. สร้าง payload สำหรับอัพเดต SR
      const createSrPayload = {
        orderNo: orderData.head.orderNo,
        warehouseFrom: formValues.warehouseFrom,
        returnDate: formValues.returnDate.toISOString(),
        trackingNo: formValues.trackingNo,
        transportType: formValues.transportType,
        srNo: srNo,
      };

      // 3. ส่ง action เพื่ออัพเดต SR
      const response = await dispatch(createSrNo(createSrPayload));
      
      // 4. ตรวจสอบผลลัพธ์
      if (response?.payload?.srNo) {
        notification.success({
          message: "สร้างเลข SR สำเร็จ",
          description: `SR Number: ${response.payload.srNo}`,
        });
        // ไม่ต้องเปลี่ยน step ทันที ให้ผู้ใช้กดปุ่ม Next เอง
      }
    } catch (error: any) {
      notification.error({
        message: "เกิดข้อผิดพลาด",
        description: error.message,
      });
    } finally {
      setStepLoading(false);
    }
  };

  // ปรับปรุง function ตรวจสอบการ disable ปุ่ม Create Return Order
  const isCreateReturnOrderDisabled = (): boolean => {
    // 1. ตรวจสอบว่ามีข้อมูล Order หรือไม่
    if (!orderData?.head?.orderNo) return true;

    // 2. ตรวจสอบว่ามีการเลือกสินค้าที่จะคืนหรือไม่
    const hasSelectedItems = Object.values(returnItems).some(qty => qty > 0);
    if (!hasSelectedItems) return true;

    // 3. ตรวจสอบว่ากรอกข้อมูลจำเป็นครบถ้วนหรือไม่
    const formValues = form.getFieldsValue();
    const requiredFields = [
      'warehouseFrom',
      'returnDate',
      'trackingNo',
      'transportType'
    ];
    
    const hasAllRequiredFields = requiredFields.every(field => {
      const value = formValues[field];
      return value !== undefined && value !== null && value !== '';
    });

    // 4. ตรวจสอบว่ามี SR Number แล้วหรือไม่
    if (orderData.head.srNo) return true;

    // 5. ตรวจสอบสถานะ loading
    if (loading || stepLoading) return true;

    // คืนค่า false ถ้าผ่านทุกเงื่อนไข (สามารถกดปุ่มได้)
    return !(hasSelectedItems && hasAllRequiredFields);
  };

  // เพิ่มฟังก์ชันเช็คการ disable ปุ่ม Create SR
  const isCreateSRDisabled = (): boolean => {
    // 1. ตรวจสอบว่ามี returnOrder หรือไม่
    if (!returnOrder) return true;

    // 2. ตรวจสอบว่ามี SR Number แล้วหรือยัง
    if (orderData?.head.srNo) return true;

    // 3. ตรวจสอบสถานะ loading
    if (loading || stepLoading) return true;

    return false;
  };

  const initializeReturnItems = (items: any[]) => {
    const initialQty = items.reduce(
      (acc, item) => ({
        ...acc,
        [item.sku]: 0,
      }),
      {}
    );
    setReturnItems(initialQty);
  };

  const getReturnQty = (sku: string): number => {
    return returnItems[sku] || 0;
  };

  const updateReturnQty = (sku: string, change: number) => {
    const currentQty = getReturnQty(sku);
    const originalQty = Math.abs(
      orderData?.lines.find((item) => item.sku === sku)?.qty || 0
    );
    const newQty = Math.max(0, Math.min(originalQty, currentQty + change));

    setReturnItems((prev) => ({
      ...prev,
      [sku]: newQty,
    }));
  };

  const renderBackButton = () => {
    let buttonText = "Back";
    let buttonIcon = <LeftOutlined style={{ color: "#fff", marginRight: 5 }} />;

    if (currentStep === "create") {
      buttonText = "Back to Search";
    } else if (currentStep === "sr") {
      buttonText = "Back to Create";
    } else if (currentStep === "confirm") {
      buttonText = "Back to SR";
    }

    return (
      <Button
        onClick={handleBack}
        style={{ background: "#98CEFF", color: "#fff" }}
        disabled={
          loading || (currentStep === "confirm" && !!orderData?.head.srNo)
        }
      >
        {buttonIcon}
        {buttonText}
      </Button>
    );
  };

  const getStepStatus = (stepKey: string) => {
    switch (stepKey) {
      case "search":
        return currentStep === "search" ? "process" : "finish";
      case "create":
        if (currentStep === "search") return "wait";
        return currentStep === "create" ? "process" : "finish";
      case "sr":
        if (currentStep === "search" || currentStep === "create") return "wait";
        return currentStep === "sr" ? "process" : "finish";
      case "preview":
        if (currentStep === "search" || currentStep === "create" || currentStep === "sr") return "wait";
        return currentStep === "preview" ? "process" : "finish";
      case "confirm":
        if (!orderData?.head.srNo) return "wait";
        return currentStep === "confirm" ? "process" : "finish";
      default:
        return "wait";
    }
  };

  const handleNext = async () => {
    setStepLoading(true);
    try {
      const steps = ['search', 'create', 'sr', 'preview', 'confirm'];
      const currentIndex = steps.indexOf(currentStep);
      const nextStep = steps[currentIndex + 1];

      if (!nextStep) return;

      if (validateStepTransition(currentStep, nextStep)) {
        await dispatch(setCurrentStep(nextStep as any));
      } else {
        notification.warning({
          message: 'ไม่สามารถดำเนินการต่อได้',
          description: 'กรุณาตรวจสอบข้อมูลให้ครบถ้วน'
        });
      }
    } finally {
      setStepLoading(false);
    }
  };

  const handleConfirm = () => {
    try {
      if (!orderData?.head.orderNo) {
        message.error("ไม่พบเลขที่ Order");
        return;
      }

      if (!auth.userID) {
        message.error("ไม่พบข้อมูลผู้ใช้งาน กรุณาเข้าสู่ระบบใหม่");
        return;
      }

      Modal.confirm({
        title: "ยืนยันคำสั่งคืนสินค้า",
        content: (
          <div>
            <p>คุณต้องการยืนยันคำสั่งคืนสินค้าใช่หรือไม่?</p>
            <p>Order No: {orderData.head.orderNo}</p>
            <p>SR No: {orderData.head.srNo}</p>
            <p style={{ color: '#1890ff' }}>
              หมายเหตุ: สถานะจะถูกอัพเดตตามสิทธิ์การใช้งานของคุณ ({auth.roleID === 2 ? 'Accounting' : auth.roleID === 3 ? 'Warehouse' : 'Staff'})
            </p>
          </div>
        ),
        okText: "ยืนยัน",
        cancelText: "ยกเลิก",
        onOk: () => {
          const confirmPayload = {
            orderNo: orderData.head.orderNo,
            roleId: auth.roleID,
            userID: auth.userID,
          };

          console.log('Confirm payload:', confirmPayload);
          
          dispatch(confirmReturn(confirmPayload));

          message.loading({
            content: 'กำลังอัพเดตสถานะ...',
            key: 'confirmStatus',
            duration: 0
          });
        }
      });
    } catch (error: any) {
      notification.error({
        message: "เกิดข้อผิดพลาด",
        description: error.message
      });
    }
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
      isCreateReturnOrderDisabled={isCreateReturnOrderDisabled}
      getStepStatus={getStepStatus}
      renderBackButton={renderBackButton}
      handleNext={handleNext}
      returnItems={returnItems}
      handleConfirm={handleConfirm}
      validateStepTransition={validateStepTransition}
      stepLoading={stepLoading}
      isCreateSRDisabled={isCreateSRDisabled} // เพิ่ม prop ใหม่
    />
  );
};

export default CreateReturnOrderMKP;
