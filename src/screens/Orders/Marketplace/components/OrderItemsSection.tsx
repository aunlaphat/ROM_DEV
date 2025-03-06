import React, { useRef } from 'react';
import { Card, Table, Space, Button, Typography, Tooltip } from 'antd';
import { MinusOutlined, PlusOutlined, QuestionCircleOutlined } from "@ant-design/icons";
import { ReturnOrderState } from '../../../../redux/orders/api';

interface OrderItemsSectionProps {
  orderData: ReturnOrderState['orderData'];
  getReturnQty: (sku: string) => number;
  updateReturnQty: (sku: string, change: number) => void; // แก้ไข type ให้รับ 2 parameters
  loading: boolean;
  currentStep?: 'search' | 'create' | 'sr' | 'preview' | 'confirm'; // เพิ่ม prop currentStep
  returnOrder?: ReturnOrderState['returnOrder']; // เพิ่ม prop returnOrder
}

interface TableItem {
  returnQty: number;
  sku: string;
  itemName: string;
  qty: number;
  price: number;
}

const OrderItemsSection: React.FC<OrderItemsSectionProps> = ({ 
  orderData, 
  getReturnQty, 
  updateReturnQty, 
  loading, 
  currentStep = 'create',
  returnOrder 
}) => {
  // เก็บ cache ด้วย useRef
  const cacheRef = useRef<{
    key: string;
    data: {
      items: TableItem[];
      totals: { items: number; amount: number };
    };
  } | null>(null);

  const processItemsData = React.useCallback(() => {
    // สร้าง key สำหรับ cache
    const dataKey = JSON.stringify({
      step: currentStep,
      returnOrder: returnOrder?.srNo,
      orderData: orderData?.head.srNo,
      items: JSON.stringify(returnOrder?.items || orderData?.lines)  // แก้ไขการสร้าง key
    });

    // ใช้ cache ถ้ามีข้อมูลเดิม
    if (cacheRef.current?.key === dataKey) {
      return cacheRef.current.data;
    }

    const processedItems: TableItem[] = (returnOrder?.items || orderData?.lines || []).map(record => ({
      sku: record.sku,
      itemName: record.itemName,
      qty: Math.abs(record.qty),
      returnQty: returnOrder ? record.returnQty : getReturnQty(record.sku),
      price: Math.abs(record.price)
    }));

    const totals = processedItems.reduce((acc, item) => ({
      items: acc.items + item.returnQty,
      amount: acc.amount + (item.price * item.returnQty)
    }), { items: 0, amount: 0 });

    // แสดง log เฉพาะเมื่อข้อมูลเปลี่ยนแปลงจริงๆ
    console.group(`🔄 Return Items [${currentStep.toUpperCase()}]`);
    console.log('📋 Order Info:', {
      'Order': orderData?.head.orderNo,
      'SR': orderData?.head.srNo || '-',
      'Source': returnOrder ? 'DB' : 'Original',
      'Step': currentStep,
    });

    // อัพเดท cache
    cacheRef.current = {
      key: dataKey,
      data: { items: processedItems, totals }
    };

    return cacheRef.current.data;
  }, [returnOrder, orderData, getReturnQty, currentStep]);

  // ล้าง cache เมื่อ step หรือข้อมูลหลักเปลี่ยน
  React.useEffect(() => {
    cacheRef.current = null;
  }, [currentStep, returnOrder?.srNo, orderData?.head.srNo]);

  // เพิ่ม cleanup effect
  React.useEffect(() => {
    return () => {
      cacheRef.current = null;  // ล้าง cache เมื่อ component unmount
    };
  }, []);

  // ดึงข้อมูลโดยตรงจาก processItemsData
  const { items, totals } = processItemsData();

  const columns = [
    { title: 'SKU', dataIndex: 'sku', width: '15%' },
    { title: 'Item Name', dataIndex: 'itemName', width: '30%' },
    { 
      title: 'QTY', 
      width: '15%',
      align: 'center' as const,
      render: (_: unknown, record: TableItem) => record.qty
    },
    { 
      title: 'ReturnQTY', 
      width: '15%',
      align: 'center' as const,
      render: (_: unknown, record: TableItem) => {
        if (currentStep === 'create') {
          return (
            <Space>
              <Button
                type="text"
                icon={<MinusOutlined />}
                onClick={(e) => {
                  e.stopPropagation(); // ป้องกันการ propagate event
                  updateReturnQty(record.sku, -1);
                }}
                disabled={record.returnQty <= 0}
                danger={record.returnQty > 0}
              />
              <Typography.Text>{record.returnQty}</Typography.Text>
              <Button
                type="text"
                icon={<PlusOutlined />}
                onClick={(e) => {
                  e.stopPropagation(); // ป้องกันการ propagate event
                  updateReturnQty(record.sku, 1);
                }}
                disabled={record.returnQty >= record.qty}
                style={{ color: record.returnQty < record.qty ? '#52c41a' : '#d9d9d9' }}
              />
            </Space>
          );
        }
        return <Typography.Text strong>{record.returnQty}</Typography.Text>;
      }
    },
    { 
      title: 'Price', 
      width: '15%',
      align: 'right' as const,
      render: (_: unknown, record: TableItem) => `฿${(record.price * record.returnQty).toLocaleString()}`
    }
  ];

  const TableSummary = React.useMemo(() => (
    <Table.Summary.Row>
      <Table.Summary.Cell index={0} colSpan={3}>
        <Typography.Text strong>รวมทั้งหมด</Typography.Text>
      </Table.Summary.Cell>
      <Table.Summary.Cell index={1} align="center">
        <Typography.Text strong>{totals.items} ชิ้น</Typography.Text>
      </Table.Summary.Cell>
      <Table.Summary.Cell index={2} align="right" colSpan={2}>
        <Typography.Text strong type="danger">
          ฿{totals.amount.toLocaleString()}
        </Typography.Text>
      </Table.Summary.Cell>
    </Table.Summary.Row>
  ), [totals]);

  return (
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
        dataSource={items}
        rowKey="sku"
        pagination={false}
        loading={loading}
        summary={() => TableSummary}
        onRow={(record) => ({
          style: {
            backgroundColor: record.returnQty === 0 ? '#fafafa' : 'inherit',
            opacity: orderData?.head.srNo ? 0.8 : 1
          }
        })}
      />
    </Card>
  );
};

export default OrderItemsSection;
