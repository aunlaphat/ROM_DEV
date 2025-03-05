import {
  Layout,
  Card,
  Typography,
  Row,
  Col,
  DatePicker,
  Button,
  Table,
  Space,
  Dropdown,
  ConfigProvider,
} from "antd";
import {
  DownloadOutlined,
  SearchOutlined,
  RedoOutlined,
  BarChartOutlined,
} from "@ant-design/icons";
import { useState } from "react";
import type { Dayjs } from "dayjs";
import * as XLSX from "xlsx";
import { Icon } from "../../resources/icon";

const { Content } = Layout;
const { Title, Text } = Typography;
const { RangePicker } = DatePicker;

const mockData = [
  {
    key: "1",
    SO: "SO123456",
    Tracking_Order: "RT123456",
    SKU: "G090108-EF05",
    SKU_Name: "Bewell Official Store",
    QTY: 20,
    Amount: "599.00",
    Site: "DCOM",
    Warehouse: "RBN",
    Location: "Bangkok",
    Ship_date: "2024-09-01",
    Channel: "ECOM",
  },
  {
    key: "2",
    SO: "SO123457",
    Tracking_Order: "RT123457",
    SKU: "G090108-EF04",
    SKU_Name: "Bewell Shop",
    QTY: 50,
    Amount: "1,299.00",
    Site: "DCOM",
    Warehouse: "RBN",
    Location: "Chiang Mai",
    Ship_date: "2024-09-05",
    Channel: "ECOM",
  },
  {
    key: "3",
    SO: "SO123458",
    Tracking_Order: "RT123458",
    SKU: "G090108-EF06",
    SKU_Name: "Bewell Accessories",
    QTY: 15,
    Amount: "399.00",
    Site: "DCOM",
    Warehouse: "RBN",
    Location: "Phuket",
    Ship_date: "2024-09-10",
    Channel: "ECOM",
  },
  {
    key: "4",
    SO: "SO123459",
    Tracking_Order: "RT123459",
    SKU: "G090108-EF07",
    SKU_Name: "Bewell Health",
    QTY: 30,
    Amount: "899.00",
    Site: "DCOM",
    Warehouse: "RBN",
    Location: "Khon Kaen",
    Ship_date: "2024-09-15",
    Channel: "ECOM",
  },
  {
    key: "5",
    SO: "SO123460",
    Tracking_Order: "RT123460",
    SKU: "G090108-EF08",
    SKU_Name: "Bewell Travel",
    QTY: 10,
    Amount: "499.00",
    Site: "DCOM",
    Warehouse: "RBN",
    Location: "Pattaya",
    Ship_date: "2024-09-20",
    Channel: "ECOM",
  },
  {
    key: "6",
    SO: "SO123461",
    Tracking_Order: "RT123461",
    SKU: "G090108-EF09",
    SKU_Name: "Bewell Sleep",
    QTY: 25,
    Amount: "799.00",
    Site: "DCOM",
    Warehouse: "RBN",
    Location: "Hua Hin",
    Ship_date: "2024-09-25",
    Channel: "ECOM",
  },
];

export const Report = () => {
  const [dates, setDates] = useState<[Dayjs, Dayjs] | null>(null);
  const [filteredData, setFilteredData] = useState(mockData);

  const columns = [
    { title: "SO", dataIndex: "SO", key: "SO" },
    {
      title: "Tracking Order",
      dataIndex: "Tracking_Order",
      key: "Tracking_Order",
    },
    { title: "SKU", dataIndex: "SKU", key: "SKU" },
    { title: "SKU Name", dataIndex: "SKU_Name", key: "SKU_Name" },
    { title: "QTY", dataIndex: "QTY", key: "QTY" },
    { title: "Amount", dataIndex: "Amount", key: "Amount" },
    { title: "Warehouse", dataIndex: "Warehouse", key: "Warehouse" },
    { title: "Ship Date", dataIndex: "Ship_date", key: "Ship_date" },
    { title: "Channel", dataIndex: "Channel", key: "Channel" },
  ];

  const handleSearch = () => {
    if (dates && dates[0] && dates[1]) {
      console.log(
        "🔍 Searching data between",
        dates[0].format("YYYY-MM-DD"),
        "and",
        dates[1].format("YYYY-MM-DD")
      );
      // Filtering logic
    }
  };

  const handleReset = () => {
    setDates(null);
    setFilteredData([]);
  };

  const handleExportExcel = () => {
    const worksheet = XLSX.utils.json_to_sheet(filteredData);
    const workbook = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(workbook, worksheet, "Report Data");
    XLSX.writeFile(workbook, "report_data.xlsx");
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
        Report
      </div>
      <Layout
        style={{
          margin: "24px",
          padding: 20,
          background: "#fff",
          borderRadius: "8px",
        }}
      >
        {/* Header Section
      <Card bordered={false} style={{ marginBottom: 24 }}>
        <Space direction="vertical" size="small">
          <Title level={4}>
            {Icon.analytics({ style: { fontSize: 24, color: "#1890ff" } })} Report
          </Title>
          <Text type="secondary">
            {Icon.chart({ style: { fontSize: 16 } })} แสดงข้อมูลรายงานการคืนสินค้าตามช่วงเวลา
          </Text>
        </Space>
      </Card> */}

        {/* Filter Section */}
        <Card bordered={false} style={{ marginBottom: 24 }}>
          <Row gutter={[16, 16]} align="middle">
            <Col xs={24} sm={12} md={8}>
              <Text strong>ช่วงวันที่</Text>
              <RangePicker
                style={{ width: "100%", marginTop: 8 }}
                value={dates}
                onChange={(dates) => setDates(dates as [Dayjs, Dayjs])}
              />
            </Col>
            <Col xs={24} sm={12} md={16} style={{ textAlign: "right" }}>
              <Space>
                <Button
                  type="primary"
                  icon={<SearchOutlined />}
                  onClick={handleSearch}
                >
                  ค้นหา
                </Button>
                <Button icon={<RedoOutlined />} onClick={handleReset}>
                  รีเซ็ต
                </Button>
                <Dropdown.Button
                  type="primary"
                  icon={<DownloadOutlined />}
                  menu={{
                    items: [
                      {
                        key: "excel",
                        label: "Export Excel",
                        onClick: handleExportExcel,
                      },
                      {
                        key: "csv",
                        label: "Export CSV",
                      },
                    ],
                  }}
                >
                  ส่งออก
                </Dropdown.Button>
              </Space>
            </Col>
          </Row>
        </Card>

        {/* Table Section */}
        <Card bordered={false}>
          <Table
            columns={columns}
            dataSource={filteredData}
            rowKey="SO"
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showTotal: (total) => `ทั้งหมด ${total} รายการ`,
            }}
            scroll={{ x: "max-content" }}
            bordered
          />
        </Card>
      </Layout>
    </ConfigProvider>
  );
};
