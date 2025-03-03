import {
  Card,
  Typography,
  Layout,
  Table,
  Button,
  Input,
  Modal,
  Form,
  Select,
  Row,
  Col,
  Popconfirm,
  notification,
  Space,
} from "antd";
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  SearchOutlined,
  UserOutlined,
} from "@ant-design/icons";
import { useState, useEffect } from "react";
import { GET, POST, PATCH, DELETE } from "../../services/index";
import { FETCHUSERS, ADDUSER, EDITUSER, DELETEUSER } from "../../services/path";
import { User, Role, Warehouse, ApiResponse } from "./types";
import { AvatarGenerator } from "../../components/avatar/AvatarGenerator";

const { Content } = Layout;
const { Search } = Input;
const { Title } = Typography;

export const ManageUser = () => {
  const [data, setData] = useState<User[]>([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [roles, setRoles] = useState<Role[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchText, setSearchText] = useState("");
  const [filteredData, setFilteredData] = useState<User[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 100
  });

  useEffect(() => {
    loadUsers();
    loadRoles();
    loadWarehouses();
  }, []);

  useEffect(() => {
    const filtered = data.filter(
      (user) =>
        user.userName.toLowerCase().includes(searchText.toLowerCase()) ||
        user.roleName.toLowerCase().includes(searchText.toLowerCase()) ||
        user.warehouseName.toLowerCase().includes(searchText.toLowerCase())
    );
    setFilteredData(filtered);
  }, [searchText, data]);

  const loadRoles = async () => {
    try {
      const response = await GET("roles");
      setRoles(response.data);
    } catch (error) {
      console.error("Failed to fetch roles", error);
    }
  };

  const loadWarehouses = async () => {
    try {
      const response = await GET("warehouses");
      setWarehouses(response.data);
    } catch (error) {
      console.error("Failed to fetch warehouses", error);
    }
  };

  const loadUsers = async (page = pagination.current, pageSize = pagination.pageSize) => {
    try {
      setLoading(true);
      const offset = (page - 1) * pageSize;
      
      const response = await GET(`manage-users/?isActive=true&limit=${pageSize}&offset=${offset}`);
      const apiResponse = response.data as ApiResponse<User[]>;

      if (apiResponse.success && Array.isArray(apiResponse.data)) {
        setData(apiResponse.data);
        setFilteredData(apiResponse.data);
        setPagination({
          ...pagination,
          current: page,
          pageSize: pageSize,
          //total: apiResponse.total || apiResponse.data.length // ต้องมีการส่ง total จาก API
          total: 100
        });
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
        await PATCH(`${EDITUSER}${currentUser.userID}`, values);
        notification.success({ message: "User updated successfully" });
      } else {
        await POST(ADDUSER, values);
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
      await DELETE(`${DELETEUSER}${record.userID}`, {});
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

  const handleTableChange = (pagination: any) => {
    loadUsers(pagination.current, pagination.pageSize);
  };

  const columns = [
    {
      title: "ผู้ใช้งาน",
      key: "avatar",
      width: 250,
      render: (record: User) => (
        <Space>
          <AvatarGenerator userID={record.userName} userName={record.userName} size="large" />{" "}
          <div>
            <div style={{ fontWeight: "bold" }}>{record.userName}</div>
            <div style={{ color: "#666" }}>{record.nickName}</div>
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
    <Layout
      className="site-layout-background"
      style={{
        margin: "24px",
        padding: 24,
        minHeight: 360,
        background: "#f5f5f5",
        borderRadius: "8px",
      }}
    >
      <Card bordered={false}>
        <Row
          justify="space-between"
          align="middle"
          style={{ marginBottom: 24 }}
        >
          <Col>
            <Title level={4} style={{ margin: 0 }}>
              <UserOutlined /> Mange Users
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
            ...pagination,
            showSizeChanger: true,
            showTotal: (total) => `ทั้งหมด ${total} รายการ`,
            showQuickJumper: true,
            pageSizeOptions: ['10', '20', '50', '100']
          }}
          onChange={handleTableChange}
          size="middle"
          bordered
        />
      </Card>
    </Layout>
  );
};
