import React, { useEffect, useState } from "react";
import {
  notification,
  Alert,
  Popconfirm,
  Layout,
  Button,
  ConfigProvider,
  Form,
  Row,
  Col,
  Select,
  FormProps,
  Input,
  DatePicker,
  Table,
  Modal,
  message,
  Tooltip,
  Spin,
  Card,
  Space,
  Descriptions,
  Typography,
} from "antd";
import { useNavigate } from "react-router-dom";
import {
  DeleteOutlined,
  LeftOutlined,
  PlusCircleOutlined,
  QuestionCircleOutlined,
  SearchOutlined,
  MinusOutlined,
  PlusOutlined,
} from "@ant-design/icons";
import { useDispatch, useSelector } from 'react-redux';
import { searchOrder, createSrNo, confirmReturn } from '../../redux/orderMKP/action';
import { RootState } from "../../redux/store";
import { ReturnOrderState } from '../../redux/orderMKP/api';  // แก้ไข path import
import { title } from "process";
import { WAREHOUSES, TRANSPORT_TYPES } from '../../constants/warehouse';

interface OrderData {
  soNo: string;
  orderNo: string;
  srCreate: string | null;
  soStatus: string;
  mkpStatus: string;
  locationTo: string;
}

// เพิ่ม interface สำหรับข้อมูลที่จะแสดง
interface OrderDetails {
  soNo: string;
  orderNo: string; 
  srNo: string | null;
  soStatus: string;
  mkpStatus: string;
  locationTo: string;
}

export const CreateReturnOrderMKP = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const { loading, error, orderData } = useSelector((state: RootState) => state.returnOrder as ReturnOrderState);
  const [selectedSalesOrder, setSelectedSalesOrder] = useState('');
  const [isChecked, setIsChecked] = useState(false);
  const [formValid, setFormValid] = useState(false);
  const [returnItems, setReturnItems] = useState<{[key: string]: number}>({});
  
  // Define columns for the table
  const columns = [
    { 
      title: 'SKU', 
      dataIndex: 'sku',
      width: '15%'
    },
    { 
      title: 'Item Name', 
      dataIndex: 'itemName',
      width: '30%'
    },
    { 
      title: 'Original QTY', 
      dataIndex: 'qty',
      width: '15%',
      align: 'center' as const,
      render: (qty: number) => Math.abs(qty)
    },
    { 
      title: 'Return QTY', 
      width: '15%',
      align: 'center' as const,
      render: (_: any, record: any) => getReturnQty(record.sku)
    },
    { 
      title: 'Price', 
      dataIndex: 'price',
      width: '15%',
      align: 'right' as const,
      render: (price: number, record: any) => {
        const returnQty = getReturnQty(record.sku);
        const totalPrice = Math.abs(price) * returnQty;
        return `฿${totalPrice.toLocaleString()}`;
      }
    },
    {
      title: 'Action',
      width: '10%',
      align: 'center' as const,
      render: (_: any, record: any) => {
        const currentQty = getReturnQty(record.sku);
        const maxQty = Math.abs(record.qty);
        
        return (
          <Space>
            <Button
              type="text"
              icon={<MinusOutlined />}
              onClick={() => updateReturnQty(record.sku, -1)}
              disabled={currentQty <= 0 || !!orderData?.head.srNo}
              danger={currentQty > 0}
            />
            <Typography.Text
              style={{
                margin: '0 8px',
                opacity: currentQty === 0 ? 0.45 : 1
              }}
            >
              {currentQty}
            </Typography.Text>
            <Button
              type="text"
              icon={<PlusOutlined />}
              onClick={() => updateReturnQty(record.sku, 1)}
              disabled={currentQty >= maxQty || !!orderData?.head.srNo}
              style={{ 
                color: currentQty >= maxQty ? undefined : '#52c41a',
                opacity: currentQty >= maxQty ? 0.45 : 1
              }}
            />
          </Space>
        );
      }
    }
  ];

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
    
    // จะ initialize returnItems เมื่อได้ข้อมูลจาก reducer
    if (orderData?.lines) {
      initializeReturnItems(orderData.lines);
    }
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
          items: returnItemsList
        };
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

  // แยกการแสดง SR Number เป็น component แยก
  const SRNumberDisplay = () => {
    const srNumber = orderData?.head.srNo;
    return (
      <Form.Item 
        label="SR Number"
        help={srNumber ? 'SR ถูกสร้างเรียบร้อยแล้ว' : 'ยังไม่มีการสร้าง SR'}
      >
        <Input
          value={srNumber || '-'}
          disabled
          style={{
            color: srNumber ? 'green' : 'inherit',
            fontWeight: srNumber ? 'bold' : 'normal'
          }}
        />
      </Form.Item>
    );
  };

  // แยกส่วนแสดงผลข้อมูลหลักออกมาเป็น component ย่อย
  const OrderDetailsSection = () => (
    <Card title="รายละเอียดคำสั่งซื้อ" style={{ marginBottom: 24 }}>
      <Row gutter={[16, 16]}>
        <Col span={8}>
          <Form.Item label="Sale Order">
            <Input value={orderData?.head.soNo} disabled />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item label="Tracking Order">
            <Input value={orderData?.head.orderNo} disabled />
          </Form.Item>
        </Col>
        <Col span={8}>
          <SRNumberDisplay />
        </Col>
        <Col span={8}>
          <Form.Item label="SO Status">
            <Input 
              value={orderData?.head.salesStatus} 
              disabled 
              style={{ color: orderData?.head.salesStatus === 'open order' ? 'green' : 'inherit' }}
            />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item label="MKP Status">
            <Input 
              value={orderData?.head.mkpStatus} 
              disabled
              style={{ color: orderData?.head.mkpStatus === 'complete' ? 'blue' : 'inherit' }}
            />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item label="Location to">
            <Input value="Return" disabled />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item 
            label="Warehouse From" 
            name="warehouseFrom"
            rules={[{ required: true, message: 'กรุณาเลือกคลังสินค้า' }]}
          >
            <Select
              options={WAREHOUSES}
              placeholder="เลือกคลังสินค้า"
              disabled={loading || !!orderData?.head.srNo}
            />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item 
            label="Return Date" 
            name="returnDate"
            rules={[{ required: true, message: 'กรุณาเลือกวันที่' }]}
          >
            <DatePicker 
              style={{ width: '100%' }}
              format="DD/MM/YYYY HH:mm"
              showTime
              placeholder="เลือกวันที่และเวลา"
              disabled={loading || !!orderData?.head.srNo}
            />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item 
            label="Tracking No" 
            name="trackingNo"
            rules={[{ required: true, message: 'กรุณากรอกเลขพัสดุ' }]}
          >
            <Input 
              placeholder="กรอกเลขพัสดุ"
              disabled={loading || !!orderData?.head.srNo}
            />
          </Form.Item>
        </Col>
        <Col span={8}>
          <Form.Item 
            label="Transport Type" 
            name="transportType"
            rules={[{ required: true, message: 'กรุณาเลือกขนส่ง' }]}
          >
            <Select
              options={TRANSPORT_TYPES}
              placeholder="เลือกประเภทขนส่ง"
              disabled={loading || !!orderData?.head.srNo}
            />
          </Form.Item>
        </Col>
      </Row>
    </Card>
  );

  const OrderItemsSection = () => (
    <Card 
      title={
        <Space>
          รายการสินค้า
          <Tooltip title="จำนวนสินค้าที่สามารถคืนได้">
            <QuestionCircleOutlined style={{ color: '#1890ff' }} />
          </Tooltip>
        </Space>
      }
    >
      <Table
        columns={columns}
        dataSource={orderData?.lines}
        rowKey="sku"
        pagination={false}
        onRow={(record) => ({
          style: getReturnQty(record.sku) === 0 ? {
            backgroundColor: '#fafafa',
            color: '#d9d9d9'
          } : {}
        })}
        summary={(pageData) => {
          const totalPrice = pageData.reduce(
            (sum, item) => sum + (Math.abs(item.price) * getReturnQty(item.sku)),
            0
          );
          const totalItems = pageData.reduce(
            (sum, item) => sum + getReturnQty(item.sku),
            0
          );
          
          return (
            <Table.Summary fixed>
              <Table.Summary.Row>
                <Table.Summary.Cell index={0} colSpan={3}>
                  รวมทั้งหมด
                </Table.Summary.Cell>
                <Table.Summary.Cell index={1} align="center">
                  {totalItems} ชิ้น
                </Table.Summary.Cell>
                <Table.Summary.Cell index={2} align="right" colSpan={2}>
                  ฿{totalPrice.toLocaleString()}
                </Table.Summary.Cell>
              </Table.Summary.Row>
            </Table.Summary>
          );
        }}
      />
    </Card>
  );

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
      [item.sku]: Math.abs(item.qty)
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
    <ConfigProvider>
      <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
        Create Return Order MKP
      </div>
      
      <Layout>
        {/* Search Section */}
        <Layout.Content style={{ margin: "24px", padding: 20, background: "#fff", borderRadius: "8px" }}>
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
        </Layout.Content>

        {/* Error Alert */}
        {error && (
          <Alert
            message="เกิดข้อผิดพลาด"
            description={error}
            type="error"
            showIcon
            style={{ margin: "0 24px" }}
          />
        )}

        {/* Order Details and Items */}
        {isChecked && orderData && (
          <Layout.Content style={{ margin: "24px", padding: "20px", background: "#fff", borderRadius: "8px" }}>
            <OrderDetailsSection />
            <OrderItemsSection />
            
            {/* Action Buttons */}
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
          </Layout.Content>
        )}
      </Layout>
    </ConfigProvider>
  );
};
