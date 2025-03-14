// src/screens/ManageUser/index.tsx
import React, { useState, useEffect } from "react";
import { Card, Layout, notification } from "antd";
import { useUsers } from "../../redux/users/hooks";
import { UserResponse, AddUserRequest, EditUserRequest } from "../../redux/users/types";

// Components
import UserHeader from "./components/UserHeader";
import UserTable from "./components/UserTable";
import UserForm from "./components/UserForm";

const { Content } = Layout;

export const ManageUser: React.FC = () => {
  // State
  const [searchText, setSearchText] = useState("");
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedUser, setSelectedUser] = useState<UserResponse | null>(null);

  // Redux hooks
  const {
    users,
    roles,
    warehouses,
    loading,
    error,
    pagination,
    fetchUsers,
    fetchRoles,
    fetchWarehouses,
    addUser,
    editUser,
    deleteUser,
  } = useUsers();

  // เมื่อโหลดหน้าแรก
  useEffect(() => {
    // โหลดข้อมูลเริ่มต้น
    fetchUsers({ isActive: true });
    fetchRoles();
    fetchWarehouses();
  }, [fetchUsers, fetchRoles, fetchWarehouses]);

  // ติดตาม error
  useEffect(() => {
    if (error) {
      notification.error({
        message: "เกิดข้อผิดพลาด",
        description: error,
      });
    }
  }, [error]);

  // กรองผู้ใช้ตามคำค้นหา
  const filteredUsers = users.filter(
    (user) =>
      user.userName.toLowerCase().includes(searchText.toLowerCase()) ||
      user.fullNameTH.toLowerCase().includes(searchText.toLowerCase()) ||
      user.roleName.toLowerCase().includes(searchText.toLowerCase()) ||
      user.warehouseName.toLowerCase().includes(searchText.toLowerCase())
  );

  // จัดการการเปลี่ยนหน้าตาราง
  const handleTableChange = (pagination: any) => {
    fetchUsers({
      isActive: true,
      limit: pagination.pageSize,
      offset: (pagination.current - 1) * pagination.pageSize,
    });
  };

  // เปิด modal สำหรับเพิ่มผู้ใช้
  const handleAddUser = () => {
    setSelectedUser(null);
    setIsModalVisible(true);
  };

  // เปิด modal สำหรับแก้ไขผู้ใช้
  const handleEditUser = (user: UserResponse) => {
    setSelectedUser(user);
    setIsModalVisible(true);
  };

  // ปิด modal
  const handleCancelModal = () => {
    setIsModalVisible(false);
    setSelectedUser(null);
  };

  // บันทึกข้อมูลผู้ใช้ (เพิ่มหรือแก้ไข)
const handleSaveUser = (values: AddUserRequest | EditUserRequest) => {
  console.log('handleSaveUser called with values:', values);
  console.log('selectedUser:', selectedUser);
  
  if (selectedUser) {
    // ถ้ามี selectedUser แสดงว่าอยู่ในโหมดแก้ไข
    console.log('Edit mode detected, calling editUser');
    editUser(selectedUser.userID, values as EditUserRequest);
  } else {
    // ถ้าไม่มี selectedUser แสดงว่าอยู่ในโหมดเพิ่ม
    console.log('Add mode detected, calling addUser');
    addUser(values as AddUserRequest);
  }

  setIsModalVisible(false);
  setSelectedUser(null);
};

  // ลบผู้ใช้
  const handleDeleteUser = (userID: string) => {
    deleteUser(userID);
  };

  // จัดการค้นหา
  const handleSearch = (value: string) => {
    setSearchText(value);
  };

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
        {/* ส่วนหัวและค้นหา */}
        <UserHeader
          onSearch={handleSearch}
          onAddUser={handleAddUser}
          searchPlaceholder="ค้นหาชื่อผู้ใช้, ชื่อ-สกุล, บทบาท..."
        />

        {/* ตารางผู้ใช้ */}
        <UserTable
          users={filteredUsers}
          loading={loading}
          pagination={pagination}
          onEdit={handleEditUser}
          onDelete={handleDeleteUser}
          onChange={handleTableChange}
        />
      </Card>

      {/* Modal สำหรับเพิ่ม/แก้ไขผู้ใช้ */}
      <UserForm
        visible={isModalVisible}
        user={selectedUser}
        roles={roles}
        warehouses={warehouses}
        onSave={handleSaveUser}
        onCancel={handleCancelModal}
        confirmLoading={loading}
      />
    </Layout>
  );
};

export default ManageUser;