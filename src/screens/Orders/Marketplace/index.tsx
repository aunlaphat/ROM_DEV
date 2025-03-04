import React, { useEffect, useState } from "react";
import { Layout, Button, Form, Row, Col, Input, Alert, Modal, message } from "antd";
import { LeftOutlined } from "@ant-design/icons";
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from "react-router-dom";
import { searchOrder, createSrNo, confirmReturn } from '../../../redux/orders/action';
import { RootState } from "../../../redux/store";
import { ReturnOrderState } from '../../../redux/orders/api';
import OrderDetailsSection from './components/OrderDetailsSection';
import OrderItemsSection from './components/OrderItemsSection';

const { Content } = Layout;

const CreateReturnOrderMKP = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const { loading, error, orderData } = useSelector((state: RootState) => state.returnOrder as ReturnOrderState);
  const [selectedSalesOrder, setSelectedSalesOrder] = useState('');
  const [isChecked, setIsChecked] = useState(false);
  const [returnItems, setReturnItems] = useState<{[key: string]: number}>({});

  // Handler functions
  const handleBack = () => navigate('/create-return-order-mkp');
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedSalesOrder(e.target.value.trim());
  };
  const handleSearch = async () => {
    if (!selectedSalesOrder) {
      message.error('กรุณากรอกเลข SO/Order');
      return;
    }
  
    const isSoNo = selectedSalesOrder.startsWith('SO');
    const searchPayload = {
      [isSoNo ? 'soNo' : 'orderNo']: selectedSalesOrder.trim()
    };
    
    dispatch(searchOrder(searchPayload));
    setIsChecked(true);
  };
  
  // Alternative approach using useEffect
  useEffect(() => {
    if (orderData?.lines) {
      initializeReturnItems(orderData.lines);
    }
  }, [orderData]);
  
  const handleCancel = () => {
    form.resetFields();
    setSelectedSalesOrder('');
    setIsChecked(false);
  };
  const handleCreateSR = () => {
    if (!orderData?.head.orderNo) {
      message.error('ไม่พบเลขที่ Order');
      return;
    }

    const formValues = form.getFieldsValue();
    const returnItemsList = orderData.lines
      .filter(item => getReturnQty(item.sku) > 0)
      .map(item => ({
        ...item,
        returnQty: getReturnQty(item.sku)
      }));

    if (returnItemsList.length === 0) {
      message.error('กรุณาระบุจำนวนสินค้าที่ต้องการคืน');
      return;
    }

    Modal.confirm({
      title: 'ยืนยันการสร้าง SR',
      content: `ต้องการสร้าง SR สำหรับ ${returnItemsList.length} รายการ ใช่หรือไม่?`,
      okText: 'ยืนยัน',
      cancelText: 'ยกเลิก',
      onOk: () => {
        const payload = {
          orderNo: orderData.head.orderNo,
          ...formValues,
          returnDate: formValues.returnDate.toISOString(), // แปลงค่า Date ให้เป็น string
          items: returnItemsList
        };
        console.log('Payload:', payload); // เพิ่มการแสดงผล payload เพื่อช่วยในการดีบัก
        dispatch(createSrNo(payload));
      }
    });
  };

  // ปรับปรุงฟังก์ชัน helper
  const isCreateSRDisabled = (): boolean => {
    if (!orderData) return true;
    if (!orderData.head) return true;
    if (loading) return true;
    if (orderData.head.srNo !== null) return true;
    return !validateAdditionalFields();
  };

  // เพิ่มฟังก์ชันตรวจสอบการกรอกข้อมูลครบถ้วน
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
    const initialQty = items.reduce((acc, item) => ({
      ...acc,
      [item.sku]: 0 // เริ่มต้นเป็น 0 เพื่อให้ผู้ใช้กรอกจำนวนที่ต้องการคืน
    }), {});
    setReturnItems(initialQty);
  };

  const getReturnQty = (sku: string): number => {
    return returnItems[sku] || 0;
  };

  const updateReturnQty = (sku: string, change: number) => {
    const currentQty = getReturnQty(sku);
    const originalQty = Math.abs(orderData?.lines.find(item => item.sku === sku)?.qty || 0);
    const newQty = Math.max(0, Math.min(originalQty, currentQty + change));
    
    setReturnItems(prev => ({
      ...prev,
      [sku]: newQty
    }));
  };

  return (
    <Layout>
      <Content style={{ margin: "24px", padding: 20, background: "#fff", borderRadius: "8px" }}>
        <Button onClick={handleBack} style={{ background: '#98CEFF', color: '#fff' }}>
          <LeftOutlined style={{ color: '#fff', marginRight: 5 }} />
          Back
        </Button>
        
        <Form layout="vertical" form={form} style={{ marginTop: '40px' }}>
          <Row gutter={30} justify="center" align="middle">
            <Col>
              <Form.Item
                label={<span style={{ color: '#657589' }}>กรอกเลข SO/Order</span>}
                name="selectedSalesOrder"
                rules={[{ required: true, message: 'กรุณากรอกเลข SO/Order!' }]}
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
      </Content>

      {error && (
        <Alert
          message="เกิดข้อผิดพลาด"
          description={error}
          type="error"
          showIcon
          style={{ margin: "0 24px" }}
        />
      )}

      {isChecked && orderData && (
        <Form layout="vertical" form={form}>
          <Content style={{ margin: "24px", padding: "20px", background: "#fff", borderRadius: "8px" }}>
            <OrderDetailsSection orderData={orderData} loading={loading} />
            <OrderItemsSection 
              orderData={orderData} 
              getReturnQty={getReturnQty} 
              updateReturnQty={updateReturnQty} 
              loading={loading} 
            />
            
            <Row justify="end" style={{ marginTop: 24 }}>
              <Button
                type="primary"
                onClick={handleCreateSR}
                loading={loading}
                disabled={isCreateSRDisabled()}
                style={{ marginRight: 8 }}
              >
                Create SR
              </Button>
              <Button onClick={handleCancel}>
                Cancel
              </Button>
            </Row>
          </Content>
        </Form>
      )}
    </Layout>
  );
};

export default CreateReturnOrderMKP;
