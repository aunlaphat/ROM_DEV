import React, { useState, Key } from 'react';
import { ConfigProvider, Layout, Row, Col, Select, Form, Button, Table, Empty, Modal, notification } from 'antd';
import type { TableColumnsType } from 'antd';
import { DeleteOutlined } from '@ant-design/icons';

const options = [
  { value: '1', label: 'Not Identified' },
  { value: '2', label: 'G097171-ARM01-BL Bewell Sport armband size M For' },
  { value: '3', label: 'G097171-ARM01-GY Bewell Sport armband size M For' },
  { value: '4', label: 'G097171-ARM02-BL Bewell Sport armband size L' },
  { value: '5', label: 'G097171-ARM03-GR Bewell Sport armband size M with light' },
  { value: '6', label: 'Cancelled' },
];

const shopOptions = [
  { value: '1', label: 'Shop A' },
  { value: '2', label: 'Shop B' },
  { value: '3', label: 'Shop C' },
];

const itemOptions = [
  { value: '1', label: 'Item 1' },
  { value: '2', label: 'Item 2' },
  { value: '3', label: 'Item 3' },
];

const variationOptions = [
  { value: '1', label: 'Variation 1' },
  { value: '2', label: 'Variation 2' },
  { value: '3', label: 'Variation 3' },
];

interface DataType {
  key: Key;
  searchSkuProductName: string;
  shopName: string;
  itemId: string;
  variationSku: string;
}

const Sync = () => {
  const [form] = Form.useForm();
  const [selectedValue, setSelectedValue] = useState<string | undefined>();
  const [shopName, setShopName] = useState<string | undefined>();
  const [itemId, setItemId] = useState<string | undefined>();
  const [variationSku, setVariationSku] = useState<string | undefined>();
  const [tableData, setTableData] = useState<DataType[]>([]);

  const handleSelectChange = (value: string) => {
    const selectedOption = options.find(option => option.value === value);
    if (selectedOption) {
      setSelectedValue(selectedOption.label);
    }
  };

  const handleShopNameChange = (value: string) => {
    const selectedOption = shopOptions.find(option => option.value === value);
    if (selectedOption) {
      setShopName(selectedOption.label);
    }
  };

  const handleItemIdChange = (value: string) => {
    const selectedOption = itemOptions.find(option => option.value === value);
    if (selectedOption) {
      setItemId(selectedOption.label);
    }
  };

  const handleVariationSkuChange = (value: string) => {
    const selectedOption = variationOptions.find(option => option.value === value);
    if (selectedOption) {
      setVariationSku(selectedOption.label);
    }
  };

  const handleAdd = () => {
    form.validateFields().then(() => {
      const newData: DataType = {
        key: Date.now(),
        searchSkuProductName: selectedValue || '',
        shopName: shopName || '',
        itemId: itemId || '',
        variationSku: variationSku || '',
      };
  
      setTableData([...tableData, newData]);
      form.resetFields();
      setSelectedValue(undefined);
      setShopName(undefined);
      setItemId(undefined);
      setVariationSku(undefined);
  
      // Show success notification
      notification.success({
        message: 'Add Successfully',
        description: 'The new item has been added to the table.',
      });
    }).catch(info => {
      // Display error notification if validation fails
      notification.error({
        message: 'Validation Failed',
        description: 'Please ensure all required fields are selected.',
      });
    });
  };
  

  const handleSubmit = () => {
    if (tableData.length === 0) {
      notification.error({
        message: 'No Data',
        description: 'Please add some data before submitting.',
      });
      return;
    }
  
    Modal.confirm({
      title: 'Confirm Submission',
      content: 'Are you sure you want to submit? This will clear all data.',
      okText: 'Yes',
      cancelText: 'No',
      onOk: () => {
        console.log('Form submitted with data:', tableData);
        setTableData([]); // Clear the table data
  
        // Show success notification
        notification.success({
          message: 'Data submitted successfully.',
          description: 'Your data has been submitted and cleared.',
        });
      },
    });
  };
  

  const handleDelete = (key: Key) => {
    setTableData(prevData => prevData.filter(item => item.key !== key));
  };

const columns: TableColumnsType<DataType> = [
  {
    title: 'Search SKU Product Name',
    dataIndex: 'searchSkuProductName',
  },
  {
    title: 'Shop Name',
    dataIndex: 'shopName',
  },
  {
    title: 'Item ID',
    dataIndex: 'itemId',
  },
  {
    title: 'Variation SKU',
    dataIndex: 'variationSku',
  },
  {
    title: 'Action',
    key: 'action',
    render: (_, record: DataType) => (
      <Button
        type="text"
        icon={<DeleteOutlined />}
        danger
        onClick={() => handleDelete(record.key)}
      />
    ),
  },
];

  return (
    <ConfigProvider>
      <div style={{ marginLeft: '28px', fontSize: '25px', fontWeight: 'bold', color: 'DodgerBlue' }}>
        Sync MKP/Shopee/Dcom
      </div>
      <Layout.Content
        style={{
          margin: '24px',
          padding: 36,
          minHeight: 360,
          background: '#fff',
          borderRadius: '8px',
        }}
      >
        <Form
          layout="vertical"
          form={form}
          style={{ width: '100%', marginTop: '40px' }}
        >
          <Row justify="center" align="middle" style={{ width: '100%' }}>
            <Col>
              <Form.Item
                label="Search Sku Product Name"
                name="searchSkuProductName"
                rules={[{ required: true, message: 'Please select the SKU!' }]}
              >
                <Select
                  showSearch
                  style={{ width: 550, height: 50 }}
                  placeholder="Search to Select"
                  optionFilterProp="label"
                  value={selectedValue}
                  onChange={handleSelectChange}
                  filterSort={(optionA, optionB) =>
                    (optionA?.label ?? '').toLowerCase().localeCompare((optionB?.label ?? '').toLowerCase())
                  }
                  options={options}
                />
              </Form.Item>
            </Col>
          </Row>
        </Form>

        <Form
          layout="vertical"
          form={form}
          style={{ display: 'flex', justifyContent: 'center', width: '100%' }}
        >
          <Row gutter={20} justify="center" style={{ marginTop: '20px' }}>
            <Col span={7}>
              <Form.Item
                label="Shop Name"
                name="shopName"
                style={{ marginRight: '16px' }}
                rules={[{ required: true, message: 'Please select the shop!' }]}
              >
                <Select
                  showSearch
                  style={{ width: '250px', height: '40px' }}
                  placeholder="Select a shop"
                  filterOption={(input, option) =>
                    (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                  }
                  onChange={handleShopNameChange}
                  options={shopOptions}
                />
              </Form.Item>
            </Col>

            <Col span={7}>
              <Form.Item
                label="Item ID"
                name="itemId"
                style={{ marginRight: '16px' }}
                rules={[{ required: true, message: 'Please select the item ID!' }]}
              >
                <Select
                  showSearch
                  style={{ width: '250px', height: '40px' }}
                  placeholder="Select an item ID"
                  filterOption={(input, option) =>
                    (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                  }
                  onChange={handleItemIdChange}
                  options={itemOptions}
                />
              </Form.Item>
            </Col>

            <Col span={7}>
              <Form.Item
                label="Variation SKU"
                name="variationSku"
                style={{ marginRight: '16px' }}
                rules={[{ required: true, message: 'Please select the variation SKU!' }]}
              >
                <Select
                  showSearch
                  style={{ width: '250px', height: '40px' }}
                  placeholder="Select a variation SKU"
                  filterOption={(input, option) =>
                    (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                  }
                  onChange={handleVariationSkuChange}
                  options={variationOptions}
                />
              </Form.Item>
            </Col>

            <Col span={3}>
              <Button
                type="primary"
                style={{ marginTop: '30px', height: '40px', width: '100px' }}
                onClick={handleAdd}
              >
                Add
              </Button>
            </Col>
          </Row>
        </Form>
        

<Layout.Content
  style={{
    margin: '24px',
    padding: 36,
    minHeight: 360,
    background: '#fff',
    borderRadius: '8px',
  }}
>
  <Table
    columns={columns}
    dataSource={tableData}
    pagination={false}
    locale={{ emptyText: <Empty description="No Data" /> }}
    style={{ marginTop: '20px' }}
  />

  <Row justify="center" style={{ marginTop: '20px' }}>
    <Col>
      <Button
        type="default"
        style={{
          height: '40px',
          width: '100px',
          backgroundColor: '#4DA02F',
          color: 'white',
          border: 'none',
        }}
        onClick={handleSubmit}
      >
        Submit
      </Button>
    </Col>
  </Row>
</Layout.Content>
</Layout.Content>

</ConfigProvider>
);
};

export default Sync;

