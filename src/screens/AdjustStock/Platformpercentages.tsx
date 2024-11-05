import React, { useState } from 'react';
import { Space, ConfigProvider, Row, Col, Card, Button, Select, Form, Popconfirm, Layout, notification } from 'antd';
import Bewell1 from '../../assets/images/Bewell1.png';
import Bewell2 from '../../assets/images/Bewell2.png';
import ShopeeLogo from '../../assets/images/shopee logo.png';
import LazadaLogo from '../../assets/images/lazada logo.png';
import TiktokLogo from '../../assets/images/tiktok logo.png';
import NocNocLogo from '../../assets/images/nocnoc logo.png';
import { Doughnut } from 'react-chartjs-2';
import { Chart as ChartJS, ArcElement, Tooltip, Legend, ChartOptions, TooltipItem } from 'chart.js';
import ChartDataLabels from 'chartjs-plugin-datalabels';

ChartJS.register(ArcElement, Tooltip, Legend, ChartDataLabels);

const logoMap: { [key: string]: string } = {
  'Shopee': ShopeeLogo,
  'Lazada': LazadaLogo,
  'Tiktok': TiktokLogo,
  'NocNoc': NocNocLogo,
};

const bewellLogoMap: { [key: string]: string } = {
  'Bewell1': Bewell1,
  'Bewell2': Bewell2,
};

const initialShopData = [
  {
    logo: 'Shopee',
    data: [
      { logo: 'Bewell1', shopname: 'Bewell Official Store', data: 20, color: '#FF6384' },
      { logo: 'Bewell2', shopname: 'Bewell Shop', data: 50, color: '#FFCE56' },
    ],
  },
  {
    logo: 'Lazada',
    data: [
      { logo: 'Bewell2', shopname: 'Bewell Official', data: 20, color: '#9747FF' },
      { logo: 'Bewell2', shopname: 'Bewell Shop', data: 50, color: '#36A2EB' },
    ],
  },
  {
    logo: 'Tiktok',
    data: [
      { logo: 'Bewell1', shopname: 'Bewellstyle', data: 20, color: '#358C5B' },
    ],
  },
  {
    logo: 'NocNoc',
    data: [
      { logo: 'Bewell1', shopname: 'Bewell Official Store', data: 20, color: '#ff8f00' },
    ],
  },
];

const Doughnutdata = (shopData: typeof initialShopData) => {
  const total = shopData.flatMap(group => group.data.map(item => item.data)).reduce((a, b) => a + b, 0);
  const remaining = 300 - total;

  return {
    labels: shopData.flatMap(group => group.data.map(item => item.shopname)).concat('Remaining'),
    datasets: [{
      label: 'Votes',
      data: shopData.flatMap(group => group.data.map(item => item.data)).concat(remaining),
      backgroundColor: shopData.flatMap(group => group.data.map(item => item.color)).concat('#D3D3D3'),
      hoverBackgroundColor: shopData.flatMap(group => group.data.map(item => item.color)).concat('#A9A9A9'),
    }],
  };
}; 

const options: ChartOptions<'doughnut'> = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    tooltip: {
      callbacks: {
        label: function (tooltipItem: TooltipItem<'doughnut'>) {
          let label = tooltipItem.label || '';
          if (label) {
            label += ': ';
          }
          if (tooltipItem.raw !== undefined) {
            label += `${tooltipItem.raw}%`;
          }
          return label;
        },
      },
    },
    datalabels: {
      color: '#fff',
      display: true,
      formatter: (value: any) => value,
      font: {
        weight: 'bold',
        size: 16,
      },
      anchor: 'center',
      align: 'center',
    },
  },
};

const Platform = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [form] = Form.useForm();
  const [visible, setVisible] = useState(false);
  const [shop, setShop] = useState(initialShopData);

  const toggleCollapsed = () => {
    setCollapsed(prev => !prev);
  };

  const handleConfirm = () => {
    form
      .validateFields()
      .then((values) => {
        const totalValue = shop.reduce((total, shopGroup, index) => {
          const groupTotal = shopGroup.data.reduce((groupSum, _, idx) => groupSum + (values[index]?.data[idx] || 0), 0);
          return total + groupTotal;
        }, 0);
  
        if (totalValue > 300) {
          notification.error({
            message: 'Error',
            description: 'Total value exceeds 300%',
          });
        } else {
          const updatedShop = shop.map((shopGroup, index) => ({
            ...shopGroup,
            data: shopGroup.data.map((value, idx) => ({
              ...value,
              data: values[index]?.data[idx] || value.data,
            })),
          }));
  
          console.log('Updated Shop Data:', updatedShop);
  
          updatedShop.forEach(shopGroup => {
            console.log(`Shop Logo: ${shopGroup.logo}`);
            shopGroup.data.forEach(item => {
              console.log(`Shop Name: ${item.shopname}, Value: ${item.data}`);
            });
          });
  
          setShop(updatedShop);
          setCollapsed(false);
          setVisible(false);

          notification.success({
            message: 'Success',
            description: 'Data has been saved successfully!',
          });
        }
      })
      .catch((error) => {
        console.log('Validation Failed:', error);
      });
  };

  const handleCancel = () => {
    setVisible(false);
  };

  const onFinish = () => {
    setVisible(true);
  };

  return (
    <ConfigProvider>
      <div style={{ marginLeft: '28px', fontSize: '25px', fontWeight: 'bold', color: 'DodgerBlue' }}>
        Platform percentages
      </div>
      <Layout.Content
        style={{
          margin: '24px 24px',
          padding: 36,
          minHeight: 360,
          background: '#fff',
          borderRadius: '8px',
        }}
      >
        <Space
          style={{
            borderColor: '#F5F5F5',
            backgroundColor: '#F5F5F5',
            height: '100%',
            width: '100%',
            display: 'flex',
            justifyContent: 'center',
            borderRadius: '0.5%',
            marginTop: '24px',
          }}
        >
          <Doughnut data={Doughnutdata(shop)} options={options} style={{ height: '700px', width: '500px' }} />
        </Space>

        <div style={{ marginLeft: '16px', fontSize: '25px', fontWeight: 'bold', padding: 'inherit' }}>
          Platform Bewell Store
        </div>

        <Row gutter={16}>
          {shop.map((shopGroup, index) => (
            <Col key={index} span={12}>
              <Card
                style={{
                  boxShadow: "0px 4px 10px rgba(0, 0, 0, 0.1)",
                  borderRadius: '10px',
                  backgroundColor: '#FFFFFF',
                  height: '250px',
                  flexDirection: 'column',
                  marginBottom: '16px',
                }}
                title={
                  <div style={{ textAlign: 'center' }}>
                    <img src={logoMap[shopGroup.logo]} alt={`${shopGroup.logo} Logo`} style={{ width: '150px', height: 'auto', margin: '20px' }} />
                  </div>
                }
                bordered={false}
              >
                <div style={{ flex: 1 }}>
                  {shopGroup.data.map((value, idx) => (
                    <Row key={idx} gutter={16} style={{ alignItems: 'center', marginBottom: '12px' }}>
                      <Col span={6} style={{ textAlign: 'center' }}>
                        {value.color && (
                          <Button
                            type="primary"
                            style={{
                              backgroundColor: value.color,
                              borderColor: value.color,
                              borderRadius: '50%',
                              width: '30px',
                              height: '30px',
                            }}
                          />
                        )}
                      </Col>
                      <Col span={12} style={{ textAlign: 'center', fontSize: '18px' }}>
                        {value.shopname}
                      </Col>
                      <Col span={6} style={{ textAlign: 'center' }}>
                        {value.data}%
                      </Col>
                    </Row>
                  ))}
                </div>
              </Card>
            </Col>
          ))}
        </Row>

        <div style={{ display: 'flex', justifyContent: 'center', width: '100%', marginBottom: 16, padding: '40px' }}>
          <Button
            type="primary"
            onClick={toggleCollapsed}
            style={{ marginBottom: 16, padding: '10px 24px', fontSize: '20px', width: '150px', height: '60px' }}
            >
              แก้ไข
            </Button>
          </div>
  
          {collapsed && (
            <Form
              form={form}
              layout="vertical"
              onFinish={onFinish}
            >
              <Row gutter={20}>
                {shop.map((shopGroup, index) => (
                  <Col key={index} span={12}>
                    <Card
                      style={{
                        boxShadow: "0px 4px 10px rgba(0, 0, 0, 0.1)",
                        borderRadius: '10px',
                        backgroundColor: '#F5F5F5',
                        height: '400px',
                        flexDirection: 'column',
                        marginBottom: '16px',
                      }}
                      title={
                        <div style={{ textAlign: 'center' }}>
                          <img 
                            src={logoMap[shopGroup.logo]} 
                            alt={`${shopGroup.logo} Logo`} 
                            style={{ width: '150px', height: 'auto', margin: '20px' }} 
                          />
                        </div>
                      }
                      bordered={false}
                    >
                      <Form.Item name={[index, 'data']} noStyle>
                        {shopGroup.data.map((value, idx) => (
                          <Row key={idx} gutter={16} style={{ alignItems: 'center', marginBottom: '12px' }}>
                            <Col span={24}>
                              <Row style={{ backgroundColor: '#FFFFFF', borderRadius: '5px', padding: '10px', alignItems: 'center', height: '90px' }}>
                                <Col span={5} style={{ textAlign: 'center' }}>
                                  <img 
                                    src={bewellLogoMap[value.logo]} 
                                    alt={`${value.logo} Logo`} 
                                    style={{ width: '50px', height: '50px' }} 
                                  />
                                </Col>
                                <Col span={12} style={{ textAlign: 'center', fontSize: '18px' }}>
                                  {value.shopname}
                                </Col>
                                <Col span={6} style={{ textAlign: 'center' }}>
                                  <Form.Item name={[index, 'data', idx]} initialValue={value.data} noStyle>
                                    <Select style={{ width: '100%' }}>
                                      {Array.from({ length: 10 }, (_, i) => (i + 1) * 10).map(optionValue => (
                                        <Select.Option key={optionValue} value={optionValue}>
                                          {optionValue}%
                                        </Select.Option>
                                      ))}
                                    </Select>
                                  </Form.Item>
                                </Col>
                              </Row>
                            </Col>
                          </Row>
                        ))}
                      </Form.Item>
                    </Card>
                  </Col>
                ))}
              </Row>
  
              <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '30vh' }}>
                <Popconfirm
                  title="แน่ใจว่าต้องการบันทึกข้อมูล?"
                  visible={visible}
                  onConfirm={handleConfirm}
                  onCancel={handleCancel}
                  okText="ใช่"
                  cancelText="ยกเลิก"
                >
                  <Button
                    type="primary"
                    onClick={() => setVisible(true)}
                    style={{ fontSize: '20px', width: '150px', height: '60px' }}
                  >
                    บันทึก
                  </Button>
                </Popconfirm>
              </div>
            </Form>
          )}
        </Layout.Content>
      </ConfigProvider>
    );
  };
  
  export default Platform;