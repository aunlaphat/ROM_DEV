import React, { useEffect, useRef, useState } from 'react';
import { Steps, Col, ConfigProvider, Form, Layout, Row, Select, Button, Table, Modal, Input, notification, Divider, Popconfirm } from 'antd';
import Webcam from 'react-webcam';
import { CameraOutlined, RedoOutlined, DeleteOutlined, ScanOutlined, CheckCircleOutlined, WarningOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons';
import { QrReader, QrReaderProps } from 'react-qr-reader';

const orderOptions = [
    { value: '1', label: 'SOA2409-12345' },
    { value: '2', label: 'SOA2409-12346' },
    { value: '3', label: 'SOA2409-12347' },
    { value: '4', label: 'SOA2409-12348' },
    { value: '5', label: 'SOA2409-12349' },
    { value: '6', label: 'SOA2409-12350' },
];

interface CustomQrReaderProps extends QrReaderProps {
    onScan: (result: string | null) => void;
    onError: (error: any) => void;
}

interface DataType {
    key: string;
    sku: string;
    name: string;
    qty: number;
    receivedQty: number;
    amount: string;
    image: string | null;
}

const SaleReturn: React.FC = () => {
    const [scanResult, setScanResult] = useState<string | null>(null);
    const [showScanner, setShowScanner] = useState<boolean>(false);
    const [skuInput, setSkuInput] = useState<string>('');
    const [currentStep, setCurrentStep] = useState(0);
    const [currentRecordKey, setCurrentRecordKey] = useState<string | null>(null);
    const [skuName, setSkuName] = useState<string | null>(null);
    const [showWebcam, setShowWebcam] = useState(false);
    const [showSteps, setShowSteps] = useState(false);
    const [showTable, setShowTable] = useState(false);
    const [current, setCurrent] = useState(0);
    const [images, setImages] = useState<{ [key: string]: string | null }>({
        step1: null,
        step2: null,
    });
    const webcamRef = useRef<Webcam>(null);
    
    const onChange = (current: number) => {
        setCurrentStep(current);
    };
    const [cameraFacingMode, setCameraFacingMode] = useState<'user' | 'environment'>('environment');

    const toggleCamera = () => {
        setCameraFacingMode((prevMode) => (prevMode === 'user' ? 'environment' : 'user'));
    };
    
    
   
    const handleSelectChange = (value: string) => {
        setShowSteps(true);
        setCurrentStep(0);
        setShowTable(false);
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
        if (currentStep < 4) { // Assuming there are 5 steps
            setCurrentStep(prevStep => prevStep + 1);
         } else {
            setShowTable(true);
        }
    };
    


    const handleDelete = (recordKey: string) => {
        setData((prevData) => prevData.filter(item => item.key !== recordKey));
    };

    const columns = [
        { title: 'SKU', dataIndex: 'sku', key: 'sku' ,id:'sku'},
        { title: 'Name', dataIndex: 'name', key: 'name' ,id:'name'},
        { title: 'QTY', dataIndex: 'qty', key: 'qty' ,id:'qty'},
        { title: 'จำนวนรับเข้า', dataIndex: 'receivedQty', key: 'receivedQty' ,id:'receivedQty'},
        { title: 'Amount', dataIndex: 'amount', key: 'amount',id:'amount' },
        {
            title: 'Return',
            id:'Return' ,
            dataIndex: 'Return',
            key: 'Return',
            render: (_: any, record: DataType) => {
                // Check if the current record has receivedQty equal to qty
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
    render: (_: any, record: DataType) => {
        // ตรวจสอบว่า record.image มีค่าหรือไม่
        const stepImage = images[`step${parseInt(record.key, 10)+2}`]; // ปรับให้ตรงกับ key ของ record
        return stepImage ? (
            <img src={stepImage} alt="Return" style={{ width: '100px' }} />
        ) : (
            <Button
                style={{ background: '#02C39A' }}
                icon={<CameraOutlined />}
                type="primary"
                onClick={() => handleTakePhoto(record.key, record.sku)} // ใช้ sku แทน name
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
        render: (_: any, record: DataType) => (
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
              onConfirm={() => handleDelete(record.key)} // เรียกใช้ฟังก์ชัน handleDelete เมื่อกดยืนยัน
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

    const [data, setData] = useState<DataType[]>([
        {
            key: '1',
            sku: 'G090108-EF05',
            name: 'Bewell Foot Rest EF-05',
            qty: 3,
            receivedQty: 0,
            amount: '500',
            image: null,
        },
        {
            key: '2',
            sku: 'G091116-PC08-GY',
            name: 'Bewell Cooling Blanket Single PC-08(Gray)',
            qty: 2,
            receivedQty: 0,
            amount: '600',
            image: null,
        },
        {
            key: '3',
            sku: 'G091116-PC09-BL',
            name: 'Bewell Cooling Blanket King PC-08(Blue)',
            qty: 3,
            receivedQty: 0,
            amount: '700',
            image: null,
        },
    ]);

    const handleTakePhoto = (recordKey: string, sku: string) => {
        setCurrentRecordKey(recordKey);
        setSkuName(sku); // Use sku instead of name
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
        setSkuName(currentItem?.sku || null); // Use sku instead of name
        setShowWebcam(true);
    };

    const handleCapturePhoto = () => {
        if (webcamRef.current && currentRecordKey) {
            const imageSrc = webcamRef.current.getScreenshot();
            if (imageSrc) {
                setData((prevData) =>
                    prevData.map((item) =>
                        item.key === currentRecordKey ? { ...item, image: imageSrc } : item
                    )
                );
                setSkuName(data.find(item => item.key === currentRecordKey)?.sku || null);
            }
            setShowWebcam(false); // Close webcam after capturing photo
        }
    };
    
    const handleCancelModal = () => {
        setShowWebcam(false);
        setCurrentRecordKey(null);
        setSkuName(null);
        setImages((prevImages) => ({ ...prevImages, [`step${currentStep + 1}`]: null })); // Clear image if needed
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
    
        setSkuInput(''); // ล้างช่องกรอกข้อมูลหลังจากการแจ้งเตือน
    };
    
    notification.config({
        placement: 'topRight',
        duration: 3, // ระยะเวลาการแสดง notification
        maxCount: 3, // จำนวน notification สูงสุดที่จะแสดงพร้อมกัน
    });

    const handleScanSku = () => {
        // เปิด/ปิดการใช้งานกล้องเพื่อสแกน
        setShowScanner(true);
    };

    const handleScanResult = (result: string | null) => {
        if (result) {
            console.log(result);
            // จัดการผลลัพธ์ที่ได้รับ
        }
    };

    const handleScanError = (error: any) => {
        console.log(error);
    };

    const handleSubmit = () => {
        const allReceived = data.every(item => item.receivedQty === item.qty);
        
    
        if (allReceived ) {
            console.log('ข้อมูลที่ส่ง:', data);
            setData([]); // ลบข้อมูลในตาราง
            notification.success({
                message: 'ส่งข้อมูลสำเร็จ',
                description: 'ข้อมูลถูกส่งเรียบร้อยแล้ว',
                placement: 'topRight',
            });
        } else {
            let warningMessage = 'กรุณายืนยันจำนวนรับเข้าทั้งหมดก่อนส่งข้อมูล';
            if (!allReceived ) {
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
        <ConfigProvider>
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                Sale Return
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
                    <Row align="middle" justify="center" style={{ marginTop: "20px",width:'100%'  }}>
                        <Col span={12} >
                            <Form.Item
                            id="selectedSalesOrder"
                                layout="vertical"
                                label={<span style={{ color: '#657589' }}>กรอกเลข Order</span>}
                                name="selectedSalesOrder"
                                rules={[{ required: true, message: 'กรุณากรอกเลข Order!' }]}
                            >
                                <Select
                                    showSearch
                                    style={{ height: 40}}
                                    placeholder="Search to Select"
                                    optionFilterProp="label"
                                    onChange={handleSelectChange}
                                    options={orderOptions}
                                />
                            </Form.Item>
                        </Col>
                    </Row>

                    {showSteps && !showTable && (
                        <>
                            <Steps  current={currentStep} onChange={onChange}
                              style={{ marginTop: '20px', width: '100%' }}>
                        <Steps.Step title="ถ่ายก่อนเปิดสินค้า" />
                        <Steps.Step title="ถ่ายหลังเปิดสินค้า" />
                        {data.map((item, index) => (
                <Steps.Step
                    key={item.key}
                    title={
                        <div style={{ textAlign: 'center', fontSize: '13px', fontWeight: 'normal' }}>
                          
                            <div style={{ color: '#1890FF'  }}>{item.sku}</div>
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
                        
                        <Button onClick={toggleCamera} style={{ marginTop: 10 }}>
            เปลี่ยนกล้อง
        </Button>
                        {currentStep > 0 && (  // แสดงปุ่ม "ย้อนกลับ" เมื่อ currentStep > 0
                <Button
                id="back"
                    style={{marginRight: '20px', width: '100px', color: '#35465B'}}
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
                                    style={{ display: 'flex', alignItems: 'center' }}>
                                    ถ่ายรูป
                                </Button>
                            ) : (
                                <>
                                    <Button
                                    id="Retakepicture"
                                        icon={<RedoOutlined />}
                                        style={{ marginRight: '20px', width: '100px', color: '#35465B' }}
                                        type="default"
                                        onClick={retakePhoto}
                                    >
                                        ถ่ายใหม่
                                    </Button>
                                    
                                      
                                    <Button
                                    id="Next"
                                        style={{ width: '100px' }}
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
                    label={<span style={{ color: '#657589' }}>กรอก SKU เพื่อรับเข้า</span>}
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
                            visible={showWebcam}
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
                            visible={showScanner}
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
    
