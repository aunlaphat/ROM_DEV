import { ConfigProvider, Layout, Table, Button, Space, Avatar, Modal, Form, Input, Select, Row, Popconfirm, notification } from "antd";
import { EditOutlined, DeleteOutlined, PlusOutlined } from "@ant-design/icons";
import { useState, useEffect } from "react";
import { GET, POST, PATCH, DELETE } from "../../services/index";
import { FETCHUSERS, FETCHUSERBYID, ADDUSER, EDITUSER, DELETEUSER } from "../../services/path";

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

interface User {
    userID: string;
    userName: string;
    roleName: string;
    warehouseName: string;
}

const ManageUser = () => {
    const [data, setData] = useState<User[]>([]);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [form] = Form.useForm();
    const [currentUser, setCurrentUser] = useState<User | null>(null);

    useEffect(() => {
        loadUsers();
    }, []);

    const loadUsers = async () => {
        try {
            const rawData = await getUsers(true, 100, 0);
            
            if (!Array.isArray(rawData)) {
                console.error("Invalid response format:", rawData);
                return;
            }
    
            setData(rawData);
        } catch (error) {
            console.error("Failed to fetch users", error);
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

    return (
        <ConfigProvider>
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                Manage User
            </div>
            <Layout>
                <Layout.Content style={{ margin: "24px", padding: 36, minHeight: 360, background: "#fff", borderRadius: "8px", overflow: "auto" }}>
                    <Row gutter={20} justify="end">
                        <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsModalVisible(true)} style={{ display: "flex", justifyContent: "flex-end", margin: "16px 28px", background: "#72BBFF" }}>
                            Add New
                        </Button>
                    </Row>
                    <Table columns={[
                        {
                            title: "UserID",
                            dataIndex: "userID",
                            key: "userID",
                            width: "250px",
                        },
                        {
                            title: "Username",
                            dataIndex: "userName",
                            key: "userName",
                            width: "250px",
                        },
                        {
                            title: "Role",
                            dataIndex: "roleName",
                            key: "roleName",
                        },
                        {
                            title: "Warehouse",
                            dataIndex: "warehouseName",
                            key: "warehouseName",
                        },
                        {
                            title: "Action",
                            key: "action",
                            width: "250px",
                            render: (text: any, record: User) => (
                                <>
                                    <Button icon={<EditOutlined />} onClick={() => { setCurrentUser(record); setIsModalVisible(true); }} type="primary" style={{ color: "#FFFFFF", background: "#D9D9D9" }}>
                                        แก้ไข
                                    </Button>
                                    <Popconfirm title={`Are you sure to delete ${record.userName}?`} onConfirm={() => handleDelete(record)} okText="Yes" cancelText="No">
                                        <Button icon={<DeleteOutlined />} type="primary" style={{ color: "red", background: "#F9D3D3", marginLeft: "20px" }}>
                                            ลบ
                                        </Button>
                                    </Popconfirm>
                                </>
                            ),
                        },
                    ]} dataSource={data} rowKey="userID" pagination={false} />
                </Layout.Content>
            </Layout>
        </ConfigProvider>
    );
};

export default ManageUser;
