import { Link } from "react-router-dom";
import { Button, ConfigProvider, Form, Layout, Row, Table, Tabs, Tooltip, Modal, Input, Col, Select, InputNumber, Popconfirm, notification } from "antd";
import { DatePicker } from "antd";
import { DeleteOutlined, FormOutlined, PlusCircleOutlined } from '@ant-design/icons'; // นำเข้า FormOutlined
import React, { useState, useEffect, useRef } from "react";
import dayjs, { Dayjs } from "dayjs";
import utc from "dayjs/plugin/utc";
import isSameOrAfter from "dayjs/plugin/isSameOrAfter";
import isSameOrBefore from "dayjs/plugin/isSameOrBefore";
import isBetween from "dayjs/plugin/isBetween";
import '../Return.css';
import Webcam from "react-webcam";
import api from "../../utils/axios/axiosInstance"; 
import { useSelector } from 'react-redux';
import { RootState } from "../../redux/types";
import { Order, OrderDetail, OrderLine, SKUData, SelectedRecord} from '../../types/types';

dayjs.extend(utc);
dayjs.extend(isSameOrAfter);
dayjs.extend(isSameOrBefore);
dayjs.extend(isBetween);

const ConfirmReturnTrade = () => {
    const [form] = Form.useForm();
    const [dates, setDates] = useState<[Dayjs, Dayjs] | null>(null);
    const { RangePicker } = DatePicker;
    const [activeTabKey, setActiveTabKey] = useState<string>("1");
    const [filteredData, setFilteredData] = useState<Order[]>([]);
    const [isNewModalVisible, setIsNewModalVisible] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [selectedRecord, setSelectedRecord] = useState<Order | null>(null);

    const [codeR, setCodeR] = useState<string | undefined>(undefined);
    const [nameR, setNameR] = useState<string | undefined>(undefined);
    const [qty, setQty] = useState<number | null>(null);  
    const [price, setPrice] = useState<number | null>(null); 

    const [newEntries, setNewEntries] = useState([]);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(5);

    const [currentPageConfirm, setCurrentPageConfirm] = useState(1);
    const [pageSizeConfirm, setPageSizeConfirm] = useState(5);

    const [editingSKU, setEditingSKU] = useState<string | null>(null);
    const [editedValues, setEditedValues] = useState<{ QTY?: number; Price?: number }>({});
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedOrderNo, setSelectedOrderNo] = useState<string | null>(null);

      // ดึงข้อมูลผู้ใช้ที่เข้าสู่ระบบ
      const auth = useSelector((state: RootState) => state.auth);
      const userID = auth?.user?.userID;

    const handleEdit = async (order: string, activeTabKey: string) => {
        try {
            const response = await api.get(`/api/return-order/get-lines/${order}`);
            const orderLines: OrderLine[] = response.data.data || [];
            const initialOrder: Order = {
                Order: order,
                SO_INV: "", // หรือค่าเริ่มต้นที่เหมาะสม
                Customer: "", // หรือค่าเริ่มต้นที่เหมาะสม
                SR: "", // หรือค่าเริ่มต้นที่เหมาะสม
                ReturnTracking: "", // หรือค่าเริ่มต้นที่เหมาะสม
                Channel: "", // หรือค่าเริ่มต้นที่เหมาะสม
                Date_Create: "", // หรือค่าเริ่มต้นที่เหมาะสม
                Warehouse: "", // หรือค่าเริ่มต้นที่เหมาะสม
                Transport: "",
                data: [],
                // data:, // หรือค่าเริ่มต้นที่เหมาะสม
                // // ... properties อื่นๆ ที่ Order มี (เติมให้ครบ)
            };

            const updatedRecord = {
                ...initialOrder, 
                Order: order,
                data: orderLines.map((line: OrderLine) => ({
                    OrderNo: order,
                    SKU: line.sku,
                    Name: line.itemName,
                    QTY: line.qty,
                    Price: line.price,
                    Action: '',
                    Type: line.Type,
                })),
            };
            setSelectedRecord(updatedRecord); // เก็บข้อมูล record ที่เลือกพร้อมกับ OrderLine
            setIsModalVisible(true); 
        } catch (error) {
            console.error('Failed to fetch order lines:', error);
            notification.error({
            message: 'Error',
            description: 'Failed to fetch order lines.',
            });
        }
    };
    
    const handleOk = () => { // ปรับให้เมื่อกดแล้วเปลี่ยนสถานะจาก waiting => confirm สถานะ StatusConfirm เปลี่ยน => แสดงชื่อคนกดที่ ConfirmBy
      handleUpdate();
      setIsModalVisible(false); // ปิด Modal
    };

    const handleCancel = () => {
        setIsModalVisible(false);
        setSelectedRecord(null);
    };

    const handleEditLine = (orderNo: string, sku: string, currentQTY: number, currentPrice: string) => {
        setSelectedOrderNo(orderNo);
        setEditingSKU(sku);
        setEditedValues({ 
            QTY: currentQTY, 
            Price: Number(currentPrice)  // ✅ แปลงเป็น number
        });
        setIsModalOpen(true);
    };

    // 📌 ฟังก์ชันอัปเดตค่าที่แก้ไข
    const handleUpdateLine = async () => {
        if (!selectedOrderNo || !editingSKU) return;
    
        try {
            const token = localStorage.getItem('access_token')
            const response = await api.patch(`/api/return-order/update-line/${selectedOrderNo}/${editingSKU}`, {  // Use PATCH
                ActualQTY: editedValues.QTY,
                Price: editedValues.Price,
                UpdateBy: userID,
            }, {  headers: {
                Authorization: `Bearer ${token}`,
              },});
    
            if (response.status === 200) {
                notification.success({
                  message: "ส่งข้อมูลสำเร็จ",
                  description: "ข้อมูลของคุณถูกส่งเรียบร้อยแล้ว!",
                });
                setIsModalOpen(false);
                setEditingSKU(null);
                // Reload data or update state as necessary
            } else {
                notification.error({
                    message: 'Error',
                    description: 'An error occurred while updating.',
                });
            }
        } catch (error) {
            notification.error({
                message: 'Error',
                description: 'Cannot connect to the server.',
            });
        }
    };
    

    const fetchData = async (statusCheckID: number) => {
        try {
            const endpoint = statusCheckID === 1 
                ? '/api/trade-return/get-waiting' 
                : '/api/trade-return/get-confirm';
            const response = await api.get(endpoint);
            const data = response.data.data.map((item: OrderDetail) => ({
                Order: item.orderNo,
                SO_INV: item.soNo,
                Customer: item.customerId,
                SR: item.srNo,
                ReturnTracking: item.trackingNo,
                // Transport: item.logistic,
                Channel: item.channelName,
                Date_Create: dayjs(item.createDate).utc().format('YYYY-MM-DD'),
                Warehouse: item.warehouseName,
                data: [], // Assuming you have a way to get SKUData
            }));

            setFilteredData(data);
        } catch (error) {
            console.error('Failed to fetch data:', error);
            notification.error({
                message: 'Error',
                description: 'Failed to fetch data.',
            });
        }
    };

    useEffect(() => {
        fetchData(activeTabKey === "1" ? 1 : 2);
    }, [activeTabKey]);
    
    useEffect(() => {
        handleSearch(activeTabKey === "1" ? 1 : 2);
    }, []);

    const handleSearch = async (statusCheckID: number) =>  {
        // if (!dates || !dates[0] || !dates[1]) {
        //     if (isManualSearch) { // แจ้งเตือนเฉพาะตอนกดปุ่ม
        //       console.log("No date selected, skipping search.");
        //       notification.warning({
        //         message: 'Warning',
        //         description: 'Please select ragnge date before searching.',
        //       });
        //     }
        //     return;
        //   }

        if (dates && dates[0] && dates[1]) {
            const startDate = dates[0].format('YYYY-MM-DD');
            const endDate = dates[1].format('YYYY-MM-DD');

            try {
                const endpoint = statusCheckID === 1 
                    ? '/api/trade-return/search-waiting' 
                    : '/api/trade-return/search-confirm';

                const response = await api.get(endpoint, {
                    params: {
                    startDate,
                    endDate,
                    },
                });

                const data = response.data.data || []; // ตรวจสอบว่ามีข้อมูลหรือไม่
                const filtered = data.filter((item: OrderDetail) => {
                    const itemDate = dayjs(item.createDate).utc().format('YYYY-MM-DD');
                    return dayjs(itemDate).isSameOrAfter(startDate) && dayjs(itemDate).isSameOrBefore(endDate);
                }).map((item: OrderDetail) => ({
                    Order: item.orderNo,
                    SO_INV: item.soNo,
                    Customer: item.customerId,
                    SR: item.srNo,
                    ReturnTracking: item.trackingNo,
                    // Transport: item.logistic,
                    Channel: item.channelName,
                    Date_Create: dayjs(item.createDate).utc().format('YYYY-MM-DD'),
                    Warehouse: item.warehouseName,
                    data: [], // Assuming you have a way to get SKUData
                }));

            //   if (filtered.length === 0) {
            //     notification.warning({
            //       message: 'Data not found',
            //       description: 'Please select new date range again!',
            //     });
            //     // setDates(null); 
            //     return;
            //   } 
        
                setFilteredData(filtered);
            } catch (error) {
            console.error('Failed to fetch data:', error);
            notification.error({
                message: 'Error',
                description: 'Failed to fetch data.',
            });
            }
        }
    };

    const handlePageChange = (page: number, pageSize: number) => {
        setCurrentPage(page);
        setPageSize(pageSize); // ถ้าผู้ใช้เลือกจำนวนรายการต่อหน้าใหม่
    };

    const handlePageChangeConfirm  = (page: number, pageSize: number) => {
        setCurrentPageConfirm(page);
        setPageSizeConfirm(pageSize); // ถ้าผู้ใช้เลือกจำนวนรายการต่อหน้าใหม่
    };
    
    const handleAdd = () => {
      if (selectedRecord) {
        const newData: SKUData = {
          OrderNo: "",
          SKU: codeR || '',
          Name: nameR || '',
          QTY: qty || 0,
          Price: price ? price.toFixed(2) : '0.00',
          Action: 'delete',
          Type: 'addon',  // กำหนด Type เป็น addon สำหรับข้อมูลที่เพิ่มใหม่
        };
    
        const updatedData = [...selectedRecord.data, newData];
        setSelectedRecord({ ...selectedRecord, data: updatedData });
    
        setCodeR(undefined);
        setNameR(undefined);
        setQty(null);
        setPrice(null);
      }
    };
  
    const handleDelete = (skuToDelete: string) => {
      if (selectedRecord) {
        // ลบเฉพาะรายการที่มี Type เป็น 'addon' และ SKU ตรงกับ skuToDelete
        const updatedData = selectedRecord.data.filter(
          item => !(item.SKU === skuToDelete && item.Type === 'addon')
        );
    
        // อัปเดต selectedRecord ด้วยข้อมูลที่ผ่านการกรองแล้ว
        setSelectedRecord({ ...selectedRecord, data: updatedData });
      }
    };
    
    const onTabChange = (key: string) => {
        setActiveTabKey(key);
        handleSearch(key === "1" ? 1 : 2); // เรียก handleSearch พร้อมกับ StatusCheckID ที่ถูกต้อง
    };

    const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
        if (!dates || !dates[0] || !dates[1]) {
            setDates(null); // รีเซ็ตค่า Select date เมื่อไม่มีการเลือกช่วงวันที่
            return;
        }

        if (dates) {
            setDates(dates as [Dayjs, Dayjs]);
            // handleSearch(activeTabKey === "1" ? 1 : 2); // เรียก handleSearch พร้อมกับ StatusCheckID ที่ถูกต้อง
        }
    };

    const handleUpdate = () => {
      // บันทึก newEntries ลงใน selectedRecord
      if (selectedRecord) {
          setSelectedRecord({
              ...selectedRecord,
              data: [...selectedRecord.data, ...newEntries]
          });
          setNewEntries([]); // รีเซ็ต newEntries หลังจากบันทึก
          setActiveTabKey('2'); // เปลี่ยนแท็บไปที่ "Confirm"
      }
    };

    const columns = [
        { title: "Order", dataIndex: "Order", id:"Order", key: "Order",     render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span>  },
        { title: "SO/INV", dataIndex: "SO_INV", id:"SO_INV", key: "SO_INV", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "Customer", dataIndex: "Customer", id:"Customer", key: "Customer" ,render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "SR", dataIndex: "SR", id:"SR", key: "SR",render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span>  },
        { title: "Return Tracking", id:"ReturnTracking", dataIndex: "ReturnTracking", key: "ReturnTracking",render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span>  },
        { title: "Channel",  id:"Channel",  dataIndex: "Channel", key: "Channel",render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span>  },
        { title: "Date Create",  id:"Date_Create",  dataIndex: "Date_Create", key: "Date_Create" ,render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "Warehouse",  id:"Warehouse",  dataIndex: "Warehouse", key: "Warehouse",render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span>  },
        {
            title: "Action",
            id:"Action", 
            dataIndex: "Action",
            key: "Action",
            render: (_: any, record: Order) => (
                <Tooltip title="Edit">
                    <Button 
                        type="link" 
                        icon={<FormOutlined />} 
                        onClick={() => handleEdit(record.Order, activeTabKey)}
                        style={{ color: 'gray', textAlign: "center" }}
                    />
                </Tooltip>
            ),
        },
    ];

    const columnsconfirm = [
      { title: "Order", dataIndex: "Order", key: "Order", id:"Order",  },
      { title: "SO/INV", dataIndex: "SO_INV", key: "SO_INV" , id:"SO_INV", },
      { title: "Customer", dataIndex: "Customer", key: "Customer" , id:"Customer", },
      { title: "SR", dataIndex: "SR", key: "SR", id:"SR",  },
      { title: "Return Tracking", dataIndex: "ReturnTracking", key: "ReturnTracking" , id:"ReturnTracking", },
      { title: "Channel", dataIndex: "Channel", key: "Channel", id:"Channel", width: 80, },
      { title: "Date Create", dataIndex: "Date_Create", key: "Date_Create", id:"Date_Create",  },
      { title: "Warehouse", dataIndex: "Warehouse", key: "Warehouse", id:"Warehouse", width: 80, },
      {
          title: "Action",
          id:"Action",
          dataIndex: "Action",
          key: "Action",
          width: 80,
          render: (_: any, record: Order) => (
              <Tooltip title="Edit">
                  <Button 
                      type="link" 
                      icon={<FormOutlined />} 
                      onClick={() => handleEdit(record.Order, activeTabKey)}
                      style={{ color: 'gray', textAlign: "center" }}
                  />
              </Tooltip>
          ),
      },
    ];

return (
    <ConfigProvider>
        <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
            Confirm Return Trade
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
            <Tabs
            id="card"
                onChange={onTabChange}
                type="card"
                items={[
                    { label: "Waiting", key: "1" },
                    { label: "Confirm", key: "2" },
                ]}
            />

            {activeTabKey === "1" && (
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
                            onClick={() => handleSearch(1)}
                        >
                            Search
                        </Button>
                    </Col>
                </Row>
                <div>
                    <Table 
                        id="Table1"
                        components={{
                            header: {
                                cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                  <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B', padding: "12px", textAlign: 'center' }} />
                                ),
                              },
                              body: {
                                  cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                    <td {...props} style={{ padding: "12px", textAlign: 'center'}} />
                                  ),
                              }
                        }}
                        pagination={false} 
                        style={{
                            width: "100%",
                            tableLayout: "auto",
                            border: "1px solid #ddd",
                            borderRadius: "8px",
                          }}
                        scroll={{ x: 'max-content' }} // ป้องกันตารางล้น จะเลื่อนไปทางซ้าย-ขวาแทนการบีบข้อความ
                        // dataSource={filteredData}
                        dataSource={filteredData.slice((currentPage - 1) * pageSize, currentPage * pageSize)}
                        columns={columns} 
                        rowKey={(record: any, index) => (index as number).toString()}
                    />
                     <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginTop: 20 }}>
                        <span style={{ fontSize: "14px", fontWeight: "bold", color: "#555" }}>
                        ทั้งหมด <span style={{ color: "#007bff" }}>{filteredData.length}</span> รายการ
                        </span>
                    </div>

                    <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginTop: 20, gap: 10 }}>
                        <Button
                        onClick={() => handlePageChange(1, pageSize)}
                        disabled={currentPage === 1}
                        >
                        {"<<"}
                        </Button>
                        <Button
                        onClick={() => handlePageChange(currentPage - 1, pageSize)}
                        disabled={currentPage === 1}
                        >
                        {"<"}
                        </Button>
                        <span style={{ fontSize: "14px", fontWeight: "bold" }}>
                        [ {currentPage} to {Math.ceil(filteredData.length / pageSize)} ]
                        </span>
                        <Button
                        onClick={() => handlePageChange(currentPage + 1, pageSize)}
                        disabled={currentPage === Math.ceil(filteredData.length / pageSize)}
                        >
                        {">"}
                        </Button>
                        <Button
                        onClick={() => handlePageChange(Math.ceil(filteredData.length / pageSize), pageSize)}
                        disabled={currentPage === Math.ceil(filteredData.length / pageSize)}
                        >
                        {">>"}
                        </Button>

                        <select
                            value={pageSize}
                            onChange={(e) => handlePageChange(1, Number(e.target.value))}
                            style={{
                                fontSize: "14px",
                                fontWeight: "bold",
                                padding: "4px 10px",
                                border: "1px solid #ddd",
                                borderRadius: "6px",
                                cursor: "pointer",
                            }}
                        >
                            <option value="5">5 รายการ</option>
                            <option value="10">10 รายการ</option>
                            <option value="20">20 รายการ</option>
                        </select>
                    </div>
                </div>
            </>
            )}
            {activeTabKey === '2' && (
                <>
                    <Row gutter={8} align="middle" justify="center" style={{ marginTop: "20px" }}>
                        <Col>
                            <Form.Item
                            id="Select date2"
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
                            id="Search2"
                                type="primary"
                                style={{ height: "40px", width: "100px", background: "#32ADE6" }}
                                onClick={() => handleSearch(2)}
                            >
                                Search
                            </Button>
                        </Col>
                    </Row>
                    <div>
                        <Table 
                            id="Table2"
                            components={{
                                header: {
                                  cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                    <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B', padding: "12px", textAlign: 'center' }} />
                                  ),
                                },
                                body: {
                                    cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                      <td {...props} style={{ padding: "12px", textAlign: 'center'}} />
                                    ),
                                }
                            }}
                            pagination={false} 
                            style={{
                                width: "100%",
                                tableLayout: "auto",
                                border: "1px solid #ddd",
                                borderRadius: "8px",
                              }}
                            scroll={{ x: 'max-content' }} 
                            dataSource={filteredData.slice((currentPageConfirm  - 1) * pageSizeConfirm, currentPageConfirm  * pageSizeConfirm)}
                            columns={columnsconfirm} 
                            rowKey={(record: any, index) => (index as number).toString()}
                        />
                    <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginTop: 20 }}>
                        <span style={{ fontSize: "14px", fontWeight: "bold", color: "#555" }}>
                        ทั้งหมด <span style={{ color: "#007bff" }}>{filteredData.length}</span> รายการ
                        </span>
                    </div>

                    <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginTop: 20, gap: 10 }}>
                        <Button
                            onClick={() => handlePageChangeConfirm(1, pageSizeConfirm)}
                            disabled={currentPageConfirm  === 1}
                        >
                        {"<<"}
                        </Button>
                        <Button
                            onClick={() => handlePageChangeConfirm(currentPageConfirm  - 1, pageSizeConfirm)}
                            disabled={currentPageConfirm  === 1}
                        >
                        {"<"}
                        </Button>
                        <span style={{ fontSize: "14px", fontWeight: "bold" }}>
                        [ {currentPageConfirm } to {Math.ceil(filteredData.length / pageSizeConfirm)} ]
                        </span>
                        <Button
                            onClick={() => handlePageChangeConfirm(currentPageConfirm  + 1, pageSizeConfirm)}
                            disabled={currentPageConfirm  === Math.ceil(filteredData.length / pageSizeConfirm)}
                        >
                        {">"}
                        </Button>
                        <Button
                            onClick={() => handlePageChangeConfirm(Math.ceil(filteredData.length / pageSizeConfirm), pageSizeConfirm)}
                            disabled={currentPageConfirm  === Math.ceil(filteredData.length / pageSizeConfirm)}
                        >
                        {">>"}
                        </Button>

                        <select
                            value={pageSize}
                            onChange={(e) => handlePageChangeConfirm(1, Number(e.target.value))}
                            style={{
                                fontSize: "14px",
                                fontWeight: "bold",
                                padding: "4px 10px",
                                border: "1px solid #ddd",
                                borderRadius: "6px",
                                cursor: "pointer",
                            }}
                        >
                            <option value="5">5 รายการ</option>
                            <option value="10">10 รายการ</option>
                            <option value="20">20 รายการ</option>
                        </select>
                    </div>
                    </div>
                </>
            )}
            </Layout.Content>
        </Layout>

        {activeTabKey=='1' && (
            <Modal
                closable={false}
                width={800}
                title="Edit Order"
                visible={isModalVisible}
                onOk={handleOk}
                footer={
                    <div style={{ display: 'flex', justifyContent: 'center' }}>

                        <Button id="Confirm" onClick={handleOk} style={{ marginLeft: 8, backgroundColor: '#14C11B', color: '#FFF' }}>
                            Confirm
                        </Button>
                        <Button id="Cancel" onClick={handleCancel} style={{ marginLeft: 8, background: '#D9D9D9', color: '#909090' }}>
                            Cancel
                        </Button>
                    </div>
                }
            >
            {selectedRecord && (
                <>
                    <Form layout="vertical" style={{ marginTop: 20 }}>
                        <Row gutter={16}>
                            <Col span={8}>
                            <Form.Item id="Order1" label={<span style={{ color: '#657589' }}>Order</span>}>
                                    <Input style={{ height: 40 }} value={selectedRecord.Order} readOnly disabled />
                                </Form.Item>
                            </Col>
                            <Col span={8}>
                            <Form.Item id="So/inv1" label={<span style={{ color: '#657589' }}>SO/INV</span>}>
                                    <Input style={{ height: 40 }} value={selectedRecord.SO_INV} disabled />
                                </Form.Item>
                            </Col>
                            <Col span={8}>
                            <Form.Item id="SR1" label={<span style={{ color: '#657589' }}>SR</span>}>
                                    <Input style={{ height: 40 }} value={selectedRecord.SR} disabled />
                                </Form.Item>
                            </Col>
                        </Row>
                    </Form>

                    {/* Table to display product data */}
                    <Table
                        id="Table3"
                        components={{
                        header: {
                            cell: (props: React.HTMLAttributes<HTMLElement>) => (
                            <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                            ),
                        },
                        }}
                        columns={[
                            {  title: 'SKU', dataIndex: 'SKU', render: (text) => <span style={{ color: '#35465B' }}>{text}</span>  },
                            { title: 'Name', dataIndex: 'Name', render: (text) => <span style={{ color: '#35465B' }}>{text}</span>  },
                            { title: 'QTY', dataIndex: 'QTY', render: (text) => <span style={{ color: '#35465B' }}>{text}</span>  },
                            { title: 'Price', dataIndex: 'Price', render: (text) => <span style={{ color: '#35465B' }}>{text}</span>  },
                            {
                              title: 'Action',
                              dataIndex: 'Action',
                              render: (_, record) => 
                                <Tooltip title="Edit">
                                    <Button
                                        type="link"
                                        icon={<FormOutlined style={{ color: 'blue' }} />}
                                        onClick={() => handleEditLine(record.OrderNo, record.SKU, record.QTY, record.Price)} // ✅ ใช้ record.OrderNo
                                    />
                                </Tooltip>
                            },
                         ]}
                        dataSource={selectedRecord.data} // Use updated data with new entries
                        rowKey="SKU"
                        pagination={false}
                    />
                    // 📌 Modal สำหรับแก้ไขข้อมูล
                    <Modal
                        title="แก้ไขข้อมูลสินค้า"
                        open={isModalOpen}
                        onCancel={() => setIsModalOpen(false)}
                        onOk={handleUpdateLine}
                        okText="บันทึก"
                        cancelText="ยกเลิก"
                    >
                        <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
                            <div>
                                <label>QTY:</label>
                                <Input
                                    type="number"
                                    value={editedValues.QTY}
                                    onChange={(e) => setEditedValues({ ...editedValues, QTY: Number(e.target.value) })}
                                />
                            </div>
                            <div>
                                <label>Price:</label>
                                <Input
                                    type="number"
                                    value={editedValues.Price}
                                    onChange={(e) => setEditedValues({ ...editedValues, Price: Number(e.target.value) })}
                                />
                            </div>
                        </div>
                    </Modal>
                </>
            )}
            </Modal>
        )}
        {activeTabKey=='2' && (
            <Modal
                width={800}
                title="Confrim"
                visible={isModalVisible}
                onOk={handleOk}
                onCancel={handleCancel}
                footer={
                    <div style={{ display: 'flex', justifyContent: 'center' }}>
                        
                    
                    </div>
                }
            >
            {selectedRecord && (
                <>
                    <Form layout="vertical" style={{ marginTop: 20 }}>
                        <Row gutter={16} align="middle" justify="center" style={{ marginTop: "20px" }}>
                            <Col span={12}>
                                <Form.Item id="Order2" label={<span style={{ color: '#657589' }}>Order</span>}>
                                    <Input style={{ height: 40 }} value={selectedRecord.Order} readOnly disabled />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item id="SR2" label={<span style={{ color: '#657589' }}>SR</span>}>
                                    <Input style={{ height: 40 }} value={selectedRecord.SR} disabled />
                                </Form.Item>
                            </Col>
                        
                        </Row>
                    </Form>

                    {/* Table to display product data */}
                    <Table
                        id="Table4"
                        components={{
                          header: {
                            cell: (props: React.HTMLAttributes<HTMLElement>) => (
                              <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                            ),
                          },
                        }}
                        columns={[
                            { title: 'SKU', dataIndex: 'SKU' },
                            { title: 'Name', dataIndex: 'Name' },
                            { title: 'QTY', dataIndex: 'QTY' },
                            { title: 'Price', dataIndex: 'Price' },
                            {
                                title: 'Action',
                                dataIndex: 'Action',
                                render: (_, record) => 
                                        <Popconfirm
                                            title="Are you sure to delete this item?"
                                            onConfirm={() => handleDelete(record.SKU)}
                                            okText="Yes"
                                            cancelText="No"
                                        >
                                            <Button 
                                                type="link"
                                                icon={<DeleteOutlined style={{ color: 'red' }} />} 
                                            />
                                        </Popconfirm>
                             
                            },
                         
                        ]}
                        dataSource={selectedRecord.data} // Use updated data with new entries
                        rowKey="SKU"
                        pagination={false}
                    />
                    </>
            )}
            </Modal>
        )}
        </ConfigProvider>
    );
};

export default ConfirmReturnTrade ;
