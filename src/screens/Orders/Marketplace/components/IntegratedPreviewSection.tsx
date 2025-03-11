// src/screens/Orders/Marketplace/components/IntegratedPreviewSection.tsx
import React, { useEffect, useMemo } from "react";
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
  Button,
  Tooltip,
  Empty
} from "antd";
import {
  FileTextOutlined,
  NumberOutlined,
  ShoppingCartOutlined,
  DollarOutlined,
  CheckCircleOutlined,
  EditOutlined,
  ArrowLeftOutlined,
} from "@ant-design/icons";
import { OrderData } from "../../../../redux/orders/types/state";
import { WAREHOUSES, TRANSPORT_TYPES } from "../../../../constants/warehouse";

const { Title, Text } = Typography;

interface IntegratedPreviewSectionProps {
  orderData: OrderData | null;
  returnItems: { [key: string]: number };
  form: any;
  onEdit: (step: string) => void;
  onNext: () => void;
  loading: boolean;
  stepLoading: boolean;
  getReturnQty: (sku: string) => number;
  isConfirmStep?: boolean;
}

const IntegratedPreviewSection: React.FC<IntegratedPreviewSectionProps> = ({
  orderData,
  returnItems,
  form,
  onEdit,
  onNext,
  loading,
  stepLoading,
  getReturnQty,
  isConfirmStep = false,
}) => {
  // ดึงข้อมูลจากฟอร์ม
  const formValues = form.getFieldsValue();

  // ฟังก์ชันสำหรับป้องกันกรณีการอ่าน returnQty ผิดพลาด
  const safeGetReturnQty = (sku: string): number => {
    try {
      const qty = getReturnQty(sku);
      return typeof qty === 'number' ? qty : 0;
    } catch (error) {
      console.error(`Error getting return qty for ${sku}:`, error);
      // Fallback to returnItems directly
      return returnItems[sku] || 0;
    }
  };

  // แปลงข้อมูลสำหรับตาราง - ใช้ safeGetReturnQty แทน getReturnQty
  const returnItemsList = useMemo(() => {
    if (!orderData?.lines) return [];

    return orderData.lines
      .filter(item => {
        const qty = safeGetReturnQty(item.sku);
        return qty > 0;
      })
      .map(item => {
        const returnQty = safeGetReturnQty(item.sku);
        return {
          key: item.sku,
          sku: item.sku,
          itemName: item.itemName,
          returnQty: returnQty,
          price: Math.abs(item.price),
          totalPrice: Math.abs(item.price) * returnQty,
        };
      });
  }, [orderData, safeGetReturnQty]);

  // ฟังก์ชันสำหรับคำนวณราคารวมทั้งหมด - ใช้ safeGetReturnQty
  const calculateTotalAmount = () => {
    if (!orderData?.lines) return 0;

    return orderData.lines.reduce((sum, item) => {
      const qty = safeGetReturnQty(item.sku);
      return sum + Math.abs(item.price) * qty;
    }, 0);
  };

  // ฟังก์ชันสำหรับคำนวณจำนวนสินค้าที่คืนทั้งหมด - ใช้ safeGetReturnQty
  const calculateTotalItems = () => {
    if (!orderData?.lines) return 0;

    return orderData.lines.reduce(
      (sum, item) => sum + safeGetReturnQty(item.sku),
      0
    );
  };

  // คำนวณค่ารวมใหม่
  const finalTotalAmount = calculateTotalAmount();
  const finalTotalItems = calculateTotalItems();
  const totalTypes = returnItemsList.length;

  // เพิ่ม Debugging
  useEffect(() => {
    console.log(`IntegratedPreviewSection mounted in ${isConfirmStep ? 'confirm' : 'preview'} step`);
    console.log("OrderData:", orderData);
    console.log("ReturnItems received:", returnItems);
    
    if (orderData?.lines) {
      console.log("OrderData lines:", orderData.lines.length);
      
      // ตรวจสอบว่า safeGetReturnQty ทำงานถูกต้องหรือไม่
      const returnItemsCheck = orderData.lines.map(item => ({
        sku: item.sku,
        returnQty: safeGetReturnQty(item.sku),
      }));
      
      console.log("Items with returnQty:", returnItemsCheck.filter(item => item.returnQty > 0));
      console.log("Final values - Total Amount:", finalTotalAmount, "Total Items:", finalTotalItems);
    }
  }, [orderData, returnItems, safeGetReturnQty, isConfirmStep, finalTotalAmount, finalTotalItems]);

  // ค้นหาชื่อคลังสินค้าและขนส่ง
  const getWarehouseName = (id: string) => {
    const warehouse = WAREHOUSES.find((w) => w.value === id);
    return warehouse ? warehouse.label : id;
  };

  const getTransportName = (id: string) => {
    const transport = TRANSPORT_TYPES.find((t) => t.value === id);
    return transport ? transport.label : id;
  };

  return (
    <div className="integrated-preview-section">
      <Row gutter={24} style={{ marginBottom: 24 }}>
        <Col span={24}>
          <Card style={{ marginBottom: 24, textAlign: "center" }}>
            <Title level={4} style={{ marginBottom: 16 }}>
              <Space>
                <CheckCircleOutlined style={{ color: "#52c41a" }} />
                ตรวจสอบข้อมูลการคืนสินค้า
              </Space>
            </Title>
            <Text type="secondary">
              โปรดตรวจสอบข้อมูลให้ถูกต้องก่อนดำเนินการต่อ
              ถ้าข้อมูลไม่ถูกต้องสามารถย้อนกลับไปแก้ไขได้
            </Text>
          </Card>
        </Col>
      </Row>

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
              value={finalTotalItems}
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
              value={finalTotalAmount}
              prefix={<DollarOutlined />}
              suffix="บาท"
              precision={2}
              valueStyle={{ color: "#52c41a" }}
            />
          </Card>
        </Col>
      </Row>

      {/* ข้อมูลรายละเอียด */}
      <Card
        title={
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
            }}
          >
            <Text strong>ข้อมูลการคืนสินค้า</Text>
            <Button
              type="link"
              icon={<EditOutlined />}
              onClick={() => onEdit("create")}
              disabled={loading || stepLoading}
            >
              แก้ไขข้อมูล
            </Button>
          </div>
        }
        style={{ marginBottom: 24 }}
      >
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
            <Text strong style={{ color: "#52c41a" }}>
              {orderData?.head.srNo || "-"}
            </Text>
          </Descriptions.Item>
          <Descriptions.Item label="คลังสินค้า">
            {getWarehouseName(formValues.warehouseFrom) || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="วันที่คืนสินค้า">
            {formValues.returnDate
              ? formValues.returnDate.format("DD/MM/YYYY HH:mm")
              : "-"}
          </Descriptions.Item>
          <Descriptions.Item label="เลขพัสดุ">
            {formValues.trackingNo || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="ประเภทขนส่ง">
            {getTransportName(formValues.transportType) || "-"}
          </Descriptions.Item>
          <Descriptions.Item label="เหตุผลการคืน" span={2}>
            {formValues.reason || "-"}
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

      {/* ส่วนแสดงข้อมูลสินค้าที่คืน */}
      <Card
        title={
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
            }}
          >
            <Space>
              <CheckCircleOutlined />
              <span>รายการสินค้าที่จะคืน</span>
              <Badge
                count={totalTypes}
                style={{ backgroundColor: "#1890ff" }}
              />
            </Space>
            <Button
              type="link"
              icon={<EditOutlined />}
              onClick={() => onEdit("create")}
              disabled={loading || stepLoading}
            >
              แก้ไขรายการสินค้า
            </Button>
          </div>
        }
      >
        {returnItemsList.length > 0 ? (
          <Table
            columns={[
              {
                title: "SKU",
                dataIndex: "sku",
                width: "15%",
                render: (sku) => (
                  <Text copyable style={{ fontFamily: "monospace" }}>
                    {sku}
                  </Text>
                ),
              },
              {
                title: "ชื่อสินค้า",
                dataIndex: "itemName",
                width: "35%",
                ellipsis: true,
              },
              {
                title: "จำนวนที่คืน",
                dataIndex: "returnQty",
                width: "15%",
                align: "center",
                render: (qty) => <Tag color="blue">{qty} ชิ้น</Tag>,
              },
              {
                title: "ราคาต่อชิ้น",
                dataIndex: "price",
                width: "15%",
                align: "right",
                render: (price) => `฿${price.toLocaleString()}`,
              },
              {
                title: "ราคารวม",
                dataIndex: "totalPrice",
                width: "20%",
                align: "right",
                render: (total) => (
                  <Text strong style={{ color: "#52c41a" }}>
                    ฿{total.toLocaleString()}
                  </Text>
                ),
              },
            ]}
            dataSource={returnItemsList}
            pagination={false}
            summary={() => (
              <Table.Summary.Row style={{ backgroundColor: "#fafafa" }}>
                <Table.Summary.Cell index={0} colSpan={2}>
                  <Text strong>รวมทั้งหมด</Text>
                </Table.Summary.Cell>
                <Table.Summary.Cell index={2} align="center">
                  <Text strong>{finalTotalItems} ชิ้น</Text>
                </Table.Summary.Cell>
                <Table.Summary.Cell index={3} align="right">
                  <Text></Text>
                </Table.Summary.Cell>
                <Table.Summary.Cell index={4} align="right">
                  <Text strong style={{ color: "#52c41a", fontSize: "16px" }}>
                    ฿{finalTotalAmount.toLocaleString()}
                  </Text>
                </Table.Summary.Cell>
              </Table.Summary.Row>
            )}
          />
        ) : (
          <Empty 
            description="ไม่พบรายการสินค้าที่ต้องการคืน" 
            image={Empty.PRESENTED_IMAGE_SIMPLE} 
          />
        )}
      </Card>

      {/* ไม่แสดงปุ่มดำเนินการต่อใน confirm step เพราะเราจะแสดงปุ่มยืนยันส่วนกลางแทน */}
      {!isConfirmStep && (
        <div
          style={{ marginTop: 24, display: "flex", justifyContent: "center" }}
        >
          <Space size="middle">
            <Button
              icon={<ArrowLeftOutlined />}
              onClick={() => onEdit("sr")}
              disabled={loading || stepLoading}
            >
              ย้อนกลับ
            </Button>
            <Button
              type="primary"
              icon={<CheckCircleOutlined />}
              onClick={onNext}
              loading={loading || stepLoading}
              size="large"
            >
              ยืนยันข้อมูลและดำเนินการต่อ
            </Button>
          </Space>
        </div>
      )}
    </div>
  );
};

export default IntegratedPreviewSection;