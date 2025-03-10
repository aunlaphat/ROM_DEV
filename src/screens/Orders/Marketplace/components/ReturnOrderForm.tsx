// src/screens/Orders/Marketplace/components/ReturnOrderForm.tsx
// ตัวอย่างการปรับปรุง ReturnOrderForm.tsx ให้ใช้ layout และ styles ใหม่

import React from "react";
import {
  Layout,
  Button,
  Form,
  Row,
  Col,
  Input,
  Alert,
  Spin,
  Space,
  Divider,
  Typography,
  Tooltip,
  Card,
  Tag,
} from "antd";
import {
  LeftOutlined,
  SearchOutlined,
  ArrowRightOutlined,
  CheckOutlined,
  QuestionCircleOutlined,
  FileTextOutlined,
  NumberOutlined,
  InfoCircleOutlined,
  CheckCircleOutlined,
} from "@ant-design/icons";
import OrderDetailsSection from "./OrderDetailsSection";
import OrderItemsSection from "./OrderItemsSection";
import ReturnOrderSteps from "./ReturnOrderSteps";
import PreviewSection from "./PreviewSection";
import IntegratedPreviewSection from "./IntegratedPreviewSection";
import styles from "../styles/ReturnOrderContainer.module.css"; // นำเข้า CSS module ใหม่

const { Content } = Layout;
const { Title, Text } = Typography;

interface ReturnOrderFormProps {
  currentStep: "search" | "create" | "sr" | "preview" | "confirm";
  orderData: any;
  loading: boolean;
  error: string | null;
  form: any;
  selectedSalesOrder: string;
  handleInputChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleSearch: () => void;
  handleCreateReturnOrder: () => void;
  handleCreateSr: () => void;
  handleCancel: () => void;
  handleNext: () => void;
  handleConfirm: () => void;
  getReturnQty: (sku: string) => number;
  updateReturnQty: (sku: string, change: number) => void;
  isCreateReturnOrderDisabled: () => boolean;
  getStepStatus: (stepKey: string) => "process" | "finish" | "wait";
  renderBackButton: () => {
    onClick: () => void;
    disabled: boolean;
    text: string;
  };
  returnItems: { [key: string]: number };
  validateStepTransition: (fromStep: string, toStep: string) => boolean;
  stepLoading: boolean;
  isCreateSRDisabled: () => boolean;
}

const ReturnOrderForm: React.FC<ReturnOrderFormProps> = ({
  currentStep,
  orderData,
  loading,
  error,
  form,
  selectedSalesOrder,
  handleInputChange,
  handleSearch,
  handleCreateReturnOrder,
  handleCreateSr,
  handleCancel,
  handleNext,
  handleConfirm,
  getReturnQty,
  updateReturnQty,
  isCreateReturnOrderDisabled,
  getStepStatus,
  renderBackButton,
  returnItems,
  validateStepTransition,
  stepLoading = false,
  isCreateSRDisabled,
}) => {
  // ฟังก์ชันแสดงหัวข้อของแต่ละขั้นตอน
  const getStepTitle = (): string => {
    const titles = {
      search: "ค้นหาคำสั่งซื้อ",
      create: "สร้างคำสั่งคืนสินค้า",
      sr: "สร้างเลข SR",
      preview: "ตรวจสอบข้อมูลการคืนสินค้า",
      confirm: "ยืนยันคำสั่งคืนสินค้า",
    };
    return titles[currentStep] || "";
  };

  // ฟังก์ชันแสดงคำอธิบายของแต่ละขั้นตอน
  const getStepDescription = (): string => {
    const descriptions = {
      search: "กรอกเลข SO หรือ Order Number เพื่อค้นหาคำสั่งซื้อ",
      create: "กรอกข้อมูลและเลือกจำนวนสินค้าที่ต้องการคืน",
      sr: orderData?.head.srNo
        ? `เลข SR Number: ${orderData.head.srNo} ได้ถูกสร้างเรียบร้อยแล้ว`
        : 'กดปุ่ม "Create SR" เพื่อสร้างเลข SR Number',
      preview: "ตรวจสอบข้อมูลคำสั่งคืนสินค้าก่อนยืนยัน",
      confirm: "ยืนยันข้อมูลเพื่อเสร็จสิ้นกระบวนการคืนสินค้า",
    };
    return descriptions[currentStep] || "";
  };

  // แสดงปุ่มดำเนินการในแต่ละขั้นตอน
  const renderActionButtons = () => {
    const backBtn = renderBackButton();

    switch (currentStep) {
      // ขั้นตอนค้นหา
      case "search":
        return (
          <Button
            type="primary"
            icon={<SearchOutlined />}
            onClick={handleSearch}
            loading={loading || stepLoading}
            disabled={!selectedSalesOrder}
            size="large"
          >
            ค้นหา
          </Button>
        );

      // ขั้นตอนสร้างคำสั่งคืนสินค้า
      case "create":
        return (
          <Space>
            <Button
              onClick={backBtn.onClick}
              icon={<LeftOutlined />}
              disabled={backBtn.disabled}
            >
              {backBtn.text}
            </Button>
            <Button
              type="primary"
              onClick={handleCreateReturnOrder}
              loading={loading || stepLoading}
              disabled={isCreateReturnOrderDisabled()}
              size="large"
            >
              สร้างคำสั่งคืนสินค้า
            </Button>
            <Button onClick={handleCancel} disabled={loading || stepLoading}>
              ยกเลิก
            </Button>
          </Space>
        );

      // ขั้นตอนสร้าง SR
      case "sr":
        return (
          <Space>
            <Button
              onClick={backBtn.onClick}
              icon={<LeftOutlined />}
              disabled={backBtn.disabled}
            >
              {backBtn.text}
            </Button>

            {orderData?.head.srNo ? (
              // ถ้ามี SR Number แล้ว แสดงปุ่ม Next
              <Button
                type="primary"
                onClick={handleNext}
                loading={loading || stepLoading}
                disabled={!validateStepTransition("sr", "preview")}
                icon={<ArrowRightOutlined />}
                size="large"
              >
                ดำเนินการต่อ
              </Button>
            ) : (
              // ถ้ายังไม่มี SR Number ให้แสดงปุ่ม Create SR
              <Button
                type="primary"
                onClick={handleCreateSr}
                loading={loading || stepLoading}
                disabled={isCreateSRDisabled()}
                size="large"
              >
                สร้างเลข SR
              </Button>
            )}

            <Button onClick={handleCancel} disabled={loading || stepLoading}>
              ยกเลิก
            </Button>
          </Space>
        );

      // ขั้นตอนตรวจสอบข้อมูล
      case "preview":
        return (
          <Space>
            <Button
              onClick={backBtn.onClick}
              icon={<LeftOutlined />}
              disabled={backBtn.disabled}
            >
              {backBtn.text}
            </Button>
            <Button
              type="primary"
              onClick={handleNext}
              loading={loading || stepLoading}
              disabled={!validateStepTransition("preview", "confirm")}
              icon={<ArrowRightOutlined />}
              size="large"
            >
              ดำเนินการต่อ
            </Button>
            <Button onClick={handleCancel} disabled={loading || stepLoading}>
              ยกเลิก
            </Button>
          </Space>
        );

      // ขั้นตอนยืนยัน
      case "confirm":
        return (
          <Space>
            <Button
              onClick={backBtn.onClick}
              icon={<LeftOutlined />}
              disabled={backBtn.disabled}
            >
              {backBtn.text}
            </Button>
            <Button
              type="primary"
              onClick={handleConfirm}
              loading={loading || stepLoading}
              icon={<CheckOutlined />}
              size="large"
            >
              ยืนยันคำสั่งคืนสินค้า
            </Button>
            <Button onClick={handleCancel} disabled={loading || stepLoading}>
              ยกเลิก
            </Button>
          </Space>
        );

      default:
        return null;
    }
  };

  return (
    <Spin
      spinning={loading || stepLoading}
      tip={stepLoading ? "กำลังดำเนินการ..." : "กำลังโหลดข้อมูล..."}
    >
      <div className={styles.returnOrderMainContainer}>
        {/* ส่วนหัวและขั้นตอน */}
        <div className={styles.createReturnHeader}>
          <Title level={3} className={styles.createReturnHeaderTitle}>
            ระบบคืนสินค้า - {getStepTitle()}
          </Title>
          <Text type="secondary" className={styles.createReturnHeaderSubtitle}>
            {getStepDescription()}
          </Text>
        </div>
        
        <div className={styles.stepProgress}>
          <ReturnOrderSteps
            currentStep={currentStep}
            orderData={orderData}
            getStepStatus={getStepStatus}
          />
        </div>

        {/* ขั้นตอนค้นหา */}
        {currentStep === "search" && (
          <div className={`${styles.orderFormCard} ${styles.panelShadow}`}>
            <div className={`${styles.p24} ${styles.textCenter}`}>
              <Space direction="vertical" size="large" style={{ width: '100%', maxWidth: '800px', margin: '0 auto' }}>
                <div>
                  <Title level={2} style={{ color: "#1890ff", marginBottom: 16 }}>
                    <SearchOutlined /> ค้นหาคำสั่งซื้อที่ต้องการคืน
                  </Title>
                  <Text type="secondary" style={{ fontSize: '16px' }}>
                    คุณสามารถค้นหาด้วยเลข SO หรือ Order Number เพื่อเริ่มกระบวนการคืนสินค้า
                  </Text>
                </div>

                {/* ตัวอย่างการค้นหา */}
                <Card 
                  bordered={false} 
                  className={styles.bgLight}
                  style={{ marginBottom: 32 }}
                >
                  <Space direction="vertical" size="small">
                    <Text strong>ตัวอย่างเลขที่ใช้ค้นหา:</Text>
                    <Space>
                      <Tag icon={<FileTextOutlined />} color="blue">SO12345678</Tag>
                      <Text type="secondary">หรือ</Text>
                      <Tag icon={<NumberOutlined />} color="cyan">OR98765432</Tag>
                    </Space>
                  </Space>
                </Card>

                <Form layout="vertical" form={form}>
                  <Form.Item
                    name="selectedSalesOrder"
                    rules={[{ required: true, message: "กรุณากรอกเลข SO หรือ Order!" }]}
                  >
                    <Input.Search
                      size="large"
                      placeholder="กรอกเลข SO หรือ Order Number"
                      value={selectedSalesOrder}
                      onChange={handleInputChange}
                      onSearch={handleSearch}
                      disabled={loading}
                      style={{
                        width: "100%",
                        height: "56px",
                        fontSize: "16px"
                      }}
                      enterButton={
                        <Button 
                          type="primary" 
                          size="large"
                          style={{ height: "56px", width: "120px" }}
                          icon={<SearchOutlined />}
                          loading={loading || stepLoading}
                        >
                          ค้นหา
                        </Button>
                      }
                    />
                  </Form.Item>
                </Form>

                <Card
                  title={
                    <Space>
                      <InfoCircleOutlined style={{ color: "#1890ff" }} />
                      <Text strong>คำแนะนำในการค้นหา</Text>
                    </Space>
                  }
                  bordered={false}
                  className={styles.bgLight}
                >
                  <Row gutter={[16, 16]}>
                    <Col xs={24} md={12}>
                      <Space align="start">
                        <CheckCircleOutlined style={{ color: "#52c41a" }} />
                        <div>
                          <Text strong>เลข SO</Text>
                          <br />
                          <Text type="secondary">ขึ้นต้นด้วย SO ตามด้วยตัวเลข 8 หลัก</Text>
                        </div>
                      </Space>
                    </Col>
                    <Col xs={24} md={12}>
                      <Space align="start">
                        <CheckCircleOutlined style={{ color: "#52c41a" }} />
                        <div>
                          <Text strong>Order Number</Text>
                          <br />
                          <Text type="secondary">ขึ้นต้นด้วย OR ตามด้วยตัวเลข 8 หลัก</Text>
                        </div>
                      </Space>
                    </Col>
                  </Row>
                </Card>

                {error && (
                  <Alert
                    message="ไม่พบข้อมูลที่ค้นหา"
                    description={
                      <div>
                        <Text>กรุณาตรวจสอบว่า:</Text>
                        <ul>
                          <li>เลขที่คุณกรอกถูกต้องและครบถ้วน</li>
                          <li>Order อยู่ในสถานะที่สามารถคืนสินค้าได้</li>
                        </ul>
                        <Text type="secondary">
                          หากยังมีปัญหา กรุณาติดต่อเจ้าหน้าที่
                        </Text>
                      </div>
                    }
                    type="error"
                    showIcon
                    action={
                      <Button size="small" onClick={handleCancel}>
                        ลองใหม่
                      </Button>
                    }
                  />
                )}
              </Space>
            </div>
          </div>
        )}

        {/* ขั้นตอนที่ไม่ใช่การค้นหา */}
        {currentStep !== "search" && orderData && (
          <Form layout="vertical" form={form} className={styles.orderFormContainer}>
            {/* แสดงข้อมูลรายการคืนสินค้า */}
            <OrderDetailsSection
              orderData={orderData}
              loading={loading}
            />

            <OrderItemsSection
              orderData={orderData}
              getReturnQty={getReturnQty}
              updateReturnQty={updateReturnQty}
              loading={loading}
              currentStep={currentStep}
            />

            {/* ขั้นตอนตรวจสอบข้อมูล */}
            {currentStep === "preview" && orderData && (
              <div className={styles.mb24}>
                <IntegratedPreviewSection
                  orderData={orderData}
                  returnItems={returnItems}
                  form={form}
                  onEdit={(step) => {
                    if (step === "create") {
                      // กลับไปขั้นตอนการสร้างคำสั่ง
                      const backBtn = renderBackButton();
                      if (!backBtn.disabled) {
                        backBtn.onClick();
                      }
                    }
                  }}
                  onNext={handleNext}
                  loading={loading}
                  stepLoading={stepLoading}
                  getReturnQty={getReturnQty}
                />
              </div>
            )}

            {/* ปุ่มดำเนินการ */}
            {currentStep !== "preview" && (
              <div className={`${styles.orderFormCard} ${styles.panelShadow}`}>
                <div className={`${styles.p24} ${styles.textCenter}`}>
                  {renderActionButtons()}
                </div>
              </div>
            )}
          </Form>
        )}

        {/* แสดงข้อผิดพลาด */}
        {error && currentStep !== "search" && (
          <Alert
            message="เกิดข้อผิดพลาด"
            description={error}
            type="error"
            showIcon
            style={{ marginTop: 16 }}
            action={
              <Button size="small" danger onClick={handleCancel}>
                ลองใหม่
              </Button>
            }
          />
        )}
      </div>
    </Spin>
  );
};

export default ReturnOrderForm;