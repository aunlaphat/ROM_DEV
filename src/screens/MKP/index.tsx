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
} from "antd";
import { useNavigate } from "react-router-dom";
import {
  DeleteOutlined,
  LeftOutlined,
  PlusCircleOutlined,
  QuestionCircleOutlined,
  SearchOutlined,
} from "@ant-design/icons";
import { useDispatch, useSelector } from 'react-redux';
import { searchOrder, createSrNo, confirmReturn } from '../../redux/orderMKP/action';
import { RootState } from "../../redux/store";
import { ReturnOrderState } from '../../redux/orderMKP/api';  // แก้ไข path import
import { title } from "process";

interface OrderData {
  soNo: string;
  orderNo: string;
  srCreate: string | null;
  soStatus: string;
  mkpStatus: string;
  locationTo: string;
}

interface OrderLineItem {
  sku: string;
  itemName: string;
  qty: number;
  price: number;
}

export const CreateReturnOrderMKP = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const { loading, error, orderData } = useSelector((state: RootState) => state.returnOrder as ReturnOrderState);
  const [selectedSalesOrder, setSelectedSalesOrder] = useState('');
  const [isChecked, setIsChecked] = useState(false);
  const [formValid, setFormValid] = useState(false);
  
  // Define columns for the table
  const columns = [
    { title: 'SKU', dataIndex: 'sku' },
    { title: 'Item Name', dataIndex: 'itemName' },
    { title: 'Quantity', dataIndex: 'qty' },
    { 
      title: 'Price', 
      dataIndex: 'price',
      render: (price: number) => `฿${price.toLocaleString()}`
    }
  ];

  // Handler functions
  const handleBack = () => navigate('/home');
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedSalesOrder(e.target.value.trim());
  };
  const handleCheck = () => {
    if (!selectedSalesOrder) {
      message.error('กรุณากรอกเลข SO/Order');
      return;
    }
    dispatch(searchOrder(selectedSalesOrder));
    setIsChecked(true);
  };
  const handleCancel = () => {
    form.resetFields();
    setSelectedSalesOrder('');
    setIsChecked(false);
  };
  const handleCreateSR = () => {
    if (orderData?.head.orderNo) {
      dispatch(createSrNo(orderData.head.orderNo));
    }
  };

  return (
    <ConfigProvider>
      <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
        Create Return Order MKP
      </div>
      <Layout>
        <Layout.Content style={{
          margin: "24px",
          padding: 20,
          minHeight: 200,
          background: "#fff",
          borderRadius: "8px",
          display: 'flex',
        }}>
          <Button
            id="backButton"
            onClick={handleBack}
            style={{ background: '#98CEFF', color: '#fff' }}
          >
            <LeftOutlined style={{ color: '#fff', marginRight: 5 }} />
            Back
          </Button>
          <Form
            layout="vertical"
            form={form}
            style={{ width: '100%', marginTop: '40px' }}
          >
            <Row gutter={30} justify="center" align="middle" style={{ width: '100%' }}>
              <Col>
                <Form.Item
                  id="salesOrderFormItem"
                  label={<span style={{ color: '#657589' }}>กรอกเลข SO/Order ที่ต้องการสร้าง SR</span>}
                  name="selectedSalesOrder"
                  rules={[
                    { required: true, message: 'กรุณากรอกเลข SO/Order ที่ต้องการสร้าง SR!' }
                  ]}
                >
                  <Input
                    id="salesOrderInput"
                    style={{ height: 40, width: 300 }}
                    placeholder="กรอก SO/Order"
                    value={selectedSalesOrder}
                    onChange={handleInputChange}
                    disabled={loading}
                  />
                </Form.Item>
              </Col>
              <Col>
                <Button
                  id="checkButton"
                  type="primary"
                  style={{ width: 100, height: 40, marginTop: 4 }}
                  onClick={handleCheck}
                  loading={loading}
                >
                  Check
                </Button>
              </Col>
            </Row>
          </Form>
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

        {isChecked && orderData && (
          <Layout.Content 
            id="checkedContent"
            style={{
              marginRight: 24,
              marginLeft: 24,
              padding: 36,
              minHeight: 360,
              background: "#fff",
              borderRadius: "8px",
            }}
          >
            <Form form={form} layout="vertical">
              <Row gutter={16}>
                <Col span={8}>
                  <Form.Item label="Sale Order">
                    <Input value={orderData.head.soNo} disabled />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item label="SR Number">
                    <Input value={orderData.head.srNo || '-'} disabled />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item label="Status">
                    <Input value={orderData.head.salesStatus} disabled />
                  </Form.Item>
                </Col>
              </Row>
              
              <Table
                style={{ marginTop: 24 }}
                columns={columns}
                dataSource={orderData.lines}
                rowKey="sku"
                pagination={false}
              />

              <Row justify="center" style={{ marginTop: 24 }}>
                <Button
                  type="primary"
                  onClick={handleCreateSR}
                  style={{ marginRight: 8 }}
                  disabled={!formValid}
                >
                  Create SR
                </Button>
                <Button onClick={handleCancel}>
                  Cancel
                </Button>
              </Row>
            </Form>
          </Layout.Content>
        )}
      </Layout>
    </ConfigProvider>
  );
};
