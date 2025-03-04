import React from 'react';
import { Card, Row, Col, Form, Input, Select, DatePicker } from 'antd';
import { ReturnOrderState } from '../../../../redux/orders/api';
import { WAREHOUSES, TRANSPORT_TYPES } from '../../../../constants/warehouse';

interface OrderDetailsSectionProps {
  orderData: ReturnOrderState['orderData'];
  loading: boolean;
}

const OrderDetailsSection: React.FC<OrderDetailsSectionProps> = ({ orderData, loading }) => {
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

  return (
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
};

export default OrderDetailsSection;
