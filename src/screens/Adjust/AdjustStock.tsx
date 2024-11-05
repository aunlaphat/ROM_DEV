import { CloudSyncOutlined } from "@ant-design/icons";
import {
  Button,
  Col,
  ConfigProvider,
  Form,
  Input,
  InputNumber,
  Layout,
  Row,
  Select,
  Table,
  Tooltip,
} from "antd";
import ShopeeLogo from '../../assets/images/shopee logo.png';
import LazadaLogo from '../../assets/images/lazada logo.png';
import BewellLogo from '../../assets/images/Bewelllogonew.jpeg';
import { useState } from "react";

const options = [
  { value: "Not Identified", label: "Not Identified" },
  { value: "G097171-ARM01-BL", label: "G097171-ARM01-BL Bewell Sport armband size M For" },
  { value: "G097171-ARM01-GY", label: "G097171-ARM01-GY Bewell Sport armband size M For" },
  { value: "G097171-ARM02-BL", label: "G097171-ARM02-BL Bewell Sport armband size L" },
  { value: "G097171-ARM03-GR", label: "G097171-ARM03-GR Bewell Sport armband size M with light" },
  { value: "Cancelled", label: "Cancelled" },
];

const initialShopData = [
  { SKU: "G097171-ARM01-BL", Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", Brand: "Bewell" },
  { SKU: "G097171-ARM01-GY", Name: "Bewell Sport armband size M For", Brand: "Bewell" },
  { SKU: "G097171-ARM02-BL", Name: "Sport armband size L", Brand: "Bewell" },
  { SKU: "G097171-ARM03-GR", Name: "Bewell Sport armband size M with light", Brand: "Bewell" },
];

const data = {
  Shopee: [
    { Shopname: "Bewell Official Store", Item_id: "234378629871", Variation_id: "165812913030", Platform_SKU: "G097172-ET02-6F", SKU_Shop: "Link", Stock: "100" },
    { Shopname: "Bewell Official Store", Item_id: "20180005271", Variation_id: "69660249759", Platform_SKU: "G097172-ET02-6F/G097172-SWE06-L", SKU_Shop: "Link", Stock: "120" },
    { Shopname: "Bewell Official Store", Item_id: "234378629871", Variation_id: "165812913030", Platform_SKU: "G097172-ET02-6F", SKU_Shop: "Link", Stock: "130" ,  },
  ],
  Lazada: [
    { Shopname: "Bewell Official Store", Seller_SKU: "G097172-ET02-6F", SKU_Shop: "Link", SQ: "100" },
    { Shopname: "Bewell Official Store", Seller_SKU: "G097172-ET02-6F", SKU_Shop: "Link", SQ: "60" },
    { Shopname: "Bewell Official Store", Seller_SKU: "G097172-ET02-6F", SKU_Shop: "Link", SQ: "70" },
  ],
  Bewell: [
    { Product_id: "61197", Variation_id: "61211", SKU_Shop: "Link", SQ: "40" },
    { Product_id: "61198", Variation_id: "61211", SKU_Shop: "Link", SQ: "30" },
    { Product_id: "61196", Variation_id: "61213", SKU_Shop: "Link", SQ: "20" },
  ]
};

const columns = {
  Shopee: [
    { title: 'Shop Name', dataIndex: 'Shopname', key: 'Shopname' },
    { title: 'Item ID', dataIndex: 'Item_id', key: 'Item_id' },
    { title: 'Variation id', dataIndex: 'Variation_id', key: 'Variation_id' },
    { title: 'Platform SKU', dataIndex: 'Platform_SKU', key: 'Platform_SKU' },
    { title: 'SKU Shop',
      dataIndex: 'SKU_Shop',
      key: 'SKU_Shop',render: () => (
      <a
        href="https://www.bewellstyle.com/product/back-seat-cushion-2/"
        target="_blank"
        rel="noopener noreferrer"
      >
        View Product
      </a>
    ),
    },
    { title: 'Stock', dataIndex: 'Stock', key: 'Stock' },
    { title: 'Status adjust',
      dataIndex: 'Status_adjust',
      key: 'Status_adjust',render: () => (
        <Select
        style={{ width: 150, height: 50,  borderRadius: "50", }}
        showSearch
        placeholder="สถานะปัจจุบัน"
        options={[
          { value: "1", label: <span style={{ color: "#4DA02F" }}>เปิดขาย</span> },
          { value: "2", label: <span style={{ color: "#FF8900" }}>ปิดขายชั่วคราว</span> },
          
        ]}
      />
    ),
    },

  ],
  Lazada: [
    { title: 'Shop Name', dataIndex: 'Shopname', key: 'Shopname' },
    { title: 'Seller SKU', dataIndex: 'Seller_SKU', key: 'Seller_SKU' },
    { title: 'SKU Shop',
      dataIndex: 'SKU_Shop',
      key: 'SKU_Shop',render: () => (
      <a
        href="https://www.bewellstyle.com/product/back-seat-cushion-2/"
        target="_blank"
        rel="noopener noreferrer"
      >
        View Product
      </a>
    ),
    },
    { title: 'Stock Quantity', dataIndex: 'SQ', key: 'SQ' },
    { title: 'Status adjust',
      dataIndex: 'Status_adjust',
      key: 'Status_adjust',render: () => (
        <Select
        style={{ width: 150, height: 50,  borderRadius: "50", }}
        showSearch
        placeholder="สถานะปัจจุบัน"
        options={[
          { value: "1", label: <span style={{ color: "#4DA02F" }}>เปิดขาย</span> },
          { value: "2", label: <span style={{ color: "#FF8900" }}>ปิดขายชั่วคราว</span> },
          
        ]}
      />
    ),
    },
  ],
  Bewell: [
    { title: 'Product ID', dataIndex: 'Product_id', key: 'Product_id' },
    { title: 'Variation ID', dataIndex: 'Variation_id', key: 'Variation_id' },
    { title: 'SKU Shop',
      dataIndex: 'SKU_Shop',
      key: 'SKU_Shop',render: () => (
      <a
        href="https://www.bewellstyle.com/product/back-seat-cushion-2/"
        target="_blank"
        rel="noopener noreferrer"
      >
        View Product
      </a>
    ),
    },
    { title: 'Stock Quantity', dataIndex: 'SQ', key: 'SQ' },
    { title: 'Status adjust',
      dataIndex: 'Status_adjust',
      key: 'Status_adjust',render: () => (
        <Select
        style={{ width: 150, height: 50,  borderRadius: "50", }}
        showSearch
        placeholder="สถานะปัจจุบัน"
        options={[
          { value: "1", label: <span style={{ color: "#4DA02F" }}>เปิดขาย</span> },
          { value: "2", label: <span style={{ color: "#FF8900" }}>ปิดขายชั่วคราว</span> },
          
        ]}
      />
    ),
    },
  ],
};

const Adjust = () => {
  const [form] = Form.useForm();
  const [selectedValue, setSelectedValue] = useState<string | undefined>();
  const [selectedChannel, setSelectedChannel] = useState<string | undefined>();

  const handleSelectChange = (value: string) => {
    const selectedOption = initialShopData.find((val) => val.SKU === value);
    if (selectedOption) {
      form.setFieldsValue({
        SKU: selectedOption.SKU,
        Name: selectedOption.Name,
        Brand: selectedOption.Brand,
      });
      setSelectedValue(value);

      // Set selected channel only if the SKU matches "G097171-ARM01-BL"
      if (value === "G097171-ARM01-BL") {
        // Assuming we want to display all available channels for the selected SKU
        setSelectedChannel("Shopee");  // Default to Shopee, can be adjusted
      } else {
        setSelectedChannel(undefined);  // Hide tables for other SKUs
      }
    }
  };

  return (
    <ConfigProvider>
       <div style={{ marginLeft: '28px', fontSize: '25px', fontWeight: 'bold', color: 'DodgerBlue' }}>
       Adjust Stock
      </div>
      <Layout>
        <Layout.Content
          style={{
            margin: "24px",
            padding: 36,
            minHeight: 360,
            background: "#fff",
            borderRadius: "8px",
          }}
        >
          <Form layout="vertical" style={{ marginTop: "20px" }} form={form}>
            <Row align="middle" justify="start" style={{ width: "100%" }}>
              <Col>
                <Form.Item
                  label="Search Sku Product Name"
                  name="Search Sku Product Name"
                  rules={[{ required: true, message: "Please select the SKU!" }]}
                >
                  <Select
                    showSearch
                    style={{ width: '100%',height:'50px', tableLayout: 'fixed' }} // Ensure the table takes full width and is fixed layout
                   
                    placeholder="Search to Select"
                    optionFilterProp="label"
                    value={selectedValue}
                    onChange={handleSelectChange}
                    filterSort={(optionA, optionB) =>
                      (optionA?.label ?? "")
                        .toLowerCase()
                        .localeCompare((optionB?.label ?? "").toLowerCase())
                    }
                    options={options}
                  />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={20}>
              <Col span={8}>
                <Form.Item
                  label="SKU ID"
                  name="SKU"
                  rules={[{ required: true, message: "Please select the SKU!" }]}
                >
                  <Input disabled={true} style={{ width: "100%", height: 50, color: '#666' }} />
                </Form.Item>
              </Col>

              <Col span={6}>
                <Form.Item
                  label="Brand"
                  name="Brand"
                  rules={[{ required: true,message: "Please select the SKU!" }]}
                  >
                    <Input disabled={true} style={{ width: "100%", height: 50, color: '#666' }} />
                  </Form.Item>
                </Col>
                <Col span={10}>
                  <Form.Item
                    label="Name"
                    name="Name"
                    rules={[{ required: true, message: "Please select the SKU!" }]}
                  >
                    <Input disabled={true} style={{ width: "100%", height: 50, color: '#666' }} />
                  </Form.Item>
                </Col>
              </Row>
  
              <Row gutter={10}>
                <Col span={6}>
                  <Form.Item
                    label="Available Stock"
                    name="Available Stock"
                    rules={[{ required: true, message: "Please select the SKU!" }]}
                  >
                    <Input style={{ width: "100%", height: 50 }} disabled={true} />
                  </Form.Item>
                </Col>
  
                <Col span={6} style={{ display: "flex", justifyContent: "center" }}>
                  <Form.Item
                    label="Minimum Stock"
                    name="minimumStockValue"
                    style={{ flex: 1 }}
                    rules={[
                      {
                        required: true,
                        message: "Please enter the minimum stock!",
                      },
                    ]}
                  >
                    <InputNumber
                      min={1}
                      max={70}
                      style={{ width: "100%", height: "50px", lineHeight: '50px' }}
                    />
                  </Form.Item>
                </Col>
  
                <Col span={6}>
                  <Form.Item
                    label="สถานะปัจจุบัน"
                    name="สถานะปัจจุบัน"
                    rules={[{ required: true, message: "Please select the SKU!" }]}
                  >
                    <Select
                      style={{ width: "100%", height: 50 }}
                      showSearch
                      placeholder="สถานะปัจจุบัน"
                      options={[
                        { value: "1", label: <span style={{ color: "#4DA02F" }}>เปิดขาย</span> },
                        { value: "2", label: <span style={{ color: "#FF8900" }}>ปิดขายชั่วคราว</span> },
                        { value: "3", label: <span style={{ color: "#E8A500" }}>Flash sale</span> },
                        { value: "4", label: <span style={{ color: "#2EA4F9" }}>Clearance</span> },
                      ]}
                    />
                  </Form.Item>
                </Col>
  
                <Col span={6}>
                  <Form.Item
                    label="ซิงค์ร้าน"
                    name="ซิงค์ร้าน"
                  >
                    <Tooltip title="Sync">
                      <Button
                        style={{ width: "100%", height: 50 }}
                        icon={<CloudSyncOutlined style={{ fontSize: "24px", color: "#0F7ECE" }} />}
                      >
                        Sync
                      </Button>
                    </Tooltip>
                  </Form.Item>
                </Col>
              </Row>
            </Form>
          </Layout.Content>
  
          {selectedValue === "G097171-ARM01-BL" && (
            <>
              <Layout.Content
                style={{
                  margin: '24px',
                  padding: 36,
                  minHeight: 360,
                  background: '#fff',
                  borderRadius: '8px',
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 24 }}>
                  <img
                    src={ShopeeLogo}
                    alt="Shopee Logo"
                    style={{ width: '150px', height: 'auto', marginRight: '8px' }}
                  />
                </div>
                <Table
                  dataSource={data.Shopee}
                  columns={columns.Shopee}
                  rowKey="Item_id"
                  pagination={false}
                  style={{ width: '100%', tableLayout: 'fixed' }} // Ensure the table takes full width and is fixed layout
            scroll={{ x: 'max-content' }}
                />
              </Layout.Content>
  
              <Layout.Content
                style={{
                  margin: '24px',
                  padding: 36,
                  minHeight: 360,
                  background: '#fff',
                  borderRadius: '8px',
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 24 }}>
                  <img
                    src={LazadaLogo}
                    alt="Lazada Logo"
                    style={{ width: '150px', height: 'auto', marginRight: '8px' }}
                  />
                </div>
                <Table
                  dataSource={data.Lazada}
                  columns={columns.Lazada}
                  rowKey="Seller_SKU"
                  pagination={false}
                  style={{ width: '100%', tableLayout: 'fixed' }} // Ensure the table takes full width and is fixed layout
                   scroll={{ x: 'max-content' }}
                />
              </Layout.Content>
  
              <Layout.Content
                style={{
                  margin: '24px',
                  padding: 36,
                  minHeight: 360,
                  background: '#fff',
                  borderRadius: '8px',
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 24 }}>
                  <img
                    src={BewellLogo}
                    alt="Bewell Logo"
                    style={{ width: '100px', height: 'auto', marginRight: '8px' }}
                  />
                </div>
                <Table
                  dataSource={data.Bewell}
                  columns={columns.Bewell}
                  rowKey="Product_id"
                  pagination={false}
                  style={{ width: '100%', tableLayout: 'fixed' }} // Ensure the table takes full width and is fixed layout
            scroll={{ x: 'max-content' }}
                />
              </Layout.Content>
            </>
          )}
        </Layout>
      </ConfigProvider>
    );
  };
  
  export default Adjust;
  
