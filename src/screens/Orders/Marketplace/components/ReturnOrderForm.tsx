import React from "react";
import { Layout, Button, Form, Row, Col, Input, Alert, Spin } from "antd";
import OrderDetailsSection from './OrderDetailsSection';
import OrderItemsSection from './OrderItemsSection';
import ReturnOrderSteps from "./ReturnOrderSteps";
import PreviewSection from './PreviewSection';

interface ReturnOrderFormProps {
  currentStep: 'search' | 'create' | 'sr' | 'preview' | 'confirm';
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
  handleNext: () => void; // เพิ่ม handler สำหรับปุ่มดำเนินการต่อ
  handleConfirm: () => void; // เพิ่ม handler สำหรับปุ่ม Confirm
  getReturnQty: (sku: string) => number;
  updateReturnQty: (sku: string, change: number) => void;
  isCreateReturnOrderDisabled: () => boolean;
  getStepStatus: (stepKey: string) => 'process' | 'finish' | 'wait';
  renderBackButton: () => JSX.Element;
  returnItems: { [key: string]: number }; // เพิ่ม returnItems
  validateStepTransition: (fromStep: string, toStep: string) => boolean;
  stepLoading: boolean; // เพิ่ม stepLoading prop
  isCreateSRDisabled: () => boolean; // เพิ่ม prop ใหม่
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
  handleCreateSr, // เพิ่ม handleCreateSr
  handleCancel,
  handleNext, // เพิ่ม handleNext
  handleConfirm,
  getReturnQty,
  updateReturnQty,
  isCreateReturnOrderDisabled,
  getStepStatus,
  renderBackButton,
  returnItems, // เพิ่ม returnItems
  validateStepTransition,
  stepLoading = false,
  isCreateSRDisabled,
}) => {
  // ใช้ validateStepTransition ในการตรวจสอบก่อนแสดงปุ่มต่างๆ
  const renderActionButtons = () => {
    switch (currentStep) {
      case 'create':
        return (
          <Button
            type="primary"
            onClick={handleCreateReturnOrder}
            loading={loading || stepLoading}
            disabled={isCreateReturnOrderDisabled()} // ใช้ฟังก์ชันที่ปรับปรุงแล้ว
            style={{ marginRight: 8 }}
          >
            Create Return Order
          </Button>
        );
      
      case 'sr':
        // ถ้ามี SR Number แล้วให้แสดงปุ่ม Next แทน Create SR
        if (orderData?.head.srNo) {
          return (
            <Button
              type="primary"
              onClick={handleNext}
              loading={loading || stepLoading}
              disabled={!validateStepTransition('sr', 'preview')}
              style={{ marginRight: 8 }}
            >
              Next
            </Button>
          );
        }
        // ถ้ายังไม่มี SR Number ให้แสดงปุ่ม Create SR
        return (
          <Button
            type="primary"
            onClick={handleCreateSr}
            loading={loading || stepLoading}
            disabled={isCreateSRDisabled()}
            style={{ marginRight: 8 }}
          >
            Create SR
          </Button>
        );
      
      case 'preview':
        return (
          <Button
            type="primary"
            onClick={handleNext}
            loading={loading}
            disabled={!validateStepTransition('preview', 'confirm')}
          >
            Next
          </Button>
        );
      
      case 'confirm':
        return (
          <Button
            type="primary"
            onClick={handleConfirm}
            loading={loading}
          >
            ยืนยันคำสั่งคืนสินค้า
          </Button>
        );
      
      default:
        return null;
    }
  };

  return (
    <Spin spinning={loading || stepLoading} tip={
      stepLoading ? 'กำลังดำเนินการ...' : 'กำลังโหลดข้อมูล...'
    }>
      <Layout>
        <Layout.Content
          style={{
            margin: "24px",
            padding: 20,
            background: "#f5f5f5",
            borderRadius: "8px",
          }}
        >
          <Row gutter={[0, 24]}>
            <Col span={24}>
              <ReturnOrderSteps
                currentStep={currentStep}
                orderData={orderData}
                getStepStatus={getStepStatus}
              />
            </Col>
            <Col span={24}>{renderBackButton()}</Col>
          </Row>

          {currentStep === "search" && (
            <Form layout="vertical" form={form} style={{ marginTop: "40px" }}>
              <Row gutter={30} justify="center" align="middle">
                <Col>
                  <Form.Item
                    label={
                      <span style={{ color: "#657589" }}>กรอกเลข SO/Order</span>
                    }
                    name="selectedSalesOrder"
                    rules={[
                      { required: true, message: "กรุณากรอกเลข SO/Order!" },
                    ]}
                  >
                    <Input
                      style={{ width: 300 }}
                      placeholder="กรอกเลข SO/Order"
                      value={selectedSalesOrder}
                      onChange={handleInputChange}
                      disabled={loading}
                    />
                  </Form.Item>
                </Col>
                <Col>
                  <Button
                    type="primary"
                    onClick={handleSearch}
                    loading={loading}
                    style={{ marginTop: 4 }}
                  >
                    ค้นหา
                  </Button>
                </Col>
              </Row>
            </Form>
          )}

          {currentStep !== "search" && orderData && (
            <Form layout="vertical" form={form}>
              <Layout.Content
                style={{
                  margin: "24px",
                  padding: "20px",
                  background: "#fff",
                  borderRadius: "8px",
                }}
              >
                <OrderDetailsSection orderData={orderData} loading={loading} />
                <OrderItemsSection
                  orderData={orderData}
                  getReturnQty={getReturnQty}
                  updateReturnQty={updateReturnQty}
                  loading={loading}
                  currentStep={currentStep} // เพิ่ม prop currentStep
                />

                {currentStep === "preview" && orderData && (
                  <PreviewSection
                    orderData={orderData}
                    returnItems={returnItems}
                  />
                )}

                <Row justify="end" style={{ marginTop: 24 }}>
                  {renderActionButtons()}
                  <Button 
                    onClick={handleCancel}
                    disabled={loading || stepLoading}
                  >
                    Cancel
                  </Button>
                </Row>
              </Layout.Content>
            </Form>
          )}
        </Layout.Content>

        {error && (
          <Alert
            message="เกิดข้อผิดพลาด"
            description={error}
            type="error"
            showIcon
            style={{ margin: "0 24px" }}
          />
        )}
      </Layout>
    </Spin>
  );
};

export default ReturnOrderForm
