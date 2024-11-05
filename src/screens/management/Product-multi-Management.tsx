import React, { useState } from "react";
import { ConfigProvider, InputNumber, Layout, Row, Col, Select, Form, Button, Tag, Popconfirm, notification } from "antd";

const options = [
  { value: '1', label: 'Not Identified' },
  { value: '2', label: 'G097171-ARM01-BL Bewell Sport armband size M For' },
  { value: '3', label: 'G097171-ARM01-GY Bewell Sport armband size M For' },
  { value: '4', label: 'G097171-ARM02-BL Bewell Sport armband size L' },
  { value: '5', label: 'G097171-ARM03-GR Bewell Sport armband size M with light' },
  { value: '6', label: 'Cancelled' },
];

const Management = () => {
  const [form] = Form.useForm();
  const [selectedLabels, setSelectedLabels] = useState<string[]>([]);
  const [selectedValue, setSelectedValue] = useState<string | undefined>();

  const handleSelectChange = (value: string) => {
    const selectedOption = options.find(option => option.value === value);
    if (selectedOption && !selectedLabels.includes(selectedOption.label)) {
      setSelectedLabels([...selectedLabels, selectedOption.label]);
      setSelectedValue(undefined); // Clear the selection after adding
    }
  };

  const handleRemove = (label: string) => {
    setSelectedLabels(selectedLabels.filter(item => item !== label));
  };

  const handleSubmit = (values: any) => {
    const data = {
      SKU: selectedLabels,
      ...values,
    };
    console.log(data);
    
    // Show success notification
    notification.success({
      message: 'Success',
      description: 'Data has been successfully submitted.',
    });

    // Reset the form and state if needed
    form.resetFields();
    setSelectedLabels([]);
    setSelectedValue(undefined);
  };

  const handleConfirm = () => {
    form.validateFields()
      .then(values => handleSubmit(values))
      .catch(errorInfo => {
        console.log('Form validation failed:', errorInfo);
      });
  };

  const handleCancel = () => {
    // Do nothing if canceled
  };

  return (
    <ConfigProvider>
      <div style={{ marginLeft: '28px', fontSize: '25px', fontWeight: 'bold', color: 'DodgerBlue' }}>
        Product multi-management
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
                label="Select SKU"
                name="selectedSku"
                rules={[{ required: true, message: 'Please select the SKU!' }]}
              >
                <Select
                  showSearch
                  style={{ width: 550, height: 50  }}
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
          <Row style={{ marginTop: '20px' }}>
            <Col span={11} style={{ display: 'flex', justifyContent: 'center' }}>
              <Form.Item
                label="Minimum Stock"
                name="minimumStockValue"
                style={{ flex: 1, marginRight: '16px' }}
                rules={[{ required: true, message: 'Please enter the minimum stock!' }]}
              >
                <InputNumber
                  min={1}
                  max={70}
                  style={{ width: '250px', height: '40px', lineHeight: '40px' }}
                  disabled={selectedLabels.length === 0}
                />
              </Form.Item>
            </Col>

            <Col span={11} style={{ display: 'flex', justifyContent: 'center' }}>
              <Form.Item
                label="Status"
                name="status"
                style={{ flex: 1, marginRight: '16px' }}
                rules={[{ required: true, message: 'Please select the status!' }]}
              >
                <Select
                  showSearch
                  style={{ width: '250px', height: '40px' }}
                  options={[
                    { value: 'active', label: <span style={{ color: 'green' }}>Active</span> },
                    { value: 'inactive', label: <span style={{ color: 'red' }}>Inactive</span> },
                    { value: 'disabled', label: 'Disabled', disabled: true },
                  ]}
                />
              </Form.Item>
            </Col>

            <Col span={2}>
              <Popconfirm
                title="Are you sure you want to submit the form?"
                onConfirm={handleConfirm}
                onCancel={handleCancel}
                okText="Yes"
                cancelText="No"
              >
                <Button type="primary" style={{ marginTop: '30px', height: '40px', width: '100px' }}>
                  Submit
                </Button>
              </Popconfirm>
            </Col>
          </Row>
        </Form>
      </Layout.Content>

      {/* Display the selected tags below */}
      <Layout.Content
        style={{
          margin: '24px',
          padding: 36,
          minHeight: 360,
          background: '#fff',
          borderRadius: '8px',
        }}
      >
        <span style={{ display: 'inline-block' }}>
          {selectedLabels.map(label => (
            <Tag
              key={label}
              closable
              onClose={() => handleRemove(label)}
              style={{
                cursor: 'pointer',
                padding: '4px 12px',
                whiteSpace: 'nowrap',
                display: 'inline-block',
                marginBottom: '20px'
              }}
            >
              {label}
            </Tag>
          ))}
        </span>
      </Layout.Content>
    </ConfigProvider>
  );
};

export default Management;
