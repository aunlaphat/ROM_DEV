import React, { useState } from 'react';
import { Upload, Button, Row, Col, ConfigProvider, Layout, notification, Popconfirm, Tooltip, Table } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import * as XLSX from 'xlsx';
import icon from '../../assets/images/document-text.png';

const dataSource = [
  { Order: '12345', SKU: 'SKU001', QTY: '10', Amount: '100.00', SO: 'SO001', Order_status: 'Cancel', SO_status: 'Invoice', key: '12345', SR_Create: 'NULL' },
  { Order: '12346', SKU: 'SKU002', QTY: '5', Amount: '50.00', SO: 'SO002', Order_status: 'Completed', SO_status: 'Invoice', key: '12346', SR_Create: 'NULL' },
  { Order: '12347', SKU: 'SKU002', QTY: '5', Amount: '50.00', SO: 'SO002', Order_status: 'Completed', SO_status: 'Invoice', key: '12346', SR_Create: 'NULL' },
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
  const [importedData, setImportedData] = useState<any[]>([]);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [isSRCreated, setIsSRCreated] = useState(false);

  const handleDownloadTemplate = () => {
    const templateData = [
      { Order: '' },
    ];

    const ws = XLSX.utils.json_to_sheet(templateData);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Template');
    XLSX.writeFile(wb, 'Order_Template.xlsx');
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
  
      const filteredData = transformedData.filter((row: any) => row.Order); // Filter out empty orders

      // Map imported data to dataSource only for existing orders
      const matchedData = filteredData
        .map((item) => {
          const found = dataSource.find(source => source.Order === item.Order.toString());
          return found ? { ...found } : null; // Spread found item or return null
        })
        .filter(Boolean); // Filter out null values
  
      setImportedData(matchedData);
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

  const handleDelete = () => {
    const remainingData = importedData.filter((record: any) => !selectedRowKeys.includes(record.Order));
    setImportedData(remainingData);
    setSelectedRowKeys([]);
    notification.success({ message: 'Deleted successfully' });
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
    notification.success({ message: 'SR created successfully' });
  };

  const handleConfirm = () => {
    const remainingData = importedData.filter((record: any) => !selectedRowKeys.includes(record.Order));
    setImportedData(remainingData);
    setSelectedRowKeys([]);
    setIsSRCreated(false);
    notification.success({ message: 'Data confirmed successfully' });
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
          Import Return order 
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
          <Row gutter={20} justify="end" style={{marginBottom:'30px'}}>
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

          <Table
            rowSelection={rowSelection}
            dataSource={importedData} // Use the mapped importedData
            columns={columns}
            pagination={false}
            rowKey="Order"
            style={{ width: '100%', tableLayout: 'fixed'}} 
            scroll={{ x: 'max-content' }}
          />

          <Row gutter={16} justify="center" style={{ marginTop: '20px' }}>
            <Col>
              <Popconfirm 
                title={isSRCreated ? 'Confirm data submission?' : 'Confirm SR creation?'}
                onConfirm={isSRCreated ? handleConfirm : handleCreateSR}
                okText="Yes"
                cancelText="No"
              >
                <Button type="primary" disabled={selectedRowKeys.length === 0}>
                  {isSRCreated ? 'Confirm Data' : 'Create SR'}
                </Button>
              </Popconfirm>
            </Col>
            <Col>
              <Popconfirm
                title="Confirm data reset?"
                onConfirm={() => {
                  setImportedData([]);
                  setSelectedRowKeys([]);
                  notification.success({ message: 'Data reset successfully' });
                }}
                okText="Yes"
                cancelText="No"
              >
                <Button type="default" disabled={importedData.length === 0}>
                  Reset
                </Button>
              </Popconfirm>
            </Col>
          </Row>
        </Layout.Content>
      </ConfigProvider>
    </Layout>
  );
};

export default ImportOrder;
