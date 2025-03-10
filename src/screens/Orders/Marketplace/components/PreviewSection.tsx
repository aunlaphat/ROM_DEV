// src/screens/Orders/Marketplace/components/PreviewSection.tsx

import React from "react";
import {
  Card,
  Descriptions,
  Table,
  Typography,
  Divider,
  Row,
  Col,
  Statistic,
  Tag,
  Space,
  Badge,
} from "antd";
import {
  FileTextOutlined,
  NumberOutlined,
  ShoppingCartOutlined,
  DollarOutlined,
  CheckCircleOutlined,
} from "@ant-design/icons";
import { ReturnOrderState } from "../types";

const { Title, Text } = Typography;

// interface สำหรับ ข้อมูลแต่ละรายการใน dataSource
interface OrderItemData {
  key: string;
  sku: string;
  itemName: string;
  returnQty: number;
  price: number;
  totalPrice: number;
}

interface PreviewSectionProps {
  orderData: ReturnOrderState["orderData"];
  returnItems: { [key: string]: number };
}

const PreviewSection: React.FC<PreviewSectionProps> = ({
  orderData,
  returnItems,
}) => {
  // แปลงข้อมูลสำหรับ Table และกำหนดค่าเริ่มต้นเป็นอาร์เรย์ว่างหากมีค่าเป็น undefined
  const dataSource: OrderItemData[] = orderData?.lines
    ? orderData.lines
        .filter((item) => returnItems[item.sku] > 0)
        .map((item) => ({
          key: item.sku,
          sku: item.sku,
          itemName: item.itemName,
          returnQty: returnItems[item.sku],
          price: item.price,
          totalPrice: item.price * returnItems[item.sku],
        }))
    : [];

  // คำนวณมูลค่ารวม - ระบุประเภทข้อมูลของ sum และ item
  const totalAmount = dataSource.reduce(
    (sum: number, item: OrderItemData) => sum + item.totalPrice,
    0
  );
  
  // คำนวณจำนวนรวม - ระบุประเภทข้อมูลของ sum และ item
  const totalItems = dataSource.reduce(
    (sum: number, item: OrderItemData) => sum + item.returnQty, 
    0
  );
  
  // จำนวนรายการ
  const totalTypes = dataSource.length;

  return (
    <div>
      <Title level={4}>
        <Space>
          <FileTextOutlined />
          ตรวจสอบข้อมูลการคืนสินค้า
        </Space>
      </Title>

      <Divider />

      {/* ส่วนสรุปข้อมูล */}
      <Row gutter={24} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={6}>
          <Card size="small" style={{ textAlign: "center", height: "100%" }}>
            <Statistic
              title="เลขที่คำสั่งซื้อ"
              value={orderData?.head.soNo || "-"}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: "#1890ff" }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card size="small" style={{ textAlign: "center", height: "100%" }}>
            <Statistic
              title="เลข SR Number"
              value={orderData?.head.srNo || "-"}
              prefix={<NumberOutlined />}
              valueStyle={{ color: "#52c41a" }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card size="small" style={{ textAlign: "center", height: "100%" }}>
            <Statistic
              title="จำนวนสินค้าที่คืน"
              value={totalItems}
              suffix="ชิ้น"
              prefix={<ShoppingCartOutlined />}
              valueStyle={{ color: "#fa8c16" }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card size="small" style={{ textAlign: "center", height: "100%" }}>
            <Statistic
              title="มูลค่ารวม"
              value={totalAmount}
              prefix={<DollarOutlined />}
              suffix="บาท"
              precision={2}
              valueStyle={{ color: "#52c41a" }}
            />
          </Card>
        </Col>
      </Row>

      {/* ข้อมูลรายละเอียด */}
      <Card title="ข้อมูลการคืนสินค้า" style={{ marginBottom: 24 }}>
        <Descriptions
          bordered
          column={{ xxl: 3, xl: 3, lg: 3, md: 2, sm: 1, xs: 1 }}
        >
          <Descriptions.Item label="เลขที่ Order">
            {orderData?.head.orderNo || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="เลขที่ SO">
            {orderData?.head.soNo || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="เลข SR">
            {orderData?.head.srNo || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="คลังสินค้า">
            {orderData?.head.locationTo || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="สถานะ SO">
            <Tag
              color={
                orderData?.head.salesStatus === "open order" ? "green" : "blue"
              }
            >
              {orderData?.head.salesStatus || "-"}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="สถานะ MKP">
            <Tag
              color={
                orderData?.head.mkpStatus === "complete" ? "blue" : "orange"
              }
            >
              {orderData?.head.mkpStatus || "-"}
            </Tag>
          </Descriptions.Item>
        </Descriptions>
      </Card>

      {/* รายการสินค้าที่คืน */}
      <Card
        title={
          <Space>
            <CheckCircleOutlined />
            <span>รายการสินค้าที่จะคืน</span>
            <Badge count={totalTypes} style={{ backgroundColor: "#1890ff" }} />
          </Space>
        }
      >
        <Table
          columns={[
            {
              title: "SKU",
              dataIndex: "sku",
              render: (sku) => (
                <Text copyable style={{ fontFamily: "monospace" }}>
                  {sku}
                </Text>
              ),
            },
            {
              title: "ชื่อสินค้า",
              dataIndex: "itemName",
              ellipsis: true,
            },
            {
              title: "จำนวนที่คืน",
              dataIndex: "returnQty",
              align: "center",
              render: (qty) => <Tag color="blue">{qty} ชิ้น</Tag>,
            },
            {
              title: "ราคาต่อชิ้น",
              dataIndex: "price",
              align: "right",
              render: (price) => `฿${Math.abs(price).toLocaleString()}`,
            },
            {
              title: "ราคารวม",
              dataIndex: "totalPrice",
              align: "right",
              render: (total) => (
                <Text strong style={{ color: "#52c41a" }}>
                  ฿{Math.abs(total).toLocaleString()}
                </Text>
              ),
            },
          ]}
          dataSource={dataSource}
          pagination={false}
          summary={() => (
            <Table.Summary.Row style={{ backgroundColor: "#fafafa" }}>
              <Table.Summary.Cell index={0} colSpan={2}>
                <Text strong>รวมทั้งหมด</Text>
              </Table.Summary.Cell>
              <Table.Summary.Cell index={2} align="center">
                <Text strong>{totalItems} ชิ้น</Text>
              </Table.Summary.Cell>
              <Table.Summary.Cell index={3} align="right">
                <Text></Text>
              </Table.Summary.Cell>
              <Table.Summary.Cell index={4} align="right">
                <Text strong style={{ color: "#52c41a", fontSize: "16px" }}>
                  ฿{Math.abs(totalAmount).toLocaleString()}
                </Text>
              </Table.Summary.Cell>
            </Table.Summary.Row>
          )}
        />
      </Card>
    </div>
  );
};

export default PreviewSection;