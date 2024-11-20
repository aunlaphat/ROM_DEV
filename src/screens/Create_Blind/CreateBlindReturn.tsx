import { Button, Radio, Select, Space, Col, ConfigProvider, Form, Layout, Row, Input, InputNumber, Table, notification, message, Tooltip } from "antd";
import { DeleteOutlined, PlusCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import type { RadioChangeEvent } from 'antd';
import React, { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";

interface TableDataItem {
    key: string; // ค่าที่ใช้เป็น key
    SKU: string;
    SKU_Name: string;
    QTY: number;
}


const SKUName = [
    { Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", SKU: "G097171-ARM01-BL" },
    { Name: "Bewell Sport armband size M For", SKU: "G097171-ARM01-GY" },
    { Name: "Sport armband size L", SKU: "G097171-ARM02-BL" },
    { Name: "Bewell Sport armband size M with light", SKU: "G097171-ARM03-GR" },
];

// Create options for SKU and Name
const skuOptions = SKUName.map(item => ({
    value: item.SKU,
    label: item.SKU
}));

const nameOptions = SKUName.map(item => ({
    value: item.Name,
    label: item.Name
}));

const CreateBlind = () => {
    const [showInput, setShowInput] = useState(false);
    const [selectedSKU, setSelectedSKU] = useState<string | undefined>(undefined);
    const [value, setValue] = useState<number>(0);
    const [selectedName, setSelectedName] = useState<string | undefined>(undefined);
    const [form] = Form.useForm();
    const [form2] = Form.useForm();
    const [formValid, setFormValid] = useState(false);
    const [formskuValid, setFormskuValid] = useState(false);
    const [qty, setQty] = useState<number | null>(null);
    const [key, setKey] = useState<null>(null);
    const [tableData, setTableData] = useState<TableDataItem[]>([]); // ใช้ interface ที่กำหนด
    const navigate = useNavigate();
    
    const onChange = (e: RadioChangeEvent) => {
        setValue(e.target.value);
        handleFormValidation();
        setShowInput(e.target.value === 1);
    };

    const handleNavigateToTakepicture = () => {
        navigate('/Takepicture'); // เส้นทางนี้ควรตรงกับการตั้งค่า Route ใน App.js หรือไฟล์ routing ของคุณ
    };

    const handleNameChange = (value: string) => {
        const selectedOption = SKUName.find((item) => item.Name === value);
        if (selectedOption) {
            form2.setFieldsValue({
                SKU: selectedOption.SKU,
                SKU_Name: selectedOption.Name,
            });
            setSelectedName(value);
            setSelectedSKU(selectedOption.SKU);
        }
    };

    const handleSKUChange = (value: string) => {
        const selectedOption = SKUName.find((item) => item.SKU === value);
        if (selectedOption) {
            form2.setFieldsValue({
                SKU: selectedOption.SKU,
                SKU_Name: selectedOption.Name,
            });
            setSelectedSKU(value);
            setSelectedName(selectedOption.Name);
        }
    };
    


    const handleAdd = () => {
        // ตรวจสอบว่า SKU ซ้ำหรือไม่
        const isDuplicate = tableData.some(item => item.SKU === selectedSKU);
        
        if (isDuplicate) {
            notification.warning({
                message: "SKU ซ้ำ",
                description: "SKU นี้ถูกเพิ่มไปแล้ว",
            });
            return;  // หยุดการทำงานหากพบว่า SKU ซ้ำ
        }
        
        // หาก SKU ไม่ซ้ำ ให้เพิ่มข้อมูลใหม่
        const newKey = selectedSKU!;
        setTableData([...tableData, { SKU: selectedSKU!, SKU_Name: selectedName!, QTY: qty!, key: newKey }]);
        
        // รีเซ็ตฟอร์ม
        form.resetFields(['SKU', 'QTY', 'SKU_Name']);
        
        notification.success({
            message: "เพิ่มข้อมูลสำเร็จ",
        });
    };
    
    const handleDelete = (key: string) => {
        setTableData(tableData.filter(item => item.key !== key));
        notification.success({
            message: "ลบข้อมูลสำเร็จ",
        });
    };

   
    const columns = [
        {
            id:'SKU',
            title: 'SKU',
            dataIndex: 'SKU',
        },
        {
            id:'Name',
            title: 'Name',
            dataIndex: 'SKU_Name',
        },
        {
            id:'QTY',
            title: 'QTY',
            dataIndex: 'QTY',
        },
        {
            id:'Action',
            title: "Action",
            dataIndex: "Action",
            render: (_: any, record: { key: string }) => (
                <DeleteOutlined
                    style={{ cursor: 'pointer', color: 'red', fontSize: '20px' }}
                    onClick={() => handleDelete(record.key)}
                />
            ),
        },
    ];
    const formatAccountNumber = (value: string) => {
        value = value.replace(/\D/g, ""); // Remove non-digit characters
        if (value.length > 3) {
            value = value.slice(0, 3) + "-" + value.slice(3);
        }
        if (value.length > 9) {
            value = value.slice(0, 9) + "-" + value.slice(9, 10); // ปรับให้ slice ที่ 10
        }
        return value;
    };

    const checkFormValidity = () => {
        const errors = form.getFieldsError().filter(({ errors }) => errors.length);
        return errors.length === 0; // ถ้าไม่มีข้อผิดพลาด
    };

    const handleSubmit = () => {
        if (value === 2) {
            // กรณีเลือก "No"
            handleNavigateToTakepicture();
        } else if (value === 1 && tableData.length === 0) {
            notification.warning({
                message: "กรุณาเพิ่มข้อมูลในตาราง",
                description: "กรุณากรอก SKU และจำนวนอย่างน้อย 1 รายการในตารางก่อนที่จะดำเนินการต่อ!",
            });
        } else {
            handleNavigateToTakepicture();
        }
    };
    
    const handleFormValidation = () => {
        // alert('test1');
        const username = form.getFieldValue('Username');
        const phonenumber = form.getFieldValue('Phonenumber');
        const Address = form.getFieldValue('Address');
        const Tracking = form.getFieldValue('Tracking');
        const Transport = form.getFieldValue('TransportType');
        // alert(Transport);
        // console.log(allFields);
        if(username !== undefined && (phonenumber !== undefined && phonenumber?.length === 10) && Address !== undefined && Tracking !== undefined && Transport!==undefined){
            setFormValid(true); 
        }else{
            setFormValid(false); 
        }
        
    };

    const handleValuesChange = () => {
        // alert(value);
        // alert('test');
        handleFormValidation();// Trigger alert on any form field change
    };


    return (
        <ConfigProvider>
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                Create Build Return
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
                    <Form form={form} 
                     onValuesChange={handleValuesChange} 
                    layout="vertical">
                        <Row gutter={16} align="middle" justify="center" style={{ marginTop: "20px", width: '100%' }}>
                            <Col span={8}>
                                <Form.Item
                                id="Username"
                                    label={<span style={{ color: '#657589' }}>กรอกชื่อลูกค้า</span>}
                                    name="Username"
                                    rules={[{ required: true, message: 'กรุณากรอกชื่อลูกค้า Order!' }]}
                                >
                                    <Input style={{ height: 40 }} placeholder="กรอกชื่อลูกค้า" />
                                </Form.Item>
                            </Col>
                            <Col span={8}>
                                <Form.Item
                                  id="Phonenumber"
                                    label={<span style={{ color: '#657589' }}>กรอกเบอร์โทร</span>}
                                    name="Phonenumber"
                                    rules={[{
                                        required: true, message: 'กรุณากรอกเบอร์โทร!'

                                    },
                                    
                                        {
                                            len: 10,
                                            message: 'กรุณากรอกเบอร์โทรให้ครบ 10 หลัก!',
                                        
                                    }

                                    ]}
                                >
                                    
                                    <Input
                                        type="number"
                                        style={{ height: 40 }}
                                        placeholder="กรอกเบอร์โทร"
                                        maxLength={10}
                                        onChange={(e) => {
                                            let value = e.target.value;
                                
                                            // Limit the value to 10 characters and remove non-numeric characters
                                            if (value.length > 10) {
                                                value = value.slice(0, 10);
                                            }
                                
                                            // Optionally format the value (e.g., adding spaces or dashes)
                                            // For example, here we remove all non-numeric characters for simplicity
                                            value = value.replace(/\D/g, '');
                                
                                            // Set the value back to the input field
                                            e.target.value = value;
                                        }}


                                    />
                                </Form.Item>
                            </Col>
                            <Col span={8}>
                                <Form.Item
                                 id="Address"
                                    label={<span style={{ color: '#657589' }}>กรอกที่อยู่</span>}
                                    name="Address"
                                    rules={[{ required: true, message: 'กรุณากรอกที่อยู่!' }]}
                                >
                                    <Input style={{ height: 40 }} placeholder="กรอกที่อยู่" />
                                </Form.Item>
                            </Col>
                        </Row>

                        <Row gutter={16} align="middle" justify="center" style={{ marginTop: "20px", width: '100%' }}>
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
                                    name="Tracking"
                                    rules={[{ required: true, message: "กรอกเลข Tracking" }]}>
                                    <Input style={{ height: 40 }} placeholder="กรอกเลข Tracking" />
                                </Form.Item>
                            </Col>
                            <Col span={8}>
                                <Form.Item
                                  id="Ordernumber"
                                    label={<span style={{ color: '#657589' }}>กรอกเลข Order</span>}
                                    name="Ordernumber"

                                >
                                    <Input style={{ height: 40 }} placeholder="กรอกเลข Order" />
                                </Form.Item>
                            </Col>
                            <Col span={8}>
                                <Form.Item
                                   id="TransportType"
                                    label={<span style={{ color: '#657589' }}>Transport Type:</span>}
                                    name="TransportType"
                                    rules={[{ required: true, message: 'กรุณาเลือก Transport Type' }]}
                                >
                                    <Select
                                        style={{ width: '100%', height: '40px' }}
                                        showSearch
                                        placeholder="Transport Type"
                                        options={[
                                            { value: '1', label: 'SPX Express' },
                                            { value: '2', label: 'J&T Express' },
                                            { value: '3', label: 'Flash Express' },
                                            { value: '4', label: 'Shopee' },
                                            { value: '5', label: 'NocNoc' },
                                        ]}
                                    />
                                </Form.Item>
                            </Col>
                        </Row>
                    </Form>
                    <Form
    form={form2}
    layout="vertical"
    onValuesChange={() => {
        const { SKU, SKU_Name, QTY } = form2.getFieldsValue();
        console.log(SKU,SKU_Name,QTY);
        // ตรวจสอบว่าทุกฟิลด์มีค่าไม่เป็น undefined และไม่เป็นค่าว่าง
        const isFormValid = SKU !== undefined && SKU_Name !== undefined && QTY !== undefined;
        setFormskuValid(isFormValid);
    }}
>
    <Row gutter={16} align="middle" justify="center" style={{ marginTop: "20px", width: '100%' }}>
        {/* ... ส่วนอื่น ๆ ของฟอร์ม ... */}
    </Row>

    <Row align="middle" justify="start" style={{ marginTop: "20px", width: '100%' }}>
        <div style={{ marginRight: "10px" }}>แกะกล่อง</div>
        <Radio.Group   id="Radio" onChange={onChange} value={value}>
            <Radio value={1}>Yes</Radio>
            <Radio value={2}>No</Radio>
        </Radio.Group>
    </Row>

    {showInput && (
        <Row gutter={16} style={{ marginTop: "20px", width: '100%' }}>
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
                    placeholder="เลือก SKU"
                    value={selectedSKU}
                    onChange={handleSKUChange}
                    options={skuOptions}
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
                    placeholder="เลือก SKU Name"
                    value={selectedName}
                    onChange={handleNameChange}
                    options={nameOptions}
                />
            </Form.Item>
        </Col>
        <Col span={4}>
            <Form.Item
                id="qty" 
                label={<span style={{ color: '#657589' }}>QTY:</span>}
                name="QTY"
                rules={[{ required: true, message: 'กรุณากรอกจำนวน' }]}
            >
                <InputNumber
                    min={1}
                    max={100}
                    value={qty}
                    onChange={(value) => setQty(value)}
                    style={{ width: '100%', height: 40, lineHeight: '40px' }}
                />
            </Form.Item>
                            </Col>
                            <Col span={4}>
            <Button
                id="Add" 
                type="primary"
                disabled={!formskuValid}  // ปิดการใช้งานเมื่อ form ไม่ valid
                onClick={handleAdd}
                style={{ width: '100%', height: '40px', marginTop: 30 }}
            >
                <PlusCircleOutlined /> {/* เพิ่มไอคอนที่นี่ */}
                Add
            </Button>
        </Col>
                        </Row>
                        
                        )}

                        {showInput && (
                            <Table
                            id="Table" 
                                style={{ marginTop: '20px' }}
                                columns={columns}
                                dataSource={tableData}
                                rowKey="SKU"
                                pagination={false}
                            />
                        )}
                        
                    </Form>
                    <Row align="middle" justify="center" style={{ marginTop: "20px", width: '100%' }}> 
                    <Button
                        id="Confirm" 
            type="primary"
            onClick={handleSubmit}
            disabled={(value === 2 && formValid === true ) ? false : (value === 1 && tableData.length > 0&& formValid === true ) ? false :true}
            // disabled={tableData.length === 0 || !formValid || value === undefined}  // ปิดการใช้งานเมื่อไม่มีข้อมูลในตาราง, form ไม่ valid หรือไม่มีการเลือก value
        >
            ยืนยัน
        </Button>
    </Row>
                </Layout.Content>
            </Layout>
        </ConfigProvider>
    );
};

export default CreateBlind;
