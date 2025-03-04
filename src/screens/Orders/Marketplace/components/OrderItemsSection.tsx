import React from 'react';
import { Card, Table, Space, Button, Typography, Tooltip } from 'antd';
import { MinusOutlined, PlusOutlined, QuestionCircleOutlined } from "@ant-design/icons";
import { ReturnOrderState } from '../../../../redux/orders/api';

interface OrderItemsSectionProps {
  orderData: ReturnOrderState['orderData'];
  getReturnQty: (sku: string) => number;
  updateReturnQty: (sku: string, change: number) => void;
  loading: boolean;
}

const OrderItemsSection: React.FC<OrderItemsSectionProps> = ({ orderData, getReturnQty, updateReturnQty, loading }) => {
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
};

export default OrderItemsSection;
