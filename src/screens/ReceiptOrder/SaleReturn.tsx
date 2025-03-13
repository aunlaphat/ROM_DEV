import React, { useEffect, useRef, useState } from 'react';
import { Steps, Col, ConfigProvider, Form, Layout, Row, Select, Button, Table, Modal, Input, notification, Divider, Popconfirm } from 'antd';
import Webcam from 'react-webcam';
import { CameraOutlined, RedoOutlined, DeleteOutlined, ScanOutlined, CheckCircleOutlined, WarningOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons';
import { QrReader, QrReaderProps } from 'react-qr-reader';
import api from "../../utils/axios/axiosInstance"; 
import { useSelector } from 'react-redux';
import { RootState } from "../../redux/types";
import { CustomQrReaderProps, ReceiptOrder, ReceiptOrderLine } from '../../types/types';

const SaleReturn: React.FC = () => {
    const [orderOptions, setOrderOptions] = useState<{ value: string; label: string }[]>([]);
    const [selectedOrderNo, setSelectedOrderNo] = useState<string | null>(null); 
    const [skuInput, setSkuInput] = useState<string>('');
    const [skuName, setSkuName] = useState<string | null>(null);
   
    const [scanResult, setScanResult] = useState<string | null>(null);
    const [showScanner, setShowScanner] = useState<boolean>(false);
    const [currentStep, setCurrentStep] = useState(0);
    const [currentRecordKey, setCurrentRecordKey] = useState<string | null>(null);
    const [showWebcam, setShowWebcam] = useState(false);
    const [showSteps, setShowSteps] = useState(false);
    const [showTable, setShowTable] = useState(false);
    const [current, setCurrent] = useState(0);
    const [images, setImages] = useState<{ [key: string]: string | null }>({
        step1: null,
        step2: null,
    });
    const webcamRef = useRef<Webcam>(null);
    const [data, setData] = useState<ReceiptOrderLine[]>([]);
    
    const onChange = (current: number) => {
        setCurrentStep(current);
    };
    const [cameraFacingMode, setCameraFacingMode] = useState<'user' | 'environment'>('environment');

    const toggleCamera = () => {
        setCameraFacingMode((prevMode) => (prevMode === 'user' ? 'environment' : 'user'));
    };

    // ดึงข้อมูลผู้ใช้ที่เข้าสู่ระบบ
    const auth = useSelector((state: RootState) => state.auth);
    const userID = auth?.user?.userID;
    const token = localStorage.getItem("access_token");

    useEffect(() => {
        const fetchOrderOptions = async () => {
            try {
                const response = await api.get('/api/import-order/get-order-tracking');
                const options = response.data.data.map((item: any, index: number) => ({
                    key: `${item.orderNo}-${index}`, // เพิ่ม index เพื่อให้คีย์เป็นเอกลักษณ์
                    value: item.orderNo,
                    label: item.orderNo,
                }));
                setOrderOptions(options);
            } catch (error) {
                console.error('Failed to fetch order options:', error);
                notification.error({
                    message: 'Error',
                    description: 'Failed to fetch order options.',
                });
            }
        };

        fetchOrderOptions();
    }, []);
   
    const handleSelectChange = async (value: string) => {
        setSelectedOrderNo(value); 
        setShowSteps(true);
        setCurrentStep(0);
        setShowTable(false);

        try {
            const response = await api.get(`/api/import-order/search-order-tracking?search=${value}`);
          // ตรวจสอบว่ามีข้อมูล orderLines หรือไม่
          if (response.data && response.data.data && response.data.data.length > 0 && response.data.data[0].orderLines) {
            // กรองข้อมูล orderLines ที่มี SKU ขึ้นต้นด้วย "G"
            const filteredOrderLines = response.data.data[0].orderLines.filter((item: any) => item.sku.startsWith('G'));
            const orderData = filteredOrderLines.map((item: any, index: number) => ({
                key: `${index + 1}`,
                sku: item.sku,
                itemName: item.itemName,
                qty: item.qty,
                receivedQty: 0,
                price: item.price,
                image: null,
            }));
            setData(orderData);
            // setShowTable(true);
        } else {
            setData([]);
            setShowTable(false);
            notification.warning({
                message: 'คำเตือน',
                description: 'ไม่พบข้อมูล Order Lines สำหรับ Order นี้',
            });
        }
        } catch (error) {
            console.error('Failed to fetch order data:', error);
            notification.error({
                message: 'Error',
                description: 'Failed to fetch order data.',
            });
        }
    };

    const capturePhoto = () => {
        if (webcamRef.current) {
            const src = webcamRef.current.getScreenshot();
            if (src) {
                setImages((prevImages) => ({
                    ...prevImages,
                    [`step${currentStep + 1}`]: src,
                }));
            }
        }
    };

    const retakePhoto = () => {
        setImages((prevImages) => ({
            ...prevImages,
            [`step${currentStep + 1}`]: null,
        }));
    };

    const handleNextStep = () => {
        if (currentStep < data.length + 2) {
            if (currentStep >= 2) { // ตรวจสอบว่าอยู่ในขั้นตอนที่ 3 (SKU) หรือไม่
                const recordKey = data[currentStep - 2].key; // คำนวณ recordKey จาก currentStep
                const imageSrc = images[`step${currentStep + 1}`];
                if (imageSrc) {
                    setData(prevData => {
                        return prevData.map(record => {
                            if (record.key === recordKey) {
                                return { ...record, image: imageSrc };
                            }
                            return record;
                        });
                    });
                }
            }
            setCurrentStep(prevStep => prevStep + 1);
        } else {
            setShowTable(true);
        }
    };
    
    const handleDelete = (recordKey: string) => {
        setData((prevData) => prevData.filter(item => item.key !== recordKey));
    };

    const columns = [
        { title: 'รหัสสินค้า', dataIndex: 'sku', key: 'sku' ,id:'sku'},
        { title: 'ชื่อสินค้า', dataIndex: 'itemName', key: 'itemName' ,id:'itemName'},
        { title: 'จำนวนสินค้า', dataIndex: 'qty', key: 'qty' ,id:'qty'},
        { title: 'จำนวนสินค้าที่คืน', dataIndex: 'receivedQty', key: 'receivedQty' ,id:'receivedQty'},
        { title: 'ราคาสินค้ารวม', dataIndex: 'price', key: 'price',id:'price' },
        {
            title: 'Return',
            id:'Return' ,
            dataIndex: 'Return',
            key: 'Return',
            render: (_: any, record: ReceiptOrderLine) => {
                const isConfirmed = record.receivedQty === record.qty;
                return isConfirmed ? (
                    <Button type="primary" style={{ background: '#D1FAD3', width: '73px', borderRadius: '20px' }}>
                        <CheckOutlined style={{ color: 'green', fontSize: '15px' }} />
                    </Button>
                ) : (
                    <Button type="primary" style={{ background: '#FDCACA', width: '73px', borderRadius: '20px' }}>
                        <CloseOutlined style={{ color: 'red', fontSize: '15px' }} />
                    </Button>
                );
            }
        },
        {
            title: 'Image',
            id:'Image',
            dataIndex: 'image',
            key: 'image',
            render: (_: any, record: ReceiptOrderLine) => {
                // ตรวจสอบว่า record.image มีค่าหรือไม่
                const stepImage = images[`step${parseInt(record.key, 10)+2}`]; 
                return stepImage ? (
                    <img src={stepImage} alt="Return" style={{ width: '100px' }} />
                ) : (
                    <Button
                        style={{ background: '#02C39A' }}
                        icon={<CameraOutlined />}
                        type="primary"
                        onClick={() => handleTakePhoto(record.key, record.sku)} 
                    >
                        กดเพื่อถ่ายรูป
                    </Button>
                );
    }
    },
    
    {
        title: 'Action',
        id:'Action',
        key: 'action',
        render: (_: any, record: ReceiptOrderLine) => (
          <>
            <Button
              style={{
                marginRight: '10px',
                marginBottom: '10px',
                background: '#BADEFF',
                color: '#1890FF',
              }}
              id="Takepicture"
              type="primary"
              onClick={() => handleRetakePhoto(record.key)}
              icon={<RedoOutlined />}
            >
              ถ่ายรูปใหม่
            </Button>
      
            <Popconfirm
             id="Delectpopconfirm"
              title="คุณแน่ใจหรือว่าต้องการลบรายการนี้?"
              onConfirm={() => handleDelete(record.key)}
              okText="ยืนยัน"
              cancelText="ยกเลิก"
            >
              <Button
                id="Cancel"
                type="primary"
                style={{ color: '#E53939', background: '#F9D3D3' }}
                icon={<DeleteOutlined />}
              >
                Delete
              </Button>
            </Popconfirm>
          </>
        ),
      },
    ];

    const handleTakePhoto = (recordKey: string, sku: string) => {
        setCurrentRecordKey(recordKey);
        setSkuName(sku); 
        console.log("CurrentRecordKey (take photo):", recordKey);
        setShowWebcam(true);
    };

    const handleBackStep = () => {
        if (currentStep > 0) {  // ตรวจสอบว่า currentStep มากกว่า 0 เพื่อย้อยกลับ
            setCurrentStep((prevStep) => prevStep - 1);
        }
    };

    const handleRetakePhoto = (recordKey: string) => {
        setCurrentRecordKey(recordKey);
        const currentItem = data.find(item => item.key === recordKey);
        setSkuName(currentItem?.sku || null); 
        console.log("CurrentRecordKey (retake photo):", recordKey);
        setShowWebcam(true);
    };

    const handleCapturePhoto = () => {
        if (webcamRef.current && currentRecordKey) {
            console.log("CurrentRecordKey:", currentRecordKey);
            const imageSrc = webcamRef.current.getScreenshot();
            if (imageSrc) {
                setData((prevData) =>
                    prevData.map((item) =>
                        item.key === currentRecordKey ? { ...item, image: imageSrc } : item
                    )
                );
                console.log("Data after capture:", data);
            }
            setShowWebcam(false);
        }
    };
    
    const handleCancelModal = () => {
        setShowWebcam(false);
        setCurrentRecordKey(null);
        setSkuName(null);
        setImages((prevImages) => ({ ...prevImages, [`step${currentStep + 1}`]: null })); 
    };
    
    const handleConfirmReceived = () => {
        const foundSKU = data.find(item => item.sku === skuInput.trim()); // ใช้ trim() เพื่อจัดการกับช่องว่าง
        if (foundSKU) {
            // เช็คว่าจำนวนที่รับเข้าจะไม่เกิน qty
            if (foundSKU.receivedQty < foundSKU.qty) {
                setData(prevData =>
                    prevData.map(item =>
                        item.sku === skuInput
                            ? { ...item, receivedQty: item.receivedQty + 1 }
                            : item
                    )
                );
                notification.success({
                    message: 'สำเร็จ',
                    description: `อัปเดตการรับเข้า SKU: ${skuInput} สำเร็จ`,
                    placement: 'topRight',
                });
            } else {
                notification.warning({
                    message: 'ข้อควรระวัง',
                    description: `ไม่สามารถรับเข้าได้ เพราะจำนวนที่รับเข้ามากกว่าหรือเท่ากับ QTY ของ SKU: ${skuInput}`,
                    placement: 'topRight',
                });
            }
        } else {
            notification.error({
                message: 'ข้อผิดพลาด',
                description: 'กรุณากรอก SKU ให้ถูกต้องหรือครบถ้วน',
                placement: 'topRight',
            });
        }
        setSkuInput(''); 
    };
    notification.config({
        placement: 'topRight',
        duration: 3, // ระยะเวลาการแสดง notification
        maxCount: 3, // จำนวน notification สูงสุดที่จะแสดงพร้อมกัน
    });

    const handleScanSku = () => {
        setShowScanner(true); // เปิด/ปิดการใช้งานกล้องเพื่อสแกน
    };

    const handleScanResult = (result: string | null) => {
        if (result) {
            console.log(result); // จัดการผลลัพธ์ที่ได้รับ
        }
    };

    const handleScanError = (error: any) => {
        console.log(error);
    };

    const handleSubmit = async () => {
      const allReceived = data.every(item => item.receivedQty === item.qty);
      if (allReceived) {
        try {
          // ดึงโทเค็นจาก Local Storage
          const token = localStorage.getItem('access_token');
    
          const base64ToFile = async (base64String: string, filename: string): Promise<File> => {
            const res = await fetch(base64String);
            const blob = await res.blob();
            console.log("Blob from base64ToFile:", blob); 
            return new File([blob], filename, { type: 'image/jpeg' }); // เปลี่ยน type ตามประเภทรูปภาพ
          };
    
          const uploadImage = async (image: string, imageTypeID: number, sku: string | null) => {
            try {
                const timestamp = new Date().toISOString().replace(/[:.]/g, '-'); // แปลงวันที่และเวลาเป็นสตริงที่ใช้ได้ในชื่อไฟล์
                let filename = '';
            
                if (imageTypeID === 1) {
                  filename = `beforeopen_${timestamp}.jpg`;
                } else if (imageTypeID === 2) {
                  filename = `afteropen_${timestamp}.jpg`;
                } else if (imageTypeID === 3 && sku) {
                  filename = `${sku}_${timestamp}.jpg`;
                } else {
                  filename = `image_${timestamp}.jpg`; // กรณีที่ไม่ตรงกับเงื่อนไขใดๆ
                }
            
                const file = await base64ToFile(image, filename);
                const formData = new FormData();
                formData.append('file', file);
                formData.append('orderNo', selectedOrderNo!);
                formData.append('imageTypeID', imageTypeID.toString());
                if (sku) {
                    formData.append('sku', sku);
                }
                console.log("FormData:", formData); // check data
        
                const response = await api.post('/api/import-order/upload-photo', formData, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'multipart/form-data',
                    },
                });
                console.log("Response data from uploadImage:", response.data); 
                console.log("Response from uploadImage: ", response); 
        
                if (response.data && response.data.data && response.data.data.filePath) {
                    return response.data.data.filePath; 
                } else {
                    console.error("filePath not found in response:", response);
                    return null; 
                }
            } catch (error) {
                console.error("Error uploading image:", error);
                return null;
            }
         };

          // อัปโหลดภาพและบันทึกพาธสำหรับแต่ละภาพ
          const step1ImagePath = images.step1 ? await uploadImage(images.step1, 1, null) : null;
          const step2ImagePath = images.step2 ? await uploadImage(images.step2, 2, null) : null;
          const itemImages = await Promise.all(
            data.map(async (item) => {
                console.log("Images for item:", item.key, item.image);
                if (item.image) { 
                    const filePath = await uploadImage(item.image, 3, item.sku);
                    return { ...item, filePath };
                }
                return item;
            })
        );

        console.log("itemImages:", itemImages);
    
          // เตรียมข้อมูลสำหรับส่งไปยัง API
          const requestData = {
            Identifier: selectedOrderNo, // ใช้ OrderNo หรือ TrackingNo ที่เลือก
            ImportLines: [
              ...itemImages.map(item => ({
                SKU: item.sku,
                QTY: item.qty,
                ReturnQTY: item.receivedQty,
                Price: item.price,
                ImageTypeID: 3, 
                FilePath: item.filePath,
              })),
              {
                SKU: null,
                QTY: null,
                ReturnQTY: null,
                Price: null,
                ImageTypeID: 1, 
                FilePath: step1ImagePath,
              },
              {
                SKU: null,
                QTY: null,
                ReturnQTY: null,
                Price: null,
                ImageTypeID: 2,
                FilePath: step2ImagePath,
              },
            ],
            CreateBy: userID, 
          };
    
          console.log(requestData);
          // ส่งข้อมูลไปยัง API
          const response = await api.post(`/api/import-order/confirm-receipt/${selectedOrderNo}`, requestData, {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          });
    
          // ตรวจสอบผลลัพธ์จาก API
          if (response.status === 200) {
            notification.success({
              message: 'ส่งข้อมูลสำเร็จ',
              description: 'ข้อมูลถูกส่งเรียบร้อยแล้ว',
              placement: 'topRight',
            });
    
            // รีเซ็ตฟอร์มและตาราง
            setData([]);
            setSelectedOrderNo(null);
            setShowTable(false);
            setShowSteps(false);
            setCurrentStep(0);
            setImages({
              step1: null,
              step2: null,
            });
          } else {
            notification.error({
              message: 'เกิดข้อผิดพลาด',
              description: 'ไม่สามารถส่งข้อมูลได้ กรุณาลองใหม่อีกครั้ง',
              placement: 'topRight',
            });
          }
        } catch (error) {
          console.error("Error submitting data:", error);
          notification.error({
            message: "เกิดข้อผิดพลาด",
            description: "ไม่สามารถส่งข้อมูลได้ กรุณาลองใหม่อีกครั้ง",
            placement: 'topRight',
          });
        }
      } else {
        let warningMessage = 'กรุณายืนยันจำนวนรับเข้าทั้งหมดก่อนส่งข้อมูล';
        if (!allReceived) {
          warningMessage += '\nกรุณาถ่ายรูปให้ครบก่อนส่งข้อมูล';
        }
    
        notification.warning({
          message: 'ข้อมูลไม่ครบถ้วน',
          description: warningMessage,
          placement: 'topRight',
        });
      }
    };

    const handleKeyDown = (event: { key: string; }) => {
        if (event.key === 'Enter' && skuInput.trim()) {
            handleConfirmReceived();
        }
    };

    useEffect(() => {
        // เพิ่ม event listener เมื่อ component ถูก mount
        window.addEventListener('keydown', handleKeyDown);

        // ลบ event listener เมื่อ component ถูก unmount
        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        };
    }, [skuInput]);
    

    return (
        <ConfigProvider
  
        >
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "#023E8A" }}>
                ถ่าย/สแกนรับเข้าการคืนสินค้าทั่วไป
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
                <Row align="middle" justify="center" style={{ marginTop: "20px", width:'100%'  }}>
                    <Col span={7} >
                        <Form.Item
                        id="selectedSalesOrder"
                            layout="vertical"
                            label={<span style={{ color: '#657589' }}>เลือกเลขออเดอร์</span>}
                            name="selectedSalesOrder"
                            rules={[{ required: true, message: 'กรุณากรอกเลข Order!' }]}
                        >
                            <Select
                                showSearch
                                style={{ height: 40 }}
                                placeholder="Search Order Number"
                                optionFilterProp="label"
                                onChange={handleSelectChange}
                                options={orderOptions}
                            />
                        </Form.Item>
                    </Col>
                </Row>

                {showSteps && !showTable && (
                    <>
                    <Steps current={currentStep} 
                        onChange={onChange}
                        style={{ marginTop: '20px', width: '100%' }}
                        className="custom-steps"
                    >
                        <Steps.Step title="ถ่ายก่อนเปิดสินค้า" className="custom-steps"/>
                        <Steps.Step title="ถ่ายหลังเปิดสินค้า"/>
                        {data.map((item, index) => (
                        <Steps.Step
                            key={item.key}
                            title={
                                <div style={{ textAlign: 'center', fontSize: '13px', fontWeight: 'normal' }}>
                                    <div>{item.sku}</div>
                                </div>
                            }
                        />
                        ))}
                    </Steps>

                    <Row justify="center" align="middle" style={{ marginTop: 20 }}>
                        <Col span={24} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', marginTop: '20px' }}>
                            {!images[`step${currentStep + 1}`] ? (
                                <Webcam
                                id="webcam"
                                    audio={false}
                                    ref={webcamRef}
                                    screenshotFormat="image/jpeg"
                                    style={{ width: '100%', maxWidth: '400px' }}
                                />
                            ) : (
                                <img
                                    src={images[`step${currentStep + 1}`]!}
                                    alt={`Captured ${currentStep + 1}`}
                                    style={{ width: '100%', maxWidth: '400px', border: '1px solid #ccc' }}
                                />
                            )}
                        </Col>

                        <Col span={24} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', marginTop: '10px' }}>
                            <Button 
                                onClick={toggleCamera} 
                                style={{ marginTop: 10 , marginRight: '20px'}}
                            >
                                เปลี่ยนกล้อง
                            </Button>
                            {currentStep > 0 && (  // แสดงปุ่ม "ย้อนกลับ" เมื่อ currentStep > 0
                            <Button
                                id="back"
                                style={{marginTop: '10px', marginRight: '20px', width: '100px'}}
                                type="default"
                                onClick={handleBackStep}
                            >
                                ย้อนกลับ
                            </Button>
                            )}
                            {!images[`step${currentStep + 1}`] ? (
                                    <Button
                                        id="Takepicture"
                                        type="primary"
                                        onClick={capturePhoto}
                                        style={{ marginTop: '10px', display: 'flex', width: '100px', alignItems: 'center' }}>
                                        ถ่ายรูป
                                    </Button>
                            ) : (
                                <>
                                    <Button
                                        id="Retakepicture"
                                        icon={<RedoOutlined />}
                                        style={{ marginTop: '10px', marginRight: '20px', width: '100px'}}
                                        type="default"
                                        onClick={retakePhoto}
                                    >
                                        ถ่ายใหม่
                                    </Button>
                                    
                                    <Button
                                        id="Next"
                                        style={{ marginTop: '10px', width: '100px' }}
                                        type="primary"
                                        onClick={handleNextStep}
                                    >
                                        ถัดไป
                                    </Button>
                                </>
                            )}
                        </Col>
                    </Row>
                    </>
                )}
                {showTable && (
                    <>
                    <Divider orientation="left" style={{color: '#657589'}}>รับเข้า SKU</Divider>
                        <Row justify="center" align="middle" style={{ marginTop: 20 }}>
                            <Col span={24} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', marginTop: '10px' }}>
                                <Form.Item  
                                    id="inputsku"
                                    layout="vertical"
                                    label={<span style={{ color: '#657589' }}>กรอก/สแกน รหัสสินค้าเพื่อยืนยันการรับเข้า</span>}
                                    style={{ marginBottom: 0 }}
                                >
                                    <Input
                                        placeholder="ระบุ SKU"
                                        value={skuInput}
                                        onChange={(e) => setSkuInput(e.target.value)}
                                        style={{ width: '375px', marginRight: '10px', height: 40 }}
                                    />
                                    <Button 
                                        id="สแกนSKU"
                                        icon={<ScanOutlined />} 
                                        onClick={handleScanSku} 
                                        style={{ height: 40, marginLeft: 5 }}
                                    >
                                        สแกน SKU
                                    </Button>
                                </Form.Item>
                            </Col>
                            <Row gutter={40} style={{ marginTop: '20px' }}> 
                                <Col span={12} >
                                    <Button id="ยืนยันการรับเข้า" type="primary" onClick={handleConfirmReceived} disabled={!skuInput} style={{}}>
                                        ยืนยันการรับเข้า
                                    </Button>
                                </Col>
                                <Col span={12} >
                                    <Button id="ส่งข้อมูล" type="primary" onClick={handleSubmit} style={{background:'#14C11B'}} >
                                        ส่งข้อมูล
                                    </Button>
                                </Col>
                            </Row>
                        </Row>
                    <Table
                        dataSource={data}
                        columns={columns}
                        pagination={false}
                        rowKey="key"
                        style={{ marginTop: '50px' }}
                    />
                    </>
                )}
                {showWebcam && (
                    <Modal
                        open={showWebcam}
                        footer={null}
                        onCancel={handleCancelModal}
                        title={
                            <div style={{ textAlign: 'center', fontSize: '16px',fontWeight: 'normal'  }}> {/* ขนาดปกติสามารถปรับได้ */}
                                ถ่ายรูป SKU:<br />
                                <div style={{color:'#1890FF'}}>{skuName}</div> {/* เน้น SKU */}
                            </div>
                        }
                        
                    >
                        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '300px' }}>
                            <Webcam
                            id="Webcam2"
                                audio={false}
                                ref={webcamRef}
                                screenshotFormat="image/jpeg"
                                style={{ width: '100%', maxWidth: '400px' }}
                            />
                        </div>
                        <Row justify="center" style={{ marginTop: '10px' }}>
                            <Button id="ถ่ายรูป" type="primary" onClick={handleCapturePhoto}>
                                ถ่ายรูป
                            </Button>
                        </Row>
                    </Modal>
                )}
                {showScanner && (
                    <Modal
                        open={showScanner}
                        footer={null}
                        onCancel={() => setShowScanner(false)}
                        title="Scan SKU"
                    >
                            
                    </Modal>
                )}
                        
                </Layout.Content>
            </Layout>
        </ConfigProvider>
        );
    };
        
    export default SaleReturn;
    
