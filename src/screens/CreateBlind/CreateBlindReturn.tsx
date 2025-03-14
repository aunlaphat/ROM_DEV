import { Modal, Popconfirm, Button, Radio, Select, Space, Col, ConfigProvider, Form, Layout, Row, Input, InputNumber, Table, notification, message, Tooltip } from "antd";
import { DeleteOutlined, PlusCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import type { RadioChangeEvent } from 'antd';
import { debounce } from "lodash";
import React, { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";
import { DataItemBlind, Product, TRANSPORT_TYPES } from '../../types/types';
import { useSelector } from 'react-redux';
import { RootState } from "../../redux/types";
import api from "../../utils/axios/axiosInstance"; 
import {FETCHSKU, SEARCHPRODUCT, CREATETRADE, SEARCHORDERTRACK} from '../../services/path';
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

    // ดึงข้อมูลผู้ใช้ที่เข้าสู่ระบบ
      const auth = useSelector((state: RootState) => state.auth);
      const userID = auth?.user?.userID;
      const token = localStorage.getItem("access_token");
    
      const [isModalVisible, setIsModalVisible] = useState(false);
    
      const [currentPage, setCurrentPage] = useState<number>(1);
      const [pageSize, setPageSize] = useState<number>(5);
    
      // ฟังก์ชันสำหรับเปลี่ยนหน้า
      const handlePageChange = (page: number, pageSize: number) => {
        setCurrentPage(page);
        setPageSize(pageSize); // ถ้าผู้ใช้เลือกจำนวนรายการต่อหน้าใหม่
      };
    
      // คำนวณจำนวนหน้าทั้งหมดจากจำนวนรายการทั้งหมด
      const totalPages = Math.ceil(dataSource.length / pageSize);
    
      // ตรวจสอบว่า pagination ควรแสดงหรือไม่ (ให้แสดงเสมอแม้ว่า dataSource จะมีน้อยกว่า pageSize)
      const showPagination = dataSource.length > 0;
    
      const showModal = () => {
          setIsModalVisible(true);
      };
    
      const handleOk = () => {
        setIsModalVisible(false);
        handleSubmit(); 
    };
  
    const handleCancel = () => {
        setIsModalVisible(false);
    };
    
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

      /*** Logistic Type ***/
      const [isOtherTransport, setIsOtherTransport] = useState(false);
      const [transportValue, setTransportValue] = useState<string | undefined>(undefined);
      const handleTransportChange = (value: string) => {
        if (value === 'OTHER') {
          setIsOtherTransport(true);
          form.resetFields(['Logistic']);
          setTransportValue(''); 
        } else {
          setIsOtherTransport(false);
          setTransportValue(value);
        }
      };

    /*** SKU&NameAlias ***/
    const debouncedSearchSKU = debounce(async (value: string, searchType: string) => {
        setLoading(true);
        try {
        const response = await api.get(SEARCHPRODUCT, {
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
        const response = await api.get(FETCHSKU, {
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
        { title: "รหัสสินค้า", dataIndex: "SKU", id: "SKU" },
        { title: "ชื่อสินค้า", dataIndex: "Name", id: "Name" },
        { title: "จำนวนที่ได้รับ", dataIndex: "QTY", id: "QTY" },
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

    const handleSubmit = async () => {
        form.validateFields()
            .then(async (values) => {
                // Check if 'แกะกล่อง' is selected as "Yes"
                if (value === 1) {
                    // If "Yes", check if there's at least one item in dataSource
                    if (dataSource.length === 0) {
                        notification.warning({
                            message: "กรุณาเพิ่มข้อมูลในตาราง",
                            description: "กรุณากรอกรายการสินค้าอย่างน้อย 1 รายการก่อนส่งข้อมูล !",
                        });
                        return;
                    }
                }
    
                // Proceed with submitting data
                const requestData = {
                    OrderNo: values.OrderNo,
                    ChannelID: 3, 
                    TrackingNo: values.Tracking,
                    Logistic: values.Logistic,
                    ReturnDate: values.Date,
                    StatusReturnID: 3, 
                    BeforeReturnOrderLines: dataSource.map(item => ({
                        SKU: item.SKU,
                        ItemName: item.Name,
                        QTY: item.QTY,
                        TrackingNo: values.Tracking,
                        CreateBy: userID, 
                    })),
                    CreateBy: userID, 
                };
                console.log("Order Data:", requestData);  // Log the data being sent
                try {
                    const response = await api.post(CREATETRADE, requestData, {
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                    });
    
                    if (response.status === 200) {
                        notification.success({
                            message: 'ส่งข้อมูลสำเร็จ',
                            description: 'ข้อมูลของคุณถูกส่งเรียบร้อยแล้ว!',
                        });
    
                        // After successful submission, navigate to the sale-return page
                        const orderResponse = await api.get(SEARCHORDERTRACK(values.OrderNo));
                        if (orderResponse.data && orderResponse.data.data && orderResponse.data.data.length > 0) {
                            navigate("/sale-return", { state: { orderNo: values.OrderNo } });
                        }
                    }
                } catch (error) {
                    notification.error({
                        message: 'เกิดข้อผิดพลาด',
                        description: 'ไม่สามารถส่งข้อมูลได้ กรุณาลองใหม่อีกครั้ง',
                    });
                }
            })
            .catch((info) => {
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
                สร้างรายการคืนจากหน้าคลังสินค้า
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
                                          label={<span style={{ color: "#657589" }}>ประเภทขนส่ง</span>}
                                          name="Logistic"
                                          rules={[{ required: true, message: "กรุณาเลือกขนส่ง" }]}
                                        >
                                          {isOtherTransport ? (
                                            <Input
                                              placeholder="กรอกประเภทขนส่ง"
                                              value={transportValue}
                                              onChange={(e) => setTransportValue(e.target.value)}
                                              style={{ height: 40, width: "100%" }}
                                            />
                                          ) : (
                                            <Select
                                              options={TRANSPORT_TYPES}
                                              placeholder="เลือกประเภทขนส่ง"
                                              showSearch
                                              optionFilterProp="label"
                                              style={{ height: 40, width: "100%" }}
                                              onChange={handleTransportChange}
                                              value={transportValue}
                                              listHeight={160} 
                                              virtual 
                                            />
                                          )}
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
                label={<span style={{ color: '#657589' }}>รหัสสินค้า</span>}
                name="SKU"
                // rules={[{ required: true, message: "กรุณาเลือก/ค้นหา รหัสสินค้า (SKU)" }]}
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
                label={<span style={{ color: '#657589' }}>ชื่อสินค้า</span>}
                name="SKU_Name"
                // rules={[{ required: true, message: "กรุณาเลือก/ค้นหา ชื่อสินค้า" }]}
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
                label={<span style={{ color: '#657589' }}>จำนวนที่ได้รับ</span>}
                name="QTY"
                // rules={[{ required: true, message: 'กรุณากรอกจำนวนที่ได้รับ' }]}
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
         <div >
         <Table
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
           dataSource={dataSource.slice((currentPage - 1) * pageSize, currentPage * pageSize)} // แสดงเฉพาะจำนวนรายการที่เลือก
           columns={columns}
           rowKey="key"
           pagination={false} // ปิด pagination ใน Table
           style={{
             width: "100%",
             tableLayout: "auto",
             border: "1px solid #ddd",
             borderRadius: "8px",
           }}
           scroll={{ x: "max-content" }}
           bordered={false}
         />
       
       {showPagination && (
         <div>
           {/* showTotal แสดงอยู่เหนือ showPagination */}
           <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginTop: 20 }}>
               <span style={{ fontSize: '14px', fontWeight: 'bold', color: '#555' }}>
                   ทั้งหมด <span style={{ color: '#007bff' }}>{dataSource.length}</span> รายการ
               </span>
           </div>
           <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginTop: 20, gap: 10 }}>
             {/* ปุ่มไปหน้าแรก */}
             <button
                 onClick={() => handlePageChange(1, pageSize)}
                 disabled={currentPage === 1}
                 style={{
                     fontSize: "14px",
                     // fontWeight: "bold",
                     padding: "4px 10px",
                     border: "1px solid #ddd",
                     borderRadius: "6px",
                     background: currentPage === 1 ? "#f5f5f5" : "#fff",
                     cursor: currentPage === 1 ? "not-allowed" : "pointer",
                 }}
             >
                 {"<<"}
             </button>

             {/* ปุ่มไปหน้าก่อน */}
             <button
                 onClick={() => handlePageChange(currentPage - 1, pageSize)}
                 disabled={currentPage === 1}
                 style={{
                     fontSize: "14px",
                     // fontWeight: "bold",
                     padding: "4px 10px",
                     border: "1px solid #ddd",
                     borderRadius: "6px",
                     background: currentPage === 1 ? "#f5f5f5" : "#fff",
                     cursor: currentPage === 1 ? "not-allowed" : "pointer",
                 }}
             >
                 {"<"}
             </button>

             {/* แสดงเลขหน้าแบบ [ 1 / 9 ] */}
             <span style={{ fontSize: "14px", fontWeight: 'bold' }}>
                 [ {currentPage} to {Math.ceil(dataSource.length / pageSize)} ]
             </span>

             {/* ปุ่มไปหน้าถัดไป */}
             <button
                 onClick={() => handlePageChange(currentPage + 1, pageSize)}
                 disabled={currentPage === Math.ceil(dataSource.length / pageSize)}
                 style={{
                     fontSize: "14px",
                     // fontWeight: "bold",
                     padding: "4px 10px",
                     border: "1px solid #ddd",
                     borderRadius: "6px",
                     background: currentPage === Math.ceil(dataSource.length / pageSize) ? "#f5f5f5" : "#fff",
                     cursor: currentPage === Math.ceil(dataSource.length / pageSize) ? "not-allowed" : "pointer",
                 }}
             >
                 {">"}
             </button>

             {/* ปุ่มไปหน้าสุดท้าย */}
             <button
                 onClick={() => handlePageChange(Math.ceil(dataSource.length / pageSize), pageSize)}
                 disabled={currentPage === Math.ceil(dataSource.length / pageSize)}
                 style={{
                     fontSize: "14px",
                     // fontWeight: "bold",
                     padding: "4px 10px",
                     border: "1px solid #ddd",
                     borderRadius: "6px",
                     background: currentPage === Math.ceil(dataSource.length / pageSize) ? "#f5f5f5" : "#fff",
                     cursor: currentPage === Math.ceil(dataSource.length / pageSize) ? "not-allowed" : "pointer",
                 }}
             >
                 {">>"}
             </button>

             {/* เลือกจำนวนรายการต่อหน้า */}
             <select
                 value={pageSize}
                 onChange={(e) => handlePageChange(1, Number(e.target.value))}
                 className="paginate"
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
         )}
       </div>
    )}
                        
    </Form>
    <Row justify="center" gutter={16}>
              <Button
                id="Submit"
                onClick={showModal}
                className="submit-trade"
              >
                Submit
              </Button>
              <Modal
                title="คุณแน่ใจหรือไม่ว่าต้องการส่งข้อมูล?"
                open={isModalVisible}
                onOk={handleOk}
                onCancel={handleCancel}
                okText="ใช่"
                cancelText="ไม่"
                centered
                style={{ textAlign: 'center'}}
                footer={
                  <div style={{ textAlign: "center" }}> {/* ทำให้ปุ่มอยู่ตรงกลาง */}
                    <Button key="ok" type="default" onClick={handleOk} style={{ marginRight: 8 }} className="button-yes">
                      Yes
                    </Button>
                    <Button key="cancel" type="dashed" onClick={handleCancel} className="button-no">
                      No
                    </Button>
                  </div>
                }

              >
              </Modal>
          </Row>
            </Layout.Content>
        </Layout>
    </ConfigProvider>
    );
};

export default CreateBlind;
