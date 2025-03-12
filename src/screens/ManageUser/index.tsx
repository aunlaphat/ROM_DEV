// src/screens/ManageUser/index.tsx
import { useState, useEffect } from "react";
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
import { useUsers } from "../../redux/users/hooks";
import { AvatarGenerator } from "../../components/avatar/AvatarGenerator";

const { Content } = Layout;
const { Search } = Input;
const { Title } = Typography;
const { Option } = Select;

export const ManageUser = () => {
  // State
  const [searchText, setSearchText] = useState("");
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<any>(null);
  const [form] = Form.useForm();

  // Redux hooks
  const {
    users,
    roles,
    warehouses,
    loading,
    pagination,
    fetchUsers,
    fetchRoles,
    fetchWarehouses,
    addUser,
    editUser,
    deleteUser,
  } = useUsers();

  // Effects
  useEffect(() => {
    // Load initial data
    fetchUsers({ isActive: true });
    fetchRoles();
    fetchWarehouses();
  }, [fetchUsers, fetchRoles, fetchWarehouses]);

  // Filtered users based on search
  const filteredUsers = users.filter(
    (user) =>
      user.userName.toLowerCase().includes(searchText.toLowerCase()) ||
      user.fullNameTH.toLowerCase().includes(searchText.toLowerCase()) ||
      user.roleName.toLowerCase().includes(searchText.toLowerCase()) ||
      user.warehouseName.toLowerCase().includes(searchText.toLowerCase())
  );

  // Table pagination change handler
  const handleTableChange = (pagination: any) => {
    fetchUsers({
      isActive: true,
      limit: pagination.pageSize,
      offset: (pagination.current - 1) * pagination.pageSize,
    });
  };

  // Modal handlers
  const showModal = (user: any = null) => {
    setEditingUser(user);
    if (user) {
      // Edit mode - pre-fill form ใช้ name แทน id
      const selectedRole = roles.find((role) => role.roleID === user.roleID);
      const selectedWarehouse = warehouses.find(
        (warehouse) => warehouse.warehouseID === user.warehouseID
      );

      form.setFieldsValue({
        userID: user.userID,
        roleID: user.roleID,
        roleName: selectedRole?.roleName || "", // เพิ่ม field roleName
        warehouseID: user.warehouseID,
        warehouseName: selectedWarehouse?.warehouseName || "", // เพิ่ม field warehouseName
      });
    } else {
      // Add mode - reset form
      form.resetFields();
    }
    setIsModalVisible(true);
  };

  const handleCancel = () => {
    setIsModalVisible(false);
    setEditingUser(null);
    form.resetFields();
  };

  const handleSave = async () => {
    try {
      const values = await form.validateFields();

      const selectedRole = roles.find((role) => role.roleID === values.roleID);
      const selectedWarehouse = warehouses.find(
        (warehouse) => warehouse.warehouseID === values.warehouseID
      );

      if (editingUser) {
        // Edit mode
        editUser(editingUser.userID, {
          userID: values.userID,
          roleID: values.roleID,
          roleName: selectedRole?.roleName, // ส่ง name ด้วย (ถ้าต้องการ)
          warehouseID: Number(values.warehouseID),
          warehouseName: selectedWarehouse?.warehouseName, // ส่ง name ด้วย (ถ้าต้องการ)
        });
      } else {
        // Add mode
        addUser({
          userID: values.userID,
          roleID: values.roleID,
          roleName: selectedRole?.roleName, // ส่ง name ด้วย (ถ้าต้องการ)
          warehouseID: Number(values.warehouseID),
          warehouseName: selectedWarehouse?.warehouseName, // ส่ง name ด้วย (ถ้าต้องการ)
        });
      }

      setIsModalVisible(false);
      setEditingUser(null);
      form.resetFields();
    } catch (error) {
      console.error("Form validation failed:", error);
    }
  };

  // Delete handler
  const handleDelete = (userID: string) => {
    deleteUser(userID);
  };

  // Table columns
  const columns = [
    {
      title: "ผู้ใช้งาน",
      key: "avatar",
      width: 250,
      render: (record: any) => (
        <Space>
          <AvatarGenerator
            userID={record.userID}
            userName={record.userName}
            size="large"
          />
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
      render: (_: any, record: any) => (
        <Space>
          <Button
            type="primary"
            icon={<EditOutlined />}
            onClick={() => showModal(record)}
            ghost
          >
            แก้ไข
          </Button>
          <Popconfirm
            title="ยืนยันการลบผู้ใช้"
            description={`คุณต้องการลบผู้ใช้ ${record.userName} ใช่หรือไม่?`}
            onConfirm={() => handleDelete(record.userID)}
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
              <UserOutlined /> จัดการผู้ใช้งาน
            </Title>
          </Col>
          <Col>
            <Space>
              <Search
                placeholder="ค้นหาผู้ใช้งาน..."
                allowClear
                onSearch={(value) => setSearchText(value)}
                style={{ width: 300 }}
              />
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => showModal()}
              >
                เพิ่มผู้ใช้งาน
              </Button>
            </Space>
          </Col>
        </Row>

        <Table
          columns={columns}
          dataSource={filteredUsers}
          rowKey="userID"
          loading={loading}
          pagination={{
            ...pagination,
            showSizeChanger: true,
            showTotal: (total) => `ทั้งหมด ${total} รายการ`,
            showQuickJumper: true,
            pageSizeOptions: ["10", "20", "50", "100"],
          }}
          onChange={handleTableChange}
          size="middle"
          bordered
        />
      </Card>

      {/* Modal for Add/Edit User */}
      <Modal
        title={editingUser ? "แก้ไขผู้ใช้งาน" : "เพิ่มผู้ใช้งานใหม่"}
        open={isModalVisible}
        onOk={handleSave}
        onCancel={handleCancel}
        okText={editingUser ? "บันทึกการแก้ไข" : "เพิ่มผู้ใช้"}
        cancelText="ยกเลิก"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="userID"
            label="รหัสผู้ใช้"
            rules={[{ required: true, message: "กรุณาระบุรหัสผู้ใช้" }]}
          >
            <Input placeholder="กรุณาระบุรหัสผู้ใช้" disabled={!!editingUser} />
          </Form.Item>

          <Form.Item
            name="roleID"
            label="บทบาท"
            rules={[{ required: true, message: "กรุณาเลือกบทบาท" }]}
          >
            <Select placeholder="เลือกบทบาท">
              {roles.map((role) => (
                <Option key={role.roleID} value={role.roleID}>
                  {role.roleID} : {role.roleName}
                </Option>
              ))}
            </Select>
          </Form.Item>

          {/* เพิ่ม hidden field สำหรับเก็บ roleName */}
          <Form.Item name="roleName" hidden={true}>
            <Input />
          </Form.Item>

          <Form.Item
            name="warehouseID"
            label="คลังสินค้า"
            rules={[{ required: true, message: "กรุณาเลือกคลังสินค้า" }]}
          >
            <Select placeholder="เลือกคลังสินค้า">
              {warehouses.map((warehouse) => (
                <Option
                  key={warehouse.warehouseID}
                  value={warehouse.warehouseID}
                >
                  {warehouse.warehouseID} : {warehouse.warehouseName}
                </Option>
              ))}
            </Select>
          </Form.Item>

          {/* เพิ่ม hidden field สำหรับเก็บ warehouseName */}
          <Form.Item name="warehouseName" hidden={true}>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  );
};
