import { Popconfirm, Button, Col, ConfigProvider, DatePicker, Form, Input, InputNumber, Layout, Row, Select, Table, notification, Modal, Upload, Divider, Tooltip } from "antd";
import { SearchOutlined, DeleteOutlined, LeftOutlined, PlusCircleOutlined, UploadOutlined, CloseOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { useEffect, useState } from "react";
import * as XLSX from 'xlsx';
import Popup from 'reactjs-popup';
import icon from '../../assets/images/document-text.png';
const { Option } = Select;

interface Address {
  province: string;
  district: string;
  subDistrict: string;
  postalCode: string;
}

const data: Address[] = [
  { province: 'กรุงเทพมหานคร', district: 'เขตคลองเตย', subDistrict: 'แขวงคลองเตย', postalCode: '10110' },
  { province: 'กรุงเทพมหานคร', district: 'เขตคลองเตย', subDistrict: 'แขวงคลองตัน', postalCode: '10110' },
  { province: 'กรุงเทพมหานคร', district: 'เขตคลองเตย', subDistrict: 'แขวงพระโขนง', postalCode: '	10110' },
  { province: 'เชียงใหม่', district: 'เมืองเชียงใหม่', subDistrict: 'สุเทพ', postalCode: '50200' },
  { province: 'เชียงใหม่', district: 'เมืองเชียงใหม่', subDistrict: 'ศรีภูมิ', postalCode: '50200' },
  { province: 'เชียงใหม่', district: 'เมืองเชียงใหม่', subDistrict: 'ช้างเผือก', postalCode: '50300' },
  // เพิ่มข้อมูลเพิ่มเติมตามต้องการ
];
const SKUName = [
  { Name: "Bewell Better Back 2 Size M Nodel H01 (Gray)", SKU: "G097171-ARM01-BL" },
  { Name: "Bewell Sport armband size M For", SKU: "G097171-ARM01-GY" },
  { Name: "Sport armband size L", SKU: "G097171-ARM02-BL" },
  { Name: "Bewell Sport armband size M with light", SKU: "G097171-ARM03-GR" },
];
interface Customer {
  Key: number;
  Customer_account: string;
  Invoice_name: string;
  Address: string;
  Tax: string;
  Customer_name: string;
}
interface DataItem {
  key: number;
  SKU: string; // หรือประเภทอื่น ๆ ที่คุณต้องการ
  Name: string;
  QTY: number;
}

const Customeraccount = [
  {
    Key: 1,
    Customer_account: "DC-NMI-0033",
    Customer_name: "ปิยวลี ศรีสุวรรณ",
    Invoice_name: "ปิยะวลี",
    Address: '990 อาคารอับดุลราฮิม ชั้น 18-19 ถ.พระราม 4 สีลม บางรัก กรุงเทพมหานคร 10500',
    Tax: '-',

  },
  {
    Key: 2,
    Customer_account: "DC-NMI-0034",
    Customer_name: "สมชาย สินทรัพย์",
    Invoice_name: "สมชาย",
    Address: '500 อาคารกลาง ถนนสีลม บางรัก กรุงเทพมหานคร 10500',
    Tax: '-',
  },
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


const CreateTradeReturn = () => {
  const [isSaving, setIsSaving] = useState(false);
  const [invoiceAddress, setInvoiceAddress] = useState('');
  const [open, setOpen] = useState(false);
  const [selectedSKU, setSelectedSKU] = useState<string | undefined>(undefined);
  const [selectedName, setSelectedName] = useState<string | undefined>(undefined);
  const [form] = Form.useForm();
  const [formValid, setFormValid] = useState(false);
  const [formaddress] = Form.useForm();
  const [selectedAccount, setSelectedAccount] = useState<Customer | null>(null);
  const [dataSource, setDataSource] = useState<DataItem[]>([]);
  const [price, setPrice] = useState<number | null>(null); // Allow null
  const [qty, setQty] = useState<number | null>(null);  // Allow null
  const [isInvoiceEnabled, setIsInvoiceEnabled] = useState(false);
  const [province, setProvince] = useState<string | undefined>(undefined);
  const [district, setDistrict] = useState<string | undefined>(undefined);
  const [subDistrict, setSubDistrict] = useState<string | undefined>(undefined);
  const [postalCode, setPostalCode] = useState<string | undefined>(undefined);

  const handleProvinceChange = (value: string) => {
    setProvince(value);
    setDistrict(undefined);
    setSubDistrict(undefined);
    setPostalCode(undefined);
  };


  const handleDistrictChange = (value: string) => {
    setDistrict(value);
    setSubDistrict(undefined);
    setPostalCode(undefined);
  };

  const handleSubDistrictChange = (value: string) => {
    setSubDistrict(value);
  };

  // Automatically set postal code when sub-district changes
  useEffect(() => {
    if (subDistrict) {
      const selected = data.find(item => item.subDistrict === subDistrict);
      setPostalCode(selected?.postalCode || undefined);
    }
  }, [subDistrict]);



  const provinces = Array.from(new Set(data.map((item) => item.province)));
  const districts = Array.from(new Set(data.filter((item) => item.province === province).map((item) => item.district)));
  const subDistricts = Array.from(new Set(data.filter((item) => item.district === district).map((item) => item.subDistrict)));



  const handleAccountChange = (value: string) => {
    const selectedCustomer = Customeraccount.find(account => account.Customer_account === value);
    setSelectedAccount(selectedCustomer || null);
    setIsInvoiceEnabled(!!selectedCustomer);

    form.setFieldsValue({
      Invoice_name: selectedCustomer?.Invoice_name || '',
      Address: selectedCustomer?.Address || '',
      Customer_name: selectedCustomer?.Customer_name || '',
      Tax: selectedCustomer?.Tax || '',
    });
  };

  const handleInvoiceChange = (value: string) => {
    const selectedCustomer = Customeraccount.find(account => account.Invoice_name === value);
    setSelectedAccount(selectedCustomer || null);

    form.setFieldsValue({
      Address: selectedCustomer?.Address || '',

    });
  };

  const handleOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
    // ให้ฟอร์มรีเซ็ตเฉพาะในกรณีที่ไม่ได้กดบันทึก

  };


  const handleSelectChange = (value: any) => {
    // เมื่อเลือกจังหวัดแล้วปิด Popup
    // setOpen(false);
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
  const handleSubmit = () => {
    if (dataSource.length === 0) {
      notification.warning({
        message: 'ไม่สามารถส่งข้อมูลได้',
        description: 'กรุณาเพิ่มข้อมูลในตารางก่อนส่ง!',
      });
      return; // หยุดการทำงานของฟังก์ชัน
    }

    // ส่งข้อมูลและรีเซ็ตฟอร์มและตาราง
    console.log('Table Data:', dataSource);

    // รีเซ็ตฟอร์มและตาราง
    form.resetFields();
    setDataSource([]); // หรือปรับเป็นค่าเริ่มต้นที่คุณต้องการได้

    notification.success({
      message: 'ส่งข้อมูลสำเร็จ',
      description: 'ข้อมูลของคุณถูกส่งเรียบร้อยแล้ว!',
    });
  };

  const handleSelectAddress = (address: string) => {
    setInvoiceAddress(address); // ตั้งค่าที่อยู่ที่เลือก
    form.setFieldsValue({ Invoice_address: address }); // ตั้งค่าให้ฟอร์ม
    // setOpen(false); // ปิด Popup
  };


  const handleSave = async () => {
    try {
      const values = await formaddress.validateFields();
      console.log('Form Values:', values);

      // Update form values in the main form
      form.setFieldsValue({
        Invoice_name: values.Invoice_name, // Update the Invoice name
        Address: values.AddressNew + " " + values.SubDistrict + " " + values.district + " " + values.province + " " + values.PostalCode,
      });

      setIsSaving(true);
      setSelectedAccount(values); // Save the address data if needed
      setIsSaving(false);

      // Reset modal form fields
      formaddress.resetFields(); // Reset the fields in the modal

      handleClose(); // Close modal after save

    } catch (error) {
      console.error('Failed to save:', error);
    }
  };


  const onSearch = (value: string) => {
    console.log('search:', value);
  };


  const columns = [
    { title: 'SKU', dataIndex: 'SKU', key: 'SKU' ,id:'SKU' },
    { title: 'Name', dataIndex: 'Name', key: 'Name',id:'Name'  },
    { title: 'QTY', dataIndex: 'QTY', key: 'QTY',id:'QTY'  },
    {
      title: 'Action',
      id:'Action',
      dataIndex: 'Action',
      key: 'Action',
      render: (_: any, record: { key: number }) => (
        <Popconfirm
          title="คุณแน่ใจหรือไม่ว่าต้องการลบข้อมูลนี้?"
          onConfirm={() => handleDelete(record.key)} // เรียกใช้ฟังก์ชัน handleDelete เมื่อกดยืนยัน
          okText="ใช่"
          cancelText="ไม่"
        >
          <DeleteOutlined
            style={{ cursor: 'pointer', color: 'red', fontSize: '20px' }}
          />
        </Popconfirm>
      )
    },

  ];
  const handleDownloadTemplate = () => {
    const templateColumns = columns.filter(col => col.key !== 'Action'); // กรองออก action column
    const ws = XLSX.utils.json_to_sheet([]);
    XLSX.utils.sheet_add_aoa(ws, [templateColumns.map(col => col.title)]);

    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Template');

    XLSX.writeFile(wb, 'Template.xlsx');
  };

  const handleUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const data = new Uint8Array(e.target?.result as ArrayBuffer);
      const workbook = XLSX.read(data, { type: 'array' });
      const worksheet = workbook.Sheets[workbook.SheetNames[0]];
      const json = XLSX.utils.sheet_to_json<DataItem>(worksheet);

      // กรองข้อมูลเฉพาะที่มี SKU และ QTY
      const filteredData = json.filter(item => item.SKU && item.QTY);

      // อัปเดต dataSource ด้วยข้อมูลที่กรอง
      setDataSource(filteredData);

      notification.success({
        message: 'อัปโหลดสำเร็จ',
        description: 'ข้อมูลจากไฟล์ Excel ถูกนำเข้าเรียบร้อยแล้ว!',
      });
    };
    reader.readAsArrayBuffer(file);
  };

  const uploadProps = {
    beforeUpload: (file: File) => {
      handleUpload(file);
      return false; // ป้องกันไม่ให้ Ant Design ทำการอัปโหลด
    },
  };

  const handleAdd = () => {
    form.validateFields()
      .then(values => {
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
        const newData = {
          key: dataSource.length + 1,
          SKU: values.SKU,
          Name: values.SKU_Name,
          QTY: values.QTY,
          Price: values.Price,
          Tax: values.Tax || '', // ถ้าไม่มี Tax ให้ใส่ค่าว่าง
        };

        setDataSource([...dataSource, newData]); // เพิ่มข้อมูลใหม่ไปยัง dataSource

        // แสดงข้อความเมื่อเพิ่มข้อมูลสำเร็จ
        notification.success({
          message: "เพิ่มสำเร็จ",
          description: "ข้อมูลของคุณถูกเพิ่มเรียบร้อยแล้ว!",
        });

        // ล้างฟิลด์ในฟอร์มหลังจากเพิ่มข้อมูลเสร็จ
        form.resetFields(['SKU', 'SKU_Name', 'QTY', 'Price']);
      })
      .catch(info => {
        console.log('Validate Failed:', info);
        notification.warning({
          message: "มีข้อสงสัย",
          description: "กรุณากรอกข้อมูลให้ครบก่อนเพิ่ม!",
        });
      });
  };


  const handleDelete = (key: number) => {
    setDataSource(dataSource.filter(item => item.key !== key));
    notification.success({
      message: 'ลบข้อมูลสำเร็จ',
      description: 'ข้อมูลของคุณถูกลบออกเรียบร้อยแล้ว.',
    });
  };


  const onChange = () => {
    const values = form.getFieldsValue();
    const { Date, SKU, QTY } = values;

    // Set form validity based on required fields
    setFormValid(Date && SKU && QTY);
  };


  return (
    <ConfigProvider>
      <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
        Create Trade Return
      </div>
      <Layout>
        <Layout.Content

          style={{
            margin: "24px",
            padding: 36,
            minHeight: 360,
            background: "#fff",
            borderRadius: "8px",



          }}
        >
          <Form form={form} layout="vertical" style={{ width: '100%', padding: '30px' }}>
            <div>
              <Divider style={{ color: '#657589', fontSize: '22px', margin: 30 }} orientation="left">Sale Order Information</Divider>
              <Row gutter={16}  >
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
                    <Input style={{ height: 40 }} />
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

                    />
                  </Form.Item>
                </Col>
                <Col span={8}></Col>
              </Row>

              <Row gutter={16} align="middle" justify="center" style={{ marginTop: '20px', width: '100%' }}>
                <Col span={8}>
                  <Form.Item
                    id="Doc"
                    label={
                      <span style={{ color: '#657589' }}>
                        กรอกเอกสารอ้างอิง:&nbsp;
                        <Tooltip title="ตัวอย่างเอกสาร SOA2410-00234">
                          <QuestionCircleOutlined style={{ color: '#657589' }} />
                        </Tooltip>
                      </span>
                    } name="Doc">
                    <Input style={{ width: '100%', height: '40px' }} placeholder="ตัวอย่างเอกสาร SOA2410-00234" />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    id="CustomerAccount"
                    label={<span style={{ color: '#657589' }}>Customer account:</span>} name="Customer_account" rules={[{ required: true, message: 'Customer account is required' }]}>
                    <Select style={{ width: '100%', height: '40px', borderWidth: '1px' }} showSearch placeholder="Customer account" optionFilterProp="children" onChange={handleAccountChange}>
                      {Customeraccount.map(account => (
                        <Select.Option key={account.Key} value={account.Customer_account}>
                          {account.Customer_account}
                        </Select.Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    id="CustomerName"
                    label={<span style={{ color: '#657589' }}>Customer name:</span>} name="Customer_name" rules={[{ required: true, message: 'Customer name is required' }]}>
                    <Input style={{ width: '100%', height: '40px', borderWidth: '1px' }} placeholder="Customer name" value={selectedAccount?.Customer_name} disabled />
                  </Form.Item>
                </Col>
              </Row>

              <Row gutter={16} style={{ marginTop: '10px' }}>
                <Col span={8}>
                  <Form.Item
                    id="InvoiceName"
                    label={<span style={{ color: '#657589' }}>Invoice name:</span>} name="Invoice_name" rules={[{ required: true, message: 'Invoice name is required' }]}>
                    <Select style={{ width: '100%', height: '40px', borderWidth: '1px' }} showSearch placeholder="Invoice name" optionFilterProp="children" disabled={!isInvoiceEnabled} onChange={handleInvoiceChange}>
                      {Customeraccount.map(account => (
                        <Option key={account.Key} value={account.Invoice_name}>
                          {account.Invoice_name}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>

                <Col span={8}>
                  <Form.Item
                    id="Taxid"
                    label={<span style={{ color: '#657589' }}>Tax ID(เลขผู้เสียภาษี):</span>} name="Tax" rules={[{ required: selectedAccount ? true : false, message: 'Tax ID is required' }]}>
                    <Input style={{ width: '100%', height: '40px', borderWidth: '1px' }} placeholder="Tax ID" value={selectedAccount?.Tax} disabled />
                  </Form.Item>
                </Col>
              </Row>

              <Divider style={{ color: '#657589', fontSize: '22px', margin: 30 }} orientation="left">Address Information</Divider>
              <Row gutter={16} style={{ marginTop: '10px' }}>
                <Col span={18}>
                  <Form.Item
                    id="InvoiceAddress"
                    label={<span style={{ color: '#657589' }}>Invoice address:</span>} name="Address" rules={[{ required: true, message: 'Invoice address is required' }]}>
                    <Select style={{ width: '100%', height: '40px', borderWidth: '1px' }} suffixIcon={null} disabled value={selectedAccount ? selectedAccount.Address : undefined}>
                      {selectedAccount ? (
                        <Option value={selectedAccount.Address}>{selectedAccount.Address}</Option>
                      ) : (
                        <Option value="">ไม่มีข้อมูล</Option>
                      )}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={6}>
                  <Button
                    id="NewInvoiceAddress"
                    type="primary" onClick={handleOpen} style={{ width: '100%', height: '40px', marginTop: 30 }}>New invoice address</Button>
                </Col>
              </Row>

              <Divider style={{ color: '#657589', fontSize: '22px', margin: 30 }} orientation="left"> SKU information</Divider>
              <Row gutter={16} style={{ marginTop: '10px', width: '100%' }}>
                <Col span={6}>
                  <Form.Item
                    id="Sku"
                    label={<span style={{ color: '#657589' }}>กรอก SKU :</span>} name="SKU" rules={[{ required: true, message: "กรุณากรอก SKU" }]}>
                    <Select showSearch style={{ width: '100%', height: '40px' }} placeholder="Search to Select" optionFilterProp="label" value={selectedSKU} onChange={handleSKUChange} options={skuOptions} dropdownStyle={{ minWidth: 200 }} />
                  </Form.Item>
                </Col>

                <Col span={7}>
                  <Form.Item
                    id="Skuname"
                    label={<span style={{ color: '#657589' }}>กรอก SKU Name:</span>} name="SKU_Name" rules={[{ required: true, message: "กรุณาเลือก SKU Name" }]}>
                    <Select showSearch style={{ width: '100%', height: '40px' }} placeholder="Search to Select" optionFilterProp="label" value={selectedName} onChange={handleNameChange} options={nameOptions} dropdownStyle={{ minWidth: 300 }} />
                  </Form.Item>
                </Col>

                <Col span={4}>
                  <Form.Item
                    id="qty"
                    label={<span style={{ color: '#657589' }}>QTY:</span>} name="QTY" rules={[{ required: true, message: "กรุณากรอก QTY" }]}>
                    <InputNumber min={1} max={100} value={qty} onChange={(value) => setQty(value)} style={{ width: '100%', height: '40px', lineHeight: '40px' }} />
                  </Form.Item>
                </Col>

                <Col span={4}>
                  <Form.Item
                    id="price"
                    label={<span style={{ color: '#657589' }}>Price:</span>} name="Price" rules={[{ required: true, message: "กรุณากรอก Price" }]}>
                    <InputNumber min={1} max={100000} value={price} onChange={(value) => setPrice(value)} step={0.01} style={{ width: '100%', height: '40px', lineHeight: '40px' }} />
                  </Form.Item>
                </Col>

                <Col span={3}>

                  <Button
                    id="add"
                    type="primary"
                    style={{ width: '100%', height: '40px', marginTop: 30 }}
                    onClick={handleAdd} // เรียกใช้ฟังก์ชัน handleAdd
                  >
                    <PlusCircleOutlined />
                    Add
                  </Button>
                </Col>
              </Row>
            </div>
          </Form>
          <Modal
            open={open}
            onClose={handleClose}
            closeIcon={false}
            footer={null}
          >
            <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '10px', overflow: "auto", }}>
              <Button id="Closeicon" type="text" onClick={handleClose} icon={<CloseOutlined style={{ fontSize: '24px' }} />} danger />
            </div>
            <div style={{ fontSize: '20px', color: '#35465B', }}>
              New Invoice Address
            </div>
            <Form
              form={formaddress}
              layout="vertical"
              style={{ width: '100%', display: 'flex', padding: 20 }}
              onFinish={handleSave}
            >
              <Row gutter={16} style={{ marginTop: '10px', justifyContent: 'center' }}>

                <Col >
                  <Form.Item
                    id="Invoicename"
                    label={<span style={{ color: '#657589' }}>Invoice name:</span>}
                    name="Invoice_name"
                    rules={[{ required: true, message: "กรอก Invoice" }]}
                  >
                    <Input
                      style={{ width: '400px', height: '40px' }}
                      placeholder="Invoice name"
                      value={selectedAccount?.Invoice_name} // แสดงค่า Tax ID ถ้ามี
                      disabled={!selectedAccount} // ปิดการใช้งานถ้าไม่มีลูกค้าที่เลือก
                    />
                  </Form.Item>
                </Col>


                <Col>
                  <Form.Item
                    id="addressnew"
                    label={<span style={{ color: '#657589' }}>บ้านเลขที่:</span>}
                    name="AddressNew"
                    rules={[{ required: true, message: "กรอก บ้านเลขที่" }]}
                  >
                    <Input style={{ width: '400px', height: '40px' }} placeholder="กรอกบ้านเลขที่" />
                  </Form.Item>
                </Col>

                {/* Province */}
                <Col>
                  <Form.Item
                    id="SelectProvince"
                    label={<span style={{ color: '#657589' }}>จังหวัด:</span>} name="province" rules={[{ required: true, message: 'เลือกจังหวัด' }]}>
                    <Select placeholder="Select Province" onChange={handleProvinceChange} style={{ width: '400px', height: '40px' }}>
                      {provinces.map((item) => (
                        <Option key={item} value={item}>{item}</Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>

                {/* District */}
                <Col>
                  <Form.Item
                    id="SelectDistrict"
                    label={<span style={{ color: '#657589' }}>เขต:</span>} name="district" rules={[{ required: true, message: 'เลือกเขต' }]}>
                    <Select placeholder="Select District" onChange={handleDistrictChange} style={{ width: '400px', height: '40px' }}>
                      {districts.map((item) => (
                        <Option key={item} value={item}>{item}</Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>

                {/* SubDistrict */}
                <Col>
                  <Form.Item
                    id="SelectSubDistrict"
                    label={<span style={{ color: '#657589' }}>แขวง:</span>} name="SubDistrict" rules={[{ required: true, message: 'เลือกแขวง' }]}>
                    <Select placeholder="Select SubDistrict" onChange={handleSubDistrictChange} style={{ width: '400px', height: '40px' }}>
                      {subDistricts.map((item) => (
                        <Option key={item} value={item}>{item}</Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>

                {/* Postal Code */}
                <Col>
                  <Form.Item
                    id="PostalCode"
                    label={<span style={{ color: '#657589' }}>รหัสไปรษณีย์:</span>} name="PostalCode" rules={[{ required: true, message: 'กรุณาระบุรหัสไปรษณีย์' }]}>
                    <Select placeholder="Postal Code" value={postalCode} style={{ width: '400px', height: '40px' }}>
                      {postalCode && <Option key={postalCode} value={postalCode}>{postalCode}</Option>}
                    </Select>
                  </Form.Item>
                </Col>

                {/* Save Button */}
                <Col>

                  <Button
                    id="save"
                    type="primary"
                    htmlType="submit"
                    disabled={isSaving}
                  >
                    Save
                  </Button>

                </Col>
              </Row>
            </Form>
          </Modal>



          <Row gutter={20} style={{ marginBottom: 20, marginLeft: 20 }}>
            <Col  >
              <Button 
              id=" Download Template"
              onClick={handleDownloadTemplate}>
                <img src={icon} alt="Download Icon" style={{ width: 16, height: 16, marginRight: 8 }} />
                Download Template
              </Button>


            </Col>

            <Col>
              <Upload {...uploadProps} showUploadList={false}>
                <Button id=" Import Excel" icon={<UploadOutlined />} style={{ background: '#7161EF', color: '#FFF', marginBottom: 10 }}>
                  Import Excel
                </Button>
              </Upload>
            </Col>

          </Row>
          <div>
            <Table
              dataSource={dataSource}
              columns={columns}
              rowKey="key"
              pagination={false} // Disable pagination if necessary
              style={{ width: '100%', tableLayout: 'fixed' }} // Ensure the table takes full width and is fixed layout
              scroll={{ x: 'max-content' }}

            />
          </div>
          <Row justify="center" gutter={16}>
            <Popconfirm
            id="popconfirmSubmit"
              title="คุณแน่ใจหรือไม่ว่าต้องการส่งข้อมูล?"
              onConfirm={handleSubmit} // เรียกใช้ฟังก์ชัน handleSubmit เมื่อกดยืนยัน
              okText="ใช่"
              cancelText="ไม่"
            >
              <Button
              id="Submit"
                style={{ color: '#fff', backgroundColor: '#14C11B', width: 100, height: 40, margin: 20 }}
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

export default CreateTradeReturn;
