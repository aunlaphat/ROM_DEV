import React, { useState, useEffect } from 'react';
import { notification, Form, Input, InputNumber, DatePicker, Button, Row, Col, Table, ConfigProvider, Layout, Select, Modal, message, Popconfirm, Divider, Tooltip } from 'antd';
import moment from 'moment';
import { DeleteOutlined, LeftOutlined, PlusCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { Name } from 'ajv';
import { debounce } from "lodash";
import api from "../../utils/axios/axiosInstance"; 
const { Option } = Select;

interface FormValues {
    SKU: string;
    QTY: number;
    SKU_Name: string;
}

// interface DataItem {
//     key: number;
//     SKU: string; 
//     Name: string;
//     QTY: number;
//   }
  
interface Product {
    Key: string;
    sku: string;
    nameAlias: string;
    size: string;
}

interface Warehouse {
    Key: string;
    WarehouseID: number;
    WarehouseName: string;
    Location: string;
}

// interface User {
//     ID: number;
//     Name: string;
//     role: 'Warehouse' | 'Accounting';
// }

interface DataSourceItem extends FormValues {
    key: number;
    warehouse_form?: string;
    location_form?: string;
    warehouse_to?: string;
}
// const SKUName = [
//     { Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", SKU: "G097171-ARM01-BL" },
//     { Name: "Bewell Sport armband size M For", SKU: "G097171-ARM01-GY" },
//     { Name: "Sport armband size L", SKU: "G097171-ARM02-BL" },
//     { Name: "Bewell Sport armband size M with light", SKU: "G097171-ARM03-GR" },
// ];

// const MockUser: User[] = [
//     { ID: 1, Name: "User 1", role: "Warehouse" },
//     { ID: 2, Name: "User 2", role: "Accounting" },

// ];

// สร้าง options สำหรับ SKU
// const skuOptions = SKUName.map(item => ({
//     value: item.SKU,  // SKU เป็นค่า value
//     label: item.SKU   // SKU เป็น label เพื่อแสดงใน dropdown
// }));

// // สร้าง options สำหรับ SKU Name
// const nameOptions = SKUName.map(item => ({
//     value: item.Name, // Name เป็นค่า value
//     label: item.Name  // Name เป็น label เพื่อแสดงใน dropdown
// }));

const IJPage: React.FC = () => {
    const [selectedSKU, setSelectedSKU] = useState<string | undefined>(undefined);
    const [selectedName, setSelectedName] = useState<string | undefined>(undefined);
    const [skuOptions, setSkuOptions] = useState<Product[]>([]); // To store SKU options
    const [nameOptions, setNameOptions] = useState<Product[]>([]); // To store Name Alias options
    const [qty, setQty] = useState<number | null>(null);  // Allow null

    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);

    const [loading, setLoading] = useState(false);

    const [form] = Form.useForm();
    const [dataSource, setDataSource] = useState<DataSourceItem[]>([]);
    const [formValid, setFormValid] = useState(false);
    const [isSubmitted, setIsSubmitted] = useState(false);
    const navigate = useNavigate();
    const [formDisabled, setFormDisabled] = useState(false);
    const [ij, setIJ] = useState<string>('');
    const [remark, setRemark] = useState<string>('');
    const [submittedRemark, setSubmittedRemark] = useState<string>('');

    const handleIJChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setIJ(e.target.value);
    };

    const handleRemarkChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setRemark(e.target.value);
    };

    const handleError = (error: any) => {
        notification.warning({
            message: "กรุณากรอกข้อมูลให้ครบ",
    
        });
    };

    const onChange = () => {
        const values = form.getFieldsValue();
        const { Date, SKU, QTY } = values;

        // Set form validity based on required fields
        setFormValid(Date && SKU && QTY);
    };
    // ค้นหา Product (SKU หรือ NAMEALIAS)
    const debouncedSearchSKU = debounce(async (value: string, searchType: string) => {
        setLoading(true);
        try {
        const response = await api.get("/api/constants/search-product", {
            params: {
            keyword: value,
            searchType,
            offset: 0,
            limit: 5,
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

    // เมื่อเลือก Name Alias แล้วใช้ `/api/constants/get-sku` เพื่อหา SKU
    const handleNameChange = async (value: string) => {
        const [nameAlias, size] = value.split("+"); // แยกค่า nameAlias และ size โดยใช้ `+`

        try {
        setLoading(true);
        const response = await api.get("/api/constants/get-sku", {
            params: { nameAlias, size },
        });

        // เก็บผลลัพธ์จาก API เพื่อแสดงหลาย SKU
        const products = response.data.data;

        if (products.length > 0) {
            setSkuOptions(products.map((product: Product) => ({
            sku: product.sku,
            nameAlias: product.nameAlias,
            size: product.size,
            })));
            form.setFieldsValue({
            SKU: products[0].sku, // ตั้งค่า SKU ตัวแรกที่พบ
            });
        } else {
            console.warn("No SKU found for:", nameAlias, size);
            setSkuOptions([]); 
            setNameOptions([]); 
            form.setFieldsValue({ SKU: "", SKU_Name: "" }); // เคลียร์ค่าในฟอร์ม
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
        form.setFieldsValue({
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
        } else { // เคลียร์ค่าเมื่อไม่มี SKU ที่ตรงกัน
        setSkuOptions([]); 
        setNameOptions([]); 
        setSelectedSKU("");
        setSelectedName("");
        }
    };

    useEffect(() => {
        const fetchWarehouses = async () => {
            try {
                const response = await api.get('/api/constants/get-warehouse');
                setWarehouses(response.data.data);
            } catch (error) {
                console.error('Failed to fetch warehouses:', error);
                notification.error({
                    message: 'Error',
                    description: 'Failed to fetch warehouses.',
                });
            }
        };

        fetchWarehouses();
    }, []);

    const handleAdd = () => {
        // ตรวจสอบการกรอกข้อมูลที่จำเป็น เช่น วันที่คืน, ประเภทการขนส่ง, SKU, ชื่อสินค้า, และ QTY
        form.validateFields(['Date', 'Logistic', 'SKU', 'SKU_Name', 'QTY'])
            .then((values) => {
                // ถ้าข้อมูลในฟิลด์เหล่านี้ไม่ครบ จะมีข้อความเตือนขึ้น
                if (!values.Date || !values.Logistic || !values.SKU || !values.SKU_Name || !values.QTY) {
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
            id:'SKU'
        },
        {
            title: 'Name',
            dataIndex: 'SKU_Name',
            id:'SKU_Name'
        },
        {
            title: 'QTY',
            dataIndex: 'QTY',
            id:'SKU_Name'
        },
        {
            title: "Warehouse Form",
            id:'warehouse_form',
            dataIndex: "warehouse_form",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '100%' }}
                    onChange={(value) => handleChange(value, record.key, "warehouse_form")}
                    options={warehouses.map(warehouse => ({
                        value: warehouse.WarehouseID,
                        label: warehouse.WarehouseName,
                    }))}
                    dropdownStyle={{ minWidth: 120 }}
                    dropdownMatchSelectWidth={false}
                    maxTagTextLength={50} // กำหนดความยาวสูงสุดของข้อความในตัวเลือก
                    disabled={formDisabled}
                />

            ),
        },
        {
            title: "Location Form",
            id:'location_form',
            dataIndex: "location_form",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '100%' }}
                    onChange={(value) => handleChange(value, record.key, "location_form")}
                    options={warehouses.map(warehouse => ({
                        value: warehouse.Location,
                        label: warehouse.Location,
                    }))}
                    dropdownStyle={{ minWidth: 120 }}
                    dropdownMatchSelectWidth={false}
                    maxTagTextLength={50} // กำหนดความยาวสูงสุดของข้อความในตัวเลือก
                    disabled={formDisabled}
                />
            ),
        },
        {
            title: "Warehouse to",
            id:'warehouse_to',
            dataIndex: "warehouse_to",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '100%' }}
                    onChange={(value) => handleChange(value, record.key, "warehouse_to")}
                    options={warehouses.map(warehouse => ({
                        value: warehouse.WarehouseID,
                        label: warehouse.WarehouseName,
                    }))}
                    dropdownStyle={{ minWidth: 100 }}
                    dropdownMatchSelectWidth={false}
                    maxTagTextLength={50} // กำหนดความยาวสูงสุดของข้อความในตัวเลือก
                    disabled={formDisabled}
                />
            ),
        },
        {
            title: "Action",
            id:'Action',
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

    // const handleSelectChange = (value: string) => {
    //     const selectedOption = SKUName.find((val) => val.SKU === value);

    //     if (selectedOption) {
    //         form.setFieldsValue({
    //             SKU: selectedOption.SKU,
    //             SKU_Name: selectedOption.Name,
    //         });
    //         setSelectedValue(value);
    //     }
    // };
    // const handleSKUChange = (value: string) => {
    //     const selectedOption = SKUName.find((val) => val.SKU === value);
    //     if (selectedOption) {
    //         form.setFieldsValue({
    //             SKU: selectedOption.SKU,
    //             SKU_Name: selectedOption.Name,
    //         });
    //         setSelectedSKU(value);
    //         setSelectedName(selectedOption.Name); // อัปเดต selectedName
    //     }
    // };

    // const handleNameChange = (value: string) => {
    //     const selectedOption = SKUName.find((val) => val.Name === value);
    //     if (selectedOption) {
    //         form.setFieldsValue({
    //             SKU: selectedOption.SKU,
    //             SKU_Name: selectedOption.Name,
    //         });
    //         setSelectedName(value);
    //         setSelectedSKU(selectedOption.SKU); // อัปเดต selectedSKU
    //     }
    // };



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
                                        rules={[
                                        { required: true, message: "กรอก Logistic" },
                                        ]}
                                    >
                                        <Input style={{ height: 40 }} />
                                    </Form.Item>
                                    </Col>
                                </Row>
                                <Divider style={{ color: '#657589', fontSize: '22px', marginTop: 30, marginBottom: 30 }} orientation="left"> SKU Information </Divider>
                                <Row gutter={16} >
                                    <Col span={8}>
                                        <Form.Item
                                        id="SKU"
                                            label={<span style={{ color: '#657589' }}>กรอก SKU:</span>}
                                            name="SKU"
                                            rules={[{ required: true, message: "กรุณากรอก SKU" }]}
                                        >
                                            <Select
                                                showSearch
                                                style={{ width: '100%', height: '40px' }}
                                                dropdownStyle={{ minWidth: 200 }}
                                                listHeight={160}
                                                placeholder="Search by SKU"
                                                // optionFilterProp="label"
                                                value={selectedSKU} // แสดง SKU ที่ถูกเลือก
                                                onSearch={handleSearchSKU}
                                                onChange={handleSKUChange}
                                                loading={loading}
                                                virtual
                                                // options={skuOptions} // แสดง SKU ใน dropdown
                                                // disabled={formDisabled}
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
                                            label={<span style={{ color: '#657589' }}>Name:</span>}
                                            name="SKU_Name"
                                            rules={[{ required: true, message: "กรุณาเลือก SKU Name" }]}
                                        >
                                            <Select
                                                showSearch
                                                style={{ width: '100%', height: '40px' }}
                                                dropdownStyle={{ minWidth: 300 }}
                                                listHeight={160}
                                                placeholder="Search by Product Name"
                                                // optionFilterProp="label"
                                                value={selectedName} // แสดง SKU Name ที่ถูกเลือก
                                                onSearch={handleSearchNameAlias}
                                                onChange={handleNameChange}
                                                loading={loading}
                                                virtual // ทำให้ค้นหาไวขึ้น
                                                // options={nameOptions} // แสดง SKU Name ใน dropdown
                                                // disabled={formDisabled}
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
                                             name="QTY"
                                             label={<span style={{ color: '#657589' }}>QTY:</span>}
                                             rules={[{ required: true, message: "กรุณากรอก QTY" }]}
                                        >
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


