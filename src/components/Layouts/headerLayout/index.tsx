import React from "react";
import {
  Dropdown,
  MenuProps,
  Modal,
  Divider,
  Space,
  Typography,
  Tooltip,
} from "antd";
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  DownOutlined,
} from "@ant-design/icons";
import { HeaderBarStyle, TopBarDropDown, TopBarUser } from "../../layouts/headerLayout/style";
import { Icon } from "../../../resources";
import { useAuth } from "../../../hooks/useAuth"; // เปลี่ยนการ import
import { useSelector } from "react-redux";
import { TextSmall } from "../../text";
import { logger } from "../../../utils/logger";

interface HeaderBarProps {
  collapsed: boolean;
  toggle: () => void;
}

const HeaderBar: React.FC<HeaderBarProps> = ({ collapsed, toggle }) => {
  // ใช้ hook ที่เรา refactor แล้ว
  const { 
    logout, 
    userID, 
    userName, 
    fullName, 
    roleName 
  } = useAuth();

  const handleLogout = () => {
    Modal.confirm({
      title: "ยืนยันการออกจากระบบ",
      content: "คุณต้องการออกจากระบบใช่หรือไม่? 🤔",
      okText: "ออกจากระบบ",
      cancelText: "ยกเลิก",
      onOk: async () => {
        try {
          logger.log("info", "User logging out", { userID, userName });
          await logout();
        } catch (error) {
          logger.error("Logout Failed", error);
        }
      },
    });
  };

  const menuItems: MenuProps["items"] = [
    {
      key: "logout",
      label: "ออกจากระบบ",
      icon: Icon.logout(),
      onClick: handleLogout,
    },
  ];

  return (
    <HeaderBarStyle
      className="site-layout-background"
      style={{ padding: 0, backgroundColor: "white" }}
    >
      {React.createElement(collapsed ? MenuUnfoldOutlined : MenuFoldOutlined, {
        className: "trigger",
        onClick: toggle,
        style: { margin: "0 20px", color: "black" },
      })}

      <TopBarUser>
        <Tooltip title="โปรไฟล์">
          <AvatarGenerator 
            userName={userName || 'User'} 
            userID={userID || '0'} 
            size="large" 
          />
        </Tooltip>

        <Space direction="horizontal" size="small">
          <Typography.Text strong>{userID || 'N/A'}</Typography.Text>
          <Typography.Text>{fullName || 'N/A'}</Typography.Text>
          <Divider type="vertical" />
          <Typography.Text keyboard type="success">
            {roleName || 'N/A'}
          </Typography.Text>
        </Space>
      </TopBarUser>

      <TopBarDropDown>
        <Dropdown menu={{ items: menuItems }} trigger={["click"]}>
          <Space style={{ cursor: "pointer", padding: "0 8px" }}>
            <DownOutlined />
          </Space>
        </Dropdown>
      </TopBarDropDown>
    </HeaderBarStyle>
  );
};

export default HeaderBar;