import React, { useEffect, useState } from 'react';
import { notification, Alert, Popconfirm, Layout, Button, ConfigProvider, Form, Row, Col, Select, FormProps, Input, DatePicker, Table, Modal, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { DeleteOutlined, LeftOutlined, PlusCircleOutlined } from '@ant-design/icons';


const SRPage = () => {

  const navigate = useNavigate();
  const [dataSource, setDataSource] = useState<DataSourceItem[]>([]);
  const [importedData, setImportedData] = useState<any[]>([]);
  const [form] = Form.useForm();
  const [selectedSalesOrder, setSelectedSalesOrder] = useState('');
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [selectedSku, setSelectedSku] = useState<string | undefined>(undefined);
  const [isChecked, setIsChecked] = useState(false);
  const [formValid, setFormValid] = useState(false);
  const [selectedData, setSelectedData] = useState<DataItem[]>([]);
  const [randomNumber, setRandomNumber] = useState<number>(0);
  const [error, setError] = useState<string | null>(null);
  const [state, setState] = useState<string | number>("");


  useEffect(() => {
    form.setFieldsValue(checkSR[0]); // Set initial values
  }, [form]);

  interface DataSourceItem extends FormValues {
    key: number;

    warehouse_to?: string;
    // เพิ่มคุณสมบัติอื่นๆ ตามต้องการ
  }
  const handleCreateSR = async () => {
    try {
      await form.validateFields(); // Validate form fields

      const allDataFilled = dataSource.every(item => item.warehouse_to);
      if (!allDataFilled) {
        notification.warning({
          message: "กรุณากรอกข้อมูลให้ครบก่อนสร้าง SR",
          description: "กรุณาเลือก Warehouse Form ในตารางด้วย",
        });
        return;
      }

      // Generate random number
      const randomNumber = generateRandomNumber();
      form.setFieldsValue({ SR_Create: randomNumber }); // Set random number
      setRandomNumber(randomNumber); // Save it to state
      notification.success({
        message: 'สำเร็จ',
        description: `Create SR สำเร็จ! เลขสุ่มที่สร้างคือ: ${randomNumber}`,
    });

      setIsSubmitted(true);
    } catch (error) {
      handleError(error);
    }
    
  };


  // Function to generate a random number (4-digit)
  const generateRandomNumber = () => {
    return Math.floor(Math.random() * 10000);
  };



  const handleError = (error: any) => {
    notification.warning({
      message: "กรุณากรอกข้อมูลให้ครบ",

    });
  };
  const handleChange = (value: string, key: number, field: string) => {
    const updatedDataSource = selectedData.map((item) => {
        if (item.key === key) {
            return { ...item, [field]: value };
        }
        return item;
    });
    console.log("updatedDataSource:", updatedDataSource);
    setSelectedData(updatedDataSource);
};


  




  const handleBack = () => {
    navigate('/CreateReturn'); // Navigate to CreateReturn page
  };

  const handleCheck = () => {
    if (!selectedSalesOrder) {
      message.error("กรุณาเลือก Sales Order ก่อน");
      return;
    }
  
    const selectedOption = options.find(option => option.value === selectedSalesOrder);
  
    if (!selectedOption) {
      message.error("ไม่พบ Sales Order ที่เลือก");
      return;
    }
  
    const relatedOrder = checkSR.find(order => order.Sales_Order === selectedOption.label);
    const relatedData = data[selectedOption.label] || []; // ดึงข้อมูลที่ตรงกันจาก data
  
    if (relatedOrder) {
      // ตั้งค่าฟิลด์ข้อมูลจาก checkSR และ data รวมกัน
      form.setFieldsValue({
        Sales_Order: relatedOrder.Sales_Order,
        Tracking_Order: relatedOrder.Tracking_Order,
        SR_Create: relatedOrder.SR_Create,
        SO_Status: relatedOrder.SO_Status,
        MKP_Status: relatedOrder.MKP_Status,
        // อัปเดตฟิลด์เพิ่มเติมจาก relatedData
        ...relatedData
      });
      setSelectedData(relatedData); // เก็บข้อมูลที่อัปเดตใน state
      setIsChecked(true);
    } else {
      message.error("ไม่พบข้อมูลที่ตรงกัน");
    }
  };
  
  

  const handleSubmitData = () => {
    Modal.confirm({
      title: 'ยืนยันการส่งข้อมูล',
      content: 'คุณต้องการส่งข้อมูลนี้ใช่หรือไม่?',
      okText: 'ใช่',
      cancelText: 'ไม่',
      onOk: () => {
        console.log("Sending data:", selectedData);
        setSelectedData([]);
        form.resetFields();
        setIsSubmitted(false);
        notification.success({
          message: 'ส่งข้อมูล สำเร็จ',
          description: 'ข้อมูลทั้งหมดได้ถูกส่งรียบร้อยแล้ว',
      });
      },
    });
  };

  const options = [
    { value: '1', label: 'SOA2409-12345' },
    { value: '2', label: 'SOA2409-12346' },
    { value: '3', label: 'SOA2409-12347' },
    { value: '4', label: 'SOA2409-12348' },
    { value: '5', label: 'SOA2409-12349' },
    { value: '6', label: 'SOA2409-12350' },
  ];

  const checkSR = [
    { Sales_Order: "SOA2409-12345", Tracking_Order: "2409901234896", SR_Create: "Null", SO_Status: "invoice", MKP_Status: 'Cancel' },
    { Sales_Order: "SOA2409-12346", Tracking_Order: "2409901234897", SR_Cerate: "Null", SO_Status: "invoice", MKP_Status: 'Cancel' },
    { Sales_Order: "SOA2409-12347", Tracking_Order: "2409901234898", SR_Create: "Null", SO_Status: "invoice", MKP_Status: 'Cancel' },
    { Sales_Order: "SOA2409-12348", Tracking_Order: "2409901234899", SR_Create: "Null", SO_Status: "invoice", MKP_Status: 'Cancel' },
    { Sales_Order: "SOA2409-12349", Tracking_Order: "2409901234900", SR_Create: "Null", SO_Status: "invoice", MKP_Status: 'Cancel' },
    { Sales_Order: "SOA2409-12350", Tracking_Order: "2409901234901", SR_Create: "Null", SO_Status: "invoice", MKP_Status: 'Cancel' },
  ];

  const data: Record<string, { key: number,SKU: string; SKU_Name: string; QTY: string; Price: string; Location_to: string; Warehouse_Form: string }[]> = {
    "SOA2409-12345": [
      {key: 1, SKU: "G097171-ARM01-BL", SKU_Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", QTY: "2", Price: "2,000", Location_to: "Return", Warehouse_Form: '' },

    ],
    "SOA2409-12346": [
      {key: 2, SKU: "G097171-ARM02-BL", SKU_Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", QTY: "3", Price: "3,000", Location_to: "Return", Warehouse_Form: '' },

    ],
    "SOA2409-12347": [
      {key: 3, SKU: "G097171-ARM03-BL", SKU_Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", QTY: "4", Price: "4,000", Location_to: "Return", Warehouse_Form: '' },

    ],
    "SOA2409-12348": [
      {key: 4, SKU: "G097171-ARM04-BL", SKU_Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", QTY: "5", Price: "5,000", Location_to: "Return", Warehouse_Form: '' },

    ],

  };

  const columns = [
    { title: 'SKU', dataIndex: 'SKU' },
    { title: 'SKU_Name', dataIndex: 'SKU_Name' },
    { title: 'QTY', dataIndex: 'QTY', },
    { title: 'Price', dataIndex: 'Price', },
    { title: 'Location_to', dataIndex: 'Location_to', },

    {
      title: 'Warehouse Form',
      dataIndex: 'Warehouse_Form',
      key: 'Warehouse_Formt',
      render: (text: any, record: DataItem, index: number) => (
        <Form.Item
          name={['selectedData', index, 'Warehouse_Form']}
          rules={[{ required: true, message: 'กรุณาเลือก Warehouse Form!' }]}
          style={{ marginBottom: 0 }}
        >
          <Select
            style={{ width: 150, height: 40, borderRadius: 50 }}
            showSearch
            placeholder="Warehouse Form"
            options={[
              { value: 'MMT', label: 'MMT' },
              { value: 'RBN', label: 'RBN' },
            ]}
            onChange={(value)=>handleChange(value,record.key,"Warehouse_Form")}



          />
        </Form.Item>
      ),
    },


  ];
  const handleonChange = () => {
    const values = form.getFieldsValue();
    console.log("values-----------",values)
    setFormValid(
        values.Date && values.TrackingNumber && values.TransportType
    );
};

  console.log("formValid:", formValid);
  console.log("dataSource length:", dataSource.length);

  useEffect(() => {
    form.setFieldsValue(checkSR[0]); // Set initial values
    const data = form.getFieldsValue()
    console.log(selectedData[0]?.Warehouse_Form)
    

  }, [form]);

  interface DataItem {
    key: number;
    SKU: string;
    SKU_Name: string;
    QTY: string;
    Price: string;
    Warehouse_Form: string;
  }

  interface FormValues {
    TransportType: any;
    Date: any;
    SKU: string;
    QTY: number;
    SKU_Name: string;
  }
  const handleCancel = () => {
    form.resetFields();       // รีเซ็ตค่าในฟอร์มทั้งหมด
    setSelectedData([]);        // รีเซ็ตข้อมูล dataSource
    setIsSubmitted(false);    // รีเซ็ตสถานะ isSubmitted
    notification.success({
      message: 'Cancel สำเร็จ',
      description: 'ข้อมูลทั้งหมดได้ถูกยกเลิกเรียบร้อยแล้ว',
  });
   
};
  return (
    <ConfigProvider>

      <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
        Create SR Return
      </div>
      <Layout>
        <Layout.Content style={{
          margin: "24px",
          padding: 20,
          minHeight: 200,
          background: "#fff",
          borderRadius: "8px",
          display: 'flex',

        }}>

          <Button
            onClick={handleBack}
            style={{ background: '#98CEFF', color: '#fff' }}
          >
            <LeftOutlined style={{ color: '#fff', marginRight: 5 }} />
            Back
          </Button>
          <Form
            layout="vertical"
            form={form}
            style={{ width: '100%', marginTop: '40px' }}
          >
            <Row gutter={30} justify="center" align="middle" style={{ width: '100%' }}>
              <Col>
                <Form.Item
                  label={<span style={{ color: '#657589' }}>กรอกเลข SO/Order ที่ต้องการสร้างSR</span>}
                  name="selectedSalesOrder"
                  rules={[{ required: true, message: 'กรุณากรอกเลข SO/Order ที่ต้องการสร้างSR!' }]}
                >
                  <Select
                    showSearch
                    style={{ height: 40, width:300 }}
                    placeholder="Search to Select"
                    optionFilterProp="label"
                    value={selectedSalesOrder}
                    onChange={(value) => {
                      setSelectedSalesOrder(value);
                    }}
                    options={options}
                  />
                </Form.Item>

              </Col>
              <Col>
                <Button type="primary" style={{ width: 100, height: 40, marginTop: 4 }} onClick={handleCheck}>
                  Check
                </Button>

              </Col>
            </Row>
          </Form>
        </Layout.Content>
        <Layout.Content style={{
          
          marginRight:24 ,
          marginLeft:24 ,
          padding: 36,
          minHeight: 360,
          background: "#fff",
          borderRadius: "8px",
          justifyContent: 'center', // Center content horizontally
          alignItems: 'center', // Center content vertically
        }}>
          {isChecked && (
            <div>
              <Form
                form={form}
                layout="vertical"
                onValuesChange={handleonChange}

                style={{ width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}
              >
                <div style={{ width: '100%', maxWidth: '800px' }}>
                  <Row gutter={16} style={{ marginTop: '10px', justifyContent: 'center' }}>
                    <Col span={8}>
                      <Form.Item label={<span style={{ color: '#657589' }}>Sale Order:</span>} name="Sales_Order">
                        <Input style={{ width: '100%', height: '40px' }} disabled />
                      </Form.Item>
                    </Col>
                    <Col span={8}>
                      <Form.Item label={<span style={{ color: '#657589' }}>Tracking Order:</span>} name="Tracking_Order">
                        <Input style={{ width: '100%', height: '40px' }} disabled />
                      </Form.Item>
                    </Col>
                    <Col span={8}>
                      <Form.Item
                        label={<span style={{ color: '#657589' }}>SR_Create:</span>}
                        name="SR_Create"
                      >
                        <Input style={{ width: '100%', height: '40px' }} placeholder="SR Create" disabled />
                      </Form.Item>
                    </Col>
                    <Col span={8}>
                      <Form.Item label={<span style={{ color: '#657589' }}>SO Status:</span>} name="SO_Status">
                        <Input style={{ width: '100%', height: '40px' }} disabled />
                      </Form.Item>
                    </Col>
                    <Col span={8}>
                      <Form.Item label={<span style={{ color: '#657589' }}>MKP Status:</span>} name="MKP_Status">
                        <Input style={{ width: '100%', height: '40px' }} disabled />
                      </Form.Item>
                    </Col>
                    <Col span={8}>
                    </Col>
                    <Col span={8}>
                      <Form.Item label={<span style={{ color: '#657589' }}>วันที่คืน:</span>} name="Date" rules={[{ required: true, message: 'กรุณาเลือกวันที่คืน' }]}>
                        <DatePicker style={{ width: '100%', height: '40px' }} placeholder="เลือกวันที่คืน" />
                      </Form.Item>
                    </Col>
                    <Col span={8}>
                      <Form.Item label={<span style={{ color: '#657589' }}>กรอกเลข Tracking:</span>} name="TrackingNumber" rules={[{ required: true, message: 'กรุณากรอกเลข Tracking!' }]}>
                        <Input style={{ width: '100%', height: '40px' }} placeholder="กรอกเลข Tracking" />
                      </Form.Item>
                    </Col>
                    <Col span={8}>
                      <Form.Item
                        label={<span style={{ color: '#657589' }}>Transport Type:</span>}
                        name="TransportType"
                        rules={[{ required: true, message: 'กรุณาเลือกTransportType ' }]}
                      >
                        <Select
                          style={{ width: '100%', height: '40px', borderWidth: '1px', }}


                          showSearch
                          placeholder="TransportType"
                          optionFilterProp="label"

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
                </div>
              </Form>

              <Table
              components={{
                header: {
                  cell: (props: React.HTMLAttributes<HTMLElement>) => (
                    <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                  ),
                },
              }}
                style={{ width: '100%', tableLayout: 'fixed', marginTop: '50px' }} // Ensure the table takes full width and is fixed layout
                scroll={{ x: 'max-content' }}
                dataSource={selectedData}
                columns={columns}
                pagination={false}
                rowKey={(record) => record.SKU}
              />

              <Row justify="center" style={{ marginTop: '20px' }}>
                {!isSubmitted ? (
                  <Button
                    type="primary"
                    onClick={handleCreateSR} // เรียกใช้ฟังก์ชันนี้เมื่อ !isSubmitted
                    style={{ width: 100, height: 40 ,marginRight: '10px' }}
                    disabled={!formValid || selectedData.length === 0 || !selectedData.every(item => item.Warehouse_Form)}
                  >
                    Create IJ
                  </Button>
                ) : (
                  <Button
                    type="primary"
                    onClick={handleSubmitData} // เรียกใช้ฟังก์ชันนี้เมื่อ isSubmitted


                    style={{ width: 100, height: 40 }}
                    disabled={!formValid || selectedData.length === 0 || !selectedData.every(item => item.Warehouse_Form)}
                  >
                    ส่งข้อมูล
                  </Button>
                )}


<Popconfirm
        title="ต้องการยกเลิกหรือไม่?"
        description="คุณแน่ใจหรือไม่ว่าต้องการยกเลิกข้อมูลทั้งหมด?"
        onConfirm={handleCancel} // ยืนยันการยกเลิก
        okText="ใช่"
        cancelText="ไม่"
    >
        <Button 
            type="default" 
            style={{ width: 100, height: 40 }}
        >
            Cancel
        </Button>
    </Popconfirm>
              </Row>
            </div>

          )}

        </Layout.Content>
      </Layout>
    </ConfigProvider>
  );
};



export default SRPage;

