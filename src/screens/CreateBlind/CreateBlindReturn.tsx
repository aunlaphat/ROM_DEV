import { Popconfirm, Button, Radio, Select, Space, Col, ConfigProvider, Form, Layout, Row, Input, InputNumber, Table, notification, message, Tooltip } from "antd";
import { DeleteOutlined, PlusCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import type { RadioChangeEvent } from 'antd';
import { debounce } from "lodash";
import React, { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";
import { DataItemBlind, Product } from '../../types/types';
import api from "../../utils/axios/axiosInstance"; 
const { Option } = Select;

const CreateBlind = () => {
    const [showInput, setShowInput] = useState(false);
    const [value, setValue] = useState<number>(0);
    const [form] = Form.useForm();
    const [form2] = Form.useForm();
    const [formValid, setFormValid] = useState(false);
    const [formskuValid, setFormskuValid] = useState(false);
    const [key, setKey] = useState<null>(null);
    const [tableData, setTableData] = useState<DataItemBlind[]>([]); 
    const [loading, setLoading] = useState(false);
    const [dataSource, setDataSource] = useState<DataItemBlind[]>([]);

    const [skuOptions, setSkuOptions] = useState<Product[]>([]); 
    const [nameOptions, setNameOptions] = useState<Product[]>([]); 
    const [selectedSKU, setSelectedSKU] = useState<string | undefined>(undefined);
    const [selectedName, setSelectedName] = useState<string | undefined>(undefined);
    const [price, setPrice] = useState<number | null>(null); 
    const [qty, setQty] = useState<number | null>(null); 
    
    const onChange = (e: RadioChangeEvent) => {
        setValue(e.target.value);
        handleFormValidation();
        setShowInput(e.target.value === 1);
    };

    const navigate = useNavigate();
    // const handleNavigateToTakepicture = () => {
    //     navigate('/Takepicture'); // เส้นทางนี้ควรตรงกับการตั้งค่า Route ใน App.js หรือไฟล์ routing ของคุณ
    // };

    const handleNavigateToTakepicture = (orderNumber: string, dataSource: DataItemBlind[], value: number) => {
        navigate('/Takepicture', { state: { orderNumber, dataSource, value } });
    };

    /*** SKU&NameAlias ***/
    const debouncedSearchSKU = debounce(async (value: string, searchType: string) => {
        setLoading(true);
        try {
        const response = await api.get("/api/constants/search-product", {
            params: {
            keyword: value,
            searchType,
            offset: 0,
            limit: 50,
            },
        });

        const products = response.data.data;
        if (searchType === "SKU") {
            setSkuOptions(products.map((product: Product) => ({
            sku: product.sku,
            nameAlias: product.nameAlias,
            size: product.size,
            })));
        } else if (searchType === "NAMEALIAS") {
            setNameOptions(products.map((product: Product) => ({
            sku: product.sku,
            nameAlias: product.nameAlias,
            size: product.size,
            })));
        }
        } catch (error) {
        console.error("Error fetching products:", error);
        notification.error({
            message: "Error",
            description: "There was an error fetching product data.",
        });
        } finally {
        setLoading(false);
        }
    }, 1000);

    const handleSearchSKU = (value: string) => {
        debouncedSearchSKU(value, "SKU");
    };

    const handleSearchNameAlias = (value: string) => {
        debouncedSearchSKU(value, "NAMEALIAS");
    };

    const handleNameChange = async (value: string) => {
        if (!value) {
            form.setFieldsValue({ SKU: "", SKU_Name: "" });
            setSkuOptions([]);
            setNameOptions([]);
            return;
        }

        const [nameAlias, size] = value.split("+"); // แยกค่า nameAlias และ size โดยใช้ `+`

        try {
        setLoading(true);
        const response = await api.get("/api/constants/get-sku", {
            params: { nameAlias, size },
        });

        const products = response.data.data;
        if (products.length > 0) {
            setSkuOptions(products.map((product: Product) => ({
            sku: product.sku,
            nameAlias: product.nameAlias,
            size: product.size,
            })));
            form2.setFieldsValue({
            SKU: products[0].sku, 
            });
        } else {
            console.warn("No SKU found for:", nameAlias, size);
            setSkuOptions([]); 
            setNameOptions([]); 
            form2.setFieldsValue({ SKU: "", SKU_Name: "" });
        }
        } catch (error) {
            console.error("Error fetching SKU:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleSKUChange = (value: string) => {
        const selected = skuOptions.find((option) => option.sku === value);
        if (selected) {
            form2.setFieldsValue({
                SKU: selected.sku,
                SKU_Name: selected.nameAlias,
            });
            setSelectedSKU(selected.sku);
            setSelectedName(selected.nameAlias);

            // อัปเดต nameOptions ตาม SKU ที่เลือก
            const filteredNameOptions = skuOptions
            .filter((option) => option.sku === selected.sku) // กรองเฉพาะ SKU ที่ตรงกับที่เลือก
            .map((option) => ({
            ...option,  // คัดลอกค่าเดิม
            Key: option.sku,  // เพิ่มคีย์ Key ที่ต้องการ
            }));
            setNameOptions(filteredNameOptions);  // อัปเดต nameOptions
        } else { 
            setSkuOptions([]); 
            setNameOptions([]); 
            setSelectedSKU("");
            setSelectedName("");
        }
    };

    useEffect(() => {
        if (selectedSKU) {
            const selected = skuOptions.find((option) => option.sku === selectedSKU);
            if (selected) {
                form2.setFieldsValue({
                    SKU: selected.sku,
                    SKU_Name: selected.nameAlias,
                });
            }
        }
    }, [selectedSKU, skuOptions, form2]);
    
    const handleAdd = () => {
        form2.validateFields()
            .then((values) => {
                // ตรวจสอบว่า SKU_Name มีค่าและมีเครื่องหมาย '+'
                if (values.SKU_Name || values.SKU_Name.includes('+')) {
                    const [nameAlias, size] = values.SKU_Name.split('+');
                    // ตรวจสอบว่า SKU ที่กรอกมีอยู่ใน dataSource หรือไม่
                    const isSKUExist = dataSource.some((item) => item.SKU === values.SKU);
    
                    if (isSKUExist) {
                        notification.warning({
                            message: "มีข้อผิดพลาด",
                            description: "SKU นี้ถูกเพิ่มไปแล้วในรายการ!",
                        });
                        return; 
                    }
    
                    // ถ้า SKU ยังไม่ซ้ำ เพิ่มข้อมูลใหม่
                    const newData = {
                        key: dataSource.length + 1,
                        SKU: values.SKU,
                        Name: nameAlias,
                        QTY: values.QTY,
                    };
                    setDataSource([...dataSource, newData]); // เพิ่มข้อมูลใหม่ไปยัง dataSource
                    notification.success({
                        message: "เพิ่มสำเร็จ",
                        description: "ข้อมูลของคุณถูกเพิ่มเรียบร้อยแล้ว!",
                    });
    
                    // รีเซ็ตฟอร์ม
                    form2.resetFields(['SKU', 'SKU_Name', 'QTY']);
                    setSkuOptions([]);
                    setNameOptions([]);
                    setSelectedSKU("");
                    setSelectedName("");
                } else {
                    notification.warning({
                        message: "มีข้อผิดพลาด",
                        description: "กรุณาเลือก SKU Name ที่ถูกต้อง!",
                    });
                }
            })
            .catch((info) => {
                console.log("Validate Failed:", info);
                notification.warning({
                    message: "มีข้อสงสัย",
                    description: "กรุณากรอกข้อมูลให้ครบก่อนเพิ่ม!",
                });
            });
    };

    const handleDelete = (key: number) => {
        setDataSource(dataSource.filter((item) => item.key !== key));
        notification.success({
          message: "ลบข้อมูลสำเร็จ",
          description: "ข้อมูลของคุณถูกลบออกเรียบร้อยแล้ว.",
        });
    };

    const columns = [
        { title: "SKU", dataIndex: "SKU", id: "SKU" },
        { title: "Name", dataIndex: "Name", id: "Name" },
        { title: "QTY", dataIndex: "QTY", id: "QTY" },
        {
            title: "Action",
            id:'Action',
            dataIndex: "Action",
            render: (_: any, record: { key: number }) => (
                <Popconfirm
                title="คุณแน่ใจหรือไม่ว่าต้องการลบข้อมูลนี้?"
                onConfirm={() => handleDelete(record.key)} // เรียกใช้ฟังก์ชัน handleDelete เมื่อกดยืนยัน
                okText="ใช่"
                cancelText="ไม่"
              >
                <DeleteOutlined
                  style={{ cursor: "pointer", color: "red", fontSize: "20px" }}
                />
              </Popconfirm>
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
        form.validateFields()
        .then((values) => {
            if (value === 1) {
                form2.validateFields()
                    .then((values2) => {
                        if (dataSource.length === 0) {
                            notification.warning({
                                message: "กรุณาเพิ่มข้อมูลในตาราง",
                                description: "กรุณากรอกรายการสินค้าอย่างน้อย 1 รายการก่อนส่งข้อมูล !",
                            });
                            return; 
                        }
                    })
                    .catch((info) => {
                        console.log("Validate Failed:", info);
                        notification.warning({
                            message: "กรุณาเพิ่มข้อมูลในตาราง",
                            description: "กรุณากรอกรายการสินค้าอย่างน้อย 1 รายการก่อนส่งข้อมูล !",
                        });
                    });
            } else {
                console.log("Table Data:", dataSource);
                form.resetFields(); 
                form2.resetFields();
                setDataSource([]);
                notification.success({
                    message: "ส่งข้อมูลสำเร็จ",
                    description: "ข้อมูลของคุณถูกส่งเรียบร้อยแล้ว!",
                });
                // handleNavigateToTakepicture(); // ไปหน้าต่อไป
                handleNavigateToTakepicture(values.OrderNo, dataSource, value);
            }
        })
        .catch((info) => {
            console.log("Validate Failed:", info);
            notification.warning({
                message: "คำเตือน",
                description: "กรุณากรอกข้อมูลให้ครบก่อนส่งข้อมูล !",
            });
        });
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
                                    rules={[{ required: true, message: 'กรุณากรอกเบอร์โทร!' }, { len: 10, message: 'กรุณากรอกเบอร์โทรให้ครบ 10 หลัก!', }]}
                                >
                                    <Input
                                        type="number"
                                        style={{ height: 40 }}
                                        placeholder="กรอกเบอร์โทร"
                                        maxLength={10}
                                        onChange={(e) => {
                                            let value = e.target.value;
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
                                    id="OrderNo"
                                    label={<span style={{ color: '#657589' }}>กรอกเลข Order</span>}
                                    name="OrderNo"
                                    rules={[{ required: true, message: "กรอกเลข Order" }]}
                                >
                                    <Input style={{ height: 40 }} placeholder="กรอกเลข Order" />
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
                                  name="Tracking"
                                  rules={[{ required: true, message: "กรอกเลข Tracking" }]}>
                                  <Input style={{ height: 40 }} placeholder="กรอกเลข Tracking" />
                                </Form.Item>
                            </Col>
                            <Col span={8}>
                                <Form.Item
                                    id="Logistic"
                                    label={
                                        <span style={{ color: "#657589" }}>
                                            กรอก Logistic:&nbsp;
                                            <Tooltip title="ผู้ให้บริการขนส่ง">
                                            <QuestionCircleOutlined
                                                style={{ color: "#657589" }}
                                            />
                                            </Tooltip>
                                        </span>
                                    }
                                    name="Logistic"
                                    rules={[{ required: true, message: "กรอก Logistic" },]}
                                >
                                    <Input style={{ height: 40 }} />
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
                id="SKU" 
                label={<span style={{ color: '#657589' }}>SKU:</span>}
                name="SKU"
                rules={[{ required: true, message: "กรุณากรอก SKU" }]}
            >
                <Select
                      showSearch
                      style={{ width: "100%", height: "40px" }}
                      placeholder="Search by SKU"
                      value={selectedSKU} // ใช้ค่าที่เลือก
                      onSearch={handleSearchSKU} // ใช้สำหรับค้นหา SKU
                      onChange={handleSKUChange} // เมื่อเลือก SKU
                      loading={loading}
                      listHeight={160} 
                      virtual 
                      dropdownStyle={{ minWidth: 200 }}
                >
                      {skuOptions.map((option) => (
                        <Option 
                          key={`${option.sku}-${option.size}`} 
                          value={option.sku}
                        >
                          {option.sku}
                      </Option>
                      ))}
                </Select>
            </Form.Item>
        </Col>
        <Col span={8}>
            <Form.Item
                id="SKU_Name" 
                label={<span style={{ color: '#657589' }}>SKU Name:</span>}
                name="SKU_Name"
                rules={[{ required: true, message: "กรุณาเลือก SKU Name" }]}
            >
                <Select
                      showSearch
                      style={{ width: "100%", height: "40px" }}
                      placeholder="Search by Product Name"
                      value={selectedName} // ใช้ค่าที่เลือก
                      onSearch={handleSearchNameAlias} // ใช้สำหรับค้นหา Name Alias
                      onChange={handleNameChange} // เมื่อเลือก Name Alias
                      loading={loading}
                      listHeight={160} // ปรับให้พอดีกับ 4 รายการ
                      virtual // ทำให้ค้นหาไวขึ้น
                      dropdownStyle={{ minWidth: 300 }}
                >
                      {nameOptions.map((option) => (
                        <Option 
                          key={`${option.nameAlias}-${option.size}`} 
                          value={`${option.nameAlias}+${option.size}`}
                        >
                          {option.nameAlias}
                        </Option>
                      ))}
                </Select>
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
            dataSource={dataSource}
            rowKey="SKU"
            pagination={false}
        />
    )}
                        
    </Form>
            <Row align="middle" justify="center" style={{ marginTop: "20px", width: '100%' }}> 
                {/* <Button
                    id="Confirm" 
                    type="primary"
                    onClick={handleSubmit}
                    // disabled={
                    //     (value === 2 && formValid === true ) 
                    //     ? false 
                    //     : (value === 1 && dataSource.length > 0&& formValid === true ) 
                    //     ? false 
                    //     :true
                    // }
                    // disabled={!((value === 2 && formValid) || (value === 1 && dataSource.length > 0 && formValid))}
                    // disabled={tableData.length === 0 || !formValid || value === undefined}  // ปิดการใช้งานเมื่อไม่มีข้อมูลในตาราง, form ไม่ valid หรือไม่มีการเลือก value
                >
                    ยืนยัน
                </Button> */}
                <Popconfirm
                    id="popconfirmSubmit"
                    title="คุณแน่ใจหรือไม่ว่าต้องการส่งข้อมูล?"
                    onConfirm={handleSubmit} 
                    okText="ใช่"
                    cancelText="ไม่"
                    >
                    <Button
                        id="Confirm"
                        type="primary"
                    >
                        Submit
                    </Button>
                </Popconfirm>
            </Row>
            </Layout.Content>
        </Layout>
    </ConfigProvider>
    );
};

export default CreateBlind;
