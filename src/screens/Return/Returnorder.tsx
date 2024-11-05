import { Button, Card, Col, ConfigProvider, Divider, Form, Input, Layout, List, Row, Table, message } from "antd";
import React, { useEffect, useState } from "react";
import { Tabs, DatePicker } from "antd";
import dayjs, { Dayjs } from "dayjs";
import isSameOrAfter from "dayjs/plugin/isSameOrAfter";
import isSameOrBefore from "dayjs/plugin/isSameOrBefore";
import isBetween from "dayjs/plugin/isBetween";
import Popup from "reactjs-popup";
import { CloseOutlined, CheckOutlined, FileImageOutlined } from '@ant-design/icons'; // Import reactjs-popup
import '../Return.css'
import { Image } from 'antd';

dayjs.extend(isSameOrAfter);
dayjs.extend(isSameOrBefore);
dayjs.extend(isBetween);

interface Order {
  Ordernumber: string;
  Date: string;
}


const OtherReturn = () => {
  const [form] = Form.useForm();

  const { TabPane } = Tabs;
  const { RangePicker } = DatePicker;
  const [activeTabKey, setActiveTabKey] = useState<string>("1");

  const [data] = useState<Order[]>([
    { Ordernumber: "12345678", Date: "2024-09-01" },
    { Ordernumber: "12345677", Date: "2024-09-15" },
    { Ordernumber: "12345676", Date: "2024-09-29" },
  ]);

  const [dataij] = useState<Order[]>([
    { Ordernumber: "11123456", Date: "2024-09-01" },
    { Ordernumber: "11123457", Date: "2024-09-15" },
    { Ordernumber: "11123458", Date: "2024-09-29" },
  ]);
  const detailorder = [
    {
      SO_INV: "SOC2407-12345",
      ReturnTracking: "ABCDEFG",
      ReturnOrderNumber: "532453245",
      Customer: "TC-NMI-0007",
      SR: "NULL",
      Channel: "Other",
      Warehouse: "RBN",
      StatusWH: "Waiting Return",
      StatusAcc: "Done",
      Date: "2024-09-02",
      ReturnReason: "-"
    },
    ,
    
  ];
  const fields = [
    { label: "SO/INV", name: "SO_INV" },
    { label: "Return Tracking", name: "ReturnTracking" },
    { label: "Return Order Number", name: "ReturnOrderNumber" },
    { label: "Customer", name: "Customer" },
    { label: "SR", name: "SR" },
    { label: "Channel", name: "Channel" },
    { label: "Warehouse", name: "Warehouse" },
    { label: "StatusWH", name: "StatusWH" },
    { label: "Status Acc", name: "StatusAcc" },
    { label: "Date", name: "Date" },
    { label: "Return Reason", name: "ReturnReason" },
  ];
  const fieldPairs = [];
  for (let i = 0; i < fields.length; i += 2) {
    fieldPairs.push(fields.slice(i, i + 2));
  }

  
  const datacolumn = [
    {
      SKU: "G090108-EF05", QTY: "1", จำนวนเข้ารับจริง: "1", Warehouse: "Ecom", Location: "HOLDSTAY",
    },
    {
      SKU: "G090108-EF04", QTY: "1", จำนวนเข้ารับจริง: "1", Warehouse: "Ecom", Location: "HOLDSTAY",
    },
    {
      SKU: "G090108-EF03", QTY: "1", จำนวนเข้ารับจริง: "1", Warehouse: "Ecom", Location: "HOLDSTAY",
    },
    ,
  ];
  const datacolumn1 = [
    {
      SKU: "", QTY: "", จำนวนเข้ารับจริง: "", Warehouse: "", Location: "",
    },

    ,
  ];

  const columns = [
    {
      title: 'SKU',
      dataIndex: 'SKU',
      key: 'sku',
    },
    {
      title: 'QTY',
      dataIndex: 'QTY',
      key: 'qty',
    },
    {
      title: 'จำนวนเข้ารับจริง',
      dataIndex: 'จำนวนเข้ารับจริง',
      key: 'จำนวนเข้ารับจริง',
    },
    {
      title: 'Warehouse',
      dataIndex: 'Warehouse',
      key: 'warehouse',
    },
    {
      title: 'Location',
      dataIndex: 'Location',
      key: 'location',
    },


  ];






  const [filteredData, setFilteredData] = useState<Order[]>(data);
  const [filteredDataIJ, setFilteredDataIJ] = useState<Order[]>(dataij);
  const [dates, setDates] = useState<[Dayjs, Dayjs] | null>(null);
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);

  const handleCopy = (item: string) => {
    navigator.clipboard.writeText(item);
    message.success(`Copied: ${item}`);
  };

  const onTabChange = (key: string) => {
    setActiveTabKey(key);
  };

  const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null, dateStrings: string[]) => {
    if (dates) {
      setDates(dates as [Dayjs, Dayjs]);
    }
  };

  const handleSearch = () => {
    if (dates && dates[0] && dates[1]) {
      const startDate = dates[0].startOf("day");
      const endDate = dates[1].endOf("day");

      const filtered = data.filter((item) => {
        const itemDate = dayjs(item.Date);
        return itemDate.isBetween(startDate, endDate, null, "[]");
      });

      const filteredIJ = dataij.filter((item) => {
        const itemDate = dayjs(item.Date);
        return itemDate.isBetween(startDate, endDate, null, "[]");
      });

      setFilteredData(filtered);
      setFilteredDataIJ(filteredIJ);
    } else {
      setFilteredData(data);
      setFilteredDataIJ(dataij);
    }
  };

  const handleOrderClick = (order: Order) => {
    setSelectedOrder(order);
  };

  const closeModal = () => setSelectedOrder(null);

  return (
    <ConfigProvider>
      <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
      Home
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
          <Tabs
            onChange={onTabChange}
            type="card"
            items={[
              { label: "Blind Return", key: "1" },
              { label: "Booked Return", key: "2" },
              { label: "Waiting Action", key: "3" },
              { label: "Unsuccess", key: "4" },
              { label: "Success", key: "5" },
            ]}
          />

          {activeTabKey === "1" && (
            <>
              <Row gutter={8} align="middle" justify="center" style={{ marginTop: "20px" }}>
                <Col>
                  <Form.Item
                    layout="vertical"
                    label="Select date"
                    name="Select date"
                    rules={[{ required: true, message: "Please select the Select date!" }]}
                  >
                    <RangePicker
                      value={dates}
                      style={{ height: "40px" }}
                      onChange={handleDateChange}
                    />
                  </Form.Item>
                </Col>
                <Col style={{ marginTop: "4px" }}>
                  <Button
                    type="primary"
                    style={{ height: "40px", width: "100px", background: "#32ADE6" }}
                    onClick={handleSearch}
                  >
                    Search
                  </Button>
                </Col>
              </Row>

              <Divider orientation="left">Sale Return</Divider>

              <List
                size="large"
                bordered
                dataSource={filteredData}
                renderItem={(item) => (
                  <List.Item key={item.Ordernumber}>
                    <code
                      style={{ cursor: 'pointer' }}
                      onClick={() => handleOrderClick(item)}
                    >
                      Order Number: {item.Ordernumber}, Date: {item.Date}
                    </code>
                    <Button
                      type="link"
                      onClick={() => handleCopy(item.Ordernumber)}
                    >
                      Copy
                    </Button>
                  </List.Item>
                )}
              />

              <Divider orientation="left">IJ Return</Divider>

              <List
                size="large"
                bordered
                dataSource={filteredDataIJ}
                renderItem={(item) => (
                  <List.Item key={item.Ordernumber}>
                    <code
                      style={{ cursor: 'pointer' }}
                      onClick={() => handleOrderClick(item)}
                    >
                      Order Number: {item.Ordernumber}, Date: {item.Date}
                    </code>
                    <Button
                      type="link"
                      onClick={() => handleCopy(item.Ordernumber)}
                    >
                      Copy
                    </Button>
                  </List.Item>
                )}
              />
            </>
          )}

          {activeTabKey === '2' && (
            <div style={{ textAlign: 'center', marginTop: '50px', fontSize: '18px', color: 'grey' }}>
              ไม่มีข้อมูลหน้า 2
            </div>
          )}
          {activeTabKey === '3' && (
            <div style={{ textAlign: 'center', marginTop: '50px', fontSize: '18px', color: 'grey' }}>
              ไม่มีข้อมูลหน้า 3
            </div>
          )}
          {activeTabKey === '4' && (
            <div style={{ textAlign: 'center', marginTop: '50px', fontSize: '18px', color: 'grey' }}>
              ไม่มีข้อมูลหน้า 4
            </div>
          )}
          {activeTabKey === '5' && (
            <div style={{ textAlign: 'center', marginTop: '50px', fontSize: '18px', color: 'grey' }}>
              ไม่มีข้อมูลหน้า 5
            </div>
          )}

          {/* Modal Popup */}
          <Popup
            open={!!selectedOrder}
            closeOnDocumentClick={false}
            onClose={closeModal}
            modal
            overlayStyle={{
              background: 'rgba(30, 30, 30, 0.4)', // Set background color to #1E1E1E with 40% opacity
            }}
          >
            <div
              style={{
                padding: '20px',
                textAlign: 'center',
                borderRadius: '10px',
                background: '#fff',
                width: '100%', // Set width to 80% of the viewport
                height: '100%', // Set height to 80% of the viewport
                maxWidth: '1200px', // Maximum width
                maxHeight: '750px', // Maximum height
                position: 'fixed', // Fixed position
                top: '50%', // Center vertically
                left: '50%', // Center horizontally
                transform: 'translate(-50%, -50%)', // Adjust position to center
              }}
            >
              {selectedOrder && (
                <Row justify="space-between" align="middle">
                  <Row style={{ marginTop: '10px', width: '50%' }}>
                    <Col span={10}>
                      <Button style={{ fontSize: '15px', color: '#35465B', background: '#D9D9D9', height: '40px', marginTop: '10px' }}>
                        Order Number: {selectedOrder.Ordernumber}
                      </Button>
                    </Col>

                  </Row>
                  <Col>
                    <Button
                      onClick={closeModal}
                      type="text" // ใช้ type="text" เพื่อให้ดูเหมือนไอคอน
                      icon={<CloseOutlined style={{ fontSize: '30px' }} />} // ใช้ไอคอน CloseOutlined
                    />
                  </Col>
                </Row>
              )}

              {selectedOrder && (
                <Tabs defaultActiveKey="1" style={{ marginTop: '20px' }}>
                  <TabPane tab="Order Detail" key="1">
                    <Row justify="space-between" align="middle">
                      <Col span={12} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                        <div
                          style={{
                            width: '90%', // ปรับเป็นเปอร์เซ็นต์เพื่อให้ยืดหยุ่น
                            maxWidth: '900px', // ความกว้างสูงสุด
                            height: '570px', // ใช้ auto เพื่อให้ยืดหยุ่นตามเนื้อหา
                            maxHeight: '650px', // จำกัดความสูง
                            border: '1px solid #EDEDED', // ขอบของกรอบ
                            borderRadius: '10px', // มุมมนของกรอบ
                            display: 'flex', // จัดเรียงเนื้อหาภายในกรอบ
                            flexDirection: 'column', // เปลี่ยนเป็นแนวตั้ง
                            alignItems: 'flex-start', // จัดให้ชิดซ้าย
                            padding: '20px', // เพิ่ม padding
                            overflow: 'auto', // เพิ่ม scroll หากเนื้อหามากเกินไป
                          }}
                        >
                          <h2 style={{ margin: 0, color: '#35465B', fontSize: '20px' }}>Order Detail</h2>
                          <Form initialValues={detailorder[0]} layout="vertical" style={{ padding: '10px', width: '100%' }}>
      {fieldPairs.map((pair, index) => (
        <Row gutter={16} style={{ marginTop: '10px', width: '100%' }} key={index}>
          {pair.map((field) => (
            <Col span={12} key={field.name}>
              <Form.Item label={field.label} name={field.name}>
                <Input placeholder={field.label} disabled />
              </Form.Item>
            </Col>
          ))}
        </Row>
      ))}
    </Form>
                        </div>
                      </Col>
                      <Col span={12} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                        <div
                          style={{
                            width: '90%', // ปรับเป็นเปอร์เซ็นต์เพื่อให้ยืดหยุ่น
                            maxWidth: '900px', // ความกว้างสูงสุด
                            height: '570px', // ใช้ auto เพื่อให้ยืดหยุ่นตามเนื้อหา
                            maxHeight: '650px', // จำกัดความสูง
                            border: '1px solid #EDEDED', // ขอบของกรอบ
                            borderRadius: '10px', // มุมมนของกรอบ
                            display: 'flex', // จัดเรียงเนื้อหาภายในกรอบ
                            flexDirection: 'column', // เปลี่ยนเป็นแนวตั้ง
                            alignItems: 'flex-start', // จัดให้ชิดซ้าย
                            padding: '20px', // เพิ่ม padding
                            overflow: 'auto', // เพิ่ม scroll หากเนื้อหามากเกินไป
                          }}
                        >
                          <h2 style={{ margin: 0, color: '#35465B', fontSize: '20px' }}>คืนครบ</h2>
                          <Form style={{ padding: '5px', width: '100%' }}>
                            <Table
                              components={{
                                header: {
                                  cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                    <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                                  ),
                                },
                              }}
                              className="custom-table"
                              columns={columns}
                              dataSource={datacolumn}
                            />
                          </Form>
                          <h2 style={{ margin: 0, color: '#35465B', fontSize: '20px' }}>คืนไม่ครบ</h2>
                          <Form style={{ padding: '5px', width: '100%' }}>
                            <Table
                              components={{
                                header: {
                                  cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                    <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                                  ),
                                },
                              }}
                              className="custom-table"
                              columns={columns}
                              dataSource={datacolumn1}
                            />
                          </Form>
                        </div>
                      </Col>
                    </Row>
                  </TabPane>
                  <TabPane tab="Return Images" key="2">
                    {/* เนื้อหาสำหรับ Return Images */}
                    <Row justify="space-between" align="middle">
                      <Col span={12} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                        <div
                          style={{
                            width: '90%', // ปรับเป็นเปอร์เซ็นต์เพื่อให้ยืดหยุ่น
                            maxWidth: '900px', // ความกว้างสูงสุด
                            height: '570px', // ใช้ auto เพื่อให้ยืดหยุ่นตามเนื้อหา
                            maxHeight: '650px', // จำกัดความสูง
                            border: '1px solid #EDEDED', // ขอบของกรอบ
                            borderRadius: '10px', // มุมมนของกรอบ
                            display: 'flex', // จัดเรียงเนื้อหาภายในกรอบ
                            flexDirection: 'column', // เปลี่ยนเป็นแนวตั้ง
                            alignItems: 'flex-start', // จัดให้ชิดซ้าย
                            padding: '20px', // เพิ่ม padding
                            overflow: 'auto', // เพิ่ม scroll หากเนื้อหามากเกินไป
                          }}
                        >
                         <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'center' }}>
  {                   datacolumn?.map((item: any) => (
    <Card
      key={item.SKU} // เพิ่ม key เพื่อให้ React รู้จักแต่ละ Card
      style={{ margin: '10px', textAlign: 'center', width: '200px' }} // กำหนดความกว้างให้คงที่
      cover={
        <Image
          width={200}
          src="https://zos.alipayobjects.com/rmsportal/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png"
          preview // Enable image preview
        />
      }
    >
      <span style={{ fontSize: '20px', color: '#35465B' }}>{item.SKU}</span>
    </Card>
  ))}
</div>

                        </div>

                      </Col>
                      <Col span={12} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                        <div
                          style={{
                            width: '90%', // ปรับเป็นเปอร์เซ็นต์เพื่อให้ยืดหยุ่น
                            maxWidth: '900px', // ความกว้างสูงสุด
                            height: '570px', // ใช้ auto เพื่อให้ยืดหยุ่นตามเนื้อหา
                            maxHeight: '650px', // จำกัดความสูง
                            border: '1px solid #EDEDED', // ขอบของกรอบ
                            borderRadius: '10px', // มุมมนของกรอบ
                            display: 'flex', // จัดเรียงเนื้อหาภายในกรอบ
                            flexDirection: 'column', // เปลี่ยนเป็นแนวตั้ง
                            alignItems: 'flex-start', // จัดให้ชิดซ้าย
                            padding: '20px', // เพิ่ม padding
                            overflow: 'auto', // เพิ่ม scroll หากเนื้อหามากเกินไป
                          }}
                        >

<div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'center' }}>
  {datacolumn?.map((item: any) => (
    <Card
      key={item.SKU} // เพิ่ม key เพื่อให้ React รู้จักแต่ละ Card
      style={{ margin: '10px', textAlign: 'center', width: '200px' }} // กำหนดความกว้างให้คงที่
      cover={
        <Image
          width={200}
          src="https://zos.alipayobjects.com/rmsportal/jkjgkEfvpUPVyRjUImniVslZfWPnJuuZ.png"
          preview // Enable image preview
        />
      }
    >
      <span style={{ fontSize: '20px', color: '#35465B' }}>{item.SKU}</span>
    </Card>
  ))}
</div>

                        </div>
                      </Col>
                    </Row>



                  </TabPane>
                </Tabs>
              )}
            </div>
          </Popup>

        </Layout.Content>
      </Layout>
    </ConfigProvider>
  );
};

export default OtherReturn;
