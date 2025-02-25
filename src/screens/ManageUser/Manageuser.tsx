import { Card, Typography, Layout, Table, Button, Input, Modal, Form, Select, Row, Col, Popconfirm, notification, Space } from "antd";
import { EditOutlined, DeleteOutlined, PlusOutlined, SearchOutlined, UserOutlined } from "@ant-design/icons";
import Avatar, { genConfig } from 'react-nice-avatar'
import { useState, useEffect } from "react";
import { GET, POST, PATCH, DELETE } from "../../services/index";
import { FETCHUSERS, ADDUSER, EDITUSER, DELETEUSER } from "../../services/path";
import { User, Role, Warehouse, ApiResponse } from "./types";

const { Content } = Layout;
const { Search } = Input;
const { Title } = Typography;

const getUsers = async (isActive = true, limit = 100, offset = 0) => {
    const response = await GET(`manage-users/?isActive=${isActive}&limit=${limit}&offset=${offset}`);
    return response.data;
};

const addUser = async (userData: any) => {
    return POST(ADDUSER, userData);
};

const editUser = async (userID: string, userData: any) => {
    return PATCH(`${EDITUSER}${userID}`, userData);
};

const deleteUser = async (userID: string) => {
    return DELETE(`${DELETEUSER}${userID}`, {});
};

const ManageUser = () => {
    const [data, setData] = useState<User[]>([]);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [form] = Form.useForm();
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [roles, setRoles] = useState<Role[]>([]);
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [loading, setLoading] = useState(false);
    const [searchText, setSearchText] = useState('');
    const [filteredData, setFilteredData] = useState<User[]>([]);

    useEffect(() => {
        loadUsers();
        loadRoles();
        loadWarehouses();
    }, []);

    useEffect(() => {
        const filtered = data.filter(user => 
            user.userName.toLowerCase().includes(searchText.toLowerCase()) ||
            user.roleName.toLowerCase().includes(searchText.toLowerCase()) ||
            user.warehouseName.toLowerCase().includes(searchText.toLowerCase())
        );
        setFilteredData(filtered);
    }, [searchText, data]);

    const loadRoles = async () => {
        try {
            const response = await GET('roles');
            setRoles(response.data);
        } catch (error) {
            console.error("Failed to fetch roles", error);
        }
    };

    const loadWarehouses = async () => {
        try {
            const response = await GET('warehouses');
            setWarehouses(response.data);
        } catch (error) {
            console.error("Failed to fetch warehouses", error);
        }
    };

    const loadUsers = async () => {
        try {
            setLoading(true);
            const response = await getUsers(true, 100, 0);
            const { success, data } = response as ApiResponse<User[]>;
            
            if (success && Array.isArray(data)) {
                setData(data);
                setFilteredData(data);
            } else {
                notification.error({ message: "Failed to load users" });
            }
        } catch (error) {
            console.error("Failed to fetch users", error);
            notification.error({ message: "Failed to load users" });
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async () => {
        try {
            const values = await form.validateFields();
            if (currentUser) {
                await editUser(currentUser.userID, values);
                notification.success({ message: "User updated successfully" });
            } else {
                await addUser(values);
                notification.success({ message: "User added successfully" });
            }
            setIsModalVisible(false);
            form.resetFields();
            loadUsers();
        } catch (error) {
            console.error("Failed to save user", error);
        }
    };

    const handleDelete = async (record: User) => {
        try {
            await deleteUser(record.userID);
            notification.success({ message: "User deleted successfully" });
            loadUsers();
        } catch (error) {
            console.error("Failed to delete user", error);
        }
    };

    const handleOpenModal = (user: User | null) => {
        setCurrentUser(user);
        if (user) {
            form.setFieldsValue(user);
        } else {
            form.resetFields();
        }
        setIsModalVisible(true);
    };

    const handleSearch = (value: string) => {
        setSearchText(value);
    };

    const renderFormItems = (isEditing: boolean) => (
        <>
            <Row gutter={16}>
                <Col span={12}>
                    <Form.Item
                        name="userName"
                        label="ชื่อผู้ใช้"
                        rules={[{ required: true, message: 'กรุณากรอกชื่อผู้ใช้' }]}
                    >
                        <Input disabled={isEditing} />
                    </Form.Item>
                </Col>
                <Col span={12}>
                    <Form.Item
                        name="nickName"
                        label="ชื่อเล่น"
                        rules={[{ required: true, message: 'กรุณากรอกชื่อเล่น' }]}
                    >
                        <Input disabled={isEditing} />
                    </Form.Item>
                </Col>
            </Row>
            <Row gutter={16}>
                <Col span={12}>
                    <Form.Item
                        name="fullNameTH"
                        label="ชื่อ-นามสกุล"
                        rules={[{ required: true, message: 'กรุณากรอกชื่อ-นามสกุล' }]}
                    >
                        <Input disabled={isEditing} />
                    </Form.Item>
                </Col>
                <Col span={12}>
                    <Form.Item
                        name="departmentNo"
                        label="แผนก"
                        rules={[{ required: true, message: 'กรุณากรอกแผนก' }]}
                    >
                        <Input disabled={isEditing} />
                    </Form.Item>
                </Col>
            </Row>
            <Row gutter={16}>
                <Col span={12}>
                    <Form.Item
                        name="roleID"
                        label="บทบาท"
                        rules={[{ required: true, message: 'กรุณาเลือกบทบาท' }]}
                    >
                        <Select>
                            {roles.map((role: any) => (
                                <Select.Option key={role.id} value={role.id}>
                                    {role.name}
                                </Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                </Col>
                <Col span={12}>
                    <Form.Item
                        name="warehouseID"
                        label="คลังสินค้า"
                        rules={[{ required: true, message: 'กรุณาเลือกคลังสินค้า' }]}
                    >
                        <Select>
                            {warehouses.map((warehouse: any) => (
                                <Select.Option key={warehouse.id} value={warehouse.id}>
                                    {warehouse.name}
                                </Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                </Col>
            </Row>
            <Form.Item
                name="description"
                label="คำอธิบาย"
            >
                <Input.TextArea disabled={isEditing} />
            </Form.Item>
        </>
    );

    // Function to generate consistent avatar config for a user
    const getAvatarConfig = (userName: string) => {
        // Use username as seed to generate consistent avatars
        const seed = userName.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
        
        return genConfig({
            sex: seed % 2 ? 'man' : 'woman',
            hairStyle: seed % 2 ? 'normal' : 'thick',
            hatStyle: 'none',
            eyeStyle: seed % 3 ? 'circle' : 'oval',
            glassesStyle: seed % 4 ? 'none' : 'round',
            noseStyle: 'short',
            mouthStyle: 'laugh',
            shirtStyle: seed % 2 ? 'hoody' : 'polo',
            bgColor: `hsl(${seed % 360}, 70%, 90%)`,
        });
    };

    const columns = [
        {
            title: "ผู้ใช้งาน",
            key: "avatar",
            width: 250,
            render: (record: User) => (
                <Space>
                    <Avatar style={{ width: '40px', height: '40px' }} {...getAvatarConfig(record.userName)} />
                    <div>
                        <div style={{ fontWeight: 'bold' }}>{record.userName}</div>
                        <div style={{ color: '#666' }}>{record.nickName}</div>
                    </div>
                </Space>
            ),
        },
        {
            title: "ชื่อ-นามสกุล",
            dataIndex: "fullNameTH",
            key: "fullNameTH",
            width: "200px",
        },
        {
            title: "แผนก",
            dataIndex: "departmentNo",
            key: "departmentNo",
            width: "100px",
        },
        {
            title: "บทบาท",
            dataIndex: "roleName",
            key: "roleName",
            width: "150px",
        },
        {
            title: "คลังสินค้า",
            dataIndex: "warehouseName",
            key: "warehouseName",
            width: "120px",
        },
        {
            title: "การดำเนินการ",
            key: "action",
            width: 200,
            render: (_: any, record: User) => (
                <Space>
                    <Button
                        type="primary"
                        icon={<EditOutlined />}
                        onClick={() => handleOpenModal(record)}
                        ghost
                    >
                        แก้ไข
                    </Button>
                    <Popconfirm
                        title="ยืนยันการลบผู้ใช้"
                        description={`คุณต้องการลบผู้ใช้ ${record.userName} ใช่หรือไม่?`}
                        onConfirm={() => handleDelete(record)}
                        okText="ยืนยัน"
                        cancelText="ยกเลิก"
                        okButtonProps={{ danger: true }}
                    >
                        <Button danger icon={<DeleteOutlined />}>
                            ลบ
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <Layout className="site-layout-background" style={{ padding: 24, background: '#F5F5F5' }}>
            <Card bordered={false}>
                <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
                    <Col>
                        <Title level={4} style={{ margin: 0 }}>
                            <UserOutlined /> จัดการผู้ใช้งาน
                        </Title>
                    </Col>
                    <Col>
                        <Space>
                            <Search
                                placeholder="ค้นหาผู้ใช้งาน..."
                                allowClear
                                onSearch={handleSearch}
                                style={{ width: 300 }}
                            />
                            <Button
                                type="primary"
                                icon={<PlusOutlined />}
                                onClick={() => handleOpenModal(null)}
                            >
                                เพิ่มผู้ใช้งาน
                            </Button>
                        </Space>
                    </Col>
                </Row>

                <Table
                    columns={columns}
                    dataSource={filteredData}
                    rowKey="userID"
                    loading={loading}
                    pagination={{
                        showSizeChanger: true,
                        showTotal: (total) => `ทั้งหมด ${total} รายการ`,
                        pageSize: 10,
                        showQuickJumper: true,
                    }}
                    size="middle"
                    bordered
                />

                <Modal
                    title={
                        <Space align="center">
                            <Avatar 
                                style={{ width: '32px', height: '32px' }} 
                                {...getAvatarConfig(currentUser?.userName || 'new')} 
                            />
                            {currentUser ? "แก้ไขผู้ใช้งาน" : "เพิ่มผู้ใช้งาน"}
                        </Space>
                    }
                    open={isModalVisible}
                    onOk={handleSave}
                    onCancel={() => {
                        setIsModalVisible(false);
                        form.resetFields();
                    }}
                    okText={currentUser ? "บันทึก" : "เพิ่ม"}
                    cancelText="ยกเลิก"
                    confirmLoading={loading}
                    width={600}
                >
                    <Form
                        form={form}
                        layout="vertical"
                        validateTrigger="onBlur"
                    >
                        {renderFormItems(!!currentUser)}
                    </Form>
                </Modal>
            </Card>
        </Layout>
    );
};

export default ManageUser;
