import React, { useState } from 'react';
import {notification , Form, Input, InputNumber, DatePicker, Button, Row, Col, Table, ConfigProvider, Layout, Select, Modal, message, Popconfirm, Divider } from 'antd';
import moment from 'moment';
import { DeleteOutlined, LeftOutlined, PlusCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

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

    const onChange = () => {
        const values = form.getFieldsValue();
        setFormValid(
            values.IJ && values.Date && values.SKU && values.QTY
        );
    };


    const handleAdd = () => {
        form.validateFields()
          .then((values) => {
            if (!values.SKU) {
              notification.warning({
                message: "มีข้อสงสัย",
                description: "กรุณากรอกข้อมูลให้ครบก่อนเพิ่ม!",
              });
              return;
            } 
            
            setDataSource((prevData) => [
              ...prevData,
              { key: Date.now(), ...values }
            ]);
      
            // เพิ่มการแจ้งเตือนเมื่อเพิ่มข้อมูลสำเร็จ
            notification.success({
              message: "Add ข้อมูลสำเร็จ",
            });
      
            form.resetFields(['SKU', 'QTY', 'SKU_Name']);
          })
          .catch((errorInfo) => {
            // แสดง notification warning เมื่อเกิดข้อผิดพลาด
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
                    style={{ width: '70%' }}
                    onChange={(value) => handleChange(value, record.key, "warehouse_form")}
                    options={Warehouse}
                />
            ),
        },
        {
            title: "Location Form",
            dataIndex: "location_form",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '70%' }}
                    onChange={(value) => handleChange(value, record.key, "location_form")}
                    options={Location}
                />
            ),
        },
        {
            title: "Warehouse to",
            dataIndex: "warehouse_to",
            render: (_: any, record: DataSourceItem) => (
                <Select
                    style={{ width: '70%' }}
                    onChange={(value) => handleChange(value, record.key, "warehouse_to")}
                    options={Warehouseto}
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
        { value: "WT1" }, { value: "WT2" }, { value: "WT3" }, { value: "WT4" },
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
  };
  

    const handleSubmitData = () => {
        Modal.confirm({
            title: 'ยืนยันการส่งข้อมูล',
            content: 'คุณต้องการส่งข้อมูลนี้ใช่หรือไม่?',
            okText: 'ใช่',
            cancelText: 'ไม่',
            onOk: () => {
                console.log("Sending data:", dataSource);
                setDataSource([]);
                form.resetFields();
                setIsSubmitted(false);
                notification.success({
                    message: 'ส่งข้อมูล สำเร็จ',
                    description: 'ข้อมูลทั้งหมดได้ถูกส่งรียบร้อยแล้ว',
                });
            },
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
        <ConfigProvider>
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                Create IJ Return
            </div>
            <Layout>
                <Layout.Content
                    style={{
                        margin: "24px",
        padding: 36,
        minHeight: 360,
        background: "#fff",
        borderRadius: "8px",
       
        justifyContent: 'center', // Center content horizontally
        alignItems: 'center', // Center content vertically
                        
                    }}
                >
                    <div>
                    <Button
            onClick={handleBack}
            style={{ background: '#98CEFF', color: '#fff' }}
          >
            <LeftOutlined style={{ color: '#fff', marginRight: 5 }} />
            Back
          </Button>
        
                    <Form
        form={form}
        layout="vertical"
        onValuesChange={onChange}
        style={{ padding: '20px', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'center' }}
    >
        <div style={{ width: '100%', maxWidth: '800px' }}> {/* Adjust max-width here */}

        <Divider style={{color: '#657589', fontSize:'22px',marginTop:30,marginBottom:30}} orientation="left"> IJ document Information </Divider>
            <Row gutter={16} >

                <Col span={8}>
                    <Form.Item
                        label={<span style={{ color: '#657589' }}>กรอกเอกสารอ้างอิง IJ (Optional):</span>}
                        name="IJ"
                        rules={[{ required: true, message: "กรุณากรอกเอกสารอ้างอิง IJ" }]}
                    >
                        <Input style={{ width: '100%', height: '40px', }} placeholder="กรอกเอกสารอ้างอิง" />
                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item
                        label={<span style={{ color: '#657589' }}>IJ_Create: (Optional):</span>}
                        name="IJ_Create"
                    >
                        <Input style={{ width: '100%', height: '40px',}} placeholder="IJ Create" disabled />
                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item
                        label={<span style={{ color: '#657589' }}>Remark: (Optional):</span>}
                        name="Remark"
                    >
                        <Input style={{ width: '100%', height: '40px',}} showCount maxLength={200} onChange={onChange} />
                    </Form.Item>
                </Col>
                </Row>
                <Row gutter={16} >
                <Divider style={{color: '#657589', fontSize:'22px',marginTop:30,marginBottom:30}} orientation="left"> Transport Information </Divider>
                <Col span={8}>
                    <Form.Item
                        label={<span style={{ color: '#657589' }}>วันที่คืน:</span>}
                        name="Date"
                        rules={[{ required: true, message: 'กรุณาเลือกวันที่คืน' }]}
                    >
                        <DatePicker style={{ width: '100%', height: '40px',   }} placeholder="เลือกวันที่คืน" />
                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item
                        label={<span style={{ color: '#657589' }}>กรอกเลขTracking:</span>}
                        name="TrackingNumber"
                    >
                        <Input style={{ width: '100%', height: '40px', }} placeholder="เลขTracking" disabled />

                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item
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
                <Divider style={{color: '#657589', fontSize:'22px',marginTop:30,marginBottom:30}} orientation="left"> SKU Information </Divider>
                <Row gutter={16} >
                <Col span={8}>
                    <Form.Item
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
                    />
                    </Form.Item>
                </Col>
                <Col span={8}>
                    <Form.Item
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
                    />
                    </Form.Item>
                </Col>
                <Col span={4}>
                    <Form.Item
                        label={<span style={{ color: '#657589' }}>QTY:</span>}
                        name="QTY"
                        rules={[{ required: true, message: 'กรุณากรอกจำนวน' }]}
                    >
                        <InputNumber min={1} max={100} defaultValue={0} style={{ width: '100%', height: '40px', lineHeight: '40px', }} />
                       
                    </Form.Item>
                </Col>
                <Col span={4}>
               
            <Button 
            type="primary" 
            disabled={!formValid || isSubmitted}  // ปิดการใช้งานเมื่อ form ไม่ valid หรือเมื่อกด Create IJ
            onClick={handleAdd} 
            style={{ width: '100%', height: '40px', marginTop:30 }}
        >
               <PlusCircleOutlined /> {/* เพิ่มไอคอนที่นี่ */}
            Add
        </Button>
        </Col>
            </Row>
            
           
            
        </div>
    </Form>

    <div >
        <Table 
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
                  <Button 
                      type="primary" 
                      onClick={isSubmitted ? handleSubmitData : handleCreateIJ} 
                      style={{ width: 100, height: 40, marginRight: '10px' }} // เพิ่มช่องว่างทางขวา
                      disabled={!formValid || dataSource.length === 0} 
                  >
                      {isSubmitted ? "ส่งข้อมูล" : "Create IJ"}
                  </Button>
                  
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



                </Layout.Content>
            </Layout>
        </ConfigProvider>
    );
};

export default IJPage;


