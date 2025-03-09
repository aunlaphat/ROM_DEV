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
} from "antd";
import {
  LeftOutlined,
  SearchOutlined,
  ArrowRightOutlined,
  CheckOutlined,
  QuestionCircleOutlined,
} from "@ant-design/icons";
import OrderDetailsSection from "./OrderDetailsSection";
import OrderItemsSection from "./OrderItemsSection";
import ReturnOrderSteps from "./ReturnOrderSteps";
import PreviewSection from "./PreviewSection";

const { Content } = Layout;
const { Title, Text } = Typography;

// สร้าง interface สำหรับ OrderLineItem เพื่อให้มี type ที่ชัดเจน
interface OrderLineItem {
  sku: string;
  price: number;
  qty: number;
  itemName: string;
}

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

  // ฟังก์ชันคำนวณจำนวนสินค้าที่คืนทั้งหมด - แก้ไขโดยระบุประเภทข้อมูล
  const getTotalReturnItems = (): number => {
    return Object.values(returnItems).reduce((sum: number, qty: number) => sum + qty, 0);
  };

  // ฟังก์ชันคำนวณมูลค่าสินค้าที่คืนทั้งหมด - แก้ไขโดยระบุประเภทข้อมูล
  const getTotalReturnAmount = (): number => {
    if (!orderData?.lines) return 0;

    return orderData.lines.reduce(
      (sum: number, item: OrderLineItem) => sum + Math.abs(item.price) * (returnItems[item.sku] || 0),
      0
    );
  };

  return (
    <Spin
      spinning={loading || stepLoading}
      tip={stepLoading ? "กำลังดำเนินการ..." : "กำลังโหลดข้อมูล..."}
    >
      <Layout className="return-order-container">
        <Content
          style={{
            margin: "24px",
            padding: 24,
            background: "#f5f5f5",
            borderRadius: "8px",
            minHeight: "80vh",
          }}
        >
          {/* ส่วนหัวและขั้นตอน */}
          <Row gutter={[0, 24]}>
            <Col span={24}>
              <div
                style={{
                  background: "#fff",
                  padding: "16px 24px",
                  borderRadius: "8px",
                  boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
                }}
              >
                <Title level={3} style={{ marginBottom: 12, color: "#1890ff" }}>
                  ระบบคืนสินค้า - {getStepTitle()}
                </Title>
                <Text type="secondary">{getStepDescription()}</Text>
              </div>
            </Col>
            <Col span={24}>
              <ReturnOrderSteps
                currentStep={currentStep}
                orderData={orderData}
                getStepStatus={getStepStatus}
              />
            </Col>
          </Row>

          <Divider style={{ margin: "24px 0" }} />

          {/* ขั้นตอนค้นหา */}
          {currentStep === "search" && (
            <div
              className="search-container"
              style={{
                maxWidth: "600px",
                margin: "40px auto",
                background: "#fff",
                padding: "36px",
                borderRadius: "8px",
                boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
              }}
            >
              <Title
                level={4}
                style={{ textAlign: "center", marginBottom: 24 }}
              >
                ค้นหาคำสั่งซื้อ
                <Tooltip title="กรอกเลข SO ที่ขึ้นต้นด้วย SO หรือเลข Order">
                  <QuestionCircleOutlined
                    style={{ fontSize: "16px", marginLeft: 8 }}
                  />
                </Tooltip>
              </Title>

              <Form layout="vertical" form={form} style={{ marginTop: "20px" }}>
                <Row gutter={16} justify="center">
                  <Col span={18}>
                    <Form.Item
                      label={
                        <span style={{ fontWeight: 500 }}>
                          กรอกเลข SO/Order
                        </span>
                      }
                      name="selectedSalesOrder"
                      rules={[
                        { required: true, message: "กรุณากรอกเลข SO/Order!" },
                      ]}
                    >
                      <Input
                        size="large"
                        placeholder="เช่น SO12345678 หรือ OR98765432"
                        value={selectedSalesOrder}
                        onChange={handleInputChange}
                        disabled={loading}
                        prefix={<SearchOutlined style={{ color: "#1890ff" }} />}
                        onPressEnter={handleSearch}
                      />
                    </Form.Item>
                  </Col>
                </Row>

                <Row justify="center" style={{ marginTop: 16 }}>
                  <Button
                    type="primary"
                    size="large"
                    icon={<SearchOutlined />}
                    onClick={handleSearch}
                    loading={loading || stepLoading}
                    disabled={!selectedSalesOrder}
                    style={{ width: "120px" }}
                  >
                    ค้นหา
                  </Button>
                </Row>
              </Form>
            </div>
          )}

          {/* ขั้นตอนที่ไม่ใช่การค้นหา */}
          {currentStep !== "search" && orderData && (
            <Form layout="vertical" form={form}>
              <Layout.Content
                style={{
                  padding: "24px",
                  background: "#fff",
                  borderRadius: "8px",
                  boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
                }}
              >
                {/* แสดงข้อมูลรายการคืนสินค้า */}
                <Row gutter={[16, 16]}>
                  <Col span={24}>
                    <OrderDetailsSection
                      orderData={orderData}
                      loading={loading}
                    />
                  </Col>

                  <Col span={24}>
                    <OrderItemsSection
                      orderData={orderData}
                      getReturnQty={getReturnQty}
                      updateReturnQty={updateReturnQty}
                      loading={loading}
                      currentStep={currentStep}
                    />
                  </Col>

                  {/* แสดงข้อมูลสรุป */}
                  {getTotalReturnItems() > 0 && currentStep !== "preview" && (
                    <Col span={24}>
                      <div
                        style={{
                          background: "#f6f8fa",
                          padding: "16px",
                          borderRadius: "8px",
                          marginTop: "12px",
                        }}
                      >
                        <Row justify="space-between" align="middle">
                          <Col>
                            <Text strong>สรุปรายการสินค้าที่จะคืน:</Text>
                          </Col>
                          <Col>
                            <Space>
                              <Text>
                                จำนวนสินค้า:{" "}
                                <Text strong>{getTotalReturnItems()} ชิ้น</Text>
                              </Text>
                              <Divider type="vertical" />
                              <Text>
                                มูลค่ารวม:{" "}
                                <Text strong style={{ color: "#52c41a" }}>
                                  ฿{getTotalReturnAmount().toLocaleString()}
                                </Text>
                              </Text>
                            </Space>
                          </Col>
                        </Row>
                      </div>
                    </Col>
                  )}

                  {/* ขั้นตอนตรวจสอบข้อมูล */}
                  {currentStep === "preview" && orderData && (
                    <Col span={24}>
                      <Divider style={{ margin: "12px 0 24px" }} />
                      <PreviewSection
                        orderData={orderData}
                        returnItems={returnItems}
                      />
                    </Col>
                  )}

                  {/* ปุ่มดำเนินการ */}
                  <Col span={24}>
                    <Divider style={{ margin: "24px 0" }} />
                    <Row justify="end">{renderActionButtons()}</Row>
                  </Col>
                </Row>
              </Layout.Content>
            </Form>
          )}
        </Content>

        {/* แสดงข้อผิดพลาด */}
        {error && (
          <Alert
            message="เกิดข้อผิดพลาด"
            description={error}
            type="error"
            showIcon
            style={{
              margin: "16px 24px",
              boxShadow: "0 2px 8px rgba(0,0,0,0.08)",
            }}
            action={
              <Button size="small" danger onClick={handleCancel}>
                ลองใหม่
              </Button>
            }
          />
        )}
      </Layout>
    </Spin>
  );
};

export default ReturnOrderForm;