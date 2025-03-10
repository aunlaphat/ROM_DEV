import React from "react";
import { AvatarGenerator } from "../../avatar/AvatarGenerator";
import {
  Dropdown,
  MenuProps,
  Modal,
  Space,
  Typography,
  Tooltip,
  Divider,
} from "antd";
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  DownOutlined,
} from "@ant-design/icons";
import {
  HeaderBarStyle,
  TopBarDropDown,
  TopBarUser,
} from "../../Layouts/headerLayout/style";
import { Icon } from "../../../resources/icon";
import { useAuth } from "../../../hooks/useAuth";
import { logger } from "../../../utils/logger";

const HeaderBar = ({ collapsed, toggle }: any) => {
  // ใช้ข้อมูลทั้งหมดจาก useAuth Context แทนที่จะใช้ useSelector โดยตรง
  const { logout, userID, userName, fullNameTH, roleName } = useAuth();

  const handleLogout = () => {
    Modal.confirm({
      title: "Confirm Logout",
      content: "คุณต้องการออกจากระบบใช่หรือไม่? 🤔",
      okText: "ออกจากระบบ",
      cancelText: "ยกเลิก",
      onOk: async () => {
        try {
          logger.log("info", "User logging out", { userID, userName });
          await logout();
        } catch (error) {
          console.error("Logout Failed", error);
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

      <TopBarUser
        style={{ display: "flex", alignItems: "center", gap: "12px" }}
      >
        <Tooltip title="Profile">
          <AvatarGenerator userName={userName} userID={userID} size="large" />
        </Tooltip>

        <Space direction="horizontal" size="small">
          <Typography.Text strong>{userID || "N/A"}</Typography.Text>
          <Typography.Text>{fullNameTH || "N/A"}</Typography.Text>
          <Divider type="vertical" />
          <Typography.Text keyboard type="success">
            {roleName || "N/A"}
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