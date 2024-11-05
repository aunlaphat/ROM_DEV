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
import Webcam from "react-webcam";

dayjs.extend(isSameOrAfter);
dayjs.extend(isSameOrBefore);
dayjs.extend(isBetween);

interface Order {
    Order: string;
    SO_INV: string;
    Customer: string;
    SR: string;
    Transport: string;
    ReturnTracking: string;
    Channel: string;
    Date_Create: string;
    Warehouse: string;
    data: SKUData[];  // ใช้ SKUData ที่มีการกำหนด Type ที่ถูกต้อง
    codeR?: string;
    nameR?: string;
}




interface SKUData {
    SKU: string;
    Name: string;
    QTY: number;
    Price: string;
    Action: string;
    Type: 'system' | 'addon';  // เพิ่ม Type เพื่อระบุว่ามาจากในระบบหรือเป็น addon
}
interface SelectedRecord {
    data: SKUData[];
}
const DraftandConfirm = () => {





    const columnsdata: Order[] = [
        {
            Order: "12345678",
            SO_INV: "SO123456",
            Customer: "TC-NMI-0007",
            SR: "SR001",
            ReturnTracking: "RT123456",
            Transport: "SPX",

            Channel: "OTHER",
            Date_Create: "2024-09-01",
            Warehouse: "RBN",
            data: [
                { SKU: 'G090108-EF05', Name: 'Bewell Official Store', QTY: 20, Price: '599.00', Action: '', Type: 'system' },
                { SKU: 'G090108-EF04', Name: 'Bewell Shop', QTY: 50, Price: '599.00', Action: '', Type: 'system' },
            ],
        },
        {
            Order: "12345677",
            SO_INV: "SO123457",
            Customer: "TC-NMI-0008",
            SR: "SR002",
            ReturnTracking: "RT123457",
            Transport: "Flash Express",

            Channel: "OTHER",
            Date_Create: "2024-09-15",
            Warehouse: "RBN",
            data: [
                { SKU: 'G090108-EF05', Name: 'Bewell Official Store', QTY: 20, Price: '599.00', Action: '', Type: 'system' },
                { SKU: 'G090108-EF04', Name: 'Bewell Shop', QTY: 50, Price: '599.00', Action: '', Type: 'system' },
            ],
        },
        {
            Order: "12345676",
            SO_INV: "SO123458",
            Customer: "TC-NMI-0009",
            SR: "SR003",
            ReturnTracking: "RT123458",
            Transport: "SPX",

            Channel: "OTHER",
            Date_Create: "2024-09-29",
            Warehouse: "RBN",
            data: [
                { SKU: 'G090108-EF05', Name: 'Bewell Official Store', QTY: 20, Price: '599.00', Action: '', Type: 'system' },
                { SKU: 'G090108-EF04', Name: 'Bewell Shop', QTY: 50, Price: '599.00', Action: '', Type: 'system' },
            ],
        },
    ];


    const columns = [
        { title: "Order", dataIndex: "Order", key: "Order", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "SO/INV", dataIndex: "SO_INV", key: "SO_INV", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "Customer", dataIndex: "Customer", key: "Customer", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "SR", dataIndex: "SR", key: "SR", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "Return Tracking", dataIndex: "ReturnTracking", key: "ReturnTracking", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "Transport", dataIndex: "Transport", key: "Transport", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },

        { title: "Channel", dataIndex: "Channel", key: "Channel", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "Date Create", dataIndex: "Date_Create", key: "Date_Create", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        { title: "Warehouse", dataIndex: "Warehouse", key: "Warehouse", render: (text: string) => <span style={{ color: '#35465B' }}>{text}</span> },
        {
            title: "Action",
            dataIndex: "Action",
            key: "Action",
            render: (_: any, record: Order) => (
                <Tooltip title="Edit">
                    <Button
                        type="link"
                        icon={<FormOutlined />}
                        onClick={() => handleEdit(record, activeTabKey)}
                        style={{ color: 'gray' }}
                    />
                </Tooltip>
            ),
        },
    ];
    const columnsconfirm = [
        { title: "Order", dataIndex: "Order", key: "Order" },
        { title: "SO/INV", dataIndex: "SO_INV", key: "SO_INV" },
        { title: "Customer", dataIndex: "Customer", key: "Customer" },
        { title: "SR", dataIndex: "SR", key: "SR" },
        { title: "Return Tracking", dataIndex: "ReturnTracking", key: "ReturnTracking" },
        { title: "Transport", dataIndex: "Transport", key: "Transport" },

        { title: "Channel", dataIndex: "Channel", key: "Channel" },
        { title: "Date Create", dataIndex: "Date_Create", key: "Date_Create" },
        { title: "Warehouse", dataIndex: "Warehouse", key: "Warehouse" },
        {
            title: "Action",
            dataIndex: "Action",
            key: "Action",
            render: (_: any, record: Order) => (
                <Tooltip title="Edit">
                    <Button
                        type="link"
                        icon={<FormOutlined />}
                        onClick={() => handleEdit(record, activeTabKey)}
                        style={{ color: 'gray' }}
                    />
                </Tooltip>
            ),
        },
    ];

    const codeROptions = [
        { value: 'R01', label: 'R01' },
        { value: 'R02', label: 'R02' },
    ];

    const codeNameOptions = [
        { value: 'ส่วนลด', label: 'ส่วนลด' },
        { value: 'ของแถม', label: 'แถม' },
    ];
    const { Option } = Select;
    const [hover, setHover] = useState(false);
    const [dates, setDates] = useState<[Dayjs, Dayjs] | null>(null);
    const { RangePicker } = DatePicker;
    const [activeTabKey, setActiveTabKey] = useState<string>("1");
    const [filteredData, setFilteredData] = useState<Order[]>(columnsdata);
    const [isNewModalVisible, setIsNewModalVisible] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [selectedRecord, setSelectedRecord] = useState<Order | null>(null);
    const [codeR, setCodeR] = useState<string | undefined>(undefined);
    const [nameR, setNameR] = useState<string | undefined>(undefined);
    const [qty, setQty] = useState<number | null>(null);  // Allow null
    const [price, setPrice] = useState<number | null>(null); // Allow null
    const [newEntries, setNewEntries] = useState([]);

    const handleEdit = (record: Order, activeTabKey: string) => {
        setSelectedRecord(record); // เก็บข้อมูล record ที่เลือก
        setIsModalVisible(true); // แสดง Modal

    };


    const handleOk = () => {
        // Logic for saving the edited record can go here
        handleUpdate();
        setIsModalVisible(false); // ปิด Modal
    };

    const handleCancel = () => {
        setIsModalVisible(false);
        setSelectedRecord(null);
    };


    const handleSearch = () => {
        if (dates && dates[0] && dates[1]) {
            const startDate = dates[0].startOf("day");
            const endDate = dates[1].endOf("day");

            const filtered = columnsdata.filter((item) => {
                const itemDate = dayjs(item.Date_Create);
                return itemDate.isBetween(startDate, endDate, null, "[]");
            });

            setFilteredData(filtered);
        }
        const handleDelete = (index: number) => {
            if (selectedRecord) {
                // สร้างข้อมูลใหม่โดยการกรองข้อมูลที่ไม่ต้องการออก
                const updatedData = selectedRecord.data.filter((_, idx) => idx !== index);

                // อัพเดต selectedRecord ด้วยข้อมูลใหม่
                setSelectedRecord({ ...selectedRecord, data: updatedData });
            }
        };
    };
    const handleAdd = () => {
        if (selectedRecord) {
            const newData: SKUData = {
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
    };

    const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
        if (dates) {
            setDates(dates as [Dayjs, Dayjs]);
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

            // เปลี่ยนแท็บไปที่ "Confirm draft"
            setActiveTabKey('2');
        }
    };

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
                        onChange={onTabChange}
                        type="card"
                        items={[
                            { label: "Draft", key: "1" },
                            { label: "Confirm Draft", key: "2" },
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
                            <div>
                                <Table
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
                    )}

                    {activeTabKey === '2' && (
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
                            <div>
                                <Table
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
                                    columns={columnsconfirm}
                                    rowKey="Order"
                                />
                            </div>
                        </>
                    )}
                </Layout.Content>
            </Layout>

            {activeTabKey == '1' && (
                <Modal
                    closable={false}
                    width={800}
                    title="Edit Order"
                    visible={isModalVisible}
                    onOk={handleOk}

                    footer={
                        <div style={{ display: 'flex', justifyContent: 'center' }}>

                            <Button onClick={handleOk} style={{ marginLeft: 8, backgroundColor: '#14C11B', color: '#FFF' }}>
                                Update
                            </Button>
                            <Button onClick={handleCancel} style={{ marginLeft: 8, background: '#D9D9D9', color: '#909090' }}>
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
                                        <Form.Item label={<span style={{ color: '#657589' }}>Order</span>}>
                                            <Input style={{ height: 40 }} value={selectedRecord.Order} readOnly disabled />
                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>SO/INV</span>}>
                                            <Input style={{ height: 40 }} value={selectedRecord.SO_INV} disabled />
                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>SR</span>}>
                                            <Input style={{ height: 40 }} value={selectedRecord.SR} disabled />
                                        </Form.Item>
                                    </Col>
                                </Row>
                                <Row gutter={16}>
                                    <Col span={5}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>กรอกโค้ด R</span>}>
                                            <Select
                                                style={{ height: 40 }}
                                                value={codeR}
                                                onChange={setCodeR}
                                                showSearch
                                                placeholder="เลือกโค้ด R"
                                            >
                                                {codeROptions.map((code) => (
                                                    <Option key={code.value} value={code.value}>
                                                        {code.label}
                                                    </Option>
                                                ))}
                                            </Select>
                                        </Form.Item>
                                    </Col>
                                    <Col span={5}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>ชื่อของโค้ด R</span>}>
                                            <Select
                                                style={{ height: 40 }}
                                                value={nameR}
                                                onChange={setNameR}
                                                showSearch
                                                placeholder="เลือกชื่อโค้ด R"
                                            >
                                                {codeNameOptions.map((name) => (
                                                    <Option key={name.value} value={name.value}>
                                                        {name.label}
                                                    </Option>
                                                ))}
                                            </Select>
                                        </Form.Item>
                                    </Col>
                                    <Col span={5}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>QTY:</span>}>
                                            <InputNumber
                                                min={1}
                                                max={100}
                                                value={qty}
                                                onChange={(value) => setQty(value)} // Set directly from InputNumber
                                                style={{ width: '100%', height: '40px', lineHeight: '40px', }}
                                            />
                                        </Form.Item>
                                    </Col>
                                    <Col span={5}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>Price:</span>}>
                                            <InputNumber
                                                min={1}
                                                max={100000}
                                                value={price}
                                                onChange={(value) => setPrice(value)} // Set directly from InputNumber
                                                step={0.01}
                                                style={{ width: '100%', height: '40px', lineHeight: '40px', }}
                                            />
                                        </Form.Item>
                                    </Col>
                                    <Col span={4}>
                                        <Button
                                            type="primary"
                                            style={{ width: '100%', height: '40px', marginTop: 30 }}
                                            onClick={handleAdd}
                                        >
                                            <PlusCircleOutlined />
                                            Add
                                        </Button>
                                    </Col>
                                </Row>
                            </Form>

                            {/* Table to display product data */}
                            <Table
                                components={{
                                    header: {
                                        cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                            <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                                        ),
                                    },
                                }}
                                columns={[
                                    {
                                        title: 'SKU',
                                        dataIndex: 'SKU',
                                        render: (text) => <span style={{ color: '#35465B' }}>{text}</span>
                                    },
                                    { title: 'Name', dataIndex: 'Name', render: (text) => <span style={{ color: '#35465B' }}>{text}</span> },
                                    { title: 'QTY', dataIndex: 'QTY', render: (text) => <span style={{ color: '#35465B' }}>{text}</span> },
                                    { title: 'Price', dataIndex: 'Price', render: (text) => <span style={{ color: '#35465B' }}>{text}</span> },
                                    {
                                        title: 'Action',
                                        dataIndex: 'Action',
                                        render: (_, record) =>
                                            record.Type === 'addon' ? (
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
                                            ) : null
                                    },
                                ]}
                                dataSource={selectedRecord.data} // Use updated data with new entries
                                rowKey="SKU"
                                pagination={false}
                            />
                        </>
                    )}
                </Modal>
            )
            }
            {activeTabKey == '2' && (
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


                                    <Col span={8}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>Order</span>}>
                                            <Input style={{ height: 40 }} value={selectedRecord.Order} readOnly disabled />
                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>SO</span>}>
                                            <Input style={{ height: 40 }} value={selectedRecord.SO_INV} disabled />
                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>SR</span>}>
                                            <Input style={{ height: 40 }} value={selectedRecord.SR} disabled />
                                        </Form.Item>
                                    </Col>

                                </Row>
                            </Form>


                            {/* Table to display product data */}
                            <Table
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
                                            record.Type === 'addon' ? (
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
                                            ) : null
                                    },
                                ]}
                                dataSource={selectedRecord.data} // Use updated data with new entries
                                rowKey="SKU"
                                pagination={false}
                            />
                        </>

                    )}
                </Modal>
            )
            }
        </ConfigProvider>
    );
};

export default DraftandConfirm;
