import React, { useEffect, useState } from "react";
import {
  Layout,
  Button,
  Form,
  Row,
  Col,
  Input,
  Alert,
  Modal,
  message,
  Steps,
  notification,
} from "antd";
import {
  LeftOutlined,
  SearchOutlined,
  FormOutlined,
  CheckCircleOutlined,
} from "@ant-design/icons";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import {
  searchOrder,
  createSrNo,
  confirmReturn,
  createReturnOrder,
  setCurrentStep,
} from "../../../redux/orders/action"; // เพิ่ม setCurrentStep
import { RootState } from "../../../redux/store";
import {
  CreateBeforeReturnOrderRequest,
  generateSrNo,
  ReturnOrderState,
} from "../../../redux/orders/api";
import ReturnOrderForm from "./components/ReturnOrderForm"; // เพิ่ม import สำหรับ ReturnOrderForm
import { useAuth } from "../../../hooks/useAuth"; // เพิ่ม import สำหรับ AuthContext

const { Content } = Layout;

const CreateReturnOrderMKP = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const { loading, error, orderData, currentStep, returnOrder } = useSelector(
    (state: RootState) => state.returnOrder as ReturnOrderState
  ); // เพิ่ม returnOrder
  const [selectedSalesOrder, setSelectedSalesOrder] = useState("");
  const [isChecked, setIsChecked] = useState(false);
  const [returnItems, setReturnItems] = useState<{ [key: string]: number }>({});
  const auth = useAuth(); // เพิ่ม hook เพื่อดึงข้อมูล auth

  // Handler functions
  const handleBack = () => {
    switch (currentStep) {
      case "create":
        // ถ้าอยู่ในขั้นตอน create ให้กลับไปหน้าค้นหา
        dispatch(setCurrentStep("search"));
        form.resetFields();
        setSelectedSalesOrder("");
        setIsChecked(false);
        break;
      case "sr":
        // ถ้าอยู่ในขั้นตอน sr ให้กลับไปขั้นตอน create
        dispatch(setCurrentStep("create"));
        break;
      case "confirm":
        // ถ้าอยู่ในขั้นตอน confirm ให้กลับไปขั้นตอน sr
        dispatch(setCurrentStep("sr"));
        break;
      default:
        // ถ้าอยู่ในขั้นตอน search ให้กลับไปหน้าก่อนหน้า
        navigate("/home");
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

    const isSoNo = selectedSalesOrder.startsWith("SO");
    const searchPayload = {
      [isSoNo ? "soNo" : "orderNo"]: selectedSalesOrder.trim(),
    };

    dispatch(searchOrder(searchPayload));
    setIsChecked(true);
    dispatch(setCurrentStep("create"));
  };

  // Alternative approach using useEffect
  useEffect(() => {
    if (orderData?.lines) {
      initializeReturnItems(orderData.lines);
    }
  }, [orderData]);

  // เพิ่ม useEffect เพื่อตรวจสอบการสร้าง returnOrder สำเร็จ
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

      // สร้าง items array ตาม interface CreateBeforeReturnOrderItemRequest
      const returnItemsList = orderData.lines
        .filter((item) => getReturnQty(item.sku) > 0)
        .map((item) => ({
          orderNo: orderData.head.orderNo,
          sku: item.sku,
          itemName: item.itemName,
          qty: Math.abs(item.qty),
          returnQty: getReturnQty(item.sku),
          price: Math.abs(item.price),
          trackingNo: formValues.trackingNo, // optional
        }));

      if (returnItemsList.length === 0) {
        message.error("กรุณาระบุจำนวนสินค้าที่ต้องการคืน");
        return;
      }

      // ตรวจสอบค่า warehouseID
      const warehouseID = Number(formValues.warehouseFrom);
      if (isNaN(warehouseID)) {
        message.error("กรุณาเลือกคลังสินค้าที่ถูกต้อง");
        return;
      }

      // สร้าง payload ตาม interface CreateBeforeReturnOrderRequest
      const createReturnPayload: CreateBeforeReturnOrderRequest = {
        orderNo: orderData.head.orderNo,
        soNo: orderData.head.soNo,
        channelID: auth.channelID || 1, // ตรวจสอบค่า channelID
        customerID: auth.customerID || "Customer-002", // ตรวจสอบค่า customerID
        reason: formValues.reason || "Return",
        warehouseID: warehouseID,
        returnDate: formValues.returnDate.toISOString(),
        trackingNo: formValues.trackingNo,
        logistic: formValues.transportType,
        soStatus: orderData.head.salesStatus,
        mkpStatus: orderData.head.mkpStatus,
        items: returnItemsList,
      };

      // แสดง modal ยืนยัน
      Modal.confirm({
        title: "ยืนยันการสร้างคำสั่งคืนสินค้า",
        content: (
          <div>
            <p>จำนวนรายการที่จะคืน: {returnItemsList.length} รายการ</p>
            <p>Tracking No: {formValues.trackingNo}</p>
            <p>ขนส่ง: {formValues.transportType}</p>
            <p>วันที่คืน: {formValues.returnDate.format("DD/MM/YYYY HH:mm")}</p>
          </div>
        ),
        okText: "สร้างคำสั่งคืนสินค้า",
        cancelText: "ยกเลิก",
        onOk: async () => {
          await dispatch(createReturnOrder(createReturnPayload));
          // ไม่ต้อง dispatch setCurrentStep ที่นี่ เพราะ reducer จะจัดการให้
        },
      });
    } catch (error: any) {
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

      const formValues = form.getFieldsValue();

      const createSrPayload = {
        orderNo: orderData.head.orderNo,
        warehouseFrom: formValues.warehouseFrom,
        returnDate: formValues.returnDate.toISOString(),
        trackingNo: formValues.trackingNo,
        transportType: formValues.transportType,
        srNo: await generateSrNo(orderData.head.orderNo),
      };

      dispatch(createSrNo(createSrPayload));
    } catch (error: any) {
      notification.error({
        message: "เกิดข้อผิดพลาด",
        description: error.message,
      });
    }
  };

  // ปรับปรุงฟังก์ชัน helper
  const isCreateReturnOrderDisabled = (): boolean => {
    if (!orderData) return true;
    if (!orderData.head) return true;
    if (loading) return true;
    if (orderData.head.srNo !== null) return true;
    return !validateAdditionalFields();
  };

  // เพิ่มฟังก์ช์ตรวจสอบการกรอกข้อมูลครบถ้วน
  const validateAdditionalFields = (): boolean => {
    const values = form.getFieldsValue();
    return !!(
      values.warehouseFrom &&
      values.returnDate &&
      values.trackingNo &&
      values.transportType
    );
  };

  // Helper functions สำหรับจัดการจำนวนสินค้าที่จะคืน
  const initializeReturnItems = (items: any[]) => {
    const initialQty = items.reduce(
      (acc, item) => ({
        ...acc,
        [item.sku]: 0, // เริ่มต้นเป็น 0 เพื่อให้ผู้ใช้กรอกจำนวนที่ต้องการคืน
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
        } // Convert to boolean
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
        if (orderData?.head.srNo) return "finish";
        return currentStep === "create" ? "process" : "finish";
      case "sr":
        if (currentStep === "search" || currentStep === "create") return "wait";
        if (orderData?.head.srNo) return "finish";
        return currentStep === "sr" ? "process" : "finish";
      case "confirm":
        if (!orderData?.head.srNo) return "wait";
        return currentStep === "confirm" ? "process" : "finish";
      default:
        return "wait";
    }
  };

  const handleNext = () => {
    // เพิ่ม logging เพื่อดู step ปัจจุบัน
    console.log('Current step:', currentStep);
    
    // เช็คว่าอยู่ที่ step preview แล้วไป confirm
    if (currentStep === 'preview') {
      dispatch(setCurrentStep('confirm'));
    }
  };

  const handleConfirm = () => {
    try {
      if (!orderData?.head.orderNo) {
        message.error("ไม่พบเลขที่ Order");
        return;
      }

      // ตรวจสอบ auth ก่อนดำเนินการ
      if (!auth.userID) { // เปลี่ยนจาก auth.user?.userID เป็น auth.userID
        message.error("ไม่พบข้อมูลผู้ใช้งาน กรุณาเข้าสู่ระบบใหม่");
        return;
      }

      // แสดง modal ยืนยันก่อนที่จะ confirm
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
          // ส่งข้อมูลที่จำเป็นสำหรับการอัพเดตสถานะ
          const confirmPayload = {
            orderNo: orderData.head.orderNo,
            roleId: auth.roleID, // เปลี่ยนจาก auth.user?.roleID เป็น auth.roleID
            userID: auth.userID, // เปลี่ยนจาก auth.user.userID เป็น auth.userID
          };

          // log payload เพื่อตรวจสอบ
          console.log('Confirm payload:', confirmPayload);
          
          dispatch(confirmReturn(confirmPayload));

          // แสดง loading message
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
      handleCreateSr={handleCreateSr} // เพิ่ม handleCreateSr
      handleCancel={handleCancel}
      getReturnQty={getReturnQty}
      updateReturnQty={updateReturnQty}
      isCreateReturnOrderDisabled={isCreateReturnOrderDisabled}
      getStepStatus={getStepStatus}
      renderBackButton={renderBackButton}
      handleNext={handleNext}
      returnItems={returnItems}
      handleConfirm={handleConfirm}
    />
  );
};

export default CreateReturnOrderMKP;
