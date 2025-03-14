// src/screens/ManageUser/components/UserHeader.tsx
import React from "react";
import { Row, Col, Typography, Button, Input, Space } from "antd";
import { PlusOutlined, UserOutlined, SearchOutlined } from "@ant-design/icons";

const { Title } = Typography;
const { Search } = Input;

interface UserHeaderProps {
  onSearch: (value: string) => void;
  onAddUser: () => void;
  searchPlaceholder?: string;
}

/**
 * Component แสดงส่วนหัวของหน้าจัดการผู้ใช้
 */
const UserHeader: React.FC<UserHeaderProps> = ({
  onSearch,
  onAddUser,
  searchPlaceholder = "ค้นหาผู้ใช้งาน...",
}) => {
  return (
    <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
      <Col>
        <Title level={4} style={{ margin: 0 }}>
          <UserOutlined /> จัดการผู้ใช้งาน
        </Title>
      </Col>
      <Col>
        <Space>
          <Search
            placeholder={searchPlaceholder}
            allowClear
            onSearch={onSearch}
            style={{ width: 300 }}
            prefix={<SearchOutlined />}
          />
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={onAddUser}
          >
            เพิ่มผู้ใช้งาน
          </Button>
        </Space>
      </Col>
    </Row>
  );
};

export default UserHeader;