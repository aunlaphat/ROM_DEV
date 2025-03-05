import React, { useState } from "react";
import {
  Layout,
  Tabs,
  DatePicker,
  Input,
  Button,
  Table,
  Space,
  Typography,
  Card,
  Row,
  Col,
  Statistic,
  ConfigProvider,
} from "antd";
import {
  SearchOutlined,
  DownloadOutlined,
  HomeOutlined,
  ProjectOutlined,
  TeamOutlined,
  ShopOutlined,
} from "@ant-design/icons";
import { useSelector } from "react-redux";
import type { ColumnsType } from "antd/es/table";
import { Icon } from "../../resources/icon";

const { Content } = Layout;
const { RangePicker } = DatePicker;
const { Title, Link, Text } = Typography;

// 🔹 Mock Data
const saleReturnData = [
  {
    key: "1",
    order: "12345678",
    soInv: "SOC2407-12345",
    customerNo: "TC-NMI-0007",
    sr: "NULL",
    returnOrder: "532453245",
    returnTracking: "ABCDEF",
  },
];

const ijReturnData = [
  {
    key: "1",
    refIj: "12345678",
    ij: "IJ24091234",
    returnTracking: "ABCDEF",
    transport: "J&T",
    date: "24-07-2567",
    warehouse: "RBN",
  },
];

// 🔹 Columns Definition
const saleReturnColumns: ColumnsType<any> = [
  {
    title: "Order",
    dataIndex: "order",
    render: (text) => <Link>{text}</Link>,
  },
  { title: "SO/INV", dataIndex: "soInv" },
  { title: "Customer No", dataIndex: "customerNo" },
  { title: "SR", dataIndex: "sr" },
  { title: "Return Order Number", dataIndex: "returnOrder" },
  { title: "Return Tracking", dataIndex: "returnTracking" },
];

const ijReturnColumns: ColumnsType<any> = [
  {
    title: "Ref IJ",
    dataIndex: "refIj",
    render: (text) => <Link>{text}</Link>,
  },
  { title: "IJ", dataIndex: "ij" },
  { title: "Return Tracking", dataIndex: "returnTracking" },
  { title: "Transport", dataIndex: "transport" },
  { title: "Date", dataIndex: "date" },
  { title: "Warehouse", dataIndex: "warehouse" },
];

const Home: React.FC = () => {
  const [dateRange, setDateRange] = useState([]);
  const [searchValue, setSearchValue] = useState("");
  const user = useSelector((state: any) => state.auth.user);

  // const stats = [
  //   {
  //     title: "รายการรอดำเนินการ",
  //     value: 25,
  //     icon: Icon.pending({ style: { fontSize: 24, color: "#1890ff" } }),
  //     color: "#1890ff",
  //   },
  //   {
  //     title: "ผู้ใช้งานทั้งหมด",
  //     value: 150,
  //     icon: Icon.team({ style: { fontSize: 24, color: "#52c41a" } }),
  //     color: "#52c41a",
  //   },
  //   {
  //     title: "รายการออเดอร์ที่เข้าวันนี้",
  //     value: 999,
  //     icon: Icon.BarCode({ style: { fontSize: 24, color: "#722ed1" } }),
  //     color: "#722ed1"
  //   }
  // ];

  const handleSearch = () => {
    console.log("Search:", { dateRange, searchValue });
  };

  return (
    <ConfigProvider>
      <div
        style={{
          marginLeft: "28px",
          fontSize: "25px",
          fontWeight: "bold",
          color: "DodgerBlue",
        }}
      >
        Home
      </div>
      <Layout
        style={{
          margin: "24px",
          padding: 20,
          background: "#fff",
          borderRadius: "8px",
        }}
      >
        {/* Welcome Section
      <Card bordered={false} style={{ marginBottom: 24 }}>
        <Space direction="vertical" size="small">
          <Title level={4}>
            <HomeOutlined /> Home Page
          </Title>
          <Text>ยินดีต้อนรับ, คุณ {user?.fullName}</Text>
          <Text type="secondary">บทบาท: {user?.roleName}</Text>
        </Space>
      </Card> */}

        {/* Statistics Section
      <Row gutter={[24, 24]}>
        {stats.map((stat, index) => (
          <Col xs={24} sm={12} lg={8} key={index}>
            <Card bordered={false}>
              <Statistic
                title={
                  <Space>
                    {stat.icon}
                    <Text strong>{stat.title}</Text>
                  </Space>
                }
                value={stat.value}
                valueStyle={{ color: stat.color }}
              />
            </Card>
          </Col>
        ))}
      </Row> */}

        <Content>
          {/* 🔹 Tabs */}
          <Tabs defaultActiveKey="1" type="card">
            <Tabs.TabPane tab="Blind Return" key="1" />
            <Tabs.TabPane tab="Booked Return" key="2" />
            <Tabs.TabPane tab="Waiting Action" key="3" />
            <Tabs.TabPane tab="Unsuccess" key="4" />
            <Tabs.TabPane tab="Success" key="5" />
          </Tabs>

          {/* 🔹 Search Section */}
          <Card style={{ marginBottom: 16 }}>
            <Row gutter={[16, 16]} align="middle">
              <Col>
                <RangePicker onChange={(dates) => setDateRange(dates as any)} />
              </Col>
              <Col>
                <Input
                  placeholder="ค้นหา Order/Ref IJ"
                  value={searchValue}
                  onChange={(e) => setSearchValue(e.target.value)}
                  allowClear
                  prefix={<SearchOutlined />}
                />
              </Col>
              <Col>
                <Button type="primary" onClick={handleSearch}>
                  Search
                </Button>
              </Col>
            </Row>
          </Card>

          {/* 🔹 Sale Return Table */}
          <Card
            title="Sale Return"
            extra={
              <Button type="primary" icon={<DownloadOutlined />}>
                Export
              </Button>
            }
          >
            <Table
              columns={saleReturnColumns}
              dataSource={saleReturnData}
              pagination={false}
            />
          </Card>

          {/* 🔹 IJ Return Table */}
          <Card
            title="IJ Return"
            style={{ marginTop: 16 }}
            extra={
              <Button type="primary" icon={<DownloadOutlined />}>
                Export
              </Button>
            }
          >
            <Table
              columns={ijReturnColumns}
              dataSource={ijReturnData}
              pagination={false}
            />
          </Card>
        </Content>
      </Layout>
    </ConfigProvider>
  );
};

export default Home;
