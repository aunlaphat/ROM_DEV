import React, { useState } from 'react';
import { Upload, Button, Row, Col, ConfigProvider, Layout, Select, Form, Table, message, Modal, Tooltip } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import * as XLSX from 'xlsx';
import icon from '../../assets/images/document-text.png';

const options = [
  { value: '12345', label: '12345' },
  { value: '12346', label: '12346' },
  { value: '12347', label: '12347' },
];

const dataSource = [
  { Order: '12345', SKU: 'SKU001', QTY: '10', Amount: '100.00', SO: 'SO001', Order_status: 'Cancel', SO_status: 'Invoice', key: '12345', SR_Create: 'NULL' },
  { Order: '12346', SKU: 'SKU002', QTY: '5', Amount: '50.00', SO: 'SO002', Order_status: 'Completed', SO_status: 'Invoice', key: '12346', SR_Create: 'NULL' },
];

const columns = [
  { title: 'Order', dataIndex: 'Order', key: 'Order' },
  { title: 'SKU', dataIndex: 'SKU', key: 'SKU' },
  { title: 'QTY', dataIndex: 'QTY', key: 'QTY' },
  { title: 'Amount', dataIndex: 'Amount', key: 'Amount' },
  { title: 'SO', dataIndex: 'SO', key: 'SO' },
  {
    title: 'Order Status',
    dataIndex: 'Order_status',
    key: 'Order_status',
    render: (text: string) => (
      <Button
        type={text === 'Cancel' ? 'primary' : 'default'}
        style={{
          backgroundColor: text === 'Cancel' ? '#FDCACA' : undefined,
          borderRadius: '20px',
          borderColor: 'transparent',
          color: text === 'Cancel' ? '#FC6F6F' : undefined,
        }}
      >
        {text === 'Cancel' ? 'Cancel' : text}
      </Button>
    ),
  },
  {
    title: 'SO Status',
    dataIndex: 'SO_status',
    key: 'SO_status',
    render: (text: string) => (
      <Button
        type="primary"
        style={{
          backgroundColor: '#E9F3FE',
          color: '#657589',
          borderColor: 'transparent',
          borderRadius: '20px',
        }}
      >
        {text}
      </Button>
    ),
  },
  { title: 'SR_Create', dataIndex: 'SR_Create', key: 'SR_Create' }
];

const ImportOrder = () => {
  const [form] = Form.useForm();
  const [selectedValue, setSelectedValue] = useState<string | undefined>();
  const [filteredData, setFilteredData] = useState<any[]>([]);
  const [importedData, setImportedData] = useState<any[]>([]);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [isSRCreated, setIsSRCreated] = useState(false); // เพิ่มการประกาศตัวแปรที่นี่

  const handleSelectChange = (value: string) => {
    setSelectedValue(value);
  };

  const handleCheck = () => {
    if (selectedValue) {
      const dataToShow = dataSource.filter(item => item.Order === selectedValue);
      setImportedData(dataToShow);
    }
  };

  const handleDownloadTemplate = () => {
    const ws = XLSX.utils.json_to_sheet([]);
    XLSX.utils.sheet_add_aoa(ws, [columns.map(col => col.title)]);

    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Template');

    XLSX.writeFile(wb, 'Template.xlsx');
  };

  const handleImportExcel = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const binaryStr = e.target?.result;
      const workbook = XLSX.read(binaryStr, { type: 'binary' });
      const firstSheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[firstSheetName];
      const importedData = XLSX.utils.sheet_to_json(worksheet);

      const transformedData = importedData.map((row: any) => ({
        Order: row['Order'] || row['Order ID'],
        SKU: row['SKU'],
        QTY: row['QTY'],
        Amount: row['Amount'],
        SO: row['SO'],
        Order_status: row['Order Status'] || row['Order Status'],
        SO_status: row['SO Status'] || row['SO Status'],
        SR_Create: row['SR_Create'] || row['SR_Create'],
      }));

      const filteredData = transformedData.filter((row: any) => Object.values(row).some(value => value !== '' && value !== undefined));
      setImportedData(filteredData);
    };
    
    reader.readAsBinaryString(file);
  };

  const uploadProps = {
    beforeUpload: (file: File) => {
      handleImportExcel(file);
      return false;
    },
    showUploadList: false,
  };

  const showDeleteConfirm = () => {
    Modal.confirm({
      title: 'ยืนยันการลบ',
      content: 'คุณแน่ใจว่าต้องการลบข้อมูลที่เลือก?',
      onOk: handleDelete,
      onCancel: () => {},
    });
  };

  const handleDelete = () => {
    const selectedData = importedData.filter((record: any) => selectedRowKeys.includes(record.Order));
    console.log('ข้อมูลที่ลบ:', selectedData);
  
    const remainingData = importedData.filter((record: any) => !selectedRowKeys.includes(record.Order));
    setImportedData(remainingData);
    setSelectedRowKeys([]);
    message.success('ลบข้อมูลสำเร็จ');
  };
  
  const showConfirmModal = () => {
    if (isSRCreated) {
      // แสดงข้อความยืนยันการส่งข้อมูล
      Modal.confirm({
        title: 'ยืนยันการส่งข้อมูล',
        content: 'คุณแน่ใจว่าต้องการยืนยันการส่งข้อมูล?',
        onOk: handleConfirm, // ถ้ายืนยันให้เรียกใช้ handleConfirm
        onCancel: () => {},
      });
    } else {
      // ถ้ายังไม่มีการสร้าง SR ให้แสดงข้อความยืนยันการสร้าง SR
      Modal.confirm({
        title: 'ยืนยันการสร้าง SR',
        content: 'คุณแน่ใจว่าต้องการส่งข้อมูลสร้าง SR?',
        onOk: handleCreateSR,
        onCancel: () => {},
      });
    }
  };
  
  const handleCreateSR = () => {
    const updatedData = importedData.map((record: any) => {
      if (selectedRowKeys.includes(record.Order)) {
        const randomNum = Math.floor(1000 + Math.random() * 9000);
        return {
          ...record,
          SR_Create: `SRA2409-${randomNum}`,
        };
      }
      return record;
    });
    setImportedData(updatedData);
    setIsSRCreated(true);
    message.success('SR สร้างสำเร็จ');
  };

  const handleConfirm = () => {
    const selectedData = importedData.filter((record: any) => selectedRowKeys.includes(record.Order));
    console.log('ข้อมูลที่มี SR:', selectedData);
    
    // ลบข้อมูลที่เลือกออกจาก importedData
    const remainingData = importedData.filter((record: any) => !selectedRowKeys.includes(record.Order));
    setImportedData(remainingData);
    
    setSelectedRowKeys([]); // รีเซ็ต selectedRowKeys หลังจากยืนยัน
    setIsSRCreated(false); // รีเซ็ต isSRCreated เพื่อให้ปุ่มกลับเป็น "สร้าง SR"
    message.success('ยืนยันข้อมูลสำเร็จ');
  };

  const rowSelection = {
    selectedRowKeys,
    onChange: (newSelectedRowKeys: React.Key[]) => {
      setSelectedRowKeys(newSelectedRowKeys);
    },
  };

  return (
    <Layout>
      <ConfigProvider>
        <div style={{ marginLeft: '28px', fontSize: '25px', fontWeight: 'bold', color: 'DodgerBlue' }}>
          Import Order
        </div>
        <Layout.Content
          style={{
            margin: "24px",
            padding: 36,
            minHeight: 360,
            background: "#fff",
            borderRadius: "8px",
          }}
        >
         <Row gutter={20} justify="end">
  <Col>
    <Tooltip title="Download the template file from table">
      <Button onClick={handleDownloadTemplate}>
        <img src={icon} alt="Download Icon" style={{ width: 16, height: 16, marginRight: 8 }} />
        Download Template
      </Button>
    </Tooltip>
  </Col>
  <Col>
    <Tooltip title="Import Excel file">
      <Upload {...uploadProps}>
        <Button icon={<UploadOutlined />} style={{ background: '#7161EF', color: '#FFF' }}>
          Import Excel
        </Button>
      </Upload>
    </Tooltip>
  </Col>
</Row>


          <Form layout="vertical" form={form} style={{ marginTop: '40px' }}>
            <Row justify="center" align="middle" >
              <Col span={8}>
                <Form.Item
                  label="กรอกเลข Order"
                  name="กรอกเลข Order"
                  rules={[{ required: true, message: 'Please select the SKU!' }]}
                >
                  <Row gutter={16}>
                    <Col span={22}>
                      <Select
                        showSearch
                        style={{maxWidth: '500px', height: '50px'}}
                        placeholder="กรอกเลข Order"
                        optionFilterProp="label"
                        value={selectedValue}
                        onChange={handleSelectChange}
                        options={options}
                      />
                    </Col>
                    <Col span={2}>
                    <Tooltip title="ตรวจสอบเลขOrder">
                      <Button type="primary" style={{ height: '50px', width: '100px'                      }} onClick={handleCheck}>
                        ตรวจสอบ
                      </Button>
                      </Tooltip>
                    </Col>
                  </Row>
                </Form.Item>
              </Col>
            </Row>
          </Form>

          <Table
            rowSelection={rowSelection}
            dataSource={importedData}
            columns={columns}
            pagination={false}
            rowKey="Order"
            style={{ width: '100%', tableLayout: 'fixed' }} // Ensure the table takes full width and is fixed layout
            scroll={{ x: 'max-content' }}
          />

<Row gutter={16} justify="center" style={{ marginTop: '20px' }}>
  <Col>
    <Button type="primary" onClick={showConfirmModal} disabled={selectedRowKeys.length === 0}>
      {isSRCreated ? 'ยืนยันข้อมูล' : 'สร้าง SR'}
    </Button>
  </Col>
  <Col>
    <Button danger onClick={showDeleteConfirm} disabled={selectedRowKeys.length === 0}>
      ลบ
    </Button>
  </Col>
</Row>
        </Layout.Content>
      </ConfigProvider>
    </Layout>
  );
};

export default ImportOrder;

