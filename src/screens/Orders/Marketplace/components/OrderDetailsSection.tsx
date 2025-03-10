// src/screens/Orders/Marketplace/components/OrderDetailsSection.tsx
import React from 'react';
import { Card, Row, Col, Form, Input, Select, DatePicker, Tooltip, Divider, Typography, Badge, Space } from 'antd';
import { InfoCircleOutlined, CheckCircleOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { WAREHOUSES, TRANSPORT_TYPES } from '../../../../constants/warehouse';
import { ReturnOrderState } from '../types';

const { Text, Title } = Typography;

interface OrderDetailsSectionProps {
  orderData: ReturnOrderState['orderData'];
  loading: boolean;
}

const OrderDetailsSection: React.FC<OrderDetailsSectionProps> = ({ orderData, loading }) => {
  // แสดงข้อมูลเลข SR
  const SRNumberDisplay = () => {
    const srNumber = orderData?.head.srNo;
    
    return (
      <Form.Item 
        label={
          <Space>
            <Text strong>เลข SR Number</Text>
            <Tooltip title="เลข Sale Return Number สำหรับการคืนสินค้า">
              <InfoCircleOutlined />
            </Tooltip>
          </Space>
        }
        help={srNumber ? 'SR ถูกสร้างเรียบร้อยแล้ว' : 'ยังไม่มีการสร้าง SR'}
      >
        <Input
          value={srNumber || '-'}
          disabled
          prefix={srNumber ? <CheckCircleOutlined /> : <ClockCircleOutlined />}
          style={{
            color: srNumber ? 'green' : 'inherit',
            fontWeight: srNumber ? 'bold' : 'normal',
            backgroundColor: srNumber ? '#f6ffed' : undefined,
            borderColor: srNumber ? '#b7eb8f' : undefined
          }}
        />
      </Form.Item>
    );
  };

  // แสดงสถานะของ SO
  const renderSOStatus = () => {
    const status = orderData?.head.salesStatus;
    let color = '';
    
    switch (status) {
      case 'open order':
        color = 'green';
        break;
      case 'close order':
        color = 'red';
        break;
      default:
        color = 'blue';
    }
    
    return (
      <Space>
        <Badge color={color} />
        <Text style={{ color }}>{status}</Text>
      </Space>
    );
  };
  
  // แสดงสถานะของ MKP
  const renderMKPStatus = () => {
    const status = orderData?.head.mkpStatus;
    let color = '';
    
    switch (status) {
      case 'complete':
        color = 'blue';
        break;
      case 'pending':
        color = 'orange';
        break;
      default:
        color = 'gray';
    }
    
    return (
      <Space>
        <Badge color={color} />
        <Text style={{ color }}>{status}</Text>
      </Space>
    );
  };

  return (
    <Card 
      title={
        <Space>
          <Title level={5} style={{ margin: 0 }}>
            รายละเอียดคำสั่งซื้อ
          </Title>
          {orderData && (
            <Text type="secondary">
              ({orderData.head.soNo} / {orderData.head.orderNo})
            </Text>
          )}
        </Space>
      } 
      style={{ marginBottom: 24, boxShadow: '0 1px 3px rgba(0,0,0,0.05)' }}
    >
      <Row gutter={[24, 16]}>
        {/* แถวที่ 1 */}
        <Col xs={24} md={8}>
          <Form.Item 
            label={<Text strong>Sale Order</Text>}
            tooltip="เลขที่คำสั่งซื้อในระบบ"
          >
            <Input 
              value={orderData?.head.soNo} 
              disabled 
              style={{ fontFamily: 'monospace' }}
              addonBefore="SO"
            />
          </Form.Item>
        </Col>
        <Col xs={24} md={8}>
          <Form.Item 
            label={<Text strong>Order Number</Text>}
            tooltip="เลขที่ Order ในระบบ"
          >
            <Input 
              value={orderData?.head.orderNo} 
              disabled 
              style={{ fontFamily: 'monospace' }}
            />
          </Form.Item>
        </Col>
        <Col xs={24} md={8}>
          <SRNumberDisplay />
        </Col>

        {/* แถวที่ 2 */}
        <Col xs={24} md={8}>
          <Form.Item 
            label={<Text strong>SO Status</Text>}
            tooltip="สถานะของคำสั่งซื้อ"
          >
            <Input 
              value={orderData?.head.salesStatus} 
              disabled 
              addonAfter={orderData && renderSOStatus()}
            />
          </Form.Item>
        </Col>
        <Col xs={24} md={8}>
          <Form.Item 
            label={<Text strong>MKP Status</Text>}
            tooltip="สถานะของการทำรายการใน Marketplace"
          >
            <Input 
              value={orderData?.head.mkpStatus} 
              disabled
              addonAfter={orderData && renderMKPStatus()}
            />
          </Form.Item>
        </Col>
        <Col xs={24} md={8}>
          <Form.Item 
            label={<Text strong>Location to</Text>}
            tooltip="สถานที่ปลายทางสำหรับการคืนสินค้า"
          >
            <Input 
              value={orderData?.head.locationTo || "Return"} 
              disabled 
              addonAfter={<Badge status="processing" text="Return" />}
            />
          </Form.Item>
        </Col>

        <Col span={24}>
          <Divider plain>ข้อมูลการคืนสินค้า</Divider>
        </Col>

        {/* แถวที่ 3 */}
        <Col xs={24} md={8}>
          <Form.Item 
            label={
              <Space>
                <Text strong>คลังสินค้า</Text>
                <Badge count="จำเป็น" style={{ backgroundColor: '#ff4d4f' }} />
              </Space>
            }
            name="warehouseFrom"
            tooltip="เลือกคลังสินค้าที่จะรับสินค้าคืน"
            rules={[{ required: true, message: 'กรุณาเลือกคลังสินค้า' }]}
          >
            <Select
              options={WAREHOUSES}
              placeholder="เลือกคลังสินค้า"
              disabled={loading || !!orderData?.head.srNo}
              showSearch
              optionFilterProp="label"
              style={{ width: '100%' }}
            />
          </Form.Item>
        </Col>
        <Col xs={24} md={8}>
          <Form.Item 
            label={
              <Space>
                <Text strong>วันที่คืนสินค้า</Text>
                <Badge count="จำเป็น" style={{ backgroundColor: '#ff4d4f' }} />
              </Space>
            }
            name="returnDate"
            tooltip="กำหนดวันที่และเวลาที่จะคืนสินค้า"
            rules={[{ required: true, message: 'กรุณาเลือกวันที่' }]}
          >
            <DatePicker 
              style={{ width: '100%' }}
              format="DD/MM/YYYY HH:mm"
              showTime={{ format: 'HH:mm' }}
              placeholder="เลือกวันที่และเวลา"
              disabled={loading || !!orderData?.head.srNo}
            />
          </Form.Item>
        </Col>
        <Col xs={24} md={8}>
          <Form.Item 
            label={
              <Space>
                <Text strong>เลขพัสดุ</Text>
                <Badge count="จำเป็น" style={{ backgroundColor: '#ff4d4f' }} />
              </Space>
            }
            name="trackingNo"
            tooltip="เลขติดตามพัสดุสำหรับการคืนสินค้า"
            rules={[{ required: true, message: 'กรุณากรอกเลขพัสดุ' }]}
          >
            <Input 
              placeholder="กรอกเลขพัสดุ"
              disabled={loading || !!orderData?.head.srNo}
            />
          </Form.Item>
        </Col>

        {/* แถวที่ 4 */}
        <Col xs={24} md={8}>
          <Form.Item 
            label={
              <Space>
                <Text strong>ประเภทขนส่ง</Text>
                <Badge count="จำเป็น" style={{ backgroundColor: '#ff4d4f' }} />
              </Space>
            }
            name="transportType"
            tooltip="เลือกประเภทการขนส่ง"
            rules={[{ required: true, message: 'กรุณาเลือกขนส่ง' }]}
          >
            <Select
              options={TRANSPORT_TYPES}
              placeholder="เลือกประเภทขนส่ง"
              disabled={loading || !!orderData?.head.srNo}
              showSearch
              optionFilterProp="label"
              style={{ width: '100%' }}
            />
          </Form.Item>
        </Col>
        <Col xs={24} md={16}>
          <Form.Item 
            label={<Text strong>เหตุผลการคืน</Text>}
            name="reason"
            tooltip="ระบุเหตุผลในการคืนสินค้า"
          >
            <Input 
              placeholder="ระบุเหตุผลในการคืนสินค้า (ถ้ามี)"
              disabled={loading || !!orderData?.head.srNo}
            />
          </Form.Item>
        </Col>
      </Row>
    </Card>
  );
};

export default OrderDetailsSection;