import { Button, Col, ConfigProvider, Layout, Row } from "antd";
import { SizeType } from "antd/es/config-provider/SizeContext";
import { useState } from "react";
import { useNavigate } from 'react-router-dom';

const CreateReturn = () => {
    const navigate = useNavigate();
    const [size, setSize] = useState<SizeType>('large'); // default is 'middle'

    const handleNavigateToSR = () => {
        navigate('/SR'); // เส้นทางนี้ควรตรงกับการตั้งค่า Route ใน App.js หรือไฟล์ routing ของคุณ
    };
    const handleNavigateToIJ = () => {
        navigate('/IJ'); // เส้นทางนี้ควรตรงกับการตั้งค่า Route ใน App.js หรือไฟล์ routing ของคุณ
    };

    return (
        <ConfigProvider>
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                Create Return
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
                    <Row style={{ marginTop: 50 }}>
                        <Col span={12} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                            <Button
                                style={{
                                    width: '283px',
                                    height: '179px',
                                    background: '#1EA39A',
                                    color: '#fff',
                                    fontSize: '30px',
                                }}
                                onClick={handleNavigateToSR}
                            >
                                SR
                            </Button>
                        </Col>
                        <Col span={12} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                            <Button
                                style={{
                                    width: '283px',
                                    height: '179px',
                                    background: '#2386E1',
                                    color: '#fff',
                                    fontSize: '30px',
                                }}
                                onClick={handleNavigateToIJ}
                            >
                                IJ
                            </Button>
                        </Col>
                    </Row>
                </Layout.Content>
            </Layout>
        </ConfigProvider>
    );
};

export default CreateReturn;
