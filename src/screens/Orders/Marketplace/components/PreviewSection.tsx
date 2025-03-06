import React from 'react';
import { Card, Descriptions, Table, Typography, Space, Divider } from 'antd';
import { ReturnOrderState } from '../../../../redux/orders/api';

interface PreviewSectionProps {
  orderData: ReturnOrderState['orderData'];
  returnItems: { [key: string]: number };
  returnOrder?: ReturnOrderState['returnOrder']; // เพิ่ม prop returnOrder
}

const PreviewSection: React.FC<PreviewSectionProps> = ({ orderData, returnItems, returnOrder }) => {
  // ปรับปรุงการสร้าง dataSource โดยใช้ข้อมูลจาก returnOrder ถ้ามี
  const dataSource = (returnOrder?.items || [])
    .filter(item => item.returnQty > 0)
    .map(item => ({
      key: item.sku,
      sku: item.sku,
      itemName: item.itemName,
      originalQty: Math.abs(item.qty),
      returnQty: item.returnQty,
      price: Math.abs(item.price),
      totalPrice: Math.abs(item.price) * item.returnQty
    }));

  return (
    <Card>
      <Typography.Title level={4}>ตรวจสอบข้อมูลการคืนสินค้า</Typography.Title>
      <Divider />

      {/* ส่วนที่ 1: ข้อมูลหลัก */}
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        <Card type="inner" title="ข้อมูลเอกสาร">
          <Descriptions column={3}>
            <Descriptions.Item label="SO Number" span={1}>{orderData?.head.soNo}</Descriptions.Item>
            <Descriptions.Item label="Order Number" span={1}>{orderData?.head.orderNo}</Descriptions.Item>
            <Descriptions.Item label="SR Number" span={1}>
              <Typography.Text strong type="success">
                {orderData?.head.srNo}
              </Typography.Text>
            </Descriptions.Item>
          </Descriptions>
        </Card>

        {/* ส่วนที่ 2: ข้อมูลการจัดส่ง */}
        <Card type="inner" title="ข้อมูลการคืนสินค้า">
          <Descriptions column={2}>
            <Descriptions.Item label="คลังสินค้า">{orderData?.head.locationTo}</Descriptions.Item>
            <Descriptions.Item label="สถานะ">
              <Space>
                SO: <Typography.Text type="warning">{orderData?.head.salesStatus}</Typography.Text>
                MKP: <Typography.Text type="warning">{orderData?.head.mkpStatus}</Typography.Text>
              </Space>
            </Descriptions.Item>
          </Descriptions>
        </Card>

        {/* ส่วนที่ 3: รายการสินค้า */}
        <Card type="inner" title="รายการสินค้าที่คืน">
          <Table 
            dataSource={dataSource}
            columns={[
              { title: 'SKU', dataIndex: 'sku', width: '15%' },
              { title: 'ชื่อสินค้า', dataIndex: 'itemName', width: '35%' },
              { 
                title: 'จำนวนเดิม', 
                dataIndex: 'originalQty',
                width: '10%',
                align: 'center',
              },
              { 
                title: 'จำนวนที่คืน', 
                dataIndex: 'returnQty',
                width: '15%',
                align: 'center',
                render: (qty) => (
                  <Typography.Text strong type="warning">
                    {qty}
                  </Typography.Text>
                )
              },
              { 
                title: 'ราคาต่อหน่วย',
                dataIndex: 'price',
                width: '15%',
                align: 'right',
                render: (price) => `฿${price.toLocaleString()}`
              },
              { 
                title: 'ราคารวม',
                width: '15%',
                align: 'right',
                render: (_, record) => (
                  <Typography.Text strong>
                    ฿{record.totalPrice.toLocaleString()}
                  </Typography.Text>
                )
              }
            ]}
            pagination={false}
            summary={(pageData) => {
              const totalAmount = pageData.reduce((sum, item) => sum + item.totalPrice, 0);
              const totalItems = pageData.reduce((sum, item) => sum + item.returnQty, 0);
              
              return (
                <>
                  <Table.Summary.Row>
                    <Table.Summary.Cell index={0} colSpan={2}>
                      <Typography.Text strong>รวมทั้งหมด</Typography.Text>
                    </Table.Summary.Cell>
                    <Table.Summary.Cell index={2} align="center">
                      <Typography.Text strong>{totalItems} ชิ้น</Typography.Text>
                    </Table.Summary.Cell>
                    <Table.Summary.Cell index={3} colSpan={2} align="right">
                      <Typography.Text strong type="danger">
                        ฿{totalAmount.toLocaleString()}
                      </Typography.Text>
                    </Table.Summary.Cell>
                  </Table.Summary.Row>
                </>
              );
            }}
          />
        </Card>
      </Space>
    </Card>
  );
};

export default PreviewSection;
