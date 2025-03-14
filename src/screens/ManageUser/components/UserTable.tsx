// src/screens/ManageUser/components/UserTable.tsx
import React from "react";
import { Table, Space } from "antd";
import { AvatarGenerator } from "../../../components/avatar/AvatarGenerator";
import UserActions from "./UserActions";
import { UserResponse } from "../../../redux/users/types";

interface UserTableProps {
  users: UserResponse[];
  loading: boolean;
  pagination: {
    current: number;
    pageSize: number;
    total: number;
  };
  onEdit: (user: UserResponse) => void;
  onDelete: (userID: string) => void;
  onChange: (pagination: any) => void;
}

/**
 * Component แสดงตารางรายการผู้ใช้
 */
const UserTable: React.FC<UserTableProps> = ({
  users,
  loading,
  pagination,
  onEdit,
  onDelete,
  onChange,
}) => {
  const columns = [
    {
      title: "ผู้ใช้งาน",
      key: "avatar",
      width: 250,
      render: (record: UserResponse) => (
        <div style={{ display: "flex", alignItems: "center", gap: "12px" }}>
          <AvatarGenerator
            userID={record.userID}
            userName={record.userName}
            size="large"
          />
          <div>
            <div style={{ fontWeight: "bold" }}>{record.userName}</div>
            <div style={{ color: "#666" }}>{record.nickName}</div>
          </div>
        </div>
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
      render: (_: any, record: UserResponse) => (
        <UserActions
          user={record}
          onEdit={() => onEdit(record)}
          onDelete={() => onDelete(record.userID)}
        />
      ),
    },
  ];

  return (
    <Table
      columns={columns}
      dataSource={users}
      rowKey="userID"
      loading={loading}
      pagination={{
        ...pagination,
        showSizeChanger: true,
        showTotal: (total) => `ทั้งหมด ${total} รายการ`,
        showQuickJumper: true,
        pageSizeOptions: ["10", "20", "50", "100"],
      }}
      onChange={onChange}
      size="middle"
      bordered
    />
  );
};

export default UserTable;