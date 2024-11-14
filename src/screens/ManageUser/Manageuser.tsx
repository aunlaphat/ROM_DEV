import { ConfigProvider, Layout, Table, Button, Space, Avatar, Modal, Form, Input, Select, Row, Popconfirm, notification } from "antd";
import { EditOutlined, DeleteOutlined, PlusOutlined } from "@ant-design/icons";
import { useState } from "react";
import image2 from '../../assets/images/image (2).png';
import image3 from '../../assets/images/image (3).png';

interface User {
    Name: string;
    Role: string;
    Warehouse: string;
    Image: string;
}

const initialData: User[] = [
    {
        Name: "Piyawalee",
        Role: "Admin",
        Warehouse: "BES",
        Image: image2,
    },
    {
        Name: "Narawit",
        Role: "User",
        Warehouse: "BES",
        Image: image3,
    },
];

const ManageUser = () => {
    const [data, setData] = useState<User[]>(initialData);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [form] = Form.useForm();
    const [currentUser, setCurrentUser] = useState<User | null>(null);

    const handleAddNew = () => {
        setCurrentUser(null);
        setIsModalVisible(true);
    };

    const handleCancel = () => {
        setIsModalVisible(false);
        form.resetFields();
    };

    const handleSave = () => {
        form.validateFields().then(values => {
            const newUser: User = {
                ...values,
                Image: image2,
            };

            if (currentUser) {
                const updatedUser: User = { ...currentUser, ...values };
                setData(prevData => prevData.map(item => (item.Name === currentUser.Name ? updatedUser : item)));
                notification.success({
                    message: 'Update Successful',
                    description: `User ${updatedUser.Name} has been updated.`,
                });
            } else {
                setData(prevData => [...prevData, newUser]);
                notification.success({
                    message: 'User Added',
                    description: `User ${newUser.Name} has been added.`,
                });
            }

            setIsModalVisible(false);
            form.resetFields();
        }).catch(info => {
            console.log('Validation Failed:', info);
        });
    };

    const handleEdit = (record: User) => {
        setCurrentUser(record);
        form.setFieldsValue(record);
        setIsModalVisible(true);
    };

    const handleDelete = (record: User) => {
        setData(data.filter(item => item.Name !== record.Name));
        notification.success({
            message: 'Delete Successful',
            description: `User ${record.Name} has been deleted.`,
        });
    };

    const columns = [
        {
            title: 'User Name',
            key: 'userDetails',
            width: '250px',
            render: (record: User) => (
                <div style={{ display: 'flex', alignItems: 'center' }}>
                    <Avatar src={record.Image} alt="avatar" style={{ marginRight: '8px', width: 56, height: 56 }} />
                    <div>
                        <div style={{ color: '#35465B', fontWeight: 'bold', fontSize: '16px', }}>{record.Name}</div>
                        <div style={{ color: '#35465B', marginTop: '10px', fontSize: '14px' }}>
                            {record.Role} <span style={{ marginLeft: '8px' }}>{record.Warehouse}</span>
                        </div>

                    </div>
                </div>
            ),
        },
        {
            title: 'Action',
            key: 'action',
            width: '250px',
            render: (text: any, record: User) => (
                <>
                    <Button
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                        type="primary"
                        style={{ color: '#FFFFFF', background: '#D9D9D9' }}
                    >
                        แก้ไข
                    </Button>
                    <Popconfirm
                        title={`Are you sure to delete ${record.Name}?`}
                        onConfirm={() => handleDelete(record)}
                        okText="Yes"
                        cancelText="No"
                    >
                        <Button
                            icon={<DeleteOutlined />}
                            type="primary"
                            style={{ color: 'red', background: '#F9D3D3', marginLeft: '20px' }}
                        >
                            ลบ
                        </Button>
                    </Popconfirm>
                </>
            ),
        },
    ];

    return (
        <ConfigProvider>
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                Manage User
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
                    <Row gutter={20} justify="end">
                        <Button
                        id=" Add New"
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={handleAddNew}
                            style={{ display: "flex", justifyContent: "flex-end", margin: "16px 28px", background: '#72BBFF' }}
                        >
                            Add New
                        </Button>
                    </Row>
                    <Table
                    id="Table1"
                        style={{ padding: 40 }}
                        components={{
                            header: {
                                cell: (props: React.HTMLAttributes<HTMLElement>) => (
                                    <th {...props} style={{ backgroundColor: '#E9F3FE', color: '#35465B' }} />
                                ),
                            },
                        }}
                        columns={columns}
                        dataSource={data}
                        rowKey="Name"
                        pagination={false}
                    />
                </Layout.Content>
            </Layout>
            <Modal
                width={600}
                title={currentUser ? "Edit User" : "Add New"}
                visible={isModalVisible}
                onCancel={handleCancel}
                footer={[
                    <div style={{ display: 'flex', justifyContent: 'center', width: '100%' }}>
                        <Popconfirm
                            title={currentUser ? "Are you sure you want to update?" : "Are you sure you want to save?"}
                            onConfirm={handleSave}
                            okText="Yes"
                            cancelText="No"
                        >
                            <Button id="Update,บันทัก" key="submit" type="primary" style={{ background: '#14C11B' }}>
                                {currentUser ? "Update" : "บันทึก"}
                            </Button>
                        </Popconfirm>
                    </div>
                ]}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                    id="Name"
                        name="Name"
                        label="Username"
                        rules={[{ required: true, message: 'Please enter the username' }]}
                    >
                        <Input style={{ height: 40 }} placeholder="Enter username" />
                    </Form.Item>
                    <Form.Item
                    id="Role"
                        name="Role"
                        label="Role"
                        rules={[{ required: true, message: 'Please select a role' }]}
                    >
                        <Select placeholder="Select role" style={{ height: 40 }}>
                            <Select.Option value="Admin">Admin</Select.Option>
                            <Select.Option value="User">User</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                     id="Warehouse"
                        name="Warehouse"
                        label="Warehouse"
                        rules={[{ required: true, message: 'Please select a warehouse' }]}
                    >
                        <Select style={{ height: 40 }} placeholder="Select warehouse">
                            <Select.Option value="RBN">RBN</Select.Option>
                            <Select.Option value="Fair">Fair</Select.Option>
                            <Select.Option value="Online">Online</Select.Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>
        </ConfigProvider>
    );
};

export default ManageUser;
