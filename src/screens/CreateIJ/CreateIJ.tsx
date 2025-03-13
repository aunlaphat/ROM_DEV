import React, { useState, useEffect } from 'react';
import { Upload, notification, Form, Input, InputNumber, DatePicker, Button, Row, Col, Table, ConfigProvider, Layout, Select, Modal, message, Popconfirm, Divider, Tooltip } from 'antd';
import moment from 'moment';
import { DeleteOutlined, LeftOutlined, PlusCircleOutlined, QuestionCircleOutlined, UploadOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import * as XLSX from "xlsx";
import { Name } from 'ajv';
import { debounce } from "lodash";
import api from "../../utils/axios/axiosInstance"; 
import icon from "../../assets/images/document-text.png";
import { useSelector } from 'react-redux';
import { RootState } from "../../redux/types";
import {FETCHSKU, SEARCHPRODUCT, FETCHWAREHOUSE, CREATETRADE} from '../../services/path';
import { TRANSPORT_TYPES, FormValues, Product, Warehouse, DataSourceItem } from '../../types/types';
const { Option } = Select;

const IJPage: React.FC = () => {
    const [selectedSKU, setSelectedSKU] = useState<string | undefined>(undefined);
    const [selectedName, setSelectedName] = useState<string | undefined>(undefined);
    const [skuOptions, setSkuOptions] = useState<Product[]>([]); // To store SKU options
    const [nameOptions, setNameOptions] = useState<Product[]>([]); // To store Name Alias options
    const [qty, setQty] = useState<number | null>(null);  
    const [returnQty, setReturnQty] = useState<number | null>(null);
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

    const handleOk = async () => {
        setIsModalVisible(false);
        await handleCreateIJ(); 
    };

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
        const response = await api.get(SEARCHPRODUCT, {
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
        const response = await api.get(FETCHSKU, {
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

    useEffect(() => {
        const fetchWarehouses = async () => {
            try {
                const response = await api.get(FETCHWAREHOUSE);
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
        form.validateFields(['Date', 'Logistic', 'SKU', 'SKU_Name', 'QTY', 'ReturnQTY'])
            .then((values) => {
                // ถ้าข้อมูลในฟิลด์เหล่านี้ไม่ครบ จะมีข้อความเตือนขึ้น
                if (!values.Date || !values.Logistic || !values.SKU || !values.SKU_Name || !values.QTY || !values.ReturnQTY) {
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
                form.resetFields(['SKU', 'SKU_Name', 'QTY', 'ReturnQTY']);
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

    const handleChange = (value: string, key: React.Key, dataIndex: string) => {
        console.log("handleChange value:", value); // เพิ่ม console.log ตรงนี้
        setDataSource(prevDataSource =>
            prevDataSource.map(item =>
                item.key === key ? { ...item, [dataIndex]: value } : item
            )
        );
    };

    const columns = [
        { title: 'รหัสสินค้า', dataIndex: 'SKU', id:'SKU', key: 'SKU', },
        { title: 'ชื่อสินค้า', dataIndex: 'SKU_Name', id:'SKU_Name', key: 'SKU_Name', },
        { title: 'จำนวนเริ่มต้น', dataIndex: 'QTY',id:'QTY', key: 'QTY',},
        { title: 'จำนวนที่คืน', dataIndex: 'QTY', id:'QTY', key: 'QTY',},
        { title: "คลังต้นทาง", id:'warehouse_form', dataIndex: "warehouse_form", key: "warehouse_form",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '100%' }}
                    onChange={(value) => handleChange(value, record.key, "warehouse_form")}
                    options={warehouses.map(warehouse => ({
                        value: warehouse.WarehouseID,
                        label: warehouse.WarehouseName,
                    }))}
                    dropdownStyle={{ minWidth: 120 }}
                    popupMatchSelectWidth={false}
                    maxTagTextLength={50} // กำหนดความยาวสูงสุดของข้อความในตัวเลือก
                    disabled={formDisabled}
                />

            ),
        },
        { title: "ที่อยู่ต้นทาง", id:'location_form', dataIndex: "location_form", key: "location_form",
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
        { title: "คลังปลายทาง", id:'warehouse_to', dataIndex: "warehouse_to", key: "warehouse_to",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '100%' }}
                    onChange={(value) => handleChange(value, record.key, "warehouse_to")}
                    options={warehouses.map(warehouse => ({
                        value: warehouse.WarehouseID,
                        label: warehouse.WarehouseName,
                    }))}
                    dropdownStyle={{ minWidth: 100 }}
                    popupMatchSelectWidth={false}
                    maxTagTextLength={50} // กำหนดความยาวสูงสุดของข้อความในตัวเลือก
                    disabled={formDisabled}
                />
            ),
        },
        { title: "Action", id:'Action', dataIndex: "Action", key: "Action",
            render: (_: any, record: DataSourceItem) => (
                <DeleteOutlined
                    style={{ cursor: 'pointer', color: 'red', fontSize: '20px' }}
                    onClick={() => handleDelete(record.key)}
                />
            ),
        },
    ];

      const handleDownloadTemplate = () => {
        const templateColumns = [
            { title: "ลำดับ", dataIndex: "key", key: "key" }, // เพิ่ม Column "ลำดับ"
            ...columns.filter((col) => col.key), // เพิ่ม Column อื่นๆ
        ];
        const ws = XLSX.utils.json_to_sheet([]);
        XLSX.utils.sheet_add_aoa(ws, [templateColumns.map((col) => col.title)]);
    
        const mappedDataSource = dataSource.map((item) => {
            return {
                key: item.key, // เพิ่ม key
                SKU: item.SKU,
                Name: item.SKU_Name,
                QTY: item.QTY,
                ReturnQTY: item.ReturnQTY,
                // PricePerUnit: item.PricePerUnit,
                // Price: item.Price,
            };
        });
    
        XLSX.utils.sheet_add_json(ws, mappedDataSource, { origin: "A2", skipHeader: true });
    
        const wb = XLSX.utils.book_new();
        XLSX.utils.book_append_sheet(wb, ws, "Template");
        XLSX.writeFile(wb, "Template.xlsx");
    };
    
    const handleUpload = (file: File) => {
      const reader = new FileReader();
      reader.onload = (e) => {
          const data = new Uint8Array(e.target?.result as ArrayBuffer);
          const workbook = XLSX.read(data, { type: "array" });
          const worksheet = workbook.Sheets[workbook.SheetNames[0]];
    
          const json = XLSX.utils.sheet_to_json<any>(worksheet);
    
          console.log("JSON Data:", json);
    
          if (Array.isArray(json) && json.length > 0) {
              const mappedData: DataSourceItem[] = json.map((row, index) => {
                  const dataItem: DataSourceItem = {
                      key: index + 1,
                      SKU: row["รหัสสินค้า"] as string,
                      SKU_Name: row["ชื่อสินค้า"] as string,
                      QTY: row["จำนวนเริ่มต้น"] as number,
                      ReturnQTY: row["จำนวนที่คืน"] as number,
                    //   PricePerUnit: row["ราคาต่อหน่วย"] as number,
                    //   Price: row["ราคารวม"] as number,
                  };
                  return dataItem;
              }).filter((item) => item.SKU && item.QTY);
    
              setDataSource(mappedData);
    
              notification.success({
                  message: "อัปโหลดสำเร็จ",
                  description: "ข้อมูลจากไฟล์ Excel ถูกนำเข้าเรียบร้อยแล้ว!",
              });
          } else {
              notification.error({
                  message: "อัปโหลดล้มเหลว",
                  description: "ไฟล์ Excel ไม่มีข้อมูลที่ถูกต้อง!",
              });
          }
      };
      reader.readAsArrayBuffer(file);
    };
    
      const uploadProps = {
        beforeUpload: (file: File) => {
          handleUpload(file);
          return false; // ป้องกันไม่ให้ Ant Design ทำการอัปโหลด
        },
      };

    const generateRandomNumber = () => {
        return Math.floor(Math.random() * 10000);
    };

    const handleCreateIJ = async () => {
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
            } else {
                const randomNumber = generateRandomNumber(); // สร้างเลขสุ่ม
                form.setFieldsValue({ IJ_Create: randomNumber }); // ตั้งค่าเลขสุ่มให้กับฟิลด์ในฟอร์ม
                setDataSource((prevData) => prevData.map(item => ({ ...item, IJ_Create: randomNumber })));
    
                notification.success({
                    message: 'สำเร็จ',
                    description: `Create IJ สำเร็จ! เลขสุ่มที่สร้างคือ: ${randomNumber}`,
                });
                handleSubmitData(); // เรียกใช้ฟังก์ชัน handleSubmit
                setIsSubmitted(true);
            }

        } catch (error) {
            handleError(error); 
        }
        setFormDisabled(true);
    };


    const handleSubmitData = async () => {
        const combinedRemark = `${remark} - IJ: ${ij}`;
        console.log("IJ:", ij);
        console.log("Remark:", combinedRemark);
    
        try {
            // ตรวจสอบว่ามีข้อมูลในตารางอย่างน้อยหนึ่งรายการ
        if (dataSource.length === 0) {
            notification.warning({
              message: "ไม่สามารถส่งข้อมูลได้",
              description: "กรุณาเพิ่มข้อมูลในตารางก่อนส่ง!",
            });
            return; 
          }
      
          // ดึงค่าจากฟอร์ม
          const values = await form.validateFields();
          console.log("Form Values:", values);
          console.log("dataSource:", dataSource); // เพิ่ม console.log ตรงนี้
            const orderData = {
                // OrderNo: values.Order,
                SoNo: values.IJ,
                SrNo: String(values.IJ_Create),
                ChannelID: 4, 
                TrackingNo: values.TrackingNumber,
                Logistic: values.Logistic,
                ReturnDate: values.Date,
                StatusReturnID: 4, 
                CreateBy: userID,
                
                BeforeReturnOrderLines: dataSource.map(item => ({
                    SKU: item.SKU,
                    ItemName: item.SKU_Name,
                    QTY: item.QTY,
                    ReturnQTY: item.ReturnQTY, 
                    WarehouseID: item.warehouse_to,
                    TrackingNo: values.TrackingNumber
                })),
                Reason: values.Remark, // Adding remark with IJ information
            };
            console.log("Order Data:", orderData);
            const token = localStorage.getItem('access_token')
            const response = await api.post(CREATETRADE, orderData, {
                headers: {
                  Authorization: `Bearer ${token}`,
                },
            });
    
            if (response.status === 200) {
                notification.success({
                    message: 'ส่งข้อมูล สำเร็จ',
                    description: 'ข้อมูลทั้งหมดได้ถูกส่งเรียบร้อยแล้ว',
                });
    
                setDataSource([]);
                form.resetFields();
                setRemark('');
                setIJ('');
                setIsSubmitted(false);
                setFormDisabled(false); // เปิดใช้งานฟอร์มใหม่อีกครั้ง
    
            } else {
                notification.error({
                    message: 'เกิดข้อผิดพลาด',
                    description: 'ไม่สามารถส่งข้อมูลได้ กรุณาลองใหม่อีกครั้ง',
                });
            }
        } catch (error) {
            console.error("Error sending data:", error);
            notification.error({
                message: 'เกิดข้อผิดพลาด',
                description: 'ไม่สามารถส่งข้อมูลได้ กรุณาลองใหม่อีกครั้ง',
            });
        }
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
        setIsModalVisible(false);

    };

    return (
        <ConfigProvider>
            <div  id="titleContainer" style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                สร้างรายการคืนภายใน (IJ)
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
                                                    <Tooltip title="เลขพัสดุจากขนส่ง">
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
                                   <Divider
                                     style={{ color: "#657589", fontSize: "22px", margin: 30 }}
                                     orientation="left"
                                   > 
                                     {" "}
                                     SKU information
                                   </Divider>
                                   <Row gutter={16} style={{ marginTop: "10px", width: "100%", justifyContent: "center"}}>
                                     <Col span={7}>
                                       <Form.Item
                                         id="SKU"
                                         label={<span style={{ color: "#657589" }}>รหัสสินค้า</span>}
                                         name="SKU"
                                         // rules={[{ required: true, message: "กรุณากรอก SKU" }]}
                                       >
                                         <Select
                                           showSearch
                                           style={{ width: "100%", height: "40px" }}
                                           dropdownStyle={{ minWidth: 200 }}
                                           listHeight={160}
                                           placeholder="Search by SKU"
                                           value={selectedSKU} // ใช้ค่าที่เลือก
                                           onSearch={handleSearchSKU} // ใช้สำหรับค้นหา SKU
                                           onChange={handleSKUChange} // เมื่อเลือก SKU
                                           loading={loading}
                                           virtual
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
                       
                                     <Col span={7}>
                                       <Form.Item
                                         id="SKU_Name"
                                         label={
                                           <span style={{ color: "#657589" }}>ชื่อสินค้า</span>
                                         }
                                         name="SKU_Name"
                                         // rules={[{ required: true, message: "กรุณาเลือก SKU Name" }]}
                                       >
                                         <Select
                                           showSearch
                                           style={{ width: "100%", height: "40px" }}
                                           dropdownStyle={{ minWidth: 300 }}
                                           listHeight={160}
                                           placeholder="Search by Product Name"
                                           value={selectedName} // ใช้ค่าที่เลือก
                                           onSearch={handleSearchNameAlias} // ใช้สำหรับค้นหา Name Alias
                                           onChange={handleNameChange} // เมื่อเลือก Name Alias
                                           loading={loading}
                                           virtual 
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
                                   </Row>
                                   <Row gutter={16} style={{ marginTop: "10px", width: "100%", justifyContent: "center" }}>
                                     <Col span={4}>
                                       <Form.Item
                                         id="qty"
                                         label={<span style={{ color: "#657589" }}>จำนวนเริ่มต้น</span>}
                                         name="QTY"
                                         // rules={[{ required: true, message: "กรุณากรอก QTY" }]}
                                       >
                                         <InputNumber
                                           min={1}
                                           max={100}
                                           value={qty}
                                           onChange={(value) => setQty(value)}
                                           style={{
                                             width: "100%",
                                             height: "40px",
                                             lineHeight: "40px",
                                           }}
                                         />
                                       </Form.Item>
                                     </Col>
                                     <Col span={4}>
                                       <Form.Item
                                            id="returnQTY"
                                            label={<span style={{ color: "#657589" }}>จำนวนที่คืน</span>}
                                            name="ReturnQTY"
                                            // rules={[{ required: true, message: "กรุณากรอก QTY" }]}
                                       >
                                         <InputNumber
                                           min={1}
                                           max={100}
                                           value={returnQty}
                                           onChange={(value) => setReturnQty(value)}
                                           style={{
                                             width: "100%",
                                             height: "40px",
                                             lineHeight: "40px",
                                           }}
                                         />
                                       </Form.Item>
                                     </Col>
                                     {/* <Col span={4}>
                                       <Form.Item
                                         id="pricePerUnit"
                                         label={<span style={{ color: "#657589" }}>ราคาต่อหน่วย</span>}
                                         name="PricePerUnit"
                                         // rules={[{ required: true, message: "กรุณากรอก Price" }]}
                                       >
                                         <InputNumber
                                           min={1}
                                           max={100000}
                                           value={pricePerUnit}
                                           onChange={(value) => setPricePerUnit(value)}
                                           step={0.01}
                                           style={{
                                             width: "100%",
                                             height: "40px",
                                             lineHeight: "40px",
                                           }}
                                         />
                                       </Form.Item>
                                     </Col>
                                     <Col span={4}>
                                       <Form.Item
                                         id="price"
                                         label={<span style={{ color: "#657589" }}>ราคารวม</span>}
                                         name="Price"
                                         // rules={[{ required: true, message: "กรุณากรอก Price" }]}
                                       >
                                         <InputNumber
                                           min={1}
                                           max={100000}
                                           value={price}
                                           // onChange={(value) => setPrice(value)}
                                           step={0.01}
                                           disabled
                                           style={{
                                             width: "100%",
                                             height: "40px",
                                             lineHeight: "40px",
                                           }}
                                         />
                                       </Form.Item>
                                     </Col> */}
                                     <Col span={4}>
                                       <Button
                                         id="add"
                                         type="primary"
                                         style={{ width: "100%", height: "40px", marginTop: 30 }}
                                         onClick={handleAdd} 
                                       >
                                         <PlusCircleOutlined />
                                         Add
                                       </Button>
                                     </Col>
                                </Row>
                            </div>
                        </Form>

                        <Row gutter={16} style={{ marginBottom: 20 }}>
                            <Col>
                                <Button id=" Download Template" onClick={handleDownloadTemplate}>
                                <img
                                    src={icon}
                                    alt="Download Icon"
                                    style={{ width: 16, height: 16, marginRight: 8 }}
                                />
                                Download Template
                                </Button>
                            </Col>
                            <Col>
                                <Upload {...uploadProps} showUploadList={false}>
                                <Button
                                    id=" Import Excel"
                                    icon={<UploadOutlined />}
                                    style={{
                                    background: "#7161EF",
                                    color: "#FFF",
                                    marginBottom: 10,
                                    }}
                                >
                                    Import Excel
                                </Button>
                                </Upload>
                            </Col>
                        </Row>

                <div>
                        <Table
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
                  </div>
                </Layout.Content>
            </Layout>
        </ConfigProvider>
    );
};

export default IJPage;


