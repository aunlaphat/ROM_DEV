import { Link } from "react-router-dom";
import { Button, ConfigProvider, Form, Layout, Row, Table, Tabs, Tooltip, Modal, Input, Col, Select, InputNumber, Popconfirm } from "antd";
import { DatePicker } from "antd";
import { DeleteOutlined, FormOutlined, PlusCircleOutlined } from '@ant-design/icons'; // นำเข้า FormOutlined
import React, { useState } from "react";
import dayjs, { Dayjs } from "dayjs";
import isSameOrAfter from "dayjs/plugin/isSameOrAfter";
import isSameOrBefore from "dayjs/plugin/isSameOrBefore";
import isBetween from "dayjs/plugin/isBetween";
import '../Return.css';
import * as XLSX from "xlsx";

dayjs.extend(isSameOrAfter);
dayjs.extend(isSameOrBefore);
dayjs.extend(isBetween);

interface Order {
    SO: string; // SO
    Tracking_Order: string; // Tracking Order
    SKU: string; // SKU
    SKU_Name: string; // SKU Name
    QTY: number; // Quantity
    Amount: string; // Amount
    Warehouse: string; // Warehouse
    Location: string; // Location
    Ship_date: string; // Ship Date
    Site:string;
    Channel:string;
  }



  

const Report = () => {
    
    

    const columnsdata: Order[] = [
        {
          SO: "SO123456",
          Tracking_Order: "RT123456",
          SKU: "G090108-EF05",
          SKU_Name: "Bewell Official Store",
          QTY: 20,
          Amount: '599.00',
          Site: "DCOM",
          Warehouse: "RBN",
          Location: "Location A",
          Ship_date: "2024-09-01",
          Channel:"ECOM",

        },
        {
            SO: "SO123457",
            Tracking_Order: "RT123457",
            SKU: "G090108-EF04",
            SKU_Name: "Bewell Shop",
            QTY: 50,
            Amount: '599.00',
            Site: "DCOM",
            Warehouse: "RBN",
            Location: "Location B",
            Ship_date: "2024-09-15",
            Channel:"ECOM",
          },
          {
            SO: "SO123458",
            Tracking_Order: "RT123458",
            SKU: "G090108-EF05",
            SKU_Name: "Bewell Official Store",
            QTY: 20,
            Amount: '599.00',
            Site: "DCOM",
            Warehouse: "RBN",
            Location: "Location C",
            Ship_date: "2024-09-29",
            Channel:"ECOM",
          },
];


const columns = [
    { title: "SO", dataIndex: "SO", key: "SO" },
    { title: "Tracking Order", dataIndex: "Tracking_Order", key: "Tracking_Order" },
    { title: "SKU", dataIndex: "SKU", key: "SKU" },
    { title: "SKU Name", dataIndex: "SKU_Name", key: "SKU_Name" },
    { title: "QTY", dataIndex: "QTY", key: "QTY" },
    { title: "Amount", dataIndex: "Amount", key: "Amount" },
    { title: "Site", dataIndex: "Site", key: "Site" },
    { title: "Warehouse", dataIndex: "Warehouse", key: "Warehouse" },
    { title: "Location", dataIndex: "Location", key: "Location" },
    { title: "Ship Date", dataIndex: "Ship_date", key: "Ship_date" },
    { title: "Channel", dataIndex: "Channel", key: "Channel" },
   
  ];
    
  
   
    
    const [dates, setDates] = useState<[Dayjs, Dayjs] | null>(null);
    const { RangePicker } = DatePicker;
    const [filteredData, setFilteredData] = useState<Order[]>(columnsdata);
   
    

      

    const handleExportExcel = () => {
        // สร้างเวิร์กชีทจากข้อมูลที่กรองแล้ว
        const worksheet = XLSX.utils.json_to_sheet(filteredData);
        const workbook = XLSX.utils.book_new();
        XLSX.utils.book_append_sheet(workbook, worksheet, "Report Data");
    
        // เขียนไฟล์ Excel และบันทึก
        XLSX.writeFile(workbook, "report_data.xlsx");
    };
  


    const handleSearch = () => {
        if (dates && dates[0] && dates[1]) {
            const startDate = dates[0].startOf("day");
            const endDate = dates[1].endOf("day");

            const filtered = columnsdata.filter((item) => {
                const itemDate = dayjs(item.Ship_date);

                return itemDate.isBetween(startDate, endDate, null, "[]");
            });

            setFilteredData(filtered);
        }
        
    };
    
    
    

    const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
        if (dates) {
            setDates(dates as [Dayjs, Dayjs]);
        }
    };
    

    return (
        
        <ConfigProvider>
             
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
           Report
            </div>
            <Layout>
                <Layout.Content
                    style={{
                        margin: "24px",
                        padding: 36,
                        minHeight: 360,
                        background: "#fff",
                        borderRadius: "8px",
                        overflow: "auto",
                    }}
                >
                    

              
                        <>
                            <Row gutter={8} align="middle" justify="center" style={{ marginTop: "20px" }}>
                                <Col>
                                    <Form.Item
                                    id="Select date"
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
                                    id="Search"
                                        type="primary"
                                        style={{ height: "40px", width: "100px", background: "#32ADE6" }}
                                        onClick={handleSearch}
                                    >
                                        Search
                                    </Button>
                                </Col>
                            </Row>
                            <div>
                            <Col style={{ marginTop: "4px" }}>
                        <Button
                         id="Export Excel"
                            type="primary"
                            style={{ height: "40px", width: "120px", background: "#B6AEF3", marginBottom:'20px' }}
                            onClick={handleExportExcel}
                        >
                            Export Excel
                        </Button>
                    </Col>
                                <Table 
                             id="Table1"
                                 components={{
                                    header: {
                                      cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                        <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                                      ),
                                    },
                                  }}
                                 pagination={false} // Disable pagination if necessary
                                 style={{ width: '100%', tableLayout: 'fixed' }} // Ensure the table takes full width and is fixed layout
                                 scroll={{ x: 'max-content' }}
                                
                                    dataSource={filteredData}
                                    columns={columns} 
                                    rowKey="Order"
                                />
                            </div>
                        </>
                   

                    
                </Layout.Content>
            </Layout>

           
    
                       
        </ConfigProvider>
    );
};

export default Report ;
