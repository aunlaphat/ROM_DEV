// src/screens/Orders/Marketplace/components/OrderItemsSection.tsx
import React, { useState, useMemo } from 'react';
import { Card, Table, Space, Button, Typography, Tooltip, Input, Badge, Empty, Row, Col, Tag } from 'antd';
import { 
  MinusOutlined, 
  PlusOutlined, 
  QuestionCircleOutlined, 
  SearchOutlined, 
  EditOutlined,
  CheckCircleOutlined
} from "@ant-design/icons";
import { ReturnOrderState } from '../types';

const { Text } = Typography;

interface OrderItemsSectionProps {
  orderData: ReturnOrderState['orderData'];
  getReturnQty: (sku: string) => number;
  updateReturnQty: (sku: string, change: number) => void;
  loading: boolean;
  currentStep?: 'search' | 'create' | 'sr' | 'preview' | 'confirm';
}

const OrderItemsSection: React.FC<OrderItemsSectionProps> = ({ 
  orderData, 
  getReturnQty, 
  updateReturnQty, 
  loading, 
  currentStep = 'create' 
}) => {
  const [searchText, setSearchText] = useState('');
  const [expandedRowKeys, setExpandedRowKeys] = useState<React.Key[]>([]);

  // กรองข้อมูลตามคำค้นหา
  const filteredData = useMemo(() => {
    if (!orderData?.lines || !searchText) return orderData?.lines;
    
    return orderData.lines.filter(item => 
      item.sku.toLowerCase().includes(searchText.toLowerCase()) || 
      item.itemName.toLowerCase().includes(searchText.toLowerCase())
    );
  }, [orderData, searchText]);

  // ได้ข้อมูลที่มี items ที่มีการคืนสินค้า
  const hasReturnItems = useMemo(() => {
    if (!orderData?.lines) return false;
    return orderData.lines.some(item => getReturnQty(item.sku) > 0);
  }, [orderData, getReturnQty]);

  // ตรวจสอบว่าผู้ใช้สามารถแก้ไขจำนวนได้หรือไม่
  const canEdit = useMemo(() => {
    return currentStep === 'create';
  }, [currentStep]);

  // คอลัมน์สำหรับตาราง
  const columns = [
    { 
      title: 'SKU', 
      dataIndex: 'sku',
      width: '15%',
      render: (sku: string) => (
        <Text copyable={{ text: sku, tooltips: ['คัดลอก SKU', 'คัดลอกแล้ว!'] }}>
          {sku}
        </Text>
      )
    },
    { 
      title: 'รายการสินค้า', 
      dataIndex: 'itemName',
      width: '30%',
      render: (name: string, record: any) => (
        <Space direction="vertical" size={0}>
          <Text strong>{name}</Text>
          <Text type="secondary" style={{ fontSize: '12px' }}>
            {record.warehouse ? `คลัง: ${record.warehouse}` : ''}
          </Text>
        </Space>
      )
    },
    { 
      title: 'จำนวนทั้งหมด', 
      dataIndex: 'qty',
      width: '12%',
      align: 'center' as const,
      render: (qty: number) => (
        <Tag color="blue" style={{ fontSize: '14px', padding: '0 8px' }}>
          {Math.abs(qty)}
        </Tag>
      )
    },
    { 
      title: (
        <Space>
          <Text>จำนวนที่คืน</Text>
          <Tooltip title="จำนวนที่ต้องการคืน (ไม่เกินจำนวนทั้งหมด)">
            <QuestionCircleOutlined style={{ color: '#1890ff' }} />
          </Tooltip>
        </Space>
      ), 
      width: '18%',
      align: 'center' as const,
      render: (_: any, record: any) => {
        const returnQty = getReturnQty(record.sku);
        const totalQty = Math.abs(record.qty);
        const percent = Math.round((returnQty / totalQty) * 100);
        
        // แสดงปุ่ม +/- เฉพาะในขั้นตอน create
        if (canEdit) {
          return (
            <Space>
              <Button
                type="text"
                size="small"
                icon={<MinusOutlined />}
                onClick={() => updateReturnQty(record.sku, -1)}
                disabled={returnQty <= 0}
                danger={returnQty > 0}
              />
              <Space direction="vertical" size={0} style={{ width: '40px' }}>
                <Text strong style={{ fontSize: '16px' }}>{returnQty}</Text>
                {returnQty > 0 && (
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    {percent}%
                  </Text>
                )}
              </Space>
              <Button
                type="text"
                size="small"
                icon={<PlusOutlined />}
                onClick={() => updateReturnQty(record.sku, 1)}
                disabled={returnQty >= totalQty}
                style={{ color: returnQty < totalQty ? '#52c41a' : undefined }}
              />
            </Space>
          );
        }
        
        // แสดงแค่ตัวเลขในขั้นตอนอื่นๆ
        if (returnQty > 0) {
          return (
            <Badge 
              count={returnQty} 
              overflowCount={999}
              style={{ 
                backgroundColor: '#52c41a',
                fontSize: '14px',
                padding: '0 8px',
                minWidth: '40px'
              }}
            />
          );
        }
        
        return (
          <Text type="secondary">0</Text>
        );
      }
    },
    { 
      title: 'ราคาต่อชิ้น', 
      dataIndex: 'price',
      width: '12%',
      align: 'right' as const,
      render: (price: number) => (
        <Text>฿{Math.abs(price).toLocaleString()}</Text>
      )
    },
    { 
      title: 'รวม', 
      width: '13%',
      align: 'right' as const,
      render: (_: any, record: any) => {
        const returnQty = getReturnQty(record.sku);
        const totalPrice = Math.abs(record.price) * returnQty;
        
        if (returnQty > 0) {
          return (
            <Text strong style={{ color: '#52c41a' }}>
              ฿{totalPrice.toLocaleString()}
            </Text>
          );
        }
        
        return (
          <Text type="secondary">฿0</Text>
        );
      }
    }
  ];

  // Function สำหรับการขยายแถว
  const handleRowExpand = (expanded: boolean, record: any) => {
    setExpandedRowKeys(expanded ? [record.sku] : []);
  };
  
  // Render เนื้อหาเมื่อขยายแถว
  const expandedRowRender = (record: any) => {
    const returnQty = getReturnQty(record.sku);
    const totalQty = Math.abs(record.qty);
    
    return (
      <div style={{ padding: '10px 40px' }}>
        <Row gutter={[24, 16]}>
          <Col span={8}>
            <Text type="secondary">SKU:</Text>
            <div>
              <Text copyable strong>{record.sku}</Text>
            </div>
          </Col>
          <Col span={8}>
            <Text type="secondary">จำนวนทั้งหมด:</Text>
            <div>
              <Text strong>{totalQty} ชิ้น</Text>
            </div>
          </Col>
          <Col span={8}>
            <Text type="secondary">จำนวนที่คืน:</Text>
            <div>
              <Text strong>{returnQty} ชิ้น</Text>
              {returnQty > 0 && (
                <Tag color="success" style={{ marginLeft: 8 }}>
                  {Math.round((returnQty / totalQty) * 100)}%
                </Tag>
              )}
            </div>
          </Col>
          <Col span={8}>
            <Text type="secondary">ราคาต่อชิ้น:</Text>
            <div>
              <Text strong>฿{Math.abs(record.price).toLocaleString()}</Text>
            </div>
          </Col>
          <Col span={8}>
            <Text type="secondary">ราคารวม:</Text>
            <div>
              <Text strong style={{ color: '#52c41a' }}>
                ฿{(Math.abs(record.price) * returnQty).toLocaleString()}
              </Text>
            </div>
          </Col>
          {canEdit && (
            <Col span={8}>
              <Text type="secondary">การดำเนินการ:</Text>
              <div>
                <Space>
                  <Button 
                    size="small" 
                    type="primary"
                    icon={<EditOutlined />}
                    onClick={() => updateReturnQty(record.sku, totalQty - returnQty)}
                    disabled={returnQty >= totalQty}
                  >
                    คืนทั้งหมด
                  </Button>
                  <Button 
                    size="small" 
                    danger 
                    onClick={() => updateReturnQty(record.sku, -returnQty)}
                    disabled={returnQty <= 0}
                  >
                    ล้างค่า
                  </Button>
                </Space>
              </div>
            </Col>
          )}
        </Row>
      </div>
    );
  };

  return (
    <Card 
      title={
        <Space align="center">
          <CheckCircleOutlined 
            style={{ 
              color: hasReturnItems ? '#52c41a' : '#d9d9d9',
              fontSize: '18px'
            }}
          />
          <Text strong style={{ fontSize: '16px', marginRight: 8 }}>
            รายการสินค้า
          </Text>
          {canEdit ? (
            <Text type="secondary" style={{ fontSize: '14px' }}>
              (เลือกจำนวนสินค้าที่ต้องการคืน)
            </Text>
          ) : (
            <Text type="secondary" style={{ fontSize: '14px' }}>
              (รายการสินค้าที่จะคืน)
            </Text>
          )}
        </Space>
      }
      extra={
        <Input
          placeholder="ค้นหา SKU หรือชื่อสินค้า"
          prefix={<SearchOutlined style={{ color: '#1890ff' }} />}
          onChange={(e) => setSearchText(e.target.value)}
          style={{ width: 250 }}
          allowClear
        />
      }
      style={{ boxShadow: '0 1px 3px rgba(0,0,0,0.05)' }}
    >
      <Table
        columns={columns}
        dataSource={filteredData}
        rowKey="sku"
        pagination={{ 
          pageSize: 5,
          showSizeChanger: true,
          showTotal: (total) => `ทั้งหมด ${total} รายการ`,
          pageSizeOptions: ['5', '10', '20']
        }}
        loading={loading}
        expandable={{
          expandedRowRender,
          expandedRowKeys,
          onExpand: handleRowExpand,
          expandRowByClick: false
        }}
        locale={{
          emptyText: (
            <Empty
              image={Empty.PRESENTED_IMAGE_SIMPLE}
              description={
                searchText ? "ไม่พบรายการที่ค้นหา" : "ไม่มีรายการสินค้า"
              }
            />
          ),
        }}
        onRow={(record) => ({
          style: {
            backgroundColor: getReturnQty(record.sku) > 0 ? 'rgba(82, 196, 26, 0.05)' : undefined,
            cursor: 'pointer'
          },
          onClick: () => handleRowExpand(
            !expandedRowKeys.includes(record.sku),
            record
          ),
        })}
        summary={(pageData) => {
          // คำนวณตัวเลขสรุปสำหรับข้อมูลในหน้าปัจจุบัน
          const totalPrice = pageData.reduce(
            (sum, item) => sum + (Math.abs(item.price) * getReturnQty(item.sku)),
            0
          );
          const totalItems = pageData.reduce(
            (sum, item) => sum + getReturnQty(item.sku),
            0
          );
          const totalQty = pageData.reduce(
            (sum, item) => sum + Math.abs(item.qty),
            0
          );
          
          return (
            <Table.Summary fixed>
              <Table.Summary.Row style={{ backgroundColor: '#fafafa', fontWeight: 'bold' }}>
                <Table.Summary.Cell index={0} colSpan={2}>
                  <Text strong>รวมทั้งหมด</Text>
                </Table.Summary.Cell>
                <Table.Summary.Cell index={2} align="center">
                  <Text>{totalQty} ชิ้น</Text>
                </Table.Summary.Cell>
                <Table.Summary.Cell index={3} align="center">
                  <Text style={{ color: totalItems > 0 ? '#52c41a' : undefined }}>
                    {totalItems} ชิ้น
                  </Text>
                </Table.Summary.Cell>
                <Table.Summary.Cell index={4} align="right">
                  <Text></Text>
                </Table.Summary.Cell>
                <Table.Summary.Cell index={5} align="right">
                  <Text strong style={{ color: totalPrice > 0 ? '#52c41a' : undefined }}>
                    ฿{totalPrice.toLocaleString()}
                  </Text>
                </Table.Summary.Cell>
              </Table.Summary.Row>
            </Table.Summary>
          );
        }}
      />
    </Card>
  );
};

export default OrderItemsSection;