import React from 'react';
import { Card, Descriptions, Table, Typography } from 'antd';
import { ReturnOrderState } from '../../../../redux/orders/api';

interface PreviewSectionProps {
  orderData: ReturnOrderState['orderData'];
  returnItems: { [key: string]: number };
}

const PreviewSection: React.FC<PreviewSectionProps> = ({ orderData, returnItems }) => {
  // แปลงข้อมูลสำหรับ Table
  const dataSource = orderData?.lines
    .filter(item => returnItems[item.sku] > 0)
    .map(item => ({
      key: item.sku,
      sku: item.sku,
      itemName: item.itemName,
      returnQty: returnItems[item.sku], // ใช้ค่าจาก returnItems โดยตรง
      price: item.price,
      totalPrice: item.price * returnItems[item.sku]
    }));

  return (
    <Card title="ตรวจสอบข้อมูลการคืนสินค้า">
      <Descriptions title="ข้อมูลการคืนสินค้า" bordered>
        <Descriptions.Item label="SO Number">{orderData?.head.soNo}</Descriptions.Item>
        <Descriptions.Item label="Order Number">{orderData?.head.orderNo}</Descriptions.Item>
        <Descriptions.Item label="SR Number">{orderData?.head.srNo}</Descriptions.Item>
        <Descriptions.Item label="คลังสินค้า">{orderData?.head.locationTo}</Descriptions.Item>
        <Descriptions.Item label="สถานะ SO">{orderData?.head.salesStatus}</Descriptions.Item>
        <Descriptions.Item label="สถานะ MKP">{orderData?.head.mkpStatus}</Descriptions.Item>
      </Descriptions>

      <Table 
        title={() => <Typography.Title level={5}>รายการสินค้าที่จะคืน</Typography.Title>}
        columns={[
          { title: 'SKU', dataIndex: 'sku' },
          { title: 'ชื่อสินค้า', dataIndex: 'itemName' },
          { title: 'จำนวนที่คืน', dataIndex: 'returnQty' },
          { 
            title: 'ราคารวม', 
            render: (_, record: any) => `฿${(record.price * record.returnQty).toLocaleString()}` 
          }
        ]}
        dataSource={dataSource}
        pagination={false}
        summary={(pageData) => {
          const totalAmount = pageData.reduce((sum, item: any) => sum + (item.price * item.returnQty), 0);
          const totalItems = pageData.reduce((sum, item: any) => sum + item.returnQty, 0);
          
          return (
            <Table.Summary.Row>
              <Table.Summary.Cell index={0} colSpan={2}>รวมทั้งหมด</Table.Summary.Cell>
              <Table.Summary.Cell index={2}>{totalItems} ชิ้น</Table.Summary.Cell>
              <Table.Summary.Cell index={3}>฿{totalAmount.toLocaleString()}</Table.Summary.Cell>
            </Table.Summary.Row>
          );
        }}
      />
    </Card>
  );
};

export default PreviewSection;
