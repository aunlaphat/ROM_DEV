import React, { useState, useRef } from 'react';
import { Steps, Col, ConfigProvider, Layout, Row, Button, Divider, Popconfirm, notification } from 'antd';
import Webcam from 'react-webcam';
import { RedoOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom'; // Import useNavigate

const SaleReturn: React.FC = () => {
  const navigate = useNavigate(); // Get navigate function
  const [currentStep, setCurrentStep] = useState(0);
  const [images, setImages] = useState<{ [key: string]: string | null }>({
    step1: null,
    step2: null,
    step3: null,
  });
  const webcamRef = useRef<Webcam>(null);

  const onChange = (current: number) => {
    setCurrentStep(current);
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
    if (currentStep < 2) {
      setCurrentStep((prevStep) => prevStep + 1);
    }
  };

  const handleSubmit = () => {
    // Reset all steps and images
    setCurrentStep(0);
    setImages({
      step1: null,
      step2: null,
      step3: null,
    });

    // Show notification on submit success
    notification.success({
      message: 'ส่งข้อมูลสำเร็จ',
      description: 'ข้อมูลของคุณถูกส่งเรียบร้อยแล้ว!',
    });

    // Navigate to CreateBlindReturn page
    navigate('/CreateBlindReturn'); // Change this path according to your routing
  };

  return (
    <ConfigProvider>
      <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
        Create Blind Return
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
          <Steps current={currentStep} onChange={onChange} style={{ marginTop: '20px', width: '100%' }}>
            <Steps.Step />
            <Steps.Step />
            <Steps.Step />
          </Steps>

          <Row justify="center" align="middle" style={{ marginTop: 20 }}>
            <Col span={24} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', marginTop: '20px' }}>
              {!images[`step${currentStep + 1}`] ? (
                <Webcam
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
              {!images[`step${currentStep + 1}`] ? (
                <Button
                  type="primary"
                  onClick={capturePhoto}
                  style={{ display: 'flex', alignItems: 'center', marginTop: '10px' }}>
                  ถ่ายรูป
                </Button>
              ) : (
                <>
                  <Button
                    icon={<RedoOutlined />}
                    style={{ marginRight: '20px', width: '100px', color: '#35465B' }}
                    type="default"
                    onClick={retakePhoto}
                  >
                    ถ่ายใหม่
                  </Button>

                  {currentStep < 2 && (
                    <Button
                      style={{ width: '100px' }}
                      type="primary"
                      onClick={handleNextStep}
                    >
                      ถัดไป
                    </Button>
                  )}
                </>
              )}
            </Col>
          </Row>

          {/* Display gallery and Submit button only after finishing step 3 */}
          {currentStep === 2 && images.step3 && (
            <>
              <Divider orientation="left" style={{ color: '#657589' }}>รูปทั้งหมด</Divider>
              <Row gutter={[16, 16]} justify="center">
                {Object.keys(images).map((key, index) => (
                  images[key] && (
                    <Col key={index} span={6}>
                      <img
                        src={images[key]!}
                        alt={`Step ${index + 1}`}
                        style={{ width: '100%', border: '1px solid #ccc', borderRadius: '8px' }}
                      />
                      <div style={{ textAlign: 'center', marginTop: '8px' }}>
                        รูปถ่ายขั้นตอนที่ {index + 1}
                      </div>
                    </Col>
                  )
                ))}
              </Row>

              <Row justify="center" style={{ marginTop: '20px' }}>
                <Popconfirm
                  title="คุณแน่ใจหรือไม่ว่าต้องการส่งข้อมูล?"
                  onConfirm={handleSubmit}
                  okText="ใช่"
                  cancelText="ไม่"
                >
                  <Button type="primary">Submit</Button>
                </Popconfirm>
              </Row>
            </>
          )}
        </Layout.Content>
      </Layout>
    </ConfigProvider>
  );
};

export default SaleReturn;
