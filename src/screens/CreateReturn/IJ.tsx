import React, { useState } from 'react';
import { notification, Form, Input, InputNumber, DatePicker, Button, Row, Col, Table, ConfigProvider, Layout, Select, Modal, message, Popconfirm, Divider, Tooltip } from 'antd';
import moment from 'moment';
import { DeleteOutlined, LeftOutlined, PlusCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { Name } from 'ajv';

const handleError = (error: any) => {
    notification.warning({
        message: "กรุณากรอกข้อมูลให้ครบ",

    });
};

interface FormValues {
    SKU: string;
    QTY: number;
    SKU_Name: string;
}
interface User {
    ID: number;
    Name: string;
    role: 'Warehouse' | 'Accounting';
}


interface DataSourceItem extends FormValues {
    key: number;
    warehouse_form?: string;
    location_form?: string;
    warehouse_to?: string;
}
const SKUName = [
    { Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", SKU: "G097171-ARM01-BL" },
    { Name: "Bewell Sport armband size M For", SKU: "G097171-ARM01-GY" },
    { Name: "Sport armband size L", SKU: "G097171-ARM02-BL" },
    { Name: "Bewell Sport armband size M with light", SKU: "G097171-ARM03-GR" },
];

const MockUser: User[] = [
    { ID: 1, Name: "User 1", role: "Warehouse" },
    { ID: 2, Name: "User 2", role: "Accounting" },

];


// สร้าง options สำหรับ SKU
const skuOptions = SKUName.map(item => ({
    value: item.SKU,  // SKU เป็นค่า value
    label: item.SKU   // SKU เป็น label เพื่อแสดงใน dropdown
}));

// สร้าง options สำหรับ SKU Name
const nameOptions = SKUName.map(item => ({
    value: item.Name, // Name เป็นค่า value
    label: item.Name  // Name เป็น label เพื่อแสดงใน dropdown
}));

const IJPage: React.FC = () => {
    const [selectedSKU, setSelectedSKU] = useState<string | undefined>(undefined);
    const [selectedName, setSelectedName] = useState<string | undefined>(undefined);
    const [form] = Form.useForm();
    const [dataSource, setDataSource] = useState<DataSourceItem[]>([]);
    const [formValid, setFormValid] = useState(false);
    const [isSubmitted, setIsSubmitted] = useState(false);
    const navigate = useNavigate();
    const [formDisabled, setFormDisabled] = useState(false);
    const [ij, setIJ] = useState<string>('');
    const [remark, setRemark] = useState<string>('');
    const [submittedRemark, setSubmittedRemark] = useState<string>('');
    const [qty, setQty] = useState<number | null>(null);  // Allow null

    const handleIJChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setIJ(e.target.value);
    };

    const handleRemarkChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setRemark(e.target.value);
    };


    const onChange = () => {
        const values = form.getFieldsValue();
        const { Date, SKU, QTY } = values;

        // Set form validity based on required fields
        setFormValid(Date && SKU && QTY);
    };


    const handleAdd = () => {
        // ตรวจสอบการกรอกข้อมูลที่จำเป็น เช่น วันที่คืน, ประเภทการขนส่ง, SKU, ชื่อสินค้า, และ QTY
        form.validateFields(['Date', 'TransportType', 'SKU', 'SKU_Name', 'QTY'])
            .then((values) => {
                // ถ้าข้อมูลในฟิลด์เหล่านี้ไม่ครบ จะมีข้อความเตือนขึ้น
                if (!values.Date || !values.TransportType || !values.SKU || !values.SKU_Name || !values.QTY) {
                    notification.warning({
                        message: "มีข้อสงสัย",
                        description: "กรุณากรอกข้อมูลที่จำเป็นให้ครบก่อนเพิ่ม!",
                    });
                    return;
                }

                // ตรวจสอบว่า SKU ที่กรอกมีอยู่ใน dataSource หรือไม่
                const isSKUExist = dataSource.some(item => item.SKU === values.SKU);

                if (isSKUExist) {
                    // แสดงข้อความเตือนว่า SKU ซ้ำ
                    notification.warning({
                        message: "มีข้อผิดพลาด",
                        description: "SKU นี้ถูกเพิ่มไปแล้วในรายการ!",
                    });
                    return; // ไม่ทำการเพิ่มข้อมูล
                }

                // ถ้า SKU ยังไม่ซ้ำ เพิ่มข้อมูลใหม่
                setDataSource((prevData) => [
                    ...prevData,
                    { key: Date.now(), ...values }, // ใช้ Date.now() เพื่อสร้าง key ใหม่
                ]);

                notification.success({
                    message: "Add ข้อมูลสำเร็จ",
                });

                // รีเซ็ตฟิลด์ที่กรอกไว้แล้ว
                form.resetFields(['SKU', 'SKU_Name', 'QTY']);
            })
            .catch((errorInfo) => {
                // หากการตรวจสอบฟอร์มไม่ผ่าน จะโชว์ข้อความเตือน
                notification.warning({
                    message: "มีข้อสงสัย",
                    description: "กรุณากรอกข้อมูลให้ครบก่อนเพิ่ม!",
                });
            });
    };


    const handleDelete = (key: number) => {
        Modal.confirm({
            title: 'ยืนยันการลบ',
            content: 'คุณต้องการลบรายการนี้ใช่หรือไม่?',
            okText: 'ใช่',
            okType: 'danger',
            cancelText: 'ไม่',
            onOk() {
                setDataSource((prevData) => prevData.filter(item => item.key !== key));
            },
        });
    };

    const handleChange = (value: string, key: number, field: string) => {
        const updatedDataSource = dataSource.map((item) => {
            if (item.key === key) {
                return { ...item, [field]: value };
            }
            return item;
        });
        setDataSource(updatedDataSource);
    };

    const columns = [
        {
            title: 'SKU',
            dataIndex: 'SKU',
        },
        {
            title: 'Name',
            dataIndex: 'SKU_Name',
        },
        {
            title: 'QTY',
            dataIndex: 'QTY',
        },
        {
            title: "Warehouse Form",
            dataIndex: "warehouse_form",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '100%' }}
                    onChange={(value) => handleChange(value, record.key, "warehouse_form")}
                    options={Warehouse}
                    dropdownStyle={{ minWidth: 120 }}
                    dropdownMatchSelectWidth={false}
                    maxTagTextLength={50} // กำหนดความยาวสูงสุดของข้อความในตัวเลือก
                    disabled={formDisabled}
                />

            ),
        },
        {
            title: "Location Form",
            dataIndex: "location_form",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '100%' }}
                    onChange={(value) => handleChange(value, record.key, "location_form")}
                    options={Location}
                    dropdownStyle={{ minWidth: 120 }}
                    dropdownMatchSelectWidth={false}
                    maxTagTextLength={50} // กำหนดความยาวสูงสุดของข้อความในตัวเลือก
                    disabled={formDisabled}
                />
            ),
        },
        {
            title: "Warehouse to",
            dataIndex: "warehouse_to",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '100%' }}
                    onChange={(value) => handleChange(value, record.key, "warehouse_to")}
                    options={Warehouseto}
                    dropdownStyle={{ minWidth: 100 }}
                    dropdownMatchSelectWidth={false}
                    maxTagTextLength={50} // กำหนดความยาวสูงสุดของข้อความในตัวเลือก
                    disabled={formDisabled}
                />
            ),
        },
        {
            title: "Action",
            dataIndex: "Action",
            render: (_: any, record: DataSourceItem) => (
                <DeleteOutlined
                    style={{ cursor: 'pointer', color: 'red', fontSize: '20px' }}
                    onClick={() => handleDelete(record.key)}
                />
            ),
        },
    ];

    const Warehouse = [
        { value: "W1" }, { value: "W2" }, { value: "W3" }, { value: "W4" },
    ];

    const Warehouseto = [
        { value: "BES" }, { value: "MMT_BEWELL" },
    ];

    const Location = [
        { value: "L1" }, { value: "L2" }, { value: "L3" }, { value: "L4" },
    ];


    const [selectedValue, setSelectedValue] = useState<string | undefined>();

    const handleSelectChange = (value: string) => {
        const selectedOption = SKUName.find((val) => val.SKU === value);

        if (selectedOption) {
            form.setFieldsValue({
                SKU: selectedOption.SKU,
                SKU_Name: selectedOption.Name,
            });
            setSelectedValue(value);
        }
    };
    const handleSKUChange = (value: string) => {
        const selectedOption = SKUName.find((val) => val.SKU === value);
        if (selectedOption) {
            form.setFieldsValue({
                SKU: selectedOption.SKU,
                SKU_Name: selectedOption.Name,
            });
            setSelectedSKU(value);
            setSelectedName(selectedOption.Name); // อัปเดต selectedName
        }
    };

    const handleNameChange = (value: string) => {
        const selectedOption = SKUName.find((val) => val.Name === value);
        if (selectedOption) {
            form.setFieldsValue({
                SKU: selectedOption.SKU,
                SKU_Name: selectedOption.Name,
            });
            setSelectedName(value);
            setSelectedSKU(selectedOption.SKU); // อัปเดต selectedSKU
        }
    };



    const onSearch = (value: string) => {
        console.log('search:', value);
    };

    const generateRandomNumber = () => {
        return Math.floor(Math.random() * 10000);
    };

    const handleCreateIJ = () => {

        try {
            const isDataValid = dataSource.every((record: DataSourceItem) =>
                record.warehouse_form && record.location_form && record.warehouse_to
            );

            if (!isDataValid) {
                notification.warning({
                    message: "กรุณากรอกข้อมูลให้ครบก่อนสร้าง IJ",
                    description: "กรุณาเลือก Warehouse Form, Location Form, Warehouse to",
                });
                return;
            }

            const randomNumber = generateRandomNumber();
            form.setFieldsValue({ IJ_Create: randomNumber });
            setDataSource((prevData) => prevData.map(item => ({ ...item, IJ_Create: randomNumber })));

            // แสดง notification ว่าสร้างสำเร็จและเลขสุ่มที่สร้างขึ้น
            notification.success({
                message: 'สำเร็จ',
                description: `Create IJ สำเร็จ! เลขสุ่มที่สร้างคือ: ${randomNumber}`,
            });

            setIsSubmitted(true);
        } catch (error) {
            handleError(error); // ใช้ handleError ในกรณีที่เกิด error
        }
        setFormDisabled(true);
    };


    const handleSubmitData = () => {
        const combinedRemark = `${remark} - IJ: ${ij}`;
        console.log("IJ:", ij);
        console.log("Remark:", combinedRemark);

        console.log("Sending data:", dataSource);

        // Reset all form fields and state
        setDataSource([]);
        form.resetFields();
        setRemark('');
        setIJ('');
        setIsSubmitted(false);
        setFormDisabled(false); // เปิดใช้งานฟอร์มใหม่อีกครั้ง

        notification.success({
            message: 'ส่งข้อมูล สำเร็จ',
            description: 'ข้อมูลทั้งหมดได้ถูกส่งเรียบร้อยแล้ว',
        });
    };



    const handleBack = () => {
        navigate('/CreateReturn'); // Navigate to CreateReturn page
    };
    const handleCancel = () => {
        form.resetFields();       // รีเซ็ตค่าในฟอร์มทั้งหมด
        setDataSource([]);        // รีเซ็ตข้อมูล dataSource
        setIsSubmitted(false);
        notification.success({
            message: 'Cancel สำเร็จ',
            description: 'ข้อมูลทั้งหมดได้ถูกยกเลิกเรียบร้อยแล้ว',
        });   // รีเซ็ตสถานะ isSubmitted

    };


    return (
        <ConfigProvider >
            <div  id="titleContainer" style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                Create IJ Return
            </div>
            <Layout id="layout">
                <Layout.Content
                id="contentContainer"
                    style={{
                        margin: "24px",
                        padding: 36,
                        minHeight: 360,
                        background: "#fff",
                        borderRadius: "8px",
                        overflow: "auto",
                    }}
                >
                    <div id="mainContent">
                        <Button
                         id="backButton"
                            onClick={handleBack}
                            style={{ background: '#98CEFF', color: '#fff' }}
                        >
                            <LeftOutlined style={{ color: '#fff', marginRight: 5 }} />
                            Back
                        </Button>

                        <Form
                        id="form"
                            form={form}
                            layout="vertical"
                            onValuesChange={onChange}
                            style={{ padding: '20px', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}
                        >
                            <div id="formContainer" style={{ width: '100%', maxWidth: '800px' }}> {/* Adjust max-width here */}

                                <Divider style={{ color: '#657589', fontSize: '22px', marginTop: 30, marginBottom: 30 }} orientation="left"> IJ document Information </Divider>
                                <Row gutter={16} >

                                    <Col span={8}>
                                        <Form.Item
                                          id="ijDocumentInput"
                                            label={<span style={{ color: '#657589' }}>กรอกเอกสารอ้างอิง IJ (ไม่บังคับ):</span>}
                                            name="IJ"

                                        >
                                            <Input id="Doc" style={{ width: '100%', height: '40px', }} placeholder="กรอกเอกสารอ้างอิง" onChange={handleRemarkChange} disabled={formDisabled} />
                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item
                                         id="ijCreateInput"
                                            label={
                                                <span style={{ color: '#657589' }}>
                                                    IJ Create:&nbsp;
                                                    <Tooltip title="กด create IJ ระบบจะส่งคำสั่งสร้าง เข้า AX แล้วจะได้เลข IJ">
                                                        <QuestionCircleOutlined style={{ color: '#657589' }} />
                                                    </Tooltip>
                                                </span>
                                            }
                                            name="IJ_Create"
                                        >
                                            <Input  style={{ width: '100%', height: '40px', }} placeholder="IJ Create" disabled={true} />
                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item
                                            id="remarkInput"
                                            label={<span style={{ color: '#657589' }}>Remark (ไม่บังคับ):</span>}
                                            name="Remark"
                                        >
                                            <Input style={{ width: '100%', height: '40px', }} showCount maxLength={200} onChange={handleIJChange} disabled={formDisabled} />
                                        </Form.Item>
                                    </Col>
                                </Row>
                                <Row gutter={16} >
                                    <Divider style={{ color: '#657589', fontSize: '22px', marginTop: 30, marginBottom: 30 }} orientation="left"> Transport Information </Divider>
                                    <Col span={8}>
                                        <Form.Item
                                        id="Date"
                                            label={<span style={{ color: '#657589' }}>วันที่คืน:</span>}
                                            name="Date"
                                            rules={[{ required: true, message: 'กรุณาเลือกวันที่คืน' }]}
                                        >
                                            <DatePicker style={{ width: '100%', height: '40px', }} placeholder="เลือกวันที่คืน" disabled={formDisabled} />
                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item
                                         id="Tracking"
                                            label={
                                                <span style={{ color: '#657589' }}>
                                                    กรอกเลข Tracking:&nbsp;
                                                    <Tooltip title="เลขTracking จากขนส่ง">
                                                        <QuestionCircleOutlined style={{ color: '#657589' }} />
                                                    </Tooltip>
                                                </span>
                                            }


                                            name="TrackingNumber"
                                            rules={[{ required: true, message: 'กรุณากรอกเลข Tracking' }]}
                                        >
                                            <Input style={{ width: '100%', height: '40px', }} placeholder="เลขTracking" disabled={formDisabled} />

                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item
                                         id="TransportType"
                                            label={<span style={{ color: '#657589' }}>Transport Type:</span>}
                                            name="TransportType"
                                            rules={[{ required: true, message: "กรุณาเลือก Transport Type" }]}
                                        >
                                            <Select
                                                style={{ width: '100%', height: '40px' }}
                                                showSearch
                                                placeholder="TransportType"
                                                optionFilterProp="label"
                                                onChange={onChange}
                                                onSearch={onSearch}
                                                options={[
                                                    { value: 'SPX Express', label: 'SPX Express' },
                                                    { value: 'J&T Express', label: 'J&T Express' },
                                                    { value: 'Flash Express', label: 'Flash Express' },
                                                    { value: 'Shopee', label: 'Shopee' },
                                                    { value: 'NocNoc', label: 'NocNoc' },

                                                ]}
                                                disabled={formDisabled}
                                            />
                                        </Form.Item>
                                    </Col>
                                </Row>
                                <Divider style={{ color: '#657589', fontSize: '22px', marginTop: 30, marginBottom: 30 }} orientation="left"> SKU Information </Divider>
                                <Row gutter={16} >
                                    <Col span={8}>
                                        <Form.Item
                                        id="Sku"
                                            label={<span style={{ color: '#657589' }}>กรอก SKU:</span>}
                                            name="SKU"
                                            rules={[{ required: true, message: "กรุณากรอก SKU" }]}
                                        >
                                            <Select
                                                showSearch
                                                style={{ width: '100%', height: '40px' }}
                                                placeholder="Search to Select"
                                                optionFilterProp="label"
                                                value={selectedSKU} // แสดง SKU ที่ถูกเลือก
                                                onChange={handleSKUChange}
                                                options={skuOptions} // แสดง SKU ใน dropdown
                                                disabled={formDisabled}

                                            />
                                        </Form.Item>
                                    </Col>
                                    <Col span={8}>
                                        <Form.Item
                                          id="SkuName"
                                            label={<span style={{ color: '#657589' }}>Name:</span>}
                                            name="SKU_Name"
                                            rules={[{ required: true, message: "กรุณาเลือก SKU Name" }]}
                                        >
                                            <Select
                                                showSearch
                                                style={{ width: '100%', height: '40px' }}
                                                placeholder="Search to Select"
                                                optionFilterProp="label"
                                                value={selectedName} // แสดง SKU Name ที่ถูกเลือก
                                                onChange={handleNameChange}
                                                options={nameOptions} // แสดง SKU Name ใน dropdown
                                                disabled={formDisabled}
                                            />
                                        </Form.Item>
                                    </Col>
                                    <Col span={4}>
                                        <Form.Item label={<span style={{ color: '#657589' }}>QTY:</span>}
                                             id="qty"
                                           name="QTY"
                                            rules={[{ required: true, message: "กรุณากรอก QTY" }]}>
                                            <InputNumber
                                                min={1}
                                                max={100}
                                                value={qty}
                                                onChange={(value) => setQty(value)} // Set directly from InputNumber
                                                style={{ width: '100%', height: '40px', lineHeight: '40px', }}
                                            />
                                        </Form.Item>
                                    </Col>
                                    <Col span={4}>

                                        <Button
                                          id="addsku"
                                            type="primary"
                                            disabled={!formValid || isSubmitted}
                                            onClick={handleAdd}
                                            style={{
                                                width: '100%',
                                                height: '40px',
                                                marginTop: 30,
                                                display: 'flex',
                                                alignItems: 'center',
                                                justifyContent: 'center',
                                            }}
                                        >
                                            <PlusCircleOutlined style={{ marginLeft: 1 }} /> {/* ลดระยะห่างระหว่างไอคอนและข้อความ */}
                                            Add SKU
                                        </Button>

                                    </Col>
                                </Row>



                            </div>
                        </Form>

                        <div >
                            <Table
                            id="table"
                                components={{
                                    header: {
                                        cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                            <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                                        ),
                                    },
                                }}
                                dataSource={dataSource}
                                columns={columns}
                                rowKey="key"
                                pagination={false} // Disable pagination if necessary
                                style={{ width: '100%', tableLayout: 'fixed' }} // Ensure the table takes full width and is fixed layout
                                scroll={{ x: 'max-content' }}

                            />
                        </div>
                    </div>
                    <Row justify="center" style={{ marginTop: '20px' }}>
                        <Popconfirm
                        id="popconfirmforreateIJ,sendData"
                            title={isSubmitted ? "ยืนยันการส่งข้อมูล" : "ยืนยันการสร้าง IJ"}
                            description={isSubmitted ? "คุณต้องการส่งข้อมูลนี้ใช่หรือไม่?" : "คุณต้องการสร้าง IJ ใช่หรือไม่?"}
                            onConfirm={isSubmitted ? handleSubmitData : handleCreateIJ} // ฟังก์ชันตามสถานะ isSubmitted
                            okText="ใช่"
                            cancelText="ไม่"
                        >
                            <Button
                            id="createIJ,sendData"
                                type="primary"
                                style={{ width: 100, height: 40, marginRight: '10px' }} // เพิ่มช่องว่างทางขวา
                                disabled={!formValid || dataSource.length === 0}
                            >
                                {isSubmitted ? "ส่งข้อมูล" : "Create IJ"}
                            </Button>
                        </Popconfirm>

                        <Popconfirm
                          id="popconfirmforcancel"
                            title="ต้องการยกเลิกหรือไม่?"
                            description="คุณแน่ใจหรือไม่ว่าต้องการยกเลิกข้อมูลทั้งหมด?"
                            onConfirm={handleCancel} // ยืนยันการยกเลิก
                            okText="ใช่"
                            cancelText="ไม่"
                        >
                            <Button
                            id="cancel"
                                type="default"
                                style={{ width: 100, height: 40 }}
                            >
                                Cancel
                            </Button>
                        </Popconfirm>
                    </Row>




                </Layout.Content>
            </Layout>
        </ConfigProvider>
    );
};

export default IJPage;


