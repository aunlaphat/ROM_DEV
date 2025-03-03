import React from "react";
import { AvatarGenerator } from "../../avatar/AvatarGenerator";
import {
  Button,
  Dropdown,
  MenuProps,
  Tag,
  Modal,
  notification,
  Card,
  Divider,
  Space,
  Typography,
  Tooltip,
  Avatar,
} from "antd";
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  DownOutlined,
} from "@ant-design/icons";
import { HeaderBarStyle, TopBarDropDown, TopBarUser } from "./style";
import { Icon } from "../../../resources";
import { useAuthLogin } from "../../../hooks/useAuth";
import { useSelector } from "react-redux";
import { TextSmall } from "../../text";
import { logger } from "../../../utils/logger";

const HeaderBar = ({ collapsed, toggle }: any) => {
  const { onLogout } = useAuthLogin();
  const user = useSelector((state: any) => state.auth.user);

  const userId = user?.userID || "N/A";
  const userName = user?.userName || "N/A";
  const userFullName = user?.fullName || "N/A";
  const roleName = user?.roleName || "N/A";

  const handleLogout = () => {
    Modal.confirm({
      title: "Confirm Logout",
      content: "คุณต้องการออกจากระบบใช่หรือไม่? 🤔",
      okText: "ออกจากระบบ",
      cancelText: "ยกเลิก",
      onOk: async () => {
        try {
          logger.auth("info", "User logging out", { userId, userName });

          await onLogout();
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
          <AvatarGenerator userName={userName} userID={userId} size="large" />
        </Tooltip>

        <Space direction="horizontal" size="small">
          <Typography.Text strong>{userId}</Typography.Text>
          <Typography.Text>{userFullName}</Typography.Text>
          <Divider type="vertical" />
          <Typography.Text keyboard type="success">
            {roleName}
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
