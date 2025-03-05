import React from 'react';
import { Steps, Badge, Tooltip, Space, Typography } from 'antd';
import { SearchOutlined, FormOutlined, NumberOutlined, CheckCircleOutlined, EyeOutlined } from "@ant-design/icons";
import { ReturnOrderState } from '../../../../redux/orders/api';

const { Text } = Typography;

interface ReturnOrderStepsProps {
  currentStep: 'search' | 'create' | 'sr' | 'confirm' | 'preview';
  orderData: ReturnOrderState['orderData'];
  getStepStatus: (stepKey: string) => 'process' | 'finish' | 'wait'; // เพิ่ม getStepStatus
}

const ReturnOrderSteps: React.FC<ReturnOrderStepsProps> = ({ currentStep, orderData, getStepStatus }) => {
  const renderStepTitle = (title: string, subTitle?: string) => (
    <Space direction="vertical" size={0}>
      <Text strong>{title}</Text>
      {subTitle && (
        <Text type="secondary" style={{ fontSize: '12px' }}>
          {subTitle}
        </Text>
      )}
    </Space>
  );

  const steps = [
    {
      key: 'search',
      title: renderStepTitle('Search Order', 'ค้นหา SO/Order'),
      icon: <SearchOutlined />,
      description: (
        <Badge
          status={currentStep === 'search' ? 'processing' : 'success'}
          text="ค้นหา Order ที่ต้องการคืนสินค้า"
        />
      ),
      disabled: currentStep === 'confirm' && !!orderData?.head.srNo
    },
    {
      key: 'create',
      title: renderStepTitle('Create Return Order', 'สร้างคำสั่งคืนสินค้า'),
      icon: <FormOutlined />,
      description: (
        <Badge
          status={getStepStatus('create') === 'process' ? 'processing' : 
                 getStepStatus('create') === 'finish' ? 'success' : 'default'}
          text="กรอกข้อมูลและเลือกสินค้าที่ต้องการคืน"
        />
      ),
      disabled: currentStep === 'confirm' && !!orderData?.head.srNo
    },
    {
      key: 'sr',
      title: renderStepTitle('Generate SR', 'สร้างเลข SR'),
      icon: <NumberOutlined />,
      description: (
        <Tooltip title={orderData?.head.srNo ? `SR Number: ${orderData.head.srNo}` : 'รอการสร้าง SR'}>
          <Badge
            status={orderData?.head.srNo ? 'success' : 
                   getStepStatus('sr') === 'process' ? 'processing' : 'default'}
            text={
              orderData?.head.srNo ? 
              <Text type="success">SR: {orderData.head.srNo}</Text> : 
              'รอการสร้างเลข SR'
            }
          />
        </Tooltip>
      )
    },
    {
      key: 'preview',
      title: renderStepTitle('Preview', 'ตรวจสอบข้อมูล'),
      icon: <EyeOutlined />,
      description: (
        <Badge
          status={getStepStatus('preview') === 'process' ? 'processing' : 
                 getStepStatus('preview') === 'finish' ? 'success' : 'default'}
          text="ตรวจสอบข้อมูลก่อนยืนยัน"
        />
      )
    },
    {
      key: 'confirm',
      title: renderStepTitle('Confirm', 'ยืนยันคำสั่งคืนสินค้า'),
      icon: <CheckCircleOutlined />,
      description: (
        <Badge
          status={getStepStatus('confirm') === 'process' ? 'processing' : 'default'}
          text="ยืนยันและเสร็จสิ้น"
        />
      )
    }
  ];

  return (
    <div style={{ padding: '24px 0' }}>
      <Steps
        type="navigation"
        current={steps.findIndex(s => s.key === currentStep)}
        items={steps}
        style={{ 
          padding: '24px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          borderRadius: '8px',
          background: '#fff'
        }}
        onChange={(current) => {
          // Handle step navigation if needed
          console.log('Step changed:', current);
        }}
      />
    </div>
  );
};

export default ReturnOrderSteps;
