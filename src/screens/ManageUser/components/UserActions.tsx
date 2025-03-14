// src/screens/ManageUser/components/UserActions.tsx
import React from "react";
import { Button, Popconfirm, Space, Tooltip } from "antd";
import { EditOutlined, DeleteOutlined } from "@ant-design/icons";
import { UserResponse } from "../../../redux/users/types";

interface UserActionsProps {
  user: UserResponse;
  onEdit: () => void;
  onDelete: () => void;
}

/**
 * Component แสดงปุ่มการดำเนินการสำหรับแต่ละรายการผู้ใช้
 */
const UserActions: React.FC<UserActionsProps> = ({ user, onEdit, onDelete }) => {
  return (
    <Space size="middle">
      <Tooltip title={`แก้ไขผู้ใช้: ${user.userName}`}>
        <Button
          type="primary"
          icon={<EditOutlined />}
          onClick={onEdit}
          ghost
        >
          แก้ไข
        </Button>
      </Tooltip>
      
      <Popconfirm
        title="ยืนยันการลบผู้ใช้"
        description={`คุณต้องการลบผู้ใช้ ${user.userName} ใช่หรือไม่?`}
        onConfirm={onDelete}
        okText="ยืนยัน"
        cancelText="ยกเลิก"
        okButtonProps={{ danger: true }}
        placement="topRight"
      >
        <Button 
          danger 
          icon={<DeleteOutlined />}
        >
          ลบ
        </Button>
      </Popconfirm>
    </Space>
  );
};

export default UserActions;